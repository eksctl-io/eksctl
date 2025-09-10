package builder_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/weaveworks/eksctl/pkg/goformation"
	gfneks "github.com/weaveworks/eksctl/pkg/goformation/cloudformation/eks"
	gfnt "github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"

	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/require"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	cft "github.com/weaveworks/eksctl/pkg/cfn/template"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap/fakes"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
	vpcfakes "github.com/weaveworks/eksctl/pkg/vpc/fakes"
)

func TestManagedPolicyResources(t *testing.T) {
	iamRoleTests := []struct {
		addons                  api.NodeGroupIAMAddonPolicies
		attachPolicy            api.InlineDocument
		attachPolicyARNs        []string
		expectedNewPolicies     []string
		expectedManagedPolicies []*gfnt.Value
		description             string
	}{
		{
			expectedManagedPolicies: makePartitionedPolicies("AmazonEKSWorkerNodePolicy", "AmazonEKS_CNI_Policy", "AmazonEC2ContainerRegistryPullOnly", "AmazonSSMManagedInstanceCore"),
			description:             "Default policies",
		},
		{
			addons: api.NodeGroupIAMAddonPolicies{
				ImageBuilder: api.Enabled(),
			},
			expectedManagedPolicies: makePartitionedPolicies("AmazonEKSWorkerNodePolicy", "AmazonEKS_CNI_Policy",
				"AmazonEC2ContainerRegistryPullOnly", "AmazonEC2ContainerRegistryPowerUser", "AmazonSSMManagedInstanceCore"),
			description: "ImageBuilder enabled",
		},
		{
			addons: api.NodeGroupIAMAddonPolicies{
				CloudWatch: api.Enabled(),
			},
			expectedManagedPolicies: makePartitionedPolicies("AmazonEKSWorkerNodePolicy", "AmazonEKS_CNI_Policy",
				"AmazonEC2ContainerRegistryPullOnly", "AmazonSSMManagedInstanceCore", "CloudWatchAgentServerPolicy"),
			description: "CloudWatch enabled",
		},
		{
			addons: api.NodeGroupIAMAddonPolicies{
				AutoScaler: api.Enabled(),
			},
			expectedNewPolicies:     []string{"PolicyAutoScaling"},
			expectedManagedPolicies: makePartitionedPolicies("AmazonEKSWorkerNodePolicy", "AmazonEKS_CNI_Policy", "AmazonEC2ContainerRegistryPullOnly", "AmazonSSMManagedInstanceCore"),
			description:             "AutoScaler enabled",
		},
		{
			attachPolicy: cft.MakePolicyDocument(cft.MapOfInterfaces{
				"Effect": "Allow",
				"Action": []string{
					"s3:Get*",
				},
				"Resource": "*",
			}),
			expectedNewPolicies:     []string{"Policy1"},
			expectedManagedPolicies: makePartitionedPolicies("AmazonEKSWorkerNodePolicy", "AmazonEKS_CNI_Policy", "AmazonEC2ContainerRegistryPullOnly", "AmazonSSMManagedInstanceCore"),
			description:             "Custom inline policies",
		},
		{
			attachPolicyARNs:        []string{"AmazonEKSWorkerNodePolicy", "AmazonEKS_CNI_Policy"},
			expectedManagedPolicies: subs(prefixPolicies("AmazonEKSWorkerNodePolicy", "AmazonEKS_CNI_Policy")),
			description:             "Custom managed policies",
		},
		// should not attach any additional policies
		{
			attachPolicyARNs:        []string{"CloudWatchAgentServerPolicy"},
			expectedManagedPolicies: subs(prefixPolicies("CloudWatchAgentServerPolicy")),
			description:             "Custom managed policies",
		},
		// no duplicate values
		{
			attachPolicyARNs: []string{"AmazonEC2ContainerRegistryPowerUser"},
			addons: api.NodeGroupIAMAddonPolicies{
				ImageBuilder: api.Enabled(),
			},
			expectedManagedPolicies: subs(prefixPolicies("AmazonEC2ContainerRegistryPowerUser")),
			description:             "Duplicate policies",
		},
		{
			attachPolicyARNs: []string{"CloudWatchAgentServerPolicy", "AmazonEC2ContainerRegistryPowerUser"},
			addons: api.NodeGroupIAMAddonPolicies{
				ImageBuilder: api.Enabled(),
				CloudWatch:   api.Enabled(),
			},
			expectedManagedPolicies: subs(prefixPolicies("CloudWatchAgentServerPolicy", "AmazonEC2ContainerRegistryPowerUser")),
			description:             "Multiple duplicate policies",
		},
	}

	for i, tt := range iamRoleTests {
		t.Run(fmt.Sprintf("%d: %s", i, tt.description), func(t *testing.T) {
			require := require.New(t)
			clusterConfig := api.NewClusterConfig()

			ng := api.NewManagedNodeGroup()
			err := api.SetManagedNodeGroupDefaults(ng, clusterConfig.Metadata, false)
			require.NoError(err)
			ng.IAM.WithAddonPolicies = tt.addons
			ng.IAM.AttachPolicy = tt.attachPolicy
			ng.IAM.AttachPolicyARNs = prefixPolicies(tt.attachPolicyARNs...)

			p := mockprovider.NewMockProvider()
			fakeVPCImporter := new(vpcfakes.FakeImporter)
			bootstrapper := &fakes.FakeBootstrapper{}
			bootstrapper.UserDataStub = func() (string, error) {
				return "", nil
			}
			mockSubnetsAndAZInstanceSupport(clusterConfig, p,
				[]string{"us-west-2a"},
				[]string{}, // local zones
				[]ec2types.InstanceType{api.DefaultNodeType})
			stack := builder.NewManagedNodeGroup(p.EC2(), clusterConfig, ng, nil, bootstrapper, false, fakeVPCImporter)
			err = stack.AddAllResources(context.Background())
			require.Nil(err)

			bytes, err := stack.RenderJSON()
			require.NoError(err)

			template, err := goformation.ParseJSON(bytes)
			require.NoError(err)

			role, err := template.GetIAMRoleWithName(builder.GetIAMRoleName())
			require.NoError(err)

			require.ElementsMatch(tt.expectedManagedPolicies, role.ManagedPolicyArns.Raw().(gfnt.Slice))

			policyNames := make([]string, 0)
			for name := range template.GetAllIAMPolicyResources() {
				policyNames = append(policyNames, name)
			}
			require.ElementsMatch(tt.expectedNewPolicies, policyNames)

			// assert custom inline policy matches
			if tt.attachPolicy != nil {
				policy, err := template.GetIAMPolicyWithName("Policy1")
				require.NoError(err)

				// convert to json for comparison since interfaces are not identical
				expectedPolicy, err := json.Marshal(tt.attachPolicy)
				require.NoError(err)
				actualPolicy, err := json.Marshal(policy.PolicyDocument)
				require.NoError(err)
				require.Equal(string(expectedPolicy), string(actualPolicy))
			}
		})
	}

}

