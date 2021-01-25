package irsa

import (
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

func (m *Manager) Get(namespace, name string) ([]*api.ClusterIAMServiceAccount, error) {
	remoteServiceAccounts, err := m.stackManager.GetIAMServiceAccounts()
	if err != nil {
		return nil, errors.Wrap(err, "getting iamserviceaccounts")
	}

	if namespace != "" {
		remoteServiceAccounts = filterByNamespace(remoteServiceAccounts, namespace)
	}

	if name != "" {
		remoteServiceAccounts = filterByName(remoteServiceAccounts, name)
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
