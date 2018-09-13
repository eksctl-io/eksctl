package nodebootstrap

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/kubicorn/kubicorn/pkg/logger"

	"github.com/weaveworks/eksctl/pkg/cloudconfig"
	"github.com/weaveworks/eksctl/pkg/eks/api"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
)

//go:generate ${GOPATH}/bin/go-bindata -pkg ${GOPACKAGE} -prefix assets -modtime 1 -o assets.go assets

const (
	configDir            = "/etc/eksctl/"
	kubeletDropInUnitDir = "/etc/systemd/system/kubelet.service.d/"
)

type configFile struct {
	content string
	isAsset bool
}

type configFiles = map[string]map[string]configFile

func getAsset(name string) (string, os.FileInfo, error) {
	data, err := Asset(name)
	if err != nil {
		return "", nil, errors.Wrapf(err, "decoding embedded file %q", name)
	}
	info, err := AssetInfo(name)
	if err != nil {
		return "", nil, errors.Wrapf(err, "getting info for embedded file %q", name)
	}
	return string(data), info, nil
}

func addFilesAndScripts(config *cloudconfig.CloudConfig, files configFiles, scripts []string) error {
	for dir, fileNames := range files {
		for fileName, file := range fileNames {
			f := cloudconfig.File{
				Path: dir + fileName,
			}
			if file.isAsset {
				data, info, err := getAsset(fileName)
				if err != nil {
					return err
				}
				f.Content = data
				f.Permissions = fmt.Sprintf("%04o", uint(info.Mode()))
			} else {
				f.Content = file.content
			}
			config.AddFile(f)
		}
	}
	for _, scriptName := range scripts {
		data, _, err := getAsset(scriptName)
		if err != nil {
			return err
		}
		config.RunScript(scriptName, data)
	}
	return nil
}

func makeAmazonLinux2Config(config *cloudconfig.CloudConfig, spec *api.ClusterConfig) (configFiles, error) {
	if spec.MaxPodsPerNode == 0 {
		spec.MaxPodsPerNode = maxPodsPerNodeType[spec.NodeType]
	}
	// TODO: use componentconfig or kubelet config file – https://github.com/weaveworks/eksctl/issues/156
	kubeletParams := []string{
		fmt.Sprintf("MAX_PODS=%d", spec.MaxPodsPerNode),
		// TODO: this will need to change when we provide options for using different VPCs and CIDRs – https://github.com/weaveworks/eksctl/issues/158
		"CLUSTER_DNS=10.100.0.10",
	}

	metadata := []string{
		fmt.Sprintf("AWS_DEFAULT_REGION=%s", spec.Region),
		fmt.Sprintf("AWS_EKS_CLUSTER_NAME=%s", spec.ClusterName),
		fmt.Sprintf("AWS_EKS_ENDPOINT=%s", spec.Endpoint),
	}

	clientConfig, _, _ := kubeconfig.New(spec, "kubelet", configDir+"ca.crt")
	kubeconfig.AppendAuthenticator(clientConfig, spec, kubeconfig.AWSIAMAuthenticator)

	clientConfigData, err := clientcmd.Write(*clientConfig)
	if err != nil {
		return nil, errors.Wrap(err, "serialising kubeconfig for nodegroup")
	}

	files := configFiles{
		kubeletDropInUnitDir: {
			"10-eksclt.al2.conf": {isAsset: true},
		},
		configDir: {
			"metadata.env": {content: strings.Join(metadata, "\n")},
			"kubelet.env":  {content: strings.Join(kubeletParams, "\n")},
			// TODO: https://github.com/weaveworks/eksctl/issues/161
			"ca.crt":          {content: string(spec.CertificateAuthorityData)},
			"kubeconfig.yaml": {content: string(clientConfigData)},
		},
	}

	return files, nil
}

func NewUserDataForAmazonLinux2(spec *api.ClusterConfig) (string, error) {
	config := cloudconfig.New()

	scripts := []string{
		"bootstrap.al2.sh",
	}

	files, err := makeAmazonLinux2Config(config, spec)
	if err != nil {
		return "", err
	}

	if err := addFilesAndScripts(config, files, scripts); err != nil {
		return "", err
	}

	body, err := config.Encode()
	if err != nil {
		return "", errors.Wrap(err, "encoding user data")
	}

	logger.Debug("user-data = %s", string(body))
	return string(body), nil
}
