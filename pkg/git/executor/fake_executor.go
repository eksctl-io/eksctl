package executor

import "github.com/stretchr/testify/mock"

// FakeExecutor is a test double that records the arguments used for calling it
type FakeExecutor struct {
	mock.Mock
	Command string
	Dir     string
	Args    []string
}

// Exec records the arguments used to call it
func (e *FakeExecutor) Exec(command string, dir string, args ...string) error {
	e.Command = command
	e.Dir = dir
	e.Args = args
	called := e.Called(command, dir, args)
	return called.Error(0)
}
