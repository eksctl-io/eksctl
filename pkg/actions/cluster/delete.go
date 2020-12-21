package cluster

import (
	"context"
	"fmt"
	"time"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/elb"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	ssh "github.com/weaveworks/eksctl/pkg/ssh/client"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"

	"github.com/kris-nova/logger"
)

func deleteCommon(cfg *api.ClusterConfig, ctl *eks.ClusterProvider, clientSet kubernetes.Interface, waitTimeout time.Duration) error {
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
