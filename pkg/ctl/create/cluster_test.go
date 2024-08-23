package create

import (
	"context"
	"errors"
	"fmt"
	"time"

	k8sclient "k8s.io/client-go/kubernetes"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cftypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go-v2/service/outposts"
	outpoststypes "github.com/aws/aws-sdk-go-v2/service/outposts/types"
	"github.com/aws/smithy-go"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"k8s.io/client-go/rest"
	k8stest "k8s.io/client-go/testing"

	"github.com/stretchr/testify/mock"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	kubefake "k8s.io/client-go/kubernetes/fake"

	"github.com/weaveworks/eksctl/pkg/actions/accessentry"
	accessentryfakes "github.com/weaveworks/eksctl/pkg/actions/accessentry/fakes"
	karpenteractions "github.com/weaveworks/eksctl/pkg/actions/karpenter"
	karpenterfakes "github.com/weaveworks/eksctl/pkg/actions/karpenter/fakes"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/cfn/waiter"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils/filter"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/eks/fakes"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
)

const outpostARN = "arn:aws:outposts:us-west-2:1234:outpost/op-1234"

var _ = Describe("create cluster", func() {
	Describe("un-managed node group", func() {
		It("understands ssh access arguments correctly", func() {
			commandArgs := []string{"cluster", "--managed=false", "--ssh-access=false", "--ssh-public-key=dummy-key"}
			cmd := newMockEmptyCmd(commandArgs...)
			count := 0
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				createClusterCmdWithRunFunc(cmd, func(cmd *cmdutils.Cmd, ngFilter *filter.NodeGroupFilter, params *cmdutils.CreateClusterCmdParams) error {
					Expect(*cmd.ClusterConfig.NodeGroups[0].SSH.Allow).To(BeFalse())
					count++
					return nil
				})
			})
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
		})
		DescribeTable("create cluster successfully",
			func(args ...string) {
				commandArgs := append([]string{"cluster"}, args...)
				cmd := newMockEmptyCmd(commandArgs...)
				count := 0
				cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
					createClusterCmdWithRunFunc(cmd, func(cmd *cmdutils.Cmd, ngFilter *filter.NodeGroupFilter, params *cmdutils.CreateClusterCmdParams) error {
						Expect(cmd.ClusterConfig.Metadata.Name).NotTo(BeNil())
						count++
						return nil
					})
				})
				_, err := cmd.execute()
				Expect(err).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			},
			Entry("without cluster name", ""),
			Entry("with cluster name as flag", "--name", "clusterName"),
			Entry("with cluster name as argument", "clusterName"),
			Entry("with cluster name with hyphen as flag", "--name", "my-cluster-name-is-fine10"),
			Entry("with cluster name with hyphen as argument", "my-Cluster-name-is-fine10"),
			// vpc networking flags
			Entry("with vpc-cidr flag", "--vpc-cidr", "10.0.0.0/20"),
			Entry("with vpc-private-subnets flag", "--vpc-private-subnets", "10.0.0.0/24"),
			Entry("with vpc-public-subnets flag", "--vpc-public-subnets", "10.0.0.0/24"),
			Entry("with vpc-from-kops-cluster flag", "--vpc-from-kops-cluster", "dummy-kops-cluster"),
			Entry("with vpc-nat-mode flag", "--vpc-nat-mode", "Single"),
			// kubeconfig flags
			Entry("with write-kubeconfig flag", "--write-kubeconfig"),
			Entry("with kubeconfig flag", "--kubeconfig", "~/.kube"),
			Entry("with authenticator-role-arn flag", "--authenticator-role-arn", "arn::dummy::123/role"),
			Entry("with auto-kubeconfig flag", "--auto-kubeconfig"),
			// common node group flags
			Entry("with node-type flag", "--node-type", "m5.large"),
			Entry("with nodes flag", "--nodes", "2"),
			Entry("with nodes-min flag", "--nodes-min", "2"),
			Entry("with nodes-max flag", "--nodes-max", "2"),
			Entry("with node-volume-size flag", "--node-volume-size", "2"),
			Entry("with node-volume-type flag", "--node-volume-type", "gp2"),
			Entry("with max-pods-per-node flag", "--max-pods-per-node", "20"),
			Entry("with ssh-access flag", "--ssh-access", "true"),
			Entry("with ssh-public-key flag", "--ssh-public-key", "dummy-public-key"),
			Entry("with enable-ssm flag", "--enable-ssm"),
			Entry("with node-ami flag", "--node-ami", "ami-dummy-123"),
			Entry("with node-ami-family flag", "--node-ami-family", "AmazonLinux2"),
			Entry("with node-private-networking flag", "--node-private-networking", "true"),
			Entry("with node-security-groups flag", "--node-security-groups", "sg-123"),
			Entry("with node-labels flag", "--node-labels", "partition=backend,nodeclass=hugememory"),
			Entry("with node-zones flag", "--node-zones", "zone1,zone2,zone3", "--zones", "zone1,zone2,zone3,zone4"),
			// commons node group IAM flags
			Entry("with asg-access flag", "--asg-access", "true"),
			Entry("with external-dns-access flag", "--external-dns-access", "true"),
			Entry("with full-ecr-access flag", "--full-ecr-access", "true"),
			Entry("with appmesh-access flag", "--appmesh-access", "true"),
			Entry("with alb-ingress-access flag", "--alb-ingress-access", "true"),
			Entry("with managed flag unset", "--managed", "false"),
		)

		DescribeTable("invalid flags or arguments",
			func(c invalidParamsCase) {
				commandArgs := append([]string{"cluster", "--managed=false"}, c.args...)
				cmd := newDefaultCmd(commandArgs...)
				_, err := cmd.execute()
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(ContainSubstring(c.error)))
			},
			Entry("with cluster name as argument and flag", invalidParamsCase{
				args:  []string{"clusterName", "--name", "clusterName"},
				error: "--name=clusterName and argument clusterName cannot be used at the same time",
			}),
			Entry("with invalid flags", invalidParamsCase{
				args:  []string{"cluster", "--invalid", "dummy"},
				error: "unknown flag: --invalid",
			}),
		)
	})

	Describe("managed node group", func() {
		DescribeTable("create cluster successfully",
			func(args ...string) {
				commandArgs := append([]string{"cluster"}, args...)
				cmd := newMockEmptyCmd(commandArgs...)
				count := 0
				cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
					createClusterCmdWithRunFunc(cmd, func(cmd *cmdutils.Cmd, ngFilter *filter.NodeGroupFilter, params *cmdutils.CreateClusterCmdParams) error {
						Expect(cmd.ClusterConfig.Metadata.Name).NotTo(BeNil())
						count++
						return nil
					})
				})
				_, err := cmd.execute()
				Expect(err).To(Not(HaveOccurred()))
				Expect(count).To(Equal(1))
			},
			Entry("without cluster name", ""),
			Entry("with cluster name as flag", "--name", "clusterName"),
			Entry("with cluster name as argument", "clusterName"),
			Entry("with cluster name with hyphen as flag", "--name", "my-cluster-name-is-fine10"),
			Entry("with cluster name with hyphen as argument", "my-Cluster-name-is-fine10"),
			// vpc networking flags
			Entry("with vpc-cidr flag", "--vpc-cidr", "10.0.0.0/20"),
			Entry("with vpc-private-subnets flag", "--vpc-private-subnets", "10.0.0.0/24"),
			Entry("with vpc-public-subnets flag", "--vpc-public-subnets", "10.0.0.0/24"),
			Entry("with vpc-from-kops-cluster flag", "--vpc-from-kops-cluster", "dummy-kops-cluster"),
			Entry("with vpc-nat-mode flag", "--vpc-nat-mode", "Single"),
			// kubeconfig flags
			Entry("with write-kubeconfig flag", "--write-kubeconfig"),
			Entry("with kubeconfig flag", "--kubeconfig", "~/.kube"),
			Entry("with authenticator-role-arn flag", "--authenticator-role-arn", "arn::dummy::123/role"),
			Entry("with auto-kubeconfig flag", "--auto-kubeconfig"),
			// common node group flags
			Entry("with node-type flag", "--node-type", "m5.large"),
			Entry("with nodes flag", "--nodes", "2"),
			Entry("with nodes-min flag", "--nodes-min", "2"),
			Entry("with nodes-max flag", "--nodes-max", "2"),
			Entry("with node-volume-size flag", "--node-volume-size", "2"),
			Entry("with ssh-access flag", "--ssh-access", "true"),
			Entry("with ssh-public-key flag", "--ssh-public-key", "dummy-public-key"),
			Entry("with enable-ssm flag", "--enable-ssm"),
			Entry("with node-ami-family flag", "--node-ami-family", "AmazonLinux2"),
			Entry("with node-private-networking flag", "--node-private-networking", "true"),
			Entry("with node-labels flag", "--node-labels", "partition=backend,nodeclass=hugememory"),
			Entry("with node-zones flag", "--node-zones", "zone1,zone2,zone3", "--zones", "zone1,zone2,zone3,zone4"),
			Entry("with asg-access flag", "--asg-access", "true"),
			Entry("with external-dns-access flag", "--external-dns-access", "true"),
			Entry("with full-ecr-access flag", "--full-ecr-access", "true"),
			Entry("with appmesh-access flag", "--appmesh-access", "true"),
			Entry("with alb-ingress-access flag", "--alb-ingress-access", "true"),
		)

		DescribeTable("invalid flags or arguments",
			func(c invalidParamsCase) {
				commandArgs := append([]string{"cluster"}, c.args...)
				cmd := newDefaultCmd(commandArgs...)
				_, err := cmd.execute()
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(ContainSubstring(c.error)))
			},
			Entry("with cluster name as argument and flag", invalidParamsCase{
				args:  []string{"clusterName", "--name", "clusterName"},
				error: "--name=clusterName and argument clusterName cannot be used at the same time",
			}),
			Entry("with invalid flags", invalidParamsCase{
				args:  []string{"cluster", "--invalid", "dummy"},
				error: "unknown flag: --invalid",
			}),
			Entry("with enableSsm disabled", invalidParamsCase{
				args:  []string{"--name=test", "--enable-ssm=false"},
				error: "SSM agent is now built into EKS AMIs and cannot be disabled",
			}),
			Entry("with node zones without zones", invalidParamsCase{
				args:  []string{"--zones=zone1,zone2", "--node-zones=zone3"},
				error: "validation for --zones and --node-zones failed: node-zones [zone3] must be a subset of zones [zone1 zone2]; \"zone3\" was not found in zones",
			}),
		)
	})

	type createClusterEntry struct {
		updateClusterConfig         func(*api.ClusterConfig)
		updateClusterParams         func(*cmdutils.CreateClusterCmdParams)
		updateMocks                 func(*mockprovider.MockProvider)
		updateKubeProvider          func(*fakes.FakeKubeProvider)
		configureKarpenterInstaller func(*karpenterfakes.FakeInstallerTaskCreator)
		mockOutposts                bool
		fullyPrivateCluster         bool

		expectedErr string
	}

	DescribeTable("doCreateCluster", func(ce createClusterEntry) {
		p := mockprovider.NewMockProvider()
		defaultProviderMocks(p, defaultOutputForCluster, ce.fullyPrivateCluster, ce.mockOutposts)
		if ce.updateMocks != nil {
			ce.updateMocks(p)
		}

		// default setting for KubeProvider
		fk := &fakes.FakeKubeProvider{}
		clientset := kubefake.NewSimpleClientset()
		fk.NewStdClientSetReturns(clientset, nil)
		fk.ServerVersionReturns("1.22", nil)
		if ce.updateKubeProvider != nil {
			ce.updateKubeProvider(fk)
		}

		ctl := &eks.ClusterProvider{
			AWSProvider: p,
			Status: &eks.ProviderStatus{
				ClusterInfo: &eks.ClusterInfo{
					Cluster: testutils.NewFakeCluster(clusterName, ""),
				},
			},
			KubeProvider: fk,
		}

		fakeInstallerTaskCreator := &karpenterfakes.FakeInstallerTaskCreator{}
		if ce.configureKarpenterInstaller != nil {
			installFunc := createKarpenterInstaller
			defer func() { createKarpenterInstaller = installFunc }()
			ce.configureKarpenterInstaller(fakeInstallerTaskCreator)
		}

		clusterConfig := api.NewClusterConfig()
		clusterConfig.Metadata.Name = clusterName
		clusterConfig.AddonsConfig.DisableDefaultAddons = true
		clusterConfig.VPC.ClusterEndpoints = api.ClusterEndpointAccessDefaults()
		clusterConfig.AccessConfig.AuthenticationMode = ekstypes.AuthenticationModeApiAndConfigMap

		if ce.updateClusterConfig != nil {
			ce.updateClusterConfig(clusterConfig)
		}
		cmd := &cmdutils.Cmd{
			ClusterConfig: clusterConfig,
			ProviderConfig: api.ProviderConfig{
				WaitTimeout: time.Second * 1,
			},
		}
		filter := filter.NewNodeGroupFilter()
		params := &cmdutils.CreateClusterCmdParams{
			Subnets: map[api.SubnetTopology]*[]string{
				api.SubnetTopologyPrivate: {},
				api.SubnetTopologyPublic:  {},
			},
		}
		if ce.updateClusterParams != nil {
			ce.updateClusterParams(params)
		}
		filter.SetExcludeAll(params.WithoutNodeGroup)
		var accessEntryCreator accessentryfakes.FakeCreatorInterface
		accessEntryCreator.CreateTasksReturns(nil)
		err := doCreateCluster(cmd, filter, params, ctl, func(_ string, _ accessentry.StackCreator) accessentry.CreatorInterface {
			return &accessEntryCreator
		})
		if ce.expectedErr != "" {
			Expect(err).To(MatchError(ContainSubstring(ce.expectedErr)))
			return
		}

		Expect(err).NotTo(HaveOccurred())

		if ce.fullyPrivateCluster {
			Expect(clusterConfig.PrivateCluster.Enabled).To(BeTrue())
			p.MockEKS().AssertNumberOfCalls(GinkgoT(), "UpdateClusterConfig", 1)
		}

		if fakeInstallerTaskCreator.CreateStub != nil {
			Expect(fakeInstallerTaskCreator.CreateCallCount()).To(Equal(1))
		}
	},

		Entry("[Cluster with NodeGroups] fails to install device plugins", createClusterEntry{
			updateClusterConfig: func(c *api.ClusterConfig) {
				nodeGroup := getDefaultNodeGroup()
				nodeGroup.InstanceType = "g3.xlarge"
				c.NodeGroups = append(c.NodeGroups, nodeGroup)
			},
			updateClusterParams: func(params *cmdutils.CreateClusterCmdParams) {
				params.InstallNvidiaDevicePlugin = true
			},
			updateMocks: updateMocksForNodegroups(cftypes.StackStatusCreateComplete, defaultOutputForNodeGroup),
			updateKubeProvider: func(fk *fakes.FakeKubeProvider) {
				rawClient, err := kubernetes.NewRawClient(kubefake.NewSimpleClientset(), &rest.Config{})
				Expect(err).To(Not(HaveOccurred()))
				fk.NewRawClientReturns(rawClient, nil)
			},
			expectedErr: "failed to create cluster",
		}),

		Entry("[Cluster with NodeGroups] CloudFormation fails to create the nodegroup", createClusterEntry{
			updateClusterConfig: func(c *api.ClusterConfig) {
				c.NodeGroups = append(c.NodeGroups, getDefaultNodeGroup())
			},
			updateMocks: updateMocksForNodegroups(cftypes.StackStatusCreateFailed, defaultOutputForNodeGroup),
			expectedErr: "failed to create cluster",
		}),

		Entry("[Cluster with NodeGroups] fails to create K8s clientset", createClusterEntry{
			updateClusterConfig: func(c *api.ClusterConfig) {
				c.NodeGroups = append(c.NodeGroups, getDefaultNodeGroup())
				c.AccessConfig.AuthenticationMode = ekstypes.AuthenticationModeConfigMap
			},
			updateMocks: updateMocksForNodegroups(cftypes.StackStatusCreateComplete, defaultOutputForNodeGroup),
			updateKubeProvider: func(fk *fakes.FakeKubeProvider) {
				fk.NewStdClientSetReturns(nil, errors.New("failed to create clientset"))
			},
			expectedErr: "failed to create clientset",
		}),

		Entry("[Cluster with nodegroups] fails when bootstrapClusterCreatorAdminPermissions is false and authenticationMode is CONFIG_MAP", createClusterEntry{
			updateClusterConfig: func(c *api.ClusterConfig) {
				c.NodeGroups = append(c.NodeGroups, getDefaultNodeGroup())
				c.AccessConfig = &api.AccessConfig{
					AuthenticationMode:                      ekstypes.AuthenticationModeConfigMap,
					BootstrapClusterCreatorAdminPermissions: api.Disabled(),
				}
			},
			updateMocks: updateMocksForNodegroups(cftypes.StackStatusCreateComplete, defaultOutputForNodeGroup),
			expectedErr: "cannot create self-managed nodegroups when authenticationMode is CONFIG_MAP and bootstrapClusterCreatorAdminPermissions is false",
		}),

		Entry("[Cluster with nodegroups] fails when bootstrapClusterCreatorAdminPermissions is false and no access entries are configured", createClusterEntry{
			updateClusterConfig: func(c *api.ClusterConfig) {
				c.NodeGroups = append(c.NodeGroups, getDefaultNodeGroup())
				c.AccessConfig = &api.AccessConfig{
					AuthenticationMode:                      ekstypes.AuthenticationModeApiAndConfigMap,
					BootstrapClusterCreatorAdminPermissions: api.Disabled(),
				}
			},
			updateMocks: updateMocksForNodegroups(cftypes.StackStatusCreateComplete, defaultOutputForNodeGroup),
			updateKubeProvider: func(fk *fakes.FakeKubeProvider) {
				fk.NewStdClientSetReturns(nil, errors.New("failed to create clientset"))
			},
			expectedErr: "cannot create self-managed nodegroups when bootstrapClusterCreatorAdminPermissions is false and no access entries are configured",
		}),

		Entry("[Cluster with nodegroups] skips error if it fails to create Clientset and cluster uses access entries", createClusterEntry{
			updateClusterConfig: func(c *api.ClusterConfig) {
				c.NodeGroups = append(c.NodeGroups, getDefaultNodeGroup())
				c.AccessConfig = &api.AccessConfig{
					AuthenticationMode:                      ekstypes.AuthenticationModeApiAndConfigMap,
					BootstrapClusterCreatorAdminPermissions: api.Disabled(),
					AccessEntries: []api.AccessEntry{
						{
							PrincipalARN: api.MustParseARN("arn:aws:iam::111122223333:role/role-1"),
						},
					},
				}
			},
			updateMocks: updateMocksForNodegroups(cftypes.StackStatusCreateComplete, defaultOutputForNodeGroup),
			updateKubeProvider: func(fk *fakes.FakeKubeProvider) {
				fk.NewStdClientSetReturns(nil, errors.New("failed to create clientset"))
			},
		}),

		Entry("[Cluster with NodeGroups] times out waiting for nodes to join the cluster", createClusterEntry{
			updateClusterConfig: func(c *api.ClusterConfig) {
				c.NodeGroups = append(c.NodeGroups, getDefaultNodeGroup())
			},
			updateMocks: updateMocksForNodegroups(cftypes.StackStatusCreateComplete, defaultOutputForNodeGroup),
			updateKubeProvider: func(fk *fakes.FakeKubeProvider) {
				node := getDefaultNode()
				clientset := kubefake.NewSimpleClientset(node)
				fk.NewStdClientSetReturns(clientset, nil)
			},
			expectedErr: "timed out waiting for at least 1 nodes to join the cluster and become ready",
		}),

		Entry("[Cluster with nodegroups] does not wait for nodes to join the cluster if cluster uses access entries and Clientset creation fails", createClusterEntry{
			updateClusterConfig: func(c *api.ClusterConfig) {
				c.NodeGroups = append(c.NodeGroups, getDefaultNodeGroup())
				c.AccessConfig = &api.AccessConfig{
					AuthenticationMode:                      ekstypes.AuthenticationModeApiAndConfigMap,
					BootstrapClusterCreatorAdminPermissions: api.Disabled(),
					AccessEntries: []api.AccessEntry{
						{
							PrincipalARN: api.MustParseARN("arn:aws:iam::111122223333:role/role-1"),
						},
					},
				}
			},
			updateMocks: updateMocksForNodegroups(cftypes.StackStatusCreateComplete, defaultOutputForNodeGroup),
			updateKubeProvider: func(fk *fakes.FakeKubeProvider) {
				fk.NewStdClientSetReturns(nil, errors.New("failed to create clientset"))
			},
		}),

		Entry("[Cluster with NodeGroups] all resources are created successfully", createClusterEntry{
			updateClusterConfig: func(c *api.ClusterConfig) {
				c.NodeGroups = append(c.NodeGroups, getDefaultNodeGroup())
			},
			updateMocks: updateMocksForNodegroups(cftypes.StackStatusCreateComplete, defaultOutputForNodeGroup),
			updateKubeProvider: func(fk *fakes.FakeKubeProvider) {
				node := getDefaultNode()
				clientset := kubefake.NewSimpleClientset(node)
				watcher := watch.NewFake()
				go func() {
					defer watcher.Stop()
					watcher.Add(node)
				}()
				clientset.PrependWatchReactor("nodes", k8stest.DefaultWatchReactor(watcher, nil))
				fk.NewStdClientSetReturns(clientset, nil)
			},
		}),

		Entry("nodegroup with an instance role ARN", createClusterEntry{
			updateClusterConfig: func(c *api.ClusterConfig) {
				ng := getDefaultNodeGroup()
				ng.IAM.InstanceRoleARN = "role-1"
				c.NodeGroups = []*api.NodeGroup{ng}
			},
			updateMocks: func(mp *mockprovider.MockProvider) {
				updateMocksForNodegroups(cftypes.StackStatusCreateComplete, defaultOutputForNodeGroup)(mp)
				mp.MockEKS().On("DescribeAccessEntry", mock.Anything, &awseks.DescribeAccessEntryInput{
					PrincipalArn: aws.String("role-1"),
					ClusterName:  aws.String(clusterName),
				}).Return(nil, &ekstypes.ResourceNotFoundException{ClusterName: aws.String(clusterName)}).Once()

				mp.MockEKS().On("CreateAccessEntry", mock.Anything, &awseks.CreateAccessEntryInput{
					PrincipalArn: aws.String("role-1"),
					ClusterName:  aws.String(clusterName),
					Type:         aws.String("EC2_LINUX"),
					Tags: map[string]string{
						api.ClusterNameLabel: clusterName,
					},
				}).Return(&awseks.CreateAccessEntryOutput{}, nil).Once()
			},
			updateKubeProvider: func(fk *fakes.FakeKubeProvider) {
				node := getDefaultNode()
				clientset := kubefake.NewSimpleClientset(node)
				watcher := watch.NewFake()
				go func() {
					defer watcher.Stop()
					watcher.Add(node)
				}()
				clientset.PrependWatchReactor("nodes", k8stest.DefaultWatchReactor(watcher, nil))
				fk.NewStdClientSetReturns(clientset, nil)
			},
		}),

		Entry("standard cluster", createClusterEntry{}),

		Entry("[Cluster with Karpenter] installs Karpenter successfully", createClusterEntry{
			updateClusterConfig: func(c *api.ClusterConfig) {
				c.Karpenter = &api.Karpenter{
					Version: "v0.18.0",
				}
			},
			configureKarpenterInstaller: func(ki *karpenterfakes.FakeInstallerTaskCreator) {
				ki.CreateStub = func(ctx context.Context) error {
					return nil
				}
				createKarpenterInstaller = func(ctx context.Context, cfg *api.ClusterConfig, ctl *eks.ClusterProvider, stackManager manager.StackManager, clientSet kubernetes.Interface, restClientGetter *kubernetes.SimpleRESTClientGetter) (karpenteractions.InstallerTaskCreator, error) {
					return ki, nil
				}
			},
		}),

		Entry("[Cluster with Karpenter] fails to create the installer", createClusterEntry{
			updateClusterConfig: func(c *api.ClusterConfig) {
				c.Karpenter = &api.Karpenter{
					Version: "v0.18.0",
				}
			},
			configureKarpenterInstaller: func(ki *karpenterfakes.FakeInstallerTaskCreator) {
				createKarpenterInstaller = func(ctx context.Context, cfg *api.ClusterConfig, ctl *eks.ClusterProvider, stackManager manager.StackManager, clientSet kubernetes.Interface, restClientGetter *kubernetes.SimpleRESTClientGetter) (karpenteractions.InstallerTaskCreator, error) {
					return ki, fmt.Errorf("failed to create karpenter installer")
				}
			},
			expectedErr: "failed to create installer",
		}),

		Entry("[Cluster with Karpenter] fails to actually install Karpenter", createClusterEntry{
			updateClusterConfig: func(c *api.ClusterConfig) {
				c.Karpenter = &api.Karpenter{
					Version: "v0.18.0",
				}
			},
			configureKarpenterInstaller: func(ki *karpenterfakes.FakeInstallerTaskCreator) {
				ki.CreateStub = func(ctx context.Context) error {
					return fmt.Errorf("failed to install karpenter")
				}
				createKarpenterInstaller = func(ctx context.Context, cfg *api.ClusterConfig, ctl *eks.ClusterProvider, stackManager manager.StackManager, clientSet kubernetes.Interface, restClientGetter *kubernetes.SimpleRESTClientGetter) (karpenteractions.InstallerTaskCreator, error) {
					return ki, nil
				}
			},
			expectedErr: "failed to install Karpenter",
		}),

		Entry("[Cluster with Karpenter] fails to install Karpenter if Clientset creation fails", createClusterEntry{
			updateClusterConfig: func(c *api.ClusterConfig) {
				c.Karpenter = &api.Karpenter{
					Version: "v0.18.0",
				}
				c.AccessConfig = &api.AccessConfig{
					AuthenticationMode:                      ekstypes.AuthenticationModeApiAndConfigMap,
					BootstrapClusterCreatorAdminPermissions: api.Disabled(),
					AccessEntries: []api.AccessEntry{
						{
							PrincipalARN: api.MustParseARN("arn:aws:iam::111122223333:role/role-1"),
						},
					},
				}
			},
			configureKarpenterInstaller: func(ki *karpenterfakes.FakeInstallerTaskCreator) {
				ki.CreateStub = func(ctx context.Context) error {
					return nil
				}
				createKarpenterInstaller = func(ctx context.Context, cfg *api.ClusterConfig, ctl *eks.ClusterProvider, stackManager manager.StackManager, clientSet kubernetes.Interface, restClientGetter *kubernetes.SimpleRESTClientGetter) (karpenteractions.InstallerTaskCreator, error) {
					return ki, nil
				}
			},
			updateKubeProvider: func(fk *fakes.FakeKubeProvider) {
				fk.NewStdClientSetStub = func(info kubeconfig.ClusterInfo) (k8sclient.Interface, error) {
					return nil, errors.New("error installing Karpenter: failed to create clientset")
				}
			},
			expectedErr: "error installing Karpenter: failed to create clientset",
		}),

		Entry("[Fully Private Cluster] updates cluster config successfully", createClusterEntry{
			updateClusterConfig: func(c *api.ClusterConfig) {
				c.VPC.ClusterEndpoints.PrivateAccess = api.Enabled()
				c.PrivateCluster = &api.PrivateCluster{
					Enabled:              true,
					SkipEndpointCreation: true,
				}
			},
			updateMocks: func(provider *mockprovider.MockProvider) {
				provider.MockEKS().On("UpdateClusterConfig", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					Expect(args).To(HaveLen(2))
					Expect(args[1]).To(BeAssignableToTypeOf(&awseks.UpdateClusterConfigInput{}))
					Expect(args[1].(*awseks.UpdateClusterConfigInput).ResourcesVpcConfig.EndpointPrivateAccess).To(Equal(api.Enabled()))
					Expect(args[1].(*awseks.UpdateClusterConfigInput).ResourcesVpcConfig.EndpointPublicAccess).To(Equal(api.Disabled()))
				}).Return(&awseks.UpdateClusterConfigOutput{
					Update: &ekstypes.Update{
						Id: aws.String("test-id"),
					},
				}, nil).Once()

				provider.MockEKS().On("DescribeUpdate", mock.Anything, mock.MatchedBy(func(input *awseks.DescribeUpdateInput) bool {
					return true
				}), mock.Anything).Return(&awseks.DescribeUpdateOutput{
					Update: &ekstypes.Update{
						Id:     aws.String("test-id"),
						Type:   ekstypes.UpdateTypeConfigUpdate,
						Status: ekstypes.UpdateStatusSuccessful,
					},
				}, nil).Once()
			},
			fullyPrivateCluster: true,
		}),

		Entry("[Fully Private Cluster] fails to update cluster config", createClusterEntry{
			updateClusterConfig: func(c *api.ClusterConfig) {
				c.VPC.ClusterEndpoints.PrivateAccess = api.Enabled()
				c.PrivateCluster = &api.PrivateCluster{
					Enabled:              true,
					SkipEndpointCreation: true,
				}
			},
			updateMocks: func(provider *mockprovider.MockProvider) {
				provider.MockEKS().On("UpdateClusterConfig", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					Expect(args).To(HaveLen(2))
					Expect(args[1]).To(BeAssignableToTypeOf(&awseks.UpdateClusterConfigInput{}))
					Expect(args[1].(*awseks.UpdateClusterConfigInput).ResourcesVpcConfig.EndpointPrivateAccess).To(Equal(api.Enabled()))
					Expect(args[1].(*awseks.UpdateClusterConfigInput).ResourcesVpcConfig.EndpointPublicAccess).To(Equal(api.Disabled()))
				}).Return(nil, fmt.Errorf("")).Once()
			},
			expectedErr:         "error disabling public endpoint access for the cluster",
			fullyPrivateCluster: true,
		}),

		Entry("[Outposts] control plane on Outposts with valid config", createClusterEntry{
			updateClusterConfig: func(c *api.ClusterConfig) {
				c.Outpost = &api.Outpost{
					ControlPlaneOutpostARN: outpostARN,
				}
			},
			mockOutposts: true,
		}),

		Entry("[Outposts] unavailable instance type specified for the control plane", createClusterEntry{
			updateClusterConfig: func(c *api.ClusterConfig) {
				c.Outpost = &api.Outpost{
					ControlPlaneOutpostARN:   outpostARN,
					ControlPlaneInstanceType: "t2.medium",
				}
			},
			mockOutposts: true,

			expectedErr: fmt.Sprintf(`instance type "t2.medium" does not exist in Outpost %q`, outpostARN),
		}),

		Entry("[Outposts] available instance type specified for the control plane", createClusterEntry{
			updateClusterConfig: func(c *api.ClusterConfig) {
				c.Outpost = &api.Outpost{
					ControlPlaneOutpostARN:   outpostARN,
					ControlPlaneInstanceType: "m5.xlarge",
				}
			},
			mockOutposts: true,
		}),

		Entry("[Outposts] nodegroups specified when the VPC will be created by eksctl", createClusterEntry{
			updateClusterConfig: func(c *api.ClusterConfig) {
				c.Outpost = &api.Outpost{
					ControlPlaneOutpostARN: outpostARN,
				}
				c.NodeGroups = []*api.NodeGroup{
					api.NewNodeGroup(),
				}
			},
			mockOutposts: true,

			expectedErr: "cannot create nodegroups on Outposts when the VPC is created by eksctl as it will not have connectivity to the API server; please rerun the command with `--without-nodegroup` and run `eksctl create nodegroup` after associating the VPC with a local gateway and ensuring connectivity to the API server",
		}),

		Entry("[Outposts] nodegroups specified on Outposts but the control plane is not on Outposts", createClusterEntry{
			updateClusterConfig: func(c *api.ClusterConfig) {
				ng := api.NewNodeGroup()
				ng.OutpostARN = outpostARN
				c.NodeGroups = []*api.NodeGroup{ng}
			},
			mockOutposts: true,

			expectedErr: "creating nodegroups on Outposts when the control plane is not on Outposts is not supported during cluster creation; " +
				"either create the nodegroups after cluster creation or consider creating the control plane on Outposts",
		}),

		Entry("[Outposts] nodegroups specified when the VPC will be created by eksctl, but with --without-nodegroup", createClusterEntry{
			updateClusterConfig: func(c *api.ClusterConfig) {
				c.Outpost = &api.Outpost{
					ControlPlaneOutpostARN: outpostARN,
				}
				c.NodeGroups = []*api.NodeGroup{
					api.NewNodeGroup(),
				}
			},
			updateClusterParams: func(params *cmdutils.CreateClusterCmdParams) {
				params.WithoutNodeGroup = true
			},
			mockOutposts: true,
		}),

		Entry("[Outposts] specified Outpost does not exist", createClusterEntry{
			updateClusterConfig: func(c *api.ClusterConfig) {
				c.Outpost = &api.Outpost{
					ControlPlaneOutpostARN: outpostARN,
				}
			},
			updateMocks: func(provider *mockprovider.MockProvider) {
				provider.MockOutposts().On("GetOutpost", mock.Anything, &outposts.GetOutpostInput{
					OutpostId: aws.String(outpostARN),
				}).Return(nil, &outpoststypes.NotFoundException{Message: aws.String("not found")})
			},
			expectedErr: "error getting Outpost details: NotFoundException: not found",
		}),

		Entry("[Outposts] specified Outpost placement group does not exist", createClusterEntry{
			updateClusterConfig: func(c *api.ClusterConfig) {
				c.Outpost = &api.Outpost{
					ControlPlaneOutpostARN: outpostARN,
					ControlPlanePlacement: &api.Placement{
						GroupName: "test",
					},
				}
			},
			updateMocks: func(provider *mockprovider.MockProvider) {
				provider.MockEC2().On("DescribePlacementGroups", mock.Anything, &ec2.DescribePlacementGroupsInput{
					GroupNames: []string{"test"},
				}).Return(nil, &smithy.OperationError{
					OperationName: "DescribePlacementGroups",
					Err:           errors.New("api error InvalidPlacementGroup.Unknown: The Placement Group 'test' is unknown"),
				})
			},
			mockOutposts: true,
			expectedErr:  `placement group "test" does not exist`,
		}),

		Entry("[Outposts] valid Outpost placement group specified", createClusterEntry{
			updateClusterConfig: func(c *api.ClusterConfig) {
				c.Outpost = &api.Outpost{
					ControlPlaneOutpostARN: outpostARN,
					ControlPlanePlacement: &api.Placement{
						GroupName: "test",
					},
				}
			},
			updateMocks: func(provider *mockprovider.MockProvider) {
				provider.MockEC2().On("DescribePlacementGroups", mock.Anything, &ec2.DescribePlacementGroupsInput{
					GroupNames: []string{"test"},
				}).Return(&ec2.DescribePlacementGroupsOutput{
					PlacementGroups: []ec2types.PlacementGroup{
						{
							GroupName: aws.String("test"),
						},
					},
				}, nil)
			},
			mockOutposts: true,
		}),
	)
})

