package nodegroup_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/tj/assert"
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

func TestCreateNodegroups(t *testing.T) {
	tests := map[string]struct {
		version   string
		pStatus   *eks.ProviderStatus
		mockCalls func(*mockprovider.MockProvider, *fakes.FakeKubeProvider, *fakes.FakeNodeGroupInitialiser, *utilFakes.FakeNodegroupFilter)

		expErr error
	}{
		"cluster version is not supported": {
			version: "1.14",
			pStatus: &eks.ProviderStatus{},
			mockCalls: func(p *mockprovider.MockProvider, k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
			},
			expErr: fmt.Errorf("invalid version, %s is no longer supported, supported values: auto, default, latest, %s\nsee also: https://docs.aws.amazon.com/eks/latest/userguide/kubernetes-versions.html", "1.14", strings.Join(api.SupportedVersions(), ", ")),
		},
		"fails ARM support check": {
			version: "1.17",
			pStatus: &eks.ProviderStatus{
				ClusterInfo: &eks.ClusterInfo{
					Cluster: testutils.NewFakeCluster("my-cluster", ""),
				},
			},
			mockCalls: func(p *mockprovider.MockProvider, k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
				k.NewRawClientReturns(nil, fmt.Errorf("err"))
			},
			expErr: fmt.Errorf("err"),
		},
		"fails to load VPC from config": {
			version: "1.17",
			pStatus: &eks.ProviderStatus{
				ClusterInfo: &eks.ClusterInfo{
					Cluster: testutils.NewFakeCluster("my-cluster", ""),
				},
			},
			mockCalls: func(p *mockprovider.MockProvider, k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
				k.NewRawClientReturns(&kubernetes.RawClient{}, nil)
				k.ServerVersionReturns("1.17", nil)
				k.LoadClusterIntoSpecFromStackReturns(&manager.StackNotFoundErr{})
			},
			expErr: errors.Wrapf(errors.New("VPC configuration required for creating nodegroups on clusters not owned by eksctl: vpc.subnets, vpc.id, vpc.securityGroup"), "loading VPC spec for cluster %q", "my-cluster"),
		},
		"cluster does not support managed nodes": {
			version: "1.17",
			pStatus: &eks.ProviderStatus{
				ClusterInfo: &eks.ClusterInfo{
					Cluster: testutils.NewFakeCluster("my-cluster", ""),
				},
			},
			mockCalls: func(p *mockprovider.MockProvider, k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
				k.NewRawClientReturns(&kubernetes.RawClient{}, nil)
				k.ServerVersionReturns("1.17", nil)
				k.LoadClusterIntoSpecFromStackReturns(nil)
				k.SupportsManagedNodesReturns(false, errors.New("bang"))
			},
			expErr: errors.New("bang"),
		},
		"fails when NodeGroupService fails to match instances": {
			version: "1.17",
			pStatus: &eks.ProviderStatus{
				ClusterInfo: &eks.ClusterInfo{
					Cluster: testutils.NewFakeCluster("my-cluster", ""),
				},
			},
			mockCalls: func(p *mockprovider.MockProvider, k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
				k.NewRawClientReturns(&kubernetes.RawClient{}, nil)
				k.ServerVersionReturns("1.17", nil)
				k.LoadClusterIntoSpecFromStackReturns(nil)
				k.SupportsManagedNodesReturns(true, nil)
				init.NewAWSSelectorSessionReturns(nil)
				init.ExpandInstanceSelectorOptionsReturns(errors.New("bang"))
			},
			expErr: errors.New("bang"),
		},
		"fails when cluster is not compatible with ng config": {
			version: "1.17",
			pStatus: &eks.ProviderStatus{
				ClusterInfo: &eks.ClusterInfo{
					Cluster: testutils.NewFakeCluster("my-cluster", ""),
				},
			},
			mockCalls: func(p *mockprovider.MockProvider, k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
				k.NewRawClientReturns(&kubernetes.RawClient{}, nil)
				k.ServerVersionReturns("1.17", nil)
				k.LoadClusterIntoSpecFromStackReturns(nil)
				k.SupportsManagedNodesReturns(true, nil)
				init.NewAWSSelectorSessionReturns(nil)
				init.ExpandInstanceSelectorOptionsReturns(nil)
				k.ValidateClusterForCompatibilityReturns(errors.New("bang"))
			},
			expErr: errors.Wrap(errors.New("bang"), "cluster compatibility check failed"),
		},
		"existing local ng stacks in config file fail to be listed": {
			version: "1.17",
			pStatus: &eks.ProviderStatus{
				ClusterInfo: &eks.ClusterInfo{
					Cluster: testutils.NewFakeCluster("my-cluster", ""),
				},
			},
			mockCalls: func(p *mockprovider.MockProvider, k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser, f *utilFakes.FakeNodegroupFilter) {
				k.NewRawClientReturns(&kubernetes.RawClient{}, nil)
				k.ServerVersionReturns("1.17", nil)
				k.LoadClusterIntoSpecFromStackReturns(nil)
				k.SupportsManagedNodesReturns(true, nil)
				init.NewAWSSelectorSessionReturns(nil)
				init.ExpandInstanceSelectorOptionsReturns(nil)
				k.ValidateClusterForCompatibilityReturns(nil)
				f.SetOnlyLocalReturns(errors.New("bang"))
			},
			expErr: errors.New("bang"),
		},
		// "test": {
		// 	version: "1.17",
		// 	pStatus: &eks.ProviderStatus{
		// 		ClusterInfo: &eks.ClusterInfo{
		// 			Cluster: testutils.NewFakeCluster("my-cluster", ""),
		// 		},
		// 	},
		// 	mockCalls: func(p *mockprovider.MockProvider, k *fakes.FakeKubeProvider, init *fakes.FakeNodeGroupInitialiser) {
		// 		k.NewRawClientReturns(&kubernetes.RawClient{}, nil)
		// 		k.ServerVersionReturns("1.17", nil)
		// 		k.LoadClusterIntoSpecFromStackReturns(nil)
		// 		k.SupportsManagedNodesReturns(true, nil)
		// 		init.NewAWSSelectorSessionReturns(nil)
		// 		init.ExpandInstanceSelectorOptionsReturns(nil)
		// 		k.ValidateClusterForCompatibilityReturns(nil)
		// 	},
		// 	expErr: errors.Wrap(errors.New("bang"), "cluster compatibility check failed"),
		// },
	}
	for k, tc := range tests {
		t.Run(k, func(t *testing.T) {
			cfg := api.NewClusterConfig()
			cfg.Metadata.Name = "my-cluster"
			cfg.Metadata.Version = tc.version
			cfg.Status = &api.ClusterStatus{
				Endpoint:                 "https://localhost/",
				CertificateAuthorityData: []byte("dGVzdAo="),
			}

			k := &fakes.FakeKubeProvider{}
			init := &fakes.FakeNodeGroupInitialiser{}
			p := mockprovider.NewMockProvider()
			ctl := &eks.ClusterProvider{
				Provider:     p,
				Status:       tc.pStatus,
				KubeProvider: k,
			}
			m := nodegroup.New(cfg, ctl, nil)
			m.MockNodeGroupService(init)
			ngFilter := &utilFakes.FakeNodegroupFilter{}
			if tc.mockCalls != nil {
				tc.mockCalls(p, k, init, ngFilter)
			}

			err := m.Create(nodegroup.CreateOpts{}, ngFilter)
			if err != nil {
				assert.EqualError(t, tc.expErr, err.Error())
			}
		})
	}
}
