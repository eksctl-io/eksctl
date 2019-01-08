package nodebootstrap

import (
	"strings"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha3"
	"github.com/weaveworks/eksctl/pkg/cloudconfig"
)

func makeUbuntu1804Config(spec *api.ClusterConfig, ng *api.NodeGroup) (configFiles, error) {
	clientConfigData, err := makeClientConfigData(spec, ng)
	if err != nil {
		return nil, err
	}

	files := configFiles{
		configDir: {
			"metadata.env": {content: strings.Join(makeMetadata(spec), "\n")},
			"kubelet.env":  {content: strings.Join(makeKubeletParamsCommon(spec, ng), "\n")},
			// TODO: https://github.com/weaveworks/eksctl/issues/161
			"ca.crt":          {content: string(spec.Status.CertificateAuthorityData)},
			"kubeconfig.yaml": {content: string(clientConfigData)},
		},
	}

	return files, nil
}

// NewUserDataForUbuntu1804 creates new user data for Ubuntu 18.04 nodes
func NewUserDataForUbuntu1804(spec *api.ClusterConfig, ng *api.NodeGroup) (string, error) {
	config := cloudconfig.New()

	scripts := []string{
		"bootstrap.ubuntu.sh",
	}

	files, err := makeUbuntu1804Config(spec, ng)
	if err != nil {
		return "", err
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
