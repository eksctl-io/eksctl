package cluster

import (
	"fmt"
	"time"

	"github.com/weaveworks/eksctl/pkg/cfn/manager"

	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"

	"github.com/pkg/errors"

	"github.com/weaveworks/eksctl/pkg/utils/waiters"

	"github.com/weaveworks/eksctl/pkg/kubernetes"

	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"

	"github.com/weaveworks/logger"

	awseks "github.com/aws/aws-sdk-go/service/eks"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	kubeclient "k8s.io/client-go/kubernetes"
)

type UnownedCluster struct {
	cfg          *api.ClusterConfig
	ctl          *eks.ClusterProvider
	clientSet    kubeclient.Interface
	stackManager manager.StackManager
}

func NewUnownedCluster(cfg *api.ClusterConfig, ctl *eks.ClusterProvider, clientSet kubeclient.Interface, stackManager manager.StackManager) *UnownedCluster {
	return &UnownedCluster{
		cfg:          cfg,
		ctl:          ctl,
		clientSet:    clientSet,
		stackManager: stackManager,
	}
}

func (c *UnownedCluster) Upgrade(dryRun bool) error {
	versionUpdateRequired, err := upgrade(c.cfg, c.ctl, dryRun)
	if err != nil {
		return err
	}

	// if no version update is required, don't log asking them to rerun with --approve
	cmdutils.LogPlanModeWarning(dryRun && versionUpdateRequired)
	return nil
}

func (c *UnownedCluster) Delete(waitInterval time.Duration, wait bool) error {
	clusterName := c.cfg.Metadata.Name

	if err := c.checkClusterExists(clusterName); err != nil {
		return err
	}

	if err := deleteSharedResources(c.cfg, c.ctl, c.clientSet); err != nil {
		return err
	}

	if err := c.deleteFargateRoleIfExists(); err != nil {
		return err
	}

	// we have to wait for nodegroups to delete before deleting the cluster
	// so the `wait` value is ignored here
	if err := c.deleteAndWaitForNodegroupsDeletion(clusterName, c.ctl.Provider.WaitTimeout(), waitInterval); err != nil {
		return err
	}

	if err := c.deleteIAMAndOIDC(wait, kubernetes.NewCachedClientSet(c.clientSet)); err != nil {
		return err
	}

	err := c.deleteCluster(clusterName, c.ctl.Provider.WaitTimeout(), wait)
	if err != nil {
		return err
	}

	if err := checkForUndeletedStacks(c.ctl.NewStackManager(c.cfg)); err != nil {
		return err
	}

	logger.Success("all cluster resources were deleted")
	return nil
}

func (c *UnownedCluster) deleteFargateRoleIfExists() error {
	stack, err := c.stackManager.GetFargateStack()
	if err != nil {
		return err
	}

	if stack != nil {
		logger.Info("deleting fargate role")
		_, err = c.stackManager.DeleteStackBySpec(stack)
		return err
	}

	logger.Debug("no fargate role found")
	return nil
}

func (c *UnownedCluster) checkClusterExists(clusterName string) error {
	_, err := c.ctl.Provider.EKS().DescribeCluster(&awseks.DescribeClusterInput{
		Name: &c.cfg.Metadata.Name,
	})
	if err != nil {
		if isNotFound(err) {
			return errors.Errorf("cluster %q not found", clusterName)
		}
		return errors.Wrapf(err, "error describing cluster %q", clusterName)
	}
	return nil
}

