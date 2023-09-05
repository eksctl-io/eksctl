package kubectl

import "os/exec"

func SetExecCommand(f func(name string, arg ...string) *exec.Cmd) {
	execCommand = f
}