var (
	clusterName         = "my-cluster"
	clusterStackName    = "eksctl-" + clusterName + "-cluster"
	nodeGroupName       = "my-nodegroup"
	nodeGroupStackName  = "eksctl-" + clusterName + "-nodegroup-" + nodeGroupName
	nodeInstanceRoleARN = "arn:aws:iam::083751696308:role/eksctl-my-cluster-cluster-nodegroup-my-nodegroup-NodeInstanceRole-1IYQ3JS8OKPX1"

	defaultOutputForNodeGroup = []cftypes.Output{
		{
			OutputKey:   aws.String(outputs.NodeGroupInstanceRoleARN),
			OutputValue: aws.String(nodeInstanceRoleARN),
		},
		{
			OutputKey:   aws.String(outputs.NodeGroupInstanceProfileARN),
			OutputValue: aws.String("arn:aws:iam::083751696308:role/eksctl-my-cluster-cluster-nodegroup-my-nodegroup-NodeInstanceProfile-1IYQ3JS8OKPX1"),
		},
		{
			OutputKey:   aws.String(outputs.NodeGroupFeaturePrivateNetworking),
			OutputValue: aws.String("ngfpn"),
		},
		{
			OutputKey:   aws.String(outputs.NodeGroupFeatureLocalSecurityGroup),
			OutputValue: aws.String("nglsg"),
		},
		{
			OutputKey:   aws.String(outputs.NodeGroupFeatureSharedSecurityGroup),
			OutputValue: aws.String("ngssg"),
		},
		{
			OutputKey:   aws.String(outputs.NodeGroupUsesAccessEntry),
			OutputValue: aws.String("true"),
		},
	}

	defaultOutputForCluster = []cftypes.Output{
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
			OutputKey:   aws.String("FeatureNATMode"),
			OutputValue: aws.String("Single"),
		},
		{
			OutputKey:   aws.String("SubnetsPrivate"),
			OutputValue: aws.String("sub-priv-1 sub-priv-2 sub-priv-3"),
		},
		{
			OutputKey:   aws.String("SubnetsPublic"),
			OutputValue: aws.String("sub-pub-1 sub-pub-2 sub-pub-3"),
		},
		{
			OutputKey:   aws.String("ServiceRoleARN"),
			OutputValue: aws.String("arn:aws:iam::123456:role/amazingrole-1"),
		},
		{
			OutputKey:   aws.String("ARN"),
			OutputValue: aws.String("arn:aws:iam::123456:role/amazingrole-2"),
		},
		{
			OutputKey:   aws.String("CertificateAuthorityData"),
			OutputValue: aws.String("dGVzdAo="),
		},
		{
			OutputKey:   aws.String("ClusterStackName"),
			OutputValue: aws.String(clusterStackName),
		},
		{
			OutputKey:   aws.String("Endpoint"),
			OutputValue: aws.String("https://endpoint.com"),
		},
	}

	getDefaultNode = func() *corev1.Node {
		return &corev1.Node{
			ObjectMeta: v1.ObjectMeta{
				Name: nodeGroupName + "-my-node",
				Labels: map[string]string{
					"alpha.eksctl.io/nodegroup-name": nodeGroupName,
				},
			},
			Status: corev1.NodeStatus{
				Conditions: []corev1.NodeCondition{
					{
						Type:   corev1.NodeReady,
						Status: corev1.ConditionTrue,
					},
				},
			}}
	}

	getDefaultNodeGroup = func() *api.NodeGroup {
		return &api.NodeGroup{
			NodeGroupBase: &api.NodeGroupBase{
				Name:             nodeGroupName,
				AMIFamily:        api.NodeImageFamilyAmazonLinux2,
				AMI:              "ami-123",
				SSH:              &api.NodeGroupSSH{Allow: api.Disabled()},
				InstanceSelector: &api.InstanceSelector{},
				SecurityGroups: &api.NodeGroupSGs{
					WithShared: aws.Bool(false),
					WithLocal:  aws.Bool(false),
				},
				ScalingConfig: &api.ScalingConfig{
					DesiredCapacity: aws.Int(1),
				},
				IAM: &api.NodeGroupIAM{},
			},
		}
	}

	updateMocksForNodegroups = func(status cftypes.StackStatus, outputs []cftypes.Output) func(mp *mockprovider.MockProvider) {
		return func(mp *mockprovider.MockProvider) {
			mp.MockEC2().On("DescribeInstanceTypeOfferings", mock.Anything, mock.Anything, mock.Anything).Return(&ec2.DescribeInstanceTypeOfferingsOutput{
				InstanceTypeOfferings: []ec2types.InstanceTypeOffering{
					{
						InstanceType: "g3.xlarge",
						Location:     aws.String("us-west-2-1b"),
						LocationType: "availability-zone",
					},
					{
						InstanceType: "g3.xlarge",
						Location:     aws.String("us-west-2-1a"),
						LocationType: "availability-zone",
					},
				},
			}, nil)
			mp.MockCloudFormation().On("CreateStack", mock.Anything, mock.Anything).Return(&cloudformation.CreateStackOutput{
				StackId: aws.String(nodeGroupStackName),
			}, nil).Once()
			// mock for when DescribeStacks is called inside DoWaitUntilStackIsCreated
			mp.MockCloudFormation().On("DescribeStacks", mock.Anything, mock.Anything, mock.Anything).Return(&cloudformation.DescribeStacksOutput{
				Stacks: []cftypes.Stack{
					{
						StackName:   aws.String(nodeGroupStackName),
						StackStatus: status,
					},
				},
			}, nil).Once()
			mp.MockCloudFormation().On("DescribeStacks", mock.Anything, mock.Anything).Return(&cloudformation.DescribeStacksOutput{
				Stacks: []cftypes.Stack{
					{
						StackName:   aws.String(nodeGroupStackName),
						StackStatus: status,
						Tags: []cftypes.Tag{
							{
								Key:   aws.String(api.NodeGroupNameTag),
								Value: aws.String(nodeGroupStackName),
							},
							{
								Key:   aws.String(api.ClusterNameTag),
								Value: aws.String(clusterStackName),
							},
						},
						Outputs: outputs,
					},
				},
			}, nil).Twice()
		}
	}
)

