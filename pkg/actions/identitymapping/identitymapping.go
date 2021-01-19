package identitymapping

import (
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	kubewrapper "github.com/weaveworks/eksctl/pkg/kubernetes"
)

type Manager struct {
	rawClient *kubewrapper.RawClient
	acm       *authconfigmap.AuthConfigMap
}

func New(rawClient *kubewrapper.RawClient, acm *authconfigmap.AuthConfigMap) *Manager {

	return &Manager{
		rawClient: rawClient,
		acm:       acm,
	}
}
