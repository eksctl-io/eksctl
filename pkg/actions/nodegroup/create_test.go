package nodegroup_test

import (
	"fmt"
	"testing"

	"github.com/tj/assert"
	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils/filter"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

func TestCreate(t *testing.T) {
	tests := map[string]struct {
		version   string
		pStatus   *eks.ProviderStatus
		mockCalls func(*mockprovider.MockProvider)

		expErr error
	}{
		// "cluster version is not supported": {
		// 	version:   "1.14",
		// 	pStatus:   &eks.ProviderStatus{},
		// 	expErr:    fmt.Errorf("invalid version, %s is no longer supported, supported values: auto, default, latest, %s\nsee also: https://docs.aws.amazon.com/eks/latest/userguide/kubernetes-versions.html", "1.14", strings.Join(api.SupportedVersions(), ", ")),
		// 	mockCalls: func(p *mockprovider.MockProvider) {},
		// },
		"no ARM support": {
			version: "1.17",
			pStatus: &eks.ProviderStatus{
				ClusterInfo: &eks.ClusterInfo{
					Cluster: testutils.NewFakeCluster("my-cluster", ""),
				},
			},
			mockCalls: func(p *mockprovider.MockProvider) {},
			expErr:    fmt.Errorf("err"),
		},
		// "creating a cluster returns no error": {
		// 	pStatus: &eks.ProviderStatus{
		// 		ClusterInfo: &eks.ClusterInfo{
		// 			Cluster: &awseks.Cluster{
		// 				Version: aws.String("1.15"),
		// 			},
		// 		},
		// 	},
		// 	mockCalls: func(p *mockprovider.MockProvider) {
		// 		p.MockEKS().On("NewRawClient").Return(fmt.Errorf("err"))
		// 	},
		// 	expErr: nil,
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

			p := mockprovider.NewMockProvider()
			m := nodegroup.New(cfg, &eks.ClusterProvider{
				Provider: p,
				Status:   tc.pStatus,
			}, nil)

			if tc.mockCalls != nil {
				tc.mockCalls(p)
			}

			err := m.Create(nodegroup.CreateOpts{}, *filter.NewNodeGroupFilter())
			if err != nil {
				assert.Equal(t, err, tc.expErr)
			}
		})
	}
}
