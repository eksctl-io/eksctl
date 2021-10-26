package addon

import (
	"strings"
	"time"

	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

func CreateAddonTasks(cfg *api.ClusterConfig, clusterProvider *eks.ClusterProvider, forceAll bool, timeout time.Duration) (*tasks.TaskTree, *tasks.TaskTree) {
	preTasks := &tasks.TaskTree{Parallel: false}
	postTasks := &tasks.TaskTree{Parallel: false}
	var preAddons []*api.Addon
	var postAddons []*api.Addon
	for _, addon := range cfg.Addons {
		if strings.ToLower(addon.Name) == "vpc-cni" {
			preAddons = append(preAddons, addon)
		} else {
			postAddons = append(postAddons, addon)
		}
	}

	preTasks.Append(
		&createAddonTask{
			info:            "create addons",
			addons:          preAddons,
			cfg:             cfg,
			clusterProvider: clusterProvider,
			forceAll:        forceAll,
			timeout:         timeout,
			wait:            false,
		},
	)

	postTasks.Append(
		&createAddonTask{
			info:            "create addons",
			addons:          postAddons,
			cfg:             cfg,
			clusterProvider: clusterProvider,
			forceAll:        forceAll,
			timeout:         timeout,
			wait:            len(cfg.NodeGroups) > 0 || len(cfg.ManagedNodeGroups) > 0,
		},
	)
	return preTasks, postTasks
}

type createAddonTask struct {
	info            string
	cfg             *api.ClusterConfig
	clusterProvider *eks.ClusterProvider
	addons          []*api.Addon
	forceAll, wait  bool
	timeout         time.Duration
}

func (t *createAddonTask) Describe() string { return t.info }

func (t *createAddonTask) Do(errorCh chan error) error {
	clientSet, err := t.clusterProvider.NewStdClientSet(t.cfg)
	if err != nil {
		return err
	}

	if err := t.clusterProvider.WaitForControlPlane(t.cfg.Metadata, clientSet); err != nil {
		return errors.Wrap(err, "failed to wait for control plane")
	}

	oidc, err := t.clusterProvider.NewOpenIDConnectManager(t.cfg)
	if err != nil {
		return err
	}

	oidcProviderExists, err := oidc.CheckProviderExists()
	if err != nil {
		return err
	}

	stackManager := t.clusterProvider.NewStackManager(t.cfg)

	addonManager, err := New(t.cfg, t.clusterProvider.Provider.EKS(), stackManager, oidcProviderExists, oidc, clientSet, t.timeout)
	if err != nil {
		return err
	}

	for _, a := range t.addons {
		if t.forceAll {
			a.Force = true
		}
		err := addonManager.Create(a, t.wait)
		if err != nil {
			go func() {
				errorCh <- err
			}()
			return err
		}
	}

	go func() {
		errorCh <- nil
	}()
	return nil
}
