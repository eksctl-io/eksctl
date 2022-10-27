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
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/cfn/manager/fakes"
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

		oidc, err := iamoidc.NewOpenIDConnectManager(nil, "456123987123", "https://oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E", "aws", nil)
		Expect(err).NotTo(HaveOccurred())
		oidc.ProviderARN = "arn:aws:iam::456123987123:oidc-provider/oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E"

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
				}, 0)

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
					}, 0)

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
					}, 0)

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
					}, 0)

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
					}, 0)
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
					}, waitTimeout)
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
					}, waitTimeout)
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
					}, 0)

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
						}, 0)

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
						}, 0)

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
						}, 0)

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
						}, 0)

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
						}, 0)

						Expect(err).NotTo(HaveOccurred())
						Expect(updateAddonInput.ResolveConflicts).To(Equal(rc))
					},
					Entry("none", ekstypes.ResolveConflictsNone),
					Entry("overwrite", ekstypes.ResolveConflictsOverwrite),
					Entry("preserve", ekstypes.ResolveConflictsPreserve),
				)
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
						}, 0)

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
						}, 0)

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
			}, 0)
			Expect(err).To(MatchError(`failed to update addon "my-addon": foo`))
			Expect(*updateAddonInput.ClusterName).To(Equal("my-cluster"))
			Expect(*updateAddonInput.AddonName).To(Equal("my-addon"))
		})
	})
})
