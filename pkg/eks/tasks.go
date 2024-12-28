package eks

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"

	"github.com/kris-nova/logger"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/weaveworks/eksctl/pkg/actions/iamidentitymapping"
	"github.com/weaveworks/eksctl/pkg/actions/identityproviders"
	"github.com/weaveworks/eksctl/pkg/actions/irsa"
	"github.com/weaveworks/eksctl/pkg/addons"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/fargate"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	instanceutils "github.com/weaveworks/eksctl/pkg/utils/instance"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
	"github.com/weaveworks/eksctl/pkg/windows"
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

// WindowsIPAMTask is a task for enabling Windows IPAM.
type WindowsIPAMTask struct {
	Info          string
	ClientsetFunc func() (kubernetes.Interface, error)
}

// Do implements Task.
func (w *WindowsIPAMTask) Do(errCh chan error) error {
	defer close(errCh)

	clientset, err := w.ClientsetFunc()
	if err != nil {
		return err
	}
	windowsIPAM := windows.IPAM{
		Clientset: clientset,
	}
	return windowsIPAM.Enable(context.TODO())
}

// Describe implements Task.
func (w *WindowsIPAMTask) Describe() string {
	return w.Info
}

// VPCControllerTask represents a task to install the VPC controller
type VPCControllerTask struct {
	Context         context.Context
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
	oidc, err := v.ClusterProvider.NewOpenIDConnectManager(v.Context, v.ClusterConfig)
	if err != nil {
		return err
	}

	stackCollection := manager.NewStackCollection(v.ClusterProvider.AWSProvider, v.ClusterConfig)

	clientSet, err := v.ClusterProvider.NewStdClientSet(v.ClusterConfig)
	if err != nil {
		return err
	}
	irsaManager := irsa.New(v.ClusterConfig.Metadata.Name, stackCollection, oidc, clientSet)
	irsa := addons.NewIRSAHelper(oidc, stackCollection, irsaManager, v.ClusterConfig.Metadata.Name)

	// TODO PlanMode doesn't work as intended
	vpcController := addons.NewVPCController(rawClient, irsa, v.ClusterConfig.Status, v.ClusterProvider.AWSProvider.Region(), v.PlanMode)
	if err := vpcController.Deploy(v.Context); err != nil {
		return fmt.Errorf("error installing VPC controller: %w", err)
	}
	return nil
}

type devicePluginTask struct {
	kind            string
	clusterProvider *ClusterProvider
	spec            *api.ClusterConfig
	mkPlugin        addons.MkDevicePlugin
	logMessage      string
}

func (n *devicePluginTask) Describe() string { return fmt.Sprintf("install %s device plugin", n.kind) }

func (n *devicePluginTask) Do(errCh chan error) error {
	defer close(errCh)
	rawClient, err := n.clusterProvider.NewRawClient(n.spec)
	if err != nil {
		return err
	}
	devicePlugin := n.mkPlugin(rawClient, n.clusterProvider.AWSProvider.Region(), false, n.spec)
	if err := devicePlugin.Deploy(); err != nil {
		return fmt.Errorf("error installing device plugin: %w", err)
	}
	logger.Info(n.logMessage)
	return nil
}

func newNvidiaDevicePluginTask(
	clusterProvider *ClusterProvider,
	spec *api.ClusterConfig,
) tasks.Task {
	t := devicePluginTask{
		kind:            "Nvidia",
		clusterProvider: clusterProvider,
		spec:            spec,
		mkPlugin:        addons.NewNvidiaDevicePlugin,
		logMessage: `as you are using the EKS-Optimized Accelerated AMI with a GPU-enabled instance type, the Nvidia Kubernetes device plugin was automatically installed.
	to skip installing it, use --install-nvidia-plugin=false.
`,
	}
	return &t
}

func newNeuronDevicePluginTask(
	clusterProvider *ClusterProvider,
	spec *api.ClusterConfig,
) tasks.Task {
	t := devicePluginTask{
		kind:            "Neuron",
		clusterProvider: clusterProvider,
		spec:            spec,
		mkPlugin:        addons.NewNeuronDevicePlugin,
		logMessage: `as you are using the EKS-Optimized Accelerated AMI with an inf1 instance type, the AWS Neuron Kubernetes device plugin was automatically installed.
	to skip installing it, use --install-neuron-plugin=false.
`,
	}
	return &t
}

func newEFADevicePluginTask(
	clusterProvider *ClusterProvider,
	spec *api.ClusterConfig,
) tasks.Task {
	t := devicePluginTask{
		kind:            "EFA",
		clusterProvider: clusterProvider,
		spec:            spec,
		mkPlugin:        addons.NewEFADevicePlugin,
		logMessage:      "as you have enabled EFA, the EFA device plugin was automatically installed.",
	}
	return &t
}

