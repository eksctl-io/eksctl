package addon

import (
	"fmt"

	kubeclient "k8s.io/client-go/kubernetes"

	"github.com/weaveworks/eksctl/pkg/utils"

	"github.com/weaveworks/eksctl/pkg/cfn/manager"

	"github.com/weaveworks/eksctl/pkg/cfn/builder"

	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
)

type Manager struct {
	clusterConfig   *api.ClusterConfig
	clusterProvider *eks.ClusterProvider
	withOIDC        bool
	oidcManager     *iamoidc.OpenIDConnectManager
	stackManager    StackManager
	clientSet       kubeclient.Interface
}

//go:generate counterfeiter -o fakes/fake_stack_manager.go . StackManager
type StackManager interface {
	CreateStack(name string, stack builder.ResourceSet, tags, parameters map[string]string, errs chan error) error
	DeleteStackByName(name string) (*manager.Stack, error)
	ListStacksMatching(nameRegex string, statusFilters ...string) ([]*manager.Stack, error)
	UpdateStack(stackName, changeSetName, description string, templateData manager.TemplateData, parameters map[string]string) error
}

func New(clusterConfig *api.ClusterConfig, clusterProvider *eks.ClusterProvider, stackManager StackManager, withOIDC bool, oidcManager *iamoidc.OpenIDConnectManager, clientSet kubeclient.Interface) (*Manager, error) {
	if err := supportedVersion(clusterConfig.Metadata.Version); err != nil {
		return nil, err
	}

	return &Manager{
		clusterConfig:   clusterConfig,
		clusterProvider: clusterProvider,
		withOIDC:        withOIDC,
		oidcManager:     oidcManager,
		stackManager:    stackManager,
		clientSet:       clientSet,
	}, nil
}

func supportedVersion(version string) error {
	supported, err := utils.IsMinVersion(api.Version1_18, version)
	if err != nil {
		return err
	}
	switch supported {
	case true:
		return nil
	default:
		return fmt.Errorf("addons not supported on %s. Must be using %s or newer", version, api.Version1_18)
	}
}

func (a *Manager) makeAddonName(name string) string {
	return fmt.Sprintf("eksctl-%s-addon-%s", a.clusterConfig.Metadata.Name, name)
}
