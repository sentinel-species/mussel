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
		tree, err := pypi.Scrape(pkg)
		if err != nil {
			_, err := fmt.Fprintf(os.Stderr, "Error scraping package: %v\n", err)
			if err != nil {
				return
			}
			os.Exit(1)
		}
		fmt.Printf("Successfully processed package: %s@%s\n", tree.Name, tree.Version)
		if len(tree.Dependencies) > 0 {
			fmt.Println("Dependencies:")
			for _, dep := range tree.Dependencies {
				fmt.Printf("  - %s %s\n", dep.Name, dep.Version)
			}
		}
	default:
		fmt.Printf("Unsupported ecosystem: %s\n", ecosystem)
		os.Exit(1)
	}
}
