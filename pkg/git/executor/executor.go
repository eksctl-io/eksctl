package executor

import (
	"os"
	"os/exec"
)

// Executor executes commands shelling out and binding the stdout and stderr to the os ones
type Executor interface {
	Exec(command string, dir string, args ...string) error
}

// ShellExecutor an executor that shells out to run commands
type ShellExecutor struct {
	envVars []string
}

// NewShellExecutor creates a new executor that runs commands
func NewShellExecutor(envVars []string) Executor {
	return ShellExecutor{
		envVars: envVars,
	}
}

// Exec execute the command inside the directory with the specified args
func (e ShellExecutor) Exec(command string, dir string, args ...string) error {
	cmd := exec.Command(command, args...)
	if len(e.envVars) > 0 {
		cmd.Env = e.envVars
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	return cmd.Run()
}
