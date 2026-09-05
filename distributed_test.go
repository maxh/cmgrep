package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	pb "github.com/maxh/cmgrep/proto"
)

var queries = []struct {
	name       string
	pattern    string
	ignoreCase bool
	extended   bool
}{
	{"frequent, all nodes", "/frequent", false, false},
	{"somewhat frequent, all nodes", "/happensrarely", false, false},
	{"rare, some nodes", "/alsoquiterare", false, false},
	{"rare, one node", "/morerare", false, false},
	{"absent everywhere", "/nothing_matches_this", false, false},
	{"matches every line", "HTTP/1.1", false, false},
	{"ignore case", "/MORERARE", true, false},
	{"extended regexp", "/(morerare|alsoquiterare)", false, true},
	{"regex metacharacter", "10\\.0\\.3\\.", false, true},
}

func generateLocalLogs(t *testing.T) string {
	script, err := filepath.Abs("scripts/gen_logs.py")
	if err != nil {
		t.Fatal(err)
	}

	dir := t.TempDir()
	cmd := exec.Command("python3", script)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("generating logs: %v\n%s", err, out)
	}
	return dir
}

func localCount(t *testing.T, dir string, node int, pattern string, ignoreCase, extended bool) int64 {
	f, err := os.Open(filepath.Join(dir, fmt.Sprintf("machine.%d.log", node)))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	n, err := countMatches(f, pattern, ignoreCase, extended)
	if err != nil {
		t.Fatal(err)
	}
	return n
}

func TestDistributedGrep(t *testing.T) {
	if testing.Short() {
		t.Skip("needs a server running on every VM")
	}

	// local logs are from fixed seed so match exactly
	dir := generateLocalLogs(t)

	for _, q := range queries {
		t.Run(q.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			results := queryAllNodes(ctx, &pb.GrepCountRequest{
				Pattern:    q.pattern,
				IgnoreCase: q.ignoreCase,
				Extended:   q.extended,
			})

			for i, r := range results {
				node := i + 1
				if r.err != nil {
					t.Errorf("machine.%d: %v", node, r.err)
					continue
				}
				want := localCount(t, dir, node, q.pattern, q.ignoreCase, q.extended)
				if r.count != want {
					t.Errorf("machine.%d: got %d, want %d", node, r.count, want)
				}
			}
		})
	}
}

func TestAbsentPatternIsNotAnError(t *testing.T) {
	if testing.Short() {
		t.Skip("needs a server running on every VM")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for i, r := range queryAllNodes(ctx, &pb.GrepCountRequest{Pattern: "zzz_no_such_pattern_zzz"}) {
		if r.err != nil {
			t.Errorf("machine.%d: %v", i+1, r.err)
		}
		if r.count != 0 {
			t.Errorf("machine.%d: got %d, want 0", i+1, r.count)
		}
	}
}
