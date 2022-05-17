package nodegroup_test

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	cftypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/stretchr/testify/mock"

	"github.com/weaveworks/eksctl/pkg/utils/tasks"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/pkg/errors"

	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	managerFakes "github.com/weaveworks/eksctl/pkg/cfn/manager/fakes"
	utilFakes "github.com/weaveworks/eksctl/pkg/ctl/cmdutils/filter/fakes"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/eks/fakes"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

type ngEntry struct {
	version       string
	opts          nodegroup.CreateOpts
	mockCalls     func(*fakes.FakeKubeProvider, *fakes.FakeNodeGroupInitialiser, *utilFakes.FakeNodegroupFilter, *mockprovider.MockProvider, *managerFakes.FakeStackManager)
	expectedCalls func(*fakes.FakeKubeProvider, *fakes.FakeNodeGroupInitialiser, *utilFakes.FakeNodegroupFilter)
	expectedErr   error
}

var _ = DescribeTable("Create", func(t ngEntry) {
	cfg := newClusterConfig()
	cfg.Metadata.Version = t.version

	p := mockprovider.NewMockProvider()
	ctl := &eks.ClusterProvider{
		AWSProvider: p,
		Status: &eks.ProviderStatus{
			ClusterInfo: &eks.ClusterInfo{
				Cluster: testutils.NewFakeCluster("my-cluster", ""),
			},
		},
	}

	m := nodegroup.New(cfg, ctl, nil)

	k := &fakes.FakeKubeProvider{}
	m.MockKubeProvider(k)

	init := &fakes.FakeNodeGroupInitialiser{}
	m.MockNodeGroupService(init)

	stackManager := &managerFakes.FakeStackManager{}
	m.SetStackManager(stackManager)

	ngFilter := utilFakes.FakeNodegroupFilter{}
	if t.mockCalls != nil {
		t.mockCalls(k, init, &ngFilter, p, stackManager)
	}

	err := m.Create(context.Background(), t.opts, &ngFilter)

	if t.expectedErr != nil {
		Expect(err).To(MatchError(ContainSubstring(t.expectedErr.Error())))
	} else {
		Expect(err).NotTo(HaveOccurred())
	}
	if t.expectedCalls != nil {
		t.expectedCalls(k, init, &ngFilter)
	}
},
	Entry("fails when cluster version is not supported", ngEntry{
		version:     "1.14",
		expectedErr: fmt.Errorf("invalid version, %s is no longer supported, supported values: auto, default, latest, %s\nsee also: https://docs.aws.amazon.com/eks/latest/userguide/kubernetes-versions.html", "1.14", strings.Join(api.SupportedVersions(), ", ")),
	}),

	Entry("when cluster is unowned, fails to load VPC from config if config is not supplied", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter, p *mockprovider.MockProvider, s *managerFakes.FakeStackManager) {
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

	Entry("fails to set instance types to instances matched by instance selector criteria", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter, p *mockprovider.MockProvider, s *managerFakes.FakeStackManager) {
			init.ExpandInstanceSelectorOptionsReturns(errors.New("err"))
			defaultProviderMocks(p, defaultOutput)
		},
		expectedCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, _ *utilFakes.FakeNodegroupFilter) {
			Expect(k.NewRawClientCallCount()).To(Equal(1))
			Expect(k.ServerVersionCallCount()).To(Equal(1))
			Expect(init.NewAWSSelectorSessionCallCount()).To(Equal(1))
			Expect(init.ExpandInstanceSelectorOptionsCallCount()).To(Equal(1))
		},
		expectedErr: errors.New("err"),
	}),

	Entry("fails when cluster is not compatible with ng config", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter, p *mockprovider.MockProvider, s *managerFakes.FakeStackManager) {
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
		expectedCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, _ *utilFakes.FakeNodegroupFilter) {
			Expect(k.NewRawClientCallCount()).To(Equal(1))
			Expect(k.ServerVersionCallCount()).To(Equal(1))
			Expect(init.NewAWSSelectorSessionCallCount()).To(Equal(1))
			Expect(init.ExpandInstanceSelectorOptionsCallCount()).To(Equal(1))
		},
		expectedErr: errors.Wrap(errors.New("shared node security group missing, to fix this run 'eksctl update cluster --name=my-cluster --region='"), "cluster compatibility check failed")}),

	Entry("fails when it cannot validate legacy subnets for ng", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter, p *mockprovider.MockProvider, s *managerFakes.FakeStackManager) {
			defaultProviderMocks(p, defaultOutput)
			init.ValidateLegacySubnetsForNodeGroupsReturns(errors.New("err"))
		},
		expectedCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, _ *utilFakes.FakeNodegroupFilter) {
			Expect(k.NewRawClientCallCount()).To(Equal(1))
			Expect(k.ServerVersionCallCount()).To(Equal(1))
			Expect(init.NewAWSSelectorSessionCallCount()).To(Equal(1))
			Expect(init.ExpandInstanceSelectorOptionsCallCount()).To(Equal(1))
			Expect(init.ValidateLegacySubnetsForNodeGroupsCallCount()).To(Equal(1))
		},
		expectedErr: errors.New("err"),
	}),

	Entry("fails when existing local ng stacks in config file is not listed", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter, p *mockprovider.MockProvider, s *managerFakes.FakeStackManager) {
			f.SetOnlyLocalReturns(errors.New("err"))
			defaultProviderMocks(p, defaultOutput)
		},
		expectedCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
			Expect(k.NewRawClientCallCount()).To(Equal(1))
			Expect(k.ServerVersionCallCount()).To(Equal(1))
			Expect(init.NewAWSSelectorSessionCallCount()).To(Equal(1))
			Expect(init.ExpandInstanceSelectorOptionsCallCount()).To(Equal(1))
			Expect(f.SetOnlyLocalCallCount()).To(Equal(1))
		},
		expectedErr: errors.Wrap(errors.New("shared node security group missing, to fix this run 'eksctl update cluster --name=my-cluster --region='"), "cluster compatibility check failed"),
	}),

	Entry("fails to evaluate whether aws-node uses IRSA", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter, p *mockprovider.MockProvider, s *managerFakes.FakeStackManager) {
			init.DoesAWSNodeUseIRSAReturns(true, errors.New("err"))
			defaultProviderMocks(p, defaultOutput)
		},
		expectedCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
			Expect(k.NewRawClientCallCount()).To(Equal(1))
			Expect(k.ServerVersionCallCount()).To(Equal(1))
			Expect(init.NewAWSSelectorSessionCallCount()).To(Equal(1))
			Expect(init.ExpandInstanceSelectorOptionsCallCount()).To(Equal(1))
			Expect(f.SetOnlyLocalCallCount()).To(Equal(1))
			Expect(init.DoesAWSNodeUseIRSACallCount()).To(Equal(1))
		},
		expectedErr: errors.New("err"),
	}),

	Entry("[happy path] creates nodegroup with no options", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter, p *mockprovider.MockProvider, s *managerFakes.FakeStackManager) {
			defaultProviderMocks(p, defaultOutput)
			createNoopTasks(s)
		},
		expectedCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
			Expect(k.NewRawClientCallCount()).To(Equal(1))
			Expect(k.ServerVersionCallCount()).To(Equal(1))
			Expect(init.NewAWSSelectorSessionCallCount()).To(Equal(1))
			Expect(init.ExpandInstanceSelectorOptionsCallCount()).To(Equal(1))
			Expect(f.SetOnlyLocalCallCount()).To(Equal(1))
			Expect(init.DoesAWSNodeUseIRSACallCount()).To(Equal(1))
		},
	}),

	Entry("[happy path] creates nodegroup with all the options", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter, p *mockprovider.MockProvider, s *managerFakes.FakeStackManager) {
			defaultProviderMocks(p, defaultOutput)
			createNoopTasks(s)
		},
		opts: nodegroup.CreateOpts{
			DryRun:                    true,
			UpdateAuthConfigMap:       true,
			InstallNeuronDevicePlugin: true,
			InstallNvidiaDevicePlugin: true,
			SkipOutdatedAddonsCheck:   true,
			ConfigFileProvided:        true,
		},
		expectedCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
			Expect(k.NewRawClientCallCount()).To(Equal(1))
			Expect(k.ServerVersionCallCount()).To(Equal(1))
			Expect(init.NewAWSSelectorSessionCallCount()).To(Equal(1))
			Expect(init.ExpandInstanceSelectorOptionsCallCount()).To(Equal(1))
			Expect(f.SetOnlyLocalCallCount()).To(Equal(1))
		},
	}),
)

