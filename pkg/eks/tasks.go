package eks

import (
	"context"
	"fmt"
	"time"

	"github.com/weaveworks/eksctl/pkg/actions/identityproviders"
	"github.com/weaveworks/eksctl/pkg/actions/irsa"

	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/logger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	"github.com/weaveworks/eksctl/pkg/addons"
	"github.com/weaveworks/eksctl/pkg/fargate"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/utils"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
)

type clusterConfigTask struct {
	info string
	spec *api.ClusterConfig
	call func(*api.ClusterConfig) error
}

func (t *clusterConfigTask) Describe() string { return t.info }

func (t *clusterConfigTask) Do(errs chan error) error {
	err := t.call(t.spec)
	close(errs)
	return err
}

// VPCControllerTask represents a task to install the VPC controller
type VPCControllerTask struct {
	Info            string
	ClusterProvider *ClusterProvider
	ClusterConfig   *api.ClusterConfig
	PlanMode        bool
}

// Describe implements Task
func (v *VPCControllerTask) Describe() string { return v.Info }

// Do implements Task
func (v *VPCControllerTask) Do(errCh chan error) error {
	defer close(errCh)
	rawClient, err := v.ClusterProvider.NewRawClient(v.ClusterConfig)
	if err != nil {
		return err
	}
	oidc, err := v.ClusterProvider.NewOpenIDConnectManager(v.ClusterConfig)
	if err != nil {
		return err
	}

	stackCollection := manager.NewStackCollection(v.ClusterProvider.Provider, v.ClusterConfig)

	clientSet, err := v.ClusterProvider.NewStdClientSet(v.ClusterConfig)
	if err != nil {
		return err
	}
	irsaManager := irsa.New(v.ClusterConfig.Metadata.Name, stackCollection, oidc, clientSet)
	irsa := addons.NewIRSAHelper(oidc, stackCollection, irsaManager, v.ClusterConfig.Metadata.Name)

	// TODO PlanMode doesn't work as intended
	vpcController := addons.NewVPCController(rawClient, irsa, v.ClusterConfig.Status, v.ClusterProvider.Provider.Region(), v.PlanMode)
	if err := vpcController.Deploy(); err != nil {
		return errors.Wrap(err, "error installing VPC controller")
	}
	return nil
}

type devicePluginTask struct {
	info            string
	clusterProvider *ClusterProvider
	spec            *api.ClusterConfig
	mkPlugin        addons.MkDevicePlugin
}

func (n *devicePluginTask) Describe() string { return n.info }

func (n *devicePluginTask) Do(errCh chan error) error {
	defer close(errCh)
	rawClient, err := n.clusterProvider.NewRawClient(n.spec)
	if err != nil {
		return err
	}
	neuronDevicePlugin := n.mkPlugin(rawClient, n.clusterProvider.Provider.Region(), false)
	if err := neuronDevicePlugin.Deploy(); err != nil {
		return errors.Wrap(err, "error installing Neuron device plugin")
	}
	return nil
}

func newNvidiaDevicePluginTask(
	clusterProvider *ClusterProvider,
	spec *api.ClusterConfig,
) tasks.Task {
	t := devicePluginTask{
		info:            "install Nvidia device plugin",
		clusterProvider: clusterProvider,
		spec:            spec,
		mkPlugin:        addons.NewNvidiaDevicePlugin,
	}
	return &t
}

func newNeuronDevicePluginTask(
	clusterProvider *ClusterProvider,
	spec *api.ClusterConfig,
) tasks.Task {
	t := devicePluginTask{
		info:            "install Neuron device plugin",
		clusterProvider: clusterProvider,
		spec:            spec,
		mkPlugin:        addons.NewNeuronDevicePlugin,
	}
	return &t
}

type restartDaemonsetTask struct {
	name            string
	namespace       string
	clusterProvider *ClusterProvider
	spec            *api.ClusterConfig
}

func (t *restartDaemonsetTask) Describe() string {
	return fmt.Sprintf(`restart daemonset "%s/%s"`, t.namespace, t.name)
}

