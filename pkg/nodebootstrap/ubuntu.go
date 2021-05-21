package nodebootstrap

import (
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

const (
	ubuntuBootScript = "bootstrap.ubuntu.sh"
)

type Ubuntu struct {
	clusterConfig *api.ClusterConfig
	ng            *api.NodeGroup
}

func NewUbuntuBootstrapper(clusterConfig *api.ClusterConfig, ng *api.NodeGroup) *Ubuntu {
	return &Ubuntu{
		clusterConfig: clusterConfig,
		ng:            ng,
	}
}

func (b *Ubuntu) UserData() (string, error) {
	body, err := linuxConfig(b.clusterConfig, ubuntuBootScript, b.ng)
	if err != nil {
		return "", errors.Wrap(err, "encoding user data")
	}

	logger.Debug("user-data = %s", body)
	return body, nil
}
