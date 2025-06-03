package cmd

import (
	"fmt"
	"os"

	"mussel/internal/pypi"
)

func Process(ecosystem string, pkg string) {
	switch ecosystem {
	case "npm":
		fmt.Println("Using npm ecosystem")
	case "pypi":
		_, err := pypi.Scrape(pkg)
		if err != nil {
			_, err := fmt.Fprintf(os.Stderr, "Error scraping package: %v\n", err)
			if err != nil {
				return
			}
			os.Exit(1)
		}
	default:
		fmt.Printf("Unsupported ecosystem: %s\n", ecosystem)
		os.Exit(1)
	}
}
