package nodegroup_test

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	utilFakes "github.com/weaveworks/eksctl/pkg/ctl/cmdutils/filter/fakes"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/eks/fakes"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
	utilsstrings "github.com/weaveworks/eksctl/pkg/utils/strings"
)

type ngEntry struct {
	version       string
	opts          nodegroup.CreateOpts
	mockCalls     func(*fakes.FakeKubeProvider, *fakes.FakeNodeGroupInitialiser, *utilFakes.FakeNodegroupFilter)
	expectedCalls func(*fakes.FakeKubeProvider, *fakes.FakeNodeGroupInitialiser, *utilFakes.FakeNodegroupFilter)
	expErr        error
}

var _ = DescribeTable("Create", func(t ngEntry) {
	cfg := newClusterConfig()
	cfg.Metadata.Version = t.version

	p := mockprovider.NewMockProvider()
	ctl := &eks.ClusterProvider{
		Provider: p,
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

	ngFilter := utilFakes.FakeNodegroupFilter{}
	if t.mockCalls != nil {
		t.mockCalls(k, init, &ngFilter)
	}

	err := m.Create(t.opts, &ngFilter)

	if t.expErr != nil {
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(ContainSubstring(t.expErr.Error())))
	} else {
		Expect(err).NotTo(HaveOccurred())
	}
	if t.expectedCalls != nil {
		t.expectedCalls(k, init, &ngFilter)
	}
},
	Entry("fails when cluster version is not supported", ngEntry{
		version: "1.14",
		expErr:  fmt.Errorf("invalid version, %s is no longer supported, supported values: auto, default, latest, %s\nsee also: https://docs.aws.amazon.com/eks/latest/userguide/kubernetes-versions.html", "1.14", strings.Join(api.SupportedVersions(), ", ")),
	}),

	Entry("fails when it does not support ARM", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
			k.NewRawClientReturns(nil, fmt.Errorf("err"))
		},
		expErr: fmt.Errorf("err"),
	}),

	Entry("when cluster is unowned, fails to load VPC from config if config is not supplied", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
			k.NewRawClientReturns(&kubernetes.RawClient{}, nil)
			k.ServerVersionReturns("1.17", nil)
			k.LoadClusterIntoSpecFromStackReturns(&manager.StackNotFoundErr{})
		},
		expErr: errors.Wrapf(errors.New("VPC configuration required for creating nodegroups on clusters not owned by eksctl: vpc.subnets, vpc.id, vpc.securityGroup"), "loading VPC spec for cluster %q", "my-cluster"),
	}),

	Entry("fails when cluster does not support managed nodes", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
			k.SupportsManagedNodesReturns(false, errors.New("err"))
		},
		expectedCalls: func(k *fakes.FakeKubeProvider, _ *fakes.FakeNodeGroupInitialiser, _ *utilFakes.FakeNodegroupFilter) {
			Expect(k.NewRawClientCallCount()).To(Equal(1))
			Expect(k.ServerVersionCallCount()).To(Equal(1))
			Expect(k.LoadClusterIntoSpecFromStackCallCount()).To(Equal(1))
			Expect(k.SupportsManagedNodesCallCount()).To(Equal(1))
		},
		expErr: errors.New("err"),
	}),

	Entry("fails to set instance types to instances matched by instance selector criteria", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
			k.SupportsManagedNodesReturns(true, nil)
			init.ExpandInstanceSelectorOptionsReturns(errors.New("err"))
		},
		expectedCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, _ *utilFakes.FakeNodegroupFilter) {
			Expect(k.NewRawClientCallCount()).To(Equal(1))
			Expect(k.ServerVersionCallCount()).To(Equal(1))
			Expect(k.LoadClusterIntoSpecFromStackCallCount()).To(Equal(1))
			Expect(k.SupportsManagedNodesCallCount()).To(Equal(1))
			Expect(init.NewAWSSelectorSessionCallCount()).To(Equal(1))
			Expect(init.ExpandInstanceSelectorOptionsCallCount()).To(Equal(1))
		},
		expErr: errors.New("err"),
	}),

	Entry("fails when cluster is not compatible with ng config", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
			k.SupportsManagedNodesReturns(true, nil)
			k.ValidateClusterForCompatibilityReturns(errors.New("err"))
		},
		expectedCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, _ *utilFakes.FakeNodegroupFilter) {
			Expect(k.NewRawClientCallCount()).To(Equal(1))
			Expect(k.ServerVersionCallCount()).To(Equal(1))
			Expect(k.LoadClusterIntoSpecFromStackCallCount()).To(Equal(1))
			Expect(k.SupportsManagedNodesCallCount()).To(Equal(1))
			Expect(init.NewAWSSelectorSessionCallCount()).To(Equal(1))
			Expect(init.ExpandInstanceSelectorOptionsCallCount()).To(Equal(1))
			Expect(k.ValidateClusterForCompatibilityCallCount()).To(Equal(1))
		},
		expErr: errors.Wrap(errors.New("err"), "cluster compatibility check failed"),
	}),

	Entry("fails when it cannot validate legacy subnets for ng", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
			k.SupportsManagedNodesReturns(true, nil)
			init.ValidateLegacySubnetsForNodeGroupsReturns(errors.New("err"))
		},
		expectedCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, _ *utilFakes.FakeNodegroupFilter) {
			Expect(k.NewRawClientCallCount()).To(Equal(1))
			Expect(k.ServerVersionCallCount()).To(Equal(1))
			Expect(k.LoadClusterIntoSpecFromStackCallCount()).To(Equal(1))
			Expect(k.SupportsManagedNodesCallCount()).To(Equal(1))
			Expect(init.NewAWSSelectorSessionCallCount()).To(Equal(1))
			Expect(init.ExpandInstanceSelectorOptionsCallCount()).To(Equal(1))
			Expect(k.ValidateClusterForCompatibilityCallCount()).To(Equal(1))
			Expect(init.ValidateLegacySubnetsForNodeGroupsCallCount()).To(Equal(1))
		},
		expErr: errors.New("err"),
	}),

	Entry("fails when existing local ng stacks in config file is not listed", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
			k.SupportsManagedNodesReturns(true, nil)
			f.SetOnlyLocalReturns(errors.New("err"))
		},
		expectedCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
			Expect(k.NewRawClientCallCount()).To(Equal(1))
			Expect(k.ServerVersionCallCount()).To(Equal(1))
			Expect(k.LoadClusterIntoSpecFromStackCallCount()).To(Equal(1))
			Expect(k.SupportsManagedNodesCallCount()).To(Equal(1))
			Expect(init.NewAWSSelectorSessionCallCount()).To(Equal(1))
			Expect(init.ExpandInstanceSelectorOptionsCallCount()).To(Equal(1))
			Expect(k.ValidateClusterForCompatibilityCallCount()).To(Equal(1))
			Expect(f.SetOnlyLocalCallCount()).To(Equal(1))
		},
		expErr: errors.New("err"),
	}),

	Entry("fails to evaluate whether aws-node uses IRSA", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
			k.SupportsManagedNodesReturns(true, nil)
			init.DoesAWSNodeUseIRSAReturns(true, errors.New("err"))
		},
		expectedCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
			Expect(k.NewRawClientCallCount()).To(Equal(1))
			Expect(k.ServerVersionCallCount()).To(Equal(1))
			Expect(k.LoadClusterIntoSpecFromStackCallCount()).To(Equal(1))
			Expect(k.SupportsManagedNodesCallCount()).To(Equal(1))
			Expect(init.NewAWSSelectorSessionCallCount()).To(Equal(1))
			Expect(init.ExpandInstanceSelectorOptionsCallCount()).To(Equal(1))
			Expect(k.ValidateClusterForCompatibilityCallCount()).To(Equal(1))
			Expect(f.SetOnlyLocalCallCount()).To(Equal(1))
			Expect(init.DoesAWSNodeUseIRSACallCount()).To(Equal(1))
		},
		expErr: errors.New("err"),
	}),

	Entry("stack manager fails to do ng tasks", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
			k.NewRawClientReturns(&kubernetes.RawClient{}, nil)
			k.SupportsManagedNodesReturns(true, nil)
			init.DoAllNodegroupStackTasksReturns(errors.New("err"))
		},
		expectedCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
			Expect(k.NewRawClientCallCount()).To(Equal(1))
			Expect(k.ServerVersionCallCount()).To(Equal(1))
			Expect(k.LoadClusterIntoSpecFromStackCallCount()).To(Equal(1))
			Expect(k.SupportsManagedNodesCallCount()).To(Equal(1))
			Expect(init.NewAWSSelectorSessionCallCount()).To(Equal(1))
			Expect(init.ExpandInstanceSelectorOptionsCallCount()).To(Equal(1))
			Expect(k.ValidateClusterForCompatibilityCallCount()).To(Equal(1))
			Expect(f.SetOnlyLocalCallCount()).To(Equal(1))
			Expect(init.DoesAWSNodeUseIRSACallCount()).To(Equal(1))
			Expect(init.DoAllNodegroupStackTasksCallCount()).To(Equal(1))
		},
		expErr: errors.New("err"),
	}),

	Entry("fails to update auth configmap", ngEntry{
		opts: nodegroup.CreateOpts{
			UpdateAuthConfigMap: true,
		},
		mockCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
			k.SupportsManagedNodesReturns(true, nil)
			k.UpdateAuthConfigMapReturns(errors.New("err"))
		},
		expectedCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
			Expect(k.NewRawClientCallCount()).To(Equal(1))
			Expect(k.ServerVersionCallCount()).To(Equal(1))
			Expect(k.LoadClusterIntoSpecFromStackCallCount()).To(Equal(1))
			Expect(k.SupportsManagedNodesCallCount()).To(Equal(1))
			Expect(init.NewAWSSelectorSessionCallCount()).To(Equal(1))
			Expect(init.ExpandInstanceSelectorOptionsCallCount()).To(Equal(1))
			Expect(k.ValidateClusterForCompatibilityCallCount()).To(Equal(1))
			Expect(f.SetOnlyLocalCallCount()).To(Equal(1))
			Expect(init.DoesAWSNodeUseIRSACallCount()).To(Equal(1))
			Expect(init.DoAllNodegroupStackTasksCallCount()).To(Equal(1))
			Expect(k.UpdateAuthConfigMapCallCount()).To(Equal(1))
		},
		expErr: errors.New("err"),
	}),

	Entry("when unable to validate existing ng for compatibility, logs but does not error", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
			k.SupportsManagedNodesReturns(true, nil)
			init.ValidateExistingNodeGroupsForCompatibilityReturns(errors.New("err"))
		},
		expectedCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
			Expect(k.NewRawClientCallCount()).To(Equal(1))
			Expect(k.ServerVersionCallCount()).To(Equal(1))
			Expect(k.LoadClusterIntoSpecFromStackCallCount()).To(Equal(1))
			Expect(k.SupportsManagedNodesCallCount()).To(Equal(1))
			Expect(init.NewAWSSelectorSessionCallCount()).To(Equal(1))
			Expect(init.ExpandInstanceSelectorOptionsCallCount()).To(Equal(1))
			Expect(k.ValidateClusterForCompatibilityCallCount()).To(Equal(1))
			Expect(f.SetOnlyLocalCallCount()).To(Equal(1))
			Expect(init.DoesAWSNodeUseIRSACallCount()).To(Equal(1))
			Expect(init.DoAllNodegroupStackTasksCallCount()).To(Equal(1))
			Expect(init.ValidateExistingNodeGroupsForCompatibilityCallCount()).To(Equal(1))
		},
		expErr: nil,
	}),

	Entry("[happy path] creates nodegroup with no options", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
			k.SupportsManagedNodesReturns(true, nil)
		},
		expectedCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
			Expect(k.NewRawClientCallCount()).To(Equal(1))
			Expect(k.ServerVersionCallCount()).To(Equal(1))
			Expect(k.LoadClusterIntoSpecFromStackCallCount()).To(Equal(1))
			Expect(k.SupportsManagedNodesCallCount()).To(Equal(1))
			Expect(init.NewAWSSelectorSessionCallCount()).To(Equal(1))
			Expect(init.ExpandInstanceSelectorOptionsCallCount()).To(Equal(1))
			Expect(k.ValidateClusterForCompatibilityCallCount()).To(Equal(1))
			Expect(f.SetOnlyLocalCallCount()).To(Equal(1))
			Expect(init.DoesAWSNodeUseIRSACallCount()).To(Equal(1))
			Expect(init.DoAllNodegroupStackTasksCallCount()).To(Equal(1))
			Expect(init.ValidateExistingNodeGroupsForCompatibilityCallCount()).To(Equal(1))
		},
		expErr: nil,
	}),

	Entry("[happy path] creates nodegroup with all the options", ngEntry{
		opts: nodegroup.CreateOpts{
			DryRun:                    true,
			UpdateAuthConfigMap:       true,
			InstallNeuronDevicePlugin: true,
			InstallNvidiaDevicePlugin: true,
			SkipOutdatedAddonsCheck:   true,
			ConfigFileProvided:        true,
		},
		mockCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
			k.SupportsManagedNodesReturns(true, nil)
		},
		expectedCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
			Expect(k.NewRawClientCallCount()).To(Equal(1))
			Expect(k.ServerVersionCallCount()).To(Equal(1))
			Expect(k.LoadClusterIntoSpecFromStackCallCount()).To(Equal(1))
			Expect(k.SupportsManagedNodesCallCount()).To(Equal(1))
			Expect(init.NewAWSSelectorSessionCallCount()).To(Equal(1))
			Expect(init.ExpandInstanceSelectorOptionsCallCount()).To(Equal(1))
			Expect(k.ValidateClusterForCompatibilityCallCount()).To(Equal(1))
			Expect(f.SetOnlyLocalCallCount()).To(Equal(1))
		},
		expErr: nil,
	}),
)

