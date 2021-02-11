package addon

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/eks/eksiface"

	kubeclient "k8s.io/client-go/kubernetes"

	"github.com/weaveworks/eksctl/pkg/utils"

	"github.com/weaveworks/eksctl/pkg/cfn/manager"

	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

type Manager struct {
	clusterConfig *api.ClusterConfig
	eksAPI        eksiface.EKSAPI
	withOIDC      bool
	oidcManager   *iamoidc.OpenIDConnectManager
	stackManager  manager.StackManager
	clientSet     kubeclient.Interface
}

func New(clusterConfig *api.ClusterConfig, eksAPI eksiface.EKSAPI, stackManager manager.StackManager, withOIDC bool, oidcManager *iamoidc.OpenIDConnectManager, clientSet kubeclient.Interface) (*Manager, error) {
	if err := supportedVersion(clusterConfig.Metadata.Version); err != nil {
		return nil, err
	}

	return &Manager{
		clusterConfig: clusterConfig,
		eksAPI:        eksAPI,
		withOIDC:      withOIDC,
		oidcManager:   oidcManager,
		stackManager:  stackManager,
		clientSet:     clientSet,
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
