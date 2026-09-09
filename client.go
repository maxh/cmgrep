package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/maxh/cmgrep/proto"
)

const hostFormat = "fa26-cs425-72%02d.cs.illinois.edu"

const numNodes = 10

type result struct {
	node  int
	count int64
	err   error
}

func queryOneNode(ctx context.Context, node int, host string, req *pb.GrepCountRequest) result {
	conn, err := grpc.NewClient(host+listenAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return result{node: node, err: err}
	}
	defer conn.Close()

	client := pb.NewCmgrepClient(conn)
	resp, err := client.GrepCount(ctx, req)
	if err != nil {
		return result{node: node, err: err}
	}
	return result{node: node, count: resp.GetCount()}
}

func queryAllNodes(ctx context.Context, req *pb.GrepCountRequest, nodeCount int) []result {
	results := make([]result, nodeCount)
	var wg sync.WaitGroup

	for i := range results {
		node := i + 1
		wg.Add(1)
		go func() {
			defer wg.Done()
			results[i] = queryOneNode(ctx, node, fmt.Sprintf(hostFormat, node), req)
		}()
	}
	wg.Wait()

	return results
}

func runClient(args []string) error {
	fs := flag.NewFlagSet("query", flag.ExitOnError)
	ignoreCase := fs.Bool("i", false, "ignore case")
	extended := fs.Bool("E", false, "extended regexp")
	nodeCount := fs.Int("nodes", numNodes, "number of nodes to query")
	fs.Parse(args)

	if *nodeCount < 1 || *nodeCount > numNodes {
		return fmt.Errorf("nodes must be between 1 and %d", numNodes)
	}

	if fs.NArg() != 1 {
		return fmt.Errorf("usage: cmgrep [-i] [-E] [--nodes=N] <pattern>")
	}

	req := &pb.GrepCountRequest{
		Pattern:    fs.Arg(0),
		IgnoreCase: *ignoreCase,
		Extended:   *extended,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	start := time.Now()
	results := queryAllNodes(ctx, req, *nodeCount)
	elapsed := time.Since(start)

	var total int64
	failed := 0
	for _, r := range results {
		if r.err != nil {
			fmt.Fprintf(os.Stderr, "machine.%d: %v\n", r.node, r.err)
			failed++
			continue
		}
		fmt.Printf("machine.%d.log:%d\n", r.node, r.count)
		total += r.count
	}
	fmt.Printf("total:%d\n", total)
	fmt.Fprintf(
		os.Stderr,
		"latency_ms:%.3f, %d machines failed\n",
		float64(elapsed)/float64(time.Millisecond),
		failed,
	)

	return nil
}
