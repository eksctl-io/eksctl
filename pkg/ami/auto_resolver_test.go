package ami_test

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	. "github.com/weaveworks/eksctl/pkg/ami"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

type returnAmi struct {
	imageID     string
	state       ec2types.ImageState
	createdDate string
}

var _ = Describe("AMI Auto Resolution", func() {

	Describe("When resolving an AMI to use", func() {

		var (
			p            *mockprovider.MockProvider
			err          error
			region       string
			version      string
			instanceType string
			imageFamily  string
			resolvedAmi  string
			expectedAmi  string
			imageState   ec2types.ImageState
		)

		Context("setting proper AWS Account IDs based on instance families", func() {
			It("should return the AWS Account ID for AL2 images", func() {
				ownerAccount, err := OwnerAccountID(api.NodeImageFamilyAmazonLinux2, region)
				Expect(ownerAccount).To(BeEquivalentTo("602401143452"))
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return the AWS Account ID for AL2 images in ap-east-1", func() {
				ownerAccount, err := OwnerAccountID(api.NodeImageFamilyAmazonLinux2, "ap-east-1")
				Expect(ownerAccount).To(BeEquivalentTo("800184023465"))
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return the AWS Account ID for Ubuntu images in ap-east-1", func() {
				ownerAccount, err := OwnerAccountID(api.NodeImageFamilyUbuntu1804, "ap-east-1")
				Expect(ownerAccount).To(BeEquivalentTo("099720109477"))
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return the Ubuntu Account ID for Ubuntu images", func() {
				ownerAccount, err := OwnerAccountID(api.NodeImageFamilyUbuntu1804, region)
				Expect(ownerAccount).To(BeEquivalentTo("099720109477"))
				Expect(err).NotTo(HaveOccurred())
			})
			It("should return the Ubuntu Account ID for Ubuntu images", func() {
				ownerAccount, err := OwnerAccountID(api.NodeImageFamilyUbuntu2204, region)
				Expect(ownerAccount).To(BeEquivalentTo("099720109477"))
				Expect(err).NotTo(HaveOccurred())
			})
			It("should return the Ubuntu Account ID for Ubuntu images", func() {
				ownerAccount, err := OwnerAccountID(api.NodeImageFamilyUbuntuPro2204, region)
				Expect(ownerAccount).To(BeEquivalentTo("099720109477"))
				Expect(err).NotTo(HaveOccurred())
			})
			It("should return the Ubuntu Account ID for Ubuntu images", func() {
				ownerAccount, err := OwnerAccountID(api.NodeImageFamilyUbuntu2404, region)
				Expect(ownerAccount).To(BeEquivalentTo("099720109477"))
				Expect(err).NotTo(HaveOccurred())
			})
			It("should return the Ubuntu Account ID for Ubuntu images", func() {
				ownerAccount, err := OwnerAccountID(api.NodeImageFamilyUbuntuPro2404, region)
				Expect(ownerAccount).To(BeEquivalentTo("099720109477"))
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return the Windows Account ID for Windows Server images", func() {
				ownerAccount, err := OwnerAccountID(api.NodeImageFamilyWindowsServer2022CoreContainer, region)
				Expect(ownerAccount).To(BeEquivalentTo("801119661308"))
				Expect(err).NotTo(HaveOccurred())
			})

		})

		Context("with a valid region and N instance type", func() {
			BeforeEach(func() {
				region = "eu-west-1"
				version = "1.15"
				expectedAmi = "ami-12345"
			})

			Context("and non-gpu instance type", func() {
				BeforeEach(func() {
					instanceType = "t2.medium"
					imageFamily = "AmazonLinux2"
				})

				Context("and ami is available", func() {
					BeforeEach(func() {
						imageState = ec2types.ImageStateAvailable

						p = mockprovider.NewMockProvider()
						addMockDescribeImages(p, "amazon-eks-node-1.15-v*", expectedAmi, imageState, "2018-08-20T23:25:53.000Z", api.NodeImageFamilyAmazonLinux2)
						resolver := NewAutoResolver(p.MockEC2())
						resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)
					})

					It("should not error", func() {
						Expect(err).NotTo(HaveOccurred())
					})

					It("should have called AWS EC2 DescribeImages", func() {
						Expect(p.MockEC2().AssertNumberOfCalls(GinkgoT(), "DescribeImages", 1)).To(BeTrue())
					})

					It("should have returned an ami id", func() {
						Expect(resolvedAmi).To(BeEquivalentTo(expectedAmi))
					})
				})

				Context("and ami is available", func() {
					BeforeEach(func() {
						imageState = ec2types.ImageStateAvailable
						imageFamily = "Ubuntu1804"

						p = mockprovider.NewMockProvider()
						addMockDescribeImages(p, "ubuntu-eks/k8s_1.15/images/*18.04*", expectedAmi, imageState, "2018-08-20T23:25:53.000Z", api.NodeImageFamilyUbuntu1804)

						resolver := NewAutoResolver(p.MockEC2())
						resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)
					})

					It("should not error", func() {
						Expect(err).NotTo(HaveOccurred())
					})

					It("should have called AWS EC2 DescribeImages", func() {
						Expect(p.MockEC2().AssertNumberOfCalls(GinkgoT(), "DescribeImages", 1)).To(BeTrue())
					})

					It("should have returned an ami id", func() {
						Expect(resolvedAmi).To(BeEquivalentTo(expectedAmi))
					})
				})

				Context("and ami is NOT available", func() {
					BeforeEach(func() {
						imageState = ec2types.ImageStatePending

						p = mockprovider.NewMockProvider()
						addMockDescribeImagesMultiple(p, "amazon-eks-node-1.15-v*", []returnAmi{})

						resolver := NewAutoResolver(p.MockEC2())
						resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)
					})

					It("should not error", func() {
						Expect(err).NotTo(HaveOccurred())
					})

					It("should have called AWS EC2 DescribeImages", func() {
						Expect(p.MockEC2().AssertNumberOfCalls(GinkgoT(), "DescribeImages", 1)).To(BeTrue())
					})

					It("should NOT have returned an ami id", func() {
						Expect(resolvedAmi).To(BeEquivalentTo(""))
					})
				})

				Context("and there are 2 ami's available", func() {
					BeforeEach(func() {
						imageState = "available"
						expectedAmi = "ami-5678"

						p = mockprovider.NewMockProvider()
						images := []returnAmi{
							{
								createdDate: "2018-08-20T23:25:53.000Z",
								imageID:     "ami-1234",
								state:       "available",
							},
							{
								createdDate: "2018-09-12T22:21:11.000Z",
								imageID:     expectedAmi,
								state:       "available",
							},
						}

						addMockDescribeImagesMultiple(p, "amazon-eks-node-1.15-v*", images)

						resolver := NewAutoResolver(p.MockEC2())
						resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)
					})

					It("should not error", func() {
						Expect(err).NotTo(HaveOccurred())
					})

					It("should have called AWS EC2 DescribeImages", func() {
						Expect(p.MockEC2().AssertNumberOfCalls(GinkgoT(), "DescribeImages", 1)).To(BeTrue())
					})

					It("should have returned an ami id", func() {
						Expect(resolvedAmi).To(BeEquivalentTo(expectedAmi))
					})
				})
			})

			Context("and gpu instance type", func() {
				BeforeEach(func() {
					instanceType = "p2.xlarge"
				})

				Context("and ami is available", func() {
					BeforeEach(func() {
						imageState = "available"

						p = mockprovider.NewMockProvider()
						addMockDescribeImages(p, "amazon-eks-gpu-node-1.15-*", expectedAmi, imageState, "2018-08-20T23:25:53.000Z", api.NodeImageFamilyAmazonLinux2)
						resolver := NewAutoResolver(p.MockEC2())
						resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)
					})

					It("should not error", func() {
						Expect(err).NotTo(HaveOccurred())
					})

					It("should have called AWS EC2 DescribeImages", func() {
						Expect(p.MockEC2().AssertNumberOfCalls(GinkgoT(), "DescribeImages", 1)).To(BeTrue())
					})

					It("should have returned an ami id", func() {
						Expect(resolvedAmi).To(BeEquivalentTo(expectedAmi))
					})
				})
			})
		})
	})
})

