package nodegroup_test

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cftypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	core "k8s.io/client-go/testing"

	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	utilFakes "github.com/weaveworks/eksctl/pkg/ctl/cmdutils/filter/fakes"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/eks/fakes"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

type ngEntry struct {
	version             string
	opts                nodegroup.CreateOpts
	mockCalls           func(*fakes.FakeKubeProvider, *utilFakes.FakeNodegroupFilter, *mockprovider.MockProvider, *fake.Clientset)
	refreshCluster      bool
	updateClusterConfig func(*api.ClusterConfig)

	expectedCalls      func(*fakes.FakeKubeProvider, *utilFakes.FakeNodegroupFilter)
	expectedErr        error
	expectedRefreshErr string
}

type stackManagerDelegate struct {
	manager.StackManager
}

func (s *stackManagerDelegate) NewUnmanagedNodeGroupTask(context.Context, []*api.NodeGroup, bool, vpc.Importer) *tasks.TaskTree {
	return &tasks.TaskTree{
		Tasks: []tasks.Task{noopTask},
	}
}

func (s *stackManagerDelegate) NewManagedNodeGroupTask(context.Context, []*api.ManagedNodeGroup, bool, vpc.Importer) *tasks.TaskTree {
	return nil
}

func (s *stackManagerDelegate) FixClusterCompatibility(_ context.Context) error {
	return nil
}

func (s *stackManagerDelegate) ClusterHasDedicatedVPC(_ context.Context) (bool, error) {
	return false, nil
}

