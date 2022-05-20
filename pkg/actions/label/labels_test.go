package label_test

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/aws/smithy-go"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	perrors "github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/weaveworks/eksctl/pkg/actions/label"
	"github.com/weaveworks/eksctl/pkg/actions/label/fakes"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("Labels", func() {
	var (
		fakeManagedService *fakes.FakeService
		mockProvider       *mockprovider.MockProvider
		manager            *label.Manager

		clusterName   string
		nodegroupName string
	)

	BeforeEach(func() {
		fakeManagedService = new(fakes.FakeService)
		mockProvider = mockprovider.NewMockProvider()
		clusterName = "foo"
		nodegroupName = "bar"
		manager = label.New(clusterName, fakeManagedService, mockProvider.EKS())
		manager.SetService(fakeManagedService)
	})

	Describe("Get", func() {
		var expectedLabels map[string]string

		BeforeEach(func() {
			expectedLabels = map[string]string{"k1": "v1"}
		})

		When("the nodegroup is owned by eksctl", func() {
			BeforeEach(func() {
				fakeManagedService.GetLabelsReturns(expectedLabels, nil)
			})

			It("returns the labels from the nodegroup stack", func() {
				summary, err := manager.Get(context.Background(), nodegroupName)
				Expect(err).NotTo(HaveOccurred())
				Expect(summary[0].Labels).To(Equal(expectedLabels))
			})

			When("the service returns an error", func() {
				BeforeEach(func() {
					fakeManagedService.GetLabelsReturns(nil, errors.New("something-terrible"))
				})

				It("fails", func() {
					summary, err := manager.Get(context.Background(), nodegroupName)
					Expect(err).To(HaveOccurred())
					Expect(summary).To(BeNil())
				})
			})
		})

		When("the nodegroup is not owned by eksctl", func() {
			var returnedLabels map[string]string

			BeforeEach(func() {
				returnedLabels = map[string]string{"k1": "v1"}
				err := &smithy.OperationError{
					Err: errors.New("ValidationError"),
				}
				fakeManagedService.GetLabelsReturns(nil, perrors.Wrapf(err, "omg %s", "what"))
			})

			It("returns the labels from the EKS api", func() {
				mockProvider.MockEKS().On("DescribeNodegroup", mock.Anything, &eks.DescribeNodegroupInput{
					ClusterName:   aws.String(clusterName),
					NodegroupName: aws.String(nodegroupName),
				}).Return(&eks.DescribeNodegroupOutput{Nodegroup: &ekstypes.Nodegroup{Labels: returnedLabels}}, nil)

				summary, err := manager.Get(context.Background(), nodegroupName)
				Expect(err).NotTo(HaveOccurred())
				Expect(summary[0].Labels).To(Equal(expectedLabels))
			})

			When("the EKS api returns an error", func() {
				It("fails", func() {
					mockProvider.MockEKS().On("DescribeNodegroup", mock.Anything, mock.Anything).Return(&eks.DescribeNodegroupOutput{}, errors.New("oh-noes"))

					summary, err := manager.Get(context.Background(), nodegroupName)
					Expect(err).To(HaveOccurred())
					Expect(summary).To(BeNil())
				})
			})
		})
	})

	Describe("Set", func() {
		var labels map[string]string

		BeforeEach(func() {
			labels = map[string]string{"k1": "v1"}
		})

		When("the nodegroup is owned by eksctl", func() {
			BeforeEach(func() {
				fakeManagedService.UpdateLabelsReturns(nil)
			})

			It("sets new labels by updating the nodegroup stack", func() {
				Expect(manager.Set(context.Background(), nodegroupName, labels)).To(Succeed())
			})

			When("the service returns an error", func() {
				BeforeEach(func() {
					fakeManagedService.UpdateLabelsReturns(errors.New("something-terrible"))
				})

				It("fails", func() {
					err := manager.Set(context.Background(), nodegroupName, labels)
					Expect(err).To(HaveOccurred())
				})
			})
		})

		When("the nodegroup is not owned by eksctl", func() {
			var eksLabels map[string]string

			BeforeEach(func() {
				eksLabels = map[string]string{"k1": "v1"}
				err := &smithy.OperationError{
					Err: errors.New("ValidationError"),
				}
				fakeManagedService.UpdateLabelsReturns(perrors.Wrapf(err, "omg %s", "what"))
			})

			It("updates the labels through the EKS api", func() {
				mockProvider.MockEKS().On("UpdateNodegroupConfig", mock.Anything, &eks.UpdateNodegroupConfigInput{
					ClusterName:   aws.String(clusterName),
					NodegroupName: aws.String(nodegroupName),
					Labels: &ekstypes.UpdateLabelsPayload{
						AddOrUpdateLabels: eksLabels,
					},
				}).Return(&eks.UpdateNodegroupConfigOutput{}, nil)

				Expect(manager.Set(context.Background(), nodegroupName, labels)).To(Succeed())
			})

			When("the EKS api returns an error", func() {
				It("fails", func() {
					mockProvider.MockEKS().On("UpdateNodegroupConfig", mock.Anything, mock.Anything).Return(&eks.UpdateNodegroupConfigOutput{}, errors.New("oh-noes"))

					err := manager.Set(context.Background(), nodegroupName, labels)
					Expect(err).To(HaveOccurred())
				})
			})
		})
	})

	Describe("Unset", func() {
		var labels []string

		BeforeEach(func() {
			labels = []string{"k1"}
		})

		When("the nodegroup is owned by eksctl", func() {
			BeforeEach(func() {
				fakeManagedService.UpdateLabelsReturns(nil)
			})

			It("removes labels by updating the nodegroup stack", func() {
				Expect(manager.Unset(context.Background(), nodegroupName, labels)).To(Succeed())
			})

			When("the service returns an error", func() {
				BeforeEach(func() {
					fakeManagedService.UpdateLabelsReturns(errors.New("something-terrible"))
				})

				It("fails", func() {
					err := manager.Unset(context.Background(), nodegroupName, labels)
					Expect(err).To(HaveOccurred())
				})
			})
		})

		When("the nodegroup is not owned by eksctl", func() {
			var eksLabels []string

			BeforeEach(func() {
				eksLabels = []string{"k1"}
				err := &smithy.OperationError{
					Err: errors.New("ValidationError"),
				}
				fakeManagedService.UpdateLabelsReturns(perrors.Wrapf(err, "omg %s", "what"))
			})

			It("removes the labels through the EKS api", func() {
				mockProvider.MockEKS().On("UpdateNodegroupConfig", mock.Anything, &eks.UpdateNodegroupConfigInput{
					ClusterName:   aws.String(clusterName),
					NodegroupName: aws.String(nodegroupName),
					Labels: &ekstypes.UpdateLabelsPayload{
						RemoveLabels: eksLabels,
					},
				}).Return(&eks.UpdateNodegroupConfigOutput{}, nil)

				Expect(manager.Unset(context.Background(), nodegroupName, labels)).To(Succeed())
			})

			When("the EKS api returns an error", func() {
				It("fails", func() {
					mockProvider.MockEKS().On("UpdateNodegroupConfig", mock.Anything, mock.Anything).Return(&eks.UpdateNodegroupConfigOutput{}, errors.New("oh-noes"))

					err := manager.Unset(context.Background(), nodegroupName, labels)
					Expect(err).To(HaveOccurred())
				})
			})
		})
	})
})
