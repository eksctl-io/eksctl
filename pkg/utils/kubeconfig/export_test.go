package kubeconfig

func SetExecCommand(f ExecCommandFunc) {
	execCommand = f
}

func SetNewKubectlClient(f NewKubectlClientFunc) {
	newKubectlClient = f
}

func SetLookupAuthenticator(f func() (string, bool)) {
	lookupAuthenticator = f
}
