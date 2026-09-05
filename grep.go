package main

import (
	"io"
	"os/exec"
	"strconv"
	"strings"
)

func countMatches(r io.Reader, pattern string, ignoreCase, extended bool) (int64, error) {
	args := []string{"-c"}
	if ignoreCase {
		args = append(args, "-i")
	}
	if extended {
		args = append(args, "-E")
	}
	args = append(args, "--", pattern)

	cmd := exec.Command("grep", args...)
	cmd.Stdin = r
	stdout, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); !ok || exitErr.ExitCode() != 1 {
			return 0, err
		}
	}

	count, err := strconv.ParseInt(strings.TrimSpace(string(stdout)), 10, 64)
	if err != nil {
		return 0, err
	}
	return count, nil
}
