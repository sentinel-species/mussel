package cmd

import (
	"errors"
	"fmt"
	"os/exec"
)

var (
	RequiredCommands = []string{"curl", "pyenv"}

	ErrCommandNotFound = errors.New("required command not found in PATH")
)

func CheckRequirements() error {
	var missing []string

	for _, cmd := range RequiredCommands {
		if _, err := exec.LookPath(cmd); err != nil {
			missing = append(missing, cmd)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("%w: %v", ErrCommandNotFound, missing)
	}
	return nil
}
