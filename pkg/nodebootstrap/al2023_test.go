package nodebootstrap_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	nodeadmapi "github.com/awslabs/amazon-eks-ami/nodeadm/api"
	nodeadm "github.com/awslabs/amazon-eks-ami/nodeadm/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap/assets"
)

type al2023Entry struct {
	overrideClusterSettings   func(*api.ClusterConfig)
	overrideNodegroupSettings func(api.NodePool)
	expectedUserData          string
}

var _ = DescribeTable("Unmanaged AL2023", func(e al2023Entry) {
	cfg, dns := makeDefaultClusterSettings()
	ng := api.NewNodeGroup()
	makeDefaultNPSettings(ng)

	if e.overrideClusterSettings != nil {
		e.overrideClusterSettings(cfg)
	}

	if e.overrideNodegroupSettings != nil {
		e.overrideNodegroupSettings(ng)
	}

	al2023BS := nodebootstrap.NewAL2023Bootstrapper(cfg, ng, dns)
	al2023BS.UserDataMimeBoundary = "//"

	userData, err := al2023BS.UserData()
	Expect(err).NotTo(HaveOccurred())
	decoded, err := base64.StdEncoding.DecodeString(userData)
	Expect(err).NotTo(HaveOccurred())
	actual := strings.ReplaceAll(string(decoded), "\r\n", "\n")
	Expect(actual).To(Equal(e.expectedUserData))
},
	Entry("default", al2023Entry{
		expectedUserData: wrapMIMEParts(xTablesLock + nodeConfig),
	}),
	Entry("ipv6", al2023Entry{
		overrideClusterSettings: func(cc *api.ClusterConfig) {
			cc.Status.KubernetesNetworkConfig.IPFamily = api.IPV6Family
			cc.Status.KubernetesNetworkConfig.ServiceIPv6CIDR = "fd00:facc:76a1::/108"
			cc.Status.KubernetesNetworkConfig.ServiceIPv4CIDR = ""
		},
		expectedUserData: wrapMIMEParts(xTablesLock + nodeConfigIPv6),
	}),
	Entry("efa enabled", al2023Entry{
		overrideNodegroupSettings: func(np api.NodePool) {
			np.BaseNodeGroup().EFAEnabled = aws.Bool(true)
		},
		expectedUserData: wrapMIMEParts(xTablesLock + nodeConfig),
	}),
)

var _ = DescribeTable("Managed AL2023", func(e al2023Entry) {
	cfg, dns := makeDefaultClusterSettings()
	mng := api.NewManagedNodeGroup()
	makeDefaultNPSettings(mng)
	mng.Taints = append(mng.Taints, api.NodeGroupTaint{
		Key:    "special",
		Value:  "true",
		Effect: "NoSchedule",
	})

	if e.overrideNodegroupSettings != nil {
		e.overrideNodegroupSettings(mng)
	}

	al2023ManagedBS := nodebootstrap.NewManagedAL2023Bootstrapper(cfg, mng, dns)
	al2023ManagedBS.UserDataMimeBoundary = "//"

	userData, err := al2023ManagedBS.UserData()
	Expect(err).NotTo(HaveOccurred())
	decoded, err := base64.StdEncoding.DecodeString(userData)
	Expect(err).NotTo(HaveOccurred())
	actual := strings.ReplaceAll(string(decoded), "\r\n", "\n")
	Expect(actual).To(Equal(e.expectedUserData))
},
	Entry("native AMI", al2023Entry{
		expectedUserData: wrapMIMEParts(xTablesLock),
	}),
	Entry("custom AMI", al2023Entry{
		overrideNodegroupSettings: func(np api.NodePool) {
			np.BaseNodeGroup().AMI = "ami-xxxx"
		},
		expectedUserData: wrapMIMEParts(xTablesLock + managedNodeConfig),
	}),
)

type al2023OverrideNodeConfigEntry struct {
	updateNodeGroup     func(*api.NodeGroup)
	expectedNodeConfigs []nodeadm.NodeConfig
}