// CreateExtraClusterConfigTasks returns all tasks for updating cluster configuration
func (c *ClusterProvider) CreateExtraClusterConfigTasks(ctx context.Context, cfg *api.ClusterConfig, preNodeGroupAddons *tasks.TaskTree, updateVPCCNITask *tasks.GenericTask) *tasks.TaskTree {
	newTasks := &tasks.TaskTree{
		Parallel:  false,
		IsSubTask: true,
	}
	if preNodeGroupAddons != nil {
		newTasks.Append(preNodeGroupAddons)
	}
	newTasks.Append(&tasks.GenericTask{
		Description: "wait for control plane to become ready",
		Doer: func() error {
			clientSet, err := c.NewRawClient(cfg)
			if err != nil {
				var unreachableErr *kubernetes.APIServerUnreachableError
				if errors.As(err, &unreachableErr) {
					logger.Warning("API server is unreachable")
				} else {
					return fmt.Errorf("error creating Clientset: %w", err)
				}
			} else if err := c.KubeProvider.WaitForControlPlane(cfg.Metadata, clientSet, c.AWSProvider.WaitTimeout()); err != nil {
				return err
			}
			return c.RefreshClusterStatus(ctx, cfg)
		},
	})
	if cfg.IsAutoModeEnabled() && cfg.VPC != nil && cfg.VPC.ID != "" {
		logger.Info("subnets supplied in subnets.private and subnets.public will be used for nodes launched by Auto Mode; please create a new NodeClass " +
			"resource if you do not want to use cluster subnets")
	}

	if api.IsEnabled(cfg.IAM.WithOIDC) {
		c.appendCreateTasksForIAMServiceAccounts(ctx, cfg, newTasks)
		if updateVPCCNITask != nil {
			newTasks.Append(updateVPCCNITask)
		}
	}

	if cfg.HasClusterCloudWatchLogging() {
		if logRetentionDays := cfg.CloudWatch.ClusterLogging.LogRetentionInDays; logRetentionDays != 0 {
			newTasks.Append(&clusterConfigTask{
				info: "update CloudWatch log retention",
				spec: cfg,
				call: func(clusterConfig *api.ClusterConfig) error {
					_, err := c.AWSProvider.CloudWatchLogs().PutRetentionPolicy(ctx, &cloudwatchlogs.PutRetentionPolicyInput{
						// The format for log group name is documented here: https://docs.aws.amazon.com/eks/latest/userguide/control-plane-logs.html
						LogGroupName:    aws.String(fmt.Sprintf("/aws/eks/%s/cluster", cfg.Metadata.Name)),
						RetentionInDays: aws.Int32(int32(logRetentionDays)),
					})
					if err != nil {
						return fmt.Errorf("error updating log retention settings: %w", err)
					}
					logger.Info("set log retention to %d days for CloudWatch logging", logRetentionDays)
					return nil
				},
			})
		}
	}

	if cfg.IsFargateEnabled() {
		manager := fargate.NewFromProvider(cfg.Metadata.Name, c.AWSProvider, c.NewStackManager(cfg))
		newTasks.Append(&fargateProfilesTask{
			info:            "create fargate profiles",
			spec:            cfg,
			clusterProvider: c,
			manager:         &manager,
			ctx:             ctx,
		})
	}

	if len(cfg.IdentityProviders) > 0 {
		newTasks.Append(identityproviders.NewAssociateProvidersTask(ctx, *cfg.Metadata, cfg.IdentityProviders, c.AWSProvider.EKS()))
	}

	if len(cfg.IAMIdentityMappings) > 0 {
		newTasks.Append(&tasks.GenericTask{
			Description: "create IAM identity mappings",
			Doer: func() error {
				clientSet, err := c.NewStdClientSet(cfg)
				if err != nil {
					return fmt.Errorf("error creating Clientset: %w", err)
				}

				rawClient, err := c.NewRawClient(cfg)
				if err != nil {
					return fmt.Errorf("error creating rawClient: %w", err)
				}
				m, err := iamidentitymapping.New(cfg, clientSet, rawClient, cfg.Metadata.Region)
				if err != nil {
					return fmt.Errorf("error initialising iamidentitymapping: %w", err)
				}

				for _, mapping := range cfg.IAMIdentityMappings {
					if err := m.Create(ctx, mapping); err != nil {
						return err
					}
				}
				return c.RefreshClusterStatus(ctx, cfg)
			},
		})
	}

	if cfg.HasWindowsNodeGroup() {
		newTasks.Append(&WindowsIPAMTask{
			Info: "enable Windows IP address management",
			ClientsetFunc: func() (kubernetes.Interface, error) {
				return c.NewStdClientSet(cfg)
			},
		})
	}

	return newTasks
}

