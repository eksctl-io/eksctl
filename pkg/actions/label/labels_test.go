package label_test

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
				summary, err := manager.Get(nodegroupName)
				Expect(err).NotTo(HaveOccurred())
				Expect(summary[0].Labels).To(Equal(expectedLabels))
			})

			When("the service returns an error", func() {
				BeforeEach(func() {
					fakeManagedService.GetLabelsReturns(nil, errors.New("something-terrible"))
				})

				It("fails", func() {
					summary, err := manager.Get(nodegroupName)
					Expect(err).To(HaveOccurred())
					Expect(summary).To(BeNil())
				})
			})
		})

		When("the nodegroup is not owned by eksctl", func() {
			var returnedLabels map[string]*string

			BeforeEach(func() {
				returnedLabels = map[string]*string{"k1": aws.String("v1")}
				fakeManagedService.GetLabelsReturns(nil, awserr.New("ValidationError", "stack not found", errors.New("omg")))
			})

			It("returns the labels from the EKS api", func() {
				mockProvider.MockEKS().On("DescribeNodegroup", &awseks.DescribeNodegroupInput{
					ClusterName:   aws.String(clusterName),
					NodegroupName: aws.String(nodegroupName),
				}).Return(&awseks.DescribeNodegroupOutput{Nodegroup: &awseks.Nodegroup{Labels: returnedLabels}}, nil)

				summary, err := manager.Get(nodegroupName)
				Expect(err).NotTo(HaveOccurred())
				Expect(summary[0].Labels).To(Equal(expectedLabels))
			})

			When("the EKS api returns an error", func() {
				BeforeEach(func() {
					fakeManagedService.GetLabelsReturns(nil, awserr.New("ValidationError", "stack not found", errors.New("omg")))
				})

				It("fails", func() {
					mockProvider.MockEKS().On("DescribeNodegroup", mock.Anything).Return(&awseks.DescribeNodegroupOutput{}, errors.New("oh-noes"))

					summary, err := manager.Get(nodegroupName)
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
				Expect(manager.Set(nodegroupName, labels)).To(Succeed())
			})

			When("the service returns an error", func() {
				BeforeEach(func() {
					fakeManagedService.UpdateLabelsReturns(errors.New("something-terrible"))
				})

				It("fails", func() {
					err := manager.Set(nodegroupName, labels)
					Expect(err).To(HaveOccurred())
				})
			})
		})

		When("the nodegroup is not owned by eksctl", func() {
			var eksLabels map[string]*string

			BeforeEach(func() {
				eksLabels = map[string]*string{"k1": aws.String("v1")}
				fakeManagedService.UpdateLabelsReturns(awserr.New("ValidationError", "stack not found", errors.New("omg")))
			})

			It("updates the labels through the EKS api", func() {
				mockProvider.MockEKS().On("UpdateNodegroupConfig", &awseks.UpdateNodegroupConfigInput{
					ClusterName:   aws.String(clusterName),
					NodegroupName: aws.String(nodegroupName),
					Labels: &awseks.UpdateLabelsPayload{
						AddOrUpdateLabels: eksLabels,
					},
				}).Return(&awseks.UpdateNodegroupConfigOutput{}, nil)

				Expect(manager.Set(nodegroupName, labels)).To(Succeed())
			})

			When("the EKS api returns an error", func() {
				BeforeEach(func() {
					fakeManagedService.GetLabelsReturns(nil, awserr.New("ValidationError", "stack not found", errors.New("omg")))
				})

				It("fails", func() {
					mockProvider.MockEKS().On("UpdateNodegroupConfig", mock.Anything).Return(&awseks.UpdateNodegroupConfigOutput{}, errors.New("oh-noes"))

					err := manager.Set(nodegroupName, labels)
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
				Expect(manager.Unset(nodegroupName, labels)).To(Succeed())
			})

			When("the service returns an error", func() {
				BeforeEach(func() {
					fakeManagedService.UpdateLabelsReturns(errors.New("something-terrible"))
				})

				It("fails", func() {
					err := manager.Unset(nodegroupName, labels)
					Expect(err).To(HaveOccurred())
				})
			})
		})

		When("the nodegroup is not owned by eksctl", func() {
			var eksLabels []*string

			BeforeEach(func() {
				eksLabels = []*string{aws.String("k1")}
				fakeManagedService.UpdateLabelsReturns(awserr.New("ValidationError", "stack not found", errors.New("omg")))
			})

			It("removes the labels through the EKS api", func() {
				mockProvider.MockEKS().On("UpdateNodegroupConfig", &awseks.UpdateNodegroupConfigInput{
					ClusterName:   aws.String(clusterName),
					NodegroupName: aws.String(nodegroupName),
					Labels: &awseks.UpdateLabelsPayload{
						RemoveLabels: eksLabels,
					},
				}).Return(&awseks.UpdateNodegroupConfigOutput{}, nil)

				Expect(manager.Unset(nodegroupName, labels)).To(Succeed())
			})

			When("the EKS api returns an error", func() {
				BeforeEach(func() {
					fakeManagedService.GetLabelsReturns(nil, awserr.New("ValidationError", "stack not found", errors.New("omg")))
				})

				It("fails", func() {
					mockProvider.MockEKS().On("UpdateNodegroupConfig", mock.Anything).Return(&awseks.UpdateNodegroupConfigOutput{}, errors.New("oh-noes"))

					err := manager.Unset(nodegroupName, labels)
					Expect(err).To(HaveOccurred())
				})
			})
		})
	})
})
