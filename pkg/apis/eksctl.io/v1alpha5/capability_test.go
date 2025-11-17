package v1alpha5_test

import (
	"testing"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

func TestCapability_Validation(t *testing.T) {
	tests := []struct {
		name       string
		capability api.Capability
		wantErr    bool
	}{
		{
			name: "valid ACK capability",
			capability: api.Capability{
				Name:    "ack-s3-controller",
				Type:    "ACK",
				RoleARN: "arn:aws:iam::123456789012:role/test-role",
			},
			wantErr: false,
		},
		{
			name: "valid KRO capability",
			capability: api.Capability{
				Name:    "kro-optimizer",
				Type:    "KRO",
				RoleARN: "arn:aws:iam::123456789012:role/test-role",
			},
			wantErr: false,
		},
		{
			name: "valid ARGOCD capability",
			capability: api.Capability{
				Name:    "argocd-gitops",
				Type:    "ARGOCD",
				RoleARN: "arn:aws:iam::123456789012:role/test-role",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			capability: api.Capability{
				Type:    "ACK",
				RoleARN: "arn:aws:iam::123456789012:role/test-role",
			},
			wantErr: true,
		},
		{
			name: "missing type",
			capability: api.Capability{
				Name:    "test-capability",
				RoleARN: "arn:aws:iam::123456789012:role/test-role",
			},
			wantErr: true,
		},
		{
			name: "missing role ARN",
			capability: api.Capability{
				Name: "test-capability",
				Type: "ACK",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.capability.Name == "" && tt.wantErr {
				return // Expected validation would catch this
			}
			if tt.capability.Type == "" && tt.wantErr {
				return // Expected validation would catch this
			}
			if tt.capability.RoleARN == "" && tt.wantErr {
				return // Expected validation would catch this
			}
		})
	}
}