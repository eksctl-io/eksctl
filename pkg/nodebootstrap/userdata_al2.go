package nodebootstrap

import (
	"strings"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cloudconfig"
	"github.com/weaveworks/eksctl/pkg/utils"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
)

func makeAmazonLinux2Config(spec *api.ClusterConfig, ng *api.NodeGroup) ([]configFile, error) {
	clientConfigData, err := makeClientConfigData(spec, kubeconfig.AWSEKSAuthenticator)
	if err != nil {
		return nil, err
	}

	if len(spec.Status.CertificateAuthorityData) == 0 {
		return nil, errors.New("invalid cluster config: missing CertificateAuthorityData")
	}

	kubeletConfigData, err := makeKubeletConfigYAML(spec, ng)
	if err != nil {
		return nil, err
	}

	files := []configFile{{
		dir:     kubeletDropInUnitDir,
		name:    "10-eksclt.al2.conf",
		isAsset: true,
	}, {
		dir:      configDir,
		name:     "metadata.env",
		contents: strings.Join(makeMetadata(spec), "\n"),
	}, {
		dir:      configDir,
		name:     "kubelet.env",
		contents: strings.Join(makeCommonKubeletEnvParams(ng), "\n"),
	}, {
		dir:      configDir,
		name:     "kubelet.yaml",
		contents: string(kubeletConfigData),
	}, {
		dir:      configDir,
		name:     "ca.crt",
		contents: string(spec.Status.CertificateAuthorityData),
	}, {
		dir:      configDir,
		name:     "kubeconfig.yaml",
		contents: string(clientConfigData),
	}, {
		dir:      configDir,
		name:     "max_pods.map",
		contents: makeMaxPodsMapping(),
	}}

	if !utils.IsGPUInstanceType(ng.InstanceType) {
		dockerConfigData, err := makeDockerConfigJSON(ng)
		if err != nil {
			return nil, err
		}

		files = append(files, configFile{dir: dockerConfigDir, name: "daemon.json", contents: string(dockerConfigData)})
	}

	return files, nil
}

// NewUserDataForAmazonLinux2 creates new user data for Amazon Linux 2 nodes
func NewUserDataForAmazonLinux2(spec *api.ClusterConfig, ng *api.NodeGroup) (string, error) {
	config := cloudconfig.New()

	files, err := makeAmazonLinux2Config(spec, ng)
	if err != nil {
		return "", err
	}

	var scripts []string

	if ng.SSH.EnableSSM != nil && *ng.SSH.EnableSSM {
		scripts = append(scripts, "install-ssm.al2.sh")
	}

	// When using GPU instance types, the daemon.json is removed and a service
	// override file used instead. We can alter the daemon command by adding
	// to the OPTIONS var in /etc/sysconfig/docker
	overrideInsert := "sed -i 's/^OPTIONS=\"/&--exec-opt native.cgroupdriver=systemd /' /etc/sysconfig/docker"
	if utils.IsGPUInstanceType(ng.InstanceType) {
		config.AddShellCommand(overrideInsert)
	}
	if api.HasMixedInstances(ng) {
		for _, it := range ng.InstancesDistribution.InstanceTypes {
			if utils.IsGPUInstanceType(it) {
				config.AddShellCommand(overrideInsert)
			}
		}
	}

	for _, command := range ng.PreBootstrapCommands {
		config.AddShellCommand(command)
	}

	if ng.OverrideBootstrapCommand != nil {
		config.AddShellCommand(*ng.OverrideBootstrapCommand)
	} else {
		if api.IsEnabled(ng.EFAEnabled) {
			scripts = append(scripts, "efa.al2.sh")
		}
		scripts = append(scripts, "bootstrap.al2.sh")
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