func mustToKubeletConfig(kubeletExtraConfig api.InlineDocument) map[string]runtime.RawExtension {
	kubeletConfig, err := nodebootstrap.ToKubeletConfig(kubeletExtraConfig)
	if err != nil {
		Expect(err).NotTo(HaveOccurred())
	}
	return kubeletConfig
}

var _ = DescribeTable("AL2023 override node config", func(e al2023OverrideNodeConfigEntry) {
	cfg, dns := makeDefaultClusterSettings()
	ng := api.NewNodeGroup()
	if e.updateNodeGroup != nil {
		e.updateNodeGroup(ng)
	}

	al2023BS := nodebootstrap.NewAL2023Bootstrapper(cfg, ng, dns)
	al2023BS.UserDataMimeBoundary = "//"

	userData, err := al2023BS.UserData()
	Expect(err).NotTo(HaveOccurred())
	decoded, err := base64.StdEncoding.DecodeString(userData)
	Expect(err).NotTo(HaveOccurred())
	reader := multipart.NewReader(bytes.NewReader(decoded), al2023BS.UserDataMimeBoundary)
	nodeConfigCounter := 0
	for {
		part, err := reader.NextPart()
		if errors.Is(err, io.EOF) {
			break
		}
		Expect(err).NotTo(HaveOccurred())
		if part.Header.Get("Content-Type") != "application/node.eks.aws" {
			continue
		}
		var nodeConfigBuf bytes.Buffer
		_, err = io.Copy(&nodeConfigBuf, part)
		Expect(err).NotTo(HaveOccurred())
		var nodeConfig nodeadm.NodeConfig
		Expect(yaml.Unmarshal(nodeConfigBuf.Bytes(), &nodeConfig)).To(Succeed())
		Expect(nodeConfigCounter).To(BeNumerically("<", len(e.expectedNodeConfigs)))
		Expect(nodeConfig).To(Equal(e.expectedNodeConfigs[nodeConfigCounter]))
		nodeConfigCounter++
	}
	Expect(nodeConfigCounter).To(BeNumerically("==", len(e.expectedNodeConfigs)))
},
	Entry("nodegroup with maxPods and taints", al2023OverrideNodeConfigEntry{
		updateNodeGroup: func(ng *api.NodeGroup) {
			ng.MaxPodsPerNode = 11
			ng.Labels = map[string]string{"alpha.eksctl.io/nodegroup-name": "al2023-mng-test"}
			ng.Taints = []api.NodeGroupTaint{
				{
					Key:    "special",
					Value:  "true",
					Effect: "NoSchedule",
				},
			}
		},

		expectedNodeConfigs: []nodeadm.NodeConfig{
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       nodeadmapi.KindNodeConfig,
					APIVersion: nodeadm.GroupVersion.String(),
				},
				Spec: nodeadm.NodeConfigSpec{
					Cluster: nodeadm.ClusterDetails{
						APIServerEndpoint:    "https://test.xxx.us-west-2.eks.amazonaws.com",
						CertificateAuthority: []byte("test CA"),
						CIDR:                 "10.100.0.0/16",
						Name:                 "al2023-test",
					},
					Kubelet: nodeadm.KubeletOptions{
						Config: mustToKubeletConfig(map[string]interface{}{
							"clusterDNS": []string{"10.100.0.10"},
						}),
						Flags: []string{
							"--node-labels=alpha.eksctl.io/nodegroup-name=al2023-mng-test",
							"--register-with-taints=special=true:NoSchedule",
						},
					},
				},
			},
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       nodeadmapi.KindNodeConfig,
					APIVersion: nodeadm.GroupVersion.String(),
				},
				Spec: nodeadm.NodeConfigSpec{
					Kubelet: nodeadm.KubeletOptions{
						Config: mustToKubeletConfig(map[string]interface{}{
							"maxPods": 11,
						}),
					},
				},
			},
		},
	}),

	Entry("nodegroup with maxPods, taints and kubeletExtraConfig", al2023OverrideNodeConfigEntry{
		updateNodeGroup: func(ng *api.NodeGroup) {
			ng.MaxPodsPerNode = 11
			ng.Labels = map[string]string{"alpha.eksctl.io/nodegroup-name": "al2023-mng-test"}
			ng.Taints = []api.NodeGroupTaint{
				{
					Key:    "special",
					Value:  "true",
					Effect: "NoSchedule",
				},
			}
			ng.KubeletExtraConfig = &api.InlineDocument{
				"shutdownGracePeriod": "5m",
				"kubeReserved": map[string]interface{}{
					"cpu":    "500m",
					"memory": "250Mi",
				},
			}
		},

		expectedNodeConfigs: []nodeadm.NodeConfig{
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       nodeadmapi.KindNodeConfig,
					APIVersion: nodeadm.GroupVersion.String(),
				},
				Spec: nodeadm.NodeConfigSpec{
					Cluster: nodeadm.ClusterDetails{
						APIServerEndpoint:    "https://test.xxx.us-west-2.eks.amazonaws.com",
						CertificateAuthority: []byte("test CA"),
						CIDR:                 "10.100.0.0/16",
						Name:                 "al2023-test",
					},
					Kubelet: nodeadm.KubeletOptions{
						Config: mustToKubeletConfig(map[string]interface{}{
							"clusterDNS":          []string{"10.100.0.10"},
							"shutdownGracePeriod": "5m",
							"kubeReserved": map[string]interface{}{
								"cpu":    "500m",
								"memory": "250Mi",
							},
						}),
						Flags: []string{
							"--node-labels=alpha.eksctl.io/nodegroup-name=al2023-mng-test",
							"--register-with-taints=special=true:NoSchedule",
						},
					},
				},
			},
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       nodeadmapi.KindNodeConfig,
					APIVersion: nodeadm.GroupVersion.String(),
				},
				Spec: nodeadm.NodeConfigSpec{
					Kubelet: nodeadm.KubeletOptions{
						Config: mustToKubeletConfig(map[string]interface{}{
							"maxPods": 11,
						}),
					},
				},
			},
		},
	}),

	Entry("nodegroup with overrideBootstrapCommand", al2023OverrideNodeConfigEntry{
		updateNodeGroup: func(ng *api.NodeGroup) {
			nodeConfig := nodeadm.NodeConfig{
				TypeMeta: metav1.TypeMeta{
					Kind:       nodeadmapi.KindNodeConfig,
					APIVersion: nodeadm.GroupVersion.String(),
				},
				Spec: nodeadm.NodeConfigSpec{
					Instance: nodeadm.InstanceOptions{
						LocalStorage: nodeadm.LocalStorageOptions{
							Strategy: nodeadm.LocalStorageRAID0,
						},
					},
					Kubelet: nodeadm.KubeletOptions{
						Config: mustToKubeletConfig(map[string]interface{}{
							"shutdownGracePeriod": "5m",
							"featureGates": map[string]interface{}{
								"DisableKubeletCloudCredentialProviders": true,
							},
						}),
					},
				},
			}
			jsonNodeConfig, err := json.Marshal(nodeConfig)
			Expect(err).NotTo(HaveOccurred())
			ng.OverrideBootstrapCommand = aws.String(string(jsonNodeConfig))
			ng.Labels = map[string]string{"alpha.eksctl.io/nodegroup-name": "al2023-mng-test"}

		},
		expectedNodeConfigs: []nodeadm.NodeConfig{
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       nodeadmapi.KindNodeConfig,
					APIVersion: nodeadm.GroupVersion.String(),
				},
				Spec: nodeadm.NodeConfigSpec{
					Cluster: nodeadm.ClusterDetails{
						APIServerEndpoint:    "https://test.xxx.us-west-2.eks.amazonaws.com",
						CertificateAuthority: []byte("test CA"),
						CIDR:                 "10.100.0.0/16",
						Name:                 "al2023-test",
					},
					Kubelet: nodeadm.KubeletOptions{
						Config: mustToKubeletConfig(map[string]interface{}{
							"clusterDNS": []string{"10.100.0.10"},
						}),
						Flags: []string{
							"--node-labels=alpha.eksctl.io/nodegroup-name=al2023-mng-test",
						},
					},
				},
			},
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       nodeadmapi.KindNodeConfig,
					APIVersion: nodeadm.GroupVersion.String(),
				},
				Spec: nodeadm.NodeConfigSpec{
					Instance: nodeadm.InstanceOptions{
						LocalStorage: nodeadm.LocalStorageOptions{
							Strategy: nodeadm.LocalStorageRAID0,
						},
					},
					Kubelet: nodeadm.KubeletOptions{
						Config: mustToKubeletConfig(map[string]interface{}{
							"shutdownGracePeriod": "5m",
							"featureGates": map[string]interface{}{
								"DisableKubeletCloudCredentialProviders": true,
							},
						}),
					},
				},
			},
		},
	}),
)

