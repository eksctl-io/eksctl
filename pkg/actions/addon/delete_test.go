package addon_test

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"

	awseks "github.com/aws/aws-sdk-go/service/eks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/actions/addon"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager/fakes"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("Delete", func() {
	var (
		manager          *addon.Manager
		withOIDC         bool
		fakeStackManager *fakes.FakeStackManager
		mockProvider     *mockprovider.MockProvider
	)

	Describe("Delete", func() {
		BeforeEach(func() {
			withOIDC = false
			fakeStackManager = new(fakes.FakeStackManager)
			mockProvider = mockprovider.NewMockProvider()

			var err error
			manager, err = addon.New(&api.ClusterConfig{Metadata: &api.ClusterMeta{
				Version: "1.18",
				Name:    "my-cluster",
			}}, mockProvider.EKS(), fakeStackManager, withOIDC, nil, nil, 5*time.Minute)
			Expect(err).NotTo(HaveOccurred())
		})

		It("deletes all associated stacks and addons", func() {
			mockProvider.MockEKS().On("DeleteAddon", &awseks.DeleteAddonInput{
				AddonName:   aws.String("my-addon"),
				ClusterName: aws.String("my-cluster"),
			}).Return(&awseks.DeleteAddonOutput{}, nil)

			fakeStackManager.ListStacksMatchingReturns([]*cloudformation.Stack{
				{
					StackName: aws.String("eksctl-my-cluster-addon-my-addon"),
				},
			}, nil)

			err := manager.Delete(&api.Addon{
				Name: "my-addon",
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeStackManager.DeleteStackByNameCallCount()).To(Equal(1))
			Expect(fakeStackManager.DeleteStackByNameArgsForCall(0)).To(Equal("eksctl-my-cluster-addon-my-addon"))
		})

		When("delete addon fails", func() {
			It("returns an error", func() {
				mockProvider.MockEKS().On("DeleteAddon", &awseks.DeleteAddonInput{
					AddonName:   aws.String("my-addon"),
					ClusterName: aws.String("my-cluster"),
				}).Return(&awseks.DeleteAddonOutput{}, fmt.Errorf("foo"))

				err := manager.Delete(&api.Addon{
					Name: "my-addon",
				})

				Expect(err).To(MatchError(`failed to delete addon "my-addon": foo`))
			})
		})

		When("list stacks fails", func() {
			It("only deletes the addon", func() {
				mockProvider.MockEKS().On("DeleteAddon", &awseks.DeleteAddonInput{
					AddonName:   aws.String("my-addon"),
					ClusterName: aws.String("my-cluster"),
				}).Return(&awseks.DeleteAddonOutput{}, nil)

				fakeStackManager.ListStacksMatchingReturns([]*cloudformation.Stack{}, fmt.Errorf("foo"))

				err := manager.Delete(&api.Addon{
					Name: "my-addon",
				})

				Expect(err).To(MatchError("failed to list stacks: foo"))
			})
		})

		When("delete stack fails", func() {
			It("returns an error", func() {
				mockProvider.MockEKS().On("DeleteAddon", &awseks.DeleteAddonInput{
					AddonName:   aws.String("my-addon"),
					ClusterName: aws.String("my-cluster"),
				}).Return(&awseks.DeleteAddonOutput{}, nil)

				fakeStackManager.DeleteStackByNameReturns(nil, fmt.Errorf("foo"))
				fakeStackManager.ListStacksMatchingReturns([]*cloudformation.Stack{
					{
						StackName: aws.String("eksctl-my-cluster-addon-my-addon"),
					},
				}, nil)

				err := manager.Delete(&api.Addon{
					Name: "my-addon",
				})

				Expect(err).To(MatchError(`failed to delete cloudformation stack "eksctl-my-cluster-addon-my-addon": foo`))
				Expect(fakeStackManager.DeleteStackByNameCallCount()).To(Equal(1))
				Expect(fakeStackManager.DeleteStackByNameArgsForCall(0)).To(Equal("eksctl-my-cluster-addon-my-addon"))
			})
		})

		When("no stack exists", func() {
			It("only deletes the addon", func() {
				mockProvider.MockEKS().On("DeleteAddon", &awseks.DeleteAddonInput{
					AddonName:   aws.String("my-addon"),
					ClusterName: aws.String("my-cluster"),
				}).Return(&awseks.DeleteAddonOutput{}, nil)

				fakeStackManager.ListStacksMatchingReturns([]*cloudformation.Stack{}, nil)

				err := manager.Delete(&api.Addon{
					Name: "my-addon",
				})

				Expect(err).NotTo(HaveOccurred())

				Expect(fakeStackManager.DeleteStackByNameCallCount()).To(Equal(0))
			})
		})

		When("when no addon exists, but the stack does", func() {
			It("only deletes the stack", func() {
				mockProvider.MockEKS().On("DeleteAddon", &awseks.DeleteAddonInput{
					AddonName:   aws.String("my-addon"),
					ClusterName: aws.String("my-cluster"),
				}).Return(&awseks.DeleteAddonOutput{}, awserr.New(awseks.ErrCodeResourceNotFoundException, "", nil))

				fakeStackManager.ListStacksMatchingReturns([]*cloudformation.Stack{
					{
						StackName: aws.String("eksctl-my-cluster-addon-my-addon"),
					},
				}, nil)
				err := manager.Delete(&api.Addon{
					Name: "my-addon",
				})

				Expect(err).NotTo(HaveOccurred())

				Expect(fakeStackManager.DeleteStackByNameCallCount()).To(Equal(1))
				Expect(fakeStackManager.DeleteStackByNameArgsForCall(0)).To(Equal("eksctl-my-cluster-addon-my-addon"))
			})
		})

		When("when no addon exists or stack exists", func() {
			It("errors", func() {
				mockProvider.MockEKS().On("DeleteAddon", &awseks.DeleteAddonInput{
					AddonName:   aws.String("my-addon"),
					ClusterName: aws.String("my-cluster"),
				}).Return(&awseks.DeleteAddonOutput{}, awserr.New(awseks.ErrCodeResourceNotFoundException, "", nil))

				fakeStackManager.ListStacksMatchingReturns([]*cloudformation.Stack{}, nil)
				err := manager.Delete(&api.Addon{
					Name: "my-addon",
				})

				Expect(err).To(MatchError("could not find addon or associated IAM stack to delete"))
				Expect(fakeStackManager.DeleteStackByNameCallCount()).To(Equal(0))
			})
		})
	})

	Describe("DeleteWithPreserve", func() {
		BeforeEach(func() {
			withOIDC = false
			fakeStackManager = new(fakes.FakeStackManager)
			mockProvider = mockprovider.NewMockProvider()

			var err error
			manager, err = addon.New(&api.ClusterConfig{Metadata: &api.ClusterMeta{
				Version: "1.18",
				Name:    "my-cluster",
			}}, mockProvider.EKS(), fakeStackManager, withOIDC, nil, nil, 5*time.Minute)
			Expect(err).NotTo(HaveOccurred())
		})

		It("deletes the addon but preserves the resources", func() {
			mockProvider.MockEKS().On("DeleteAddon", &awseks.DeleteAddonInput{
				AddonName:   aws.String("my-addon"),
				ClusterName: aws.String("my-cluster"),
				Preserve:    aws.Bool(true),
			}).Return(&awseks.DeleteAddonOutput{}, nil)

			err := manager.DeleteWithPreserve(&api.Addon{
				Name: "my-addon",
			})
			Expect(err).NotTo(HaveOccurred())
		})

		When("delete addon fails", func() {
			It("returns an error", func() {
				mockProvider.MockEKS().On("DeleteAddon", &awseks.DeleteAddonInput{
					AddonName:   aws.String("my-addon"),
					ClusterName: aws.String("my-cluster"),
					Preserve:    aws.Bool(true),
				}).Return(&awseks.DeleteAddonOutput{}, fmt.Errorf("foo"))

				err := manager.DeleteWithPreserve(&api.Addon{
					Name: "my-addon",
				})

				Expect(err).To(MatchError(`failed to delete addon "my-addon": foo`))
			})
		})
	})
})
