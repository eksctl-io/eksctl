package identitymapping

import (
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
)

type Manager struct {
	rawClient  kubernetes.RawClientInterface
	acmManager authconfigmap.AuthConfigMapManager
}

func New(rawClient kubernetes.RawClientInterface, acm authconfigmap.AuthConfigMapManager) *Manager {
	return &Manager{
		rawClient:  rawClient,
		acmManager: acm,
	}
}
