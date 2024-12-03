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
	remoteServiceAccounts, err := m.stackManager.GetIAMServiceAccounts(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting iamserviceaccounts: %w", err)
	}

	if options.Namespace != "" {
		remoteServiceAccounts = filterByNamespace(remoteServiceAccounts, options.Namespace)
	}

	if options.Name != "" {
		remoteServiceAccounts = filterByName(remoteServiceAccounts, options.Name)
	}

	return remoteServiceAccounts, nil
}

func filterByNamespace(serviceAccounts []*api.ClusterIAMServiceAccount, namespace string) []*api.ClusterIAMServiceAccount {
	var serviceAccountsMatching []*api.ClusterIAMServiceAccount
	for _, sa := range serviceAccounts {
		if sa.Namespace == namespace {
			serviceAccountsMatching = append(serviceAccountsMatching, sa)
		}
	}
	return serviceAccountsMatching
}

func filterByName(serviceAccounts []*api.ClusterIAMServiceAccount, name string) []*api.ClusterIAMServiceAccount {
	var serviceAccountsMatching []*api.ClusterIAMServiceAccount
	for _, sa := range serviceAccounts {
		if sa.Name == name {
			serviceAccountsMatching = append(serviceAccountsMatching, sa)
		}
	}
	return serviceAccountsMatching
}
