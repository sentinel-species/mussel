package types

import (
	"fmt"
	"os"
	"os/exec"
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
	// TODO: Start hack week here

	return
}