var _ = DescribeTable("Create", func(t ngEntry) {
	cfg := newClusterConfig()
	cfg.Metadata.Version = t.version
	if t.updateClusterConfig != nil {
		t.updateClusterConfig(cfg)
	}

	p := mockprovider.NewMockProvider()
	ctl := &eks.ClusterProvider{
		AWSProvider: p,
		Status: &eks.ProviderStatus{
			ClusterInfo: &eks.ClusterInfo{
				Cluster: testutils.NewFakeCluster("my-cluster", ""),
			},
		},
	}

	clientset := fake.NewSimpleClientset()
	m := nodegroup.New(cfg, ctl, clientset, nil)

	k := &fakes.FakeKubeProvider{}
	m.MockKubeProvider(k)

	stackManager := &stackManagerDelegate{
		StackManager: m.GetStackManager(),
	}
	m.SetStackManager(stackManager)

	ngFilter := utilFakes.FakeNodegroupFilter{}

	if t.mockCalls != nil {
		t.mockCalls(k, &ngFilter, p, clientset)
	}
	if t.refreshCluster {
		err := ctl.RefreshClusterStatus(context.Background(), cfg)
		if t.expectedRefreshErr != "" {
			Expect(err).To(MatchError(ContainSubstring(t.expectedRefreshErr)))
			return
		}
		Expect(err).NotTo(HaveOccurred())
	}

	err := m.Create(context.Background(), t.opts, &ngFilter)

	if t.expectedErr != nil {
		Expect(err).To(MatchError(ContainSubstring(t.expectedErr.Error())))
	} else {
		Expect(err).NotTo(HaveOccurred())
	}
	if t.expectedCalls != nil {
		t.expectedCalls(k, &ngFilter)
	}
},
	Entry("fails when cluster version is not supported", ngEntry{
		version:     "1.14",
		expectedErr: fmt.Errorf("invalid version, %s is no longer supported, supported values: auto, default, latest, %s\nsee also: https://docs.aws.amazon.com/eks/latest/userguide/kubernetes-versions.html", "1.14", strings.Join(api.SupportedVersions(), ", ")),
	}),

	Entry("when cluster is unowned, fails to load VPC from config if config is not supplied", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, f *utilFakes.FakeNodegroupFilter, p *mockprovider.MockProvider, _ *fake.Clientset) {
			k.NewRawClientReturns(&kubernetes.RawClient{}, nil)
			k.ServerVersionReturns("1.17", nil)
			p.MockCloudFormation().On("ListStacks", mock.Anything, mock.Anything).Return(&cloudformation.ListStacksOutput{
				StackSummaries: []cftypes.StackSummary{
					{
						StackName:   aws.String("eksctl-my-cluster-cluster"),
						StackStatus: "CREATE_COMPLETE",
					},
				},
			}, nil)
			p.MockCloudFormation().On("DescribeStacks", mock.Anything, mock.Anything).Return(&cloudformation.DescribeStacksOutput{
				Stacks: []cftypes.Stack{
					{
						StackName:   aws.String("eksctl-my-cluster-cluster"),
						StackStatus: "CREATE_COMPLETE",
					},
				},
			}, nil)
		},
		expectedErr: errors.Wrapf(errors.New("VPC configuration required for creating nodegroups on clusters not owned by eksctl: vpc.subnets, vpc.id, vpc.securityGroup"), "loading VPC spec for cluster %q", "my-cluster"),
	}),

	Entry("fails when cluster is not compatible with ng config", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, f *utilFakes.FakeNodegroupFilter, p *mockprovider.MockProvider, _ *fake.Clientset) {
			// no shared security group will trigger a compatibility check failure later in the call chain.
			output := []cftypes.Output{
				{
					OutputKey:   aws.String("ClusterSecurityGroupId"),
					OutputValue: aws.String("csg-1234"),
				},
				{
					OutputKey:   aws.String("SecurityGroup"),
					OutputValue: aws.String("sg-1"),
				},
				{
					OutputKey:   aws.String("VPC"),
					OutputValue: aws.String("vpc-1"),
				},
			}
			defaultProviderMocks(p, output)
		},
		expectedCalls: func(k *fakes.FakeKubeProvider, _ *utilFakes.FakeNodegroupFilter) {
			Expect(k.NewRawClientCallCount()).To(Equal(1))
			Expect(k.ServerVersionCallCount()).To(Equal(1))
		},
		expectedErr: errors.Wrap(errors.New("shared node security group missing, to fix this run 'eksctl update cluster --name=my-cluster --region='"), "cluster compatibility check failed")}),

	Entry("fails when existing local ng stacks in config file is not listed", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, f *utilFakes.FakeNodegroupFilter, p *mockprovider.MockProvider, _ *fake.Clientset) {
			f.SetOnlyLocalReturns(errors.New("err"))
			defaultProviderMocks(p, defaultOutput)
		},
		expectedCalls: func(k *fakes.FakeKubeProvider, f *utilFakes.FakeNodegroupFilter) {
			Expect(k.NewRawClientCallCount()).To(Equal(1))
			Expect(k.ServerVersionCallCount()).To(Equal(1))
			Expect(f.SetOnlyLocalCallCount()).To(Equal(1))
		},
		expectedErr: errors.New("err"),
	}),

	Entry("fails to evaluate whether aws-node uses IRSA", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, f *utilFakes.FakeNodegroupFilter, p *mockprovider.MockProvider, c *fake.Clientset) {
			c.PrependReactor("get", "serviceaccounts", func(action core.Action) (bool, runtime.Object, error) {
				return true, nil, errors.New("failed to determine if aws-node uses IRSA")
			})
			defaultProviderMocks(p, defaultOutput)
		},
		expectedCalls: func(k *fakes.FakeKubeProvider, f *utilFakes.FakeNodegroupFilter) {
			Expect(k.NewRawClientCallCount()).To(Equal(1))
			Expect(k.ServerVersionCallCount()).To(Equal(1))
			Expect(f.SetOnlyLocalCallCount()).To(Equal(1))
		},
		expectedErr: errors.New("failed to determine if aws-node uses IRSA"),
	}),

	Entry("fails to create managed nodegroups on Outposts", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, f *utilFakes.FakeNodegroupFilter, p *mockprovider.MockProvider, _ *fake.Clientset) {
			mockProviderWithOutpostConfig(p, defaultOutput, &ekstypes.OutpostConfigResponse{
				OutpostArns:              []string{"arn:aws:outposts:us-west-2:1234:outpost/op-1234"},
				ControlPlaneInstanceType: aws.String("m5a.large"),
			})
		},
		opts: nodegroup.CreateOpts{
			DryRun:                    true,
			UpdateAuthConfigMap:       true,
			InstallNeuronDevicePlugin: true,
			InstallNvidiaDevicePlugin: true,
			SkipOutdatedAddonsCheck:   true,
			ConfigFileProvided:        false,
		},
		refreshCluster: true,
		expectedCalls: func(k *fakes.FakeKubeProvider, f *utilFakes.FakeNodegroupFilter) {
			Expect(k.NewRawClientCallCount()).To(Equal(0))
			Expect(k.ServerVersionCallCount()).To(Equal(0))
			Expect(f.SetOnlyLocalCallCount()).To(Equal(0))
		},
		expectedErr: errors.New("Managed Nodegroups are not supported on Outposts; please rerun the command with --managed=false"),
	}),

	Entry("fails to create managed nodegroups on Outposts with a config file", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, f *utilFakes.FakeNodegroupFilter, p *mockprovider.MockProvider, _ *fake.Clientset) {
			mockProviderWithOutpostConfig(p, defaultOutput, &ekstypes.OutpostConfigResponse{
				OutpostArns:              []string{"arn:aws:outposts:us-west-2:1234:outpost/op-1234"},
				ControlPlaneInstanceType: aws.String("m5a.large"),
			})
		},
		opts: nodegroup.CreateOpts{
			DryRun:                    true,
			UpdateAuthConfigMap:       true,
			InstallNeuronDevicePlugin: true,
			InstallNvidiaDevicePlugin: true,
			SkipOutdatedAddonsCheck:   true,
			ConfigFileProvided:        true,
		},
		refreshCluster: true,
		expectedCalls: func(k *fakes.FakeKubeProvider, f *utilFakes.FakeNodegroupFilter) {
			Expect(k.NewRawClientCallCount()).To(Equal(0))
			Expect(k.ServerVersionCallCount()).To(Equal(0))
			Expect(f.SetOnlyLocalCallCount()).To(Equal(0))
		},
		expectedErr: errors.New("Managed Nodegroups are not supported on Outposts"),
	}),

	Entry("Outpost config does not match cluster's Outpost config", ngEntry{
		updateClusterConfig: func(c *api.ClusterConfig) {
			c.Outpost = &api.Outpost{
				ControlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
			}
		},
		mockCalls: func(k *fakes.FakeKubeProvider, f *utilFakes.FakeNodegroupFilter, p *mockprovider.MockProvider, _ *fake.Clientset) {
			mockProviderWithOutpostConfig(p, defaultOutput, &ekstypes.OutpostConfigResponse{
				OutpostArns:              []string{"arn:aws:outposts:us-west-2:1234:outpost/op-5678"},
				ControlPlaneInstanceType: aws.String("m5a.large"),
			})
		},
		opts: nodegroup.CreateOpts{
			DryRun:                    true,
			UpdateAuthConfigMap:       true,
			InstallNeuronDevicePlugin: true,
			InstallNvidiaDevicePlugin: true,
			SkipOutdatedAddonsCheck:   true,
			ConfigFileProvided:        true,
		},
		refreshCluster: true,
		expectedCalls: func(k *fakes.FakeKubeProvider, f *utilFakes.FakeNodegroupFilter) {
			Expect(k.NewRawClientCallCount()).To(Equal(0))
			Expect(k.ServerVersionCallCount()).To(Equal(0))
			Expect(f.SetOnlyLocalCallCount()).To(Equal(0))
		},
		expectedRefreshErr: fmt.Sprintf("outpost.controlPlaneOutpostARN %q does not match the cluster's Outpost ARN %q", "arn:aws:outposts:us-west-2:1234:outpost/op-1234", "arn:aws:outposts:us-west-2:1234:outpost/op-5678"),
	}),

	Entry("Outpost config set but control plane is not on Outposts", ngEntry{
		updateClusterConfig: func(c *api.ClusterConfig) {
			c.Outpost = &api.Outpost{
				ControlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
			}
		},
		mockCalls: func(k *fakes.FakeKubeProvider, f *utilFakes.FakeNodegroupFilter, p *mockprovider.MockProvider, _ *fake.Clientset) {
			defaultProviderMocks(p, defaultOutput)
		},
		opts: nodegroup.CreateOpts{
			DryRun:                    true,
			UpdateAuthConfigMap:       true,
			InstallNeuronDevicePlugin: true,
			InstallNvidiaDevicePlugin: true,
			SkipOutdatedAddonsCheck:   true,
			ConfigFileProvided:        true,
		},
		refreshCluster: true,
		expectedCalls: func(k *fakes.FakeKubeProvider, f *utilFakes.FakeNodegroupFilter) {
			Expect(k.NewRawClientCallCount()).To(Equal(0))
			Expect(k.ServerVersionCallCount()).To(Equal(0))
			Expect(f.SetOnlyLocalCallCount()).To(Equal(0))
		},
		expectedRefreshErr: "outpost.controlPlaneOutpostARN is set but control plane is not on Outposts",
	}),

	Entry("API server unreachable when creating a nodegroup on Outposts", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, f *utilFakes.FakeNodegroupFilter, p *mockprovider.MockProvider, _ *fake.Clientset) {
			k.NewRawClientReturns(nil, &kubernetes.APIServerUnreachableError{
				Err: errors.New("timeout"),
			})
			mockProviderWithOutpostConfig(p, defaultOutput, &ekstypes.OutpostConfigResponse{
				OutpostArns:              []string{"arn:aws:outposts:us-west-2:1234:outpost/op-1234"},
				ControlPlaneInstanceType: aws.String("m5a.large"),
			})
		},
		opts: nodegroup.CreateOpts{
			DryRun:                    true,
			UpdateAuthConfigMap:       true,
			InstallNeuronDevicePlugin: true,
			InstallNvidiaDevicePlugin: true,
			SkipOutdatedAddonsCheck:   true,
			ConfigFileProvided:        true,
		},
		refreshCluster: true,
		updateClusterConfig: func(c *api.ClusterConfig) {
			c.ManagedNodeGroups = nil
			c.Outpost = &api.Outpost{
				ControlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
			}
		},
		expectedCalls: func(k *fakes.FakeKubeProvider, f *utilFakes.FakeNodegroupFilter) {
			Expect(k.NewRawClientCallCount()).To(Equal(1))
		},
		expectedErr: errors.New("eksctl requires connectivity to the API server to create nodegroups;" +
			" please ensure the Outpost VPC is associated with your local gateway and you are able to connect to" +
			" the API server before rerunning the command: timeout"),
	}),

	Entry("API server unreachable in a cluster with private-only endpoint access", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, f *utilFakes.FakeNodegroupFilter, p *mockprovider.MockProvider, _ *fake.Clientset) {
			k.NewRawClientReturns(nil, &kubernetes.APIServerUnreachableError{
				Err: errors.New("timeout"),
			})
			mockProviderWithConfig(p, defaultOutput, &ekstypes.VpcConfigResponse{
				ClusterSecurityGroupId: aws.String("csg-1234"),
				EndpointPublicAccess:   false,
				EndpointPrivateAccess:  true,
				SecurityGroupIds:       []string{"sg-1"},
				SubnetIds:              []string{"sub-1", "sub-2"},
				VpcId:                  aws.String("vpc-1"),
			}, nil)
		},
		opts: nodegroup.CreateOpts{
			DryRun:                    true,
			UpdateAuthConfigMap:       true,
			InstallNeuronDevicePlugin: true,
			InstallNvidiaDevicePlugin: true,
			SkipOutdatedAddonsCheck:   true,
			ConfigFileProvided:        true,
		},
		refreshCluster: true,
		updateClusterConfig: func(c *api.ClusterConfig) {
			c.ManagedNodeGroups = nil
		},
		expectedCalls: func(k *fakes.FakeKubeProvider, f *utilFakes.FakeNodegroupFilter) {
			Expect(k.NewRawClientCallCount()).To(Equal(1))
		},
		expectedErr: errors.New("eksctl requires connectivity to the API server to create nodegroups;" +
			" please run eksctl from an environment that has access to the API server: timeout"),
	}),

	Entry("creates nodegroups on Outposts", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, f *utilFakes.FakeNodegroupFilter, p *mockprovider.MockProvider, _ *fake.Clientset) {
			mockProviderWithOutpostConfig(p, defaultOutput, &ekstypes.OutpostConfigResponse{
				OutpostArns:              []string{"arn:aws:outposts:us-west-2:1234:outpost/op-1234"},
				ControlPlaneInstanceType: aws.String("m5a.large"),
			})
		},
		opts: nodegroup.CreateOpts{
			DryRun:                    true,
			UpdateAuthConfigMap:       true,
			InstallNeuronDevicePlugin: true,
			InstallNvidiaDevicePlugin: true,
			SkipOutdatedAddonsCheck:   true,
			ConfigFileProvided:        true,
		},
		refreshCluster: true,
		updateClusterConfig: func(c *api.ClusterConfig) {
			c.ManagedNodeGroups = nil
		},
		expectedCalls: func(k *fakes.FakeKubeProvider, f *utilFakes.FakeNodegroupFilter) {
			Expect(k.NewRawClientCallCount()).To(Equal(1))
			Expect(k.ServerVersionCallCount()).To(Equal(1))
			Expect(f.SetOnlyLocalCallCount()).To(Equal(1))
		},
	}),

	Entry("[happy path] creates nodegroup with no options", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, f *utilFakes.FakeNodegroupFilter, p *mockprovider.MockProvider, _ *fake.Clientset) {
			defaultProviderMocks(p, defaultOutput)
		},
		expectedCalls: func(k *fakes.FakeKubeProvider, f *utilFakes.FakeNodegroupFilter) {
			Expect(k.NewRawClientCallCount()).To(Equal(1))
			Expect(k.ServerVersionCallCount()).To(Equal(1))
			Expect(f.SetOnlyLocalCallCount()).To(Equal(1))
		},
	}),

	Entry("[happy path] creates nodegroup with all the options", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, f *utilFakes.FakeNodegroupFilter, p *mockprovider.MockProvider, _ *fake.Clientset) {
			defaultProviderMocks(p, defaultOutput)
		},
		opts: nodegroup.CreateOpts{
			DryRun:                    true,
			UpdateAuthConfigMap:       true,
			InstallNeuronDevicePlugin: true,
			InstallNvidiaDevicePlugin: true,
			SkipOutdatedAddonsCheck:   true,
			ConfigFileProvided:        true,
		},
		expectedCalls: func(k *fakes.FakeKubeProvider, f *utilFakes.FakeNodegroupFilter) {
			Expect(k.NewRawClientCallCount()).To(Equal(1))
			Expect(k.ServerVersionCallCount()).To(Equal(1))
			Expect(f.SetOnlyLocalCallCount()).To(Equal(1))
		},
	}),
)

