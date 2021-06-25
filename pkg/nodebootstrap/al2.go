package nodebootstrap

import (
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
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
	var scripts []string

	if api.IsEnabled(b.ng.SSH.EnableSSM) {
		scripts = append(scripts, "install-ssm.al2.sh")
	}

	if api.IsEnabled(b.ng.EFAEnabled) {
		scripts = append(scripts, "efa.al2.sh")
	}

	body, err := linuxConfig(b.clusterConfig, al2BootScript, b.ng, scripts...)
	if err != nil {
		return "", errors.Wrap(err, "encoding user data")
	}

	logger.Debug("user-data = %s", body)
	return body, nil
}
