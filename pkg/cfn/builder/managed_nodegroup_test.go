package builder

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap/fakes"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
	vpcfakes "github.com/weaveworks/eksctl/pkg/vpc/fakes"
	"github.com/weaveworks/goformation/v4"
	gfneks "github.com/weaveworks/goformation/v4/cloudformation/eks"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"
)

func TestManagedPolicyResources(t *testing.T) {
	iamRoleTests := []struct {
		addons                  api.NodeGroupIAMAddonPolicies
		attachPolicyARNs        []string
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
			attachPolicyARNs:        []string{"AmazonEKSWorkerNodePolicy", "AmazonEKS_CNI_Policy"},
			expectedManagedPolicies: subs(prefixPolicies("AmazonEKSWorkerNodePolicy", "AmazonEKS_CNI_Policy")),
			description:             "Custom policies",
		},
		// should not attach any additional policies
		{
			attachPolicyARNs:        []string{"CloudWatchAgentServerPolicy"},
			expectedManagedPolicies: subs(prefixPolicies("CloudWatchAgentServerPolicy")),
			description:             "Custom policies",
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
			api.SetManagedNodeGroupDefaults(ng, clusterConfig.Metadata)
			ng.IAM.WithAddonPolicies = tt.addons
			ng.IAM.AttachPolicyARNs = prefixPolicies(tt.attachPolicyARNs...)

			p := mockprovider.NewMockProvider()
			fakeVPCImporter := new(vpcfakes.FakeImporter)
			bootstrapper := &fakes.FakeBootstrapper{}
			bootstrapper.UserDataStub = func() (string, error) {
				return "", nil
			}
			stack := NewManagedNodeGroup(p.EC2(), clusterConfig, ng, nil, bootstrapper, false, fakeVPCImporter)
			err := stack.AddAllResources()
			require.Nil(err)

			bytes, err := stack.RenderJSON()
			require.NoError(err)

			template, err := goformation.ParseJSON(bytes)
			require.NoError(err)

			role, ok := template.GetAllIAMRoleResources()["NodeInstanceRole"]
			require.True(ok)

			require.ElementsMatch(tt.expectedManagedPolicies, role.ManagedPolicyArns.Raw().(gfnt.Slice))

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
			expectedNodeRoleARN: gfnt.MakeFnGetAtt(cfnIAMInstanceRoleName, gfnt.NewString("Arn")), // creating new role
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
			api.SetManagedNodeGroupDefaults(tt.nodeGroup, clusterConfig.Metadata)
			p := mockprovider.NewMockProvider()
			fakeVPCImporter := new(vpcfakes.FakeImporter)
			bootstrapper := nodebootstrap.NewManagedBootstrapper(clusterConfig, tt.nodeGroup)
			stack := NewManagedNodeGroup(p.EC2(), clusterConfig, tt.nodeGroup, nil, bootstrapper, false, fakeVPCImporter)
			err := stack.AddAllResources()
			require.NoError(err)

			bytes, err := stack.RenderJSON()
			require.NoError(err)

			template, err := goformation.ParseJSON(bytes)
			require.NoError(err)
			ngResource, ok := template.Resources["ManagedNodeGroup"]
			require.True(ok)
			ng, ok := ngResource.(*gfneks.Nodegroup)
			require.True(ok)
			require.Equal(tt.expectedNodeRoleARN, ng.NodeRole)

			_, ok = template.GetAllIAMRoleResources()[cfnIAMInstanceRoleName]
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
