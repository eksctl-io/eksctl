package nodebootstrap

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/weaveworks/eksctl/pkg/nodebootstrap/assets"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap/utils"

	kubeletapi "k8s.io/kubelet/config/v1beta1"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cloudconfig"
)

const (
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
func NewBootstrapper(clusterConfig *api.ClusterConfig, ng *api.NodeGroup) (Bootstrapper, error) {
	clusterDNS := ng.ClusterDNS
	if clusterDNS == "" {
		var err error
		clusterDNS, err = GetClusterDNS(clusterConfig)
		if err != nil {
			return nil, err
		}
	}
	if api.IsWindowsImage(ng.AMIFamily) {
		return NewWindowsBootstrapper(clusterConfig, ng, clusterDNS), nil
	}
	switch ng.AMIFamily {
	case api.NodeImageFamilyUbuntuPro2404, api.NodeImageFamilyUbuntu2404, api.NodeImageFamilyUbuntuPro2204, api.NodeImageFamilyUbuntu2204, api.NodeImageFamilyUbuntuPro2004:
		return NewUbuntuBootstrapper(clusterConfig, ng, clusterDNS), nil
	case api.NodeImageFamilyBottlerocket:
		return NewBottlerocketBootstrapper(clusterConfig, ng), nil
	case api.NodeImageFamilyAmazonLinux2023:
		return NewAL2023Bootstrapper(clusterConfig, ng, clusterDNS), nil
	case api.NodeImageFamilyAmazonLinux2:
		return NewAL2Bootstrapper(clusterConfig, ng, clusterDNS), nil
	default:
		return nil, fmt.Errorf("unrecognized AMI family %q for creating bootstrapper", ng.AMIFamily)

	}
}

// NewManagedBootstrapper creates a new bootstrapper for managed nodegroups based on the AMI family
func NewManagedBootstrapper(clusterConfig *api.ClusterConfig, ng *api.ManagedNodeGroup) (Bootstrapper, error) {
	clusterDNS, err := GetClusterDNS(clusterConfig)
	if err != nil {
		return nil, err
	}
	if api.IsWindowsImage(ng.AMIFamily) {
		return &ManagedWindows{
			NodeGroup: ng,
		}, nil
	}
	switch ng.AMIFamily {
	case api.NodeImageFamilyAmazonLinux2023:
		return NewManagedAL2023Bootstrapper(clusterConfig, ng, clusterDNS), nil
	case api.NodeImageFamilyAmazonLinux2:
		return NewManagedAL2Bootstrapper(ng), nil
	case api.NodeImageFamilyBottlerocket:
		return NewManagedBottlerocketBootstrapper(clusterConfig, ng), nil
	case api.NodeImageFamilyUbuntuPro2004, api.NodeImageFamilyUbuntu2204, api.NodeImageFamilyUbuntuPro2204, api.NodeImageFamilyUbuntu2404, api.NodeImageFamilyUbuntuPro2404:
		return NewUbuntuBootstrapper(clusterConfig, ng, clusterDNS), nil
	}
	return nil, nil
}

// GetClusterDNS returns the DNS address to use
func GetClusterDNS(clusterConfig *api.ClusterConfig) (string, error) {
	networkConfig := clusterConfig.Status.KubernetesNetworkConfig
	if networkConfig == nil {
		return "", nil
	}

	var (
		serviceCIDR  string
		toClusterDNS func(net.IP) string
	)

	if networkConfig.ServiceIPv4CIDR != "" {
		serviceCIDR = networkConfig.ServiceIPv4CIDR
		toClusterDNS = func(parsedIP net.IP) string {
			ip := parsedIP.To4()
			ip[net.IPv4len-1] = 10
			return ip.String()
		}
	}
	if networkConfig.ServiceIPv6CIDR != "" {
		serviceCIDR = networkConfig.ServiceIPv6CIDR
		toClusterDNS = func(parsedIP net.IP) string {
			ip := parsedIP.To16()
			ip[net.IPv6len-1] = 10
			return ip.String()
		}
	}

	parsedIP, _, err := net.ParseCIDR(serviceCIDR)
	if err != nil {
		return "", fmt.Errorf("unexpected error parsing KubernetesNetworkConfig service CIDR: %q: %w", serviceCIDR, err)
	}
	return toClusterDNS(parsedIP), nil
}

func linuxConfig(clusterConfig *api.ClusterConfig, bootScriptName, bootScriptContent, clusterDNS string, np api.NodePool, scripts ...script) (string, error) {
	config := cloudconfig.New()
	ng := np.BaseNodeGroup()

	for _, command := range ng.PreBootstrapCommands {
		config.AddShellCommand(command)
	}

	var files []cloudconfig.File
	if len(scripts) == 0 {
		scripts = []script{}
	}

	if ng.OverrideBootstrapCommand != nil {
		config.AddShellCommand(*ng.OverrideBootstrapCommand)
	} else {
		scripts = append(scripts, script{name: bootScriptName, contents: bootScriptContent})
	}
	scripts = append(scripts, script{name: commonLinuxBootScript, contents: assets.BootstrapHelperSh})
	var kubeletExtraConf *api.InlineDocument
	if unmanaged, ok := np.(*api.NodeGroup); ok {
		kubeletExtraConf = unmanaged.KubeletExtraConfig
	}
	kubeletConf, err := makeKubeletExtraConf(kubeletExtraConf)
	if err != nil {
		return "", err
	}
	files = append(files, kubeletConf)
	envFile := makeBootstrapEnv(clusterConfig, np, clusterDNS)
	files = append(files, envFile)

	if err := addFilesAndScripts(config, files, scripts); err != nil {
		return "", err
	}

	body, err := config.Encode()
	if err != nil {
		return "", fmt.Errorf("encoding user data: %w", err)
	}

	return body, nil
}

func makeKubeletExtraConf(kubeletExtraConf *api.InlineDocument) (cloudconfig.File, error) {
	if kubeletExtraConf == nil {
		kubeletExtraConf = &api.InlineDocument{}
	}
	data, err := json.Marshal(kubeletExtraConf)
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

func makeBootstrapEnv(clusterConfig *api.ClusterConfig, np api.NodePool, clusterDNS string) cloudconfig.File {
	ng := np.BaseNodeGroup()
	variables := map[string]string{
		"CLUSTER_NAME":   clusterConfig.Metadata.Name,
		"API_SERVER_URL": clusterConfig.Status.Endpoint,
		"B64_CLUSTER_CA": base64.StdEncoding.EncodeToString(clusterConfig.Status.CertificateAuthorityData),
		"NODE_LABELS":    formatLabels(ng.Labels),
		"NODE_TAINTS":    utils.FormatTaints(np.NGTaints()),
	}
	if id := clusterConfig.Status.ID; id != "" {
		variables["CLUSTER_ID"] = id
	}
	if clusterConfig.IsControlPlaneOnOutposts() {
		variables["ENABLE_LOCAL_OUTPOST"] = strconv.FormatBool(true)
	}

	if ng.MaxPodsPerNode > 0 {
		variables["MAX_PODS"] = strconv.Itoa(ng.MaxPodsPerNode)
	}

	if clusterDNS != "" {
		variables["CLUSTER_DNS"] = clusterDNS
	}

	if unmanaged, ok := np.(*api.NodeGroup); ok && ng.AMIFamily == api.NodeImageFamilyAmazonLinux2 {
		variables["CONTAINER_RUNTIME"] = unmanaged.GetContainerRuntime()
	}

	return cloudconfig.File{
		Path:    configDir + envFile,
		Content: makeKeyValues(variables, "\n"),
	}
}

func makeKeyValues(kv map[string]string, separator string) string {
	var params []string
	for k, v := range kv {
		params = append(params, fmt.Sprintf("%s=%s", k, v))
	}

	return strings.Join(params, separator)
}

func formatLabels(labels map[string]string) string {
	return makeKeyValues(labels, ",")
}

type script struct {
	name     string
	contents string
}

func addFilesAndScripts(config *cloudconfig.CloudConfig, files []cloudconfig.File, scripts []script) error {
	for _, file := range files {
		config.AddFile(file)
	}

	for _, s := range scripts {
		config.RunScript(s.name, s.contents)
	}

	return nil
}
