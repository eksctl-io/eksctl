package karpenter_test

import (
	"bytes"
	"errors"
	"fmt"
	"net"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	"github.com/kris-nova/logger"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	core "k8s.io/client-go/testing"

	karpenteractions "github.com/weaveworks/eksctl/pkg/actions/karpenter"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/manager/fakes"
	managerfakes "github.com/weaveworks/eksctl/pkg/cfn/manager/fakes"
	"github.com/weaveworks/eksctl/pkg/eks"
	karpenterfakes "github.com/weaveworks/eksctl/pkg/karpenter/fakes"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
	"github.com/weaveworks/eksctl/pkg/utils/ipnet"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

type fakeTask struct {
	err error
}

func (f *fakeTask) Do(errs chan error) error {
	close(errs)
	return f.err
}
func (f *fakeTask) Describe() string { return "I'm a fake task" }

var _ = Describe("Create", func() {
	Context("Create Karpenter Installation", func() {
		var (
			clusterName            string
			p                      *mockprovider.MockProvider
			cfg                    *api.ClusterConfig
			fakeStackManager       *managerfakes.FakeStackManager
			ctl                    *eks.ClusterProvider
			fakeKarpenterInstaller *karpenterfakes.FakeChartInstaller
			fakeClientSet          *fake.Clientset
		)

		BeforeEach(func() {
			clusterName = "my-cluster"
			p = mockprovider.NewMockProvider()
			cfg = api.NewClusterConfig()
			cfg.Metadata.Name = clusterName
			cfg.VPC = vpcConfig()
			cfg.AvailabilityZones = []string{"us-west-2a", "us-west-2b"}
			cfg.Status = &api.ClusterStatus{
				ARN: "arn:aws:iam::123456789012:user/test",
			}
			cfg.Karpenter = &api.Karpenter{
				Version: "0.4.3",
			}
			fakeStackManager = &fakes.FakeStackManager{}
			fakeKarpenterInstaller = &karpenterfakes.FakeChartInstaller{}
			ctl = &eks.ClusterProvider{
				Provider: p,
				Status: &eks.ProviderStatus{
					ClusterInfo: &eks.ClusterInfo{
						Cluster: testutils.NewFakeCluster(clusterName, awseks.ClusterStatusActive),
					},
				},
			}
			fakeStackManager.CreateStackStub = func(_ string, rs builder.ResourceSet, _ map[string]string, _ map[string]string, errs chan error) error {
				go func() {
					errs <- nil
				}()
				return nil
			}
			p.MockEC2().On("DescribeSubnets", &ec2.DescribeSubnetsInput{
				SubnetIds: aws.StringSlice([]string{
					privateSubnet1,
					privateSubnet2,
					publicSubnet1,
					publicSubnet2,
				}),
			}).Return(&ec2.DescribeSubnetsOutput{
				Subnets: []*ec2.Subnet{
					{
						SubnetId: aws.String(privateSubnet1),
						VpcId:    aws.String(cfg.VPC.ID),
					},
					{
						SubnetId: aws.String(privateSubnet2),
						VpcId:    aws.String(cfg.VPC.ID),
						Tags: []*ec2.Tag{
							{
								Key:   aws.String("kubernetes.io/cluster/" + clusterName),
								Value: aws.String("shared"),
							},
						},
					},
					{
						SubnetId: aws.String(publicSubnet1),
						VpcId:    aws.String(cfg.VPC.ID),
					},
					{
						SubnetId: aws.String(publicSubnet2),
						VpcId:    aws.String(cfg.VPC.ID),
					},
				},
			}, nil)

			p.MockEC2().On("CreateTags", &ec2.CreateTagsInput{
				Resources: []*string{
					&privateSubnet1,
					&publicSubnet1,
					&publicSubnet2,
				},
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("kubernetes.io/cluster/" + clusterName),
						Value: aws.String(""),
					},
				},
			}).Return(&ec2.CreateTagsOutput{}, nil)
			fakeClientSet = fake.NewSimpleClientset()
		})
		It("can install Karpenter on an existing cluster", func() {
			fakeKarpenterInstaller.InstallReturns(nil)
			install := &karpenteractions.Installer{
				StackManager:       fakeStackManager,
				CTL:                ctl,
				Config:             cfg,
				KarpenterInstaller: fakeKarpenterInstaller,
				ClientSet:          fakeClientSet,
			}
			Expect(install.Create()).To(Succeed())
			Expect(fakeKarpenterInstaller.InstallCallCount()).To(Equal(1))
		})
		When("DescribeSubnets fails", func() {
			var (
				output *bytes.Buffer
			)
			BeforeEach(func() {
				p = mockprovider.NewMockProvider()
				p.MockEC2().On("DescribeSubnets", mock.Anything).Return(nil, errors.New("nope"))
				ctl = &eks.ClusterProvider{
					Provider: p,
					Status: &eks.ProviderStatus{
						ClusterInfo: &eks.ClusterInfo{
							Cluster: testutils.NewFakeCluster(clusterName, ""),
						},
					},
				}
				output = &bytes.Buffer{}
				logger.Writer = output
			})
			It("errors", func() {
				install := &karpenteractions.Installer{
					StackManager:       fakeStackManager,
					CTL:                ctl,
					Config:             cfg,
					KarpenterInstaller: fakeKarpenterInstaller,
					ClientSet:          fakeClientSet,
				}
				err := install.Create()
				Expect(err).To(MatchError(ContainSubstring("failed to install Karpenter on cluster")))
				Expect(output.String()).To(ContainSubstring("failed to describe subnets: nope"))
			})
		})
		When("CreateTags fails", func() {
			var (
				output *bytes.Buffer
			)
			BeforeEach(func() {
				p = mockprovider.NewMockProvider()
				p.MockEC2().On("DescribeSubnets", mock.Anything).Return(&ec2.DescribeSubnetsOutput{
					Subnets: []*ec2.Subnet{
						{
							SubnetId: aws.String(privateSubnet1),
							VpcId:    aws.String(cfg.VPC.ID),
						},
					},
				}, nil)
				p.MockEC2().On("CreateTags", mock.Anything).Return(nil, errors.New("nope"))
				ctl = &eks.ClusterProvider{
					Provider: p,
					Status: &eks.ProviderStatus{
						ClusterInfo: &eks.ClusterInfo{
							Cluster: testutils.NewFakeCluster(clusterName, ""),
						},
					},
				}
				output = &bytes.Buffer{}
				logger.Writer = output
			})
			It("errors", func() {
				install := &karpenteractions.Installer{
					StackManager:       fakeStackManager,
					CTL:                ctl,
					Config:             cfg,
					KarpenterInstaller: fakeKarpenterInstaller,
					ClientSet:          fakeClientSet,
				}
				err := install.Create()
				Expect(err).To(MatchError(ContainSubstring("failed to install Karpenter on cluster")))
				Expect(output.String()).To(ContainSubstring("failed to add tags for subnets: nope"))
			})
		})
		When("Karpenter install fails", func() {
			It("errors", func() {
				fakeKarpenterInstaller.InstallReturns(errors.New("nope"))
				install := &karpenteractions.Installer{
					StackManager:       fakeStackManager,
					CTL:                ctl,
					Config:             cfg,
					KarpenterInstaller: fakeKarpenterInstaller,
					ClientSet:          fakeClientSet,
				}
				err := install.Create()
				Expect(err).To(MatchError(ContainSubstring("nope")))
			})
		})
		When("CreateStack fails", func() {
			var (
				output *bytes.Buffer
			)
			BeforeEach(func() {
				fakeStackManager.CreateStackStub = func(_ string, rs builder.ResourceSet, _ map[string]string, _ map[string]string, errs chan error) error {
					go func() {
						errs <- nil
					}()
					return errors.New("nope")
				}
				output = &bytes.Buffer{}
				logger.Writer = output
			})
			It("errors", func() {
				install := &karpenteractions.Installer{
					StackManager:       fakeStackManager,
					CTL:                ctl,
					Config:             cfg,
					KarpenterInstaller: fakeKarpenterInstaller,
					ClientSet:          fakeClientSet,
				}
				err := install.Create()
				Expect(err).To(MatchError(ContainSubstring("failed to install Karpenter on cluster")))
				Expect(output.String()).To(ContainSubstring("failed to create stack: nope"))
			})
		})
		When("arn is invalid in the status", func() {
			BeforeEach(func() {
				cfg.Status.ARN = ""
			})
			It("errors", func() {
				install := &karpenteractions.Installer{
					StackManager:       fakeStackManager,
					CTL:                ctl,
					Config:             cfg,
					KarpenterInstaller: fakeKarpenterInstaller,
					ClientSet:          fakeClientSet,
				}
				err := install.Create()
				Expect(err).To(MatchError(ContainSubstring("unexpected or invalid ARN")))
			})
		})
		When("create service account fails", func() {
			BeforeEach(func() {
				cfg.Karpenter.CreateServiceAccount = api.Disabled()
				ft := &fakeTask{
					err: errors.New("nope"),
				}
				fakeStackManager.NewTasksToCreateIAMServiceAccountsReturns(&tasks.TaskTree{
					Tasks: []tasks.Task{ft},
				})
			})
			It("errors", func() {
				install := &karpenteractions.Installer{
					StackManager:       fakeStackManager,
					CTL:                ctl,
					Config:             cfg,
					KarpenterInstaller: fakeKarpenterInstaller,
					ClientSet:          fakeClientSet,
				}
				err := install.Create()
				Expect(err).To(MatchError(ContainSubstring("failed to create/attach service account: failed to install Karpenter on cluster")))
			})
		})
		When("fails to fetch the identity mapping config map", func() {
			BeforeEach(func() {
				fakeClientSet = fake.NewSimpleClientset()
				fakeClientSet.PrependReactor("get", "configmaps", func(action core.Action) (bool, runtime.Object, error) {
					return true, nil, errors.New("nope")
				})
			})
			It("errors", func() {
				install := &karpenteractions.Installer{
					StackManager:       fakeStackManager,
					CTL:                ctl,
					Config:             cfg,
					KarpenterInstaller: fakeKarpenterInstaller,
					ClientSet:          fakeClientSet,
				}
				err := install.Create()
				Expect(err).To(MatchError(ContainSubstring("failed to create client for auth config: getting auth ConfigMap: nope")))
			})
		})
		When("creating the iam mapping configmap fails", func() {
			BeforeEach(func() {
				fakeClientSet = fake.NewSimpleClientset()
				fakeClientSet.PrependReactor("create", "configmaps", func(action core.Action) (bool, runtime.Object, error) {
					return true, nil, errors.New("nope")
				})
			})
			It("errors", func() {
				install := &karpenteractions.Installer{
					StackManager:       fakeStackManager,
					CTL:                ctl,
					Config:             cfg,
					KarpenterInstaller: fakeKarpenterInstaller,
					ClientSet:          fakeClientSet,
				}
				err := install.Create()
				Expect(err).To(MatchError(ContainSubstring("failed to save the identity config: nope")))
			})
		})
		When("updating the iam mapping configmap fails", func() {
			BeforeEach(func() {
				fakeClientSet = fake.NewSimpleClientset()
				fakeClientSet.PrependReactor("update", "configmaps", func(action core.Action) (bool, runtime.Object, error) {
					return true, nil, errors.New("nope")
				})
				fakeClientSet.PrependReactor("get", "configmaps", func(action core.Action) (bool, runtime.Object, error) {
					return true, &corev1.ConfigMap{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "aws-auth",
							Namespace: "kube-system",
							UID:       "uid",
						},
						Data: map[string]string{
							"something": "there",
						},
					}, nil
				})
			})
			It("errors", func() {
				install := &karpenteractions.Installer{
					StackManager:       fakeStackManager,
					CTL:                ctl,
					Config:             cfg,
					KarpenterInstaller: fakeKarpenterInstaller,
					ClientSet:          fakeClientSet,
				}
				err := install.Create()
				Expect(err).To(MatchError(ContainSubstring("failed to save the identity config: nope")))
			})
		})
		When("createServiceAccount is enabled", func() {
			BeforeEach(func() {
				cfg.Karpenter.CreateServiceAccount = api.Enabled()
			})
			It("eksctl should only create the role with a specific policy", func() {
				fakeKarpenterInstaller.InstallReturns(nil)
				install := &karpenteractions.Installer{
					StackManager:       fakeStackManager,
					CTL:                ctl,
					Config:             cfg,
					KarpenterInstaller: fakeKarpenterInstaller,
					ClientSet:          fakeClientSet,
				}
				Expect(install.Create()).To(Succeed())
				Expect(fakeKarpenterInstaller.InstallCallCount()).To(Equal(1))
				accounts, _, _ := fakeStackManager.NewTasksToCreateIAMServiceAccountsArgsForCall(0)
				Expect(accounts).NotTo(BeEmpty())
				Expect(api.IsEnabled(accounts[0].RoleOnly)).To(BeTrue())
			})
		})
		When("createServiceAccount is disabled", func() {
			BeforeEach(func() {
				cfg.Karpenter.CreateServiceAccount = api.Disabled()
			})
			It("eksctl should create the service account", func() {
				fakeKarpenterInstaller.InstallReturns(nil)
				install := &karpenteractions.Installer{
					StackManager:       fakeStackManager,
					CTL:                ctl,
					Config:             cfg,
					KarpenterInstaller: fakeKarpenterInstaller,
					ClientSet:          fakeClientSet,
				}
				Expect(install.Create()).To(Succeed())
				Expect(fakeKarpenterInstaller.InstallCallCount()).To(Equal(1))
				accounts, _, _ := fakeStackManager.NewTasksToCreateIAMServiceAccountsArgsForCall(0)
				Expect(accounts).NotTo(BeEmpty())
				Expect(accounts[0].RoleOnly).To(BeNil())
				policyARN := fmt.Sprintf("arn:aws:iam::123456789012:policy/eksctl-%s-%s", builder.KarpenterManagedPolicy, cfg.Metadata.Name)
				Expect(accounts[0].AttachPolicyARNs).To(ConsistOf(policyARN))
			})
		})
	})
})

