package builder

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
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
			expectedManagedPolicies: makePartitionedPolicies("AmazonEKSWorkerNodePolicy", "AmazonEKS_CNI_Policy", "AmazonEC2ContainerRegistryReadOnly"),
			description:             "Default policies",
		},
		{
			addons: api.NodeGroupIAMAddonPolicies{
				ImageBuilder: api.Enabled(),
			},
			expectedManagedPolicies: makePartitionedPolicies("AmazonEKSWorkerNodePolicy", "AmazonEKS_CNI_Policy",
				"AmazonEC2ContainerRegistryReadOnly", "AmazonEC2ContainerRegistryPowerUser"),
			description: "ImageBuilder enabled",
		},
		{
			addons: api.NodeGroupIAMAddonPolicies{
				CloudWatch: api.Enabled(),
			},
			expectedManagedPolicies: makePartitionedPolicies("AmazonEKSWorkerNodePolicy", "AmazonEKS_CNI_Policy",
				"AmazonEC2ContainerRegistryReadOnly", "CloudWatchAgentServerPolicy"),
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
			assert := assert.New(t)
			clusterConfig := api.NewClusterConfig()

			ng := api.NewManagedNodeGroup()
			ng.IAM.WithAddonPolicies = tt.addons
			ng.IAM.AttachPolicyARNs = prefixPolicies(tt.attachPolicyARNs...)

			stack := NewManagedNodeGroup(clusterConfig, ng, "iam-test")
			err := stack.AddAllResources()
			assert.Nil(err)

			bytes, err := stack.RenderJSON()
			assert.NoError(err)

			template, err := goformation.ParseJSON(bytes)
			assert.NoError(err)

			role, ok := template.GetAllIAMRoleResources()["NodeInstanceRole"]
			assert.True(ok)

			assert.ElementsMatch(tt.expectedManagedPolicies, role.ManagedPolicyArns)

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
				ScalingConfig: &api.ScalingConfig{},
				SSH: &api.NodeGroupSSH{
					Allow: api.Disabled(),
				},
				IAM: &api.NodeGroupIAM{},
			},
			expectedNewRole:     true,
			expectedNodeRoleARN: gfnt.MakeFnGetAtt(cfnIAMInstanceRoleName, gfnt.NewString("Arn")), // creating new role
		},
		{
			description: "InstanceRoleARN is provided",
			nodeGroup: &api.ManagedNodeGroup{
				ScalingConfig: &api.ScalingConfig{},
				SSH: &api.NodeGroupSSH{
					Allow: api.Disabled(),
				},
				IAM: &api.NodeGroupIAM{
					InstanceRoleARN: "arn::DUMMY::DUMMYROLE",
				},
			},
			expectedNewRole:     false,
			expectedNodeRoleARN: gfnt.NewString("arn::DUMMY::DUMMYROLE"), // using the provided role
		},
	}

	for i, tt := range nodeRoleTests {
		t.Run(fmt.Sprintf("%d: %s", i, tt.description), func(t *testing.T) {
			assert := assert.New(t)
			stack := NewManagedNodeGroup(api.NewClusterConfig(), tt.nodeGroup, "iam-test")
			err := stack.AddAllResources()
			assert.NoError(err)

			bytes, err := stack.RenderJSON()
			assert.NoError(err)

			template, err := goformation.ParseJSON(bytes)
			assert.NoError(err)
			ngResource, ok := template.Resources["ManagedNodeGroup"]
			assert.True(ok)
			ng, ok := ngResource.(*gfneks.Nodegroup)
			assert.True(ok)
			assert.Equal(tt.expectedNodeRoleARN, ng.NodeRole)

			_, ok = template.GetAllIAMRoleResources()[cfnIAMInstanceRoleName]
			assert.Equal(tt.expectedNewRole, ok)
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
