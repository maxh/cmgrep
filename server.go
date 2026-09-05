package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/maxh/cmgrep/proto"
)

const listenAddr = ":26950"

// our VMs are fa26-cs425-7201 through fa26-cs425-7210
var hostPattern = regexp.MustCompile(`cs425-72(\d\d)`)

func nodeFromHostname() (int, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return 0, err
	}
	m := hostPattern.FindStringSubmatch(hostname)
	if m == nil {
		return 0, fmt.Errorf("no node number in hostname %q", hostname)
	}
	return strconv.Atoi(m[1])
}

type server struct {
	pb.UnimplementedCmgrepServer

	node int
}

func (s *server) GrepCount(ctx context.Context, req *pb.GrepCountRequest) (*pb.GrepCountResponse, error) {
	f, err := os.Open(fmt.Sprintf("machine.%d.log", s.node))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not open log file %v", err)
	}
	defer f.Close()

	n, err := countMatches(f, req.GetPattern(), req.GetIgnoreCase(), req.GetExtended())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "grep failed %v", err)
	}

	return &pb.GrepCountResponse{Count: n}, nil
}

func runServer(args []string) error {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	node := fs.Int("node", 0, "this machine's number")
	fs.Parse(args)

	if *node < 1 {
		n, err := nodeFromHostname()
		if err != nil {
			return fmt.Errorf("pass --node: %w", err)
		}
		*node = n
	}
	return serve(*node)
}

func serve(node int) error {
	logPath := fmt.Sprintf("machine.%d.log", node)

	if _, err := os.Stat(logPath); err != nil {
		return fmt.Errorf("log file for node %d: %w", node, err)
	}

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return fmt.Errorf("tcp port listen %w", err)
	}

	s := grpc.NewServer()
	pb.RegisterCmgrepServer(s, &server{node: node})

	log.Printf("node %d serving on %s for file %s", node, listenAddr, logPath)
	return s.Serve(lis)
}
