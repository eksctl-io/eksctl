package cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/fargate"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/elb"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	ssh "github.com/weaveworks/eksctl/pkg/ssh/client"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"

	"github.com/kris-nova/logger"
)

func deleteSharedResources(cfg *api.ClusterConfig, ctl *eks.ClusterProvider, clientSet kubernetes.Interface, waitTimeout time.Duration) error {
	clusterOperable, _ := ctl.CanOperate(cfg)
	if clusterOperable {
		if err := deleteFargateProfiles(cfg.Metadata, waitTimeout, ctl); err != nil {
			return err
		}
	}

	stackManager := ctl.NewStackManager(cfg)

	if hasDeprecatedStacks, err := deleteDeprecatedStacks(stackManager); hasDeprecatedStacks {
		if err != nil {
			return err
		}
		return nil
	}

	ssh.DeleteKeys(cfg.Metadata.Name, ctl.Provider.EC2())

	kubeconfig.MaybeDeleteConfig(cfg.Metadata)

	// only need to cleanup ELBs if the cluster has already been created.
	if clusterOperable {
		ctx, cleanup := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cleanup()

		logger.Info("cleaning up AWS load balancers created by Kubernetes objects of Kind Service or Ingress")
		if err := elb.Cleanup(ctx, ctl.Provider.EC2(), ctl.Provider.ELB(), ctl.Provider.ELBV2(), clientSet, cfg); err != nil {
			return err
		}
	}
	logger.Success("all cluster resources were deleted")
	return nil
}

func handleErrors(errs []error, subject string) error {
	logger.Info("%d error(s) occurred while deleting %s", len(errs), subject)
	for _, err := range errs {
		logger.Critical("%s\n", err.Error())
	}
	return fmt.Errorf("failed to delete %s", subject)
}

func deleteFargateProfiles(clusterMeta *api.ClusterMeta, waitTimeout time.Duration, ctl *eks.ClusterProvider) error {
	awsClient := fargate.NewClientWithWaitTimeout(
		clusterMeta.Name,
		ctl.Provider.EKS(),
		waitTimeout,
	)
	profileNames, err := awsClient.ListProfiles()
	if err != nil {
		if fargate.IsUnauthorizedError(err) {
			logger.Debug("Fargate: unauthorized error: %v", err)
			logger.Info("either account is not authorized to use Fargate or region %s is not supported. Ignoring error",
				clusterMeta.Region)
			return nil
		}
		return err
	}

	// Linearise the deleting of Fargate profiles by passing as the API
	// otherwise errors out with:
	//   ResourceInUseException: Cannot delete Fargate Profile ${name2} because
	//   cluster ${clusterName} currently has Fargate profile ${name1} in
	//   status DELETING

	for _, profileName := range profileNames {
		logger.Info("deleting Fargate profile %q", *profileName)
		// All Fargate profiles must be completely deleted by waiting for the deletion to complete, before deleting
		// the cluster itself, otherwise it can result in this error:
		//   Cannot delete because cluster <cluster> currently has Fargate profile <profile> in status DELETING
		if err := awsClient.DeleteProfile(*profileName, true); err != nil {
			return err
		}
		logger.Info("deleted Fargate profile %q", *profileName)
	}
	logger.Info("deleted %v Fargate profile(s)", len(profileNames))
	return nil
}

func deleteDeprecatedStacks(stackManager *manager.StackCollection) (bool, error) {
	tasks, err := stackManager.DeleteTasksForDeprecatedStacks()
	if err != nil {
		return true, err
	}
	if count := tasks.Len(); count > 0 {
		logger.Info(tasks.Describe())
		if errs := tasks.DoAllSync(); len(errs) > 0 {
			return true, handleErrors(errs, "deprecated stacks")
		}
		logger.Success("deleted all %s deperecated stacks", count)
		return true, nil
	}
	return false, nil
}
