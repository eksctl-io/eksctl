package addon_test

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/weaveworks/eksctl/pkg/actions/addon"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("Get", func() {
	var (
		manager            *addon.Manager
		mockProvider       *mockprovider.MockProvider
		describeAddonInput *awseks.DescribeAddonInput
	)
	BeforeEach(func() {
		var err error
		mockProvider = mockprovider.NewMockProvider()
		manager, err = addon.New(&api.ClusterConfig{Metadata: &api.ClusterMeta{
			Version: "1.18",
			Name:    "my-cluster",
		}}, mockProvider.EKS(), nil, false, nil, nil)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Get", func() {
		mockDescribeAddon := func(podIdentityAssociationIDs ...string) {
			mockProvider.MockEKS().On("DescribeAddonVersions", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(2))
				Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeAddonVersionsInput{}))
			}).Return(&awseks.DescribeAddonVersionsOutput{
				Addons: []ekstypes.AddonInfo{
					{
						AddonName: aws.String("my-addon"),
						Type:      aws.String("type"),
						AddonVersions: []ekstypes.AddonVersionInfo{
							{
								AddonVersion: aws.String("v1.0.0-eksbuild.1"),
							},
							{
								AddonVersion: aws.String("v1.1.0-eksbuild.1"),
							},
							{
								AddonVersion: aws.String("v1.1.0-eksbuild.4"),
							},
							{
								AddonVersion: aws.String("v1.2.0-eksbuild.1"),
							},
						},
					},
				},
			}, nil)

			mockProvider.MockEKS().On("DescribeAddon", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(2))
				Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeAddonInput{}))
				describeAddonInput = args[1].(*awseks.DescribeAddonInput)
			}).Return(&awseks.DescribeAddonOutput{
				Addon: &ekstypes.Addon{
					AddonName:               aws.String("my-addon"),
					AddonVersion:            aws.String("v1.1.0-eksbuild.1"),
					ServiceAccountRoleArn:   aws.String("foo"),
					Status:                  "created",
					PodIdentityAssociations: podIdentityAssociationIDs,
					Health: &ekstypes.AddonHealth{
						Issues: []ekstypes.AddonIssue{
							{
								Code:        "1",
								Message:     aws.String("foo"),
								ResourceIds: []string{"id-1"},
							},
						},
					},
				},
			}, nil)
			if len(podIdentityAssociationIDs) > 0 {
				for _, piaID := range podIdentityAssociationIDs {
					mockProvider.MockEKS().On("DescribePodIdentityAssociation", mock.Anything, &awseks.DescribePodIdentityAssociationInput{
						AssociationId: aws.String(piaID),
					}).Return(&awseks.DescribePodIdentityAssociationOutput{
						Association: &ekstypes.PodIdentityAssociation{
							AssociationId: aws.String(piaID),
						},
					}, nil).Once()
				}
			}
		}
		It("returns an addon", func() {
			mockDescribeAddon()
			summary, err := manager.Get(context.Background(), &api.Addon{
				Name: "my-addon",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(summary).To(Equal(addon.Summary{
				Name:         "my-addon",
				Version:      "v1.1.0-eksbuild.1",
				NewerVersion: "v1.1.0-eksbuild.4,v1.2.0-eksbuild.1",
				IAMRole:      "foo",
				Status:       "created",
				Issues: []addon.Issue{
					{
						Code:        "1",
						Message:     "foo",
						ResourceIDs: []string{"id-1"},
					},
				},
			}))

			Expect(*describeAddonInput.ClusterName).To(Equal("my-cluster"))
			Expect(*describeAddonInput.AddonName).To(Equal("my-addon"))
		})

		It("returns an addon with pod identity associations", func() {
			mockDescribeAddon("arn:aws:eks:us-west-2:00:podidentityassociation/cluster/a-zkgxwyqoexvjka9a3")
			mockProvider.MockEKS().On("DescribePodIdentityAssociation", mock.Anything, &awseks.DescribePodIdentityAssociationInput{
				AssociationId: aws.String("a-zkgxwyqoexvjka9a3"),
				ClusterName:   aws.String("my-cluster"),
			}).Return(&awseks.DescribePodIdentityAssociationOutput{
				Association: &ekstypes.PodIdentityAssociation{
					RoleArn:        aws.String("role-1"),
					ServiceAccount: aws.String("default"),
					Namespace:      aws.String("default"),
					AssociationId:  aws.String("a-1"),
				},
			}, nil)
			summary, err := manager.Get(context.Background(), &api.Addon{
				Name: "my-addon",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(summary).To(Equal(addon.Summary{
				Name:         "my-addon",
				Version:      "v1.1.0-eksbuild.1",
				NewerVersion: "v1.1.0-eksbuild.4,v1.2.0-eksbuild.1",
				IAMRole:      "foo",
				Status:       "created",
				Issues: []addon.Issue{
					{
						Code:        "1",
						Message:     "foo",
						ResourceIDs: []string{"id-1"},
					},
				},
				PodIdentityAssociations: []addon.PodIdentityAssociationSummary{
					{
						Namespace:      "default",
						ServiceAccount: "default",
						RoleARN:        "role-1",
						AssociationID:  "a-1",
					},
				},
			}))
		})

		It("returns an addon with namespace config", func() {
			mockProvider.MockEKS().On("DescribeAddonVersions", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(2))
				Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeAddonVersionsInput{}))
			}).Return(&awseks.DescribeAddonVersionsOutput{
				Addons: []ekstypes.AddonInfo{
					{
						AddonName: aws.String("my-addon"),
						Type:      aws.String("type"),
						AddonVersions: []ekstypes.AddonVersionInfo{
							{
								AddonVersion: aws.String("v1.0.0-eksbuild.1"),
							},
							{
								AddonVersion: aws.String("v1.1.0-eksbuild.1"),
							},
						},
					},
				},
			}, nil)

			mockProvider.MockEKS().On("DescribeAddon", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(2))
				Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeAddonInput{}))
				describeAddonInput = args[1].(*awseks.DescribeAddonInput)
			}).Return(&awseks.DescribeAddonOutput{
				Addon: &ekstypes.Addon{
					AddonName:    aws.String("my-addon"),
					AddonVersion: aws.String("v1.1.0-eksbuild.1"),
					Status:       "created",
					NamespaceConfig: &ekstypes.AddonNamespaceConfigResponse{
						Namespace: aws.String("custom-namespace"),
					},
				},
			}, nil)

			summary, err := manager.Get(context.Background(), &api.Addon{
				Name: "my-addon",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(summary).To(Equal(addon.Summary{
				Name:         "my-addon",
				Version:      "v1.1.0-eksbuild.1",
				NewerVersion: "",
				Status:       "created",
				NamespaceConfig: &api.AddonNamespaceConfig{
					Namespace: "custom-namespace",
				},
			}))

			Expect(*describeAddonInput.ClusterName).To(Equal("my-cluster"))
			Expect(*describeAddonInput.AddonName).To(Equal("my-addon"))
		})

		It("returns an addon with nil namespace config when not set", func() {
			mockProvider.MockEKS().On("DescribeAddonVersions", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(2))
				Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeAddonVersionsInput{}))
			}).Return(&awseks.DescribeAddonVersionsOutput{
				Addons: []ekstypes.AddonInfo{
					{
						AddonName: aws.String("my-addon"),
						Type:      aws.String("type"),
						AddonVersions: []ekstypes.AddonVersionInfo{
							{
								AddonVersion: aws.String("v1.0.0-eksbuild.1"),
							},
							{
								AddonVersion: aws.String("v1.1.0-eksbuild.1"),
							},
						},
					},
				},
			}, nil)

			mockProvider.MockEKS().On("DescribeAddon", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(2))
				Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeAddonInput{}))
				describeAddonInput = args[1].(*awseks.DescribeAddonInput)
			}).Return(&awseks.DescribeAddonOutput{
				Addon: &ekstypes.Addon{
					AddonName:       aws.String("my-addon"),
					AddonVersion:    aws.String("v1.1.0-eksbuild.1"),
					Status:          "created",
					NamespaceConfig: nil, // No namespace config
				},
			}, nil)

			summary, err := manager.Get(context.Background(), &api.Addon{
				Name: "my-addon",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(summary).To(Equal(addon.Summary{
				Name:            "my-addon",
				Version:         "v1.1.0-eksbuild.1",
				NewerVersion:    "",
				Status:          "created",
				NamespaceConfig: nil, // Should be nil when not set
			}))

			Expect(*describeAddonInput.ClusterName).To(Equal("my-cluster"))
			Expect(*describeAddonInput.AddonName).To(Equal("my-addon"))
		})

		It("returns an addon with namespace config containing empty namespace", func() {
			mockProvider.MockEKS().On("DescribeAddonVersions", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(2))
				Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeAddonVersionsInput{}))
			}).Return(&awseks.DescribeAddonVersionsOutput{
				Addons: []ekstypes.AddonInfo{
					{
						AddonName: aws.String("my-addon"),
						Type:      aws.String("type"),
						AddonVersions: []ekstypes.AddonVersionInfo{
							{
								AddonVersion: aws.String("v1.0.0-eksbuild.1"),
							},
						},
					},
				},
			}, nil)

			mockProvider.MockEKS().On("DescribeAddon", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(2))
				Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeAddonInput{}))
				describeAddonInput = args[1].(*awseks.DescribeAddonInput)
			}).Return(&awseks.DescribeAddonOutput{
				Addon: &ekstypes.Addon{
					AddonName:    aws.String("my-addon"),
					AddonVersion: aws.String("v1.0.0-eksbuild.1"),
					Status:       "created",
					NamespaceConfig: &ekstypes.AddonNamespaceConfigResponse{
						Namespace: aws.String(""), // Empty namespace
					},
				},
			}, nil)

			summary, err := manager.Get(context.Background(), &api.Addon{
				Name: "my-addon",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(summary).To(Equal(addon.Summary{
				Name:         "my-addon",
				Version:      "v1.0.0-eksbuild.1",
				NewerVersion: "",
				Status:       "created",
				NamespaceConfig: &api.AddonNamespaceConfig{
					Namespace: "", // Should preserve empty namespace
				},
			}))

			Expect(*describeAddonInput.ClusterName).To(Equal("my-cluster"))
			Expect(*describeAddonInput.AddonName).To(Equal("my-addon"))
		})

		It("returns an addon with namespace config when AWS response has nil namespace", func() {
			mockProvider.MockEKS().On("DescribeAddonVersions", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(2))
				Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeAddonVersionsInput{}))
			}).Return(&awseks.DescribeAddonVersionsOutput{
				Addons: []ekstypes.AddonInfo{
					{
						AddonName: aws.String("my-addon"),
						Type:      aws.String("type"),
						AddonVersions: []ekstypes.AddonVersionInfo{
							{
								AddonVersion: aws.String("v1.0.0-eksbuild.1"),
							},
						},
					},
				},
			}, nil)

			mockProvider.MockEKS().On("DescribeAddon", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(2))
				Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeAddonInput{}))
				describeAddonInput = args[1].(*awseks.DescribeAddonInput)
			}).Return(&awseks.DescribeAddonOutput{
				Addon: &ekstypes.Addon{
					AddonName:    aws.String("my-addon"),
					AddonVersion: aws.String("v1.0.0-eksbuild.1"),
					Status:       "created",
					NamespaceConfig: &ekstypes.AddonNamespaceConfigResponse{
						Namespace: nil, // Nil namespace pointer
					},
				},
			}, nil)

			summary, err := manager.Get(context.Background(), &api.Addon{
				Name: "my-addon",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(summary).To(Equal(addon.Summary{
				Name:            "my-addon",
				Version:         "v1.0.0-eksbuild.1",
				NewerVersion:    "",
				Status:          "created",
				NamespaceConfig: nil, // Should be nil when AWS namespace is nil
			}))

			Expect(*describeAddonInput.ClusterName).To(Equal("my-cluster"))
			Expect(*describeAddonInput.AddonName).To(Equal("my-addon"))
		})

		It("returns an addon with namespace config containing valid DNS-1123 namespace", func() {
			mockProvider.MockEKS().On("DescribeAddonVersions", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(2))
				Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeAddonVersionsInput{}))
			}).Return(&awseks.DescribeAddonVersionsOutput{
				Addons: []ekstypes.AddonInfo{
					{
						AddonName: aws.String("my-addon"),
						Type:      aws.String("type"),
						AddonVersions: []ekstypes.AddonVersionInfo{
							{
								AddonVersion: aws.String("v1.0.0-eksbuild.1"),
							},
						},
					},
				},
			}, nil)

			mockProvider.MockEKS().On("DescribeAddon", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(2))
				Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeAddonInput{}))
				describeAddonInput = args[1].(*awseks.DescribeAddonInput)
			}).Return(&awseks.DescribeAddonOutput{
				Addon: &ekstypes.Addon{
					AddonName:    aws.String("my-addon"),
					AddonVersion: aws.String("v1.0.0-eksbuild.1"),
					Status:       "created",
					NamespaceConfig: &ekstypes.AddonNamespaceConfigResponse{
						Namespace: aws.String("kube-system-addon-ns"), // Valid DNS-1123 name with hyphens
					},
				},
			}, nil)

			summary, err := manager.Get(context.Background(), &api.Addon{
				Name: "my-addon",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(summary).To(Equal(addon.Summary{
				Name:         "my-addon",
				Version:      "v1.0.0-eksbuild.1",
				NewerVersion: "",
				Status:       "created",
				NamespaceConfig: &api.AddonNamespaceConfig{
					Namespace: "kube-system-addon-ns",
				},
			}))

			Expect(*describeAddonInput.ClusterName).To(Equal("my-cluster"))
			Expect(*describeAddonInput.AddonName).To(Equal("my-addon"))
		})

		When("it fails to get the addon", func() {
			It("returns an error", func() {
				mockProvider.MockEKS().On("DescribeAddon", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					Expect(args).To(HaveLen(2))
					Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeAddonInput{}))
					describeAddonInput = args[1].(*awseks.DescribeAddonInput)
				}).Return(nil, fmt.Errorf("foo"))

				_, err := manager.Get(context.Background(), &api.Addon{
					Name: "my-addon",
				})
				Expect(err).To(MatchError(`failed to get addon "my-addon": foo`))
				Expect(*describeAddonInput.ClusterName).To(Equal("my-cluster"))
				Expect(*describeAddonInput.AddonName).To(Equal("my-addon"))
			})
		})
	})

	Describe("GetAll", func() {
		var listAddonsInput *awseks.ListAddonsInput
		It("returns an addon", func() {
			mockProvider.MockEKS().On("DescribeAddonVersions", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(2))
				Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeAddonVersionsInput{}))
			}).Return(&awseks.DescribeAddonVersionsOutput{
				Addons: []ekstypes.AddonInfo{
					{
						AddonName: aws.String("my-addon"),
						Type:      aws.String("type"),
						AddonVersions: []ekstypes.AddonVersionInfo{
							{
								AddonVersion: aws.String("v1.0.0-eksbuild.1"),
							},
							{
								AddonVersion: aws.String("v1.1.0-eksbuild.1"),
							},
							{
								AddonVersion: aws.String("v1.1.0-eksbuild.4"),
							},
							{
								AddonVersion: aws.String("v1.2.0-eksbuild.1"),
							},
						},
					},
				},
			}, nil)

			mockProvider.MockEKS().On("ListAddons", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(2))
				Expect(args[1]).To(BeAssignableToTypeOf(&awseks.ListAddonsInput{}))
				listAddonsInput = args[1].(*awseks.ListAddonsInput)
			}).Return(&awseks.ListAddonsOutput{
				Addons: []string{"my-addon"},
			}, nil)

			mockProvider.MockEKS().On("DescribeAddon", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(2))
				Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeAddonInput{}))
				describeAddonInput = args[1].(*awseks.DescribeAddonInput)
			}).Return(&awseks.DescribeAddonOutput{
				Addon: &ekstypes.Addon{
					AddonName:             aws.String("my-addon"),
					AddonVersion:          aws.String("v1.1.0-eksbuild.1"),
					ServiceAccountRoleArn: aws.String("foo"),
					Status:                "created",
					ConfigurationValues:   aws.String("{\"replicaCount\":3}"),
				},
			}, nil)

			summary, err := manager.GetAll(context.Background())
			Expect(err).NotTo(HaveOccurred())
			Expect(summary).To(Equal([]addon.Summary{
				{
					Name:                "my-addon",
					Version:             "v1.1.0-eksbuild.1",
					NewerVersion:        "v1.1.0-eksbuild.4,v1.2.0-eksbuild.1",
					IAMRole:             "foo",
					Status:              "created",
					ConfigurationValues: "{\"replicaCount\":3}",
				},
			}))

			Expect(*describeAddonInput.ClusterName).To(Equal("my-cluster"))
			Expect(*describeAddonInput.AddonName).To(Equal("my-addon"))
			Expect(*listAddonsInput.ClusterName).To(Equal("my-cluster"))
		})

		It("returns addons with namespace config", func() {
			mockProvider.MockEKS().On("DescribeAddonVersions", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(2))
				Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeAddonVersionsInput{}))
			}).Return(&awseks.DescribeAddonVersionsOutput{
				Addons: []ekstypes.AddonInfo{
					{
						AddonName: aws.String("addon-with-namespace"),
						Type:      aws.String("type"),
						AddonVersions: []ekstypes.AddonVersionInfo{
							{
								AddonVersion: aws.String("v1.0.0-eksbuild.1"),
							},
						},
					},
					{
						AddonName: aws.String("addon-without-namespace"),
						Type:      aws.String("type"),
						AddonVersions: []ekstypes.AddonVersionInfo{
							{
								AddonVersion: aws.String("v1.0.0-eksbuild.1"),
							},
						},
					},
				},
			}, nil)

			mockProvider.MockEKS().On("ListAddons", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(2))
				Expect(args[1]).To(BeAssignableToTypeOf(&awseks.ListAddonsInput{}))
				listAddonsInput = args[1].(*awseks.ListAddonsInput)
			}).Return(&awseks.ListAddonsOutput{
				Addons: []string{"addon-with-namespace", "addon-without-namespace"},
			}, nil)

			// Mock first addon with namespace config
			mockProvider.MockEKS().On("DescribeAddon", mock.Anything, &awseks.DescribeAddonInput{
				ClusterName: aws.String("my-cluster"),
				AddonName:   aws.String("addon-with-namespace"),
			}).Return(&awseks.DescribeAddonOutput{
				Addon: &ekstypes.Addon{
					AddonName:    aws.String("addon-with-namespace"),
					AddonVersion: aws.String("v1.0.0-eksbuild.1"),
					Status:       "created",
					NamespaceConfig: &ekstypes.AddonNamespaceConfigResponse{
						Namespace: aws.String("monitoring"),
					},
				},
			}, nil)

			// Mock second addon without namespace config
			mockProvider.MockEKS().On("DescribeAddon", mock.Anything, &awseks.DescribeAddonInput{
				ClusterName: aws.String("my-cluster"),
				AddonName:   aws.String("addon-without-namespace"),
			}).Return(&awseks.DescribeAddonOutput{
				Addon: &ekstypes.Addon{
					AddonName:       aws.String("addon-without-namespace"),
					AddonVersion:    aws.String("v1.0.0-eksbuild.1"),
					Status:          "created",
					NamespaceConfig: nil,
				},
			}, nil)

			summary, err := manager.GetAll(context.Background())
			Expect(err).NotTo(HaveOccurred())
			Expect(summary).To(Equal([]addon.Summary{
				{
					Name:         "addon-with-namespace",
					Version:      "v1.0.0-eksbuild.1",
					NewerVersion: "",
					Status:       "created",
					NamespaceConfig: &api.AddonNamespaceConfig{
						Namespace: "monitoring",
					},
				},
				{
					Name:            "addon-without-namespace",
					Version:         "v1.0.0-eksbuild.1",
					NewerVersion:    "",
					Status:          "created",
					NamespaceConfig: nil,
				},
			}))

			Expect(*listAddonsInput.ClusterName).To(Equal("my-cluster"))
		})

		When("it fails to get the addon", func() {
			It("returns an error", func() {
				mockProvider.MockEKS().On("ListAddons", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					Expect(args).To(HaveLen(2))
					Expect(args[1]).To(BeAssignableToTypeOf(&awseks.ListAddonsInput{}))
					listAddonsInput = args[1].(*awseks.ListAddonsInput)
				}).Return(&awseks.ListAddonsOutput{
					Addons: []string{"my-addon"},
				}, nil)

				mockProvider.MockEKS().On("DescribeAddon", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					Expect(args).To(HaveLen(2))
					Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeAddonInput{}))
					describeAddonInput = args[1].(*awseks.DescribeAddonInput)
				}).Return(nil, fmt.Errorf("foo"))

				_, err := manager.GetAll(context.Background())
				Expect(err).To(MatchError(`failed to get addon "my-addon": foo`))
				Expect(*describeAddonInput.ClusterName).To(Equal("my-cluster"))
				Expect(*describeAddonInput.AddonName).To(Equal("my-addon"))
				Expect(*listAddonsInput.ClusterName).To(Equal("my-cluster"))

			})
		})

		When("it fails to list addons", func() {
			It("returns an error", func() {
				mockProvider.MockEKS().On("ListAddons", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					Expect(args).To(HaveLen(2))
					Expect(args[1]).To(BeAssignableToTypeOf(&awseks.ListAddonsInput{}))
					listAddonsInput = args[1].(*awseks.ListAddonsInput)
				}).Return(&awseks.ListAddonsOutput{
					Addons: []string{"my-addon"},
				}, fmt.Errorf("foo"))

				_, err := manager.GetAll(context.Background())
				Expect(err).To(MatchError(`failed to list addons: foo`))
				Expect(*listAddonsInput.ClusterName).To(Equal("my-cluster"))
			})
		})
	})
})
