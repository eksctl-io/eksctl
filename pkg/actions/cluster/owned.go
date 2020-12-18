package cluster

import (
	"context"
	"time"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/elb"
	"github.com/weaveworks/eksctl/pkg/fargate"
	"github.com/weaveworks/eksctl/pkg/gitops"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	ssh "github.com/weaveworks/eksctl/pkg/ssh/client"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

type OwnedCluster struct {
	cfg          *api.ClusterConfig
	ctl          *eks.ClusterProvider
	stackManager *manager.StackCollection
}

func NewOwnedCluster(cfg *api.ClusterConfig, ctl *eks.ClusterProvider, stackManager *manager.StackCollection) (*OwnedCluster, error) {
	return &OwnedCluster{
		cfg:          cfg,
		ctl:          ctl,
		stackManager: stackManager,
	}, nil
}

func (c *OwnedCluster) Upgrade(dryRun bool) error {
	if err := c.ctl.LoadClusterVPC(c.cfg); err != nil {
		return errors.Wrapf(err, "getting VPC configuration for cluster %q", c.cfg.Metadata.Name)
	}

	versionUpdateRequired, err := upgrade(c.cfg, c.ctl, dryRun)
	if err != nil {
		return err
	}

	if err := c.ctl.RefreshClusterStatus(c.cfg); err != nil {
		return err
	}

	supportsManagedNodes, err := c.ctl.SupportsManagedNodes(c.cfg)
	if err != nil {
		return err
	}

	stackUpdateRequired, err := c.stackManager.AppendNewClusterStackResource(dryRun, supportsManagedNodes)
	if err != nil {
		return err
	}

	if err := c.ctl.ValidateExistingNodeGroupsForCompatibility(c.cfg, c.stackManager); err != nil {
		logger.Critical("failed checking nodegroups", err.Error())
	}

	cmdutils.LogPlanModeWarning(dryRun && (stackUpdateRequired || versionUpdateRequired))
	return nil
}

func (c *OwnedCluster) Delete(waitTimeout time.Duration, wait bool) error {
	var (
		clientSet kubernetes.Interface
		oidc      *iamoidc.OpenIDConnectManager
		err       error
	)

	clusterOperable, _ := c.ctl.CanOperate(c.cfg)
	oidcSupported := true
	if clusterOperable {
		clientSet, err = c.ctl.NewStdClientSet(c.cfg)
		if err != nil {
			return err
		}

		oidc, err = c.ctl.NewOpenIDConnectManager(c.cfg)
		if err != nil {
			if _, ok := err.(*eks.UnsupportedOIDCError); !ok {
				return err
			}
			oidcSupported = false
		}
		if err := deleteFargateProfiles(c.cfg.Metadata, waitTimeout, c.ctl); err != nil {
			return err
		}
	}

	stackManager := c.ctl.NewStackManager(c.cfg)

	ssh.DeleteKeys(c.cfg.Metadata.Name, c.ctl.Provider.EC2())

	kubeconfig.MaybeDeleteConfig(c.cfg.Metadata)

	if hasDeprecatedStacks, err := deleteDeprecatedStacks(stackManager); hasDeprecatedStacks {
		if err != nil {
			return err
		}
		return nil
	}

	{
		// only need to cleanup ELBs if the cluster has already been created.
		if clusterOperable {
			ctx, cleanup := context.WithTimeout(context.Background(), 10*time.Minute)
			defer cleanup()

			logger.Info("cleaning up AWS load balancers created by Kubernetes objects of Kind Service or Ingress")
			if err := elb.Cleanup(ctx, c.ctl.Provider.EC2(), c.ctl.Provider.ELB(), c.ctl.Provider.ELBV2(), clientSet, c.cfg); err != nil {
				return err
			}
		}

		deleteOIDCProvider := clusterOperable && oidcSupported
		tasks, err := stackManager.NewTasksToDeleteClusterWithNodeGroups(deleteOIDCProvider, oidc, kubernetes.NewCachedClientSet(clientSet), wait, func(errs chan error, _ string) error {
			logger.Info("trying to cleanup dangling network interfaces")
			if err := c.ctl.LoadClusterVPC(c.cfg); err != nil {
				return errors.Wrapf(err, "getting VPC configuration for cluster %q", c.cfg.Metadata.Name)
			}

			go func() {
				errs <- vpc.CleanupNetworkInterfaces(c.ctl.Provider.EC2(), c.cfg)
				close(errs)
			}()
			return nil
		})

		if err != nil {
			return err
		}

		if tasks.Len() == 0 {
			logger.Warning("no cluster resources were found for %q", c.cfg.Metadata.Name)
			return nil
		}

		logger.Info(tasks.Describe())
		if errs := tasks.DoAllSync(); len(errs) > 0 {
			return handleErrors(errs, "cluster with nodegroup(s)")
		}

		logger.Success("all cluster resources were deleted")
	}

	{
		if err := gitops.DeleteKey(c.cfg); err != nil {
			return err
		}
	}

	return nil
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