func TestManagedNodeRole(t *testing.T) {
	nodeRoleTests := []struct {
		description         string
		nodeGroup           *api.ManagedNodeGroup
		expectedNewRole     bool
		expectedNodeRoleARN *gfnt.Value
	}{
		{
			description: "InstanceRoleARN is not provided",
			nodeGroup: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{},
			},
			expectedNewRole:     true,
			expectedNodeRoleARN: gfnt.MakeFnGetAtt(builder.GetIAMRoleName(), gfnt.NewString("Arn")), // creating new role
		},
		{
			description: "InstanceRoleARN is provided",
			nodeGroup: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					IAM: &api.NodeGroupIAM{
						InstanceRoleARN: "arn::DUMMY::DUMMYROLE",
					},
				},
			},
			expectedNewRole:     false,
			expectedNodeRoleARN: gfnt.NewString("arn::DUMMY::DUMMYROLE"), // using the provided role
		},
		{
			description: "InstanceRoleARN is provided and normalized",
			nodeGroup: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					IAM: &api.NodeGroupIAM{
						InstanceRoleARN: "arn:aws:iam::1234567890:role/foo/bar/baz/custom-eks-role",
					},
				},
			},
			expectedNewRole:     false,
			expectedNodeRoleARN: gfnt.NewString("arn:aws:iam::1234567890:role/custom-eks-role"),
		},
	}

	for i, tt := range nodeRoleTests {
		t.Run(fmt.Sprintf("%d: %s", i, tt.description), func(t *testing.T) {
			require := require.New(t)
			clusterConfig := api.NewClusterConfig()
			clusterConfig.Status = &api.ClusterStatus{}
			err := api.SetManagedNodeGroupDefaults(tt.nodeGroup, clusterConfig.Metadata, false)
			require.NoError(err)
			p := mockprovider.NewMockProvider()
			fakeVPCImporter := new(vpcfakes.FakeImporter)
			bootstrapper, err := nodebootstrap.NewManagedBootstrapper(clusterConfig, tt.nodeGroup)
			require.NoError(err)
			mockSubnetsAndAZInstanceSupport(clusterConfig, p,
				[]string{"us-west-2a"},
				[]string{}, // local zones
				[]ec2types.InstanceType{api.DefaultNodeType})
			stack := builder.NewManagedNodeGroup(p.EC2(), clusterConfig, tt.nodeGroup, nil, bootstrapper, false, fakeVPCImporter)
			err = stack.AddAllResources(context.Background())
			require.NoError(err)

			bytes, err := stack.RenderJSON()
			require.NoError(err)

			template, err := goformation.ParseJSON(bytes)
			require.NoError(err)
			ngResource, ok := template.Resources[builder.ManagedNodeGroupResourceName]
			require.True(ok)
			ng, ok := ngResource.(*gfneks.Nodegroup)
			require.True(ok)
			require.Equal(tt.expectedNodeRoleARN, ng.NodeRole)

			_, ok = template.GetAllIAMRoleResources()[builder.GetIAMRoleName()]
			require.Equal(tt.expectedNewRole, ok)
		})
	}
}

