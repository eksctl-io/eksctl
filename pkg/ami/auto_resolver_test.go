package ami_test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	. "github.com/weaveworks/eksctl/pkg/ami"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

type returnAmi struct {
	imageID     string
	state       string
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
			imageState   string
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
						imageState = "available"

						_, p = createProviders()
						addMockDescribeImages(p, "amazon-eks-node-1.15-v*", expectedAmi, imageState, "2018-08-20T23:25:53.000Z", api.NodeImageFamilyAmazonLinux2)
						resolver := NewAutoResolver(p.MockEC2())
						resolvedAmi, err = resolver.Resolve(region, version, instanceType, imageFamily)
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
						imageState = "available"
						imageFamily = "Ubuntu1804"

						_, p = createProviders()
						addMockDescribeImages(p, "ubuntu-eks/k8s_1.15/images/*18.04*", expectedAmi, imageState, "2018-08-20T23:25:53.000Z", api.NodeImageFamilyUbuntu1804)

						resolver := NewAutoResolver(p.MockEC2())
						resolvedAmi, err = resolver.Resolve(region, version, instanceType, imageFamily)
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
						imageState = "pending"

						_, p = createProviders()
						addMockDescribeImagesMultiple(p, "amazon-eks-node-1.15-v*", []returnAmi{})

						resolver := NewAutoResolver(p.MockEC2())
						resolvedAmi, err = resolver.Resolve(region, version, instanceType, imageFamily)
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

						_, p = createProviders()
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
						resolvedAmi, err = resolver.Resolve(region, version, instanceType, imageFamily)
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

						_, p = createProviders()
						addMockDescribeImages(p, "amazon-eks-gpu-node-1.15-*", expectedAmi, imageState, "2018-08-20T23:25:53.000Z", api.NodeImageFamilyAmazonLinux2)
						resolver := NewAutoResolver(p.MockEC2())
						resolvedAmi, err = resolver.Resolve(region, version, instanceType, imageFamily)
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

func createProviders() (*eks.ClusterProvider, *mockprovider.MockProvider) {
	p := mockprovider.NewMockProvider()

	c := &eks.ClusterProvider{
		Provider: p,
	}

	return c, p
}

func addMockDescribeImages(p *mockprovider.MockProvider, expectedNamePattern string, amiID string, amiState string, createdDate string, instanceFamily string) {
	p.MockEC2().On("DescribeImages",
		mock.MatchedBy(func(input *ec2.DescribeImagesInput) bool {
			for _, filter := range input.Filters {
				if *filter.Name == "name" {
					if len(filter.Values) > 0 {
						if *filter.Values[0] == expectedNamePattern {
							return true
						}
					}
				}
			}
			return false
		}),
	).Return(&ec2.DescribeImagesOutput{
		Images: []*ec2.Image{
			{
				ImageId:      aws.String(amiID),
				State:        aws.String(amiState),
				CreationDate: aws.String(createdDate),
				Description:  aws.String(instanceFamily),
			},
		},
	}, nil)
}

func addMockDescribeImagesMultiple(p *mockprovider.MockProvider, expectedNamePattern string, returnAmis []returnAmi) {
	images := make([]*ec2.Image, len(returnAmis))
	for index, ami := range returnAmis {
		images[index] = &ec2.Image{
			ImageId:      aws.String(ami.imageID),
			State:        aws.String(ami.state),
			CreationDate: aws.String(ami.createdDate),
		}
	}

	p.MockEC2().On("DescribeImages",
		mock.MatchedBy(func(input *ec2.DescribeImagesInput) bool {
			for _, filter := range input.Filters {
				if *filter.Name == "name" {
					if len(filter.Values) > 0 {
						if *filter.Values[0] == expectedNamePattern {
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
