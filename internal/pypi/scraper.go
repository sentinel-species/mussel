package pypi

import (
	"context"
	"fmt"
	"io"
	"mussel/internal/config"
	"mussel/internal/store/database"
	"os"
	"os/exec"
	"strings"

	"mussel/internal/types"
)

type Scraper struct {
	dependencyStore *database.DependencyStore
	packageStore    *database.PackageStore
	ecosystemStore  *database.EcosystemStore
}

func NewScraper(dependencyStore *database.DependencyStore, packageStore *database.PackageStore, ecosystemStore *database.EcosystemStore) *Scraper {
	return &Scraper{
		dependencyStore: dependencyStore,
		packageStore:    packageStore,
		ecosystemStore:  ecosystemStore,
	}
}

func (s *Scraper) Scrape(pkg string) (err error) {
	err = config.Config.Pypi.Setup()
	if err != nil {
		return fmt.Errorf("failed to setup environment: %w", err)
	}

	err = config.Config.Pypi.InstallPython()
	if err != nil {
		return fmt.Errorf("failed to install python: %w", err)
	}

	err = config.Config.Pypi.NewVenv()
	if err != nil {
		return fmt.Errorf("failed to create virtual environment: %w", err)
	}

	packageVersions, err := s.getPackageVersions(pkg)
	if err != nil {
		return fmt.Errorf("failed to get package versions: %w", err)
	}

	for _, version := range strings.Split(*packageVersions, "\n") {
		err := s.dependencyTree(pkg, version)
		if err != nil {
			_, err = fmt.Fprintf(os.Stderr, "failed to process version %s: %v\n", version, err)
			if err != nil {
				return err
			}
		}

	}
	return
}

func (s *Scraper) getPackageVersions(pkg string) (packageVersions *string, err error) {
	curl := exec.Command("curl", "-s", fmt.Sprintf("https://pypi.org/pypi/%s/json", pkg))
	jq := exec.Command("jq", "-r", ".releases | keys[]")
	sort := exec.Command("sort", "-V", "-r")

	jq.Stdin, _ = curl.StdoutPipe()
	sort.Stdin, _ = jq.StdoutPipe()
	output, _ := sort.StdoutPipe()

	_ = curl.Start()
	_ = jq.Start()
	_ = sort.Start()

	bytes, err := io.ReadAll(output)
	if err != nil {
		return nil, fmt.Errorf("failed to read command output: %w", err)
	}
	packageVersions = types.Pointer(string(bytes))

	_ = curl.Wait()
	_ = jq.Wait()
	_ = sort.Wait()
	return
}

func (s *Scraper) dependencyTree(pkg string, version string) (err error) {
	fmt.Printf("Processing %s version %s\n", pkg, version)

	cmd := exec.Command(config.Config.Pypi.GetPythonPath(), "-m", "pip", "install", fmt.Sprintf("%s==%s", pkg, version))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install package %s version %s: %w (output: %s)", pkg, version, err, string(output))
	}

	_, err = s.checkDependencies(pkg)
	if err != nil {
		return err
	}

	return nil
}

func (s *Scraper) checkDependencies(pkg string) (dependencies []*types.Package, err error) {
	cmd := exec.Command(config.Config.Pypi.GetPythonPath(), "-m", "pip", "show", pkg)
	showOutput, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get package info for %s: %w (output: %s)", pkg, err, string(showOutput))
	}

	p := types.Package{Name: pkg}

	for _, line := range strings.Split(string(showOutput), "\n") {
		if strings.Contains(line, "Version:") {
			p.Version = strings.TrimSpace(strings.TrimPrefix(line, "Version:"))
		}
		if strings.HasPrefix(line, "Requires:") {
			requiresText := strings.TrimPrefix(line, "Requires: ")
			if requiresText == "" {
				continue
			}
			requires := strings.Split(requiresText, ", ")

			for _, requirement := range requires {
				dependsOn, err := s.checkDependencies(requirement)
				if err != nil {
					return nil, fmt.Errorf("error processing dependency %s: %w", requirement, err)
				}

				for _, d := range dependsOn {
					_, err = s.dependencyStore.Upsert(context.Background(), types.Dependency{
						Package:   p,
						DependsOn: types.Package{Name: d.Name, Version: d.Version},
					})
					if err != nil {
						return nil, err
					}
				}
			}
		}
	}

	return dependencies, nil
}
