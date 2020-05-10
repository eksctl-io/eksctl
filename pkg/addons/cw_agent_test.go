package addons

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockDeployer struct {
	Deployer
	CreateOrReplaceFn func(manifest []byte, plan bool) error
}

func (m *mockDeployer) CreateOrReplace(manifest []byte, plan bool) error {
	return m.CreateOrReplaceFn(manifest, plan)
}

const (
	clusterName = "dummyCluster"
	regionName  = "dummyRegion"
)

func TestCloudwatchAgent_Deploy(t *testing.T) {
	type fields struct {
		client      Deployer
		clusterName string
		region      string
		planMode    bool
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "plan successfully",
			fields: fields{
				client: &mockDeployer{
					CreateOrReplaceFn: func(manifest []byte, plan bool) error {
						str := string(manifest)
						assert.True(t, strings.Contains(str, clusterName))
						assert.True(t, strings.Contains(str, regionName))
						assert.True(t, plan)
						return nil
					},
				},
				clusterName: clusterName,
				region:      regionName,
				planMode:    true,
			},
		},
		{
			name: "deploy successfully",
			fields: fields{
				client: &mockDeployer{
					CreateOrReplaceFn: func(manifest []byte, plan bool) error {
						str := string(manifest)
						assert.True(t, strings.Contains(str, clusterName))
						assert.True(t, strings.Contains(str, regionName))
						assert.True(t, !plan)
						return nil
					},
				},
				clusterName: clusterName,
				region:      regionName,
				planMode:    false,
			},
		},
		{
			name: "fail due to some reason",
			fields: fields{
				client: &mockDeployer{
					CreateOrReplaceFn: func(manifest []byte, plan bool) error {
						str := string(manifest)
						assert.True(t, strings.Contains(str, clusterName))
						assert.True(t, strings.Contains(str, regionName))
						assert.True(t, plan)
						return fmt.Errorf("failed due to some reason")
					},
				},
				clusterName: clusterName,
				region:      regionName,
				planMode:    true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cw := &CloudwatchAgent{
				client:      tt.fields.client,
				clusterName: tt.fields.clusterName,
				region:      tt.fields.region,
				planMode:    tt.fields.planMode,
			}
			if err := cw.Deploy(); (err != nil) != tt.wantErr {
				t.Errorf("Deploy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
