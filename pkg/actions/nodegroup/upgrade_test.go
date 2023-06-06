package nodegroup_test

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/stretchr/testify/mock"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/cfn/manager/fakes"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
	"github.com/weaveworks/eksctl/pkg/version"
)

var _ = Describe("Upgrade", func() {
	supportedVersions := api.SupportedVersions()
	if len(supportedVersions) < 2 {
		Fail("Upgrade test requires at least two supported EKS versions")
	}

	var (
		eksVersion        = aws.String(supportedVersions[len(supportedVersions)-2])
		latestEKSVersion  = supportedVersions[len(supportedVersions)-1]
		eksReleaseVersion = aws.String(fmt.Sprintf("%s-20201212", *eksVersion))
	)

	var (
		clusterName, ngName string
		p                   *mockprovider.MockProvider
		cfg                 *api.ClusterConfig
		m                   *nodegroup.Manager
		fakeStackManager    *fakes.FakeStackManager
		fakeClientSet       *fake.Clientset
		options             nodegroup.UpgradeOptions
	)

	BeforeEach(func() {
		ngName = "my-nodegroup"
		clusterName = "my-cluster"
		cfg = api.NewClusterConfig()
		cfg.Metadata.Name = clusterName
		p = mockprovider.NewMockProvider()
		fakeClientSet = fake.NewSimpleClientset()
		m = nodegroup.New(cfg, &eks.ClusterProvider{AWSProvider: p}, fakeClientSet, nil)

		fakeStackManager = new(fakes.FakeStackManager)
		m.SetStackManager(fakeStackManager)
		options = nodegroup.UpgradeOptions{
			NodegroupName:     ngName,
			KubernetesVersion: latestEKSVersion,
			Wait:              false,
			ForceUpgrade:      false,
		}
	})

	Context("the nodegroup does not have a stack", func() {
		When("the nodegroup does not use a launch template", func() {
			BeforeEach(func() {
				p.MockEKS().On("DescribeNodegroup", mock.Anything, &awseks.DescribeNodegroupInput{
					ClusterName:   aws.String(clusterName),
					NodegroupName: aws.String(ngName),
				}).Return(&awseks.DescribeNodegroupOutput{
					Nodegroup: &ekstypes.Nodegroup{
						NodegroupName:  aws.String(ngName),
						ClusterName:    aws.String(clusterName),
						Status:         ekstypes.NodegroupStatusActive,
						AmiType:        "ami-type",
						Version:        eksVersion,
						LaunchTemplate: nil,
					},
				}, nil)
			})
			It("returns an error if launch template version is specified", func() {
				options.LaunchTemplateVersion = "2"
				err := m.Upgrade(context.Background(), options)
				Expect(err).To(MatchError(ContainSubstring("cannot update launch template version because the nodegroup is not configured to use one")))
			})
		})

		When("the nodegroup uses launch template", func() {
			BeforeEach(func() {
				p.MockEKS().On("DescribeNodegroup", mock.Anything, &awseks.DescribeNodegroupInput{
					ClusterName:   aws.String(clusterName),
					NodegroupName: aws.String(ngName),
				}).Return(&awseks.DescribeNodegroupOutput{
					Nodegroup: &ekstypes.Nodegroup{
						NodegroupName: aws.String(ngName),
						ClusterName:   aws.String(clusterName),
						Status:        ekstypes.NodegroupStatusActive,
						AmiType:       "ami-type",
						Version:       eksVersion,
						LaunchTemplate: &ekstypes.LaunchTemplateSpecification{
							Id:      aws.String("id-123"),
							Version: aws.String("2"),
						},
					},
				}, nil)
			})
			When("using a custom AMI", func() {
				BeforeEach(func() {
					options.ReleaseVersion = ""
					options.KubernetesVersion = ""
					p.MockEC2().On("DescribeLaunchTemplateVersions", mock.Anything, &ec2.DescribeLaunchTemplateVersionsInput{
						LaunchTemplateId: aws.String("id-123"),
						Versions:         []string{"2"},
					}).Return(&ec2.DescribeLaunchTemplateVersionsOutput{LaunchTemplateVersions: []ec2types.LaunchTemplateVersion{
						{
							LaunchTemplateData: &ec2types.ResponseLaunchTemplateData{
								ImageId:      aws.String("1234"),
								InstanceType: "big",
							},
							VersionNumber: aws.Int64(2),
						},
					}}, nil)
				})
				It("returns an error if kubernetes version is specified", func() {
					options.KubernetesVersion = api.DefaultVersion
					err := m.Upgrade(context.Background(), options)
					Expect(err).To(MatchError(ContainSubstring("cannot specify kubernetes-version or release-version when using a custom AMI")))
				})
				It("returns an error if release version is specified", func() {
					options.ReleaseVersion = *eksReleaseVersion
					err := m.Upgrade(context.Background(), options)
					Expect(err).To(MatchError(ContainSubstring("cannot specify kubernetes-version or release-version when using a custom AMI")))
				})
				It("upgrades the nodegroup version and lt by calling the EKS API", func() {
					p.MockEKS().On("UpdateNodegroupVersion", mock.Anything, &awseks.UpdateNodegroupVersionInput{
						NodegroupName: aws.String(ngName),
						ClusterName:   aws.String(clusterName),
						Force:         false,
						LaunchTemplate: &ekstypes.LaunchTemplateSpecification{
							Id:      aws.String("id-123"),
							Version: aws.String("2"),
						},
					}).Return(&awseks.UpdateNodegroupVersionOutput{}, nil)
					options.LaunchTemplateVersion = "2"
					Expect(m.Upgrade(context.Background(), options)).To(Succeed())
				})
			})
			When("not using a custom AMI", func() {
				BeforeEach(func() {
					p.MockEC2().On("DescribeLaunchTemplateVersions", mock.Anything, &ec2.DescribeLaunchTemplateVersionsInput{
						LaunchTemplateId: aws.String("id-123"),
						Versions:         []string{"2"},
					}).Return(&ec2.DescribeLaunchTemplateVersionsOutput{LaunchTemplateVersions: []ec2types.LaunchTemplateVersion{
						{
							LaunchTemplateData: &ec2types.ResponseLaunchTemplateData{
								ImageId:      nil,
								InstanceType: "big",
							},
							VersionNumber: aws.Int64(2),
						},
					}}, nil)
				})
				It("upgrades the nodegroup version and lt by calling the EKS API", func() {
					p.MockEKS().On("UpdateNodegroupVersion", mock.Anything, &awseks.UpdateNodegroupVersionInput{
						NodegroupName: aws.String(ngName),
						ClusterName:   aws.String(clusterName),
						Force:         false,
						Version:       aws.String(api.Version1_22),
						LaunchTemplate: &ekstypes.LaunchTemplateSpecification{
							Id:      aws.String("id-123"),
							Version: aws.String("2"),
						},
						ReleaseVersion: eksReleaseVersion,
					}).Return(&awseks.UpdateNodegroupVersionOutput{}, nil)
					options.KubernetesVersion = api.Version1_22
					options.LaunchTemplateVersion = "2"
					options.ReleaseVersion = *eksReleaseVersion
					Expect(m.Upgrade(context.Background(), options)).To(Succeed())
				})
			})
		})
	})

	Context("the nodegroup does have a stack", func() {
		When("ForceUpdateEnabled isn't set", func() {
			When("it uses amazonlinux2", func() {
				BeforeEach(func() {
					fakeStackManager.ListNodeGroupStacksWithStatusesReturns([]manager.NodeGroupStack{{NodeGroupName: ngName}}, nil)

					fakeStackManager.GetManagedNodeGroupTemplateReturns(al2WithoutForceTemplate, nil)

					fakeStackManager.DescribeNodeGroupStackReturns(&manager.Stack{
						Tags: []types.Tag{
							{
								Key:   aws.String(api.EksctlVersionTag),
								Value: aws.String(version.GetVersion()),
							},
						},
					}, nil)

					fakeStackManager.UpdateNodeGroupStackReturns(nil)

					p.MockEKS().On("DescribeNodegroup", mock.Anything, &awseks.DescribeNodegroupInput{
						ClusterName:   aws.String(clusterName),
						NodegroupName: aws.String(ngName),
					}).Return(&awseks.DescribeNodegroupOutput{
						Nodegroup: &ekstypes.Nodegroup{
							NodegroupName:  aws.String(ngName),
							ClusterName:    aws.String(clusterName),
							Status:         ekstypes.NodegroupStatusActive,
							AmiType:        ekstypes.AMITypesAl2X8664,
							Version:        eksVersion,
							ReleaseVersion: eksReleaseVersion,
						},
					}, nil)

					p.MockSSM().On("GetParameter", mock.Anything, &ssm.GetParameterInput{
						Name: aws.String(fmt.Sprintf("/aws/service/eks/optimized-ami/%s/amazon-linux-2/recommended/release_version", latestEKSVersion)),
					}).Return(&ssm.GetParameterOutput{
						Parameter: &ssmtypes.Parameter{
							Value: aws.String(fmt.Sprintf("%s-20201212", latestEKSVersion)),
						},
					}, nil)
				})

				It("upgrades the nodegroup with the latest al2 release_version by updating the stack", func() {
					Expect(m.Upgrade(context.Background(), options)).To(Succeed())
					Expect(fakeStackManager.GetManagedNodeGroupTemplateCallCount()).To(Equal(1))
					_, n := fakeStackManager.GetManagedNodeGroupTemplateArgsForCall(0)
					Expect(n.NodeGroupName).To(Equal(ngName))
					Expect(fakeStackManager.UpdateNodeGroupStackCallCount()).To(Equal(2))
					By("upgrading the ForceUpdateEnabled setting first")
					_, ng, template, wait := fakeStackManager.UpdateNodeGroupStackArgsForCall(0)
					Expect(ng).To(Equal(ngName))
					Expect(template).To(Equal(al2ForceFalseTemplate))
					Expect(wait).To(BeTrue())

					By("upgrading the ReleaseVersion setting next")
					_, ng, template, wait = fakeStackManager.UpdateNodeGroupStackArgsForCall(1)
					Expect(ng).To(Equal(ngName))
					Expect(template).To(Equal(al2FullyUpdatedTemplate))
					Expect(wait).To(BeTrue())
				})
			})
		})

		When("it already has ForceUpdateEnabled set to false", func() {
			When("it uses amazonlinux2 GPU nodes", func() {
				BeforeEach(func() {
					fakeStackManager.ListNodeGroupStacksWithStatusesReturns([]manager.NodeGroupStack{{NodeGroupName: ngName}}, nil)

					fakeStackManager.GetManagedNodeGroupTemplateReturns(al2ForceFalseTemplate, nil)

					fakeStackManager.DescribeNodeGroupStackReturns(&manager.Stack{
						Tags: []types.Tag{
							{
								Key:   aws.String(api.EksctlVersionTag),
								Value: aws.String(version.GetVersion()),
							},
						},
					}, nil)

					fakeStackManager.UpdateNodeGroupStackReturns(nil)

					p.MockEKS().On("DescribeNodegroup", mock.Anything, &awseks.DescribeNodegroupInput{
						ClusterName:   aws.String(clusterName),
						NodegroupName: aws.String(ngName),
					}).Return(&awseks.DescribeNodegroupOutput{
						Nodegroup: &ekstypes.Nodegroup{
							NodegroupName:  aws.String(ngName),
							ClusterName:    aws.String(clusterName),
							Status:         ekstypes.NodegroupStatusActive,
							AmiType:        ekstypes.AMITypesAl2X8664Gpu,
							Version:        eksVersion,
							ReleaseVersion: eksReleaseVersion,
						},
					}, nil)

					p.MockSSM().On("GetParameter", mock.Anything, &ssm.GetParameterInput{
						Name: aws.String(fmt.Sprintf("/aws/service/eks/optimized-ami/%s/amazon-linux-2-gpu/recommended/release_version", latestEKSVersion)),
					}).Return(&ssm.GetParameterOutput{
						Parameter: &ssmtypes.Parameter{
							Value: aws.String(fmt.Sprintf("%s-20201212", latestEKSVersion)),
						},
					}, nil)
				})

				It("upgrades the nodegroup with the latest al2 release_version by updating the stack", func() {
					Expect(m.Upgrade(context.Background(), options)).To(Succeed())
					Expect(fakeStackManager.GetManagedNodeGroupTemplateCallCount()).To(Equal(1))
					_, n := fakeStackManager.GetManagedNodeGroupTemplateArgsForCall(0)
					Expect(n.NodeGroupName).To(Equal(ngName))
					Expect(fakeStackManager.UpdateNodeGroupStackCallCount()).To(Equal(1))
					By("upgrading the ReleaseVersion and not updating the ForceUpdateEnabled setting")
					_, ng, template, wait := fakeStackManager.UpdateNodeGroupStackArgsForCall(0)
					Expect(ng).To(Equal(ngName))
					Expect(template).To(Equal(al2FullyUpdatedTemplate))
					Expect(wait).To(BeTrue())
				})
			})
		})

		When("ForceUpdateEnabled is set to true but the desired value is false", func() {
			When("it uses bottlerocket", func() {
				BeforeEach(func() {
					fakeStackManager.ListNodeGroupStacksWithStatusesReturns([]manager.NodeGroupStack{{NodeGroupName: ngName}}, nil)

					fakeStackManager.GetManagedNodeGroupTemplateReturns(brForceTrueTemplate, nil)

					fakeStackManager.DescribeNodeGroupStackReturns(&manager.Stack{
						Tags: []types.Tag{
							{
								Key:   aws.String(api.EksctlVersionTag),
								Value: aws.String(version.GetVersion()),
							},
						},
					}, nil)

					fakeStackManager.UpdateNodeGroupStackReturns(nil)

					p.MockEKS().On("DescribeNodegroup", mock.Anything, &awseks.DescribeNodegroupInput{
						ClusterName:   aws.String(clusterName),
						NodegroupName: aws.String(ngName),
					}).Return(&awseks.DescribeNodegroupOutput{
						Nodegroup: &ekstypes.Nodegroup{
							NodegroupName:  aws.String(ngName),
							ClusterName:    aws.String(clusterName),
							Status:         ekstypes.NodegroupStatusActive,
							AmiType:        ekstypes.AMITypesBottlerocketX8664,
							Version:        eksVersion,
							ReleaseVersion: eksReleaseVersion,
						},
					}, nil)

					p.MockSSM().On("GetParameter", mock.Anything, &ssm.GetParameterInput{
						Name: aws.String(fmt.Sprintf("/aws/service/bottlerocket/aws-k8s-%s/x86_64/latest/image_version", latestEKSVersion)),
					}).Return(&ssm.GetParameterOutput{
						Parameter: &ssmtypes.Parameter{
							Value: aws.String("1.5.2-1602f3a8"),
						},
					}, nil)
				})

				It("upgrades the nodegroup updating the stack with the kubernetes version", func() {
					Expect(m.Upgrade(context.Background(), options)).To(Succeed())
					Expect(fakeStackManager.GetManagedNodeGroupTemplateCallCount()).To(Equal(1))
					_, n := fakeStackManager.GetManagedNodeGroupTemplateArgsForCall(0)
					Expect(n.NodeGroupName).To(Equal(ngName))

					By("upgrading the ForceUpdateEnabled setting first")
					Expect(fakeStackManager.UpdateNodeGroupStackCallCount()).To(Equal(2))
					_, ng, template, wait := fakeStackManager.UpdateNodeGroupStackArgsForCall(0)
					Expect(ng).To(Equal(ngName))
					Expect(template).To(Equal(brForceFalseTemplate))
					Expect(wait).To(BeTrue())

					By("upgrading the Version next")
					_, ng, template, wait = fakeStackManager.UpdateNodeGroupStackArgsForCall(1)
					Expect(ng).To(Equal(ngName))
					Expect(template).To(Equal(brFulllyUpdatedTemplate))
					Expect(wait).To(BeTrue())
				})
			})
		})

		When("nodegroup is already being updated", func() {
			BeforeEach(func() {
				p.MockEKS().On("DescribeNodegroup", mock.Anything, &awseks.DescribeNodegroupInput{
					ClusterName:   aws.String(clusterName),
					NodegroupName: aws.String(ngName),
				}).Return(&awseks.DescribeNodegroupOutput{
					Nodegroup: &ekstypes.Nodegroup{
						NodegroupName:  aws.String(ngName),
						ClusterName:    aws.String(clusterName),
						Status:         ekstypes.NodegroupStatusUpdating,
						AmiType:        ekstypes.AMITypesBottlerocketX8664,
						Version:        eksVersion,
						ReleaseVersion: eksReleaseVersion,
					},
				}, nil)
			})

			It("should return an error", func() {
				Expect(m.Upgrade(context.Background(), options)).To(MatchError(ContainSubstring("nodegroup is currently being updated, please retry the command after the existing update is complete")))
			})
		})

		When("nodegroup is not ACTIVE", func() {
			BeforeEach(func() {
				p.MockEKS().On("DescribeNodegroup", mock.Anything, &awseks.DescribeNodegroupInput{
					ClusterName:   aws.String(clusterName),
					NodegroupName: aws.String(ngName),
				}).Return(&awseks.DescribeNodegroupOutput{
					Nodegroup: &ekstypes.Nodegroup{
						NodegroupName:  aws.String(ngName),
						ClusterName:    aws.String(clusterName),
						Status:         ekstypes.NodegroupStatusDegraded,
						AmiType:        ekstypes.AMITypesBottlerocketX8664,
						Version:        eksVersion,
						ReleaseVersion: eksReleaseVersion,
					},
				}, nil)
			})

			It("should return an error", func() {
				Expect(m.Upgrade(context.Background(), options)).To(MatchError(ContainSubstring(`nodegroup must be in "ACTIVE" state when upgrading a nodegroup; got state "DEGRADED"`)))
			})
		})
	})
})
