package ami_test

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	. "github.com/weaveworks/eksctl/pkg/ami"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
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

						p = mockprovider.NewMockProvider()
						addMockGetParameter(p, "/aws/service/eks/optimized-ami/1.12/amazon-linux-2/recommended/image_id", expectedAmi)
						resolver := NewSSMResolver(p.MockSSM())
						resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)
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

				Context("and ami is not available", func() {
					BeforeEach(func() {

						p = mockprovider.NewMockProvider()
						addMockFailedGetParameter(p, "/aws/service/eks/optimized-ami/1.12/amazon-linux-2/recommended/image_id")

						resolver := NewSSMResolver(p.MockSSM())
						resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)
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
					instanceType = "g5.xlarge"
				})

				Context("and ami is available", func() {
					BeforeEach(func() {

						p = mockprovider.NewMockProvider()
						addMockGetParameter(p, "/aws/service/eks/optimized-ami/1.12/amazon-linux-2-gpu/recommended/image_id", expectedAmi)
						resolver := NewSSMResolver(p.MockSSM())
						resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)
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
						p = mockprovider.NewMockProvider()
					})

					It("should return a valid Core image for 1.15", func() {
						imageFamily = "WindowsServer2019CoreContainer"
						addMockGetParameter(p, "/aws/service/ami-windows-latest/Windows_Server-2019-English-Core-EKS_Optimized-1.15/image_id", expectedAmi)

						resolver := NewSSMResolver(p.MockSSM())
						resolvedAmi, err = resolver.Resolve(context.Background(), region, "1.15", instanceType, imageFamily)

						Expect(err).NotTo(HaveOccurred())
						Expect(resolvedAmi).To(BeEquivalentTo(expectedAmi))
						Expect(p.MockSSM().AssertNumberOfCalls(GinkgoT(), "GetParameter", 1)).To(BeTrue())
					})
				})

				Context("Windows Server 2022 Core", func() {
					BeforeEach(func() {
						version = "1.23"
						p = mockprovider.NewMockProvider()
					})

					It("should return a valid AMI", func() {
						imageFamily = "WindowsServer2022CoreContainer"
						addMockGetParameter(p, "/aws/service/ami-windows-latest/Windows_Server-2022-English-Core-EKS_Optimized-1.23/image_id", expectedAmi)

						resolver := NewSSMResolver(p.MockSSM())
						resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)

						Expect(err).NotTo(HaveOccurred())
						Expect(resolvedAmi).To(BeEquivalentTo(expectedAmi))
						Expect(p.MockSSM().AssertNumberOfCalls(GinkgoT(), "GetParameter", 1)).To(BeTrue())
					})

					It("should return an error for EKS versions below 1.23", func() {
						imageFamily = "WindowsServer2022CoreContainer"

						resolver := NewSSMResolver(p.MockSSM())
						resolvedAmi, err = resolver.Resolve(context.Background(), region, "1.22", instanceType, imageFamily)

						Expect(err).To(HaveOccurred())
						Expect(err).To(MatchError(ContainSubstring("Windows Server 2022 Core requires EKS version 1.23 and above")))
					})
				})

				Context("Windows Server 2025 Core", func() {
					BeforeEach(func() {
						version = "1.35"
						p = mockprovider.NewMockProvider()
					})

					It("should return a valid AMI", func() {
						imageFamily = "WindowsServer2025CoreContainer"
						addMockGetParameter(p, "/aws/service/ami-windows-latest/Windows_Server-2025-English-Core-EKS_Optimized-1.35/image_id", expectedAmi)

						resolver := NewSSMResolver(p.MockSSM())
						resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)

						Expect(err).NotTo(HaveOccurred())
						Expect(resolvedAmi).To(BeEquivalentTo(expectedAmi))
						Expect(p.MockSSM().AssertNumberOfCalls(GinkgoT(), "GetParameter", 1)).To(BeTrue())
					})

					It("should return an error for EKS versions below 1.35", func() {
						imageFamily = "WindowsServer2025CoreContainer"

						resolver := NewSSMResolver(p.MockSSM())
						resolvedAmi, err = resolver.Resolve(context.Background(), region, "1.34", instanceType, imageFamily)

						Expect(err).To(HaveOccurred())
						Expect(err).To(MatchError(ContainSubstring("Windows Server 2025 Core requires EKS version 1.35 and above")))
					})
				})

			})

			Context("and Windows Full family", func() {
				BeforeEach(func() {
					version = "1.14"
					instanceType = "t3.xlarge"
				})

				Context("and ami is available", func() {
					BeforeEach(func() {
						version = "1.14"
						p = mockprovider.NewMockProvider()
					})

					It("should return a valid Full image for 1.14", func() {
						imageFamily = "WindowsServer2019FullContainer"
						addMockGetParameter(p, "/aws/service/ami-windows-latest/Windows_Server-2019-English-Full-EKS_Optimized-1.14/image_id", expectedAmi)

						resolver := NewSSMResolver(p.MockSSM())
						resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)

						Expect(err).NotTo(HaveOccurred())
						Expect(resolvedAmi).To(BeEquivalentTo(expectedAmi))
						Expect(p.MockSSM().AssertNumberOfCalls(GinkgoT(), "GetParameter", 1)).To(BeTrue())
					})

				})

				Context("Windows Server 2022 Full", func() {
					BeforeEach(func() {
						version = "1.23"
						p = mockprovider.NewMockProvider()
					})

					It("should return a valid AMI", func() {
						imageFamily = "WindowsServer2022FullContainer"
						addMockGetParameter(p, "/aws/service/ami-windows-latest/Windows_Server-2022-English-Full-EKS_Optimized-1.23/image_id", expectedAmi)

						resolver := NewSSMResolver(p.MockSSM())
						resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)

						Expect(err).NotTo(HaveOccurred())
						Expect(resolvedAmi).To(BeEquivalentTo(expectedAmi))
						Expect(p.MockSSM().AssertNumberOfCalls(GinkgoT(), "GetParameter", 1)).To(BeTrue())
					})

					It("should return an error for EKS versions below 1.23", func() {
						imageFamily = "WindowsServer2022FullContainer"

						resolver := NewSSMResolver(p.MockSSM())
						resolvedAmi, err = resolver.Resolve(context.Background(), region, "1.22", instanceType, imageFamily)

						Expect(err).To(HaveOccurred())
						Expect(err).To(MatchError(ContainSubstring("Windows Server 2022 Full requires EKS version 1.23 and above")))
					})
				})

				Context("Windows Server 2025 Full", func() {
					BeforeEach(func() {
						version = "1.35"
						p = mockprovider.NewMockProvider()
					})

					It("should return a valid AMI", func() {
						imageFamily = "WindowsServer2025FullContainer"
						addMockGetParameter(p, "/aws/service/ami-windows-latest/Windows_Server-2025-English-Full-EKS_Optimized-1.35/image_id", expectedAmi)

						resolver := NewSSMResolver(p.MockSSM())
						resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)

						Expect(err).NotTo(HaveOccurred())
						Expect(resolvedAmi).To(BeEquivalentTo(expectedAmi))
						Expect(p.MockSSM().AssertNumberOfCalls(GinkgoT(), "GetParameter", 1)).To(BeTrue())
					})

					It("should return an error for EKS versions below 1.34", func() {
						imageFamily = "WindowsServer2025FullContainer"

						resolver := NewSSMResolver(p.MockSSM())
						resolvedAmi, err = resolver.Resolve(context.Background(), region, "1.34", instanceType, imageFamily)

						Expect(err).To(HaveOccurred())
						Expect(err).To(MatchError(ContainSubstring("Windows Server 2025 Full requires EKS version 1.35 and above")))
					})
				})

			})

			Context("and Ubuntu2004 family", func() {
				BeforeEach(func() {
					p = mockprovider.NewMockProvider()
					instanceType = "t2.medium"
					imageFamily = "Ubuntu2004"
				})

				DescribeTable("should return an error",
					func(version string) {
						resolver := NewSSMResolver(p.MockSSM())
						resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)

						Expect(err).To(HaveOccurred())
						Expect(err).To(MatchError("Ubuntu2004 requires EKS version greater or equal than 1.21 and lower than 1.29"))
					},
					EntryDescription("When EKS version is %s"),
					Entry(nil, "1.20"),
					Entry(nil, "1.30"),
				)

				DescribeTable("should return a valid AMI",
					func(version string) {
						addMockGetParameter(p, fmt.Sprintf("/aws/service/canonical/ubuntu/eks/20.04/%s/stable/current/amd64/hvm/ebs-gp2/ami-id", version), expectedAmi)

						resolver := NewSSMResolver(p.MockSSM())
						resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)

						Expect(err).NotTo(HaveOccurred())
						Expect(p.MockSSM().AssertNumberOfCalls(GinkgoT(), "GetParameter", 1)).To(BeTrue())
						Expect(resolvedAmi).To(BeEquivalentTo(expectedAmi))
					},
					EntryDescription("When EKS version is %s"),
					Entry(nil, "1.21"),
					Entry(nil, "1.22"),
					Entry(nil, "1.23"),
					Entry(nil, "1.24"),
					Entry(nil, "1.25"),
					Entry(nil, "1.26"),
					Entry(nil, "1.27"),
					Entry(nil, "1.28"),
					Entry(nil, "1.29"),
				)

				Context("for arm instance type", func() {
					BeforeEach(func() {
						instanceType = "c6g.12xlarge"
					})
					DescribeTable("should return a valid AMI for arm64",
						func(version string) {
							addMockGetParameter(p, fmt.Sprintf("/aws/service/canonical/ubuntu/eks/20.04/%s/stable/current/arm64/hvm/ebs-gp2/ami-id", version), expectedAmi)

							resolver := NewSSMResolver(p.MockSSM())
							resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)

							Expect(err).NotTo(HaveOccurred())
							Expect(p.MockSSM().AssertNumberOfCalls(GinkgoT(), "GetParameter", 1)).To(BeTrue())
							Expect(resolvedAmi).To(BeEquivalentTo(expectedAmi))
						},
						EntryDescription("When EKS version is %s"),
						Entry(nil, "1.21"),
						Entry(nil, "1.22"),
						Entry(nil, "1.23"),
						Entry(nil, "1.24"),
						Entry(nil, "1.25"),
						Entry(nil, "1.26"),
						Entry(nil, "1.27"),
						Entry(nil, "1.28"),
						Entry(nil, "1.29"),
					)
				})
			})

			Context("and UbuntuPro2004 family", func() {
				BeforeEach(func() {
					p = mockprovider.NewMockProvider()
					instanceType = "t2.medium"
					imageFamily = "UbuntuPro2004"
				})

				DescribeTable("should return an error",
					func(version string) {
						resolver := NewSSMResolver(p.MockSSM())
						resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)

						Expect(err).To(HaveOccurred())
						Expect(err).To(MatchError("UbuntuPro2004 requires EKS version greater or equal than 1.27 and lower than 1.29"))
					},
					EntryDescription("When EKS version is %s"),
					Entry(nil, "1.26"),
					Entry(nil, "1.30"),
				)

				DescribeTable("should return a valid AMI",
					func(version string) {
						addMockGetParameter(p, fmt.Sprintf("/aws/service/canonical/ubuntu/eks-pro/20.04/%s/stable/current/amd64/hvm/ebs-gp2/ami-id", version), expectedAmi)

						resolver := NewSSMResolver(p.MockSSM())
						resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)

						Expect(err).NotTo(HaveOccurred())
						Expect(p.MockSSM().AssertNumberOfCalls(GinkgoT(), "GetParameter", 1)).To(BeTrue())
						Expect(resolvedAmi).To(BeEquivalentTo(expectedAmi))
					},
					EntryDescription("When EKS version is %s"),
					Entry(nil, "1.27"),
					Entry(nil, "1.28"),
					Entry(nil, "1.29"),
				)

				Context("for arm instance type", func() {
					BeforeEach(func() {
						instanceType = "c6g.12xlarge"
					})
					DescribeTable("should return a valid AMI for arm64",
						func(version string) {
							addMockGetParameter(p, fmt.Sprintf("/aws/service/canonical/ubuntu/eks-pro/20.04/%s/stable/current/arm64/hvm/ebs-gp2/ami-id", version), expectedAmi)

							resolver := NewSSMResolver(p.MockSSM())
							resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)

							Expect(err).NotTo(HaveOccurred())
							Expect(p.MockSSM().AssertNumberOfCalls(GinkgoT(), "GetParameter", 1)).To(BeTrue())
							Expect(resolvedAmi).To(BeEquivalentTo(expectedAmi))
						},
						EntryDescription("When EKS version is %s"),
						Entry(nil, "1.27"),
						Entry(nil, "1.28"),
						Entry(nil, "1.29"),
					)
				})
			})

			Context("and Ubuntu2204 family", func() {
				BeforeEach(func() {
					p = mockprovider.NewMockProvider()
					instanceType = "t2.medium"
					imageFamily = "Ubuntu2204"
				})

				DescribeTable("should return an error",
					func(version string) {
						resolver := NewSSMResolver(p.MockSSM())
						resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)

						Expect(err).To(HaveOccurred())
						Expect(err).To(MatchError("Ubuntu2204 requires EKS version greater or equal than 1.29"))
					},
					EntryDescription("When EKS version is %s"),
					Entry(nil, "1.21"),
					Entry(nil, "1.28"),
				)

				DescribeTable("should return a valid AMI",
					func(version string) {
						addMockGetParameter(p, fmt.Sprintf("/aws/service/canonical/ubuntu/eks/22.04/%s/stable/current/amd64/hvm/ebs-gp2/ami-id", version), expectedAmi)

						resolver := NewSSMResolver(p.MockSSM())
						resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)

						Expect(err).NotTo(HaveOccurred())
						Expect(p.MockSSM().AssertNumberOfCalls(GinkgoT(), "GetParameter", 1)).To(BeTrue())
						Expect(resolvedAmi).To(BeEquivalentTo(expectedAmi))
					},
					EntryDescription("When EKS version is %s"),
					Entry(nil, "1.29"),
					Entry(nil, "1.30"),
					Entry(nil, "1.31"),
				)

				Context("for arm instance type", func() {
					BeforeEach(func() {
						instanceType = "g5g.2xlarge"
					})
					DescribeTable("should return a valid AMI for arm64",
						func(version string) {
							addMockGetParameter(p, fmt.Sprintf("/aws/service/canonical/ubuntu/eks/22.04/%s/stable/current/arm64/hvm/ebs-gp2/ami-id", version), expectedAmi)

							resolver := NewSSMResolver(p.MockSSM())
							resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)

							Expect(err).NotTo(HaveOccurred())
							Expect(p.MockSSM().AssertNumberOfCalls(GinkgoT(), "GetParameter", 1)).To(BeTrue())
							Expect(resolvedAmi).To(BeEquivalentTo(expectedAmi))
						},
						EntryDescription("When EKS version is %s"),
						Entry(nil, "1.29"),
						Entry(nil, "1.30"),
						Entry(nil, "1.31"),
					)
				})
			})

			Context("and UbuntuPro2204 family", func() {
				BeforeEach(func() {
					p = mockprovider.NewMockProvider()
					instanceType = "t2.medium"
					imageFamily = "UbuntuPro2204"
				})

				DescribeTable("should return an error",
					func(version string) {
						resolver := NewSSMResolver(p.MockSSM())
						resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)

						Expect(err).To(HaveOccurred())
						Expect(err).To(MatchError("UbuntuPro2204 requires EKS version greater or equal than 1.29"))
					},
					EntryDescription("When EKS version is %s"),
					Entry(nil, "1.21"),
					Entry(nil, "1.28"),
				)

				DescribeTable("should return a valid AMI",
					func(version string) {
						addMockGetParameter(p, fmt.Sprintf("/aws/service/canonical/ubuntu/eks-pro/22.04/%s/stable/current/amd64/hvm/ebs-gp2/ami-id", version), expectedAmi)

						resolver := NewSSMResolver(p.MockSSM())
						resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)

						Expect(err).NotTo(HaveOccurred())
						Expect(p.MockSSM().AssertNumberOfCalls(GinkgoT(), "GetParameter", 1)).To(BeTrue())
						Expect(resolvedAmi).To(BeEquivalentTo(expectedAmi))
					},
					EntryDescription("When EKS version is %s"),
					Entry(nil, "1.29"),
					Entry(nil, "1.30"),
					Entry(nil, "1.31"),
				)

				Context("for arm instance type", func() {
					BeforeEach(func() {
						instanceType = "g5g.2xlarge"
					})
					DescribeTable("should return a valid AMI for arm64",
						func(version string) {
							addMockGetParameter(p, fmt.Sprintf("/aws/service/canonical/ubuntu/eks-pro/22.04/%s/stable/current/arm64/hvm/ebs-gp2/ami-id", version), expectedAmi)

							resolver := NewSSMResolver(p.MockSSM())
							resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)

							Expect(err).NotTo(HaveOccurred())
							Expect(p.MockSSM().AssertNumberOfCalls(GinkgoT(), "GetParameter", 1)).To(BeTrue())
							Expect(resolvedAmi).To(BeEquivalentTo(expectedAmi))
						},
						EntryDescription("When EKS version is %s"),
						Entry(nil, "1.29"),
						Entry(nil, "1.30"),
						Entry(nil, "1.31"),
					)
				})
			})

			Context("and Ubuntu2404 family", func() {
				BeforeEach(func() {
					p = mockprovider.NewMockProvider()
					instanceType = "t2.medium"
					imageFamily = "Ubuntu2404"
				})

				DescribeTable("should return an error",
					func(version string) {
						resolver := NewSSMResolver(p.MockSSM())
						resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)

						Expect(err).To(HaveOccurred())
						Expect(err).To(MatchError("Ubuntu2404 requires EKS version greater or equal than 1.31"))
					},
					EntryDescription("When EKS version is %s"),
					Entry(nil, "1.21"),
					Entry(nil, "1.30"),
				)

				DescribeTable("should return a valid AMI",
					func(version string) {
						addMockGetParameter(p, fmt.Sprintf("/aws/service/canonical/ubuntu/eks/24.04/%s/stable/current/amd64/hvm/ebs-gp3/ami-id", version), expectedAmi)

						resolver := NewSSMResolver(p.MockSSM())
						resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)

						Expect(err).NotTo(HaveOccurred())
						Expect(p.MockSSM().AssertNumberOfCalls(GinkgoT(), "GetParameter", 1)).To(BeTrue())
						Expect(resolvedAmi).To(BeEquivalentTo(expectedAmi))
					},
					EntryDescription("When EKS version is %s"),
					Entry(nil, "1.31"),
				)

				Context("for arm instance type", func() {
					BeforeEach(func() {
						instanceType = "g5g.2xlarge"
					})
					DescribeTable("should return a valid AMI for arm64",
						func(version string) {
							addMockGetParameter(p, fmt.Sprintf("/aws/service/canonical/ubuntu/eks/24.04/%s/stable/current/arm64/hvm/ebs-gp3/ami-id", version), expectedAmi)

							resolver := NewSSMResolver(p.MockSSM())
							resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)

							Expect(err).NotTo(HaveOccurred())
							Expect(p.MockSSM().AssertNumberOfCalls(GinkgoT(), "GetParameter", 1)).To(BeTrue())
							Expect(resolvedAmi).To(BeEquivalentTo(expectedAmi))
						},
						EntryDescription("When EKS version is %s"),
						Entry(nil, "1.31"),
					)
				})
			})

			Context("and UbuntuPro2404 family", func() {
				BeforeEach(func() {
					p = mockprovider.NewMockProvider()
					instanceType = "t2.medium"
					imageFamily = "UbuntuPro2404"
				})

				DescribeTable("should return an error",
					func(version string) {
						resolver := NewSSMResolver(p.MockSSM())
						resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)

						Expect(err).To(HaveOccurred())
						Expect(err).To(MatchError("UbuntuPro2404 requires EKS version greater or equal than 1.31"))
					},
					EntryDescription("When EKS version is %s"),
					Entry(nil, "1.21"),
					Entry(nil, "1.30"),
				)

				DescribeTable("should return a valid AMI",
					func(version string) {
						addMockGetParameter(p, fmt.Sprintf("/aws/service/canonical/ubuntu/eks-pro/24.04/%s/stable/current/amd64/hvm/ebs-gp3/ami-id", version), expectedAmi)

						resolver := NewSSMResolver(p.MockSSM())
						resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)

						Expect(err).NotTo(HaveOccurred())
						Expect(p.MockSSM().AssertNumberOfCalls(GinkgoT(), "GetParameter", 1)).To(BeTrue())
						Expect(resolvedAmi).To(BeEquivalentTo(expectedAmi))
					},
					EntryDescription("When EKS version is %s"),
					Entry(nil, "1.31"),
				)

				Context("for arm instance type", func() {
					BeforeEach(func() {
						instanceType = "g5g.2xlarge"
					})
					DescribeTable("should return a valid AMI for arm64",
						func(version string) {
							addMockGetParameter(p, fmt.Sprintf("/aws/service/canonical/ubuntu/eks-pro/24.04/%s/stable/current/arm64/hvm/ebs-gp3/ami-id", version), expectedAmi)

							resolver := NewSSMResolver(p.MockSSM())
							resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)

							Expect(err).NotTo(HaveOccurred())
							Expect(p.MockSSM().AssertNumberOfCalls(GinkgoT(), "GetParameter", 1)).To(BeTrue())
							Expect(resolvedAmi).To(BeEquivalentTo(expectedAmi))
						},
						EntryDescription("When EKS version is %s"),
						Entry(nil, "1.31"),
					)
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
						p = mockprovider.NewMockProvider()
						addMockGetParameter(p, "/aws/service/bottlerocket/aws-k8s-1.15/x86_64/latest/image_id", expectedAmi)
						resolver := NewSSMResolver(p.MockSSM())
						resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)
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
						p = mockprovider.NewMockProvider()
						addMockFailedGetParameter(p, "/aws/service/bottlerocket/aws-k8s-1.15/x86_64/latest/image_id")

						resolver := NewSSMResolver(p.MockSSM())
						resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)
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
						instanceType = "c6g.2xlarge"
					})

					Context("and ami is available", func() {
						BeforeEach(func() {
							p = mockprovider.NewMockProvider()
							addMockGetParameter(p, "/aws/service/bottlerocket/aws-k8s-1.15/arm64/latest/image_id", expectedAmi)
							resolver := NewSSMResolver(p.MockSSM())
							resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)
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
							p = mockprovider.NewMockProvider()
							addMockFailedGetParameter(p, "/aws/service/bottlerocket/aws-k8s-1.15/arm64/latest/image_id")

							resolver := NewSSMResolver(p.MockSSM())
							resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)
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

				Context("and gpu instance", func() {
					BeforeEach(func() {
						instanceType = "p3.2xlarge"
						version = "1.23"
					})

					Context("and ami is available", func() {
						BeforeEach(func() {
							p = mockprovider.NewMockProvider()
							addMockGetParameter(p, "/aws/service/bottlerocket/aws-k8s-1.23-nvidia/x86_64/latest/image_id", expectedAmi)
							resolver := NewSSMResolver(p.MockSSM())
							resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)
						})

						It("does not return an error", func() {
							Expect(err).NotTo(HaveOccurred())
						})
						It("calls AWS SSM GetParameter", func() {
							Expect(p.MockSSM().AssertNumberOfCalls(GinkgoT(), "GetParameter", 1)).To(BeTrue())
						})
						It("returns an ami id", func() {
							Expect(resolvedAmi).To(BeEquivalentTo(expectedAmi))
						})
					})

					Context("and ami is NOT available", func() {
						BeforeEach(func() {
							p = mockprovider.NewMockProvider()
							addMockFailedGetParameter(p, "/aws/service/bottlerocket/aws-k8s-1.23-nvidia/x86_64/latest/image_id")

							resolver := NewSSMResolver(p.MockSSM())
							resolvedAmi, err = resolver.Resolve(context.Background(), region, version, instanceType, imageFamily)
						})

						It("errors", func() {
							Expect(err).To(HaveOccurred())
						})

						It("does NOT return an ami id", func() {
							Expect(resolvedAmi).To(BeEquivalentTo(""))
						})

						It("calls AWS SSM GetParameter", func() {
							Expect(p.MockSSM().AssertNumberOfCalls(GinkgoT(), "GetParameter", 1)).To(BeTrue())
						})

					})
				})
			})
		})
	})

	Context("managed SSM parameter name", func() {
		It("should support SSM parameter generation for all AMI types but Windows", func() {
			var eksAMIType ekstypes.AMITypes
			for _, amiType := range eksAMIType.Values() {
				if amiType == ekstypes.AMITypesCustom || strings.HasPrefix(string(amiType), "WINDOWS_") ||
					// TODO: remove this condition after support for Bottlerocket FIPS AMI types.
					amiType == ekstypes.AMITypesBottlerocketArm64Fips || amiType == ekstypes.AMITypesBottlerocketX8664Fips {
					continue
				}
				ssmParameterName := MakeManagedSSMParameterName(api.Version1_31, amiType)
				Expect(ssmParameterName).NotTo(BeEmpty(), "expected to generate SSM parameter name for AMI type %s", amiType)
			}
		})
	})
})

func addMockGetParameter(p *mockprovider.MockProvider, name, amiID string) {
	p.MockSSM().On("GetParameter", mock.Anything,
		mock.MatchedBy(func(input *ssm.GetParameterInput) bool {
			return *input.Name == name
		}),
	).Return(&ssm.GetParameterOutput{
		Parameter: &ssmtypes.Parameter{
			Name:  aws.String(name),
			Type:  ssmtypes.ParameterTypeString,
			Value: aws.String(amiID),
		},
	}, nil)
}

func addMockFailedGetParameter(p *mockprovider.MockProvider, name string) {
	p.MockSSM().On("GetParameter", mock.Anything,
		mock.MatchedBy(func(input *ssm.GetParameterInput) bool {
			return *input.Name == name
		}),
	).Return(&ssm.GetParameterOutput{
		Parameter: nil,
	}, nil)
}
