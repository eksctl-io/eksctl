package nodebootstrap

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap/utils"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cloudconfig"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap/bindata"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap/legacy"
	kubeletapi "k8s.io/kubelet/config/v1beta1"
)

//go:generate ${GOBIN}/go-bindata -pkg bindata -prefix assets -nometadata -o bindata/assets.go bindata/assets

const (
	dataDir               = "bindata/assets"
	configDir             = "/etc/eksctl/"
	envFile               = "kubelet.env"
	extraKubeConfFile     = "kubelet-extra.json"
	commonLinuxBootScript = "bootstrap.helper.sh"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate -o fakes/fake_bootstrapper.go . Bootstrapper
type Bootstrapper interface {
	// UserData returns userdata for bootstrapping nodes
	UserData() (string, error)
}

// NewBootstrapper returns the correct bootstrapper for the AMI family
func NewBootstrapper(clusterSpec *api.ClusterConfig, ng *api.NodeGroup) Bootstrapper {
	if api.IsWindowsImage(ng.AMIFamily) {
		return NewWindowsBootstrapper(clusterSpec.Metadata.Name, ng)
	}
	switch ng.AMIFamily {
	case api.NodeImageFamilyUbuntu2004, api.NodeImageFamilyUbuntu1804:
		// TODO remove
		if ng.CustomAMI {
			logger.Warning("Custom AMI detected for nodegroup %s, using legacy nodebootstrap mechanism. Please refer to https://github.com/weaveworks/eksctl/issues/3563 for upcoming breaking changes", ng.Name)
			return legacy.NewUbuntuBootstrapper(clusterSpec, ng)
		}
		return NewUbuntuBootstrapper(clusterSpec, ng)
	case api.NodeImageFamilyBottlerocket:
		return NewBottlerocketBootstrapper(clusterSpec, ng)
	case api.NodeImageFamilyAmazonLinux2:
		// TODO remove
		if ng.CustomAMI {
			logger.Warning("Custom AMI detected for nodegroup %s, using legacy nodebootstrap mechanism. Please refer to https://github.com/weaveworks/eksctl/issues/3563 for upcoming breaking changes", ng.Name)
			return legacy.NewAL2Bootstrapper(clusterSpec, ng)
		}
		return NewAL2Bootstrapper(clusterSpec, ng)
	}

	return nil
}

func linuxConfig(clusterConfig *api.ClusterConfig, bootScript string, ng *api.NodeGroup, scripts ...string) (string, error) {
	config := cloudconfig.New()

	for _, command := range ng.PreBootstrapCommands {
		config.AddShellCommand(command)
	}

	var files []cloudconfig.File
	if len(scripts) == 0 {
		scripts = []string{}
	}

	if ng.OverrideBootstrapCommand != nil {
		config.AddShellCommand(*ng.OverrideBootstrapCommand)
	} else {
		scripts = append(scripts, commonLinuxBootScript, bootScript)
		kubeletConf, err := makeKubeletExtraConf(ng)
		if err != nil {
			return "", err
		}
		envFile := makeBootstrapEnv(clusterConfig, ng)
		files = append(files, kubeletConf, envFile)
	}

	if err := addFilesAndScripts(config, files, scripts); err != nil {
		return "", err
	}

	body, err := config.Encode()
	if err != nil {
		return "", errors.Wrap(err, "encoding user data")
	}

	return body, nil
}

func makeKubeletExtraConf(ng *api.NodeGroup) (cloudconfig.File, error) {
	if ng.KubeletExtraConfig == nil {
		ng.KubeletExtraConfig = &api.InlineDocument{}
	}

	data, err := json.Marshal(ng.KubeletExtraConfig)
	if err != nil {
		return cloudconfig.File{}, err
	}

	// validate that data can be decoded as legit KubeletConfiguration
	if err := json.Unmarshal(data, &kubeletapi.KubeletConfiguration{}); err != nil {
		return cloudconfig.File{}, err
	}

	return cloudconfig.File{
		Path:    configDir + extraKubeConfFile,
		Content: string(data),
	}, nil
}

func makeBootstrapEnv(clusterConfig *api.ClusterConfig, ng *api.NodeGroup) cloudconfig.File {
	variables := []string{
		fmt.Sprintf("CLUSTER_NAME=%s", clusterConfig.Metadata.Name),
		fmt.Sprintf("API_SERVER_URL=%s", clusterConfig.Status.Endpoint),
		fmt.Sprintf("B64_CLUSTER_CA=%s", base64.StdEncoding.EncodeToString(clusterConfig.Status.CertificateAuthorityData)),
		fmt.Sprintf("NODE_LABELS=%s", kvs(ng.Labels)),
		fmt.Sprintf("NODE_TAINTS=%s", utils.FormatTaints(ng.Taints)),
	}

	if ng.ClusterDNS != "" {
		variables = append(variables, fmt.Sprintf("CLUSTER_DNS=%s", ng.ClusterDNS))
	}

	return cloudconfig.File{
		Path:    configDir + envFile,
		Content: strings.Join(variables, "\n"),
	}
}

func kvs(kv map[string]string) string {
	var params []string
	for k, v := range kv {
		params = append(params, fmt.Sprintf("%s=%s", k, v))
	}

	return strings.Join(params, ",")
}

func addFilesAndScripts(config *cloudconfig.CloudConfig, files []cloudconfig.File, scripts []string) error {
	for _, file := range files {
		config.AddFile(file)
	}

	for _, scriptName := range scripts {
		data, err := getAsset(scriptName)
		if err != nil {
			return err
		}
		config.RunScript(scriptName, data)
	}

	return nil
}

func getAsset(name string) (string, error) {
	data, err := bindata.Asset(filepath.Join(dataDir, name))
	if err != nil {
		return "", errors.Wrapf(err, "decoding embedded file %q", name)
	}

	return string(data), nil
}
