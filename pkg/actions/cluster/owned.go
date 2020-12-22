package cluster

import (
	"time"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/gitops"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
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
	}

	if err := deleteSharedResources(c.cfg, c.ctl, clientSet, waitTimeout); err != nil {
		return err
	}

	deleteOIDCProvider := clusterOperable && oidcSupported
	tasks, err := c.ctl.NewStackManager(c.cfg).NewTasksToDeleteClusterWithNodeGroups(deleteOIDCProvider, oidc, kubernetes.NewCachedClientSet(clientSet), wait, func(errs chan error, _ string) error {
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

	if err := gitops.DeleteKey(c.cfg); err != nil {
		return err
	}

	return nil
}
