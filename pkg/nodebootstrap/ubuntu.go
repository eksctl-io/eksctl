package nodebootstrap

import (
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap/assets"
)

const (
	ubuntuBootScript = "bootstrap.ubuntu.sh"
)

type Ubuntu struct {
	clusterConfig *api.ClusterConfig
	np            api.NodePool
	clusterDNS    string
}

func NewUbuntuBootstrapper(clusterConfig *api.ClusterConfig, np api.NodePool, clusterDNS string) *Ubuntu {
	return &Ubuntu{
		clusterConfig: clusterConfig,
		np:            np,
		clusterDNS:    clusterDNS,
	}
}

func (b *Ubuntu) UserData() (string, error) {
	body, err := linuxConfig(b.clusterConfig, ubuntuBootScript, assets.BootstrapUbuntuSh, b.clusterDNS, b.np)
	if err != nil {
		return "", errors.Wrap(err, "encoding user data")
	}

	logger.Debug("user-data = %s", body)
	return body, nil
}
