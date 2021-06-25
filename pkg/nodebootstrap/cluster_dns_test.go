package nodebootstrap_test

import (
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
)

type clusterDNSEntry struct {
	clusterStatus *api.ClusterStatus

	expectedClusterDNS string
	expectedErr        string
}

var _ = DescribeTable("Cluster DNS", func(c clusterDNSEntry) {
	clusterDNS, err := nodebootstrap.GetClusterDNS(&api.ClusterConfig{Status: c.clusterStatus})
	if c.expectedErr != "" {
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(ContainSubstring(c.expectedErr)))
		return
	}
	Expect(err).ToNot(HaveOccurred())
	Expect(clusterDNS).To(Equal(c.expectedClusterDNS))

},
	Entry("default ServiceIPv4CIDR", clusterDNSEntry{
		clusterStatus: &api.ClusterStatus{
			KubernetesNetworkConfig: &api.KubernetesNetworkConfig{
				ServiceIPv4CIDR: "10.100.0.0/16",
			},
		},
		expectedClusterDNS: "10.100.0.10",
	}),

	Entry("custom ServiceIPv4CIDR", clusterDNSEntry{
		clusterStatus: &api.ClusterStatus{
			KubernetesNetworkConfig: &api.KubernetesNetworkConfig{
				ServiceIPv4CIDR: "172.16.0.0/12",
			},
		},
		expectedClusterDNS: "172.16.0.10",
	}),

	Entry("empty ServiceIPv4CIDR", clusterDNSEntry{
		clusterStatus:      &api.ClusterStatus{},
		expectedClusterDNS: "",
	}),

	Entry("invalid CIDR", clusterDNSEntry{
		clusterStatus: &api.ClusterStatus{
			KubernetesNetworkConfig: &api.KubernetesNetworkConfig{
				ServiceIPv4CIDR: "10.0.0.0/4000",
			},
		},

		expectedErr: "unexpected error",
	}),
)
