package nodebootstrap

import (
	"fmt"
	"strings"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cloudconfig"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
)

const ubuntu2004ResolveConfPath = "/run/systemd/resolve/resolv.conf"

func makeUbuntuConfig(spec *api.ClusterConfig, ng *api.NodeGroup) (configFiles, error) {
	clientConfigData, err := makeClientConfigData(spec, kubeconfig.HeptioAuthenticatorAWS)
	if err != nil {
		return nil, err
	}

	if len(spec.Status.CertificateAuthorityData) == 0 {
		return nil, errors.New("invalid cluster config: missing CertificateAuthorityData")
	}

	kubeletEnvParams := append(makeCommonKubeletEnvParams(ng),
		fmt.Sprintf("CLUSTER_DNS=%s", clusterDNS(spec, ng)),
	)

	// Set resolvConf for Ubuntu 20.04 only, do not override user set value
	if ng.AMIFamily == api.NodeImageFamilyUbuntu2004 {
		if ng.KubeletExtraConfig == nil {
			ng.KubeletExtraConfig = &api.InlineDocument{}
		}
		if _, ok := (*ng.KubeletExtraConfig)["resolvConf"]; !ok {
			(*ng.KubeletExtraConfig)["resolvConf"] = ubuntu2004ResolveConfPath
		}
	}

	kubeletConfigData, err := makeKubeletConfigYAML(spec, ng)
	if err != nil {
		return nil, err
	}

	files := configFiles{
		configDir: {
			"metadata.env": {content: strings.Join(makeMetadata(spec), "\n")},
			"kubelet.env":  {content: strings.Join(kubeletEnvParams, "\n")},
			"kubelet.yaml": {content: string(kubeletConfigData)},
			// TODO: https://github.com/weaveworks/eksctl/issues/161
			"ca.crt":          {content: string(spec.Status.CertificateAuthorityData)},
			"kubeconfig.yaml": {content: string(clientConfigData)},
			"max_pods.map":    {content: makeMaxPodsMapping()},
		},
	}

	return files, nil
}

// NewUserDataForUbuntu creates new user data for Ubuntu 18.04 & 20.04 nodes
func NewUserDataForUbuntu(spec *api.ClusterConfig, ng *api.NodeGroup) (string, error) {
	config := cloudconfig.New()

	files, err := makeUbuntuConfig(spec, ng)
	if err != nil {
		return "", err
	}

	scripts := []string{}

	for _, command := range ng.PreBootstrapCommands {
		config.AddShellCommand(command)
	}

	if ng.OverrideBootstrapCommand != nil {
		config.AddShellCommand(*ng.OverrideBootstrapCommand)
	} else {
		scripts = append(scripts, "bootstrap.ubuntu.sh")
	}

	if err = addFilesAndScripts(config, files, scripts); err != nil {
		return "", err
	}

	body, err := config.Encode()
	if err != nil {
		return "", errors.Wrap(err, "encoding user data")
	}

	logger.Debug("user-data = %s", body)
	return body, nil
}
