package addons

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

type serviceAccountCreator interface {
	NewTasksToCreateIAMServiceAccounts(serviceAccounts []*api.ClusterIAMServiceAccount, oidc *iamoidc.OpenIDConnectManager, clientSetGetter kubernetes.ClientSetGetter, replaceExistingRole bool) *tasks.TaskTree
}

// IRSAHelper provides methods for enabling IRSA
type IRSAHelper interface {
	IsSupported() (bool, error)
	Create(serviceAccounts []*api.ClusterIAMServiceAccount) error
}

// irsaHelper applies the annotations required for a ServiceAccount to work with IRSA
type irsaHelper struct {
	oidc *iamoidc.OpenIDConnectManager
	serviceAccountCreator
	clientSet kubernetes.ClientSetGetter
}

// NewIRSAHelper creates a new IRSAHelper
func NewIRSAHelper(oidc *iamoidc.OpenIDConnectManager, saCreator serviceAccountCreator, clientSet kubernetes.ClientSetGetter) IRSAHelper {
	return &irsaHelper{
		oidc:                  oidc,
		serviceAccountCreator: saCreator,
		clientSet:             clientSet,
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
func (h *irsaHelper) Create(serviceAccounts []*api.ClusterIAMServiceAccount) error {
	taskTree := h.NewTasksToCreateIAMServiceAccounts(serviceAccounts, h.oidc, h.clientSet, true)
	if errs := taskTree.DoAllSync(); len(errs) > 0 {
		return errors.Wrap(joinErrors(errs), "error creating IAM service account")
	}
	return nil
}

func joinErrors(errs []error) error {
	if len(errs) == 1 {
		return errs[0]
	}
	allErrs := []string{"errors:\n"}
	for _, err := range errs {
		allErrs = append(allErrs, fmt.Sprintf("- %v", err))
	}
	return errors.New(strings.Join(allErrs, "\n"))
}