func (t *restartDaemonsetTask) Do(errCh chan error) error {
	defer close(errCh)
	clientSet, err := t.clusterProvider.NewStdClientSet(t.spec)
	if err != nil {
		return err
	}
	ds := clientSet.AppsV1().DaemonSets(t.namespace)
	dep, err := ds.Get(context.TODO(), t.name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if dep.Spec.Template.Annotations == nil {
		dep.Spec.Template.Annotations = make(map[string]string)
	}
	dep.Spec.Template.Annotations["eksctl.io/restartedAt"] = time.Now().Format(time.RFC3339)
	bytes, err := runtime.Encode(unstructured.UnstructuredJSONScheme, dep)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal %q deployment", t.name)
	}
	if _, err := ds.Patch(context.TODO(), t.name, types.MergePatchType, bytes, metav1.PatchOptions{}); err != nil {
		return errors.Wrap(err, "failed to patch deployment")
	}
	logger.Info(`daemonset "%s/%s" restarted`, t.namespace, t.name)
	return nil
}

// CreateExtraClusterConfigTasks returns all tasks for updating cluster configuration not depending on the control plane availability
func (c *ClusterProvider) CreateExtraClusterConfigTasks(cfg *api.ClusterConfig, installVPCController bool) *tasks.TaskTree {
	newTasks := &tasks.TaskTree{
		Parallel:  false,
		IsSubTask: true,
	}

	newTasks.Append(&tasks.GenericTask{
		Description: "wait for control plane to become ready",
		Doer: func() error {
			clientSet, err := c.NewStdClientSet(cfg)
			if err != nil {
				return errors.Wrap(err, "error creating Clientset")
			}
			return c.WaitForControlPlane(cfg.Metadata, clientSet)
		},
	})

	if len(cfg.Metadata.Tags) > 0 {
		newTasks.Append(&clusterConfigTask{
			info: "tag cluster",
			spec: cfg,
			call: c.UpdateClusterTags,
		})
	}
	if !cfg.HasClusterCloudWatchLogging() {
		logger.Info("CloudWatch logging will not be enabled for cluster %q in %q", cfg.Metadata.Name, cfg.Metadata.Region)
		logger.Info("you can enable it with 'eksctl utils update-cluster-logging --enable-types={SPECIFY-YOUR-LOG-TYPES-HERE (e.g. all)} --region=%s --cluster=%s'", cfg.Metadata.Region, cfg.Metadata.Name)

	} else {
		newTasks.Append(&clusterConfigTask{
			info: "update CloudWatch logging configuration",
			spec: cfg,
			call: c.UpdateClusterConfigForLogging,
		})
	}
	c.maybeAppendTasksForEndpointAccessUpdates(cfg, newTasks)

	if len(cfg.VPC.PublicAccessCIDRs) > 0 {
		newTasks.Append(&clusterConfigTask{
			info: "update public access CIDRs",
			spec: cfg,
			call: c.UpdatePublicAccessCIDRs,
		})
	}

	if cfg.IsFargateEnabled() {
		manager := fargate.NewFromProvider(cfg.Metadata.Name, c.Provider)
		newTasks.Append(&fargateProfilesTask{
			info:            "create fargate profiles",
			spec:            cfg,
			clusterProvider: c,
			manager:         &manager,
		})
	}

	if api.IsEnabled(cfg.IAM.WithOIDC) {
		c.appendCreateTasksForIAMServiceAccounts(cfg, newTasks)
	}

	if len(cfg.IdentityProviders) > 0 {
		newTasks.Append(identityproviders.NewAssociateProvidersTask(*cfg.Metadata, cfg.IdentityProviders, c.Provider.EKS()))
	}

	if installVPCController {
		newTasks.Append(&VPCControllerTask{
			Info:            "install Windows VPC controller",
			ClusterConfig:   cfg,
			ClusterProvider: c,
		})
	}

	return newTasks
}

