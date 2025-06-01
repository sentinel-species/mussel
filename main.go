package main

import (
	"fmt"
	"log"
	"mussel/cmd"
	"mussel/internal/config"
	"os"
)

func main() {
	err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	err = cmd.CheckRequirements()
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}

	cfg := cmd.ParseFlags()
	cmd.Process(cfg.Ecosystem, cfg.Package)
}
