package nodebootstrap

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/kubicorn/kubicorn/pkg/logger"

	gfn "github.com/awslabs/goformation/cloudformation"

	"github.com/weaveworks/eksctl/pkg/cloudconfig"
	"github.com/weaveworks/eksctl/pkg/eks/api"
)

//go:generate go-bindata -pkg $GOPACKAGE -prefix assets -modtime 1 -o assets.go assets

const (
	configDir            = "/etc/eksctl/"
	kubeletDropInUnitDir = "/etc/systemd/system/kubelet.service.d/"
)

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

func addFilesAndScripts(config *cloudconfig.CloudConfig, files map[string][]string, scripts []string) error {
	for dir, fileNames := range files {
		for _, fileName := range fileNames {
			data, info, err := getAsset(fileName)
			if err != nil {
				return err
			}
			config.AddFile(cloudconfig.File{
				Path:        dir + fileName,
				Content:     data,
				Permissions: fmt.Sprintf("%04o", uint(info.Mode())),
			})
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

func setKubeletParams(config *cloudconfig.CloudConfig, spec *api.ClusterConfig) {
	if spec.MaxPodsPerNode == 0 {
		spec.MaxPodsPerNode = maxPodsPerNodeType[spec.NodeType]
	}
	// TODO: investigate if we can use componentconfig, or at least switch to kubelet config file
	kubeletParams := []string{
		fmt.Sprintf("MAX_PODS=%d", spec.MaxPodsPerNode),
		// TODO: this will need to change when we provide options for using different VPCs and CIDRs
		"CLUSTER_DNS=10.100.0.10",
	}

	config.AddFile(cloudconfig.File{
		Path:    configDir + "kubelet.env",
		Content: strings.Join(kubeletParams, "\n"),
	})
}

func NewUserDataForAmazonLinux2(spec *api.ClusterConfig) (*gfn.StringIntrinsic, error) {
	config := cloudconfig.New()

	scripts := []string{
		"get_metadata.sh",
		"get_credentials.sh",
	}
	files := map[string][]string{
		configDir: {
			"authenticator.sh",
			"kubeconfig.yaml",
		},
		kubeletDropInUnitDir: {
			"10-eksclt.al2.conf",
		},
	}

	setKubeletParams(config, spec)

	config.AddPackages("jq")
	config.AddCommand("pip", "install", "--upgrade", "awscli")

	if err := addFilesAndScripts(config, files, scripts); err != nil {
		return nil, err
	}

	config.AddCommand("systemctl", "daemon-reload")
	config.AddCommand("systemctl", "restart", "kubelet")

	body, err := config.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "encoding user data")
	}

	logger.Debug("user-data = %s", string(body))
	return gfn.NewString(string(body)), nil
}
