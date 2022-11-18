package v1alpha5_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var _ = Describe("Outposts validation", func() {
	type outpostsEntry struct {
		updateDefaultConfig func(*api.ClusterConfig)

		expectedErr string
	}

	DescribeTable("unsupported ClusterConfig features", func(oe outpostsEntry) {
		clusterConfig := api.NewClusterConfig()
		clusterConfig.Metadata.Version = api.Version1_21
		clusterConfig.Outpost = &api.Outpost{
			ControlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
		}
		oe.updateDefaultConfig(clusterConfig)
		err := api.ValidateClusterConfig(clusterConfig)
		Expect(err).To(MatchError(ContainSubstring(oe.expectedErr)))
	},
		Entry("fully-private cluster", outpostsEntry{
			updateDefaultConfig: func(c *api.ClusterConfig) {
				c.PrivateCluster = &api.PrivateCluster{
					Enabled: true,
				}
			},
			expectedErr: "fully-private cluster (privateCluster.enabled) is not supported for Outposts",
		}),

		Entry("Addons", outpostsEntry{
			updateDefaultConfig: func(c *api.ClusterConfig) {
				c.Addons = []*api.Addon{
					{
						Name: "vpc-cni",
					},
				}
			},
			expectedErr: "Addons are not supported on Outposts",
		}),

		Entry("Identity Providers", outpostsEntry{
			updateDefaultConfig: func(c *api.ClusterConfig) {
				c.IdentityProviders = []api.IdentityProvider{
					{
						Inner: &api.OIDCIdentityProvider{
							Name:           "pool-1",
							IssuerURL:      "url",
							ClientID:       "id",
							UsernameClaim:  "usernameClaim",
							UsernamePrefix: "usernamePrefix",
							GroupsClaim:    "groupsClaim",
							GroupsPrefix:   "groupsPrefix",
							RequiredClaims: map[string]string{"permission": "true"},
							Tags:           map[string]string{"department": "a"},
						}},
				}
			},

			expectedErr: "Identity Providers are not supported on Outposts",
		}),

		Entry("Fargate", outpostsEntry{
			updateDefaultConfig: func(c *api.ClusterConfig) {
				c.FargateProfiles = []*api.FargateProfile{
					{
						Name: "test",
					},
				}
			},

			expectedErr: "Fargate is not supported on Outposts",
		}),

		Entry("Karpenter", outpostsEntry{
			updateDefaultConfig: func(c *api.ClusterConfig) {
				c.Karpenter = &api.Karpenter{
					Version: "1.0.0",
				}
			},

			expectedErr: "Karpenter is not supported on Outposts",
		}),

		Entry("KMS encryption", outpostsEntry{
			updateDefaultConfig: func(c *api.ClusterConfig) {
				c.SecretsEncryption = &api.SecretsEncryption{
					KeyARN: "arn:aws:kms:us-west-2:000000000000:key/12345-12345",
				}
			},

			expectedErr: "KMS encryption is not supported on Outposts",
		}),

		Entry("Availability Zones", outpostsEntry{
			updateDefaultConfig: func(c *api.ClusterConfig) {
				c.AvailabilityZones = []string{"us-west-2a", "us-west-2b"}
			},

			expectedErr: "cannot specify availabilityZones on Outposts; the AZ defaults to the Outpost AZ",
		}),

		Entry("Local Zones", outpostsEntry{
			updateDefaultConfig: func(c *api.ClusterConfig) {
				c.LocalZones = []string{"us-west-2-lax-1a", "us-west-lax-1b"}
			},
			expectedErr: "cannot specify localZones on Outposts; the AZ defaults to the Outpost AZ",
		}),

		Entry("IPv6", outpostsEntry{
			updateDefaultConfig: func(c *api.ClusterConfig) {
				c.KubernetesNetworkConfig = &api.KubernetesNetworkConfig{
					IPFamily: api.IPV6Family,
				}
				c.IAM = &api.ClusterIAM{
					WithOIDC: api.Enabled(),
				}
				c.Addons = []*api.Addon{
					{
						Name: "vpc-cni",
					},
					{
						Name: "coredns",
					},
					{
						Name: "kube-proxy",
					},
				}
			},

			expectedErr: "IPv6 is not supported on Outposts",
		}),

		Entry("GitOps", outpostsEntry{
			updateDefaultConfig: func(c *api.ClusterConfig) {
				c.GitOps = &api.GitOps{}
			},

			expectedErr: "GitOps is not supported on Outposts",
		}),

		Entry("iam.withOIDC", outpostsEntry{
			updateDefaultConfig: func(c *api.ClusterConfig) {
				c.IAM = &api.ClusterIAM{
					WithOIDC: api.Enabled(),
				}
			},

			expectedErr: "iam.withOIDC is not supported on Outposts",
		}),

		Entry("nodeGroup.outpostARN set in a fully-private cluster", outpostsEntry{
			updateDefaultConfig: func(c *api.ClusterConfig) {
				c.Outpost = nil
				c.PrivateCluster = &api.PrivateCluster{
					Enabled: true,
				}
				ng := api.NewNodeGroup()
				ng.PrivateNetworking = true
				ng.Name = "test"
				ng.OutpostARN = "arn:aws:outposts:us-west-2:1234:outpost/op-1234"
				c.NodeGroups = []*api.NodeGroup{ng}
			},

			expectedErr: "nodeGroup.outpostARN is not supported on a fully-private cluster (privateCluster.enabled)",
		}),
	)

	DescribeTable("support for node AMI families", func(amiFamily string, shouldFail bool) {
		clusterConfig := api.NewClusterConfig()
		clusterConfig.Metadata.Version = api.Version1_21
		clusterConfig.Outpost = &api.Outpost{
			ControlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
		}
		err := api.ValidateNodeGroup(0, &api.NodeGroup{
			NodeGroupBase: &api.NodeGroupBase{
				Name:      "test",
				AMIFamily: amiFamily,
			},
		}, clusterConfig)
		if shouldFail {
			Expect(err).To(MatchError("only AmazonLinux2 is supported on local clusters"))
		} else {
			Expect(err).NotTo(HaveOccurred())
		}
	},
		Entry("AmazonLinux2", api.NodeImageFamilyAmazonLinux2, false),
		Entry("Bottlerocket", api.NodeImageFamilyBottlerocket, true),
		Entry("Ubuntu1804", api.NodeImageFamilyUbuntu1804, true),
		Entry("Ubuntu2004", api.NodeImageFamilyUbuntu2004, true),
		Entry("Windows2019Core", api.NodeImageFamilyWindowsServer2019CoreContainer, true),
		Entry("Windows2019Full", api.NodeImageFamilyWindowsServer2019FullContainer, true),
		Entry("Windows2022Core", api.NodeImageFamilyWindowsServer2022CoreContainer, true),
		Entry("Windows2022Full", api.NodeImageFamilyWindowsServer2022FullContainer, true),
	)

	type nodeGroupEntry struct {
		outpostInfo *api.ClusterConfig
		nodeGroup   *api.NodeGroup

		expectedErr string
	}

	DescribeTable("invalid nodegroup config", func(oe nodeGroupEntry) {
		defaultOutpostInfo := &api.ClusterConfig{
			Outpost: &api.Outpost{
				ControlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
			},
		}
		if oe.outpostInfo == nil {
			oe.outpostInfo = defaultOutpostInfo
		}
		err := api.ValidateNodeGroup(0, oe.nodeGroup, oe.outpostInfo)
		Expect(err).To(MatchError(ContainSubstring(oe.expectedErr)))
	},
		Entry("invalid Outpost ARN", nodeGroupEntry{
			nodeGroup: &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					OutpostARN: "arn:invalid",
				},
			},
			expectedErr: "invalid Outpost ARN",
		}),

		Entry("invalid service in Outpost ARN", nodeGroupEntry{
			nodeGroup: &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					OutpostARN: "arn:aws:eks:us-west-2:1234:eks/eks-1234",
				},
			},
			expectedErr: "invalid Outpost ARN",
		}),

		Entry("Outpost ARN does not match control plane's Outpost ARN", nodeGroupEntry{
			nodeGroup: &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					OutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
				},
			},
			outpostInfo: &api.ClusterConfig{
				Outpost: &api.Outpost{
					ControlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-5678",
				},
			},
			expectedErr: `nodeGroup.outpostARN must either be empty or match the control plane's Outpost ARN ("arn:aws:outposts:us-west-2:1234:outpost/op-1234" != "arn:aws:outposts:us-west-2:1234:outpost/op-5678")`,
		}),

		Entry("instanceSelector", nodeGroupEntry{
			nodeGroup: &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					InstanceSelector: &api.InstanceSelector{
						VCPUs: 4,
					},
				},
			},
			expectedErr: "cannot specify instanceSelector for a nodegroup on Outposts",
		}),

		Entry("availabilityZones", nodeGroupEntry{
			nodeGroup: &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					AvailabilityZones: []string{"us-west-2a", "us-west-2b"},
				},
			},
			expectedErr: "availabilityZones cannot be specified for a nodegroup on Outposts; the AZ defaults to the Outpost AZ",
		}),

		Entry("localZones", nodeGroupEntry{
			nodeGroup: &api.NodeGroup{
				LocalZones:    []string{"us-west-2a", "us-west-2b"},
				NodeGroupBase: &api.NodeGroupBase{},
			},
			expectedErr: "localZones cannot be specified for a nodegroup on Outposts; the AZ defaults to the Outpost AZ",
		}),
	)

	DescribeTable("unsupported volumeType", func(volumeType string, expectedErr bool) {
		ng := newNodeGroup()
		ng.VolumeType = &volumeType
		clusterConfig := api.NewClusterConfig()
		clusterConfig.Outpost = &api.Outpost{
			ControlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
		}
		err := api.ValidateNodeGroup(0, ng, clusterConfig)
		if expectedErr {
			Expect(err).To(MatchError(ContainSubstring(fmt.Sprintf(`cannot set %q for nodeGroups[0].volumeType; only %q volume types are supported on Outposts`, volumeType, api.NodeVolumeTypeGP2))))
		} else {
			Expect(err).NotTo(HaveOccurred())
		}
	},
		Entry(api.NodeVolumeTypeGP3, api.NodeVolumeTypeGP3, true),
		Entry(api.NodeVolumeTypeIO1, api.NodeVolumeTypeIO1, true),
		Entry(api.NodeVolumeTypeSC1, api.NodeVolumeTypeSC1, true),
		Entry(api.NodeVolumeTypeST1, api.NodeVolumeTypeST1, true),
		Entry(api.NodeVolumeTypeGP2, api.NodeVolumeTypeGP2, false),
	)
})
