package nodebootstrap

import (
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cloudconfig"
)

const (
	al2BootScript = "bootstrap.al2.sh"
)

type AmazonLinux2 struct {
	clusterName string
	ng          *api.NodeGroup
}

func NewAL2Bootstrapper(clusterName string, ng *api.NodeGroup) *AmazonLinux2 {
	return &AmazonLinux2{
		clusterName: clusterName,
		ng:          ng,
	}
}

func (b *AmazonLinux2) UserData() (string, error) {
	config := cloudconfig.New()

	var (
		scripts []string
		files   []cloudconfig.File
	)

	if api.IsEnabled(b.ng.SSH.EnableSSM) {
		scripts = append(scripts, "install-ssm.al2.sh")
	}

	if api.IsEnabled(b.ng.EFAEnabled) {
		scripts = append(scripts, "efa.al2.sh")
	}

	for _, command := range b.ng.PreBootstrapCommands {
		config.AddShellCommand(command)
	}

	if b.ng.OverrideBootstrapCommand != nil {
		config.AddShellCommand(*b.ng.OverrideBootstrapCommand)
	} else {
		scripts = append(scripts, commonLinuxBootScript, al2BootScript)

		// TODO: should this happen even if override is set? do more scripting
		if b.ng.KubeletExtraConfig != nil {
			kubeletConf, err := makeKubeletExtraConf(b.ng)
			if err != nil {
				return "", err
			}
			files = append(files, kubeletConf)
		}
		envFile := makeBootstrapEnv(b.clusterName, b.ng)

		files = append(files, envFile)
	}

	if err := addFilesAndScripts(config, files, scripts); err != nil {
		return "", err
	}

	body, err := config.Encode()
	if err != nil {
		return "", errors.Wrap(err, "encoding user data")
	}

	logger.Debug("user-data = %s", body)
	return body, nil
}

// TODO instance distribution?
