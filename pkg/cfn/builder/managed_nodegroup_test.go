package builder

import (
	"strconv"
	"testing"

	"github.com/awslabs/goformation/v4"
	"github.com/stretchr/testify/assert"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

func TestManagedResources(t *testing.T) {
	iamRoleTests := []struct {
		addons                  api.NodeGroupIAMAddonPolicies
		attachPolicyARNs        []string
		expectedManagedPolicies []string
	}{
		{
			expectedManagedPolicies: []string{"AmazonEKSWorkerNodePolicy", "AmazonEKS_CNI_Policy", "AmazonEC2ContainerRegistryReadOnly"},
		},
		{
			addons: api.NodeGroupIAMAddonPolicies{
				ImageBuilder: api.Enabled(),
			},
			expectedManagedPolicies: []string{"AmazonEKSWorkerNodePolicy", "AmazonEKS_CNI_Policy", "AmazonEC2ContainerRegistryReadOnly", "AmazonEC2ContainerRegistryPowerUser"},
		},
		{
			addons: api.NodeGroupIAMAddonPolicies{
				CloudWatch: api.Enabled(),
			},
			expectedManagedPolicies: []string{"AmazonEKSWorkerNodePolicy", "AmazonEKS_CNI_Policy", "AmazonEC2ContainerRegistryReadOnly", "CloudWatchAgentServerPolicy"},
		},
		{
			attachPolicyARNs:        []string{"AmazonEKSWorkerNodePolicy", "AmazonEKS_CNI_Policy"},
			expectedManagedPolicies: []string{"AmazonEKSWorkerNodePolicy", "AmazonEKS_CNI_Policy"},
		},
		// should not attach any additional policies
		{
			attachPolicyARNs:        []string{"CloudWatchAgentServerPolicy"},
			expectedManagedPolicies: []string{"CloudWatchAgentServerPolicy"},
		},
		// no duplicate values
		{
			attachPolicyARNs: []string{"AmazonEC2ContainerRegistryPowerUser"},
			addons: api.NodeGroupIAMAddonPolicies{
				ImageBuilder: api.Enabled(),
			},
			expectedManagedPolicies: []string{"AmazonEC2ContainerRegistryPowerUser"},
		},
		{
			attachPolicyARNs: []string{"CloudWatchAgentServerPolicy", "AmazonEC2ContainerRegistryPowerUser"},
			addons: api.NodeGroupIAMAddonPolicies{
				ImageBuilder: api.Enabled(),
				CloudWatch:   api.Enabled(),
			},
			expectedManagedPolicies: []string{"CloudWatchAgentServerPolicy", "AmazonEC2ContainerRegistryPowerUser"},
		},
	}

	for i, tt := range iamRoleTests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
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

func prefixPolicies(policies []string) []string {
	var prefixedPolicies []string
	for _, policy := range policies {
		prefixedPolicies = append(prefixedPolicies, "arn:aws:iam::aws:policy/"+policy)
	}
	return prefixedPolicies
}