func makePartitionedPolicies(policies ...string) []*gfnt.Value {
	var partitionedPolicies []*gfnt.Value
	for _, policy := range policies {
		partitionedPolicies = append(partitionedPolicies, gfnt.MakeFnSubString("arn:${AWS::Partition}:iam::aws:policy/"+policy))
	}
	return partitionedPolicies
}

func prefixPolicies(policies ...string) []string {
	var prefixedPolicies []string
	for _, policy := range policies {
		prefixedPolicies = append(prefixedPolicies, "arn:aws:iam::aws:policy/"+policy)
	}
	return prefixedPolicies
}

func subs(ss []string) []*gfnt.Value {
	var subs []*gfnt.Value
	for _, s := range ss {
		subs = append(subs, gfnt.NewString(s))
	}
	return subs
}

func TestManagedNodeGroupNodeRepairConfig(t *testing.T) {
	nodeRepairTests := []struct {
		description      string
		nodeRepairConfig *api.NodeGroupNodeRepairConfig
		expectedConfig   *gfneks.Nodegroup_NodeRepairConfig
	}{
		{
			description:      "nil node repair config",
			nodeRepairConfig: nil,
			expectedConfig:   nil,
		},
		{
			description: "enabled only",
			nodeRepairConfig: &api.NodeGroupNodeRepairConfig{
				Enabled: api.Enabled(),
			},
			expectedConfig: &gfneks.Nodegroup_NodeRepairConfig{
				Enabled: gfnt.NewBoolean(true),
			},
		},
		{
			description: "disabled only",
			nodeRepairConfig: &api.NodeGroupNodeRepairConfig{
				Enabled: api.Disabled(),
			},
			expectedConfig: &gfneks.Nodegroup_NodeRepairConfig{
				Enabled: gfnt.NewBoolean(false),
			},
		},
		{
			description: "all threshold and parallel parameters",
			nodeRepairConfig: &api.NodeGroupNodeRepairConfig{
				Enabled:                                 api.Enabled(),
				MaxUnhealthyNodeThresholdPercentage:     aws.Int(20),
				MaxUnhealthyNodeThresholdCount:          aws.Int(5),
				MaxParallelNodesRepairedPercentage:      aws.Int(15),
				MaxParallelNodesRepairedCount:           aws.Int(2),
			},
			expectedConfig: &gfneks.Nodegroup_NodeRepairConfig{
				Enabled:                                 gfnt.NewBoolean(true),
				MaxUnhealthyNodeThresholdPercentage:     gfnt.NewInteger(20),
				MaxUnhealthyNodeThresholdCount:          gfnt.NewInteger(5),
				MaxParallelNodesRepairedPercentage:      gfnt.NewInteger(15),
				MaxParallelNodesRepairedCount:           gfnt.NewInteger(2),
			},
		},
		{
			description: "with node repair config overrides",
			nodeRepairConfig: &api.NodeGroupNodeRepairConfig{
				Enabled: api.Enabled(),
				NodeRepairConfigOverrides: []api.NodeRepairConfigOverride{
					{
						NodeMonitoringCondition: "AcceleratedInstanceNotReady",
						NodeUnhealthyReason:     "NvidiaXID13Error",
						MinRepairWaitTimeMins:   10,
						RepairAction:            "Terminate",
					},
					{
						NodeMonitoringCondition: "NetworkNotReady",
						NodeUnhealthyReason:     "InterfaceNotUp",
						MinRepairWaitTimeMins:   20,
						RepairAction:            "Restart",
					},
				},
			},
			expectedConfig: &gfneks.Nodegroup_NodeRepairConfig{
				Enabled: gfnt.NewBoolean(true),
				NodeRepairConfigOverrides: []gfneks.Nodegroup_NodeRepairConfigOverride{
					{
						NodeMonitoringCondition: gfnt.NewString("AcceleratedInstanceNotReady"),
						NodeUnhealthyReason:     gfnt.NewString("NvidiaXID13Error"),
						MinRepairWaitTimeMins:   gfnt.NewInteger(10),
						RepairAction:            gfnt.NewString("Terminate"),
					},
					{
						NodeMonitoringCondition: gfnt.NewString("NetworkNotReady"),
						NodeUnhealthyReason:     gfnt.NewString("InterfaceNotUp"),
						MinRepairWaitTimeMins:   gfnt.NewInteger(20),
						RepairAction:            gfnt.NewString("Restart"),
					},
				},
			},
		},
		{
			description: "comprehensive configuration",
			nodeRepairConfig: &api.NodeGroupNodeRepairConfig{
				Enabled:                                 api.Enabled(),
				MaxUnhealthyNodeThresholdPercentage:     aws.Int(25),
				MaxParallelNodesRepairedCount:           aws.Int(3),
				NodeRepairConfigOverrides: []api.NodeRepairConfigOverride{
					{
						NodeMonitoringCondition: "NetworkNotReady",
						NodeUnhealthyReason:     "InterfaceNotUp",
						MinRepairWaitTimeMins:   15,
						RepairAction:            "Restart",
					},
				},
			},
			expectedConfig: &gfneks.Nodegroup_NodeRepairConfig{
				Enabled:                           gfnt.NewBoolean(true),
				MaxUnhealthyNodeThresholdPercentage: gfnt.NewInteger(25),
				MaxParallelNodesRepairedCount:     gfnt.NewInteger(3),
				NodeRepairConfigOverrides: []gfneks.Nodegroup_NodeRepairConfigOverride{
					{
						NodeMonitoringCondition: gfnt.NewString("NetworkNotReady"),
						NodeUnhealthyReason:     gfnt.NewString("InterfaceNotUp"),
						MinRepairWaitTimeMins:   gfnt.NewInteger(15),
						RepairAction:            gfnt.NewString("Restart"),
					},
				},
			},
		},
	}

	for _, tt := range nodeRepairTests {
		t.Run(tt.description, func(t *testing.T) {
			clusterConfig := api.NewClusterConfig()
			clusterConfig.Metadata.Name = "test-cluster"
			clusterConfig.Metadata.Region = "us-west-2"

			ng := &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Name:         "test-ng",
					InstanceType: "m5.large",
				},
				NodeRepairConfig: tt.nodeRepairConfig,
			}

			clusterConfig.Status = &api.ClusterStatus{}
			err := api.SetManagedNodeGroupDefaults(ng, clusterConfig.Metadata, false)
			require.NoError(t, err)
			
			p := mockprovider.NewMockProvider()
			fakeVPCImporter := new(vpcfakes.FakeImporter)
			bootstrapper, err := nodebootstrap.NewManagedBootstrapper(clusterConfig, ng)
			require.NoError(t, err)
			
			// Mock subnets and AZ instance support like other tests
			mockSubnetsAndAZInstanceSupport(clusterConfig, p,
				[]string{"us-west-2a"},
				[]string{}, // local zones
				[]ec2types.InstanceType{api.DefaultNodeType})

			stack := builder.NewManagedNodeGroup(p.EC2(), clusterConfig, ng, nil, bootstrapper, false, fakeVPCImporter)
			err = stack.AddAllResources(context.Background())
			require.NoError(t, err)

			bytes, err := stack.RenderJSON()
			require.NoError(t, err)

			template, err := goformation.ParseJSON(bytes)
			require.NoError(t, err)
			
			// Get the managed nodegroup resource
			ngResource, ok := template.Resources[builder.ManagedNodeGroupResourceName]
			require.True(t, ok, "ManagedNodeGroup resource should exist")
			managedNodeGroup, ok := ngResource.(*gfneks.Nodegroup)
			require.True(t, ok, "Resource should be a Nodegroup")

			// Test the node repair config
			if tt.expectedConfig == nil {
				require.Nil(t, managedNodeGroup.NodeRepairConfig, "NodeRepairConfig should be nil")
			} else {
				require.NotNil(t, managedNodeGroup.NodeRepairConfig, "NodeRepairConfig should not be nil")
				
				// Test enabled field
				if tt.expectedConfig.Enabled != nil {
					require.NotNil(t, managedNodeGroup.NodeRepairConfig.Enabled)
					require.Equal(t, tt.expectedConfig.Enabled.Raw(), managedNodeGroup.NodeRepairConfig.Enabled.Raw())
				} else {
					require.Nil(t, managedNodeGroup.NodeRepairConfig.Enabled)
				}

				// Test threshold percentage
				if tt.expectedConfig.MaxUnhealthyNodeThresholdPercentage != nil {
					require.NotNil(t, managedNodeGroup.NodeRepairConfig.MaxUnhealthyNodeThresholdPercentage)
					require.Equal(t, tt.expectedConfig.MaxUnhealthyNodeThresholdPercentage.Raw(), 
						managedNodeGroup.NodeRepairConfig.MaxUnhealthyNodeThresholdPercentage.Raw())
				} else {
					require.Nil(t, managedNodeGroup.NodeRepairConfig.MaxUnhealthyNodeThresholdPercentage)
				}

				// Test threshold count
				if tt.expectedConfig.MaxUnhealthyNodeThresholdCount != nil {
					require.NotNil(t, managedNodeGroup.NodeRepairConfig.MaxUnhealthyNodeThresholdCount)
					require.Equal(t, tt.expectedConfig.MaxUnhealthyNodeThresholdCount.Raw(), 
						managedNodeGroup.NodeRepairConfig.MaxUnhealthyNodeThresholdCount.Raw())
				} else {
					require.Nil(t, managedNodeGroup.NodeRepairConfig.MaxUnhealthyNodeThresholdCount)
				}

				// Test parallel percentage
				if tt.expectedConfig.MaxParallelNodesRepairedPercentage != nil {
					require.NotNil(t, managedNodeGroup.NodeRepairConfig.MaxParallelNodesRepairedPercentage)
					require.Equal(t, tt.expectedConfig.MaxParallelNodesRepairedPercentage.Raw(), 
						managedNodeGroup.NodeRepairConfig.MaxParallelNodesRepairedPercentage.Raw())
				} else {
					require.Nil(t, managedNodeGroup.NodeRepairConfig.MaxParallelNodesRepairedPercentage)
				}

				// Test parallel count
				if tt.expectedConfig.MaxParallelNodesRepairedCount != nil {
					require.NotNil(t, managedNodeGroup.NodeRepairConfig.MaxParallelNodesRepairedCount)
					require.Equal(t, tt.expectedConfig.MaxParallelNodesRepairedCount.Raw(), 
						managedNodeGroup.NodeRepairConfig.MaxParallelNodesRepairedCount.Raw())
				} else {
					require.Nil(t, managedNodeGroup.NodeRepairConfig.MaxParallelNodesRepairedCount)
				}

				// Test overrides
				require.Equal(t, len(tt.expectedConfig.NodeRepairConfigOverrides), 
					len(managedNodeGroup.NodeRepairConfig.NodeRepairConfigOverrides))
				
				for i, expectedOverride := range tt.expectedConfig.NodeRepairConfigOverrides {
					actualOverride := managedNodeGroup.NodeRepairConfig.NodeRepairConfigOverrides[i]
					require.Equal(t, expectedOverride.NodeMonitoringCondition.Raw(), 
						actualOverride.NodeMonitoringCondition.Raw())
					require.Equal(t, expectedOverride.NodeUnhealthyReason.Raw(), 
						actualOverride.NodeUnhealthyReason.Raw())
					require.Equal(t, expectedOverride.MinRepairWaitTimeMins.Raw(), 
						actualOverride.MinRepairWaitTimeMins.Raw())
					require.Equal(t, expectedOverride.RepairAction.Raw(), 
						actualOverride.RepairAction.Raw())
				}
			}
		})
	}
}
