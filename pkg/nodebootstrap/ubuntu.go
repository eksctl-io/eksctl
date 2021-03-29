package nodebootstrap

import (
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cloudconfig"
)

const (
	ubuntuBootScript = "bootstrap.ubuntu.sh"
)

type Ubuntu struct {
	clusterName string
	ng          *api.NodeGroup
}

func NewUbuntuBootstrapper(clusterName string, ng *api.NodeGroup) *Ubuntu {
	return &Ubuntu{
		clusterName: clusterName,
		ng:          ng,
	}
}

func (b *Ubuntu) UserData() (string, error) {
	config := cloudconfig.New()

	for _, command := range b.ng.PreBootstrapCommands {
		config.AddShellCommand(command)
	}

	var (
		scripts []string
		files   []cloudconfig.File
	)

	if b.ng.OverrideBootstrapCommand != nil {
		config.AddShellCommand(*b.ng.OverrideBootstrapCommand)
	} else {
		scripts = append(scripts, commonLinuxBootScript, ubuntuBootScript)

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
