package addon

import (
	"context"
	"fmt"
	"strings"
	"time"

	"k8s.io/client-go/kubernetes"

	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

func CreateAddonTasks(ctx context.Context, cfg *api.ClusterConfig, clusterProvider *eks.ClusterProvider, iamRoleCreator IAMRoleCreator, forceAll bool, timeout time.Duration) (*tasks.TaskTree, *tasks.TaskTree) {
	preTasks := &tasks.TaskTree{Parallel: false}
	postTasks := &tasks.TaskTree{Parallel: false}
	var preAddons []*api.Addon
	var postAddons []*api.Addon
	for _, addon := range cfg.Addons {
		if strings.EqualFold(addon.Name, api.VPCCNIAddon) ||
			strings.EqualFold(addon.Name, api.PodIdentityAgentAddon) {
			preAddons = append(preAddons, addon)
		} else {
			postAddons = append(postAddons, addon)
		}
	}

	preAddonsTask := createAddonTask{
		info:            "create addons",
		addons:          preAddons,
		ctx:             ctx,
		cfg:             cfg,
		clusterProvider: clusterProvider,
		forceAll:        forceAll,
		timeout:         timeout,
		wait:            false,
		iamRoleCreator:  iamRoleCreator,
	}

	preTasks.Append(&preAddonsTask)

	postAddonsTask := preAddonsTask
	postAddonsTask.addons = postAddons
	postAddonsTask.wait = cfg.HasNodes()
	postTasks.Append(&postAddonsTask)

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
	iamRoleCreator  IAMRoleCreator
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

	addonManager, err := New(t.cfg, t.clusterProvider.AWSProvider.EKS(), stackManager, oidcProviderExists, oidc, func() (kubernetes.Interface, error) {
		return t.clusterProvider.NewStdClientSet(t.cfg)
	})
	if err != nil {
		return err
	}

	// always install EKS Pod Identity Agent Addon first, if present,
	// as other addons might require IAM permissions
	for _, a := range t.addons {
		if a.CanonicalName() != api.PodIdentityAgentAddon {
			continue
		}
		if t.forceAll {
			a.Force = true
		}
		err := addonManager.Create(t.ctx, a, t.iamRoleCreator, t.timeout)
		if err != nil {
			go func() {
				errorCh <- err
			}()
			return err
		}
	}

	for _, a := range t.addons {
		if a.CanonicalName() == api.PodIdentityAgentAddon {
			continue
		}
		if t.forceAll {
			a.Force = true
		}
		err := addonManager.Create(t.ctx, a, t.iamRoleCreator, t.timeout)
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

type deleteAddonIAMTask struct {
	ctx          context.Context
	info         string
	stack        *cfntypes.Stack
	stackManager StackManager
	wait         bool
}

func (t *deleteAddonIAMTask) Describe() string { return t.info }

func (t *deleteAddonIAMTask) Do(errorCh chan error) error {
	errMsg := fmt.Sprintf("deleting addon IAM %q", *t.stack.StackName)
	if t.wait {
		if err := t.stackManager.DeleteStackBySpecSync(t.ctx, t.stack, errorCh); err != nil {
			return fmt.Errorf("%s: %w", errMsg, err)
		}
		return nil
	}
	defer close(errorCh)
	if _, err := t.stackManager.DeleteStackBySpec(t.ctx, t.stack); err != nil {
		return fmt.Errorf("%s: %w", errMsg, err)
	}
	return nil
}
