package irsa

import (
	"context"
	"fmt"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// GetOptions holds the configuration for the IRSA get action
type GetOptions struct {
	Name      string
	Namespace string
}

func (m *Manager) Get(ctx context.Context, options GetOptions) ([]*api.ClusterIAMServiceAccount, error) {
	remoteServiceAccounts, err := m.stackManager.GetIAMServiceAccounts(ctx, options.Name, options.Namespace)
	if err != nil {
		return nil, fmt.Errorf("getting iamserviceaccounts: %w", err)
	}
	return remoteServiceAccounts, nil
}
