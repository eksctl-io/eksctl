package v1alpha5_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var _ = DescribeTable("Auto Mode Validation", func(c *api.ClusterConfig, expectedErr string) {
	err := api.ValidateAutoModeConfig(c)
	if expectedErr != "" {
		Expect(err).To(MatchError(expectedErr))
	} else {
		Expect(err).NotTo(HaveOccurred())
	}
},
	Entry("Auto Mode in an Outposts cluster", &api.ClusterConfig{
		Outpost: &api.Outpost{
			ControlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
		},
		AutoModeConfig: &api.AutoModeConfig{
			Enabled: api.Enabled(),
		},
	}, "Auto Mode is not supported on Outposts"),
	Entry("both nodeRoleARN and nodePools specified", &api.ClusterConfig{
		AutoModeConfig: &api.AutoModeConfig{
			Enabled:     api.Enabled(),
			NodeRoleARN: api.MustParseARN("arn:aws:iam::1234:role/CustomNodeRole"),
			NodePools:   &[]string{},
		},
	}, "cannot specify autoModeConfig.nodeRoleARN when autoModeConfig.nodePools is empty"),
	Entry("invalid nodePools", &api.ClusterConfig{
		AutoModeConfig: &api.AutoModeConfig{
			Enabled:   api.Enabled(),
			NodePools: &[]string{"invalid"},
		},
	}, `invalid NodePool "invalid"`),
	Entry("nodeRoleARN and nodePools specified when Auto Mode is disabled", &api.ClusterConfig{
		AutoModeConfig: &api.AutoModeConfig{
			Enabled:     api.Disabled(),
			NodeRoleARN: api.MustParseARN("arn:aws:iam::1234:role/CustomNodeRole"),
			NodePools:   &[]string{api.AutoModeNodePoolGeneralPurpose},
		},
	}, "cannot set autoModeConfig.nodeRoleARN or autoModeConfig.nodePools when Auto Mode is disabled"),
)
