package v1alpha5_test

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var _ = DescribeTable("Elastic Load Balancing Validation", func(c *api.ClusterConfig, expectedErr string) {
	err := api.ValidateElasticLoadBalancing(c)
	if expectedErr != "" {
		Expect(err).To(MatchError(expectedErr))
	} else {
		Expect(err).NotTo(HaveOccurred())
	}
},
	Entry("no kubernetesNetworkConfig", &api.ClusterConfig{}, ""),
	Entry("no elasticLoadBalancing", &api.ClusterConfig{
		KubernetesNetworkConfig: &api.KubernetesNetworkConfig{},
	}, ""),
	Entry("elasticLoadBalancing enabled without gatewayApi", &api.ClusterConfig{
		KubernetesNetworkConfig: &api.KubernetesNetworkConfig{
			ElasticLoadBalancing: &api.ElasticLoadBalancing{
				Enabled: api.Enabled(),
			},
		},
	}, ""),
	Entry("gatewayApi when the capability is disabled", &api.ClusterConfig{
		KubernetesNetworkConfig: &api.KubernetesNetworkConfig{
			ElasticLoadBalancing: &api.ElasticLoadBalancing{
				Enabled: api.Disabled(),
				GatewayAPI: &api.GatewayAPI{
					CRDs: &api.GatewayAPICRDs{Managed: api.Enabled()},
				},
			},
		},
	}, "cannot set kubernetesNetworkConfig.elasticLoadBalancing.gatewayApi when the load balancing capability is disabled"),
	Entry("gatewayApi without crds.managed", &api.ClusterConfig{
		KubernetesNetworkConfig: &api.KubernetesNetworkConfig{
			ElasticLoadBalancing: &api.ElasticLoadBalancing{
				Enabled:    api.Enabled(),
				GatewayAPI: &api.GatewayAPI{},
			},
		},
	}, "kubernetesNetworkConfig.elasticLoadBalancing.gatewayApi.crds.managed must be set to either true or false"),
	Entry("gatewayApi with crds.managed", &api.ClusterConfig{
		KubernetesNetworkConfig: &api.KubernetesNetworkConfig{
			ElasticLoadBalancing: &api.ElasticLoadBalancing{
				Enabled: api.Enabled(),
				GatewayAPI: &api.GatewayAPI{
					CRDs: &api.GatewayAPICRDs{Managed: api.Enabled()},
				},
			},
		},
	}, ""),
)

var _ = DescribeTable("Elastic Load Balancing matches Auto Mode", func(c *api.ClusterConfig, expectedErr string) {
	err := api.ValidateElasticLoadBalancingMatchesAutoMode(c)
	if expectedErr != "" {
		Expect(err).To(MatchError(expectedErr))
	} else {
		Expect(err).NotTo(HaveOccurred())
	}
},
	Entry("both unset", &api.ClusterConfig{}, ""),
	Entry("unset with Auto Mode enabled", &api.ClusterConfig{
		AutoModeConfig: &api.AutoModeConfig{Enabled: api.Enabled()},
	}, ""),
	Entry("enabled with Auto Mode enabled", &api.ClusterConfig{
		AutoModeConfig: &api.AutoModeConfig{Enabled: api.Enabled()},
		KubernetesNetworkConfig: &api.KubernetesNetworkConfig{
			ElasticLoadBalancing: &api.ElasticLoadBalancing{Enabled: api.Enabled()},
		},
	}, ""),
	Entry("disabled with Auto Mode enabled", &api.ClusterConfig{
		AutoModeConfig: &api.AutoModeConfig{Enabled: api.Enabled()},
		KubernetesNetworkConfig: &api.KubernetesNetworkConfig{
			ElasticLoadBalancing: &api.ElasticLoadBalancing{Enabled: api.Disabled()},
		},
	}, "kubernetesNetworkConfig.elasticLoadBalancing.enabled must match autoModeConfig.enabled"),
	Entry("enabled without Auto Mode", &api.ClusterConfig{
		KubernetesNetworkConfig: &api.KubernetesNetworkConfig{
			ElasticLoadBalancing: &api.ElasticLoadBalancing{Enabled: api.Enabled()},
		},
	}, "kubernetesNetworkConfig.elasticLoadBalancing.enabled must match autoModeConfig.enabled"),
)

