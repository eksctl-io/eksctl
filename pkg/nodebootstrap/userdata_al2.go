package nodebootstrap

import (
	"strings"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/cloudconfig"
	"github.com/weaveworks/eksctl/pkg/eks/api"
)

func makeAmazonLinux2Config(spec *api.ClusterConfig, nodeGroupID int) (configFiles, error) {
	clientConfigData, err := makeClientConfigData(spec, nodeGroupID)
	if err != nil {
		return nil, err
	}

	files := configFiles{
		kubeletDropInUnitDir: {
			"10-eksclt.al2.conf": {isAsset: true},
		},
		configDir: {
			"metadata.env": {content: strings.Join(makeMetadata(spec), "\n")},
			"kubelet.env":  {content: strings.Join(makeKubeletParams(spec, nodeGroupID), "\n")},
			// TODO: https://github.com/weaveworks/eksctl/issues/161
			"kubelet-config.json": {isAsset: true},
			"ca.crt":              {content: string(spec.CertificateAuthorityData)},
			"kubeconfig.yaml":     {content: string(clientConfigData)},
		},
	}

	return files, nil
}

// NewUserDataForAmazonLinux2 creates new user data for Amazon Linux 2 nodes
func NewUserDataForAmazonLinux2(spec *api.ClusterConfig, nodeGroupID int) (string, error) {
	config := cloudconfig.New()

	scripts := []string{
		"bootstrap.al2.sh",
	}

	files, err := makeAmazonLinux2Config(spec, nodeGroupID)
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
