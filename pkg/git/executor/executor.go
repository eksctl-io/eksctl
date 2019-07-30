package executor

import (
	"context"
	"os"
	"os/exec"
	"time"
)

// Executor executes commands shelling out and binding the stdout and stderr to the os ones
type Executor interface {
	Exec(command string, dir string, args ...string) error
}

// ShellExecutor an executor that shells out to run commands
type ShellExecutor struct {
	parentCtx context.Context
	timeout   time.Duration
}

// NewShellExecutor creates a new executor that runs commands
func NewShellExecutor(ctx context.Context, timeout time.Duration) Executor {
	return ShellExecutor{
		parentCtx: ctx,
		timeout:   timeout,
	}
}

// Exec execute the command inside the directory with the specified args
func (e ShellExecutor) Exec(command string, dir string, args ...string) error {
	ctx, ctxCancel := context.WithTimeout(e.parentCtx, e.timeout)
	defer ctxCancel()
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	return cmd.Run()
}
