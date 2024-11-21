package nodegroup_test

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cftypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	core "k8s.io/client-go/testing"

	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"
	ngfakes "github.com/weaveworks/eksctl/pkg/actions/nodegroup/fakes"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	utilfakes "github.com/weaveworks/eksctl/pkg/ctl/cmdutils/filter/fakes"
	"github.com/weaveworks/eksctl/pkg/eks"
	eksfakes "github.com/weaveworks/eksctl/pkg/eks/fakes"
	"github.com/weaveworks/eksctl/pkg/iam"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

type ngEntry struct {
	version             string
	opts                nodegroup.CreateOpts
	mockCalls           func(mockCalls)
	refreshCluster      bool
	updateClusterConfig func(*api.ClusterConfig)

	expectedCalls      func(expectedCalls)
	expectedErr        error
	expectedRefreshErr string
}

type mockCalls struct {
	kubeProvider    *eksfakes.FakeKubeProvider
	nodeGroupFilter *utilfakes.FakeNodegroupFilter
	mockProvider    *mockprovider.MockProvider
	clientset       *fake.Clientset
}

type expectedCalls struct {
	kubeProvider         *eksfakes.FakeKubeProvider
	nodeGroupFilter      *utilfakes.FakeNodegroupFilter
	nodeGroupTaskCreator *ngfakes.FakeNodeGroupTaskCreator
	clientset            *fake.Clientset
}

type vpcSubnets struct {
	publicIDs  []string
	privateIDs []string
}

//counterfeiter:generate -o fakes/fake_nodegroup_task_creator.go . nodeGroupTaskCreator
type nodeGroupTaskCreator interface {
	NewUnmanagedNodeGroupTask(context.Context, []*api.NodeGroup, bool, bool, bool, vpc.Importer) *tasks.TaskTree
}

type stackManagerDelegate struct {
	manager.StackManager
	ngTaskCreator nodeGroupTaskCreator
}

func (s *stackManagerDelegate) NewUnmanagedNodeGroupTask(ctx context.Context, nodeGroups []*api.NodeGroup, forceAddCNIPolicy, skipEgressRules, disableAccessEntryCreation bool, vpcImporter vpc.Importer, nodeGroupParallelism int) *tasks.TaskTree {
	return s.ngTaskCreator.NewUnmanagedNodeGroupTask(ctx, nodeGroups, forceAddCNIPolicy, skipEgressRules, disableAccessEntryCreation, vpcImporter)
}

