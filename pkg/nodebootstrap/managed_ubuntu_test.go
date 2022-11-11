package nodebootstrap_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
)

var _ = Describe("Managed Ubuntu User Data", func() {
	type bootScriptEntry struct {
		clusterStatus    *api.ClusterStatus
		expectedUserData string
	}

	DescribeTable("cluster DNS in userdata", func(be bootScriptEntry) {
		clusterConfig := api.NewClusterConfig()
		clusterConfig.Metadata.Name = "cluster"
		clusterConfig.Status = be.clusterStatus
		ng := &api.ManagedNodeGroup{
			NodeGroupBase: &api.NodeGroupBase{
				AMIFamily: api.NodeImageFamilyUbuntu2004,
			},
		}

		bootstrapper, err := nodebootstrap.NewManagedBootstrapper(clusterConfig, ng)
		Expect(err).NotTo(HaveOccurred())
		userData, err := bootstrapper.UserData()
		Expect(err).NotTo(HaveOccurred())
		cloudCfg := decode(userData)
		Expect(cloudCfg.WriteFiles[1].Path).To(Equal("/etc/eksctl/kubelet.env"))
		contentLines := strings.Split(cloudCfg.WriteFiles[1].Content, "\n")
		Expect(contentLines).To(ConsistOf(strings.Split(be.expectedUserData, "\n")))
	},
		Entry("custom Kubernetes serviceIPV4CIDR", bootScriptEntry{
			clusterStatus: &api.ClusterStatus{
				KubernetesNetworkConfig: &api.KubernetesNetworkConfig{
					ServiceIPv4CIDR: "10.255.0.0/16",
				},
			},
			expectedUserData: `CLUSTER_NAME=cluster
API_SERVER_URL=
B64_CLUSTER_CA=
NODE_LABELS=
NODE_TAINTS=
CLUSTER_DNS=10.255.0.10`,
		}),
		Entry("default Kubernetes serviceIPV4CIDR", bootScriptEntry{
			clusterStatus: &api.ClusterStatus{
				KubernetesNetworkConfig: &api.KubernetesNetworkConfig{
					ServiceIPv4CIDR: "172.20.0.0/16",
				},
			},
			expectedUserData: `CLUSTER_NAME=cluster
API_SERVER_URL=
B64_CLUSTER_CA=
NODE_LABELS=
NODE_TAINTS=
CLUSTER_DNS=172.20.0.10`,
		}),
	)
})