var (
	makeDefaultClusterSettings = func() (*api.ClusterConfig, string) {
		clusterConfig := api.NewClusterConfig()
		clusterConfig.Metadata.Name = "al2023-test"
		clusterConfig.Status = &api.ClusterStatus{
			Endpoint:                 "https://test.xxx.us-west-2.eks.amazonaws.com",
			CertificateAuthorityData: []byte("test CA"),
			KubernetesNetworkConfig: &api.KubernetesNetworkConfig{
				ServiceIPv4CIDR: "10.100.0.0/16",
			},
		}
		return clusterConfig, "10.100.0.10"
	}

	makeDefaultNPSettings = func(np api.NodePool) {
		baseNg := np.BaseNodeGroup()
		baseNg.Labels = map[string]string{
			"alpha.eksctl.io/nodegroup-name": "al2023-mng-test",
		}
	}

	wrapMIMEParts = func(parts string) string {
		return `MIME-Version: 1.0
Content-Type: multipart/mixed; boundary=//

` + parts + `--//--
`
	}

	xTablesLock = fmt.Sprintf(`--//
Content-Type: text/x-shellscript
Content-Type: charset="us-ascii"

%s
`, assets.AL2023XTablesLock)

	nodeConfig = `--//
Content-Type: application/node.eks.aws

apiVersion: node.eks.aws/v1alpha1
kind: NodeConfig
metadata:
  creationTimestamp: null
spec:
  cluster:
    apiServerEndpoint: https://test.xxx.us-west-2.eks.amazonaws.com
    certificateAuthority: dGVzdCBDQQ==
    cidr: 10.100.0.0/16
    name: al2023-test
  containerd: {}
  instance:
    localStorage: {}
  kubelet:
    config:
      clusterDNS:
      - 10.100.0.10
    flags:
    - --node-labels=alpha.eksctl.io/nodegroup-name=al2023-mng-test

`
	nodeConfigIPv6 = `--//
Content-Type: application/node.eks.aws

apiVersion: node.eks.aws/v1alpha1
kind: NodeConfig
metadata:
  creationTimestamp: null
spec:
  cluster:
    apiServerEndpoint: https://test.xxx.us-west-2.eks.amazonaws.com
    certificateAuthority: dGVzdCBDQQ==
    cidr: fd00:facc:76a1::/108
    name: al2023-test
  containerd: {}
  instance:
    localStorage: {}
  kubelet:
    config:
      clusterDNS:
      - 10.100.0.10
    flags:
    - --node-labels=alpha.eksctl.io/nodegroup-name=al2023-mng-test

`
	managedNodeConfig = `--//
Content-Type: application/node.eks.aws

apiVersion: node.eks.aws/v1alpha1
kind: NodeConfig
metadata:
  creationTimestamp: null
spec:
  cluster:
    apiServerEndpoint: https://test.xxx.us-west-2.eks.amazonaws.com
    certificateAuthority: dGVzdCBDQQ==
    cidr: 10.100.0.0/16
    name: al2023-test
  containerd: {}
  instance:
    localStorage: {}
  kubelet:
    config:
      clusterDNS:
      - 10.100.0.10
    flags:
    - --node-labels=alpha.eksctl.io/nodegroup-name=al2023-mng-test
    - --register-with-taints=special=true:NoSchedule

`
)