func addMockDescribeImages(p *mockprovider.MockProvider, expectedNamePattern string, amiID string, amiState ec2types.ImageState, createdDate string, instanceFamily string) {
	p.MockEC2().On("DescribeImages",
		mock.Anything,
		mock.MatchedBy(func(input *ec2.DescribeImagesInput) bool {
			for _, filter := range input.Filters {
				if *filter.Name == "name" {
					if len(filter.Values) > 0 {
						if filter.Values[0] == expectedNamePattern {
							return true
						}
					}
				}
			}
			return false
		}),
	).Return(&ec2.DescribeImagesOutput{
		Images: []ec2types.Image{
			{
				ImageId:      aws.String(amiID),
				State:        amiState,
				CreationDate: aws.String(createdDate),
				Description:  aws.String(instanceFamily),
			},
		},
	}, nil)
}

func addMockDescribeImagesMultiple(p *mockprovider.MockProvider, expectedNamePattern string, returnAmis []returnAmi) {
	images := make([]ec2types.Image, len(returnAmis))
	for i, ami := range returnAmis {
		images[i] = ec2types.Image{
			ImageId:      aws.String(ami.imageID),
			State:        ami.state,
			CreationDate: aws.String(ami.createdDate),
		}
	}

	p.MockEC2().On("DescribeImages",
		mock.Anything,
		mock.MatchedBy(func(input *ec2.DescribeImagesInput) bool {
			for _, filter := range input.Filters {
				if *filter.Name == "name" {
					if len(filter.Values) > 0 {
						if filter.Values[0] == expectedNamePattern {
							return true
						}
					}
				}
			}
			return false
		}),
	).Return(&ec2.DescribeImagesOutput{
		Images: images,
	}, nil)
}
