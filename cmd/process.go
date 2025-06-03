package cmd

import (
	"fmt"
	"mussel/internal/config"
	db "mussel/internal/database"
	"mussel/internal/store/database"
	"os"

	"mussel/internal/pypi"
)

func Process(ecosystem string, pkg string) {
	conn, err := db.Connect(config.Config)

	dependencyStore := database.NewDependencyStore(conn)
	packageStore := database.NewPackageStore(conn)
	ecosystemStore := database.NewEcosystemStore(conn)

	if err != nil {
		switch ecosystem {
		case "npm":
			// scraper := npm.NewScraper(dependencyStore, packageStore, ecosystemStore)
			fmt.Println("Using npm ecosystem")
		case "pypi":
			scraper := pypi.NewScraper(dependencyStore, packageStore, ecosystemStore)
			err := scraper.Scrape(pkg)
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
}
