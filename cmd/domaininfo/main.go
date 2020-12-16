package main

import (
	"log"
	"os"

	"github.com/marc-barry/domaininfo/pkg/cmd/domaininfo"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Requires at least one command line argument")
	}

	if err := domaininfo.RunCommand(os.Args[1]); err != nil {
		log.Fatal(err)
	}
}
