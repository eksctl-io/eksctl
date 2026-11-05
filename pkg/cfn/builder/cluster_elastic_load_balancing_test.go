package builder_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tidwall/gjson"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

type elbEntry struct {
	autoMode        bool
	elb             *api.ElasticLoadBalancing
	expectedEnabled *bool
	expectedManaged *bool
}

var _ = DescribeTable("Cluster template elasticLoadBalancing", func(e elbEntry) {
	provider := mockprovider.NewMockProvider()
	cfg := api.NewClusterConfig()
	api.SetClusterConfigDefaults(cfg)
	cfg.VPC = vpcConfig()
	cfg.AvailabilityZones = []string{"us-west-2a", "us-west-2b"}
	cfg.KubernetesNetworkConfig = &api.KubernetesNetworkConfig{
		IPFamily:             api.IPV4Family,
		ElasticLoadBalancing: e.elb,
	}
	if e.autoMode {
		cfg.AutoModeConfig = &api.AutoModeConfig{Enabled: api.Enabled()}
	}

	crs := builder.NewClusterResourceSet(provider.EC2(), provider.STS(), provider.Region(), cfg, nil, false)
	Expect(crs.AddAllResources(context.Background())).To(Succeed())
	body, err := crs.RenderJSON()
	Expect(err).NotTo(HaveOccurred())

	elb := gjson.GetBytes(body, "Resources.ControlPlane.Properties.KubernetesNetworkConfig.ElasticLoadBalancing")
	if e.expectedEnabled == nil {
		Expect(elb.Exists()).To(BeFalse(), "unexpected ElasticLoadBalancing: %s", elb.Raw)
		return
	}
	Expect(elb.Get("Enabled").Exists()).To(BeTrue(), "missing ElasticLoadBalancing.Enabled")
	Expect(elb.Get("Enabled").Bool()).To(Equal(*e.expectedEnabled))
	managed := elb.Get("GatewayApi.Crds.Managed")
	if e.expectedManaged == nil {
		Expect(elb.Get("GatewayApi").Exists()).To(BeFalse(), "unexpected GatewayApi: %s", elb.Get("GatewayApi").Raw)
		return
	}
	Expect(managed.Exists()).To(BeTrue(), "missing GatewayApi.Crds.Managed")
	Expect(managed.Bool()).To(Equal(*e.expectedManaged))
},
	Entry("Auto Mode, ELB unset defaults to enabled", elbEntry{
		autoMode:        true,
		expectedEnabled: api.Enabled(),
	}),
	Entry("Auto Mode, ELB explicitly disabled", elbEntry{
		autoMode:        true,
		elb:             &api.ElasticLoadBalancing{Enabled: api.Disabled()},
		expectedEnabled: api.Disabled(),
	}),
	Entry("no Auto Mode, ELB explicitly disabled", elbEntry{
		elb:             &api.ElasticLoadBalancing{Enabled: api.Disabled()},
		expectedEnabled: api.Disabled(),
	}),
	Entry("no Auto Mode, ELB unset is omitted", elbEntry{}),
	Entry("Auto Mode, gatewayApi with managed CRDs", elbEntry{
		autoMode: true,
		elb: &api.ElasticLoadBalancing{
			Enabled:    api.Enabled(),
			GatewayAPI: &api.GatewayAPI{CRDs: &api.GatewayAPICRDs{Managed: api.Enabled()}},
		},
		expectedEnabled: api.Enabled(),
		expectedManaged: api.Enabled(),
	}),
	Entry("Auto Mode, gatewayApi with self-managed CRDs", elbEntry{
		autoMode: true,
		elb: &api.ElasticLoadBalancing{
			Enabled:    api.Enabled(),
			GatewayAPI: &api.GatewayAPI{CRDs: &api.GatewayAPICRDs{Managed: api.Disabled()}},
		},
		expectedEnabled: api.Enabled(),
		expectedManaged: api.Disabled(),
	}),
)
