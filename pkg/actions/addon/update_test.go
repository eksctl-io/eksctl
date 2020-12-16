package addon_test

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/cloudformation"

	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"

	"github.com/weaveworks/eksctl/pkg/cfn/manager"

	"github.com/weaveworks/eksctl/pkg/cfn/builder"

	"github.com/weaveworks/eksctl/pkg/actions/addon/fakes"

	"github.com/aws/aws-sdk-go/aws"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/weaveworks/eksctl/pkg/actions/addon"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("Update", func() {
	var (
		addonManager       *addon.Manager
		mockProvider       *mockprovider.MockProvider
		updateAddonInput   *awseks.UpdateAddonInput
		describeAddonInput *awseks.DescribeAddonInput
		fakeStackManager   *fakes.FakeStackManager
	)
	BeforeEach(func() {
		var err error
		mockProvider = mockprovider.NewMockProvider()
		fakeStackManager = new(fakes.FakeStackManager)

		fakeStackManager.CreateStackStub = func(_ string, rs builder.ResourceSet, _ map[string]string, _ map[string]string, errs chan error) error {
			go func() {
				errs <- nil
			}()
			Expect(rs).To(BeAssignableToTypeOf(&builder.IAMRoleResourceSet{}))
			rs.(*builder.IAMRoleResourceSet).OutputRole = "new-service-account-role-arn"
			return nil
		}

		oidc, err := iamoidc.NewOpenIDConnectManager(nil, "456123987123", "https://oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E", "aws")
		Expect(err).ToNot(HaveOccurred())
		oidc.ProviderARN = "arn:aws:iam::456123987123:oidc-provider/oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E"

		mockProvider.MockEKS().On("DescribeAddon", mock.Anything).Run(func(args mock.Arguments) {
			Expect(args).To(HaveLen(1))
			Expect(args[0]).To(BeAssignableToTypeOf(&awseks.DescribeAddonInput{}))
			describeAddonInput = args[0].(*awseks.DescribeAddonInput)
		}).Return(&awseks.DescribeAddonOutput{
			Addon: &awseks.Addon{
				AddonName:             aws.String("my-addon"),
				AddonVersion:          aws.String("1.0"),
				ServiceAccountRoleArn: aws.String("original-arn"),
				Status:                aws.String("created"),
			},
		}, nil)

		addonManager, err = addon.New(&api.ClusterConfig{Metadata: &api.ClusterMeta{
			Version: "1.18",
			Name:    "my-cluster",
		}}, &eks.ClusterProvider{Provider: mockProvider}, fakeStackManager, true, oidc, nil)
		Expect(err).NotTo(HaveOccurred())

	})

	When("Updating the version", func() {
		It("updates the addon and preserves the existing role", func() {
			mockProvider.MockEKS().On("UpdateAddon", mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(1))
				Expect(args[0]).To(BeAssignableToTypeOf(&awseks.UpdateAddonInput{}))
				updateAddonInput = args[0].(*awseks.UpdateAddonInput)
			}).Return(&awseks.UpdateAddonOutput{}, nil)

			err := addonManager.Update(&api.Addon{
				Name:    "my-addon",
				Version: "1.1",
				Force:   true,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(*describeAddonInput.ClusterName).To(Equal("my-cluster"))
			Expect(*describeAddonInput.AddonName).To(Equal("my-addon"))
			Expect(*updateAddonInput.ClusterName).To(Equal("my-cluster"))
			Expect(*updateAddonInput.AddonName).To(Equal("my-addon"))
			Expect(*updateAddonInput.AddonVersion).To(Equal("1.1"))
			Expect(*updateAddonInput.ServiceAccountRoleArn).To(Equal("original-arn"))
			Expect(*updateAddonInput.ResolveConflicts).To(Equal("overwrite"))
		})
	})

	When("updating the policy", func() {
		When("specifying a new serviceAccountRoleARN", func() {
			It("updates the addon", func() {
				mockProvider.MockEKS().On("UpdateAddon", mock.Anything).Run(func(args mock.Arguments) {
					Expect(args).To(HaveLen(1))
					Expect(args[0]).To(BeAssignableToTypeOf(&awseks.UpdateAddonInput{}))
					updateAddonInput = args[0].(*awseks.UpdateAddonInput)
				}).Return(&awseks.UpdateAddonOutput{}, nil)

				err := addonManager.Update(&api.Addon{
					Name:                  "my-addon",
					Version:               "1.1",
					ServiceAccountRoleARN: "new-arn",
				})

				Expect(err).NotTo(HaveOccurred())
				Expect(*describeAddonInput.ClusterName).To(Equal("my-cluster"))
				Expect(*describeAddonInput.AddonName).To(Equal("my-addon"))
				Expect(*updateAddonInput.ClusterName).To(Equal("my-cluster"))
				Expect(*updateAddonInput.AddonName).To(Equal("my-addon"))
				Expect(*updateAddonInput.AddonVersion).To(Equal("1.1"))
				Expect(*updateAddonInput.ServiceAccountRoleArn).To(Equal("new-arn"))
			})
		})

		When("attachPolicyARNs is configured", func() {
			When("its an update to an existing cloudformation", func() {
				It("uses the updates stacks role", func() {
					fakeStackManager.ListStacksMatchingReturns([]*manager.Stack{
						{
							StackName: aws.String("eksctl-my-cluster-addon-my-addon"),
							Outputs: []*cloudformation.Output{
								{
									OutputValue: aws.String("new-service-account-role-arn"),
									OutputKey:   aws.String("Role1"),
								},
							},
						},
					}, nil)

					mockProvider.MockEKS().On("UpdateAddon", mock.Anything).Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(1))
						Expect(args[0]).To(BeAssignableToTypeOf(&awseks.UpdateAddonInput{}))
						updateAddonInput = args[0].(*awseks.UpdateAddonInput)
					}).Return(&awseks.UpdateAddonOutput{}, nil)

					err := addonManager.Update(&api.Addon{
						Name:             "vpc-cni",
						Version:          "1.1",
						AttachPolicyARNs: []string{"arn-1"},
					})

					Expect(err).NotTo(HaveOccurred())

					Expect(fakeStackManager.UpdateStackCallCount()).To(Equal(1))
					stackName, changeSetName, description, templateData, _ := fakeStackManager.UpdateStackArgsForCall(0)
					Expect(stackName).To(Equal("eksctl-my-cluster-addon-vpc-cni"))
					Expect(changeSetName).To(Equal("updating-policy"))
					Expect(description).To(Equal("updating policies"))
					Expect(err).NotTo(HaveOccurred())
					Expect(string(templateData.(manager.TemplateBody))).To(ContainSubstring("arn-1"))
					Expect(string(templateData.(manager.TemplateBody))).To(ContainSubstring(":sub\":\"system:serviceaccount:kube-system:aws-node"))

					Expect(*updateAddonInput.ClusterName).To(Equal("my-cluster"))
					Expect(*updateAddonInput.AddonName).To(Equal("vpc-cni"))
					Expect(*updateAddonInput.AddonVersion).To(Equal("1.1"))
					Expect(*updateAddonInput.ServiceAccountRoleArn).To(Equal("new-service-account-role-arn"))

				})
			})

			When("its a new set of arns", func() {
				It("uses AttachPolicyARNS to create a role to attach to the addon", func() {
					mockProvider.MockEKS().On("UpdateAddon", mock.Anything).Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(1))
						Expect(args[0]).To(BeAssignableToTypeOf(&awseks.UpdateAddonInput{}))
						updateAddonInput = args[0].(*awseks.UpdateAddonInput)
					}).Return(&awseks.UpdateAddonOutput{}, nil)

					err := addonManager.Update(&api.Addon{
						Name:             "my-addon",
						Version:          "1.1",
						AttachPolicyARNs: []string{"arn-1"},
					})

					Expect(err).NotTo(HaveOccurred())

					Expect(fakeStackManager.CreateStackCallCount()).To(Equal(1))
					name, resourceSet, tags, _, _ := fakeStackManager.CreateStackArgsForCall(0)
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
					Expect(*updateAddonInput.AddonVersion).To(Equal("1.1"))
					Expect(*updateAddonInput.ServiceAccountRoleArn).To(Equal("new-service-account-role-arn"))
				})
			})
		})
	})

	When("it fails to update the addon", func() {
		It("returns an error", func() {
			mockProvider.MockEKS().On("UpdateAddon", mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(1))
				Expect(args[0]).To(BeAssignableToTypeOf(&awseks.UpdateAddonInput{}))
				updateAddonInput = args[0].(*awseks.UpdateAddonInput)
			}).Return(nil, fmt.Errorf("foo"))

			err := addonManager.Update(&api.Addon{
				Name: "my-addon",
			})
			Expect(err).To(MatchError(`failed to update addon "my-addon": foo`))
			Expect(*updateAddonInput.ClusterName).To(Equal("my-cluster"))
			Expect(*updateAddonInput.AddonName).To(Equal("my-addon"))
		})
	})
})
