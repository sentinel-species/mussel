package cmd

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Ecosystem string
	Package   string
}

func ParseFlags() *Config {
	config := &Config{}

	flag.StringVar(&config.Ecosystem, "ecosystem", "", "The package ecosystem (e.g., npm, pypi, go)")
	flag.StringVar(&config.Package, "package", "", "The package name to analyze")

	flag.Usage = func() {
		_, err := fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]\n", os.Args[0])
		if err != nil {
			return
		}
		_, err = fmt.Fprintln(flag.CommandLine.Output(), "Options:")
		if err != nil {
			return
		}
		flag.PrintDefaults()
	}

	flag.Parse()

	if config.Ecosystem == "" || config.Package == "" {
		_, err := fmt.Fprintln(os.Stderr, "Error: both ecosystem and package must be specified")
		if err != nil {
			return nil
		}
		flag.Usage()
		os.Exit(1)
	}

	return config
}
