package v1alpha5_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var _ = DescribeTable("Autonomous Mode Validation", func(c *api.ClusterConfig, expectedErr string) {
	err := api.ValidateAutonomousModeConfig(c)
	if expectedErr != "" {
		Expect(err).To(MatchError(expectedErr))
	} else {
		Expect(err).NotTo(HaveOccurred())
	}
},
	Entry("Autonomous Mode in an Outposts cluster", &api.ClusterConfig{
		Outpost: &api.Outpost{
			ControlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
		},
		AutonomousModeConfig: &api.AutonomousModeConfig{
			Enabled: api.Enabled(),
		},
	}, "Autonomous Mode is not supported on Outposts"),
	Entry("both nodeRoleARN and nodePools specified", &api.ClusterConfig{
		AutonomousModeConfig: &api.AutonomousModeConfig{
			Enabled:     api.Enabled(),
			NodeRoleARN: api.MustParseARN("arn:aws:iam::1234:role/CustomNodeRole"),
			NodePools:   &[]string{},
		},
	}, "cannot specify autonomousModeConfig.nodeRoleARN when autonomousModeConfig.nodePools is empty"),
	Entry("invalid nodePools", &api.ClusterConfig{
		AutonomousModeConfig: &api.AutonomousModeConfig{
			Enabled:   api.Enabled(),
			NodePools: &[]string{"invalid"},
		},
	}, `invalid NodePool "invalid"`),
	Entry("nodeRoleARN and nodePools specified when Autonomous Mode is disabled", &api.ClusterConfig{
		AutonomousModeConfig: &api.AutonomousModeConfig{
			Enabled:     api.Disabled(),
			NodeRoleARN: api.MustParseARN("arn:aws:iam::1234:role/CustomNodeRole"),
			NodePools:   &[]string{api.AutonomousModeNodePoolGeneralPurpose},
		},
	}, "cannot set autonomousModeConfig.nodeRoleARN or autonomousModeConfig.nodePools when Autonomous Mode is disabled"),
)
