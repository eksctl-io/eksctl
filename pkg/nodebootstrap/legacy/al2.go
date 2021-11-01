package legacy

import (
	"strings"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cloudconfig"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap/assets"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"

	//For go:embed
	_ "embed"
)

//go:embed scripts/bootstrap.legacy.al2.sh
var bootstrapLegacyAl2Sh string

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

	var scripts []script

	for _, command := range b.ng.PreBootstrapCommands {
		config.AddShellCommand(command)
	}

	if b.ng.OverrideBootstrapCommand != nil {
		config.AddShellCommand(*b.ng.OverrideBootstrapCommand)
	} else {
		if api.IsEnabled(b.ng.EFAEnabled) {
			scripts = append(scripts, script{name: "efa.al2.sh", contents: assets.EfaAl2Sh})
		}
		scripts = append(scripts, script{name: "bootstrap.legacy.al2.sh", contents: bootstrapLegacyAl2Sh})
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
		dir:      kubeletDropInUnitDir,
		name:     "10-eksctl.al2.conf",
		contents: assets.EksctlAl2Conf,
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

	return files, nil
}
