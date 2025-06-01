package pypi

import (
	"fmt"
	"io"
	"mussel/internal/config"
	"os"
	"os/exec"
	"strings"

	"mussel/internal/types"
)

func Scrape(pkg string) (tree *types.DependencyTree, err error) {
	err = config.Config.Pypi.Setup()
	if err != nil {
		return nil, fmt.Errorf("failed to setup environment: %w", err)
	}

	err = config.Config.Pypi.InstallPython()
	if err != nil {
		return nil, fmt.Errorf("failed to install python: %w", err)
	}

	err = config.Config.Pypi.NewVenv()
	if err != nil {
		return nil, fmt.Errorf("failed to create virtual environment: %w", err)
	}

	packageVersions, err := getPackageVersions(pkg)
	if err != nil {
		return nil, fmt.Errorf("failed to get package versions: %w", err)
	}

	for _, version := range strings.Split(*packageVersions, "\n") {
		tree, err = CheckDependencies(pkg, version)
		if err != nil {
			_, err = fmt.Fprintf(os.Stderr, "failed to process version %s: %v\n", version, err)
			if err != nil {
				return nil, err
			}
			continue
		}
		return tree, nil
	}
	return nil, fmt.Errorf("no valid versions found for package %s", pkg)
}

func getPackageVersions(pkg string) (packageVersions *string, err error) {
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
