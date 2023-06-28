package executor

import (
	"fmt"
	"os"
	"os/exec"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate -o fakes/fake_executor.go . Executor
// Executor executes commands shelling out and binding the stdout and stderr to the os ones
type Executor interface {
	Exec(command string, args ...string) error
	ExecWithOut(command string, args ...string) ([]byte, error)
	ExecInDir(command string, dir string, args ...string) error
}

type EnvVars map[string]string

// ShellExecutor an executor that shells out to run commands
type ShellExecutor struct {
	envVars EnvVars
}

// NewShellExecutor creates a new executor that runs commands
func NewShellExecutor(envVars EnvVars) Executor {
	return ShellExecutor{
		envVars: envVars,
	}
}

// Exec execute the command with the specified args
func (e ShellExecutor) Exec(command string, args ...string) error {
	return e.buildCmd(command, true, args...).Run()
}

// ExecWithOut executes the command with the specified args and returns the output
func (e ShellExecutor) ExecWithOut(command string, args ...string) ([]byte, error) {
	return e.buildCmd(command, false, args...).Output()
}

// ExecInDir executes the command inside the directory with the specified args
func (e ShellExecutor) ExecInDir(command string, dir string, args ...string) error {
	cmd := e.buildCmd(command, true, args...)
	cmd.Dir = dir

	return cmd.Run()
}

func (e ShellExecutor) buildCmd(command string, setOut bool, args ...string) *exec.Cmd {
	cmd := exec.Command(command, args...)

	cmd.Env = os.Environ()
	for k, v := range e.envVars {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	if setOut {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd
}
