package iam

import (
	"fmt"

	"github.com/kris-nova/logger"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
)

func (a *Manager) CreateIAMServiceAccount(iamServiceAccounts []*api.ClusterIAMServiceAccount, plan bool) error {
	tasks := a.stackManager.NewTasksToCreateIAMServiceAccounts(iamServiceAccounts, a.oidcManager, kubernetes.NewCachedClientSet(a.clientSet))

	logger.Info(tasks.Describe())
	if errs := tasks.DoAllSync(); len(errs) > 0 {
		logger.Info("%d error(s) occurred and IAM Role stacks haven't been created properly, you may wish to check CloudFormation console", len(errs))
		for _, err := range errs {
			logger.Critical("%s\n", err.Error())
		}
		return fmt.Errorf("failed to create iamserviceaccount(s)")
	}

	cmdutils.LogPlanModeWarning(plan && len(iamServiceAccounts) > 0)

	return nil
}
