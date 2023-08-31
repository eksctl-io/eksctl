package kubectl

import "os/exec"

func SetExecCommand(f func(name string, arg ...string) *exec.Cmd) {
	execCommand = f
}

func SetExecLookPath(f func(file string) (string, error)) {
	execLookPath = f
}
