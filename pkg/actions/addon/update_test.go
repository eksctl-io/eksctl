package addon_test

import (
	"bytes"
	"context"
	"fmt"
	"time"

	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/kris-nova/logger"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/weaveworks/eksctl/pkg/actions/addon"
	"github.com/weaveworks/eksctl/pkg/actions/addon/fakes"
	"github.com/weaveworks/eksctl/pkg/actions/addon/mocks"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("Update", func() {
	var (
		addonManager       *addon.Manager
		mockProvider       *mockprovider.MockProvider
		updateAddonInput   *awseks.UpdateAddonInput
		describeAddonInput *awseks.DescribeAddonInput
		fakeStackManager   *fakes.FakeStackManager
		waitTimeout        = 5 * time.Minute
	)

	makeOIDCManager := func() *iamoidc.OpenIDConnectManager {
		oidc, err := iamoidc.NewOpenIDConnectManager(nil, "456123987123", "https://oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E", "aws", nil)
		Expect(err).NotTo(HaveOccurred())
		oidc.ProviderARN = "arn:aws:iam::456123987123:oidc-provider/oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E"
		return oidc
	}

	BeforeEach(func() {
		var err error
		mockProvider = mockprovider.NewMockProvider()
		fakeStackManager = new(fakes.FakeStackManager)

		fakeStackManager.CreateStackStub = func(_ context.Context, _ string, rs builder.ResourceSetReader, _ map[string]string, _ map[string]string, errs chan error) error {
			go func() {
				errs <- nil
			}()
			Expect(rs).To(BeAssignableToTypeOf(&builder.IAMRoleResourceSet{}))
			rs.(*builder.IAMRoleResourceSet).OutputRole = "new-service-account-role-arn"
			return nil
		}

		oidc := makeOIDCManager()

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
							AddonVersion: aws.String("v1.7.5-eksbuild.1"),
						},
						{
							AddonVersion: aws.String("v1.7.5-eksbuild.2"),
						},
						{
							// not sure if all versions come with v prefix or not, so test a mix.
							AddonVersion: aws.String("v1.7.7-eksbuild.2"),
						},
						{
							AddonVersion: aws.String("v1.7.6"),
						},
						{
							AddonVersion: aws.String("v1.0.0-eksbuild.2"),
						},
					},
				},
			},
		}, nil)

		mockProvider.MockEKS().On("DescribeAddon", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			Expect(args).To(HaveLen(2))
			Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeAddonInput{}))
			describeAddonInput = args[1].(*awseks.DescribeAddonInput)
		}).Return(&awseks.DescribeAddonOutput{
			Addon: &ekstypes.Addon{
				AddonName:             aws.String("my-addon"),
				AddonVersion:          aws.String("v1.0.0-eksbuild.2"),
				ServiceAccountRoleArn: aws.String("original-arn"),
				Status:                "created",
			},
		}, nil).Once()

		addonManager, err = addon.New(&api.ClusterConfig{Metadata: &api.ClusterMeta{
			Version: "1.18",
			Name:    "my-cluster",
		}}, mockProvider.EKS(), fakeStackManager, true, oidc, nil)
		Expect(err).NotTo(HaveOccurred())
	})

	When("EKS returns an UpdateAddonOutput", func() {
		var podIdentityIAMUpdater mocks.PodIdentityIAMUpdater
		BeforeEach(func() {
			mockProvider.MockEKS().On("UpdateAddon", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(2))
				Expect(args[1]).To(BeAssignableToTypeOf(&awseks.UpdateAddonInput{}))
				updateAddonInput = args[1].(*awseks.UpdateAddonInput)
			}).Return(&awseks.UpdateAddonOutput{}, nil)
		})

		When("updating the version", func() {
			It("updates the addon and preserves the existing role", func() {
				err := addonManager.Update(context.Background(), &api.Addon{
					Name:    "my-addon",
					Version: "v1.0.0-eksbuild.2",
					Force:   true,
				}, &podIdentityIAMUpdater, 0)

				Expect(err).NotTo(HaveOccurred())
				Expect(*describeAddonInput.ClusterName).To(Equal("my-cluster"))
				Expect(*describeAddonInput.AddonName).To(Equal("my-addon"))
				Expect(*updateAddonInput.ClusterName).To(Equal("my-cluster"))
				Expect(*updateAddonInput.AddonName).To(Equal("my-addon"))
				Expect(*updateAddonInput.AddonVersion).To(Equal("v1.0.0-eksbuild.2"))
				Expect(*updateAddonInput.ServiceAccountRoleArn).To(Equal("original-arn"))
				Expect(updateAddonInput.ResolveConflicts).To(Equal(ekstypes.ResolveConflictsOverwrite))
			})

			When("the version is not set", func() {
				It("preserves the existing addon version", func() {
					output := &bytes.Buffer{}
					logger.Writer = output

					err := addonManager.Update(context.Background(), &api.Addon{
						Name:    "my-addon",
						Version: "",
					}, &podIdentityIAMUpdater, 0)

					Expect(err).NotTo(HaveOccurred())
					Expect(output.String()).To(ContainSubstring("no new version provided, preserving existing version"))
					Expect(*describeAddonInput.ClusterName).To(Equal("my-cluster"))
					Expect(*describeAddonInput.AddonName).To(Equal("my-addon"))
					Expect(*updateAddonInput.ClusterName).To(Equal("my-cluster"))
					Expect(*updateAddonInput.AddonName).To(Equal("my-addon"))
					Expect(*updateAddonInput.AddonVersion).To(Equal("v1.0.0-eksbuild.2"))
					Expect(*updateAddonInput.ServiceAccountRoleArn).To(Equal("original-arn"))
				})
			})

			When("the version is set to a numeric version", func() {
				It("discovers and uses the latest available version", func() {
					err := addonManager.Update(context.Background(), &api.Addon{
						Name:    "my-addon",
						Version: "1.7.5",
					}, &podIdentityIAMUpdater, 0)

					Expect(err).NotTo(HaveOccurred())
					Expect(*describeAddonInput.ClusterName).To(Equal("my-cluster"))
					Expect(*describeAddonInput.AddonName).To(Equal("my-addon"))
					Expect(*updateAddonInput.ClusterName).To(Equal("my-cluster"))
					Expect(*updateAddonInput.AddonName).To(Equal("my-addon"))
					Expect(*updateAddonInput.AddonVersion).To(Equal("v1.7.5-eksbuild.2"))
					Expect(*updateAddonInput.ServiceAccountRoleArn).To(Equal("original-arn"))
				})
			})

			When("the version is set to latest", func() {
				It("discovers and uses the latest available version", func() {
					err := addonManager.Update(context.Background(), &api.Addon{
						Name:    "my-addon",
						Version: "latest",
					}, &podIdentityIAMUpdater, 0)

					Expect(err).NotTo(HaveOccurred())
					Expect(*describeAddonInput.ClusterName).To(Equal("my-cluster"))
					Expect(*describeAddonInput.AddonName).To(Equal("my-addon"))
					Expect(*updateAddonInput.ClusterName).To(Equal("my-cluster"))
					Expect(*updateAddonInput.AddonName).To(Equal("my-addon"))
					Expect(*updateAddonInput.AddonVersion).To(Equal("v1.7.7-eksbuild.2"))
					Expect(*updateAddonInput.ServiceAccountRoleArn).To(Equal("original-arn"))
				})
			})

			When("the version is set to a version that does not exist", func() {
				It("returns an error", func() {
					err := addonManager.Update(context.Background(), &api.Addon{
						Name:             "my-addon",
						Version:          "1.7.8",
						AttachPolicyARNs: []string{"arn-1"},
					}, &podIdentityIAMUpdater, 0)
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(ContainSubstring("no version(s) found matching \"1.7.8\" for \"my-addon\"")))
				})
			})
		})

		When("wait is true", func() {
			When("the addon update succeeds", func() {
				BeforeEach(func() {
					mockProvider.MockEKS().On("DescribeAddon", mock.Anything, mock.Anything, mock.Anything).
						Return(&awseks.DescribeAddonOutput{
							Addon: &ekstypes.Addon{
								AddonName: aws.String("my-addon"),
								Status:    ekstypes.AddonStatusActive,
							},
						}, nil)
				})

				It("creates the addon and waits for it to be running", func() {
					err := addonManager.Update(context.Background(), &api.Addon{
						Name:    "my-addon",
						Version: "v1.0.0-eksbuild.2",
						Force:   true,
					}, &podIdentityIAMUpdater, waitTimeout)
					Expect(err).NotTo(HaveOccurred())
					Expect(*describeAddonInput.ClusterName).To(Equal("my-cluster"))
					Expect(*describeAddonInput.AddonName).To(Equal("my-addon"))
					Expect(*updateAddonInput.ClusterName).To(Equal("my-cluster"))
					Expect(*updateAddonInput.AddonName).To(Equal("my-addon"))
					Expect(*updateAddonInput.AddonVersion).To(Equal("v1.0.0-eksbuild.2"))
					Expect(*updateAddonInput.ServiceAccountRoleArn).To(Equal("original-arn"))
					Expect(updateAddonInput.ResolveConflicts).To(Equal(ekstypes.ResolveConflictsOverwrite))
				})
			})

			When("the addon update fails", func() {
				BeforeEach(func() {
					mockProvider.MockEKS().On("DescribeAddon", mock.Anything, mock.Anything, mock.Anything).
						Return(&awseks.DescribeAddonOutput{
							Addon: &ekstypes.Addon{
								AddonName: aws.String("my-addon"),
								Status:    ekstypes.AddonStatusDegraded,
							},
						}, nil)
				})

				It("returns an error", func() {
					err := addonManager.Update(context.Background(), &api.Addon{
						Name:    "my-addon",
						Version: "v1.0.0-eksbuild.2",
						Force:   true,
					}, &podIdentityIAMUpdater, waitTimeout)
					Expect(err).To(MatchError(`addon status transitioned to "DEGRADED"`))
				})
			})
		})

		When("updating the policy", func() {
			When("specifying a new serviceAccountRoleARN", func() {
				It("updates the addon", func() {
					err := addonManager.Update(context.Background(), &api.Addon{
						Name:                  "my-addon",
						Version:               "v1.0.0-eksbuild.2",
						ServiceAccountRoleARN: "new-arn",
					}, &podIdentityIAMUpdater, 0)

					Expect(err).NotTo(HaveOccurred())
					Expect(*describeAddonInput.ClusterName).To(Equal("my-cluster"))
					Expect(*describeAddonInput.AddonName).To(Equal("my-addon"))
					Expect(*updateAddonInput.ClusterName).To(Equal("my-cluster"))
					Expect(*updateAddonInput.AddonName).To(Equal("my-addon"))
					Expect(*updateAddonInput.AddonVersion).To(Equal("v1.0.0-eksbuild.2"))
					Expect(*updateAddonInput.ServiceAccountRoleArn).To(Equal("new-arn"))
				})
			})

			When("attachPolicyARNs is configured", func() {
				When("its an update to an existing cloudformation", func() {
					It("updates the stack", func() {
						fakeStackManager.DescribeStackReturns(&manager.Stack{
							StackName: aws.String("eksctl-my-cluster-addon-vpc-cni"),
							Outputs: []types.Output{
								{
									OutputValue: aws.String("new-service-account-role-arn"),
									OutputKey:   aws.String(outputs.IAMServiceAccountRoleName),
								},
							},
						}, nil)

						err := addonManager.Update(context.Background(), &api.Addon{
							Name:             "vpc-cni",
							Version:          "v1.0.0-eksbuild.2",
							AttachPolicyARNs: []string{"arn-1"},
						}, &podIdentityIAMUpdater, 0)

						Expect(err).NotTo(HaveOccurred())

						Expect(fakeStackManager.UpdateStackCallCount()).To(Equal(1))
						_, options := fakeStackManager.UpdateStackArgsForCall(0)
						Expect(*options.Stack.StackName).To(Equal("eksctl-my-cluster-addon-vpc-cni"))
						Expect(options.ChangeSetName).To(ContainSubstring("updating-policy"))
						Expect(options.Description).To(Equal("updating policies"))
						Expect(options.Wait).To(BeTrue())
						Expect(err).NotTo(HaveOccurred())
						Expect(string(options.TemplateData.(manager.TemplateBody))).To(ContainSubstring("arn-1"))
						Expect(string(options.TemplateData.(manager.TemplateBody))).To(ContainSubstring(":sub\":\"system:serviceaccount:kube-system:aws-node"))

						Expect(*updateAddonInput.ClusterName).To(Equal("my-cluster"))
						Expect(*updateAddonInput.AddonName).To(Equal("vpc-cni"))
						Expect(*updateAddonInput.AddonVersion).To(Equal("v1.0.0-eksbuild.2"))
						Expect(*updateAddonInput.ServiceAccountRoleArn).To(Equal("new-service-account-role-arn"))
					})
				})

				When("its a new set of ARNs", func() {
					It("creates a role with the ARNs", func() {
						err := addonManager.Update(context.Background(), &api.Addon{
							Name:             "my-addon",
							Version:          "v1.0.0-eksbuild.2",
							AttachPolicyARNs: []string{"arn-1"},
						}, &podIdentityIAMUpdater, 0)

						Expect(err).NotTo(HaveOccurred())

						Expect(fakeStackManager.CreateStackCallCount()).To(Equal(1))
						_, name, resourceSet, tags, _, _ := fakeStackManager.CreateStackArgsForCall(0)
						Expect(name).To(Equal("eksctl-my-cluster-addon-my-addon"))
						Expect(resourceSet).NotTo(BeNil())
						Expect(tags).To(Equal(map[string]string{
							api.AddonNameTag: "my-addon",
						}))
						output, err := resourceSet.RenderJSON()
						Expect(err).NotTo(HaveOccurred())
						Expect(string(output)).To(ContainSubstring("arn-1"))

						Expect(*updateAddonInput.ClusterName).To(Equal("my-cluster"))
						Expect(*updateAddonInput.AddonName).To(Equal("my-addon"))
						Expect(*updateAddonInput.AddonVersion).To(Equal("v1.0.0-eksbuild.2"))
						Expect(*updateAddonInput.ServiceAccountRoleArn).To(Equal("new-service-account-role-arn"))
					})
				})
			})

			When("attachPolicy is configured", func() {
				When("its an update to an existing cloudformation", func() {
					It("updates the stack", func() {
						fakeStackManager.DescribeStackReturns(&manager.Stack{
							StackName: aws.String("eksctl-my-cluster-addon-vpc-cni"),
							Outputs: []types.Output{
								{
									OutputValue: aws.String("new-service-account-role-arn"),
									OutputKey:   aws.String(outputs.IAMServiceAccountRoleName),
								},
							},
						}, nil)

						err := addonManager.Update(context.Background(), &api.Addon{
							Name:    "vpc-cni",
							Version: "v1.0.0-eksbuild.2",
							AttachPolicy: api.InlineDocument{
								"foo": "policy-bar",
							},
						}, &podIdentityIAMUpdater, 0)

						Expect(err).NotTo(HaveOccurred())

						Expect(fakeStackManager.UpdateStackCallCount()).To(Equal(1))
						_, options := fakeStackManager.UpdateStackArgsForCall(0)
						Expect(*options.Stack.StackName).To(Equal("eksctl-my-cluster-addon-vpc-cni"))
						Expect(options.ChangeSetName).To(ContainSubstring("updating-policy"))
						Expect(options.Description).To(Equal("updating policies"))
						Expect(options.Wait).To(BeTrue())
						Expect(err).NotTo(HaveOccurred())
						Expect(string(options.TemplateData.(manager.TemplateBody))).To(ContainSubstring("policy-bar"))

						Expect(*updateAddonInput.ClusterName).To(Equal("my-cluster"))
						Expect(*updateAddonInput.AddonName).To(Equal("vpc-cni"))
						Expect(*updateAddonInput.AddonVersion).To(Equal("v1.0.0-eksbuild.2"))
						Expect(*updateAddonInput.ServiceAccountRoleArn).To(Equal("new-service-account-role-arn"))
					})
				})

				When("its a new set of policies", func() {
					It("creates a role with the policies", func() {
						err := addonManager.Update(context.Background(), &api.Addon{
							Name:    "my-addon",
							Version: "v1.0.0-eksbuild.2",
							AttachPolicy: api.InlineDocument{
								"foo": "policy-bar",
							},
						}, &podIdentityIAMUpdater, 0)

						Expect(err).NotTo(HaveOccurred())

						Expect(fakeStackManager.CreateStackCallCount()).To(Equal(1))
						_, name, resourceSet, tags, _, _ := fakeStackManager.CreateStackArgsForCall(0)
						Expect(name).To(Equal("eksctl-my-cluster-addon-my-addon"))
						Expect(resourceSet).NotTo(BeNil())
						Expect(tags).To(Equal(map[string]string{
							api.AddonNameTag: "my-addon",
						}))
						output, err := resourceSet.RenderJSON()
						Expect(err).NotTo(HaveOccurred())
						Expect(string(output)).To(ContainSubstring("policy-bar"))

						Expect(*updateAddonInput.ClusterName).To(Equal("my-cluster"))
						Expect(*updateAddonInput.AddonName).To(Equal("my-addon"))
						Expect(*updateAddonInput.AddonVersion).To(Equal("v1.0.0-eksbuild.2"))
						Expect(*updateAddonInput.ServiceAccountRoleArn).To(Equal("new-service-account-role-arn"))
					})
				})
			})

			When("resolveConflicts is configured", func() {
				DescribeTable("AWS EKS resolve conflicts matches value from cluster config",
					func(rc ekstypes.ResolveConflicts) {
						err := addonManager.Update(context.Background(), &api.Addon{
							Name:             "my-addon",
							Version:          "v1.0.0-eksbuild.2",
							ResolveConflicts: rc,
						}, &podIdentityIAMUpdater, 0)

						Expect(err).NotTo(HaveOccurred())
						Expect(updateAddonInput.ResolveConflicts).To(Equal(rc))
					},
					Entry("none", ekstypes.ResolveConflictsNone),
					Entry("overwrite", ekstypes.ResolveConflictsOverwrite),
					Entry("preserve", ekstypes.ResolveConflictsPreserve),
				)
			})

			When("configurationValues is configured", func() {
				It("AWS EKS configuration values matches the value from cluster config", func() {
					err := addonManager.Update(context.Background(), &api.Addon{
						Name:                "my-addon",
						Version:             "v1.0.0-eksbuild.2",
						ConfigurationValues: "{\"replicaCount\":3}",
					}, &podIdentityIAMUpdater, 0)

					Expect(err).NotTo(HaveOccurred())
					Expect(aws.ToString(updateAddonInput.ConfigurationValues)).To(Equal("{\"replicaCount\":3}"))
				})
			})

			When("wellKnownPolicies are configured", func() {
				When("its an update to an existing cloudformation", func() {
					It("updates the stack", func() {
						fakeStackManager.DescribeStackReturns(&manager.Stack{
							StackName: aws.String("eksctl-my-cluster-addon-vpc-cni"),
							Outputs: []types.Output{
								{
									OutputValue: aws.String("new-service-account-role-arn"),
									OutputKey:   aws.String(outputs.IAMServiceAccountRoleName),
								},
							},
						}, nil)

						err := addonManager.Update(context.Background(), &api.Addon{
							Name:    "vpc-cni",
							Version: "v1.0.0-eksbuild.2",
							WellKnownPolicies: api.WellKnownPolicies{
								AutoScaler: true,
							},
						}, &podIdentityIAMUpdater, 0)

						Expect(err).NotTo(HaveOccurred())

						Expect(fakeStackManager.UpdateStackCallCount()).To(Equal(1))
						_, options := fakeStackManager.UpdateStackArgsForCall(0)
						Expect(*options.Stack.StackName).To(Equal("eksctl-my-cluster-addon-vpc-cni"))
						Expect(options.ChangeSetName).To(ContainSubstring("updating-policy"))
						Expect(options.Description).To(Equal("updating policies"))
						Expect(options.Wait).To(BeTrue())
						Expect(err).NotTo(HaveOccurred())
						Expect(string(options.TemplateData.(manager.TemplateBody))).To(ContainSubstring("autoscaling:SetDesiredCapacity"))

						Expect(*updateAddonInput.ClusterName).To(Equal("my-cluster"))
						Expect(*updateAddonInput.AddonName).To(Equal("vpc-cni"))
						Expect(*updateAddonInput.AddonVersion).To(Equal("v1.0.0-eksbuild.2"))
						Expect(*updateAddonInput.ServiceAccountRoleArn).To(Equal("new-service-account-role-arn"))
					})
				})

				When("its a new set of well known policies", func() {
					It("creates a role with the well known policies", func() {
						err := addonManager.Update(context.Background(), &api.Addon{
							Name:    "my-addon",
							Version: "v1.0.0-eksbuild.2",
							WellKnownPolicies: api.WellKnownPolicies{
								AutoScaler: true,
							},
						}, &podIdentityIAMUpdater, 0)

						Expect(err).NotTo(HaveOccurred())

						Expect(fakeStackManager.CreateStackCallCount()).To(Equal(1))
						_, name, resourceSet, tags, _, _ := fakeStackManager.CreateStackArgsForCall(0)
						Expect(name).To(Equal("eksctl-my-cluster-addon-my-addon"))
						Expect(resourceSet).NotTo(BeNil())
						Expect(tags).To(Equal(map[string]string{
							api.AddonNameTag: "my-addon",
						}))
						output, err := resourceSet.RenderJSON()
						Expect(err).NotTo(HaveOccurred())
						Expect(string(output)).To(ContainSubstring("autoscaling:SetDesiredCapacity"))

						Expect(*updateAddonInput.ClusterName).To(Equal("my-cluster"))
						Expect(*updateAddonInput.AddonName).To(Equal("my-addon"))
						Expect(*updateAddonInput.AddonVersion).To(Equal("v1.0.0-eksbuild.2"))
						Expect(*updateAddonInput.ServiceAccountRoleArn).To(Equal("new-service-account-role-arn"))
					})
				})
			})
		})
	})

	When("EKS fails to return an UpdateAddonOutput", func() {
		It("returns an error", func() {
			mockProvider.MockEKS().On("UpdateAddon", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(2))
				Expect(args[1]).To(BeAssignableToTypeOf(&awseks.UpdateAddonInput{}))
				updateAddonInput = args[1].(*awseks.UpdateAddonInput)
			}).Return(nil, fmt.Errorf("foo"))

			err := addonManager.Update(context.Background(), &api.Addon{
				Name: "my-addon",
			}, &mocks.PodIdentityIAMUpdater{}, 0)
			Expect(err).To(MatchError(`failed to update addon "my-addon": foo`))
			Expect(*updateAddonInput.ClusterName).To(Equal("my-cluster"))
			Expect(*updateAddonInput.AddonName).To(Equal("my-addon"))
		})
	})

	type addonPIAEntry struct {
		addonVersion                      string
		existingPodIdentityAssociations   []string
		addonsConfig                      api.AddonsConfig
		useDefaultPodIdentityAssociations bool
		mockDescribeAddonConfiguration    bool
		mockUpdateRole                    bool
		mockUpdateAddon                   bool

		expectedErr string
	}
	DescribeTable("updating pod identity associations", func(e addonPIAEntry) {
		const clusterName = "my-cluster"
		const addonVersion = "v1.7.5-eksbuild.1"
		mockProvider := mockprovider.NewMockProvider()

		mockProvider.MockEKS().On("DescribeAddon", mock.Anything, &awseks.DescribeAddonInput{
			ClusterName: aws.String(clusterName),
			AddonName:   aws.String(api.VPCCNIAddon),
		}).Return(&awseks.DescribeAddonOutput{
			Addon: &ekstypes.Addon{
				AddonName:               aws.String(api.VPCCNIAddon),
				ClusterName:             aws.String(clusterName),
				AddonVersion:            aws.String(addonVersion),
				PodIdentityAssociations: e.existingPodIdentityAssociations,
			},
		}, nil).Once()
		mockProvider.MockEKS().On("DescribeAddonVersions", mock.Anything, &awseks.DescribeAddonVersionsInput{
			AddonName:         aws.String("vpc-cni"),
			KubernetesVersion: aws.String(api.LatestVersion),
		}).Return(&awseks.DescribeAddonVersionsOutput{
			Addons: []ekstypes.AddonInfo{
				{
					AddonName: aws.String("vpc-cni"),
					AddonVersions: []ekstypes.AddonVersionInfo{
						{
							AddonVersion:           aws.String(addonVersion),
							RequiresIamPermissions: true,
						},
					},
				},
			},
		}, nil).Twice()

		if len(e.existingPodIdentityAssociations) > 0 {
			mockProvider.MockEKS().On("DescribePodIdentityAssociation", mock.Anything, &awseks.DescribePodIdentityAssociationInput{
				AssociationId: aws.String("a-zkgxwyqoexvjka9a3"),
				ClusterName:   aws.String(clusterName),
			}).Return(&awseks.DescribePodIdentityAssociationOutput{
				Association: &ekstypes.PodIdentityAssociation{
					AssociationId:  aws.String("a-zkgxwyqoexvjka9a3"),
					Namespace:      aws.String("kube-system"),
					ServiceAccount: aws.String("aws-node"),
					RoleArn:        aws.String("role-1"),
				},
			}, nil).Once()
		}

		if e.mockDescribeAddonConfiguration {
			mockProvider.MockEKS().On("DescribeAddonConfiguration", mock.Anything, &awseks.DescribeAddonConfigurationInput{
				AddonName:    aws.String(api.VPCCNIAddon),
				AddonVersion: aws.String(addonVersion),
			}).Return(&awseks.DescribeAddonConfigurationOutput{
				AddonName:    aws.String(api.VPCCNIAddon),
				AddonVersion: aws.String(addonVersion),
				PodIdentityConfiguration: []ekstypes.AddonPodIdentityConfiguration{
					{
						ServiceAccount:             aws.String("aws-node"),
						RecommendedManagedPolicies: []string{"arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"},
					},
				},
			}, nil).Once()
		}
		var podIdentityIAMUpdater mocks.PodIdentityIAMUpdater
		addonPIAs := []ekstypes.AddonPodIdentityAssociations{
			{
				ServiceAccount: aws.String("aws-node"),
				RoleArn:        aws.String("role-1"),
			},
		}
		if e.mockUpdateRole {
			podIdentityIAMUpdater.On("UpdateRole", mock.Anything, []api.PodIdentityAssociation{
				{
					Namespace:            "",
					ServiceAccountName:   "aws-node",
					PermissionPolicyARNs: []string{"arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"},
				},
			}, "vpc-cni", mock.Anything).Return(addonPIAs, nil).Once()
		}
		if e.mockUpdateAddon {
			updateAddonInput := &awseks.UpdateAddonInput{
				AddonName:    aws.String("vpc-cni"),
				ClusterName:  aws.String("my-cluster"),
				AddonVersion: aws.String("v1.7.5-eksbuild.1"),
			}
			if len(e.existingPodIdentityAssociations) == 0 {
				updateAddonInput.PodIdentityAssociations = addonPIAs
			}
			mockProvider.MockEKS().On("UpdateAddon", mock.Anything, updateAddonInput).Return(&awseks.UpdateAddonOutput{}, nil).Once()
		}

		addonManager, err := addon.New(&api.ClusterConfig{
			Metadata: &api.ClusterMeta{
				Version: api.LatestVersion,
				Name:    clusterName,
			},
			AddonsConfig: e.addonsConfig,
		}, mockProvider.EKS(), fakeStackManager, true, makeOIDCManager(), nil)
		Expect(err).NotTo(HaveOccurred())

		err = addonManager.Update(context.Background(), &api.Addon{
			Name:                              "vpc-cni",
			UseDefaultPodIdentityAssociations: e.useDefaultPodIdentityAssociations,
			Version:                           e.addonVersion,
		}, &podIdentityIAMUpdater, 0)
		if e.expectedErr != "" {
			Expect(err).To(MatchError(e.expectedErr))
		} else {
			Expect(err).NotTo(HaveOccurred())
		}
		mockProvider.MockEKS().AssertExpectations(GinkgoT())
		podIdentityIAMUpdater.AssertExpectations(GinkgoT())
	},
		Entry("addon with version", addonPIAEntry{
			addonVersion:                      "v1.7.5-eksbuild.1",
			useDefaultPodIdentityAssociations: true,
			mockDescribeAddonConfiguration:    true,
			mockUpdateRole:                    true,
			mockUpdateAddon:                   true,
		}),
		Entry("addon without version", addonPIAEntry{
			useDefaultPodIdentityAssociations: true,
			mockDescribeAddonConfiguration:    true,
			mockUpdateRole:                    true,
			mockUpdateAddon:                   true,
		}),
		Entry("addon with existing pod identity associations", addonPIAEntry{
			existingPodIdentityAssociations: []string{"arn:aws:eks:us-west-2:00:podidentityassociation/cluster/a-zkgxwyqoexvjka9a3"},
			expectedErr: "addon vpc-cni has pod identity associations, to remove pod identity associations from an addon, " +
				"addon.podIdentityAssociations must be explicitly set to []; if the addon was migrated to use pod identity, " +
				"addon.podIdentityAssociations must be set to values obtained from " +
				"`aws eks describe-pod-identity-association --cluster-name=my-cluster --association-id=a-zkgxwyqoexvjka9a3`",
		}),
		Entry("addon with existing pod identity associations and addonsConfig.autoApplyPodIdentityAssociations", addonPIAEntry{
			existingPodIdentityAssociations: []string{"arn:aws:eks:us-west-2:00:podidentityassociation/cluster/a-zkgxwyqoexvjka9a3"},
			mockUpdateAddon:                 true,
			addonsConfig: api.AddonsConfig{
				AutoApplyPodIdentityAssociations: true,
			},
		}),
	)
})
