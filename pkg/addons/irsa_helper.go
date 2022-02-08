package addons

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/pkg/errors"

	"github.com/weaveworks/eksctl/pkg/actions/irsa"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
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
	irsaManager  *irsa.Manager
	stackManager manager.StackManager
	clusterName  string
}

// NewIRSAHelper creates a new IRSAHelper
func NewIRSAHelper(oidc *iamoidc.OpenIDConnectManager, stackManager manager.StackManager, irsaManager *irsa.Manager, clusterName string) IRSAHelper {
	return &irsaHelper{
		oidc:         oidc,
		stackManager: stackManager,
		irsaManager:  irsaManager,
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

// CreateOrUpdate creates IRSA for the specified IAM service accounts or updates it
func (h *irsaHelper) CreateOrUpdate(sa *api.ClusterIAMServiceAccount) error {
	serviceAccounts := []*api.ClusterIAMServiceAccount{sa}
	name := makeIAMServiceAccountStackName(h.clusterName, sa.Namespace, sa.Name)
	stack, err := h.stackManager.DescribeStack(&manager.Stack{StackName: &name})
	if err != nil {
		if awsError, ok := errors.Unwrap(errors.Unwrap(err)).(awserr.Error); !ok || ok &&
			awsError.Code() != "ValidationError" {
			return errors.Wrapf(err, "error checking if iamserviceaccount %s/%s exists", sa.Namespace, sa.Name)
		}
	}
	if stack == nil {
		err = h.irsaManager.CreateIAMServiceAccount(serviceAccounts, false)
	} else {
		err = h.irsaManager.UpdateIAMServiceAccounts(serviceAccounts, []*manager.Stack{stack}, false)
	}
	return err
}

func makeIAMServiceAccountStackName(clusterName, namespace, name string) string {
	return fmt.Sprintf("eksctl-%s-addon-iamserviceaccount-%s-%s", clusterName, namespace, name)
}
