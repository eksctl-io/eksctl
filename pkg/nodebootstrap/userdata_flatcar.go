package nodebootstrap

import (
	"strings"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cloudconfig"
)

func makeFlatcarConfig(spec *api.ClusterConfig, ng *api.NodeGroup) (configFiles, error) {
	clientConfigData, err := makeClientConfigData(spec, ng)
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

	files := configFiles{
		systemdUnitDir: {
			"flatcar/kubelet.service": {isAsset: true},
		},
		configDir: {
			"metadata.env": {content: strings.Join(makeMetadata(spec), "\n")},
			"kubelet.env":  {content: strings.Join(makeCommonKubeletEnvParams(spec, ng), "\n")},
			"kubelet.yaml": {content: string(kubeletConfigData)},
			// TODO: https://github.com/weaveworks/eksctl/issues/161
			"ca.crt":               {content: string(spec.Status.CertificateAuthorityData)},
			"kubeconfig.yaml":      {content: string(clientConfigData)},
			"max_pods.map":         {content: makeMaxPodsMapping()},
			"bootstrap.flatcar.sh": {isAsset: true},
		},
	}

	return files, nil
}

// NewUserDataForFlatcar creates new user data for Flatcar nodes
func NewUserDataForFlatcar(spec *api.ClusterConfig, ng *api.NodeGroup) (string, error) {
	config := cloudconfig.New()

	files, err := makeFlatcarConfig(spec, ng)
	if err != nil {
		return "", err
	}

	var scripts []string

	for _, command := range ng.PreBootstrapCommands {
		config.AddShellCommand(command)
	}

	if ng.OverrideBootstrapCommand != nil {
		config.AddShellCommand(*ng.OverrideBootstrapCommand)
	} else {
		scripts = append(scripts, "bootstrap.flatcar.sh")
	}

	if err = addFilesAndScripts(config, files, scripts); err != nil {
		return "", err
	}
	config.AddSystemdUnit("kubelet-first-time.service", true, "start", `[Unit]\nConditionPathExists=!/etc/eksctl/done\n\n[Service]\nType=oneshot\nExecStart=/etc/eksctl/bootstrap.flatcar.sh\n\n[Install]\nWantedBy=multi-user.target`)

	body, err := config.Encode()
	if err != nil {
		return "", errors.Wrap(err, "encoding user data")
	}

	logger.Debug("user-data = %s", body)
	return body, nil
}
