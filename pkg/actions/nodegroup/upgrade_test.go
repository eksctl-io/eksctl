package nodegroup_test

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
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
			KubernetesVersion: "1.21",
			Wait:              false,
			ForceUpgrade:      false,
		}
	})

	When("the nodegroup does not have a stack", func() {
		When("launchTemplate Id is set", func() {
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
						Version:       aws.String("1.20"),
						LaunchTemplate: &ekstypes.LaunchTemplateSpecification{
							Id: aws.String("id-123"),
						},
					},
				}, nil)

				p.MockEKS().On("UpdateNodegroupVersion", mock.Anything, &awseks.UpdateNodegroupVersionInput{
					NodegroupName: aws.String(ngName),
					ClusterName:   aws.String(clusterName),
					Force:         false,
					Version:       aws.String("1.21"),
					LaunchTemplate: &ekstypes.LaunchTemplateSpecification{
						Id:      aws.String("id-123"),
						Version: aws.String("v2"),
					},
				}).Return(&awseks.UpdateNodegroupVersionOutput{}, nil)
			})

			It("upgrades the nodegroup version and lt by calling the API", func() {
				options.LaunchTemplateVersion = "v2"
				Expect(m.Upgrade(context.Background(), options)).To(Succeed())
			})
		})

		When("launchTemplate name is set", func() {
			BeforeEach(func() {
				p.MockEKS().On("DescribeNodegroup", mock.Anything, &awseks.DescribeNodegroupInput{
					ClusterName:   aws.String(clusterName),
					NodegroupName: aws.String(ngName),
				}).Return(&awseks.DescribeNodegroupOutput{
					Nodegroup: &ekstypes.Nodegroup{
						NodegroupName: aws.String(ngName),
						ClusterName:   aws.String(clusterName),
						Status:        ekstypes.NodegroupStatusActive,
						AmiType:       ekstypes.AMITypesAl2X8664,
						Version:       aws.String("1.20"),
						LaunchTemplate: &ekstypes.LaunchTemplateSpecification{
							Name: aws.String("lt"),
						},
					},
				}, nil)

				p.MockEKS().On("UpdateNodegroupVersion", mock.Anything, &awseks.UpdateNodegroupVersionInput{
					NodegroupName: aws.String(ngName),
					ClusterName:   aws.String(clusterName),
					Force:         false,
					Version:       aws.String("1.21"),
					LaunchTemplate: &ekstypes.LaunchTemplateSpecification{
						Name:    aws.String("lt"),
						Version: aws.String("v2"),
					},
				}).Return(&awseks.UpdateNodegroupVersionOutput{}, nil)
			})

			It("upgrades the nodegroup version and lt by calling the API", func() {
				options.LaunchTemplateVersion = "v2"
				Expect(m.Upgrade(context.Background(), options)).To(Succeed())
			})
		})
	})

	When("the nodegroup does have a stack", func() {
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
							Version:        aws.String("1.20"),
							ReleaseVersion: aws.String("1.20-20201212"),
						},
					}, nil)

					p.MockSSM().On("GetParameter", mock.Anything, &ssm.GetParameterInput{
						Name: aws.String("/aws/service/eks/optimized-ami/1.21/amazon-linux-2/recommended/release_version"),
					}).Return(&ssm.GetParameterOutput{
						Parameter: &ssmtypes.Parameter{
							Value: aws.String("1.21-20201212"),
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
							Version:        aws.String("1.20"),
							ReleaseVersion: aws.String("1.20-20201212"),
						},
					}, nil)

					p.MockSSM().On("GetParameter", mock.Anything, &ssm.GetParameterInput{
						Name: aws.String("/aws/service/eks/optimized-ami/1.21/amazon-linux-2-gpu/recommended/release_version"),
					}).Return(&ssm.GetParameterOutput{
						Parameter: &ssmtypes.Parameter{
							Value: aws.String("1.21-20201212"),
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
							Version:        aws.String("1.20"),
							ReleaseVersion: aws.String("1.20-20201212"),
						},
					}, nil)

					p.MockSSM().On("GetParameter", mock.Anything, &ssm.GetParameterInput{
						Name: aws.String("/aws/service/bottlerocket/aws-k8s-1.21/x86_64/latest/image_version"),
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
						Version:        aws.String("1.20"),
						ReleaseVersion: aws.String("1.20-20201212"),
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
						Version:        aws.String("1.20"),
						ReleaseVersion: aws.String("1.20-20201212"),
					},
				}, nil)
			})

			It("should return an error", func() {
				Expect(m.Upgrade(context.Background(), options)).To(MatchError(ContainSubstring(`nodegroup must be in "ACTIVE" state when upgrading a nodegroup; got state "DEGRADED"`)))
			})
		})
	})
})
