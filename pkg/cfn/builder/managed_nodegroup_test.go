package builder

import (
	"fmt"
	"testing"

	"github.com/awslabs/goformation/v4"
	"github.com/stretchr/testify/assert"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

func TestManagedPolicyResources(t *testing.T) {
	iamRoleTests := []struct {
		addons                  api.NodeGroupIAMAddonPolicies
		attachPolicyARNs        []string
		expectedManagedPolicies []string
		description             string
	}{
		{
			expectedManagedPolicies: []string{"AmazonEKSWorkerNodePolicy", "AmazonEKS_CNI_Policy", "AmazonEC2ContainerRegistryReadOnly"},
			description:             "Default policies",
		},
		{
			addons: api.NodeGroupIAMAddonPolicies{
				ImageBuilder: api.Enabled(),
			},
			expectedManagedPolicies: []string{"AmazonEKSWorkerNodePolicy", "AmazonEKS_CNI_Policy", "AmazonEC2ContainerRegistryReadOnly", "AmazonEC2ContainerRegistryPowerUser"},
			description:             "ImageBuilder enabled",
		},
		{
			addons: api.NodeGroupIAMAddonPolicies{
				CloudWatch: api.Enabled(),
			},
			expectedManagedPolicies: []string{"AmazonEKSWorkerNodePolicy", "AmazonEKS_CNI_Policy", "AmazonEC2ContainerRegistryReadOnly", "CloudWatchAgentServerPolicy"},
			description:             "CloudWatch enabled",
		},
		{
			attachPolicyARNs:        []string{"AmazonEKSWorkerNodePolicy", "AmazonEKS_CNI_Policy"},
			expectedManagedPolicies: []string{"AmazonEKSWorkerNodePolicy", "AmazonEKS_CNI_Policy"},
			description:             "Custom policies",
		},
		// should not attach any additional policies
		{
			attachPolicyARNs:        []string{"CloudWatchAgentServerPolicy"},
			expectedManagedPolicies: []string{"CloudWatchAgentServerPolicy"},
			description:             "Custom policies",
		},
		// no duplicate values
		{
			attachPolicyARNs: []string{"AmazonEC2ContainerRegistryPowerUser"},
			addons: api.NodeGroupIAMAddonPolicies{
				ImageBuilder: api.Enabled(),
			},
			expectedManagedPolicies: []string{"AmazonEC2ContainerRegistryPowerUser"},
			description:             "Duplicate policies",
		},
		{
			attachPolicyARNs: []string{"CloudWatchAgentServerPolicy", "AmazonEC2ContainerRegistryPowerUser"},
			addons: api.NodeGroupIAMAddonPolicies{
				ImageBuilder: api.Enabled(),
				CloudWatch:   api.Enabled(),
			},
			expectedManagedPolicies: []string{"CloudWatchAgentServerPolicy", "AmazonEC2ContainerRegistryPowerUser"},
			description:             "Multiple duplicate policies",
		},
	}

	for i, tt := range iamRoleTests {
		t.Run(fmt.Sprintf("%d: %s", i, tt.description), func(t *testing.T) {
			assert := assert.New(t)
			clusterConfig := api.NewClusterConfig()

			ng := api.NewManagedNodeGroup()
			ng.IAM.WithAddonPolicies = tt.addons
			ng.IAM.AttachPolicyARNs = prefixPolicies(tt.attachPolicyARNs)

			stack := NewManagedNodeGroup(clusterConfig, ng, "iam-test")
			err := stack.AddAllResources()
			assert.NoError(err)

			bytes, err := stack.RenderJSON()
			assert.NoError(err)

			template, err := goformation.ParseJSON(bytes)
			assert.NoError(err)

			role, ok := template.GetAllIAMRoleResources()["NodeInstanceRole"]
			assert.True(ok)

			assert.ElementsMatch(prefixPolicies(tt.expectedManagedPolicies), role.ManagedPolicyArns)

		})
	}

}

func TestManagedNodeRole(t *testing.T) {
	nodeRoleTests := []struct {
		description      string
		nodeGroup        *api.ManagedNodeGroup
		expectedNodeRole string
	}{
		{
			description: "InstanceRoleARN is not provided",
			nodeGroup: &api.ManagedNodeGroup{
				ScalingConfig: &api.ScalingConfig{},
				SSH: &api.NodeGroupSSH{
					Allow: api.Disabled(),
				},
				IAM: &api.NodeGroupIAM{
				},
			},
			expectedNodeRole: "NodeInstanceRole", // creating new role
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
			expectedNodeRole: "arn::DUMMY::DUMMYROLE", // using the provided role
		},
	}

	for i, tt := range nodeRoleTests {
		t.Run(fmt.Sprintf("%d: %s", i, tt.description), func(t *testing.T) {
			stack := NewManagedNodeGroup(api.NewClusterConfig(), tt.nodeGroup, "iam-test")
			err := stack.AddAllResources()
			assert.NoError(t, err)

			bytes, err := stack.RenderJSON()
			assert.NoError(t, err)

			_, err = goformation.ParseJSON(bytes)
			assert.NoError(t, err)
			assert.Contains(t, string(bytes), tt.expectedNodeRole)
		})
	}
}

func prefixPolicies(policies []string) []string {
	var prefixedPolicies []string
	for _, policy := range policies {
		prefixedPolicies = append(prefixedPolicies, "arn:aws:iam::aws:policy/"+policy)
	}
	return prefixedPolicies
}
