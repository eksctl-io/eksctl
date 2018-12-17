package ami_test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	. "github.com/weaveworks/eksctl/pkg/ami"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

type returnAmi struct {
	imageId     string
	state       string
	createdDate string
}

var _ = Describe("AMI Auto Resolution", func() {

	Describe("When resolving an AMI to use", func() {

		var (
			p            *testutils.MockProvider
			err          error
			region       string
			version      string
			instanceType string
			imageFamily  string
			resolvedAmi  string
			expectedAmi  string
			imageState   string
		)

		Context("with a valid region and N instance type", func() {
			BeforeEach(func() {
				region = "eu-west-1"
				version = "1.10"
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
						addMockDescribeImages(p, "amazon-eks-node-1.10-v*", expectedAmi, imageState, "2018-08-20T23:25:53.000Z")

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
						addMockDescribeImagesMultiple(p, "amazon-eks-node-1.10-v*", []returnAmi{})

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
							returnAmi{
								createdDate: "2018-08-20T23:25:53.000Z",
								imageId:     "ami-1234",
								state:       "available",
							},
							returnAmi{
								createdDate: "2018-09-12T22:21:11.000Z",
								imageId:     expectedAmi,
								state:       "available",
							},
						}

						addMockDescribeImagesMultiple(p, "amazon-eks-node-1.10-v*", images)

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
						addMockDescribeImages(p, "amazon-eks-gpu-node-1.10-v*", expectedAmi, imageState, "2018-08-20T23:25:53.000Z")

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

func createProviders() (*eks.ClusterProvider, *testutils.MockProvider) {
	p := testutils.NewMockProvider()

	c := &eks.ClusterProvider{
		Provider: p,
	}

	return c, p
}

func addMockDescribeImages(p *testutils.MockProvider, expectedNamePattern string, amiId string, amiState string, createdDate string) {
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
			&ec2.Image{
				ImageId: aws.String(amiId),
				State:   aws.String(amiState),
			},
		},
	}, nil)
}

func addMockDescribeImagesMultiple(p *testutils.MockProvider, expectedNamePattern string, returnAmis []returnAmi) {
	images := make([]*ec2.Image, len(returnAmis))
	for index, ami := range returnAmis {
		images[index] = &ec2.Image{
			ImageId:      aws.String(ami.imageId),
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