var _ = DescribeTable("DescribeElasticLoadBalancing", func(elb *ekstypes.ElasticLoadBalancing, expected string) {
	Expect(api.DescribeElasticLoadBalancing(elb)).To(Equal(expected))
},
	Entry("disabled", &ekstypes.ElasticLoadBalancing{Enabled: aws.Bool(false)}, "disabled"),
	Entry("enabled without gatewayApi", &ekstypes.ElasticLoadBalancing{Enabled: aws.Bool(true)}, "enabled"),
	Entry("managed CRDs", &ekstypes.ElasticLoadBalancing{
		Enabled:    aws.Bool(true),
		GatewayApi: &ekstypes.GatewayApiConfiguration{Crds: &ekstypes.CrdManagementConfiguration{Managed: aws.Bool(true)}},
	}, "enabled, Gateway API CRDs managed by EKS"),
	Entry("self-managed CRDs", &ekstypes.ElasticLoadBalancing{
		Enabled:    aws.Bool(true),
		GatewayApi: &ekstypes.GatewayApiConfiguration{Crds: &ekstypes.CrdManagementConfiguration{Managed: aws.Bool(false)}},
	}, "enabled, Gateway API CRDs self-managed"),
)

var _ = DescribeTable("ElasticLoadBalancingEnabled", func(c *api.ClusterConfig, expected bool) {
	Expect(c.ElasticLoadBalancingEnabled()).To(Equal(expected))
},
	Entry("no config", &api.ClusterConfig{}, false),
	Entry("defaults to Auto Mode when unset", &api.ClusterConfig{
		AutoModeConfig: &api.AutoModeConfig{Enabled: api.Enabled()},
	}, true),
	Entry("explicitly disabled with Auto Mode enabled", &api.ClusterConfig{
		AutoModeConfig: &api.AutoModeConfig{Enabled: api.Enabled()},
		KubernetesNetworkConfig: &api.KubernetesNetworkConfig{
			ElasticLoadBalancing: &api.ElasticLoadBalancing{Enabled: api.Disabled()},
		},
	}, false),
	Entry("explicitly enabled without Auto Mode", &api.ClusterConfig{
		KubernetesNetworkConfig: &api.KubernetesNetworkConfig{
			ElasticLoadBalancing: &api.ElasticLoadBalancing{Enabled: api.Enabled()},
		},
	}, true),
)

var _ = DescribeTable("ElasticLoadBalancingRequest", func(c *api.ClusterConfig, expected *ekstypes.ElasticLoadBalancing) {
	Expect(c.ElasticLoadBalancingRequest()).To(Equal(expected))
},
	Entry("defaults to Auto Mode", &api.ClusterConfig{
		AutoModeConfig: &api.AutoModeConfig{Enabled: api.Enabled()},
	}, &ekstypes.ElasticLoadBalancing{Enabled: aws.Bool(true)}),
	Entry("omits gatewayApi when unset", &api.ClusterConfig{
		KubernetesNetworkConfig: &api.KubernetesNetworkConfig{
			ElasticLoadBalancing: &api.ElasticLoadBalancing{Enabled: api.Disabled()},
		},
	}, &ekstypes.ElasticLoadBalancing{Enabled: aws.Bool(false)}),
	Entry("includes gatewayApi.crds.managed", &api.ClusterConfig{
		KubernetesNetworkConfig: &api.KubernetesNetworkConfig{
			ElasticLoadBalancing: &api.ElasticLoadBalancing{
				Enabled: api.Enabled(),
				GatewayAPI: &api.GatewayAPI{
					CRDs: &api.GatewayAPICRDs{Managed: api.Enabled()},
				},
			},
		},
	}, &ekstypes.ElasticLoadBalancing{
		Enabled: aws.Bool(true),
		GatewayApi: &ekstypes.GatewayApiConfiguration{
			Crds: &ekstypes.CrdManagementConfiguration{Managed: aws.Bool(true)},
		},
	}),
)
