package iam

import (
	"fmt"

	"github.com/kris-nova/logger"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
	kubeclient "k8s.io/client-go/kubernetes"
)

type Manager struct {
	clusterName     string
	clusterProvider api.ClusterProvider
	oidcManager     *iamoidc.OpenIDConnectManager
	stackManager    StackManager
	clientSet       kubeclient.Interface
}

//go:generate counterfeiter -o fakes/fake_stack_manager.go . StackManager
type StackManager interface {
	ListStacksMatching(nameRegex string, statusFilters ...string) ([]*manager.Stack, error)
	UpdateStack(stackName, changeSetName, description string, templateData manager.TemplateData, parameters map[string]string) error
	CreateIAMServiceAccount(errs chan error, spec *api.ClusterIAMServiceAccount, oidc *iamoidc.OpenIDConnectManager) error
}

func New(clusterName string, clusterProvider api.ClusterProvider, stackManager StackManager, oidcManager *iamoidc.OpenIDConnectManager, clientSet kubeclient.Interface) *Manager {
	return &Manager{
		clusterName:     clusterName,
		clusterProvider: clusterProvider,
		oidcManager:     oidcManager,
		stackManager:    stackManager,
		clientSet:       clientSet,
	}
}

func doTasks(taskTree *tasks.TaskTree) error {
	logger.Info(taskTree.Describe())
	if errs := taskTree.DoAllSync(); len(errs) > 0 {
		logger.Info("%d error(s) occurred and IAM Role stacks haven't been created properly, you may wish to check CloudFormation console", len(errs))
		for _, err := range errs {
			logger.Critical("%s\n", err.Error())
		}
		return fmt.Errorf("failed to create iamserviceaccount(s)")
	}
	return nil
}

func LogPlanModeWarning(plan bool) {
	if plan {
		logger.Warning("no changes were applied, run again with '--approve' to apply the changes")
	}
}
