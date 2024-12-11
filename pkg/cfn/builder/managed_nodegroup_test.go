package builder_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"goformation/v4"
	gfneks "goformation/v4/cloudformation/eks"
	gfnt "goformation/v4/cloudformation/types"

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
			expectedManagedPolicies: makePartitionedPolicies("AmazonEKSWorkerNodePolicy", "AmazonEKS_CNI_Policy", "AmazonEC2ContainerRegistryReadOnly", "AmazonSSMManagedInstanceCore"),
			description:             "Default policies",
		},
		{
			addons: api.NodeGroupIAMAddonPolicies{
				ImageBuilder: api.Enabled(),
			},
			expectedManagedPolicies: makePartitionedPolicies("AmazonEKSWorkerNodePolicy", "AmazonEKS_CNI_Policy",
				"AmazonEC2ContainerRegistryReadOnly", "AmazonEC2ContainerRegistryPowerUser", "AmazonSSMManagedInstanceCore"),
			description: "ImageBuilder enabled",
		},
		{
			addons: api.NodeGroupIAMAddonPolicies{
				CloudWatch: api.Enabled(),
			},
			expectedManagedPolicies: makePartitionedPolicies("AmazonEKSWorkerNodePolicy", "AmazonEKS_CNI_Policy",
				"AmazonEC2ContainerRegistryReadOnly", "AmazonSSMManagedInstanceCore", "CloudWatchAgentServerPolicy"),
			description: "CloudWatch enabled",
		},
		{
			addons: api.NodeGroupIAMAddonPolicies{
				AutoScaler: api.Enabled(),
			},
			expectedNewPolicies:     []string{"PolicyAutoScaling"},
			expectedManagedPolicies: makePartitionedPolicies("AmazonEKSWorkerNodePolicy", "AmazonEKS_CNI_Policy", "AmazonEC2ContainerRegistryReadOnly", "AmazonSSMManagedInstanceCore"),
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
			expectedManagedPolicies: makePartitionedPolicies("AmazonEKSWorkerNodePolicy", "AmazonEKS_CNI_Policy", "AmazonEC2ContainerRegistryReadOnly", "AmazonSSMManagedInstanceCore"),
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
			api.SetManagedNodeGroupDefaults(ng, clusterConfig.Metadata, false)
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
			err := stack.AddAllResources(context.Background())
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
			api.SetManagedNodeGroupDefaults(tt.nodeGroup, clusterConfig.Metadata, false)
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
