package legacy

import (
	"strings"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cloudconfig"
	"github.com/weaveworks/eksctl/pkg/utils"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
)

type AL2Bootstrapper struct {
	clusterSpec *api.ClusterConfig
	ng          *api.NodeGroup
}

func NewAL2Bootstrapper(clusterSpec *api.ClusterConfig, ng *api.NodeGroup) AL2Bootstrapper {
	return AL2Bootstrapper{
		clusterSpec: clusterSpec,
		ng:          ng,
	}
}

func (b AL2Bootstrapper) UserData() (string, error) {
	config := cloudconfig.New()

	files, err := makeAmazonLinux2Config(b.clusterSpec, b.ng)
	if err != nil {
		return "", err
	}

	var scripts []string

	if b.ng.SSH.EnableSSM != nil && *b.ng.SSH.EnableSSM {
		scripts = append(scripts, "install-ssm.al2.sh")
	}

	// When using GPU instance types, the daemon.json is removed and a service
	// override file used instead. We can alter the daemon command by adding
	// to the OPTIONS var in /etc/sysconfig/docker
	overrideInsert := "sed -i 's/^OPTIONS=\"/&--exec-opt native.cgroupdriver=systemd /' /etc/sysconfig/docker"
	if utils.IsGPUInstanceType(b.ng.InstanceType) {
		config.AddShellCommand(overrideInsert)
	}
	if api.HasMixedInstances(b.ng) {
		for _, it := range b.ng.InstancesDistribution.InstanceTypes {
			if utils.IsGPUInstanceType(it) {
				config.AddShellCommand(overrideInsert)
			}
		}
	}

	for _, command := range b.ng.PreBootstrapCommands {
		config.AddShellCommand(command)
	}

	if b.ng.OverrideBootstrapCommand != nil {
		config.AddShellCommand(*b.ng.OverrideBootstrapCommand)
	} else {
		if api.IsEnabled(b.ng.EFAEnabled) {
			scripts = append(scripts, "efa.al2.sh")
		}
		scripts = append(scripts, "bootstrap.legacy.al2.sh")
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
		name:    "10-eksctl.al2.conf",
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
		dockerConfigData, err := makeDockerConfigJSON()
		if err != nil {
			return nil, err
		}

		files = append(files, configFile{dir: dockerConfigDir, name: "daemon.json", contents: dockerConfigData})
	}

	return files, nil
}