var _ = Describe("create", func() {
	When("creating an unmanaged nodegroup for an ipv6 cluster", func() {
		It("returns an error", func() {
			cfg := newClusterConfig()
			cfg.Metadata.Version = api.Version1_21
			cfg.IAM = &api.ClusterIAM{
				WithOIDC: api.Enabled(),
			}
			cfg.Addons = []*api.Addon{
				{
					Name: api.VPCCNIAddon,
				},
				{
					Name: api.KubeProxyAddon,
				},
				{
					Name: api.CoreDNSAddon,
				},
			}
			cfg.VPC.IPFamily = utilsstrings.Pointer(string(api.IPV6Family))

			p := mockprovider.NewMockProvider()
			ctl := &eks.ClusterProvider{
				Provider: p,
				Status: &eks.ProviderStatus{
					ClusterInfo: &eks.ClusterInfo{
						Cluster: testutils.NewFakeCluster("my-cluster", ""),
					},
				},
			}
			m := nodegroup.New(cfg, ctl, nil)

			k := &fakes.FakeKubeProvider{}
			k.SupportsManagedNodesReturns(true, nil)
			m.MockKubeProvider(k)

			init := &fakes.FakeNodeGroupInitialiser{}
			m.MockNodeGroupService(init)

			ngFilter := utilFakes.FakeNodegroupFilter{}

			err := m.Create(nodegroup.CreateOpts{}, &ngFilter)
			Expect(err).To(MatchError(ContainSubstring("unmanaged nodegroups are not supported with IPv6 clusters")))
		})
	})
	When("creating a nodegroups for unowned ipv6 cluster", func() {
		It("returns an error", func() {
			cfg := newClusterConfig()
			cfg.Metadata.Version = api.Version1_21
			cfg.IAM = &api.ClusterIAM{
				WithOIDC: api.Enabled(),
			}
			cfg.Addons = []*api.Addon{
				{
					Name: api.VPCCNIAddon,
				},
				{
					Name: api.KubeProxyAddon,
				},
				{
					Name: api.CoreDNSAddon,
				},
			}
			cfg.VPC.IPFamily = utilsstrings.Pointer(string(api.IPV6Family))

			p := mockprovider.NewMockProvider()
			p.MockEC2().On("DescribeVpcs", &ec2.DescribeVpcsInput{VpcIds: aws.StringSlice([]string{"custom"})}).Return(&ec2.DescribeVpcsOutput{
				Vpcs: []*ec2.Vpc{
					{
						VpcId:     aws.String("custom"),
						CidrBlock: aws.String("192.168.0.0/16"),
					},
				},
			}, nil)
			p.MockEC2().On("DescribeSubnets", &ec2.DescribeSubnetsInput{
				SubnetIds: aws.StringSlice([]string{"sn-private-1", "sn-private-2"}),
			}).Return(&ec2.DescribeSubnetsOutput{
				Subnets: []*ec2.Subnet{
					{
						SubnetId:         aws.String("sn-private-1"),
						VpcId:            aws.String("custom"),
						AvailabilityZone: aws.String("us-east-1b"),
						CidrBlock:        aws.String("192.168.0.0/24"),
					},
					{
						SubnetId:         aws.String("sn-private-2"),
						VpcId:            aws.String("custom"),
						AvailabilityZone: aws.String("us-east-1b"),
						CidrBlock:        aws.String("192.168.0.1/24"),
					},
				},
			}, nil)
			// The order is undetermined because the entries are maps and we don't sort the ids.
			p.MockEC2().On("DescribeSubnets", mock.Anything).Return(&ec2.DescribeSubnetsOutput{
				Subnets: []*ec2.Subnet{
					{
						SubnetId:         aws.String("sn-public-1"),
						VpcId:            aws.String("custom"),
						AvailabilityZone: aws.String("us-east-1b"),
						CidrBlock:        aws.String("192.168.1.0/24"),
					},
					{
						SubnetId:         aws.String("sn-public-2"),
						VpcId:            aws.String("custom"),
						AvailabilityZone: aws.String("us-east-1b"),
						CidrBlock:        aws.String("192.168.2.0/24"),
					},
					{
						SubnetId:         aws.String("sn-public-3"),
						VpcId:            aws.String("custom"),
						AvailabilityZone: aws.String("us-east-1b"),
						CidrBlock:        aws.String("192.168.3.0/24"),
					},
					{
						SubnetId:         aws.String("sn-public-4"),
						VpcId:            aws.String("custom"),
						AvailabilityZone: aws.String("us-east-1b"),
						CidrBlock:        aws.String("192.168.4.0/24"),
					},
				},
			}, nil)
			ctl := &eks.ClusterProvider{
				Provider: p,
				Status: &eks.ProviderStatus{
					ClusterInfo: &eks.ClusterInfo{
						Cluster: testutils.NewFakeCluster("my-cluster", ""),
					},
				},
			}
			cfg.VPC.Subnets = &api.ClusterSubnets{
				Public: api.AZSubnetMapping{
					"sn-public-1": api.AZSubnetSpec{ID: "sn-public-1"},
					"sn-public-2": api.AZSubnetSpec{ID: "sn-public-2"},
					"sn-public-3": api.AZSubnetSpec{ID: "sn-public-3"},
					"sn-public-4": api.AZSubnetSpec{ID: "sn-public-4"},
				},
			}
			cfg.VPC.ID = "custom"
			cfg.VPC.SecurityGroup = "sg"
			cfg.NodeGroups = nil
			m := nodegroup.New(cfg, ctl, nil)

			k := &fakes.FakeKubeProvider{}
			k.LoadClusterIntoSpecFromStackReturns(&manager.StackNotFoundErr{})
			k.SupportsManagedNodesReturns(true, nil)
			m.MockKubeProvider(k)

			init := &fakes.FakeNodeGroupInitialiser{}
			m.MockNodeGroupService(init)

			ngFilter := utilFakes.FakeNodegroupFilter{}

			err := m.Create(nodegroup.CreateOpts{}, &ngFilter)
			Expect(err).To(MatchError(ContainSubstring("nodegroups cannot be created on IPv6 unowned clusters")))
		})
	})
})

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
		VPC: api.NewClusterVPC(),
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
