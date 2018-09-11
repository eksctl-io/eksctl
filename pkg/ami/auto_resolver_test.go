package ami_test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	. "github.com/weaveworks/eksctl/pkg/ami"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/eks/api"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

var _ = Describe("AMI Auto Resolution", func() {

	Describe("When resolving an AMI to use", func() {
		var (
			p            *testutils.MockProvider
			err          error
			region       string
			instanceType string
			resolvedAmi  string
			expectedAmi  string
			imageState   string
		)

		Context("with a valid region and N instance type", func() {
			BeforeEach(func() {
				region = "eu-west-1"
				expectedAmi = "ami-12345"
			})

			Context("and non-gpu instance type", func() {
				BeforeEach(func() {
					instanceType = "t2.medium"
				})

				Context("and ami is available", func() {
					BeforeEach(func() {
						imageState = "available"

						_, p = createProviders()
						addMockDescribeImages(p, "amazon-eks-node-*", expectedAmi, imageState)

						resolver := NewAutoResolver(p.MockEC2())
						resolvedAmi, err = resolver.Resolve(region, instanceType)
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

				Context("and ami is pending", func() {
					BeforeEach(func() {
						imageState = "pending"

						_, p = createProviders()
						addMockDescribeImages(p, "amazon-eks-node-*", expectedAmi, imageState)

						resolver := NewAutoResolver(p.MockEC2())
						resolvedAmi, err = resolver.Resolve(region, instanceType)
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
			})

			Context("and gpu instance type", func() {
				BeforeEach(func() {
					instanceType = "p2.xlarge"
				})

				Context("and ami is available", func() {
					BeforeEach(func() {
						imageState = "available"

						_, p = createProviders()
						addMockDescribeImages(p, "amazon-eks-gpu-node-*", expectedAmi, imageState)

						resolver := NewAutoResolver(p.MockEC2())
						resolvedAmi, err = resolver.Resolve(region, instanceType)
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
		Spec: &api.ClusterConfig{
			Region: "eu-west-1",
		},
	}

	return c, p
}

func addMockDescribeImages(p *testutils.MockProvider, expectedNamePattern string, amiId string, amiState string) {
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
