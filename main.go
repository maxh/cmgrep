package main

import (
	"flag"
	"log"
)

func main() {
	node := flag.Int("node", 0, "this machine's number, e.g. 1 for machine.1.log")
	flag.Parse()

	if *node < 1 {
		log.Fatalf("--node is required and must be >= 1")
	}

	err := serve(*node)
	if err != nil {
		log.Fatalf("serve error %v", err)
	}
}
