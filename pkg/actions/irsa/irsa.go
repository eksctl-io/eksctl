package irsa

import (
	"fmt"

	"github.com/kris-nova/logger"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
	kubeclient "k8s.io/client-go/kubernetes"
)

type Manager struct {
	clusterName  string
	oidcManager  *iamoidc.OpenIDConnectManager
	stackManager manager.StackManager
	clientSet    kubeclient.Interface
}

func New(clusterName string, stackManager manager.StackManager, oidcManager *iamoidc.OpenIDConnectManager, clientSet kubeclient.Interface) *Manager {
	return &Manager{
		clusterName:  clusterName,
		oidcManager:  oidcManager,
		stackManager: stackManager,
		clientSet:    clientSet,
	}
}

func doTasks(taskTree *tasks.TaskTree) error {
	logger.Info(taskTree.Describe())
	if errs := taskTree.DoAllSync(); len(errs) > 0 {
		logger.Info("%d error(s) occurred and IAM Role stacks haven't been updated properly, you may wish to check CloudFormation console", len(errs))
		for _, err := range errs {
			logger.Critical("%s\n", err.Error())
		}
		return fmt.Errorf("failed to create iamserviceaccount(s)")
	}
	return nil
}

// logPlanModeWarning will log a message to inform user that they are in plan-mode
func logPlanModeWarning(plan bool) {
	if plan {
		logger.Warning("no changes were applied, run again with '--approve' to apply the changes")
	}
}
