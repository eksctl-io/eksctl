package kubeconfig

func SetExecCommand(f ExecCommandFunc) {
	execCommand = f
}
