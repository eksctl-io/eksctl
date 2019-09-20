package nodebootstrap

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"k8s.io/client-go/tools/clientcmd"
	kubeletapi "k8s.io/kubelet/config/v1beta1"

	"sigs.k8s.io/yaml"

	"github.com/weaveworks/eksctl/pkg/ami"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cloudconfig"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
)

//go:generate ${GOBIN}/go-bindata -pkg ${GOPACKAGE} -prefix assets -nometadata -o assets.go assets
//go:generate go run ./maxpods_generate.go

const (
	configDir            = "/etc/eksctl/"
	kubeletDropInUnitDir = "/etc/systemd/system/kubelet.service.d/"
)

type configFile struct {
	content string
	isAsset bool
}

type configFiles = map[string]map[string]configFile

func getAsset(name string) (string, error) {
	data, err := Asset(name)
	if err != nil {
		return "", errors.Wrapf(err, "decoding embedded file %q", name)
	}
	return string(data), nil
}

func addFilesAndScripts(config *cloudconfig.CloudConfig, files configFiles, scripts []string) error {
	for dir, fileNames := range files {
		for fileName, file := range fileNames {
			f := cloudconfig.File{
				Path: dir + fileName,
			}
			if file.isAsset {
				data, err := getAsset(fileName)
				if err != nil {
					return err
				}
				f.Content = data
			} else {
				f.Content = file.content
			}
			config.AddFile(f)
		}
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

func makeClientConfigData(spec *api.ClusterConfig, ng *api.NodeGroup) ([]byte, error) {
	clientConfig, _, _ := kubeconfig.New(spec, "kubelet", configDir+"ca.crt")
	authenticator := kubeconfig.AWSIAMAuthenticator
	if ng.AMIFamily == ami.ImageFamilyUbuntu1804 {
		authenticator = kubeconfig.HeptioAuthenticatorAWS
	}
	kubeconfig.AppendAuthenticator(clientConfig, spec, authenticator, "", "")
	clientConfigData, err := clientcmd.Write(*clientConfig)
	if err != nil {
		return nil, errors.Wrap(err, "serialising kubeconfig for nodegroup")
	}
	return clientConfigData, nil
}

func clusterDNS(spec *api.ClusterConfig, ng *api.NodeGroup) string {
	if ng.ClusterDNS != "" {
		return ng.ClusterDNS
	}
	// Default service network is 10.100.0.0, but it gets set 172.20.0.0 automatically when pod network
	// is anywhere within 10.0.0.0/8
	if spec.VPC.CIDR != nil && spec.VPC.CIDR.IP[0] == 10 {
		return "172.20.0.10"
	}
	return "10.100.0.10"
}

func makeKubeletConfigYAML(spec *api.ClusterConfig, ng *api.NodeGroup) ([]byte, error) {
	data, err := Asset("kubelet.yaml")
	if err != nil {
		return nil, err
	}

	// use a map here, as using struct will require us to add defaulting etc,
	// and we only need to add a few top-level fields
	obj := map[string]interface{}{}
	if err := yaml.UnmarshalStrict(data, &obj); err != nil {
		return nil, err
	}

	obj["clusterDNS"] = []string{
		clusterDNS(spec, ng),
	}

	// Add extra configuration from configfile
	if ng.KubeletExtraConfig != nil {
		for k, v := range *ng.KubeletExtraConfig {
			obj[k] = v
		}
	}

	data, err = yaml.Marshal(obj)
	if err != nil {
		return nil, err
	}

	// validate if data can be decoded as KubeletConfiguration
	if err := yaml.UnmarshalStrict(data, &kubeletapi.KubeletConfiguration{}); err != nil {
		return nil, errors.Wrap(err, "validating generated KubeletConfiguration object")
	}

	return data, nil
}

func kvs(kv map[string]string) string {
	var params []string
	for k, v := range kv {
		params = append(params, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(params, ",")
}

func makeCommonKubeletEnvParams(spec *api.ClusterConfig, ng *api.NodeGroup) []string {
	variables := []string{
		fmt.Sprintf("NODE_LABELS=%s", kvs(ng.Labels)),
		fmt.Sprintf("NODE_TAINTS=%s", kvs(ng.Taints)),
	}

	if ng.MaxPodsPerNode != 0 {
		variables = append(variables, fmt.Sprintf("MAX_PODS=%d", ng.MaxPodsPerNode))
	}
	return variables
}

func makeMetadata(spec *api.ClusterConfig) []string {
	return []string{
		fmt.Sprintf("AWS_DEFAULT_REGION=%s", spec.Metadata.Region),
		fmt.Sprintf("AWS_EKS_CLUSTER_NAME=%s", spec.Metadata.Name),
		fmt.Sprintf("AWS_EKS_ENDPOINT=%s", spec.Status.Endpoint),
		fmt.Sprintf("AWS_EKS_ECR_ACCOUNT=%s", api.EKSResourceAccountID(spec.Metadata.Region)),
	}
}

func makeMaxPodsMapping() string {
	var text strings.Builder
	for k, v := range maxPodsPerNodeType {
		text.WriteString(fmt.Sprintf("%s %d\n", k, v))
	}
	return text.String()
}

// NewUserData creates new user data for a given node image family
func NewUserData(spec *api.ClusterConfig, ng *api.NodeGroup) (string, error) {
	switch ng.AMIFamily {
	case ami.ImageFamilyAmazonLinux2:
		return NewUserDataForAmazonLinux2(spec, ng)
	case ami.ImageFamilyUbuntu1804:
		return NewUserDataForUbuntu1804(spec, ng)
	case ami.ImageFamilyWindowsServer2019CoreContainer, ami.ImageFamilyWindowsServer2019FullContainer:
		return newUserDataForWindows(spec, ng)
	default:
		return "", nil
	}
}
