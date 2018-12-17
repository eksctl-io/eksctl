package nodebootstrap

import (
	"strings"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/cloudconfig"
	"github.com/weaveworks/eksctl/pkg/eks/api"
)

func makeUbuntu1804Config(spec *api.ClusterConfig, nodeGroupID int) (configFiles, error) {
	clientConfigData, err := makeClientConfigData(spec, nodeGroupID)
	if err != nil {
		return nil, err
	}

	files := configFiles{
		configDir: {
			"metadata.env": {content: strings.Join(makeMetadata(spec), "\n")},
			"kubelet.env":  {content: strings.Join(makeKubeletParams(spec, nodeGroupID), "\n")},
			// TODO: https://github.com/weaveworks/eksctl/issues/161
			"ca.crt":          {content: string(spec.CertificateAuthorityData)},
			"kubeconfig.yaml": {content: string(clientConfigData)},
		},
	}

	return files, nil
}

// NewUserDataForUbuntu1804 creates new user data for Ubuntu 18.04 nodes
func NewUserDataForUbuntu1804(spec *api.ClusterConfig, nodeGroupID int) (string, error) {
	config := cloudconfig.New()

	scripts := []string{
		"bootstrap.ubuntu.sh",
	}

	files, err := makeUbuntu1804Config(spec, nodeGroupID)
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
