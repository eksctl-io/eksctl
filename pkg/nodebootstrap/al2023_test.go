package nodebootstrap_test

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
		expectedUserData: wrapMIMEParts(nodeConfig),
	}),
	Entry("efa enabled", al2023Entry{
		overrideNodegroupSettings: func(np api.NodePool) {
			np.BaseNodeGroup().EFAEnabled = aws.Bool(true)
		},
		expectedUserData: wrapMIMEParts(efaScript + nodeConfig),
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
		expectedUserData: "",
	}),
	Entry("native AMI && EFA enabled", al2023Entry{
		overrideNodegroupSettings: func(np api.NodePool) {
			np.BaseNodeGroup().EFAEnabled = aws.Bool(true)
		},
		expectedUserData: wrapMIMEParts(efaCloudhook),
	}),
	Entry("custom AMI", al2023Entry{
		overrideNodegroupSettings: func(np api.NodePool) {
			np.BaseNodeGroup().AMI = "ami-xxxx"
		},
		expectedUserData: wrapMIMEParts(managedNodeConfig),
	}),
	Entry("custom AMI && EFA enabled", al2023Entry{
		overrideNodegroupSettings: func(np api.NodePool) {
			np.BaseNodeGroup().AMI = "ami-xxxx"
			np.BaseNodeGroup().EFAEnabled = aws.Bool(true)
		},
		expectedUserData: wrapMIMEParts(efaCloudhook + managedNodeConfig),
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
		baseNg.MaxPodsPerNode = 4
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
spec:
  cluster:
    apiServerEndpoint: https://test.xxx.us-west-2.eks.amazonaws.com
    certificateAuthority: dGVzdCBDQQ==
    cidr: 10.100.0.0/16
    name: al2023-test
  kubelet:
    config:
      maxPods: 4
      clusterDNS:
      - 10.100.0.10
    flags:
    - --node-labels=alpha.eksctl.io/nodegroup-name=al2023-mng-test

`
	managedNodeConfig = `--//
Content-Type: application/node.eks.aws

apiVersion: node.eks.aws/v1alpha1
kind: NodeConfig
spec:
  cluster:
    apiServerEndpoint: https://test.xxx.us-west-2.eks.amazonaws.com
    certificateAuthority: dGVzdCBDQQ==
    cidr: 10.100.0.0/16
    name: al2023-test
  kubelet:
    config:
      maxPods: 4
      clusterDNS:
      - 10.100.0.10
    flags:
    - --node-labels=alpha.eksctl.io/nodegroup-name=al2023-mng-test
    - --register-with-taints=special=true:NoSchedule

`
)
