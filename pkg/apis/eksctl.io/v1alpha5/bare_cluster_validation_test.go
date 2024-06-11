package v1alpha5_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go-v2/aws"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

type bareClusterEntry struct {
	updateClusterConfig func(*api.ClusterConfig)
	expectErr           bool
}

var _ = DescribeTable("Bare cluster validation", func(e bareClusterEntry) {
	clusterConfig := api.NewClusterConfig()
	clusterConfig.AddonsConfig.DisableDefaultAddons = true
	clusterConfig.Addons = []*api.Addon{
		{
			Name: api.CoreDNSAddon,
		},
	}
	e.updateClusterConfig(clusterConfig)
	err := api.ValidateClusterConfig(clusterConfig)
	if e.expectErr {
		Expect(err).To(MatchError("fields nodeGroups, managedNodeGroups, fargateProfiles, karpenter, gitops, iam.serviceAccounts, " +
			"and iam.podIdentityAssociations are not supported during cluster creation in a cluster without VPC CNI; please remove these fields " +
			"and add them back after cluster creation is successful"))
	} else {
		Expect(err).NotTo(HaveOccurred())
	}

},
	Entry("nodeGroups", bareClusterEntry{
		updateClusterConfig: func(c *api.ClusterConfig) {
			ng := api.NewNodeGroup()
			ng.Name = "ng"
			ng.DesiredCapacity = aws.Int(1)
			c.NodeGroups = []*api.NodeGroup{ng}
		},
		expectErr: true,
	}),
	Entry("managedNodeGroups", bareClusterEntry{
		updateClusterConfig: func(c *api.ClusterConfig) {
			ng := api.NewManagedNodeGroup()
			ng.Name = "mng"
			ng.DesiredCapacity = aws.Int(1)
			c.ManagedNodeGroups = []*api.ManagedNodeGroup{ng}
		},
		expectErr: true,
	}),
	Entry("fargateProfiles", bareClusterEntry{
		updateClusterConfig: func(c *api.ClusterConfig) {
			c.FargateProfiles = []*api.FargateProfile{
				{
					Name: "test",
					Selectors: []api.FargateProfileSelector{
						{
							Namespace: "default",
						},
					},
				},
			}
		},
		expectErr: true,
	}),
	Entry("gitops", bareClusterEntry{
		updateClusterConfig: func(c *api.ClusterConfig) {
			c.GitOps = &api.GitOps{
				Flux: &api.Flux{},
			}
		},
		expectErr: true,
	}),
	Entry("karpenter", bareClusterEntry{
		updateClusterConfig: func(c *api.ClusterConfig) {
			c.Karpenter = &api.Karpenter{}
		},
		expectErr: true,
	}),
	Entry("iam.serviceAccounts", bareClusterEntry{
		updateClusterConfig: func(c *api.ClusterConfig) {
			c.IAM.WithOIDC = api.Enabled()
			c.IAM.ServiceAccounts = []*api.ClusterIAMServiceAccount{
				{
					ClusterIAMMeta: api.ClusterIAMMeta{
						Name:      "test",
						Namespace: "test",
					},
					AttachPolicyARNs: []string{"arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"},
				},
			}
		},
		expectErr: true,
	}),
	Entry("iam.podIdentityAssociations", bareClusterEntry{
		updateClusterConfig: func(c *api.ClusterConfig) {
			c.IAM.PodIdentityAssociations = []api.PodIdentityAssociation{
				{
					Namespace:            "test",
					PermissionPolicyARNs: []string{"arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"},
				},
			}
		},
		expectErr: true,
	}),
	Entry("no unsupported field set", bareClusterEntry{
		updateClusterConfig: func(c *api.ClusterConfig) {},
	}),
)
