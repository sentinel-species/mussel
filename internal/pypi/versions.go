package pypi

import (
	"fmt"
	"mussel/internal/config"
	"os/exec"
	"strings"

	"mussel/internal/types"
)

func CheckDependencies(pkg string, version string) (*types.DependencyTree, error) {
	fmt.Printf("Processing %s version %s\n", pkg, version)

	cmd := exec.Command(config.Config.Pypi.GetPythonPath(), "-m", "pip", "install", fmt.Sprintf("%s==%s", pkg, version))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to install package %s version %s: %w (output: %s)", pkg, version, err, string(output))
	}

	tree, err := checkDependencies(pkg)
	if err != nil {
		return nil, err
	}

	return tree, nil
}

func checkDependencies(pkg string) (*types.DependencyTree, error) {
	cmd := exec.Command(config.Config.Pypi.GetPythonPath(), "-m", "pip", "show", pkg)
	showOutput, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get package info for %s: %w (output: %s)", pkg, err, string(showOutput))
	}

	tree := &types.DependencyTree{
		Name: pkg,
	}

	for _, line := range strings.Split(string(showOutput), "\n") {
		if strings.Contains(line, "Version:") {
			tree.Version = strings.TrimSpace(strings.TrimPrefix(line, "Version:"))
		}
		if strings.HasPrefix(line, "Requires:") {
			requires := strings.TrimPrefix(line, "Requires: ")
			if requires == "" {
				continue
			}
			dependencies := strings.Split(requires, ", ")

			tree.Dependencies = []*types.DependencyTree{}
			for _, dependency := range dependencies {
				dependencyTree, err := checkDependencies(dependency)
				if err != nil {
					return nil, fmt.Errorf("error processing dependency %s: %w", dependency, err)
				}
				tree.Dependencies = append(tree.Dependencies, dependencyTree)
			}
		}
	}

	return tree, nil
}