var noopTask = &tasks.GenericTask{
	Doer: func() error {
		return nil
	},
}

func newClusterConfig() *api.ClusterConfig {
	return &api.ClusterConfig{
		TypeMeta: api.ClusterConfigTypeMeta(),
		Metadata: &api.ClusterMeta{
			Name:    "my-cluster",
			Version: api.DefaultVersion,
		},
		Status: &api.ClusterStatus{
			Endpoint:                 "https://localhost/",
			CertificateAuthorityData: []byte("dGVzdAo="),
		},
		IAM: api.NewClusterIAM(),
		VPC: api.NewClusterVPC(false),
		CloudWatch: &api.ClusterCloudWatch{
			ClusterLogging: &api.ClusterCloudWatchLogging{},
		},
		PrivateCluster: &api.PrivateCluster{},
		NodeGroups: []*api.NodeGroup{{
			NodeGroupBase: &api.NodeGroupBase{
				Name:             "my-ng",
				AMIFamily:        api.NodeImageFamilyAmazonLinux2,
				AMI:              "ami-123",
				SSH:              &api.NodeGroupSSH{Allow: api.Disabled()},
				InstanceSelector: &api.InstanceSelector{},
			}},
		},
		ManagedNodeGroups: []*api.ManagedNodeGroup{{
			NodeGroupBase: &api.NodeGroupBase{
				Name:             "my-ng",
				AMIFamily:        api.NodeImageFamilyAmazonLinux2,
				SSH:              &api.NodeGroupSSH{Allow: api.Disabled()},
				InstanceSelector: &api.InstanceSelector{},
			}},
		},
	}
}

