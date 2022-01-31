package nodebootstrap

import (
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap/assets"
)

const (
	al2BootScript = "bootstrap.al2.sh"
)

type AmazonLinux2 struct {
	clusterConfig *api.ClusterConfig
	ng            *api.NodeGroup
}

func NewAL2Bootstrapper(clusterConfig *api.ClusterConfig, ng *api.NodeGroup) *AmazonLinux2 {
	return &AmazonLinux2{
		clusterConfig: clusterConfig,
		ng:            ng,
	}
}

func (b *AmazonLinux2) UserData() (string, error) {
	var scripts []script

	if api.IsEnabled(b.ng.EFAEnabled) {
		scripts = append(scripts, script{name: "efa.al2.sh", contents: assets.EfaAl2Sh})
	}

	body, err := linuxConfig(b.clusterConfig, al2BootScript, assets.BootstrapAl2Sh, b.ng, scripts...)
	if err != nil {
		return "", errors.Wrap(err, "encoding user data")
	}

	logger.Debug("user-data = %s", body)
	return body, nil
}
