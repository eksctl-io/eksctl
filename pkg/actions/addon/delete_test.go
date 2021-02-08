package addon_test

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"

	awseks "github.com/aws/aws-sdk-go/service/eks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/weaveworks/eksctl/pkg/actions/addon"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager/fakes"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("Delete", func() {
	var (
		manager          *addon.Manager
		withOIDC         bool
		fakeStackManager *fakes.FakeStackManager
		mockProvider     *mockprovider.MockProvider
	)

	BeforeEach(func() {
		withOIDC = false
		fakeStackManager = new(fakes.FakeStackManager)
		mockProvider = mockprovider.NewMockProvider()

		var err error
		manager, err = addon.New(&api.ClusterConfig{Metadata: &api.ClusterMeta{
			Version: "1.18",
			Name:    "my-cluster",
		}}, &eks.ClusterProvider{Provider: mockProvider}, fakeStackManager, withOIDC, nil, nil)
		Expect(err).NotTo(HaveOccurred())
	})

	It("deletes all associated stacks and addons", func() {
		var deleteAddonInput *awseks.DeleteAddonInput

		mockProvider.MockEKS().On("DeleteAddon", mock.Anything).Run(func(args mock.Arguments) {
			Expect(args).To(HaveLen(1))
			Expect(args[0]).To(BeAssignableToTypeOf(&awseks.DeleteAddonInput{}))
			deleteAddonInput = args[0].(*awseks.DeleteAddonInput)
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
		Expect(*deleteAddonInput.ClusterName).To(Equal("my-cluster"))
		Expect(*deleteAddonInput.AddonName).To(Equal("my-addon"))
	})

	When("delete addon fails", func() {
		It("returns an error", func() {
			var deleteAddonInput *awseks.DeleteAddonInput

			mockProvider.MockEKS().On("DeleteAddon", mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(1))
				Expect(args[0]).To(BeAssignableToTypeOf(&awseks.DeleteAddonInput{}))
				deleteAddonInput = args[0].(*awseks.DeleteAddonInput)
			}).Return(&awseks.DeleteAddonOutput{}, fmt.Errorf("foo"))

			err := manager.Delete(&api.Addon{
				Name: "my-addon",
			})

			Expect(err).To(MatchError(`failed to delete addon "my-addon": foo`))
			Expect(*deleteAddonInput.ClusterName).To(Equal("my-cluster"))
			Expect(*deleteAddonInput.AddonName).To(Equal("my-addon"))
		})
	})

	When("list stacks fails", func() {
		It("only deletes the addon", func() {
			var deleteAddonInput *awseks.DeleteAddonInput

			mockProvider.MockEKS().On("DeleteAddon", mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(1))
				Expect(args[0]).To(BeAssignableToTypeOf(&awseks.DeleteAddonInput{}))
				deleteAddonInput = args[0].(*awseks.DeleteAddonInput)
			}).Return(&awseks.DeleteAddonOutput{}, nil)

			fakeStackManager.ListStacksMatchingReturns([]*cloudformation.Stack{}, fmt.Errorf("foo"))

			err := manager.Delete(&api.Addon{
				Name: "my-addon",
			})

			Expect(err).To(MatchError("failed to list stacks: foo"))
			Expect(*deleteAddonInput.ClusterName).To(Equal("my-cluster"))
			Expect(*deleteAddonInput.AddonName).To(Equal("my-addon"))
		})
	})

	When("delete stack fails", func() {
		It("returns an error", func() {
			var deleteAddonInput *awseks.DeleteAddonInput

			mockProvider.MockEKS().On("DeleteAddon", mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(1))
				Expect(args[0]).To(BeAssignableToTypeOf(&awseks.DeleteAddonInput{}))
				deleteAddonInput = args[0].(*awseks.DeleteAddonInput)
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
			Expect(*deleteAddonInput.ClusterName).To(Equal("my-cluster"))
			Expect(*deleteAddonInput.AddonName).To(Equal("my-addon"))
			Expect(fakeStackManager.DeleteStackByNameCallCount()).To(Equal(1))
			Expect(fakeStackManager.DeleteStackByNameArgsForCall(0)).To(Equal("eksctl-my-cluster-addon-my-addon"))
		})
	})

	When("no stack exists", func() {
		It("only deletes the addon", func() {
			var deleteAddonInput *awseks.DeleteAddonInput

			mockProvider.MockEKS().On("DeleteAddon", mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(1))
				Expect(args[0]).To(BeAssignableToTypeOf(&awseks.DeleteAddonInput{}))
				deleteAddonInput = args[0].(*awseks.DeleteAddonInput)
			}).Return(&awseks.DeleteAddonOutput{}, nil)

			fakeStackManager.ListStacksMatchingReturns([]*cloudformation.Stack{}, nil)

			err := manager.Delete(&api.Addon{
				Name: "my-addon",
			})

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeStackManager.DeleteStackByNameCallCount()).To(Equal(0))
			Expect(*deleteAddonInput.ClusterName).To(Equal("my-cluster"))
			Expect(*deleteAddonInput.AddonName).To(Equal("my-addon"))
		})
	})
})
