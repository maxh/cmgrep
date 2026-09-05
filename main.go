package main

import (
	"log"
	"os"
)

func main() {
	var err error
	if len(os.Args) > 1 && os.Args[1] == "serve" {
		err = runServer(os.Args[2:])
	} else {
		err = runClient(os.Args[1:])
	}
	if err != nil {
		log.Fatalf("%v", err)
	}
}