func defaultProviderMocks(p *mockprovider.MockProvider, output []cftypes.Output, fullyPrivateCluster, controlPlaneOnOutposts bool) {
	p.MockEC2().On("DescribeAvailabilityZones", mock.Anything, &ec2.DescribeAvailabilityZonesInput{
		Filters: []ec2types.Filter{{
			Name:   aws.String("region-name"),
			Values: []string{"us-west-2"},
		}, {
			Name:   aws.String("state"),
			Values: []string{string(ec2types.AvailabilityZoneStateAvailable)},
		}, {
			Name:   aws.String("zone-type"),
			Values: []string{string(ec2types.LocationTypeAvailabilityZone)},
		}},
	}).Return(&ec2.DescribeAvailabilityZonesOutput{
		AvailabilityZones: []ec2types.AvailabilityZone{
			{
				GroupName: aws.String("name"),
				ZoneName:  aws.String("us-west-2-1b"),
				ZoneId:    aws.String("id"),
			},
			{
				GroupName: aws.String("name"),
				ZoneName:  aws.String("us-west-2-1a"),
				ZoneId:    aws.String("id"),
			}},
	}, nil)
	p.MockCloudFormation().On("ListStacks", mock.Anything, mock.Anything, mock.Anything).Return(&cloudformation.ListStacksOutput{
		StackSummaries: []cftypes.StackSummary{
			{
				StackName:   aws.String(clusterStackName),
				StackStatus: "CREATE_COMPLETE",
			},
		},
	}, nil)

	if fullyPrivateCluster {
		output = append(output, cftypes.Output{
			OutputKey:   aws.String("ClusterFullyPrivate"),
			OutputValue: aws.String(""),
		})
	}

	p.MockCloudFormation().On("DescribeStacks", mock.Anything, mock.Anything).Return(&cloudformation.DescribeStacksOutput{
		Stacks: []cftypes.Stack{
			{
				StackName:   aws.String(clusterStackName),
				StackStatus: "CREATE_COMPLETE",
				Tags: []cftypes.Tag{
					{
						Key:   aws.String(api.ClusterNameTag),
						Value: aws.String(clusterStackName),
					},
				},
				Outputs: output,
			},
		},
	}, nil).Once()

	var outpostConfig *ekstypes.OutpostConfigResponse
	if controlPlaneOnOutposts {
		mockOutposts(p, outpostARN)
		outpostConfig = &ekstypes.OutpostConfigResponse{
			OutpostArns:              []string{outpostARN},
			ControlPlaneInstanceType: aws.String("m5.xlarge"),
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
			Name:                    aws.String(clusterName),
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
				api.ClusterNameTag: clusterStackName,
			},
			Version:       aws.String("1.22"),
			OutpostConfig: outpostConfig,
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
	p.MockEC2().On("DescribeSubnets", mock.Anything, mock.Anything, mock.Anything).Return(&ec2.DescribeSubnetsOutput{
		Subnets: []ec2types.Subnet{},
	}, nil)
	p.MockEC2().On("DescribeVpcs", mock.Anything, mock.Anything).Return(&ec2.DescribeVpcsOutput{
		Vpcs: []ec2types.Vpc{
			{
				VpcId:     aws.String("vpc-1"),
				CidrBlock: aws.String("192.168.0.0/16"),
			},
		},
	}, nil)

	waiter.ClusterCreationNextDelay = func(_ int) time.Duration {
		return 0
	}
	p.MockCloudFormation().On("CreateStack", mock.Anything, mock.Anything).Return(&cloudformation.CreateStackOutput{
		StackId: aws.String(clusterStackName),
	}, nil).Once()

	p.MockCredentialsProvider().On("Retrieve", mock.Anything).Return(aws.Credentials{
		AccessKeyID:     "key-id",
		SecretAccessKey: "secret-access-key",
		SessionToken:    "token",
	}, nil)
}

func mockOutposts(provider *mockprovider.MockProvider, outpostID string) {
	provider.MockOutposts().On("GetOutpost", mock.Anything, &outposts.GetOutpostInput{
		OutpostId: aws.String(outpostID),
	}).Return(&outposts.GetOutpostOutput{
		Outpost: &outpoststypes.Outpost{
			AvailabilityZone: aws.String("us-west-2a"),
		},
	}, nil)
	provider.MockOutposts().On("GetOutpostInstanceTypes", mock.Anything, &outposts.GetOutpostInstanceTypesInput{
		OutpostId: aws.String(outpostID),
	}, mock.Anything).Return(&outposts.GetOutpostInstanceTypesOutput{
		InstanceTypes: []outpoststypes.InstanceTypeItem{
			{
				InstanceType: aws.String("m5.xlarge"),
			},
		},
	}, nil)
	provider.MockEC2().On("DescribeInstanceTypes", mock.Anything, &ec2.DescribeInstanceTypesInput{
		InstanceTypes: []ec2types.InstanceType{"m5.xlarge"},
	}, mock.Anything).Return(&ec2.DescribeInstanceTypesOutput{
		InstanceTypes: []ec2types.InstanceTypeInfo{
			{
				InstanceType: "m5.xlarge",
				VCpuInfo: &ec2types.VCpuInfo{
					DefaultVCpus:          aws.Int32(4),
					DefaultCores:          aws.Int32(2),
					DefaultThreadsPerCore: aws.Int32(2),
				},
				MemoryInfo: &ec2types.MemoryInfo{
					SizeInMiB: aws.Int64(16384),
				},
			},
		},
	}, nil)
}
