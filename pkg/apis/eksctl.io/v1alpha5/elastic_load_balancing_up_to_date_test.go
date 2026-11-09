package v1alpha5_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go-v2/aws"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var _ = Describe("ElasticLoadBalancingUpToDate", func() {
	gatewayAPI := func(managed bool) *ekstypes.GatewayApiConfiguration {
		return &ekstypes.GatewayApiConfiguration{
			Crds: &ekstypes.CrdManagementConfiguration{Managed: aws.Bool(managed)},
		}
	}

	type entry struct {
		current   *ekstypes.ElasticLoadBalancing
		requested *ekstypes.ElasticLoadBalancing
		expected  bool
	}

	DescribeTable("compares the current and requested config", func(e entry) {
		Expect(api.ElasticLoadBalancingUpToDate(e.current, e.requested)).To(Equal(e.expected))
	},
		Entry("current is unset", entry{
			requested: &ekstypes.ElasticLoadBalancing{Enabled: aws.Bool(true)},
			expected:  false,
		}),
		Entry("enabled matches and gatewayApi is not requested", entry{
			current:   &ekstypes.ElasticLoadBalancing{Enabled: aws.Bool(true), GatewayApi: gatewayAPI(true)},
			requested: &ekstypes.ElasticLoadBalancing{Enabled: aws.Bool(true)},
			expected:  true,
		}),
		Entry("enabled differs", entry{
			current:   &ekstypes.ElasticLoadBalancing{Enabled: aws.Bool(false)},
			requested: &ekstypes.ElasticLoadBalancing{Enabled: aws.Bool(true)},
			expected:  false,
		}),
		Entry("self-managed CRDs requested on a cluster without gatewayApi", entry{
			current:   &ekstypes.ElasticLoadBalancing{Enabled: aws.Bool(true)},
			requested: &ekstypes.ElasticLoadBalancing{Enabled: aws.Bool(true), GatewayApi: gatewayAPI(false)},
			expected:  false,
		}),
		Entry("managed CRDs requested on a cluster without gatewayApi", entry{
			current:   &ekstypes.ElasticLoadBalancing{Enabled: aws.Bool(true)},
			requested: &ekstypes.ElasticLoadBalancing{Enabled: aws.Bool(true), GatewayApi: gatewayAPI(true)},
			expected:  false,
		}),
		Entry("managed differs", entry{
			current:   &ekstypes.ElasticLoadBalancing{Enabled: aws.Bool(true), GatewayApi: gatewayAPI(true)},
			requested: &ekstypes.ElasticLoadBalancing{Enabled: aws.Bool(true), GatewayApi: gatewayAPI(false)},
			expected:  false,
		}),
		Entry("managed matches", entry{
			current:   &ekstypes.ElasticLoadBalancing{Enabled: aws.Bool(true), GatewayApi: gatewayAPI(false)},
			requested: &ekstypes.ElasticLoadBalancing{Enabled: aws.Bool(true), GatewayApi: gatewayAPI(false)},
			expected:  true,
		}),
	)
})
