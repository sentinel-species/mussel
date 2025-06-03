package types

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type PypiConfig struct {
	PathPyEnvRoot  string   `yaml:"pyEnvRoot"`
	PythonVersions []string `yaml:"pythonVersions"`
	PathVenv       string   `yaml:"venv"`

	currentVersion int `yaml:"-"` // Internal field, not in YAML
}

func (pc *PypiConfig) getPythonVersion() string {
	return pc.PythonVersions[pc.currentVersion]
}

func (pc *PypiConfig) getFullPyenvVersion(majorMinor string) (majorMinorBugfix *string, err error) {
	var versions []string

	entries, err := os.ReadDir(fmt.Sprintf("%s/versions", pc.PathPyEnvRoot))
	if err != nil {
		return nil, fmt.Errorf("reading directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()

		if strings.HasPrefix(name, majorMinor+".") {
			versions = append(versions, name)
		}
	}

	if len(versions) == 0 {
		return nil, fmt.Errorf("no matching versions found for %s", majorMinor)
	}

	sort.Sort(sort.Reverse(sort.StringSlice(versions)))

	return Pointer(versions[0]), nil
}

func (pc *PypiConfig) Setup() (err error) {
	err = os.Setenv("PYENV_ROOT", pc.PathPyEnvRoot)
	if err != nil {
		return fmt.Errorf("failed to set PYENV_ROOT: %w", err)
	}

	pc.currentVersion = len(pc.PythonVersions) - 1

	return nil
}

func (pc *PypiConfig) GetPythonPath() string {
	return pc.PathVenv + "/bin/python"
}

func (pc *PypiConfig) DowngradePythonVersion() (string, error) {
	pc.currentVersion = (pc.currentVersion - 1 + len(pc.PythonVersions)) % len(pc.PythonVersions)
	return pc.getPythonVersion(), nil
}

func (pc *PypiConfig) InstallPython() (err error) {
	err = exec.Command("pyenv", "install", "--skip-existing", pc.getPythonVersion()).Run()
	if err != nil {
		return fmt.Errorf("failed to install python %s: %w", pc.getPythonVersion(), err)
	}

	err = exec.Command("pyenv", "rehash").Run()
	if err != nil {
		return fmt.Errorf("failed to rehash pyenv: %w", err)
	}

	return
}

func (pc *PypiConfig) NewVenv() (err error) {
	err = os.RemoveAll(pc.PathVenv)
	if err != nil {
		return fmt.Errorf("failed to remove old virtual environment: %w", err)
	}

	fullVersion, err := pc.getFullPyenvVersion(pc.getPythonVersion())
	if err != nil {
		return err
	}
	err = exec.Command(fmt.Sprintf("%s/versions/%s/bin/python", pc.PathPyEnvRoot, *fullVersion), "-m", "venv", pc.PathVenv).Run()
	if err != nil {
		return fmt.Errorf("failed to create virtual environment: %w", err)
	}
	return
}
