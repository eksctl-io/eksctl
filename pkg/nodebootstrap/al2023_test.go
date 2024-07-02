package nodebootstrap_test

import (
	"bytes"
	"encoding/base64"
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
	overrideNodegroupSettings func(api.NodePool)
	expectedUserData          string
}

var _ = DescribeTable("Unmanaged AL2023", func(e al2023Entry) {
	cfg, dns := makeDefaultClusterSettings()
	ng := api.NewNodeGroup()
	makeDefaultNPSettings(ng)

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
	Entry("efa enabled", al2023Entry{
		overrideNodegroupSettings: func(np api.NodePool) {
			np.BaseNodeGroup().EFAEnabled = aws.Bool(true)
		},
		expectedUserData: wrapMIMEParts(xTablesLock + efaScript + nodeConfig),
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
	Entry("native AMI && EFA enabled", al2023Entry{
		overrideNodegroupSettings: func(np api.NodePool) {
			np.BaseNodeGroup().EFAEnabled = aws.Bool(true)
		},
		expectedUserData: wrapMIMEParts(xTablesLock + efaCloudhook),
	}),
	Entry("custom AMI", al2023Entry{
		overrideNodegroupSettings: func(np api.NodePool) {
			np.BaseNodeGroup().AMI = "ami-xxxx"
		},
		expectedUserData: wrapMIMEParts(xTablesLock + managedNodeConfig),
	}),
	Entry("custom AMI && EFA enabled", al2023Entry{
		overrideNodegroupSettings: func(np api.NodePool) {
			np.BaseNodeGroup().AMI = "ami-xxxx"
			np.BaseNodeGroup().EFAEnabled = aws.Bool(true)
		},
		expectedUserData: wrapMIMEParts(xTablesLock + efaCloudhook + managedNodeConfig),
	}),
)

type al2023KubeletEntry struct {
	updateNodeGroup    func(*api.NodeGroup)
	expectedNodeConfig nodeadm.NodeConfig
}

func mustToKubeletConfig(kubeletExtraConfig api.InlineDocument) map[string]runtime.RawExtension {
	kubeletConfig, err := nodebootstrap.ToKubeletConfig(kubeletExtraConfig)
	if err != nil {
		Expect(err).NotTo(HaveOccurred())
	}
	return kubeletConfig
}

var _ = DescribeTable("AL2023 kubeletExtraConfig", func(e al2023KubeletEntry) {
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
	foundNodeConfig := false
	for {
		part, err := reader.NextPart()
		if errors.Is(err, io.EOF) {
			break
		}
		Expect(err).NotTo(HaveOccurred())
		if part.Header.Get("Content-Type") != "application/node.eks.aws" {
			continue
		}
		foundNodeConfig = true
		var nodeConfigBuf bytes.Buffer
		_, err = io.Copy(&nodeConfigBuf, part)
		Expect(err).NotTo(HaveOccurred())
		var nodeConfig nodeadm.NodeConfig
		Expect(yaml.Unmarshal(nodeConfigBuf.Bytes(), &nodeConfig)).To(Succeed())
		Expect(nodeConfig).To(Equal(e.expectedNodeConfig))
	}
	Expect(foundNodeConfig).To(BeTrue(), "expected to find NodeConfig in user data")
},
	Entry("nodegroup with maxPods and taints", al2023KubeletEntry{
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

		expectedNodeConfig: nodeadm.NodeConfig{
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
						"maxPods":    "11",
					}),
					Flags: []string{
						"--node-labels=alpha.eksctl.io/nodegroup-name=al2023-mng-test",
						"--register-with-taints=special=true:NoSchedule",
					},
				},
			},
		},
	}),

	Entry("nodegroup with maxPods, taints and kubeletExtraConfig", al2023KubeletEntry{
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

		expectedNodeConfig: nodeadm.NodeConfig{
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
						"maxPods":             "11",
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

	efaCloudhook = fmt.Sprintf(`--//
Content-Type: text/cloud-boothook
Content-Type: charset="us-ascii"

%s
`, assets.EfaManagedAL2023Boothook)

	efaScript = fmt.Sprintf(`--//
Content-Type: text/x-shellscript
Content-Type: charset="us-ascii"

%s
`, assets.EfaAl2023Sh)

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
