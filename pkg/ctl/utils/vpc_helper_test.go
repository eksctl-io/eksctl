package utils_test

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/ctl/utils/mocks"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"

	"github.com/weaveworks/eksctl/pkg/ctl/utils"
)

type vpcHelperEntry struct {
	clusterVPC *ekstypes.VpcConfigResponse
	vpc        *api.ClusterVPC
	outposts   bool
	planMode   bool

	expectedUpdates []*eks.UpdateClusterConfigInput
	expectedErr     string
}

var _ = DescribeTable("VPCHelper", func(e vpcHelperEntry) {
	const updateClusterConfigMethodName = "UpdateClusterConfig"
	var vpcUpdater mocks.VPCConfigUpdater
	vpcUpdater.On(updateClusterConfigMethodName, mock.Anything, mock.Anything).Return(nil)

	clusterMeta := &api.ClusterMeta{
		Name: "test",
	}
	cluster := &ekstypes.Cluster{
		Name:               aws.String(clusterMeta.Name),
		ResourcesVpcConfig: e.clusterVPC,
	}
	if e.outposts {
		cluster.OutpostConfig = &ekstypes.OutpostConfigResponse{
			OutpostArns: []string{"arn:aws:outposts:us-west-2:1234:outpost/op-1234"},
		}
	}
	vpcHelper := &utils.VPCHelper{
		VPCUpdater:  &vpcUpdater,
		ClusterMeta: clusterMeta,
		Cluster:     cluster,
		PlanMode:    e.planMode,
	}
	err := vpcHelper.UpdateClusterVPCConfig(context.Background(), e.vpc)
	if e.expectedErr != "" {
		Expect(err).To(MatchError(ContainSubstring(e.expectedErr)))
	} else {
		Expect(err).NotTo(HaveOccurred())
	}
	vpcUpdater.AssertNumberOfCalls(GinkgoT(), updateClusterConfigMethodName, len(e.expectedUpdates))
	for _, u := range e.expectedUpdates {
		vpcUpdater.AssertCalled(GinkgoT(), updateClusterConfigMethodName, mock.Anything, u)
	}
},
	Entry("cluster matches default config", vpcHelperEntry{
		clusterVPC: &ekstypes.VpcConfigResponse{
			EndpointPublicAccess:  true,
			EndpointPrivateAccess: false,
			PublicAccessCidrs:     []string{"0.0.0.0/0"},
		},
		vpc: &api.ClusterVPC{},
	}),

	Entry("cluster matches desired config", vpcHelperEntry{
		clusterVPC: &ekstypes.VpcConfigResponse{
			EndpointPublicAccess:  true,
			EndpointPrivateAccess: false,
			PublicAccessCidrs:     []string{"0.0.0.0/0"},
		},
		vpc: &api.ClusterVPC{
			ClusterEndpoints: &api.ClusterEndpoints{
				PublicAccess:  api.Enabled(),
				PrivateAccess: api.Disabled(),
			},
		},
	}),

	Entry("cluster endpoint access does not match desired config", vpcHelperEntry{
		clusterVPC: &ekstypes.VpcConfigResponse{
			EndpointPublicAccess:  true,
			EndpointPrivateAccess: false,
			PublicAccessCidrs:     []string{"0.0.0.0/0"},
		},
		vpc: &api.ClusterVPC{
			ClusterEndpoints: &api.ClusterEndpoints{
				PublicAccess:  api.Enabled(),
				PrivateAccess: api.Enabled(),
			},
		},

		expectedUpdates: []*eks.UpdateClusterConfigInput{
			{
				Name: aws.String("test"),
				ResourcesVpcConfig: &ekstypes.VpcConfigRequest{
					EndpointPublicAccess:  api.Enabled(),
					EndpointPrivateAccess: api.Enabled(),
				},
			},
		},
	}),

	Entry("cluster public access CIDRs do not match desired config", vpcHelperEntry{
		clusterVPC: &ekstypes.VpcConfigResponse{
			EndpointPublicAccess:  true,
			EndpointPrivateAccess: false,
			PublicAccessCidrs:     []string{"0.0.0.0/0"},
		},
		vpc: &api.ClusterVPC{
			ClusterEndpoints: &api.ClusterEndpoints{
				PublicAccess:  api.Enabled(),
				PrivateAccess: api.Disabled(),
			},
			PublicAccessCIDRs: []string{"1.1.1.1/32", "2.2.2.2/32"},
		},

		expectedUpdates: []*eks.UpdateClusterConfigInput{
			{
				Name: aws.String("test"),
				ResourcesVpcConfig: &ekstypes.VpcConfigRequest{
					PublicAccessCidrs: []string{"1.1.1.1/32", "2.2.2.2/32"},
				},
			},
		},
	}),

	Entry("cluster public access CIDRs match desired config but out of order", vpcHelperEntry{
		clusterVPC: &ekstypes.VpcConfigResponse{
			EndpointPublicAccess:  true,
			EndpointPrivateAccess: false,
			PublicAccessCidrs:     []string{"2.2.2.2/32", "1.1.1.1/32"},
		},
		vpc: &api.ClusterVPC{
			ClusterEndpoints: &api.ClusterEndpoints{
				PublicAccess:  api.Enabled(),
				PrivateAccess: api.Disabled(),
			},
			PublicAccessCIDRs: []string{"1.1.1.1/32", "2.2.2.2/32"},
		},
	}),

	Entry("both cluster endpoint access and public access CIDRs do not match desired config", vpcHelperEntry{
		clusterVPC: &ekstypes.VpcConfigResponse{
			EndpointPublicAccess:  true,
			EndpointPrivateAccess: false,
			PublicAccessCidrs:     []string{"0.0.0.0/0"},
		},
		vpc: &api.ClusterVPC{
			ClusterEndpoints: &api.ClusterEndpoints{
				PublicAccess:  api.Disabled(),
				PrivateAccess: api.Enabled(),
			},
			PublicAccessCIDRs: []string{"1.1.1.1/32", "2.2.2.2/32"},
		},

		expectedUpdates: []*eks.UpdateClusterConfigInput{
			{
				Name: aws.String("test"),
				ResourcesVpcConfig: &ekstypes.VpcConfigRequest{
					EndpointPublicAccess:  api.Disabled(),
					EndpointPrivateAccess: api.Enabled(),
				},
			},
			{
				Name: aws.String("test"),
				ResourcesVpcConfig: &ekstypes.VpcConfigRequest{
					PublicAccessCidrs: []string{"1.1.1.1/32", "2.2.2.2/32"},
				},
			},
		},
	}),

	Entry("cluster does not match desired config but in plan mode", vpcHelperEntry{
		clusterVPC: &ekstypes.VpcConfigResponse{
			EndpointPublicAccess:  true,
			EndpointPrivateAccess: false,
			PublicAccessCidrs:     []string{"0.0.0.0/0"},
		},
		vpc: &api.ClusterVPC{
			ClusterEndpoints: &api.ClusterEndpoints{
				PublicAccess:  api.Disabled(),
				PrivateAccess: api.Enabled(),
			},
			PublicAccessCIDRs: []string{"1.1.1.1/32", "2.2.2.2/32"},
		},
		planMode: true,
	}),

	Entry("updating an Outpost cluster fails", vpcHelperEntry{
		clusterVPC: &ekstypes.VpcConfigResponse{
			EndpointPublicAccess:  true,
			EndpointPrivateAccess: false,
			PublicAccessCidrs:     []string{"0.0.0.0/0"},
		},
		vpc: &api.ClusterVPC{
			ClusterEndpoints: &api.ClusterEndpoints{
				PublicAccess:  api.Disabled(),
				PrivateAccess: api.Enabled(),
			},
			PublicAccessCIDRs: []string{"1.1.1.1/32", "2.2.2.2/32"},
		},
		outposts: true,

		expectedErr: "this operation is not supported on Outposts clusters",
	}),

	Entry("cluster matches desired config when subnets and security groups are specified", vpcHelperEntry{
		clusterVPC: &ekstypes.VpcConfigResponse{
			EndpointPublicAccess:  true,
			EndpointPrivateAccess: false,
			PublicAccessCidrs:     []string{"0.0.0.0/0"},
			SecurityGroupIds:      []string{"sg-1234"},
			SubnetIds:             []string{"subnet-1234"},
		},
		vpc: &api.ClusterVPC{
			ClusterEndpoints: &api.ClusterEndpoints{
				PublicAccess:  api.Enabled(),
				PrivateAccess: api.Disabled(),
			},
			ControlPlaneSecurityGroupIDs: []string{"sg-1234"},
			ControlPlaneSubnetIDs:        []string{"subnet-1234"},
		},
	}),

	Entry("cluster security groups do not match desired config", vpcHelperEntry{
		clusterVPC: &ekstypes.VpcConfigResponse{
			EndpointPublicAccess:  true,
			EndpointPrivateAccess: false,
			PublicAccessCidrs:     []string{"0.0.0.0/0"},
			SecurityGroupIds:      []string{"sg-1234"},
			SubnetIds:             []string{"subnet-1234"},
		},
		vpc: &api.ClusterVPC{
			ClusterEndpoints: &api.ClusterEndpoints{
				PublicAccess:  api.Enabled(),
				PrivateAccess: api.Disabled(),
			},
			ControlPlaneSecurityGroupIDs: []string{"sg-1234", "sg-5678"},
			ControlPlaneSubnetIDs:        []string{"subnet-1234"},
		},

		expectedUpdates: []*eks.UpdateClusterConfigInput{
			{
				Name: aws.String("test"),
				ResourcesVpcConfig: &ekstypes.VpcConfigRequest{
					SecurityGroupIds: []string{"sg-1234", "sg-5678"},
					SubnetIds:        []string{"subnet-1234"},
				},
			},
		},
	}),

	Entry("cluster subnets do not match desired config", vpcHelperEntry{
		clusterVPC: &ekstypes.VpcConfigResponse{
			EndpointPublicAccess:  true,
			EndpointPrivateAccess: false,
			PublicAccessCidrs:     []string{"0.0.0.0/0"},
			SecurityGroupIds:      []string{"sg-1234", "sg-5678"},
			SubnetIds:             []string{"subnet-1234"},
		},
		vpc: &api.ClusterVPC{
			ClusterEndpoints: &api.ClusterEndpoints{
				PublicAccess:  api.Enabled(),
				PrivateAccess: api.Disabled(),
			},
			ControlPlaneSecurityGroupIDs: []string{"sg-1234", "sg-5678"},
			ControlPlaneSubnetIDs:        []string{"subnet-1234", "subnet-5678"},
		},

		expectedUpdates: []*eks.UpdateClusterConfigInput{
			{
				Name: aws.String("test"),
				ResourcesVpcConfig: &ekstypes.VpcConfigRequest{
					SecurityGroupIds: []string{"sg-1234", "sg-5678"},
					SubnetIds:        []string{"subnet-1234", "subnet-5678"},
				},
			},
		},
	}),

	Entry("cluster security group and subnets do not match desired config", vpcHelperEntry{
		clusterVPC: &ekstypes.VpcConfigResponse{
			EndpointPublicAccess:  true,
			EndpointPrivateAccess: false,
			PublicAccessCidrs:     []string{"0.0.0.0/0"},
			SecurityGroupIds:      []string{"sg-1234", "sg-5678"},
			SubnetIds:             []string{"subnet-1234"},
		},
		vpc: &api.ClusterVPC{
			ClusterEndpoints: &api.ClusterEndpoints{
				PublicAccess:  api.Enabled(),
				PrivateAccess: api.Disabled(),
			},
			ControlPlaneSecurityGroupIDs: []string{"sg-1234", "sg-5678"},
			ControlPlaneSubnetIDs:        []string{"subnet-1234", "subnet-5678"},
		},

		expectedUpdates: []*eks.UpdateClusterConfigInput{
			{
				Name: aws.String("test"),
				ResourcesVpcConfig: &ekstypes.VpcConfigRequest{
					SecurityGroupIds: []string{"sg-1234", "sg-5678"},
					SubnetIds:        []string{"subnet-1234", "subnet-5678"},
				},
			},
		},
	}),

	Entry("no fields match desired config", vpcHelperEntry{
		clusterVPC: &ekstypes.VpcConfigResponse{
			EndpointPublicAccess:  false,
			EndpointPrivateAccess: true,
			PublicAccessCidrs:     []string{"0.0.0.0/0"},
			SecurityGroupIds:      []string{"sg-1234"},
			SubnetIds:             []string{"subnet-1234"},
		},
		vpc: &api.ClusterVPC{
			ClusterEndpoints: &api.ClusterEndpoints{
				PublicAccess:  api.Enabled(),
				PrivateAccess: api.Disabled(),
			},
			PublicAccessCIDRs:            []string{"1.1.1.1/1"},
			ControlPlaneSecurityGroupIDs: []string{"sg-1234", "sg-5678"},
			ControlPlaneSubnetIDs:        []string{"subnet-1234", "subnet-5678"},
		},

		expectedUpdates: []*eks.UpdateClusterConfigInput{
			{
				Name: aws.String("test"),
				ResourcesVpcConfig: &ekstypes.VpcConfigRequest{
					EndpointPublicAccess:  api.Enabled(),
					EndpointPrivateAccess: api.Disabled(),
				},
			},
			{
				Name: aws.String("test"),
				ResourcesVpcConfig: &ekstypes.VpcConfigRequest{
					PublicAccessCidrs: []string{"1.1.1.1/1"},
				},
			},
			{
				Name: aws.String("test"),
				ResourcesVpcConfig: &ekstypes.VpcConfigRequest{
					SecurityGroupIds: []string{"sg-1234", "sg-5678"},
					SubnetIds:        []string{"subnet-1234", "subnet-5678"},
				},
			},
		},
	}),
)
