package addons

import (
	"fmt"

	"github.com/weaveworks/eksctl/pkg/actions/iam"

	"github.com/weaveworks/eksctl/pkg/cfn/manager"

	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
)

// IRSAHelper provides methods for enabling IRSA
type IRSAHelper interface {
	IsSupported() (bool, error)
	CreateOrUpdate(serviceAccounts *api.ClusterIAMServiceAccount) error
}

// irsaHelper applies the annotations required for a ServiceAccount to work with IRSA
type irsaHelper struct {
	oidc         *iamoidc.OpenIDConnectManager
	iamManager   *iam.Manager
	stackManager *manager.StackCollection
	clusterName  string
}

// NewIRSAHelper creates a new IRSAHelper
func NewIRSAHelper(oidc *iamoidc.OpenIDConnectManager, stackManager *manager.StackCollection, iamManager *iam.Manager, clusterName string) IRSAHelper {
	return &irsaHelper{
		oidc:         oidc,
		stackManager: stackManager,
		iamManager:   iamManager,
		clusterName:  clusterName,
	}
}

// IsSupported checks whether IRSA is supported or not
func (h *irsaHelper) IsSupported() (bool, error) {
	exists, err := h.oidc.CheckProviderExists()
	if err != nil {
		return false, errors.Wrapf(err, "error checking OIDC provider")
	}
	return exists, nil
}

// Create creates IRSA for the specified IAM service accounts
func (h *irsaHelper) CreateOrUpdate(sa *api.ClusterIAMServiceAccount) error {
	serviceAccounts := []*api.ClusterIAMServiceAccount{sa}
	stacks, err := h.stackManager.ListStacksMatching(makeIAMServiceAccountStackName(h.clusterName, sa.Namespace, sa.Name))
	if err != nil {
		return errors.Wrapf(err, "error checking if iamserviceaccount %s/%s exists", sa.Namespace, sa.Name)
	}
	if len(stacks) == 0 {
		err = h.iamManager.CreateIAMServiceAccount(serviceAccounts, false)
	} else {
		err = h.iamManager.UpdateIAMServiceAccounts(serviceAccounts, false)
	}
	return err
}

func makeIAMServiceAccountStackName(clusterName, namespace, name string) string {
	return fmt.Sprintf("eksctl-%s-addon-iamserviceaccount-%s-%s", clusterName, namespace, name)
}
