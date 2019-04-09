package nodebootstrap

import (
	"fmt"
	"os"
	"strings"

	gfn "github.com/awslabs/goformation/cloudformation"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	"k8s.io/client-go/tools/clientcmd"
	kubeletapi "k8s.io/kubelet/config/v1beta1"

	"sigs.k8s.io/yaml"

	"github.com/weaveworks/eksctl/pkg/ami"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/cloudconfig"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
)

//go:generate ${GOPATH}/bin/go-bindata -pkg ${GOPACKAGE} -prefix assets -modtime 1 -o assets.go assets
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

func makeClientConfigData(spec *api.ClusterConfig, ng *api.NodeGroup) ([]byte, error) {
	clientConfig, _, _ := kubeconfig.New(spec, "kubelet", configDir+"ca.crt")
	authenticator := kubeconfig.AWSIAMAuthenticator
	if ng.AMIFamily == ami.ImageFamilyUbuntu1804 {
		authenticator = kubeconfig.HeptioAuthenticatorAWS
	}
	kubeconfig.AppendAuthenticator(clientConfig, spec, authenticator, "")
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
	if err := yaml.Unmarshal(data, &obj); err != nil {
		return nil, err
	}

	if ng.MaxPodsPerNode == 0 {
		ng.MaxPodsPerNode = maxPodsPerNodeType[ng.InstanceType]
	}
	obj["maxPods"] = int32(ng.MaxPodsPerNode)

	obj["clusterDNS"] = []string{
		clusterDNS(spec, ng),
	}

	data, err = yaml.Marshal(obj)
	if err != nil {
		return nil, err
	}

	// validate if data can be decoded as KubeletConfiguration
	if err := yaml.Unmarshal(data, &kubeletapi.KubeletConfiguration{}); err != nil {
		return nil, errors.Wrap(err, "validating generated KubeletConfiguration object")
	}

	return data, nil
}

func makeCommonKubeletEnvParams(spec *api.ClusterConfig, ng *api.NodeGroup) []string {
	kvs := func(kv map[string]string) string {
		var params []string
		for k, v := range kv {
			params = append(params, fmt.Sprintf("%s=%s", k, v))
		}
		return strings.Join(params, ",")
	}

	return []string{
		fmt.Sprintf("NODE_LABELS=%s", kvs(ng.Labels)),
		fmt.Sprintf("NODE_TAINTS=%s", kvs(ng.Taints)),
	}
}

func makeMetadata(spec *api.ClusterConfig) []string {
	return []string{
		fmt.Sprintf("AWS_DEFAULT_REGION=%s", spec.Metadata.Region),
		fmt.Sprintf("AWS_EKS_CLUSTER_NAME=%s", spec.Metadata.Name),
		fmt.Sprintf("AWS_EKS_ENDPOINT=%s", spec.Status.Endpoint),
	}
}

// NewUserData creates new user data for a given node image family
func NewUserData(spec *api.ClusterConfig, ng *api.NodeGroup) (*gfn.Value, error) {
	if ng.OverrideUserData != nil {
		return NewOverrideUserData(*ng.OverrideUserData), nil
	}

	var userData string
	var err error
	switch ng.AMIFamily {
	case ami.ImageFamilyAmazonLinux2:
		userData, err = NewUserDataForAmazonLinux2(spec, ng)
	case ami.ImageFamilyUbuntu1804:
		userData, err = NewUserDataForUbuntu1804(spec, ng)
	default:
		err = errors.Errorf("unsupported AMI family '%s'", ng.AMIFamily)
		logger.Critical("can't get user data", err)
	}

	return gfn.NewString(userData), err
}

// NewOverrideUserData creates new user data from NodeGroup.OverrideUserData
func NewOverrideUserData(userData string) *gfn.Value {
	return gfn.MakeIntrinsic(
		gfn.FnBase64,
		gfn.MakeFnSub(
			gfn.NewString(userData),
		),
	)
}

// DefaultOverrideUserData provides canned default override user data
func DefaultOverrideUserData(ng *api.NodeGroup) *string {
	var bytes []byte
	switch ng.AMIFamily {
	case ami.ImageFamilyAmazonLinux2:
		bytes = MustAsset("simple.al2.userdata")
	default:
		err := errors.Errorf("unsupported AMI family '%s'", ng.AMIFamily)
		logger.Critical("can't get default user data override", err)
		panic(err)
	}

	str := string(bytes)
	return &str
}
