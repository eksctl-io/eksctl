package autonomousmode_test

import (
	"context"
	"fmt"

	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/autonomousmode"
	"github.com/weaveworks/eksctl/pkg/autonomousmode/mocks"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

type subnetsLoaderTest struct {
	ignoreMissingSubnets bool
	updateMocks          func(*mocks.ClusterStackDescriber, *mocks.VPCImporter)
	updateClusterConfig  func(*api.ClusterConfig)
	expectedSubnetIDs    []string
	expectedErr          string
}

var _ = DescribeTable("Subnets Loader", func(t subnetsLoaderTest) {
	var clusterStackDescriber mocks.ClusterStackDescriber
	var vpcImporter mocks.VPCImporter
	clusterConfig := api.NewClusterConfig()
	clusterConfig.Metadata.Name = "cluster"
	if t.updateMocks != nil {
		t.updateMocks(&clusterStackDescriber, &vpcImporter)
	}
	if t.updateClusterConfig != nil {
		t.updateClusterConfig(clusterConfig)
	}
	subnetsLoader := &autonomousmode.SubnetsLoader{
		ClusterStackDescriber: &clusterStackDescriber,
		VPCImporter:           &vpcImporter,
		IgnoreMissingSubnets:  t.ignoreMissingSubnets,
	}
	subnetIDs, _, err := subnetsLoader.LoadSubnets(context.Background(), clusterConfig)
	if t.expectedErr != "" {
		Expect(err).To(MatchError(t.expectedErr))
	} else {
		Expect(err).NotTo(HaveOccurred())
		Expect(subnetIDs).To(ConsistOf(t.expectedSubnetIDs))
	}
	for _, asserter := range []interface {
		AssertExpectations(t mock.TestingT) bool
	}{
		&clusterStackDescriber,
		&vpcImporter,
	} {
		asserter.AssertExpectations(GinkgoT())
	}
},
	Entry("no dedicated VPC", subnetsLoaderTest{
		updateMocks: func(c *mocks.ClusterStackDescriber, v *mocks.VPCImporter) {
			c.EXPECT().ClusterHasDedicatedVPC(mock.Anything).Return(false, nil).Once()
		},
	}),
	Entry("no cluster stack", subnetsLoaderTest{
		updateMocks: func(c *mocks.ClusterStackDescriber, v *mocks.VPCImporter) {
			c.EXPECT().ClusterHasDedicatedVPC(mock.Anything).Return(true, nil).Once()
			c.EXPECT().DescribeClusterStack(mock.Anything).Return(nil, &manager.StackNotFoundErr{
				ClusterName: "cluster",
			}).Once()
		},
	}),
	Entry("drifted stack", subnetsLoaderTest{
		updateMocks: mockDedicatedVPCWithError(&vpc.StackDriftError{
			Msg: "drift error",
		}),
		expectedErr: "loading cluster VPC: drift error; to skip patching NodeClass to use private subnets and ignore this error, " +
			"please retry the command with --ignore-missing-subnets and patch the NodeClass " +
			"resource manually if you do not want to use cluster subnets for Autonomous Mode",
	}),
	Entry("drifted stack but ignoreMissingSubnets set", subnetsLoaderTest{
		ignoreMissingSubnets: true,
		updateMocks: mockDedicatedVPCWithError(&vpc.StackDriftError{
			Msg: "drift error",
		}),
	}),
	Entry("no private subnets", subnetsLoaderTest{
		updateMocks: mockDedicatedVPC,
		updateClusterConfig: func(c *api.ClusterConfig) {
			c.VPC.Subnets = &api.ClusterSubnets{
				Private: api.AZSubnetMapping{},
			}
		},
		expectedErr: "expected to find private subnets in cluster stack",
	}),
	Entry("insufficient private subnets", subnetsLoaderTest{
		updateMocks: mockDedicatedVPC,
		updateClusterConfig: func(c *api.ClusterConfig) {
			c.VPC.Subnets = &api.ClusterSubnets{
				Private: api.AZSubnetMapping{
					"us-west-2a": api.AZSubnetSpec{
						ID: "subnet-1",
					},
				},
			}
		},
		expectedErr: fmt.Sprintf("Autonomous Mode requires at least two private subnets; got %v", []string{"subnet-1"}),
	}),
	Entry("sufficient private subnets", subnetsLoaderTest{
		updateMocks: mockDedicatedVPC,
		updateClusterConfig: func(c *api.ClusterConfig) {
			c.VPC.Subnets = &api.ClusterSubnets{
				Private: api.AZSubnetMapping{
					"us-west-2a": api.AZSubnetSpec{
						ID: "subnet-1",
					},
					"us-west-2b": api.AZSubnetSpec{
						ID: "subnet-2",
					},
				},
			}
		},
		expectedSubnetIDs: []string{"subnet-1", "subnet-2"},
	}),
)

func mockDedicatedVPC(c *mocks.ClusterStackDescriber, v *mocks.VPCImporter) {
	mockDedicatedVPCWithError(nil)(c, v)
}

func mockDedicatedVPCWithError(loadClusterVPCErr error) func(c *mocks.ClusterStackDescriber, v *mocks.VPCImporter) {
	return func(c *mocks.ClusterStackDescriber, v *mocks.VPCImporter) {
		c.EXPECT().ClusterHasDedicatedVPC(mock.Anything).Return(true, nil).Once()
		var stack cfntypes.Stack
		c.EXPECT().DescribeClusterStack(mock.Anything).Return(&stack, nil).Once()
		v.EXPECT().LoadClusterVPC(mock.Anything, mock.Anything, &stack, false).Return(loadClusterVPCErr).Once()
	}
}
