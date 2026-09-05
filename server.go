package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/maxh/cmgrep/proto"
)

const listenAddr = ":26950"
const logFileFormat = "machine.%d.log"

func logPathFor(node int) string {
	return fmt.Sprintf(logFileFormat, node)
}

type server struct {
	pb.UnimplementedCmgrepServer

	node int
}

func newServer(node int) *server {
	return &server{node: node}
}

func (s *server) GrepCount(ctx context.Context, req *pb.GrepCountRequest) (*pb.GrepCountResponse, error) {
	f, err := os.Open(logPathFor(s.node))
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

func serve(node int) error {
	logPath := logPathFor(node)

	if _, err := os.Stat(logPath); err != nil {
		return fmt.Errorf("log file for node %d: %w", node, err)
	}

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return fmt.Errorf("tcp port listen %w", err)
	}

	s := grpc.NewServer()
	pb.RegisterCmgrepServer(s, newServer(node))

	log.Printf("node %d serving on %s for file %s", node, listenAddr, logPath)
	return s.Serve(lis)
}
