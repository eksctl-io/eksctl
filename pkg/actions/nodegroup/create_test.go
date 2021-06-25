package nodegroup_test

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	utilFakes "github.com/weaveworks/eksctl/pkg/ctl/cmdutils/filter/fakes"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/eks/fakes"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

type ngEntry struct {
	version   string
	opts      nodegroup.CreateOpts
	mockCalls func(*fakes.FakeKubeProvider, *fakes.FakeNodeGroupInitialiser, *utilFakes.FakeNodegroupFilter)
	expErr    error
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
	if err != nil {
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(ContainSubstring(t.expErr.Error())))
		return
	}

	Expect(err).NotTo(HaveOccurred())
},
	Entry("fails when cluster version is not supported", ngEntry{
		version: "1.14",
		mockCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
		},
		expErr: fmt.Errorf("invalid version, %s is no longer supported, supported values: auto, default, latest, %s\nsee also: https://docs.aws.amazon.com/eks/latest/userguide/kubernetes-versions.html", "1.14", strings.Join(api.SupportedVersions(), ", ")),
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
			k.NewRawClientReturns(&kubernetes.RawClient{}, nil)
			k.ServerVersionReturns("1.17", nil)
			k.LoadClusterIntoSpecFromStackReturns(nil)
			k.SupportsManagedNodesReturns(false, errors.New("err"))
		},
		expErr: errors.New("err"),
	}),

	Entry("fails to set instance types to instances matched by instance selector criteria", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
			k.NewRawClientReturns(&kubernetes.RawClient{}, nil)
			k.ServerVersionReturns("1.17", nil)
			k.LoadClusterIntoSpecFromStackReturns(nil)
			k.SupportsManagedNodesReturns(true, nil)
			init.NewAWSSelectorSessionReturns(nil)
			init.ExpandInstanceSelectorOptionsReturns(errors.New("err"))
		},
		expErr: errors.New("err"),
	}),

	Entry("fails when cluster is not compatible with ng config", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
			k.NewRawClientReturns(&kubernetes.RawClient{}, nil)
			k.ServerVersionReturns("1.17", nil)
			k.LoadClusterIntoSpecFromStackReturns(nil)
			k.SupportsManagedNodesReturns(true, nil)
			init.NewAWSSelectorSessionReturns(nil)
			init.ExpandInstanceSelectorOptionsReturns(nil)
			k.ValidateClusterForCompatibilityReturns(errors.New("err"))
		},
		expErr: errors.Wrap(errors.New("err"), "cluster compatibility check failed"),
	}),

	Entry("fails when it cannot validate legacy subnets for ng", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
			k.NewRawClientReturns(&kubernetes.RawClient{}, nil)
			k.ServerVersionReturns("1.17", nil)
			k.LoadClusterIntoSpecFromStackReturns(nil)
			k.SupportsManagedNodesReturns(true, nil)
			init.NewAWSSelectorSessionReturns(nil)
			init.ExpandInstanceSelectorOptionsReturns(nil)
			k.ValidateClusterForCompatibilityReturns(nil)
			init.ValidateLegacySubnetsForNodeGroupsReturns(errors.New("err"))
		},
		expErr: errors.New("err"),
	}),

	Entry("fails when existing local ng stacks in config file is not listed", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
			k.NewRawClientReturns(&kubernetes.RawClient{}, nil)
			k.ServerVersionReturns("1.17", nil)
			k.LoadClusterIntoSpecFromStackReturns(nil)
			k.SupportsManagedNodesReturns(true, nil)
			init.NewAWSSelectorSessionReturns(nil)
			init.ExpandInstanceSelectorOptionsReturns(nil)
			k.ValidateClusterForCompatibilityReturns(nil)
			f.SetOnlyLocalReturns(errors.New("err"))
		},
		expErr: errors.New("err"),
	}),

	Entry("fails to evaluate whether aws-node uses IRSA", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
			k.NewRawClientReturns(&kubernetes.RawClient{}, nil)
			k.ServerVersionReturns("1.17", nil)
			k.LoadClusterIntoSpecFromStackReturns(nil)
			k.SupportsManagedNodesReturns(true, nil)
			init.NewAWSSelectorSessionReturns(nil)
			init.ExpandInstanceSelectorOptionsReturns(nil)
			k.ValidateClusterForCompatibilityReturns(nil)
			f.SetOnlyLocalReturns(nil)
			init.DoesAWSNodeUseIRSAReturns(true, errors.New("err"))
		},
		expErr: errors.New("err"),
	}),

	Entry("stack manager fails to do ng tasks", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
			k.NewRawClientReturns(&kubernetes.RawClient{}, nil)
			k.ServerVersionReturns("1.17", nil)
			k.LoadClusterIntoSpecFromStackReturns(nil)
			k.SupportsManagedNodesReturns(true, nil)
			init.NewAWSSelectorSessionReturns(nil)
			init.ExpandInstanceSelectorOptionsReturns(nil)
			k.ValidateClusterForCompatibilityReturns(nil)
			f.SetOnlyLocalReturns(nil)
			init.DoesAWSNodeUseIRSAReturns(false, nil)
			init.DoAllNodegroupStackTasksReturns(errors.New("err"))
		},
		expErr: errors.New("err"),
	}),

	Entry("fails to update auth configmap", ngEntry{
		opts: nodegroup.CreateOpts{
			DryRun:              true,
			UpdateAuthConfigMap: true,
		},
		mockCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
			k.NewRawClientReturns(&kubernetes.RawClient{}, nil)
			k.ServerVersionReturns("1.17", nil)
			k.LoadClusterIntoSpecFromStackReturns(nil)
			k.SupportsManagedNodesReturns(true, nil)
			init.NewAWSSelectorSessionReturns(nil)
			init.ExpandInstanceSelectorOptionsReturns(nil)
			k.ValidateClusterForCompatibilityReturns(nil)
			f.SetOnlyLocalReturns(nil)
			init.DoesAWSNodeUseIRSAReturns(false, nil)
			init.DoAllNodegroupStackTasksReturns(nil)
			k.UpdateAuthConfigMapReturns(errors.New("err"))
		},
		expErr: errors.New("err"),
	}),

	Entry("fails to validate existing ng for compatibility", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
			k.NewRawClientReturns(&kubernetes.RawClient{}, nil)
			k.ServerVersionReturns("1.17", nil)
			k.LoadClusterIntoSpecFromStackReturns(nil)
			k.SupportsManagedNodesReturns(true, nil)
			init.NewAWSSelectorSessionReturns(nil)
			init.ExpandInstanceSelectorOptionsReturns(nil)
			k.ValidateClusterForCompatibilityReturns(nil)
			f.SetOnlyLocalReturns(nil)
			init.DoesAWSNodeUseIRSAReturns(false, nil)
			init.DoAllNodegroupStackTasksReturns(nil)
			init.ValidateExistingNodeGroupsForCompatibilityReturns(errors.New("err"))
		},
		expErr: errors.New("err"),
	}),

	Entry("[happy path] creates nodegroup with no options", ngEntry{
		mockCalls: func(k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
			k.NewRawClientReturns(&kubernetes.RawClient{}, nil)
			k.ServerVersionReturns("1.17", nil)
			k.LoadClusterIntoSpecFromStackReturns(nil)
			k.SupportsManagedNodesReturns(true, nil)
			init.NewAWSSelectorSessionReturns(nil)
			init.ExpandInstanceSelectorOptionsReturns(nil)
			k.ValidateClusterForCompatibilityReturns(nil)
			f.SetOnlyLocalReturns(nil)
			init.DoesAWSNodeUseIRSAReturns(false, nil)
			init.DoAllNodegroupStackTasksReturns(nil)
			init.ValidateExistingNodeGroupsForCompatibilityReturns(nil)
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
			k.NewRawClientReturns(&kubernetes.RawClient{}, nil)
			k.ServerVersionReturns("1.17", nil)
			k.LoadClusterIntoSpecFromStackReturns(nil)
			k.SupportsManagedNodesReturns(true, nil)
			init.NewAWSSelectorSessionReturns(nil)
			init.ExpandInstanceSelectorOptionsReturns(nil)
			k.ValidateClusterForCompatibilityReturns(nil)
			f.SetOnlyLocalReturns(nil)
			init.DoesAWSNodeUseIRSAReturns(false, nil)
			init.DoAllNodegroupStackTasksReturns(nil)
			k.UpdateAuthConfigMapReturns(nil)
		},
		expErr: nil,
	}),
)

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
