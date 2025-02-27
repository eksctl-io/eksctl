package nodebootstrap

import (
	"fmt"

	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap/assets"
)

const (
	al2BootScript = "bootstrap.al2.sh"
)

type AmazonLinux2 struct {
	clusterConfig *api.ClusterConfig
	ng            *api.NodeGroup
	clusterDNS    string
}

func NewAL2Bootstrapper(clusterConfig *api.ClusterConfig, ng *api.NodeGroup, clusterDNS string) *AmazonLinux2 {
	return &AmazonLinux2{
		clusterConfig: clusterConfig,
		ng:            ng,
		clusterDNS:    clusterDNS,
	}
}

func (b *AmazonLinux2) UserData() (string, error) {
	var scripts []script

	body, err := linuxConfig(b.clusterConfig, al2BootScript, assets.BootstrapAl2Sh, b.clusterDNS, b.ng, scripts...)
	if err != nil {
		return "", fmt.Errorf("encoding user data: %w", err)
	}

	logger.Debug("user-data = %s", body)
	return body, nil
}
