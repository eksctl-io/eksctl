package identitymapping

import (
	kubewrapper "github.com/weaveworks/eksctl/pkg/kubernetes"
	"k8s.io/client-go/kubernetes"
)

type Manager struct {
	rawClient *kubewrapper.RawClient
	clientSet *kubernetes.Clientset
}

func New(rawClient *kubewrapper.RawClient, clientSet *kubernetes.Clientset) *Manager {
	return &Manager{
		rawClient: rawClient,
		clientSet: clientSet,
	}
}
