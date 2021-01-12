package iam

import (
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/eks"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
	kubeclient "k8s.io/client-go/kubernetes"
)

type Manager struct {
	clusterName     string
	clusterProvider *eks.ClusterProvider
	oidcManager     *iamoidc.OpenIDConnectManager
	stackManager    StackManager
	clientSet       kubeclient.Interface
}

//go:generate counterfeiter -o fakes/fake_stack_manager.go . StackManager
type StackManager interface {
	ListStacksMatching(nameRegex string, statusFilters ...string) ([]*manager.Stack, error)
	UpdateStack(stackName, changeSetName, description string, templateData manager.TemplateData, parameters map[string]string) error
	NewTasksToCreateIAMServiceAccounts(serviceAccounts []*api.ClusterIAMServiceAccount, oidc *iamoidc.OpenIDConnectManager, clientSetGetter kubernetes.ClientSetGetter) *tasks.TaskTree
}

func New(clusterName string, clusterProvider *eks.ClusterProvider, stackManager StackManager, oidcManager *iamoidc.OpenIDConnectManager, clientSet kubeclient.Interface) *Manager {
	return &Manager{
		clusterName:     clusterName,
		clusterProvider: clusterProvider,
		oidcManager:     oidcManager,
		stackManager:    stackManager,
		clientSet:       clientSet,
	}
}
