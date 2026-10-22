package v1alpha5_test

import (
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
	}, "kubernetesNetworkConfig.elasticLoadBalancing.gatewayApi is not supported by the AWS SDK version eksctl is built against"),
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
