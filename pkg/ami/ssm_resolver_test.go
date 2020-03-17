package ami_test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	. "github.com/weaveworks/eksctl/pkg/ami"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

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
		)

		Context("with a valid region and N instance type", func() {
			BeforeEach(func() {
				region = "eu-west-1"
				version = "1.12"
				expectedAmi = "ami-12345"
			})

			Context("and non-gpu instance type", func() {
				BeforeEach(func() {
					instanceType = "t2.medium"
					imageFamily = "AmazonLinux2"
				})

				Context("and AL2 ami is available", func() {
					BeforeEach(func() {

						_, p = createProviders()
						addMockGetParameter(p, "/aws/service/eks/optimized-ami/1.12/amazon-linux-2/recommended/image_id", expectedAmi)
						resolver := NewSSMResolver(p.MockSSM())
						resolvedAmi, err = resolver.Resolve(region, version, instanceType, imageFamily)
					})

					It("should not error", func() {
						Expect(err).NotTo(HaveOccurred())
					})

					It("should have called AWS SSM GetParameter", func() {
						Expect(p.MockSSM().AssertNumberOfCalls(GinkgoT(), "GetParameter", 1)).To(BeTrue())
					})

					It("should have returned an ami id", func() {
						Expect(resolvedAmi).To(BeEquivalentTo(expectedAmi))
					})
				})

				Context("and ami is NOT available", func() {
					BeforeEach(func() {

						_, p = createProviders()
						addMockFailedGetParameter(p, "/aws/service/eks/optimized-ami/1.12/amazon-linux-2/recommended/image_id")

						resolver := NewSSMResolver(p.MockSSM())
						resolvedAmi, err = resolver.Resolve(region, version, instanceType, imageFamily)
					})

					It("should return an error", func() {
						Expect(err).To(HaveOccurred())
					})

					It("should NOT have returned an ami id", func() {
						Expect(resolvedAmi).To(BeEquivalentTo(""))
					})

					It("should have called AWS SSM GetParameter", func() {
						Expect(p.MockSSM().AssertNumberOfCalls(GinkgoT(), "GetParameter", 1)).To(BeTrue())
					})

				})
			})

			Context("and gpu instance type", func() {
				BeforeEach(func() {
					instanceType = "p2.xlarge"
				})

				Context("and ami is available", func() {
					BeforeEach(func() {

						_, p = createProviders()
						addMockGetParameter(p, "/aws/service/eks/optimized-ami/1.12/amazon-linux-2-gpu/recommended/image_id", expectedAmi)
						resolver := NewSSMResolver(p.MockSSM())
						resolvedAmi, err = resolver.Resolve(region, version, instanceType, imageFamily)
					})

					It("should not error", func() {
						Expect(err).NotTo(HaveOccurred())
					})

					It("should have called AWS SSM GetParameter", func() {
						Expect(p.MockSSM().AssertNumberOfCalls(GinkgoT(), "GetParameter", 1)).To(BeTrue())
					})

					It("should have returned an ami id", func() {
						Expect(resolvedAmi).To(BeEquivalentTo(expectedAmi))
					})
				})
			})

			Context("and Windows Core family", func() {
				BeforeEach(func() {
					instanceType = "t3.xlarge"
				})

				Context("and ami is available", func() {
					BeforeEach(func() {
						version = "1.14"
						_, p = createProviders()
					})

					It("should return a valid Full image for 1.14", func() {
						imageFamily = "WindowsServer2019FullContainer"
						addMockGetParameter(p, "/aws/service/ami-windows-latest/Windows_Server-2019-English-Full-EKS_Optimized-1.14/image_id", expectedAmi)

						resolver := NewSSMResolver(p.MockSSM())
						resolvedAmi, err = resolver.Resolve(region, version, instanceType, imageFamily)

						Expect(err).NotTo(HaveOccurred())
						Expect(resolvedAmi).To(BeEquivalentTo(expectedAmi))
						Expect(p.MockSSM().AssertNumberOfCalls(GinkgoT(), "GetParameter", 1)).To(BeTrue())
					})

					It("should return a valid Core image for 1.15", func() {
						imageFamily = "WindowsServer2019CoreContainer"
						addMockGetParameter(p, "/aws/service/ami-windows-latest/Windows_Server-2019-English-Core-EKS_Optimized-1.15/image_id", expectedAmi)

						resolver := NewSSMResolver(p.MockSSM())
						resolvedAmi, err = resolver.Resolve(region, "1.15", instanceType, imageFamily)

						Expect(err).NotTo(HaveOccurred())
						Expect(resolvedAmi).To(BeEquivalentTo(expectedAmi))
						Expect(p.MockSSM().AssertNumberOfCalls(GinkgoT(), "GetParameter", 1)).To(BeTrue())
					})

					It("should return an error for versions below 1.14", func() {
						resolver := NewSSMResolver(p.MockSSM())
						resolvedAmi, err = resolver.Resolve(region, "1.13", instanceType, imageFamily)

						Expect(err).To(HaveOccurred())
					})
				})
			})

			Context("and Ubuntu family", func() {
				BeforeEach(func() {
					_, p = createProviders()
				})

				It("should return an error", func() {
					resolver := NewSSMResolver(p.MockSSM())
					resolvedAmi, err = resolver.Resolve(region, version, instanceType, imageFamily)

					Expect(err).To(HaveOccurred())
				})
			})

			Context("and Bottlerocket image family", func() {
				BeforeEach(func() {
					instanceType = "t2.medium"
					imageFamily = "Bottlerocket"
					version = "1.15"
				})

				Context("and ami is available", func() {
					BeforeEach(func() {
						_, p = createProviders()
						addMockGetParameter(p, "/aws/service/bottlerocket/aws-k8s-1.15/x86_64/latest/image_id", expectedAmi)
						resolver := NewSSMResolver(p.MockSSM())
						resolvedAmi, err = resolver.Resolve(region, version, instanceType, imageFamily)
					})

					It("should not error", func() {
						Expect(err).NotTo(HaveOccurred())
					})

					It("should have called AWS SSM GetParameter", func() {
						Expect(p.MockSSM().AssertNumberOfCalls(GinkgoT(), "GetParameter", 1)).To(BeTrue())
					})

					It("should have returned an ami id", func() {
						Expect(resolvedAmi).To(BeEquivalentTo(expectedAmi))
					})
				})


				Context("and ami is NOT available", func() {
					BeforeEach(func() {
						_, p = createProviders()
						addMockFailedGetParameter(p, "/aws/service/bottlerocket/aws-k8s-1.15/x86_64/latest/image_id")

						resolver := NewSSMResolver(p.MockSSM())
						resolvedAmi, err = resolver.Resolve(region, version, instanceType, imageFamily)
					})

					It("should return an error", func() {
						Expect(err).To(HaveOccurred())
					})

					It("should NOT have returned an ami id", func() {
						Expect(resolvedAmi).To(BeEquivalentTo(""))
					})

					It("should have called AWS SSM GetParameter", func() {
						Expect(p.MockSSM().AssertNumberOfCalls(GinkgoT(), "GetParameter", 1)).To(BeTrue())
					})

				})


				Context("for arm instance type", func() {
					BeforeEach(func() {
						instanceType = "a1.large"
					})

					Context("and ami is available", func() {
						BeforeEach(func() {
							_, p = createProviders()
							addMockGetParameter(p, "/aws/service/bottlerocket/aws-k8s-1.15/arm64/latest/image_id", expectedAmi)
							resolver := NewSSMResolver(p.MockSSM())
							resolvedAmi, err = resolver.Resolve(region, version, instanceType, imageFamily)
						})

						It("should not error", func() {
							Expect(err).NotTo(HaveOccurred())
						})

						It("should have called AWS SSM GetParameter", func() {
							Expect(p.MockSSM().AssertNumberOfCalls(GinkgoT(), "GetParameter", 1)).To(BeTrue())
						})

						It("should have returned an ami id", func() {
							Expect(resolvedAmi).To(BeEquivalentTo(expectedAmi))
						})
					})

					Context("and ami is NOT available", func() {
						BeforeEach(func() {
							_, p = createProviders()
							addMockFailedGetParameter(p, "/aws/service/bottlerocket/aws-k8s-1.15/arm64/latest/image_id")

							resolver := NewSSMResolver(p.MockSSM())
							resolvedAmi, err = resolver.Resolve(region, version, instanceType, imageFamily)
						})

						It("should return an error", func() {
							Expect(err).To(HaveOccurred())
						})

						It("should NOT have returned an ami id", func() {
							Expect(resolvedAmi).To(BeEquivalentTo(""))
						})

						It("should have called AWS SSM GetParameter", func() {
							Expect(p.MockSSM().AssertNumberOfCalls(GinkgoT(), "GetParameter", 1)).To(BeTrue())
						})

					})

				})
			})


		})
	})
})

func addMockGetParameter(p *mockprovider.MockProvider, name string, amiId string) {
	p.MockSSM().On("GetParameter",
		mock.MatchedBy(func(input *ssm.GetParameterInput) bool {
			return *input.Name == name
		}),
	).Return(&ssm.GetParameterOutput{
		Parameter: &ssm.Parameter{
			Name:  aws.String(name),
			Type:  aws.String("String"),
			Value: aws.String(amiId),
		},
	}, nil)
}

func addMockFailedGetParameter(p *mockprovider.MockProvider, name string) {
	p.MockSSM().On("GetParameter",
		mock.MatchedBy(func(input *ssm.GetParameterInput) bool {
			return *input.Name == name
		}),
	).Return(&ssm.GetParameterOutput{
		Parameter: nil,
	}, nil)
}
