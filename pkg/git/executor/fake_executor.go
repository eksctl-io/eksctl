package executor

import "github.com/stretchr/testify/mock"

type FakeExecutor struct {
	mock.Mock
	Command string
	Dir     string
	Args    []string
}

func (e *FakeExecutor) Exec(command string, dir string, args ...string) error {
	e.Command = command
	e.Dir = dir
	e.Args = args
	return e.Called(command, dir, args).Error(0)
}
