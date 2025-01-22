package addon

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"k8s.io/client-go/kubernetes"

	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

func CreateAddonTasks(ctx context.Context, cfg *api.ClusterConfig, clusterProvider *eks.ClusterProvider, iamRoleCreator IAMRoleCreator, forceAll bool, timeout time.Duration, region string) (*tasks.TaskTree, *tasks.TaskTree, *tasks.GenericTask, []string) {
	var addons []*api.Addon
	var autoDefaultAddonNames []string
	if !cfg.AddonsConfig.DisableDefaultAddons {
		addons = make([]*api.Addon, len(cfg.Addons))
		copy(addons, cfg.Addons)

		for addonName, addonInfo := range api.KnownAddons {
			if addonInfo.IsDefault && !slices.Contains(addonInfo.ExcludedRegions, region) && !slices.ContainsFunc(cfg.Addons, func(a *api.Addon) bool {
				return strings.EqualFold(a.Name, addonName)
			}) {
				if !cfg.IsAutoModeEnabled() || addonInfo.IsDefaultAutoMode {
					addons = append(addons, &api.Addon{Name: addonName})
					autoDefaultAddonNames = append(autoDefaultAddonNames, addonName)
				}

			}
		}
	} else {
		addons = cfg.Addons
		if cfg.IsAutoModeEnabled() && len(cfg.NodeGroups) == 0 && len(cfg.ManagedNodeGroups) == 0 {
			logger.Info("default EKS addons are not required for a cluster using Auto Mode; " +
				"if nodegroups are not required, consider setting `addonsConfig.disableDefaultAddons: true` during " +
				"cluster creation, or deleting default addons using `eksctl delete addon`")
		}
	}

	var (
		preAddons  []*api.Addon
		postAddons []*api.Addon
	)
	var vpcCNIAddon *api.Addon
	for _, addon := range addons {
		addonInfo, ok := api.KnownAddons[addon.Name]
		if ok && addonInfo.CreateBeforeNodeGroup {
			preAddons = append(preAddons, addon)
		} else {
			postAddons = append(postAddons, addon)
		}
		if addon.Name == api.VPCCNIAddon {
			vpcCNIAddon = addon
		}
	}
	preTasks := &tasks.TaskTree{Parallel: false}
	postTasks := &tasks.TaskTree{Parallel: false}

	makeAddonTask := func(addons []*api.Addon, wait bool) *createAddonTask {
		return &createAddonTask{
			info:            "create addons",
			addons:          addons,
			ctx:             ctx,
			cfg:             cfg,
			clusterProvider: clusterProvider,
			forceAll:        forceAll,
			timeout:         timeout,
			wait:            wait,
			iamRoleCreator:  iamRoleCreator,
		}
	}

	if len(preAddons) > 0 {
		preTasks.Append(makeAddonTask(preAddons, false))
	}
	if len(postAddons) > 0 {
		postTasks.Append(makeAddonTask(postAddons, cfg.HasNodes()))
	}
	var updateVPCCNI *tasks.GenericTask
	if vpcCNIAddon != nil && api.IsEnabled(cfg.IAM.WithOIDC) {
		updateVPCCNI = &tasks.GenericTask{
			Description: "update VPC CNI to use IRSA if required",
			Doer: func() error {
				addonManager, err := createAddonManager(ctx, clusterProvider, cfg)
				if err != nil {
					return err
				}
				// VPC CNI is being created in a separate task, so we need to wait for it to be active before updating to use IRSA
				if err := addonManager.waitForAddonToBeActive(ctx, &api.Addon{Name: api.VPCCNIAddon}, api.DefaultWaitTimeout); err != nil {
					return fmt.Errorf("waiting for %q to become active: %w", api.VPCCNIAddon, err)
				}
				return addonManager.Update(ctx, vpcCNIAddon, nil, clusterProvider.AWSProvider.WaitTimeout())
			},
		}
	}
	return preTasks, postTasks, updateVPCCNI, autoDefaultAddonNames
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
	addonManager, err := createAddonManager(t.ctx, t.clusterProvider, t.cfg)
	if err != nil {
		return err
	}

	addonManager.DisableAWSNodePatch = true
	// always install EKS Pod Identity Agent Addon first, if present,
	// as other addons might require IAM permissions
	for _, a := range t.addons {
		if a.CanonicalName() != api.PodIdentityAgentAddon {
			continue
		}
		if t.forceAll {
			a.Force = true
		}
		if !t.wait {
			t.timeout = 0
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
		if !t.wait {
			t.timeout = 0
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

func createAddonManager(ctx context.Context, clusterProvider *eks.ClusterProvider, cfg *api.ClusterConfig) (*Manager, error) {
	var (
		oidc               *iamoidc.OpenIDConnectManager
		oidcProviderExists bool
	)
	if api.IsEnabled(cfg.IAM.WithOIDC) {
		var err error
		oidc, err = clusterProvider.NewOpenIDConnectManager(ctx, cfg)
		if err != nil {
			return nil, err
		}
		oidcProviderExists, err = oidc.CheckProviderExists(ctx)
		if err != nil {
			return nil, err
		}
	}

	stackManager := clusterProvider.NewStackManager(cfg)

	return New(cfg, clusterProvider.AWSProvider.EKS(), stackManager, oidcProviderExists, oidc, func() (kubernetes.Interface, error) {
		return clusterProvider.NewStdClientSet(cfg)
	})
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

func runAllTasks(taskTree *tasks.TaskTree) error {
	logger.Debug(taskTree.Describe())
	if errs := taskTree.DoAllSync(); len(errs) > 0 {
		var allErrs []string
		for _, err := range errs {
			allErrs = append(allErrs, err.Error())
		}
		return fmt.Errorf("%s", strings.Join(allErrs, "\n"))
	}
	completedAction := func() string {
		if taskTree.PlanMode {
			return "skipped"
		}
		return "completed successfully"
	}
	logger.Debug("all tasks were %s", completedAction())
	return nil
}
