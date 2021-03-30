package nodebootstrap

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cloudconfig"
	kubeletapi "k8s.io/kubelet/config/v1beta1"
)

//go:generate ${GOBIN}/go-bindata -pkg ${GOPACKAGE} -prefix assets -nometadata -o assets.go assets

const (
	configDir             = "/etc/eksctl/"
	envFile               = "kubelet.env"
	extraConfFile         = "kubelet-extra.json"
	commonLinuxBootScript = "bootstrap.linux.sh"
)

//go:generate counterfeiter -o fakes/fake_bootstrapper.go . Bootstrapper
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
		return NewUbuntuBootstrapper(clusterSpec.Metadata.Name, ng)
	case api.NodeImageFamilyBottlerocket:
		return NewBottlerocketBootstrapper(clusterSpec, ng)
	default:
		return NewAL2Bootstrapper(clusterSpec.Metadata.Name, ng)
	}
}

func linuxConfig(bootScript, clusterName string, ng *api.NodeGroup, scripts ...string) (string, error) {
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

		if ng.KubeletExtraConfig != nil {
			kubeletConf, err := makeKubeletExtraConf(ng)
			if err != nil {
				return "", err
			}
			files = append(files, kubeletConf)
		}
		envFile := makeBootstrapEnv(clusterName, ng)

		files = append(files, envFile)
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
	data, err := json.Marshal(ng.KubeletExtraConfig)
	if err != nil {
		return cloudconfig.File{}, err
	}

	// validate that data can be decoded as legit KubeletConfiguration
	if err := json.Unmarshal(data, &kubeletapi.KubeletConfiguration{}); err != nil {
		return cloudconfig.File{}, err
	}

	return cloudconfig.File{
		Path:    configDir + extraConfFile,
		Content: string(data),
	}, nil
}

func makeBootstrapEnv(clusterName string, ng *api.NodeGroup) cloudconfig.File {
	variables := []string{
		fmt.Sprintf("NODE_LABELS=%s", kvs(ng.Labels)),
		fmt.Sprintf("NODE_TAINTS=%s", mapTaints(ng.Taints)),
		fmt.Sprintf("CLUSTER_NAME=%s", clusterName),
	}

	if ng.ClusterDNS != "" {
		variables = append(variables, fmt.Sprintf("CLUSTER_DNS=%s", ng.ClusterDNS))
	}

	return cloudconfig.File{
		Path:    configDir + envFile,
		Content: strings.Join(variables, "\n"),
	}
}

func mapTaints(kv map[string]string) string {
	var params []string
	for k, v := range kv {
		params = append(params, fmt.Sprintf("%s=:%s", k, v))
	}
	return strings.Join(params, ",")
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
	data, err := Asset(name)
	if err != nil {
		return "", errors.Wrapf(err, "decoding embedded file %q", name)
	}

	return string(data), nil
}