var defaultOutput = []cftypes.Output{
	{
		OutputKey:   aws.String("ClusterSecurityGroupId"),
		OutputValue: aws.String("csg-1234"),
	},
	{
		OutputKey:   aws.String("SecurityGroup"),
		OutputValue: aws.String("sg-1"),
	},
	{
		OutputKey:   aws.String("VPC"),
		OutputValue: aws.String("vpc-1"),
	},
	{
		OutputKey:   aws.String("SharedNodeSecurityGroup"),
		OutputValue: aws.String("sg-1"),
	},
}

func defaultProviderMocks(p *mockprovider.MockProvider, output []cftypes.Output) {
	mockProviderWithConfig(p, output, nil, nil)
}

func mockProviderWithOutpostConfig(p *mockprovider.MockProvider, describeStacksOutput []cftypes.Output, outpostConfig *ekstypes.OutpostConfigResponse) {
	mockProviderWithConfig(p, describeStacksOutput, nil, outpostConfig)
}

func mockProviderWithConfig(p *mockprovider.MockProvider, describeStacksOutput []cftypes.Output, vpcConfigRes *ekstypes.VpcConfigResponse, outpostConfig *ekstypes.OutpostConfigResponse) {
	p.MockCloudFormation().On("ListStacks", mock.Anything, mock.Anything).Return(&cloudformation.ListStacksOutput{
		StackSummaries: []cftypes.StackSummary{
			{
				StackName:   aws.String("eksctl-my-cluster-cluster"),
				StackStatus: "CREATE_COMPLETE",
			},
		},
	}, nil)
	p.MockCloudFormation().On("DescribeStacks", mock.Anything, mock.Anything).Return(&cloudformation.DescribeStacksOutput{
		Stacks: []cftypes.Stack{
			{
				StackName:   aws.String("eksctl-my-cluster-cluster"),
				StackStatus: "CREATE_COMPLETE",
				Tags: []cftypes.Tag{
					{
						Key:   aws.String(api.ClusterNameTag),
						Value: aws.String("eksctl-my-cluster-cluster"),
					},
				},
				Outputs: describeStacksOutput,
			},
		},
	}, nil)
	if vpcConfigRes == nil {
		vpcConfigRes = &ekstypes.VpcConfigResponse{
			ClusterSecurityGroupId: aws.String("csg-1234"),
			EndpointPublicAccess:   true,
			PublicAccessCidrs:      []string{"1.2.3.4/24", "1.2.3.4/12"},
			SecurityGroupIds:       []string{"sg-1", "sg-2"},
			SubnetIds:              []string{"sub-1", "sub-2"},
			VpcId:                  aws.String("vpc-1"),
		}
	}
	p.MockEKS().On("DescribeCluster", mock.Anything, mock.Anything).Return(&awseks.DescribeClusterOutput{
		Cluster: &ekstypes.Cluster{
			CertificateAuthority: &ekstypes.Certificate{
				Data: aws.String("dGVzdAo="),
			},
			Endpoint:                aws.String("endpoint"),
			Arn:                     aws.String("arn"),
			KubernetesNetworkConfig: nil,
			Logging:                 nil,
			Name:                    aws.String("my-cluster"),
			PlatformVersion:         aws.String("1.22"),
			ResourcesVpcConfig:      vpcConfigRes,
			OutpostConfig:           outpostConfig,
			Status:                  "CREATE_COMPLETE",
			Tags: map[string]string{
				api.ClusterNameTag: "eksctl-my-cluster-cluster",
			},
			Version: aws.String("1.22"),
		},
	}, nil)

	p.MockEC2().On("DescribeImages", mock.Anything, mock.Anything).
		Return(&ec2.DescribeImagesOutput{
			Images: []ec2types.Image{
				{
					ImageId:        aws.String("ami-123"),
					State:          ec2types.ImageStateAvailable,
					OwnerId:        aws.String("123"),
					RootDeviceType: ec2types.DeviceTypeEbs,
					RootDeviceName: aws.String("/dev/sda1"),
					BlockDeviceMappings: []ec2types.BlockDeviceMapping{
						{
							DeviceName: aws.String("/dev/sda1"),
							Ebs: &ec2types.EbsBlockDevice{
								Encrypted: aws.Bool(false),
							},
						},
					},
				},
			},
		}, nil)
}
