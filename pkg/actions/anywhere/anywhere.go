package anywhere

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/weaveworks/eksctl/pkg/version"
)

const (
	// BinaryFileName defines the name of the anywhere binary.
	BinaryFileName = "eksctl-anywhere"
)

// IsAnywhereCommand detects whether the user would like to execute anywhere based commands.
func IsAnywhereCommand(args []string) (bool, error) {
	if len(args) == 0 {
		return false, nil
	}

	if args[0] == "anywhere" {
		return true, nil
	}

	// if they have any args/flags before the anywhere command we should error
	// e.g. eksctl --foo=bar anywhere should error
	for _, arg := range args {
		if arg == "anywhere" {
			return false, fmt.Errorf("flags cannot be placed before the anywhere command")
		}
	}

	return false, nil
}

// RunAnywhereCommand executes the anywhere binary.
func RunAnywhereCommand(args []string) (int, error) {
	if _, err := exec.LookPath(BinaryFileName); errors.Is(err, exec.ErrNotFound) {
		return 1, fmt.Errorf(fmt.Sprintf("%q plugin was not found on your path", BinaryFileName))
	} else if err != nil {
		return 1, fmt.Errorf("failed to lookup anywhere plugin: %w", err)
	}

	cmd := exec.Command(BinaryFileName, args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	cmd.Env = append(os.Environ(), fmt.Sprintf("EKSCTL_VERSION=%s", version.GetVersion()))

	err := cmd.Run()
	if exiterr, ok := err.(*exec.ExitError); ok {
		return exiterr.ExitCode(), nil
	}
	return 0, err
}
