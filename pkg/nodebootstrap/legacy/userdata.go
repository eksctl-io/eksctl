package legacy

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap/utils"
	"k8s.io/client-go/tools/clientcmd"
	kubeletapi "k8s.io/kubelet/config/v1beta1"

	"sigs.k8s.io/yaml"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cloudconfig"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap/bindata"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
)

const (
	dataDir              = "bindata/assets"
	configDir            = "/etc/eksctl/"
	kubeletDropInUnitDir = "/etc/systemd/system/kubelet.service.d/"
	dockerConfigDir      = "/etc/docker/"
)

type configFile struct {
	dir      string
	name     string
	contents string
	isAsset  bool
}

func getAsset(name string) (string, error) {
	data, err := bindata.Asset(filepath.Join(dataDir, name))
	if err != nil {
		return "", errors.Wrapf(err, "decoding embedded file %q", name)
	}
	return string(data), nil
}

func addFilesAndScripts(config *cloudconfig.CloudConfig, files []configFile, scripts []string) error {
	for _, file := range files {
		f := cloudconfig.File{
			Path: file.dir + file.name,
		}

		if file.isAsset {
			data, err := getAsset(file.name)
			if err != nil {
				return err
			}
			f.Content = data
		} else {
			f.Content = file.contents
		}

		config.AddFile(f)
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

func makeClientConfigData(spec *api.ClusterConfig, authenticatorCMD string) ([]byte, error) {
	clientConfig := kubeconfig.
		NewBuilder(spec.Metadata, spec.Status, "kubelet").
		UseCertificateAuthorityFile(configDir + "ca.crt").
		Build()
	kubeconfig.AppendAuthenticator(clientConfig, spec.Metadata, authenticatorCMD, "", "")
	clientConfigData, err := clientcmd.Write(*clientConfig)
	if err != nil {
		return nil, errors.Wrap(err, "serialising kubeconfig for nodegroup")
	}
	return clientConfigData, nil
}

func getKubeReserved(info InstanceTypeInfo) api.InlineDocument {
	return api.InlineDocument{
		"ephemeral-storage": info.DefaultStorageToReserve(),
		"cpu":               info.DefaultCPUToReserve(),
		"memory":            info.DefaultMemoryToReserve(),
	}
}

func makeDockerConfigJSON() (string, error) {
	return bindata.AssetString(filepath.Join(dataDir, "docker-daemon.json"))
}

func makeKubeletConfigYAML(spec *api.ClusterConfig, ng *api.NodeGroup) ([]byte, error) {
	data, err := bindata.Asset(filepath.Join(dataDir, "kubelet.yaml"))
	if err != nil {
		return nil, err
	}

	// use a map here, as using struct will require us to add defaulting etc,
	// and we only need to add a few top-level fields
	obj := api.InlineDocument{}
	if err := yaml.UnmarshalStrict(data, &obj); err != nil {
		return nil, err
	}

	if ng.ClusterDNS != "" {
		obj["clusterDNS"] = []string{ng.ClusterDNS}
	}

	// Set default reservations if specs about instance is available
	if info, ok := instanceTypeInfos[ng.InstanceType]; ok {
		// This is a NodeGroup with a single instanceType defined
		if _, ok := obj["kubeReserved"]; !ok {
			obj["kubeReserved"] = api.InlineDocument{}
		}
		obj["kubeReserved"] = getKubeReserved(info)
	} else if ng.InstancesDistribution != nil {
		// This is a NodeGroup using mixed instance types
		var minCPU, minMaxPodsPerNode int64
		for _, instanceType := range ng.InstancesDistribution.InstanceTypes {
			if info, ok := instanceTypeInfos[instanceType]; ok {
				if instanceCPU := info.CPU; minCPU == 0 || instanceCPU < minCPU {
					minCPU = instanceCPU
				}
				if instanceMaxPodsPerNode := info.MaxPodsPerNode; minMaxPodsPerNode == 0 || instanceMaxPodsPerNode < minMaxPodsPerNode {
					minMaxPodsPerNode = instanceMaxPodsPerNode
				}
			}
		}
		if minCPU > 0 && minMaxPodsPerNode > 0 {
			info = InstanceTypeInfo{
				MaxPodsPerNode: minMaxPodsPerNode,
				CPU:            minCPU,
			}
			if _, ok := obj["kubeReserved"]; !ok {
				obj["kubeReserved"] = api.InlineDocument{}
			}
			obj["kubeReserved"] = getKubeReserved(info)
		}
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

func makeCommonKubeletEnvParams(ng *api.NodeGroup) []string {
	variables := []string{
		fmt.Sprintf("NODE_LABELS=%s", kvs(ng.Labels)),
		fmt.Sprintf("NODE_TAINTS=%s", utils.FormatTaints(ng.Taints)),
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