func createNoopTasks(stackManager *managerFakes.FakeStackManager) {
	noopTask := &tasks.GenericTask{
		Doer: func() error {
			return nil
		},
	}
	stackManager.NewUnmanagedNodeGroupTaskReturns(&tasks.TaskTree{
		Tasks: []tasks.Task{noopTask},
	})
	stackManager.NewClusterCompatTaskReturns(noopTask)
	stackManager.NewManagedNodeGroupTaskReturns(nil)
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
				Name: "my-ng",
			}},
		},
		ManagedNodeGroups: []*api.ManagedNodeGroup{{
			NodeGroupBase: &api.NodeGroupBase{
				Name: "my-ng",
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
				Outputs: output,
			},
		},
	}, nil)
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
			ResourcesVpcConfig: &ekstypes.VpcConfigResponse{
				ClusterSecurityGroupId: aws.String("csg-1234"),
				EndpointPublicAccess:   true,
				PublicAccessCidrs:      []string{"1.2.3.4/24", "1.2.3.4/12"},
				SecurityGroupIds:       []string{"sg-1", "sg-2"},
				SubnetIds:              []string{"sub-1", "sub-2"},
				VpcId:                  aws.String("vpc-1"),
			},
			Status: "CREATE_COMPLETE",
			Tags: map[string]string{
				api.ClusterNameTag: "eksctl-my-cluster-cluster",
			},
			Version: aws.String("1.22"),
		},
	}, nil)
}
