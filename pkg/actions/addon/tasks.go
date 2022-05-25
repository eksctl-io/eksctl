package addon

import (
	"context"
	"strings"
	"time"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

func CreateAddonTasks(ctx context.Context, cfg *api.ClusterConfig, clusterProvider *eks.ClusterProvider, forceAll bool, timeout time.Duration) (*tasks.TaskTree, *tasks.TaskTree) {
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
			ctx:             ctx,
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
			ctx:             ctx,
			cfg:             cfg,
			clusterProvider: clusterProvider,
			forceAll:        forceAll,
			timeout:         timeout,
			wait:            cfg.HasNodes(),
		},
	)
	return preTasks, postTasks
}

type createAddonTask struct {
	// Context should ideally be passed to methods and not be a struct field,
	// but the current task code requires it to be passed this way.
	ctx             context.Context
	info            string
	cfg             *api.ClusterConfig
	clusterProvider *eks.ClusterProvider
	addons          []*api.Addon
	forceAll, wait  bool
	timeout         time.Duration
}

func (t *createAddonTask) Describe() string { return t.info }

func (t *createAddonTask) Do(errorCh chan error) error {
	oidc, err := t.clusterProvider.NewOpenIDConnectManager(t.ctx, t.cfg)
	if err != nil {
		return err
	}

	oidcProviderExists, err := oidc.CheckProviderExists(t.ctx)
	if err != nil {
		return err
	}

	stackManager := t.clusterProvider.NewStackManager(t.cfg)

	clientSet, err := t.clusterProvider.NewStdClientSet(t.cfg)
	if err != nil {
		return err
	}
	addonManager, err := New(t.cfg, t.clusterProvider.AWSProvider.EKS(), stackManager, oidcProviderExists, oidc, clientSet)
	if err != nil {
		return err
	}

	for _, a := range t.addons {
		if t.forceAll {
			a.Force = true
		}
		err := addonManager.Create(t.ctx, a, t.timeout)
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
