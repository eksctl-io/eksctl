package identitymapping

import (
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
)

type Manager struct {
	rawClient kubernetes.RawClientInterface
	acm       *authconfigmap.AuthConfigMap
}

func New(rawClient kubernetes.RawClientInterface, acm *authconfigmap.AuthConfigMap) *Manager {

	return &Manager{
		rawClient: rawClient,
		acm:       acm,
	}
}
