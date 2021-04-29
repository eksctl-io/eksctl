package addon

import (
	"time"

	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

func CreateAddonTasks(cfg *api.ClusterConfig, clusterProvider *eks.ClusterProvider, forceAll bool, timeout time.Duration) *tasks.TaskTree {
	taskTree := &tasks.TaskTree{Parallel: false}

	taskTree.Append(
		&createAddonTask{
			info:            "create addons",
			addons:          cfg.Addons,
			cfg:             cfg,
			clusterProvider: clusterProvider,
			forceAll:        forceAll,
			timeout:         timeout,
		},
	)
	return taskTree
}

type createAddonTask struct {
	info            string
	cfg             *api.ClusterConfig
	clusterProvider *eks.ClusterProvider
	addons          []*api.Addon
	forceAll        bool
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
		err := addonManager.Create(a, true)
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