// ClusterTasksForNodeGroups returns all tasks dependent on node groups
func (c *ClusterProvider) ClusterTasksForNodeGroups(cfg *api.ClusterConfig, installNeuronDevicePluginParam, installNvidiaDevicePluginParam bool) *tasks.TaskTree {
	tasks := &tasks.TaskTree{
		Parallel:  true,
		IsSubTask: false,
	}
	var haveNeuronInstanceType bool
	for _, ng := range cfg.NodeGroups {
		haveNeuronInstanceType = haveNeuronInstanceType || api.HasInstanceType(ng, utils.IsInferentiaInstanceType)
	}
	var haveNvidiaInstanceType bool
	for _, ng := range cfg.NodeGroups {
		haveNvidiaInstanceType = haveNvidiaInstanceType || api.HasInstanceType(ng, func(t string) bool {
			return utils.IsGPUInstanceType(t) && !utils.IsInferentiaInstanceType(t)
		})
	}
	if haveNeuronInstanceType && installNeuronDevicePluginParam {
		tasks.Append(newNeuronDevicePluginTask(c, cfg))
	}
	if haveNvidiaInstanceType && installNvidiaDevicePluginParam {
		tasks.Append(newNvidiaDevicePluginTask(c, cfg))
	}

	var ngs []*api.NodeGroupBase
	for _, ng := range cfg.NodeGroups {
		ngs = append(ngs, ng.NodeGroupBase)
	}
	for _, ng := range cfg.ManagedNodeGroups {
		ngs = append(ngs, ng.NodeGroupBase)
	}
	for _, ng := range ngs {
		if len(ng.ASGSuspendProcesses) > 0 {
			tasks.Append(newSuspendProcesses(c, cfg, ng))
		}
	}

	return tasks
}

func (c *ClusterProvider) appendCreateTasksForIAMServiceAccounts(cfg *api.ClusterConfig, tasks *tasks.TaskTree) {
	// we don't have all the information to construct full iamoidc.OpenIDConnectManager now,
	// instead we just create a reference that gets updated when first task runs, and gets
	// used by this would be more elegant if it was all done via CloudFormation and we didn't
	// have to put wires across all the things like this; this whole function is needed because
	// we cannot manage certain EKS features with CloudFormation
	oidcPlaceholder := &iamoidc.OpenIDConnectManager{}
	tasks.Append(&clusterConfigTask{
		info: "associate IAM OIDC provider",
		spec: cfg,
		call: func(cfg *api.ClusterConfig) error {
			if err := c.RefreshClusterStatus(cfg); err != nil {
				return err
			}
			oidc, err := c.NewOpenIDConnectManager(cfg)
			if err != nil {
				return err
			}
			if err := oidc.CreateProvider(); err != nil {
				return err
			}
			*oidcPlaceholder = *oidc
			// Make sure control plane is reachable
			clientSet, err := c.NewStdClientSet(cfg)
			if err != nil {
				return errors.Wrap(err, "failed to get ClientSet")
			}
			if err := c.WaitForControlPlane(cfg.Metadata, clientSet); err != nil {
				return errors.Wrap(err, "failed to wait for control plane")
			}
			return nil
		},
	})

	clientSet := &kubernetes.CallbackClientSet{
		Callback: func() (kubernetes.Interface, error) {
			return c.NewStdClientSet(cfg)
		},
	}

	// as this is non-CloudFormation context, we need to construct a new stackManager,
	// given a clientSet getter and OpenIDConnectManager reference we can build out
	// the list of tasks for each of the service accounts that need to be created
	newTasks := c.NewStackManager(cfg).NewTasksToCreateIAMServiceAccounts(
		api.IAMServiceAccountsWithImplicitServiceAccounts(cfg),
		oidcPlaceholder,
		clientSet,
	)
	newTasks.IsSubTask = true
	tasks.Append(newTasks)
	tasks.Append(&restartDaemonsetTask{
		namespace:       "kube-system",
		name:            "aws-node",
		clusterProvider: c,
		spec:            cfg,
	})
}

func (c *ClusterProvider) maybeAppendTasksForEndpointAccessUpdates(cfg *api.ClusterConfig, tasks *tasks.TaskTree) {
	// if a cluster config doesn't have the default api endpoint access, append a new task
	// so that we update the cluster with the new access configuration.  This is a
	// non-CloudFormation context, so we create a task to send it through the EKS API.
	// A caveat is that sending the default endpoint parameters for a cluster as an update will
	// return an error from the EKS API, so we must check for this before sending the request.
	if cfg.HasClusterEndpointAccess() && api.EndpointsEqual(*cfg.VPC.ClusterEndpoints, *api.ClusterEndpointAccessDefaults()) {
		// No tasks to append here as there's no updates to make.
		logger.Info(cfg.DefaultEndpointsMsg())
	} else {
		logger.Info(cfg.CustomEndpointsMsg())

		tasks.Append(&clusterConfigTask{
			info: "update cluster VPC endpoint access configuration",
			spec: cfg,
			call: c.UpdateClusterConfigForEndpoints,
		})
	}
}
