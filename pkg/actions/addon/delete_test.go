package addon_test

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/smithy-go"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/pkg/errors"

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

			fakeStackManager.DescribeStackReturns(&types.Stack{StackName: aws.String("eksctl-my-cluster-addon-my-addon")}, nil)

			err := manager.Delete(context.TODO(), &api.Addon{
				Name: "my-addon",
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeStackManager.DeleteStackBySpecCallCount()).To(Equal(1))
			_, stack := fakeStackManager.DeleteStackBySpecArgsForCall(0)
			Expect(*stack.StackName).To(Equal("eksctl-my-cluster-addon-my-addon"))
		})

		When("delete addon fails", func() {
			It("returns an error", func() {
				mockProvider.MockEKS().On("DeleteAddon", &awseks.DeleteAddonInput{
					AddonName:   aws.String("my-addon"),
					ClusterName: aws.String("my-cluster"),
				}).Return(&awseks.DeleteAddonOutput{}, fmt.Errorf("foo"))

				err := manager.Delete(context.TODO(), &api.Addon{
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

				fakeStackManager.DescribeStackReturns(nil, fmt.Errorf("foo"))

				err := manager.Delete(context.TODO(), &api.Addon{
					Name: "my-addon",
				})

				Expect(err).To(MatchError("failed to get stack: foo"))
			})
		})

		When("delete stack fails", func() {
			It("returns an error", func() {
				mockProvider.MockEKS().On("DeleteAddon", &awseks.DeleteAddonInput{
					AddonName:   aws.String("my-addon"),
					ClusterName: aws.String("my-cluster"),
				}).Return(&awseks.DeleteAddonOutput{}, nil)

				fakeStackManager.DeleteStackBySpecReturns(nil, fmt.Errorf("foo"))
				fakeStackManager.DescribeStackReturns(&types.Stack{
					StackName: aws.String("eksctl-my-cluster-addon-my-addon"),
				}, nil)

				err := manager.Delete(context.TODO(), &api.Addon{
					Name: "my-addon",
				})

				Expect(err).To(MatchError(`failed to delete cloudformation stack "eksctl-my-cluster-addon-my-addon": foo`))
				Expect(fakeStackManager.DeleteStackBySpecCallCount()).To(Equal(1))
				_, stack := fakeStackManager.DeleteStackBySpecArgsForCall(0)
				Expect(*stack.StackName).To(Equal("eksctl-my-cluster-addon-my-addon"))
			})
		})

		When("no stack exists", func() {
			It("only deletes the addon", func() {
				mockProvider.MockEKS().On("DeleteAddon", &awseks.DeleteAddonInput{
					AddonName:   aws.String("my-addon"),
					ClusterName: aws.String("my-cluster"),
				}).Return(&awseks.DeleteAddonOutput{}, nil)

				fakeStackManager.DescribeStackReturns(nil, errors.Wrap(&smithy.OperationError{
					Err: fmt.Errorf("ValidationError"),
				}, "nope"))

				err := manager.Delete(context.TODO(), &api.Addon{
					Name: "my-addon",
				})

				Expect(err).NotTo(HaveOccurred())

				Expect(fakeStackManager.DeleteStackBySpecCallCount()).To(Equal(0))
			})
		})

		When("when no addon exists, but the stack does", func() {
			It("only deletes the stack", func() {
				mockProvider.MockEKS().On("DeleteAddon", &awseks.DeleteAddonInput{
					AddonName:   aws.String("my-addon"),
					ClusterName: aws.String("my-cluster"),
				}).Return(&awseks.DeleteAddonOutput{}, awserr.New(awseks.ErrCodeResourceNotFoundException, "", nil))

				fakeStackManager.DescribeStackReturns(&types.Stack{StackName: aws.String("eksctl-my-cluster-addon-my-addon")}, nil)

				err := manager.Delete(context.TODO(), &api.Addon{
					Name: "my-addon",
				})

				Expect(err).NotTo(HaveOccurred())

				Expect(fakeStackManager.DeleteStackBySpecCallCount()).To(Equal(1))
				_, stack := fakeStackManager.DeleteStackBySpecArgsForCall(0)
				Expect(*stack.StackName).To(Equal("eksctl-my-cluster-addon-my-addon"))
			})
		})

		When("when no addon exists or stack exists", func() {
			It("errors", func() {
				mockProvider.MockEKS().On("DeleteAddon", &awseks.DeleteAddonInput{
					AddonName:   aws.String("my-addon"),
					ClusterName: aws.String("my-cluster"),
				}).Return(&awseks.DeleteAddonOutput{}, awserr.New(awseks.ErrCodeResourceNotFoundException, "", nil))

				fakeStackManager.DescribeStackReturns(nil, errors.Wrap(&smithy.OperationError{
					Err: fmt.Errorf("ValidationError"),
				}, "nope"))
				err := manager.Delete(context.TODO(), &api.Addon{
					Name: "my-addon",
				})

				Expect(err).To(MatchError("could not find addon or associated IAM stack to delete"))
				Expect(fakeStackManager.DeleteStackBySpecCallCount()).To(Equal(0))
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