var (
	azA, azB                       = "us-west-2a", "us-west-2b"
	privateSubnet1, privateSubnet2 = "subnet-1", "subnet-2"
	publicSubnet1, publicSubnet2   = "subnet-3", "subnet-4"
)

func vpcConfig() *api.ClusterVPC {
	disable := api.ClusterDisableNAT
	return &api.ClusterVPC{
		NAT: &api.ClusterNAT{
			Gateway: &disable,
		},
		Subnets: &api.ClusterSubnets{
			Public: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
				azB: {
					ID: publicSubnet2,
					CIDR: &ipnet.IPNet{
						IPNet: net.IPNet{
							IP:   []byte{192, 168, 0, 0},
							Mask: []byte{255, 255, 224, 0},
						},
					},
				},
				azA: {
					ID: publicSubnet1,
					CIDR: &ipnet.IPNet{
						IPNet: net.IPNet{
							IP:   []byte{192, 168, 32, 0},
							Mask: []byte{255, 255, 224, 0},
						},
					},
				},
			}),
			Private: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
				azB: {
					ID: privateSubnet2,
					CIDR: &ipnet.IPNet{
						IPNet: net.IPNet{
							IP:   []byte{192, 168, 96, 0},
							Mask: []byte{255, 255, 224, 0},
						},
					},
				},
				azA: {
					ID: privateSubnet1,
					CIDR: &ipnet.IPNet{
						IPNet: net.IPNet{
							IP:   []byte{192, 168, 128, 0},
							Mask: []byte{255, 255, 224, 0},
						},
					},
				},
			}),
		},
	}
}