// LogEnabledFeatures logs enabled features
func LogEnabledFeatures(clusterConfig *api.ClusterConfig) {
	if clusterConfig.HasClusterEndpointAccess() && api.EndpointsEqual(*clusterConfig.VPC.ClusterEndpoints, *api.ClusterEndpointAccessDefaults()) {
		logger.Info(clusterConfig.DefaultEndpointsMsg())
	} else {
		logger.Info(clusterConfig.CustomEndpointsMsg())
	}

	if !clusterConfig.HasClusterCloudWatchLogging() {
		logger.Info("CloudWatch logging will not be enabled for cluster %q in %q", clusterConfig.Metadata.Name, clusterConfig.Metadata.Region)
		logger.Info("you can enable it with 'eksctl utils update-cluster-logging --enable-types={SPECIFY-YOUR-LOG-TYPES-HERE (e.g. all)} --region=%s --cluster=%s'", clusterConfig.Metadata.Region, clusterConfig.Metadata.Name)
		return
	}

	all := sets.New[string](api.SupportedCloudWatchClusterLogTypes()...)

	enabled := sets.New[string]()
	if clusterConfig.HasClusterCloudWatchLogging() {
		enabled.Insert(clusterConfig.CloudWatch.ClusterLogging.EnableTypes...)
	}

	disabled := all.Difference(enabled)

	describeEnabledTypes := "no types enabled"
	if enabled.Len() > 0 {
		describeEnabledTypes = fmt.Sprintf("enabled types: %s", strings.Join(sets.List(enabled), ", "))
	}

	describeDisabledTypes := "no types disabled"
	if disabled.Len() > 0 {
		describeDisabledTypes = fmt.Sprintf("disabled types: %s", strings.Join(sets.List(disabled), ", "))
	}

	logger.Info("configuring CloudWatch logging for cluster %q in %q (%s & %s)",
		clusterConfig.Metadata.Name, clusterConfig.Metadata.Region, describeEnabledTypes, describeDisabledTypes,
	)
}

// ClusterTasksForNodeGroups returns all tasks dependent on node groups
func (c *ClusterProvider) ClusterTasksForNodeGroups(cfg *api.ClusterConfig, installNeuronDevicePluginParam, installNvidiaDevicePluginParam bool) *tasks.TaskTree {
	tasks := &tasks.TaskTree{
		Parallel:  true,
		IsSubTask: false,
	}
	var clusterRequiresNeuronDevicePlugin, clusterRequiresNvidiaDevicePlugin, efaEnabled bool
	for _, ng := range cfg.NodeGroups {
		clusterRequiresNeuronDevicePlugin = clusterRequiresNeuronDevicePlugin ||
			api.HasInstanceType(ng, instanceutils.IsNeuronInstanceType)
		// Only AL2 requires the NVIDIA device plugin
		clusterRequiresNvidiaDevicePlugin = clusterRequiresNvidiaDevicePlugin ||
			(api.HasInstanceType(ng, instanceutils.IsNvidiaInstanceType) &&
				ng.GetAMIFamily() == api.NodeImageFamilyAmazonLinux2)
		efaEnabled = efaEnabled || api.IsEnabled(ng.EFAEnabled)
	}
	for _, ng := range cfg.ManagedNodeGroups {
		clusterRequiresNeuronDevicePlugin = clusterRequiresNeuronDevicePlugin ||
			api.HasInstanceTypeManaged(ng, instanceutils.IsNeuronInstanceType)
		// Only AL2 requires the NVIDIA device plugin
		clusterRequiresNvidiaDevicePlugin = clusterRequiresNvidiaDevicePlugin ||
			(api.HasInstanceTypeManaged(ng, instanceutils.IsNvidiaInstanceType) &&
				ng.GetAMIFamily() == api.NodeImageFamilyAmazonLinux2)
		efaEnabled = efaEnabled || api.IsEnabled(ng.EFAEnabled)
	}
	if clusterRequiresNeuronDevicePlugin {
		if installNeuronDevicePluginParam {
			tasks.Append(newNeuronDevicePluginTask(c, cfg))
		} else {
			logger.Info("as you are using the EKS-Optimized Accelerated AMI with an inf1 instance type, you will need to install the AWS Neuron Kubernetes device plugin.")
			logger.Info("\t see the following page for instructions: https://awsdocs-neuron.readthedocs-hosted.com/en/latest/neuron-deploy/tutorials/tutorial-k8s.html#tutorial-k8s-env-setup-for-neuron")
		}
	}
	if clusterRequiresNvidiaDevicePlugin {
		if installNvidiaDevicePluginParam {
			tasks.Append(newNvidiaDevicePluginTask(c, cfg))
		} else {
			logger.Info("as you are using a GPU optimized instance type you will need to install NVIDIA Kubernetes device plugin.")
			logger.Info("\t see the following page for instructions: https://github.com/NVIDIA/k8s-device-plugin")
		}
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

	if efaEnabled {
		tasks.Append(newEFADevicePluginTask(c, cfg))
	}

	return tasks
}

func (c *ClusterProvider) appendCreateTasksForIAMServiceAccounts(ctx context.Context, cfg *api.ClusterConfig, tasks *tasks.TaskTree) {
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
			oidc, err := c.NewOpenIDConnectManager(ctx, cfg)
			if err != nil {
				return err
			}
			if err := oidc.CreateProvider(ctx); err != nil {
				return err
			}
			*oidcPlaceholder = *oidc
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
		cfg.IAM.ServiceAccounts,
		oidcPlaceholder,
		clientSet,
	)
	newTasks.IsSubTask = true
	tasks.Append(newTasks)
}
