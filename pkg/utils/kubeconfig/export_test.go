package kubeconfig

import (
	"os/exec"

	"github.com/weaveworks/eksctl/pkg/utils/kubectl"
)

func SetExecCommand(f func(name string, arg ...string) *exec.Cmd) {
	execCommand = f
}

func SetExecLookPath(f func(file string) (string, error)) {
	execLookPath = f
}

func SetNewVersionManager(f func() kubectl.KubernetesVersionManager) {
	newVersionManager = f
}

func SetLookupAuthenticator(f func() (string, bool)) {
	lookupAuthenticator = f
}
