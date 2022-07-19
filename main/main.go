package main

import (
	"log"

	"github.com/sub-usd-net/setup/cmd"
)

func main() {
	run, err := cmd.New()
	if err != nil {
		log.Fatalf("Failed to initialize command: %s", err)
	}

	if err := run.Execute(); err != nil {
		log.Fatalf("Error running command: %s", err)
	}
}