func (s *stackManagerDelegate) NewManagedNodeGroupTask(context.Context, []*api.ManagedNodeGroup, bool, vpc.Importer, int) *tasks.TaskTree {
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

	k := &eksfakes.FakeKubeProvider{}
	m.MockKubeProvider(k)

	var ngTaskCreator ngfakes.FakeNodeGroupTaskCreator
	ngTaskCreator.NewUnmanagedNodeGroupTaskStub = func(_ context.Context, _ []*api.NodeGroup, _, _, _ bool, _ vpc.Importer) *tasks.TaskTree {
		return &tasks.TaskTree{
			Tasks: []tasks.Task{noopTask},
		}
	}
	stackManager := &stackManagerDelegate{
		ngTaskCreator: &ngTaskCreator,
		StackManager:  m.GetStackManager(),
	}
	m.SetStackManager(stackManager)

	var ngFilter utilfakes.FakeNodegroupFilter
	ngFilter.MatchReturns(true)

	if t.mockCalls != nil {
		t.mockCalls(mockCalls{
			kubeProvider:    k,
			nodeGroupFilter: &ngFilter,
			mockProvider:    p,
			clientset:       clientset,
		})
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
		t.expectedCalls(expectedCalls{
			kubeProvider:         k,
			nodeGroupFilter:      &ngFilter,
			nodeGroupTaskCreator: &ngTaskCreator,
			clientset:            clientset,
		})
	}
},
	Entry("when cluster is unowned, fails to load VPC from config if config is not supplied", ngEntry{
		mockCalls: func(m mockCalls) {
			m.kubeProvider.NewRawClientReturns(&kubernetes.RawClient{}, nil)
			m.kubeProvider.ServerVersionReturns("1.17", nil)
			m.mockProvider.MockCloudFormation().On("ListStacks", mock.Anything, mock.Anything, mock.Anything).Return(&cloudformation.ListStacksOutput{
				StackSummaries: []cftypes.StackSummary{
					{
						StackName:   aws.String("eksctl-my-cluster-cluster"),
						StackStatus: "CREATE_COMPLETE",
					},
				},
			}, nil)
			m.mockProvider.MockCloudFormation().On("DescribeStacks", mock.Anything, mock.Anything).Return(&cloudformation.DescribeStacksOutput{
				Stacks: []cftypes.Stack{
					{
						StackName:   aws.String("eksctl-my-cluster-cluster"),
						StackStatus: "CREATE_COMPLETE",
					},
				},
			}, nil)
		},
		expectedErr: errors.New(`loading VPC spec for cluster "my-cluster": VPC configuration required for creating nodegroups on clusters not owned by eksctl: vpc.subnets, vpc.id, vpc.securityGroup`),
	}),

	Entry("when cluster is unowned and vpc.securityGroup contains external egress rules, it fails validation", ngEntry{
		updateClusterConfig: makeUnownedClusterConfig,
		mockCalls: func(m mockCalls) {
			mockProviderForUnownedCluster(m.mockProvider, m.kubeProvider, ec2types.SecurityGroupRule{
				Description:         aws.String("Allow control plane to communicate with a custom nodegroup on a custom port"),
				FromPort:            aws.Int32(8443),
				ToPort:              aws.Int32(8443),
				GroupId:             aws.String("sg-custom"),
				IpProtocol:          aws.String("https"),
				IsEgress:            aws.Bool(true),
				SecurityGroupRuleId: aws.String("sgr-5"),
			})

		},
		expectedErr: errors.New("vpc.securityGroup (sg-custom) has egress rules that were not attached by eksctl; vpc.securityGroup should not contain any non-default external egress rules on a cluster not created by eksctl (rule ID: sgr-5)"),
	}),

	Entry("when cluster is unowned and vpc.securityGroup contains a default egress rule, it passes validation but fails if DescribeImages fails", ngEntry{
		updateClusterConfig: makeUnownedClusterConfig,
		mockCalls: func(m mockCalls) {
			mockProviderForUnownedCluster(m.mockProvider, m.kubeProvider, ec2types.SecurityGroupRule{
				Description:         aws.String(""),
				CidrIpv4:            aws.String("0.0.0.0/0"),
				FromPort:            aws.Int32(-1),
				ToPort:              aws.Int32(-1),
				GroupId:             aws.String("sg-custom"),
				IpProtocol:          aws.String("-1"),
				IsEgress:            aws.Bool(true),
				SecurityGroupRuleId: aws.String("sgr-5"),
			})
			m.mockProvider.MockEC2().On("DescribeImages", mock.Anything, mock.Anything).Return(nil, errors.New("DescribeImages error"))

		},
		expectedErr: errors.New("DescribeImages error"),
	}),

	Entry("when cluster is unowned and vpc.securityGroup contains no external egress rules, it passes validation but fails if DescribeImages fails", ngEntry{
		updateClusterConfig: makeUnownedClusterConfig,
		mockCalls: func(m mockCalls) {
			mockProviderForUnownedCluster(m.mockProvider, m.kubeProvider)
			m.mockProvider.MockEC2().On("DescribeImages", mock.Anything, mock.Anything).Return(nil, errors.New("DescribeImages error"))

		},
		expectedErr: errors.New("DescribeImages error"),
	}),

	Entry("fails when cluster is not compatible with ng config", ngEntry{
		mockCalls: func(m mockCalls) {
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
			defaultProviderMocks(m.mockProvider, output)
		},
		expectedCalls: func(e expectedCalls) {
			Expect(e.kubeProvider.NewRawClientCallCount()).To(Equal(1))
		},
		expectedErr: errors.New("cluster compatibility check failed: shared node security group missing, to fix this run 'eksctl update cluster --name=my-cluster --region='"),
	}),

	Entry("fails when nodegroup uses privateNetworking:true and there's no private subnet within vpc", ngEntry{
		updateClusterConfig: func(c *api.ClusterConfig) {
			c.NodeGroups[0].PrivateNetworking = true
		},
		mockCalls: func(m mockCalls) {
			mockProviderWithVPCSubnets(m.mockProvider, &vpcSubnets{
				publicIDs: []string{"subnet-public-1", "subnet-public-2"},
			})
		},
		expectedErr: fmt.Errorf("all private subnets from vpc-1, that the cluster was originally created on, have been deleted; to create private nodegroups within vpc-1 please manually set valid private subnets via nodeGroup.SubnetIDs"),
	}),

	Entry("fails when nodegroup uses privateNetworking:false and there's no public subnet within vpc", ngEntry{
		mockCalls: func(m mockCalls) {
			mockProviderWithVPCSubnets(m.mockProvider, &vpcSubnets{
				publicIDs: []string{"subnet-private-1", "subnet-private-2"},
			})
		},
		expectedErr: fmt.Errorf("all public subnets from vpc-1, that the cluster was originally created on, have been deleted; to create public nodegroups within vpc-1 please manually set valid public subnets via nodeGroup.SubnetIDs"),
	}),

	Entry("fails when nodegroup uses privateNetworking:true and there's no private subnet within az", ngEntry{
		updateClusterConfig: func(c *api.ClusterConfig) {
			c.NodeGroups[0].PrivateNetworking = true
			c.NodeGroups[0].AvailabilityZones = []string{"us-west-2b"}
		},
		mockCalls: func(m mockCalls) {
			mockProviderWithVPCSubnets(m.mockProvider, &vpcSubnets{
				publicIDs:  []string{"subnet-public-1", "subnet-public-2"},
				privateIDs: []string{"subnet-private-1"},
			})
		},
		expectedErr: fmt.Errorf("all private subnets from us-west-2b, that the cluster was originally created on, have been deleted; to create private nodegroups within us-west-2b please manually set valid private subnets via nodeGroup.SubnetIDs"),
	}),

	Entry("fails when nodegroup uses privateNetworking:false and there's no private subnet within az", ngEntry{
		updateClusterConfig: func(c *api.ClusterConfig) {
			c.NodeGroups[0].AvailabilityZones = []string{"us-west-2a", "us-west-2b"}
			c.VPC.Subnets = &api.ClusterSubnets{
				Private: api.AZSubnetMapping{
					"private-1": api.AZSubnetSpec{
						ID: "subnet-private-1",
					},
					"private-2": api.AZSubnetSpec{
						ID: "subnet-private-2",
					},
				},
				Public: api.AZSubnetMapping{
					"public-1": api.AZSubnetSpec{
						ID: "subnet-public-2",
					},
				},
			}
		},
		mockCalls: func(m mockCalls) {
			mockProviderWithVPCSubnets(m.mockProvider, &vpcSubnets{
				publicIDs:  []string{"subnet-public-2"},
				privateIDs: []string{"subnet-private-1", "subnet-private-2"},
			})
		},
		expectedErr: fmt.Errorf("all public subnets from us-west-2a, that the cluster was originally created on, have been deleted; to create public nodegroups within us-west-2a please manually set valid public subnets via nodeGroup.SubnetIDs"),
	}),

	Entry("fails when existing local ng stacks in config file is not listed", ngEntry{
		mockCalls: func(m mockCalls) {
			m.nodeGroupFilter.SetOnlyLocalReturns(errors.New("err"))
			defaultProviderMocks(m.mockProvider, defaultOutput)
		},
		expectedCalls: func(e expectedCalls) {
			Expect(e.kubeProvider.NewRawClientCallCount()).To(Equal(1))
			Expect(e.nodeGroupFilter.SetOnlyLocalCallCount()).To(Equal(1))
		},
		expectedErr: errors.New("err"),
	}),

	Entry("fails to evaluate whether aws-node uses IRSA", ngEntry{
		mockCalls: func(m mockCalls) {
			m.clientset.PrependReactor("get", "serviceaccounts", func(action core.Action) (bool, runtime.Object, error) {
				return true, nil, errors.New("failed to determine if aws-node uses IRSA")
			})
			defaultProviderMocks(m.mockProvider, defaultOutput)
		},
		expectedCalls: func(e expectedCalls) {
			Expect(e.kubeProvider.NewRawClientCallCount()).To(Equal(1))
			Expect(e.nodeGroupFilter.SetOnlyLocalCallCount()).To(Equal(1))
		},
		expectedErr: errors.New("failed to determine if aws-node uses IRSA"),
	}),

	Entry("fails to create managed nodegroups on Outposts", ngEntry{
		mockCalls: func(m mockCalls) {
			mockProviderWithOutpostConfig(m.mockProvider, defaultOutput, &ekstypes.OutpostConfigResponse{
				OutpostArns:              []string{"arn:aws:outposts:us-west-2:1234:outpost/op-1234"},
				ControlPlaneInstanceType: aws.String("m5a.large"),
			})
		},
		opts: nodegroup.CreateOpts{
			DryRunSettings: nodegroup.DryRunSettings{
				DryRun:    true,
				OutStream: os.Stdout,
			},
			UpdateAuthConfigMap:       api.Enabled(),
			InstallNeuronDevicePlugin: true,
			InstallNvidiaDevicePlugin: true,
			SkipOutdatedAddonsCheck:   true,
			ConfigFileProvided:        false,
		},
		refreshCluster: true,
		expectedCalls: func(e expectedCalls) {
			Expect(e.kubeProvider.NewRawClientCallCount()).To(Equal(0))
			Expect(e.kubeProvider.ServerVersionCallCount()).To(Equal(0))
			Expect(e.nodeGroupFilter.SetOnlyLocalCallCount()).To(Equal(0))
		},
		expectedErr: errors.New("Managed Nodegroups are not supported on Outposts; please rerun the command with --managed=false"),
	}),

	Entry("fails to create managed nodegroups on Outposts with a config file", ngEntry{
		mockCalls: func(m mockCalls) {
			mockProviderWithOutpostConfig(m.mockProvider, defaultOutput, &ekstypes.OutpostConfigResponse{
				OutpostArns:              []string{"arn:aws:outposts:us-west-2:1234:outpost/op-1234"},
				ControlPlaneInstanceType: aws.String("m5a.large"),
			})
		},
		opts: nodegroup.CreateOpts{
			DryRunSettings: nodegroup.DryRunSettings{
				DryRun:    true,
				OutStream: os.Stdout,
			},
			UpdateAuthConfigMap:       api.Enabled(),
			InstallNeuronDevicePlugin: true,
			InstallNvidiaDevicePlugin: true,
			SkipOutdatedAddonsCheck:   true,
			ConfigFileProvided:        true,
		},
		refreshCluster: true,
		expectedCalls: func(e expectedCalls) {
			Expect(e.kubeProvider.NewRawClientCallCount()).To(Equal(0))
			Expect(e.kubeProvider.ServerVersionCallCount()).To(Equal(0))
			Expect(e.nodeGroupFilter.SetOnlyLocalCallCount()).To(Equal(0))
		},
		expectedErr: errors.New("Managed Nodegroups are not supported on Outposts"),
	}),

	Entry("Outpost config does not match cluster's Outpost config", ngEntry{
		updateClusterConfig: func(c *api.ClusterConfig) {
			c.Outpost = &api.Outpost{
				ControlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
			}
		},
		mockCalls: func(m mockCalls) {
			mockProviderWithOutpostConfig(m.mockProvider, defaultOutput, &ekstypes.OutpostConfigResponse{
				OutpostArns:              []string{"arn:aws:outposts:us-west-2:1234:outpost/op-5678"},
				ControlPlaneInstanceType: aws.String("m5a.large"),
			})
		},
		opts: nodegroup.CreateOpts{
			DryRunSettings: nodegroup.DryRunSettings{
				DryRun:    true,
				OutStream: os.Stdout,
			},
			UpdateAuthConfigMap:       api.Enabled(),
			InstallNeuronDevicePlugin: true,
			InstallNvidiaDevicePlugin: true,
			SkipOutdatedAddonsCheck:   true,
			ConfigFileProvided:        true,
		},
		refreshCluster: true,
		expectedCalls: func(e expectedCalls) {
			Expect(e.kubeProvider.NewRawClientCallCount()).To(Equal(0))
			Expect(e.kubeProvider.ServerVersionCallCount()).To(Equal(0))
			Expect(e.nodeGroupFilter.SetOnlyLocalCallCount()).To(Equal(0))
		},
		expectedRefreshErr: fmt.Sprintf("outpost.controlPlaneOutpostARN %q does not match the cluster's Outpost ARN %q", "arn:aws:outposts:us-west-2:1234:outpost/op-1234", "arn:aws:outposts:us-west-2:1234:outpost/op-5678"),
	}),

	Entry("Outpost config set but control plane is not on Outposts", ngEntry{
		updateClusterConfig: func(c *api.ClusterConfig) {
			c.Outpost = &api.Outpost{
				ControlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
			}
		},
		mockCalls: func(m mockCalls) {
			defaultProviderMocks(m.mockProvider, defaultOutput)
		},
		opts: nodegroup.CreateOpts{
			DryRunSettings: nodegroup.DryRunSettings{
				DryRun:    true,
				OutStream: os.Stdout,
			},
			UpdateAuthConfigMap:       api.Enabled(),
			InstallNeuronDevicePlugin: true,
			InstallNvidiaDevicePlugin: true,
			SkipOutdatedAddonsCheck:   true,
			ConfigFileProvided:        true,
		},
		refreshCluster: true,
		expectedCalls: func(e expectedCalls) {
			Expect(e.kubeProvider.NewRawClientCallCount()).To(Equal(0))
			Expect(e.kubeProvider.ServerVersionCallCount()).To(Equal(0))
			Expect(e.nodeGroupFilter.SetOnlyLocalCallCount()).To(Equal(0))
		},
		expectedRefreshErr: "outpost.controlPlaneOutpostARN is set but control plane is not on Outposts",
	}),

	Entry("API server unreachable when creating a nodegroup on Outposts", ngEntry{
		mockCalls: func(m mockCalls) {
			m.kubeProvider.NewRawClientReturns(nil, &kubernetes.APIServerUnreachableError{
				Err: errors.New("timeout"),
			})
			mockProviderWithOutpostConfig(m.mockProvider, defaultOutput, &ekstypes.OutpostConfigResponse{
				OutpostArns:              []string{"arn:aws:outposts:us-west-2:1234:outpost/op-1234"},
				ControlPlaneInstanceType: aws.String("m5a.large"),
			})
		},
		opts: nodegroup.CreateOpts{
			DryRunSettings: nodegroup.DryRunSettings{
				DryRun:    true,
				OutStream: os.Stdout,
			},
			UpdateAuthConfigMap:       api.Enabled(),
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
		expectedCalls: func(e expectedCalls) {
			Expect(e.kubeProvider.NewRawClientCallCount()).To(Equal(1))
		},
		expectedErr: errors.New("eksctl requires connectivity to the API server to create nodegroups;" +
			" please ensure the Outpost VPC is associated with your local gateway and you are able to connect to" +
			" the API server before rerunning the command: timeout"),
	}),

	Entry("API server unreachable in a cluster with private-only endpoint access", ngEntry{
		mockCalls: func(m mockCalls) {
			m.kubeProvider.NewRawClientReturns(nil, &kubernetes.APIServerUnreachableError{
				Err: errors.New("timeout"),
			})
			mockProviderWithConfig(m.mockProvider, defaultOutput, nil, &ekstypes.VpcConfigResponse{
				ClusterSecurityGroupId: aws.String("csg-1234"),
				EndpointPublicAccess:   false,
				EndpointPrivateAccess:  true,
				SecurityGroupIds:       []string{"sg-1"},
				SubnetIds:              []string{"sub-1", "sub-2"},
				VpcId:                  aws.String("vpc-1"),
			}, nil, nil)
		},
		opts: nodegroup.CreateOpts{
			DryRunSettings: nodegroup.DryRunSettings{
				DryRun:    true,
				OutStream: os.Stdout,
			},
			UpdateAuthConfigMap:       api.Enabled(),
			InstallNeuronDevicePlugin: true,
			InstallNvidiaDevicePlugin: true,
			SkipOutdatedAddonsCheck:   true,
			ConfigFileProvided:        true,
		},
		refreshCluster: true,
		updateClusterConfig: func(c *api.ClusterConfig) {
			c.ManagedNodeGroups = nil
		},
		expectedCalls: func(e expectedCalls) {
			Expect(e.kubeProvider.NewRawClientCallCount()).To(Equal(1))
		},
		expectedErr: errors.New("eksctl requires connectivity to the API server to create nodegroups;" +
			" please run eksctl from an environment that has access to the API server: timeout"),
	}),

	Entry("creates nodegroups on Outposts", ngEntry{
		mockCalls: func(m mockCalls) {
			mockProviderWithOutpostConfig(m.mockProvider, defaultOutput, &ekstypes.OutpostConfigResponse{
				OutpostArns:              []string{"arn:aws:outposts:us-west-2:1234:outpost/op-1234"},
				ControlPlaneInstanceType: aws.String("m5a.large"),
			})
		},
		opts: nodegroup.CreateOpts{
			DryRunSettings: nodegroup.DryRunSettings{
				DryRun:    true,
				OutStream: os.Stdout,
			},
			UpdateAuthConfigMap:       api.Enabled(),
			InstallNeuronDevicePlugin: true,
			InstallNvidiaDevicePlugin: true,
			SkipOutdatedAddonsCheck:   true,
			ConfigFileProvided:        true,
		},
		refreshCluster: true,
		updateClusterConfig: func(c *api.ClusterConfig) {
			c.ManagedNodeGroups = nil
		},
		expectedCalls: func(e expectedCalls) {
			Expect(e.kubeProvider.NewRawClientCallCount()).To(Equal(1))
			Expect(e.nodeGroupFilter.SetOnlyLocalCallCount()).To(Equal(1))
		},
	}),

	Entry("fails to create nodegroup when authenticationMode is API and updateAuthConfigMap is false", ngEntry{
		opts: nodegroup.CreateOpts{
			UpdateAuthConfigMap: api.Disabled(),
		},
		mockCalls: func(m mockCalls) {
			mockProviderWithConfig(m.mockProvider, defaultOutput, nil, nil, nil, &ekstypes.AccessConfigResponse{
				AuthenticationMode: ekstypes.AuthenticationModeApi,
			})
			defaultProviderMocks(m.mockProvider, defaultOutput)
		},
		refreshCluster: true,
		expectedCalls: func(e expectedCalls) {
			Expect(e.kubeProvider.NewRawClientCallCount()).To(Equal(0))
			Expect(e.kubeProvider.ServerVersionCallCount()).To(Equal(0))
			Expect(e.nodeGroupFilter.SetOnlyLocalCallCount()).To(Equal(0))
		},

		expectedErr: errors.New("--update-auth-configmap is not supported when authenticationMode is set to API"),
	}),

	Entry("fails to create nodegroup when authenticationMode is API and updateAuthConfigMap is true", ngEntry{
		opts: nodegroup.CreateOpts{
			UpdateAuthConfigMap: api.Enabled(),
		},
		mockCalls: func(m mockCalls) {
			mockProviderWithConfig(m.mockProvider, defaultOutput, nil, nil, nil, &ekstypes.AccessConfigResponse{
				AuthenticationMode: ekstypes.AuthenticationModeApi,
			})
			defaultProviderMocks(m.mockProvider, defaultOutput)
		},
		refreshCluster: true,
		expectedCalls: func(e expectedCalls) {
			Expect(e.kubeProvider.NewRawClientCallCount()).To(Equal(0))
			Expect(e.kubeProvider.ServerVersionCallCount()).To(Equal(0))
			Expect(e.nodeGroupFilter.SetOnlyLocalCallCount()).To(Equal(0))
		},

		expectedErr: errors.New("--update-auth-configmap is not supported when authenticationMode is set to API"),
	}),

	Entry("creates nodegroup using access entries when authenticationMode is API_AND_CONFIG_MAP and updateAuthConfigMap is not supplied", ngEntry{
		mockCalls: func(m mockCalls) {
			mockProviderWithConfig(m.mockProvider, defaultOutput, nil, nil, nil, &ekstypes.AccessConfigResponse{
				AuthenticationMode: ekstypes.AuthenticationModeApiAndConfigMap,
			})
			defaultProviderMocks(m.mockProvider, defaultOutput)
		},
		expectedCalls: func(e expectedCalls) {
			Expect(e.kubeProvider.NewRawClientCallCount()).To(Equal(1))
			Expect(e.nodeGroupFilter.SetOnlyLocalCallCount()).To(Equal(1))
			Expect(e.nodeGroupTaskCreator.NewUnmanagedNodeGroupTaskCallCount()).To(Equal(1))
			_, _, _, _, disableAccessEntryCreation, _ := e.nodeGroupTaskCreator.NewUnmanagedNodeGroupTaskArgsForCall(0)
			Expect(disableAccessEntryCreation).To(BeFalse())
			Expect(getIAMIdentities(e.clientset)).To(HaveLen(0))
		},
	}),

	Entry("creates nodegroup using aws-auth ConfigMap when authenticationMode is CONFIG_MAP and updateAuthConfigMap is true", ngEntry{
		mockCalls: func(m mockCalls) {
			mockProviderWithConfig(m.mockProvider, defaultOutput, nil, nil, nil, &ekstypes.AccessConfigResponse{
				AuthenticationMode: ekstypes.AuthenticationModeConfigMap,
			})
			defaultProviderMocks(m.mockProvider, defaultOutput)
		},
		opts: nodegroup.CreateOpts{
			UpdateAuthConfigMap: api.Enabled(),
		},
		refreshCluster: true,
		expectedCalls:  expectedCallsForAWSAuth,
	}),

	Entry("creates nodegroup using aws-auth ConfigMap when authenticationMode is CONFIG_MAP and updateAuthConfigMap is not supplied", ngEntry{
		mockCalls: func(m mockCalls) {
			mockProviderWithConfig(m.mockProvider, defaultOutput, nil, nil, nil, &ekstypes.AccessConfigResponse{
				AuthenticationMode: ekstypes.AuthenticationModeConfigMap,
			})
			defaultProviderMocks(m.mockProvider, defaultOutput)
		},
		opts: nodegroup.CreateOpts{
			UpdateAuthConfigMap: api.Enabled(),
		},
		refreshCluster: true,
		expectedCalls:  expectedCallsForAWSAuth,
	}),

	Entry("creates nodegroup but does not use either aws-auth ConfigMap or access entries when authenticationMode is API_AND_CONFIG_MAP and updateAuthConfigMap is false", ngEntry{
		mockCalls: func(m mockCalls) {
			mockProviderWithConfig(m.mockProvider, defaultOutput, nil, nil, nil, &ekstypes.AccessConfigResponse{
				AuthenticationMode: ekstypes.AuthenticationModeApiAndConfigMap,
			})
			defaultProviderMocks(m.mockProvider, defaultOutput)
		},
		refreshCluster: true,
		opts: nodegroup.CreateOpts{
			UpdateAuthConfigMap: api.Disabled(),
		},
		expectedCalls: func(e expectedCalls) {
			Expect(e.kubeProvider.NewRawClientCallCount()).To(Equal(1))
			Expect(e.nodeGroupFilter.SetOnlyLocalCallCount()).To(Equal(1))
			Expect(e.nodeGroupTaskCreator.NewUnmanagedNodeGroupTaskCallCount()).To(Equal(1))
			_, _, _, _, disableAccessEntryCreation, _ := e.nodeGroupTaskCreator.NewUnmanagedNodeGroupTaskArgsForCall(0)
			Expect(disableAccessEntryCreation).To(BeTrue())
			Expect(getIAMIdentities(e.clientset)).To(HaveLen(0))
		},
	}),

	Entry("creates nodegroup but does not use either aws-auth ConfigMap or access entries when authenticationMode is CONFIG_MAP and updateAuthConfigMap is false", ngEntry{
		mockCalls: func(m mockCalls) {
			mockProviderWithConfig(m.mockProvider, defaultOutput, nil, nil, nil, &ekstypes.AccessConfigResponse{
				AuthenticationMode: ekstypes.AuthenticationModeConfigMap,
			})
			defaultProviderMocks(m.mockProvider, defaultOutput)
		},
		refreshCluster: true,
		opts: nodegroup.CreateOpts{
			UpdateAuthConfigMap: api.Disabled(),
		},
		expectedCalls: func(e expectedCalls) {
			Expect(e.kubeProvider.NewRawClientCallCount()).To(Equal(1))
			Expect(e.nodeGroupFilter.SetOnlyLocalCallCount()).To(Equal(1))
			Expect(e.nodeGroupTaskCreator.NewUnmanagedNodeGroupTaskCallCount()).To(Equal(1))
			_, _, _, _, disableAccessEntryCreation, _ := e.nodeGroupTaskCreator.NewUnmanagedNodeGroupTaskArgsForCall(0)
			Expect(disableAccessEntryCreation).To(BeTrue())
			Expect(getIAMIdentities(e.clientset)).To(HaveLen(0))
		},
	}),

	Entry("authorizes nodegroups using aws-auth ConfigMap when authenticationMode is API_AND_CONFIG_MAP and updateAuthConfigMap is true", ngEntry{
		mockCalls: func(m mockCalls) {
			mockProviderWithConfig(m.mockProvider, defaultOutput, nil, nil, nil, &ekstypes.AccessConfigResponse{
				AuthenticationMode: ekstypes.AuthenticationModeApiAndConfigMap,
			})
			defaultProviderMocks(m.mockProvider, defaultOutput)
		},
		refreshCluster: true,
		opts: nodegroup.CreateOpts{
			UpdateAuthConfigMap: api.Enabled(),
		},
		expectedCalls: expectedCallsForAWSAuth,
	}),

	Entry("[happy path] creates nodegroup with no options", ngEntry{
		mockCalls: func(m mockCalls) {
			defaultProviderMocks(m.mockProvider, defaultOutput)
		},
		expectedCalls: func(e expectedCalls) {
			Expect(e.kubeProvider.NewRawClientCallCount()).To(Equal(1))
			Expect(e.nodeGroupFilter.SetOnlyLocalCallCount()).To(Equal(1))
		},
	}),

	Entry("[happy path] creates nodegroup with all the options", ngEntry{
		mockCalls: func(m mockCalls) {
			defaultProviderMocks(m.mockProvider, defaultOutput)
		},
		refreshCluster: true,
		opts: nodegroup.CreateOpts{
			DryRunSettings: nodegroup.DryRunSettings{
				DryRun:    true,
				OutStream: os.Stdout,
			},
			UpdateAuthConfigMap:       api.Enabled(),
			InstallNeuronDevicePlugin: true,
			InstallNvidiaDevicePlugin: true,
			SkipOutdatedAddonsCheck:   true,
			ConfigFileProvided:        true,
		},
		expectedCalls: func(e expectedCalls) {
			Expect(e.kubeProvider.NewRawClientCallCount()).To(Equal(1))
			Expect(e.nodeGroupFilter.SetOnlyLocalCallCount()).To(Equal(1))
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
		AccessConfig:   &api.AccessConfig{},
		PrivateCluster: &api.PrivateCluster{},
		NodeGroups: []*api.NodeGroup{{
			NodeGroupBase: &api.NodeGroupBase{
				Name:             "my-ng",
				AMIFamily:        api.NodeImageFamilyAmazonLinux2,
				AMI:              "ami-123",
				SSH:              &api.NodeGroupSSH{Allow: api.Disabled()},
				InstanceSelector: &api.InstanceSelector{},
				ScalingConfig:    &api.ScalingConfig{},
				IAM: &api.NodeGroupIAM{
					InstanceRoleARN: "arn:aws:iam::1234567890:role/my-ng",
				},
			}},
		},
		ManagedNodeGroups: []*api.ManagedNodeGroup{{
			NodeGroupBase: &api.NodeGroupBase{
				Name:             "my-ng",
				AMIFamily:        api.NodeImageFamilyAmazonLinux2,
				SSH:              &api.NodeGroupSSH{Allow: api.Disabled()},
				InstanceSelector: &api.InstanceSelector{},
				ScalingConfig:    &api.ScalingConfig{},
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
	{
		OutputKey:   aws.String("SubnetsPublic"),
		OutputValue: aws.String("subnet-public-1,subnet-public-2"),
	},
	{
		OutputKey:   aws.String("SubnetsPrivate"),
		OutputValue: aws.String("subnet-private-1,subnet-private-2"),
	},
}

func getIAMIdentities(clientset kubernetes.Interface) []iam.Identity {
	acm, err := authconfigmap.NewFromClientSet(clientset)
	Expect(err).NotTo(HaveOccurred())
	identities, err := acm.GetIdentities()
	Expect(err).NotTo(HaveOccurred())
	return identities
}

func expectedCallsForAWSAuth(e expectedCalls) {
	Expect(e.kubeProvider.NewRawClientCallCount()).To(Equal(1))
	Expect(e.nodeGroupFilter.SetOnlyLocalCallCount()).To(Equal(1))
	Expect(e.nodeGroupTaskCreator.NewUnmanagedNodeGroupTaskCallCount()).To(Equal(1))
	_, _, _, _, disableAccessEntryCreation, _ := e.nodeGroupTaskCreator.NewUnmanagedNodeGroupTaskArgsForCall(0)
	Expect(disableAccessEntryCreation).To(BeTrue())
	identities := getIAMIdentities(e.clientset)
	Expect(identities).To(HaveLen(1))
	for _, id := range identities {
		roleIdentity, ok := id.(iam.RoleIdentity)
		Expect(ok).To(BeTrue())
		Expect(roleIdentity).To(Equal(iam.RoleIdentity{
			RoleARN: "arn:aws:iam::1234567890:role/my-ng",
			KubernetesIdentity: iam.KubernetesIdentity{
				KubernetesUsername: "system:node:{{EC2PrivateDNSName}}",
				KubernetesGroups:   []string{"system:bootstrappers", "system:nodes"},
			},
		}))
	}
}

func defaultProviderMocks(p *mockprovider.MockProvider, output []cftypes.Output) {
	mockProviderWithConfig(p, output, nil, nil, nil, nil)
}

func mockProviderWithOutpostConfig(p *mockprovider.MockProvider, describeStacksOutput []cftypes.Output, outpostConfig *ekstypes.OutpostConfigResponse) {
	mockProviderWithConfig(p, describeStacksOutput, nil, nil, outpostConfig, nil)
}

func mockProviderWithVPCSubnets(p *mockprovider.MockProvider, subnets *vpcSubnets) {
	mockProviderWithConfig(p, defaultOutput, subnets, nil, nil, nil)
}

func mockProviderWithConfig(p *mockprovider.MockProvider, describeStacksOutput []cftypes.Output, subnets *vpcSubnets, vpcConfigRes *ekstypes.VpcConfigResponse, outpostConfig *ekstypes.OutpostConfigResponse, accessConfig *ekstypes.AccessConfigResponse) {
	p.MockCloudFormation().On("ListStacks", mock.Anything, mock.Anything, mock.Anything).Return(&cloudformation.ListStacksOutput{
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
	if accessConfig == nil {
		accessConfig = &ekstypes.AccessConfigResponse{
			AuthenticationMode: ekstypes.AuthenticationModeApiAndConfigMap,
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
			AccessConfig:            accessConfig,
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

	if subnets == nil {
		subnets = &vpcSubnets{
			publicIDs:  []string{"subnet-public-1", "subnet-public-2"},
			privateIDs: []string{"subnet-private-1", "subnet-private-2"},
		}
	}

	subnetPublic1 := ec2types.Subnet{
		SubnetId:            aws.String("subnet-public-1"),
		CidrBlock:           aws.String("192.168.64.0/20"),
		AvailabilityZone:    aws.String("us-west-2a"),
		VpcId:               aws.String("vpc-1"),
		MapPublicIpOnLaunch: aws.Bool(true),
	}
	subnetPrivate1 := ec2types.Subnet{
		SubnetId:            aws.String("subnet-private-1"),
		CidrBlock:           aws.String("192.168.128.0/20"),
		AvailabilityZone:    aws.String("us-west-2a"),
		VpcId:               aws.String("vpc-1"),
		MapPublicIpOnLaunch: aws.Bool(false),
	}
	subnetPublic2 := ec2types.Subnet{
		SubnetId:            aws.String("subnet-public-2"),
		CidrBlock:           aws.String("192.168.80.0/20"),
		AvailabilityZone:    aws.String("us-west-2b"),
		VpcId:               aws.String("vpc-1"),
		MapPublicIpOnLaunch: aws.Bool(true),
	}
	subnetPrivate2 := ec2types.Subnet{
		SubnetId:            aws.String("subnet-private-2"),
		CidrBlock:           aws.String("192.168.32.0/20"),
		AvailabilityZone:    aws.String("us-west-2b"),
		VpcId:               aws.String("vpc-1"),
		MapPublicIpOnLaunch: aws.Bool(false),
	}

	subnetsForID := map[string]ec2types.Subnet{
		"subnet-public-1":  subnetPublic1,
		"subnet-private-1": subnetPrivate1,
		"subnet-public-2":  subnetPublic2,
		"subnet-private-2": subnetPrivate2,
	}

	mockDescribeSubnets := func(mp *mockprovider.MockProvider, vpcID string, subnetIDs []string) {
		var subnets []ec2types.Subnet
		for _, id := range subnetIDs {
			subnets = append(subnets, subnetsForID[id])
		}
		if vpcID == "" {
			mp.MockEC2().On("DescribeSubnets", mock.Anything, &ec2.DescribeSubnetsInput{
				SubnetIds: subnetIDs,
			}, mock.Anything).Return(&ec2.DescribeSubnetsOutput{Subnets: subnets}, nil)
			return
		}
		mp.MockEC2().On("DescribeSubnets", mock.Anything, &ec2.DescribeSubnetsInput{
			Filters: []ec2types.Filter{
				{
					Name:   aws.String("vpc-id"),
					Values: []string{vpcID},
				},
			},
		}, mock.Anything).Return(&ec2.DescribeSubnetsOutput{Subnets: subnets}, nil)
	}

	mockDescribeSubnets(p, "", subnets.publicIDs)
	mockDescribeSubnets(p, "", subnets.privateIDs)
	mockDescribeSubnets(p, "vpc-1", append(subnets.publicIDs, subnets.privateIDs...))

	p.MockEC2().On("DescribeVpcs", mock.Anything, mock.Anything).Return(&ec2.DescribeVpcsOutput{
		Vpcs: []ec2types.Vpc{
			{
				CidrBlock: aws.String("192.168.0.0/20"),
				VpcId:     aws.String("vpc-1"),
				CidrBlockAssociationSet: []ec2types.VpcCidrBlockAssociation{
					{
						CidrBlock: aws.String("192.168.0.0/20"),
					},
				},
			},
		},
	}, nil)
}

func mockProviderForUnownedCluster(p *mockprovider.MockProvider, k *eksfakes.FakeKubeProvider, extraSGRules ...ec2types.SecurityGroupRule) {
	k.NewRawClientReturns(&kubernetes.RawClient{}, nil)
	k.ServerVersionReturns("1.27", nil)
	p.MockCloudFormation().On("ListStacks", mock.Anything, mock.Anything, mock.Anything).Return(&cloudformation.ListStacksOutput{
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

	vpcID := aws.String("vpc-custom")
	p.MockEC2().On("DescribeVpcs", mock.Anything, mock.Anything).Return(&ec2.DescribeVpcsOutput{
		Vpcs: []ec2types.Vpc{
			{
				CidrBlock: aws.String("192.168.0.0/19"),
				VpcId:     vpcID,
				CidrBlockAssociationSet: []ec2types.VpcCidrBlockAssociation{
					{
						CidrBlock: aws.String("192.168.0.0/19"),
					},
				},
			},
		},
	}, nil)
	p.MockEC2().On("DescribeSubnets", mock.Anything, mock.Anything, mock.Anything).Return(&ec2.DescribeSubnetsOutput{
		Subnets: []ec2types.Subnet{
			{
				SubnetId:         aws.String("subnet-custom1"),
				CidrBlock:        aws.String("192.168.160.0/19"),
				AvailabilityZone: aws.String("us-west-2a"),
				VpcId:            vpcID,
			},
			{
				SubnetId:         aws.String("subnet-custom2"),
				CidrBlock:        aws.String("192.168.96.0/19"),
				AvailabilityZone: aws.String("us-west-2b"),
				VpcId:            vpcID,
			},
		},
	}, nil)

	sgID := aws.String("sg-custom")
	p.MockEC2().On("DescribeSecurityGroupRules", mock.Anything, mock.MatchedBy(func(input *ec2.DescribeSecurityGroupRulesInput) bool {
		if len(input.Filters) != 1 {
			return false
		}
		filter := input.Filters[0]
		return *filter.Name == "group-id" && len(filter.Values) == 1 && filter.Values[0] == *sgID
	}), mock.Anything).Return(&ec2.DescribeSecurityGroupRulesOutput{
		SecurityGroupRules: append([]ec2types.SecurityGroupRule{
			{
				Description:         aws.String("Allow control plane to communicate with worker nodes in group ng-1 (kubelet and workload TCP ports"),
				FromPort:            aws.Int32(1025),
				ToPort:              aws.Int32(65535),
				GroupId:             sgID,
				IpProtocol:          aws.String("tcp"),
				IsEgress:            aws.Bool(true),
				SecurityGroupRuleId: aws.String("sgr-1"),
			},
			{
				Description:         aws.String("Allow control plane to communicate with worker nodes in group ng-1 (workload using HTTPS port, commonly used with extension API servers"),
				FromPort:            aws.Int32(443),
				ToPort:              aws.Int32(443),
				GroupId:             sgID,
				IpProtocol:          aws.String("tcp"),
				IsEgress:            aws.Bool(true),
				SecurityGroupRuleId: aws.String("sgr-2"),
			},
			{
				Description:         aws.String("Allow control plane to receive API requests from worker nodes in group ng-1"),
				FromPort:            aws.Int32(443),
				ToPort:              aws.Int32(443),
				GroupId:             sgID,
				IpProtocol:          aws.String("tcp"),
				IsEgress:            aws.Bool(false),
				SecurityGroupRuleId: aws.String("sgr-3"),
			},
			{
				Description:         aws.String("Allow control plane to communicate with worker nodes in group ng-2 (workload using HTTPS port, commonly used with extension API servers"),
				FromPort:            aws.Int32(443),
				ToPort:              aws.Int32(443),
				GroupId:             sgID,
				IpProtocol:          aws.String("tcp"),
				IsEgress:            aws.Bool(true),
				SecurityGroupRuleId: aws.String("sgr-4"),
			},
		}, extraSGRules...),
	}, nil)
}

func makeUnownedClusterConfig(clusterConfig *api.ClusterConfig) {
	clusterConfig.VPC = &api.ClusterVPC{
		SecurityGroup: "sg-custom",
		Network: api.Network{
			ID: "vpc-custom",
		},
		Subnets: &api.ClusterSubnets{
			Private: api.AZSubnetMapping{
				"us-west-2a": api.AZSubnetSpec{
					ID: "subnet-custom1",
				},
				"us-west-2b": api.AZSubnetSpec{
					ID: "subnet-custom2",
				},
			},
		},
	}
	clusterConfig.NodeGroups = append(clusterConfig.NodeGroups, &api.NodeGroup{
		NodeGroupBase: &api.NodeGroupBase{
			Name: "ng",
		},
	})
}