func (c *UnownedCluster) deleteIAMAndOIDC(wait bool, clientSetGetter kubernetes.ClientSetGetter) error {
	var oidc *iamoidc.OpenIDConnectManager

	clusterOperable, _ := c.ctl.CanOperate(c.cfg)
	oidcSupported := true

	if clusterOperable {
		var err error
		oidc, err = c.ctl.NewOpenIDConnectManager(c.cfg)
		if err != nil {
			if _, ok := err.(*eks.UnsupportedOIDCError); !ok {
				return err
			}
			oidcSupported = false
		}
	}

	tasksTree := &tasks.TaskTree{Parallel: false}

	if clusterOperable && oidcSupported {
		serviceAccountAndOIDCTasks, err := c.stackManager.NewTasksToDeleteOIDCProviderWithIAMServiceAccounts(oidc, clientSetGetter)
		if err != nil {
			return err
		}

		if serviceAccountAndOIDCTasks.Len() > 0 {
			serviceAccountAndOIDCTasks.IsSubTask = true
			tasksTree.Append(serviceAccountAndOIDCTasks)
		}
	}

	deleteAddonIAMtasks, err := c.stackManager.NewTaskToDeleteAddonIAM(wait)
	if err != nil {
		return err
	}

	if deleteAddonIAMtasks.Len() > 0 {
		deleteAddonIAMtasks.IsSubTask = true
		tasksTree.Append(deleteAddonIAMtasks)
	}

	if tasksTree.Len() == 0 {
		logger.Warning("no IAM and OIDC resources were found for %q", c.cfg.Metadata.Name)
		return nil
	}

	logger.Info(tasksTree.Describe())
	if errs := tasksTree.DoAllSync(); len(errs) > 0 {
		return handleErrors(errs, "deleting cluster IAM and OIDC")
	}

	logger.Info("all IAM and OIDC resources were deleted")
	return nil
}

func (c *UnownedCluster) deleteCluster(clusterName string, waitTimeout time.Duration, wait bool) error {
	out, err := c.ctl.Provider.EKS().DeleteCluster(&awseks.DeleteClusterInput{
		Name: &clusterName,
	})

	if err != nil {
		return err
	}

	logger.Info("initiated deletion of cluster %q", clusterName)
	if out != nil {
		logger.Debug("delete cluster response: %s", out.String())
	}

	if !wait {
		logger.Info("to see the status of the deletion run `eksctl get cluster --name %s --region %s`", clusterName, c.cfg.Metadata.Region)
		return nil
	}
	newRequest := func() *request.Request {
		input := &awseks.DescribeClusterInput{
			Name: &clusterName,
		}
		req, _ := c.ctl.Provider.EKS().DescribeClusterRequest(input)
		return req
	}

	acceptors := waiters.MakeErrorCodeAcceptors(awseks.ErrCodeResourceNotFoundException)

	msg := fmt.Sprintf("waiting for cluster %q to be deleted", clusterName)

	return waiters.Wait(clusterName, msg, acceptors, newRequest, waitTimeout, nil)
}

func (c *UnownedCluster) deleteAndWaitForNodegroupsDeletion(clusterName string, waitTimeout, waitInterval time.Duration) error {
	var nodegroups []*string

	pager := func(ng *awseks.ListNodegroupsOutput, _ bool) bool {
		nodegroups = append(nodegroups, ng.Nodegroups...)
		return true
	}

	err := c.ctl.Provider.EKS().ListNodegroupsPages(&awseks.ListNodegroupsInput{
		ClusterName: &clusterName,
	}, pager)

	if err != nil {
		return err
	}

	if len(nodegroups) == 0 {
		logger.Info("no nodegroups to delete")
		return nil
	}

	for _, nodeGroupName := range nodegroups {
		out, err := c.ctl.Provider.EKS().DeleteNodegroup(&awseks.DeleteNodegroupInput{
			ClusterName:   &clusterName,
			NodegroupName: nodeGroupName,
		})

		if err != nil {
			return err
		}
		logger.Info("initiated deletion of nodegroup %q", *nodeGroupName)
		logger.Debug("delete nodegroup %q response: %s", *nodeGroupName, out.String())
	}

	condition := func() (bool, error) {
		nodeGroups, err := c.ctl.Provider.EKS().ListNodegroups(&awseks.ListNodegroupsInput{
			ClusterName: &clusterName,
		})
		if err != nil {
			return false, err
		}
		if len(nodeGroups.Nodegroups) == 0 {
			logger.Info("all nodegroups for cluster %q successfully deleted", clusterName)
			return true, nil
		}

		logger.Info("waiting for nodegroups to be deleted, %d remaining", len(nodeGroups.Nodegroups))
		return false, nil
	}

	return waiters.WaitForCondition(waitTimeout, waitInterval, fmt.Errorf("timed out waiting for nodegroup deletion after %s", waitTimeout), condition)
}

func isNotFound(err error) bool {
	awsError, ok := err.(awserr.Error)
	return ok && awsError.Code() == awseks.ErrCodeResourceNotFoundException
}
