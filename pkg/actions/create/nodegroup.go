package create

import (
	"fmt"
	"strings"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	defaultaddons "github.com/weaveworks/eksctl/pkg/addons/default"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils/filter"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/utils"
	"github.com/weaveworks/eksctl/pkg/vpc"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/printers"
)

// NodeGroupOptions controls specific steps of node group creation
type NodeGroupOptions struct {
	UpdateAuthConfigMap       bool
	InstallNeuronDevicePlugin bool
	InstallNvidiaDevicePlugin bool
}

// NodeGroup holds everything needed to create node groups
type NodeGroup struct {
	cfg             api.ClusterConfig
	nodegroupFilter filter.NodeGroupFilter
	clusterProvider eks.ClusterProvider
	options         NodeGroupOptions
}

// NewNodeGroup returns a NodeGroup create action
func NewNodeGroup(
	cfg api.ClusterConfig, nodegroupFilter filter.NodeGroupFilter, clusterProvider eks.ClusterProvider, options NodeGroupOptions,
) (*NodeGroup, error) {
	return &NodeGroup{
		cfg,
		nodegroupFilter,
		clusterProvider,
		options,
	}, nil
}

func (c *NodeGroup) Create() error {
	cfg := &c.cfg
	meta := cfg.Metadata

	printer := printers.NewJSONPrinter()

	ctl := &c.clusterProvider

	cmdutils.LogRegionAndVersionInfo(meta)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if ok, err := ctl.CanOperate(cfg); !ok {
		return err
	}

	if err := checkVersion(ctl, cfg.Metadata); err != nil {
		return err
	}

	if err := ctl.LoadClusterIntoSpec(cfg); err != nil {
		return errors.Wrapf(err, "getting existing configuration for cluster %q", cfg.Metadata.Name)
	}

	stackManager := ctl.NewStackManager(cfg)

	if err := c.nodegroupFilter.SetOnlyLocal(stackManager, cfg); err != nil {
		return err
	}

	logFiltered := cmdutils.ApplyFilter(cfg, &c.nodegroupFilter)

	clientSet, err := ctl.NewStdClientSet(cfg)
	if err != nil {
		return err
	}

	if err := checkARMSupport(ctl, clientSet, cfg); err != nil {
		return err
	}

	// EKS 1.14 clusters created with prior versions of eksctl may not support Managed Nodes
	supportsManagedNodes, err := ctl.SupportsManagedNodes(cfg)
	if err != nil {
		return err
	}

	if len(cfg.ManagedNodeGroups) > 0 && !supportsManagedNodes {
		return errors.New("Managed Nodegroups are not supported for this cluster version. Please update the cluster before adding managed nodegroups")
	}

	if err := eks.ValidateBottlerocketSupport(ctl.ControlPlaneVersion(), cmdutils.ToKubeNodeGroups(cfg)); err != nil {
		return err
	}

	nodeGroupService := eks.NewNodeGroupService(cfg, ctl.Provider)
	if err := nodeGroupService.Normalize(cmdutils.ToNodePools(cfg)); err != nil {
		return err
	}

	if err := printer.LogObj(logger.Debug, "cfg.json = \\\n%s\n", cfg); err != nil {
		return err
	}

	// TODO
	if err := ctl.ValidateClusterForCompatibility(cfg, stackManager); err != nil {
		return errors.Wrap(err, "cluster compatibility check failed")
	}

	if err := vpc.ValidateLegacySubnetsForNodeGroups(cfg, ctl.Provider); err != nil {
		return err
	}

	{
		logFiltered()
		logMsg := func(resource string, count int) {
			logger.Info("will create a CloudFormation stack for each of %d %s in cluster %q", count, resource, cfg.Metadata.Name)
		}
		if len(cfg.NodeGroups) > 0 {
			logMsg("nodegroups", len(cfg.NodeGroups))
		}

		if len(cfg.ManagedNodeGroups) > 0 {
			logMsg("managed nodegroups", len(cfg.ManagedNodeGroups))
		}

		tasks := &manager.TaskTree{
			Parallel: false,
		}
		if supportsManagedNodes {
			tasks.Append(stackManager.NewClusterCompatTask())
		}

		allNodeGroupTasks := &manager.TaskTree{
			Parallel: true,
		}
		awsNodeUsesIRSA, err := eks.DoesAWSNodeUseIRSA(ctl.Provider, clientSet)
		if err != nil {
			return errors.Wrap(err, "couldn't check aws-node for annotation")
		}

		if !awsNodeUsesIRSA && api.IsEnabled(cfg.IAM.WithOIDC) {
			logger.Debug("cluster has withOIDC enabled but is not using IRSA for CNI, will add CNI policy to node role")
		}

		nodeGroupTasks := stackManager.NewUnmanagedNodeGroupTask(cfg.NodeGroups, supportsManagedNodes, !awsNodeUsesIRSA)
		if nodeGroupTasks.Len() > 0 {
			allNodeGroupTasks.Append(nodeGroupTasks)
		}
		managedTasks := stackManager.NewManagedNodeGroupTask(cfg.ManagedNodeGroups, !awsNodeUsesIRSA)
		if managedTasks.Len() > 0 {
			allNodeGroupTasks.Append(managedTasks)
		}

		tasks.Append(allNodeGroupTasks)
		logger.Info(tasks.Describe())
		errs := tasks.DoAllSync()
		if len(errs) > 0 {
			logger.Info("%d error(s) occurred and nodegroups haven't been created properly, you may wish to check CloudFormation console", len(errs))
			logger.Info("to cleanup resources, run 'eksctl delete nodegroup --region=%s --cluster=%s --name=<name>' for each of the failed nodegroup", cfg.Metadata.Region, cfg.Metadata.Name)
			for _, err := range errs {
				if err != nil {
					logger.Critical("%s\n", err.Error())
				}
			}
			return fmt.Errorf("failed to create nodegroups for cluster %q", cfg.Metadata.Name)
		}
	}

	if err := PostNodeCreationTasks(ctl, cfg, clientSet, c.options); err != nil {
		return err
	}

	if err := ctl.ValidateExistingNodeGroupsForCompatibility(cfg, stackManager); err != nil {
		logger.Critical("failed checking nodegroups", err.Error())
	}

	return nil
}

func PostNodeCreationTasks(ctl *eks.ClusterProvider, cfg *api.ClusterConfig, clientSet kubernetes.Interface, options NodeGroupOptions) error {
	tasks := ctl.ClusterTasksForNodeGroups(cfg, options.InstallNeuronDevicePlugin, options.InstallNvidiaDevicePlugin)
	logger.Info(tasks.Describe())
	errs := tasks.DoAllSync()
	if len(errs) > 0 {
		logger.Info("%d error(s) occurred and nodegroups haven't been created properly, you may wish to check CloudFormation console", len(errs))
		logger.Info("to cleanup resources, run 'eksctl delete nodegroup --region=%s --cluster=%s --name=<name>' for each of the failed nodegroups", cfg.Metadata.Region, cfg.Metadata.Name)
		for _, err := range errs {
			if err != nil {
				logger.Critical("%s\n", err.Error())
			}
		}
		return fmt.Errorf("failed to create nodegroups for cluster %q", cfg.Metadata.Name)
	}

	for _, ng := range cfg.NodeGroups {
		if options.UpdateAuthConfigMap {
			// authorise nodes to join
			if err := authconfigmap.AddNodeGroup(clientSet, ng); err != nil {
				return err
			}

			// wait for nodes to join
			if err := ctl.WaitForNodes(clientSet, ng); err != nil {
				return err
			}
		}

		ShowDevicePluginMessageForNodeGroup(ng, options.InstallNeuronDevicePlugin, options.InstallNvidiaDevicePlugin)
	}
	logger.Success("created %d nodegroup(s) in cluster %q", len(cfg.NodeGroups), cfg.Metadata.Name)

	for _, ng := range cfg.ManagedNodeGroups {
		if err := ctl.WaitForNodes(clientSet, ng); err != nil {
			if cfg.PrivateCluster.Enabled {
				logger.Info("error waiting for nodes to join the cluster; this command was likely run from outside the cluster's VPC as the API server is not reachable, nodegroup(s) should still be able to join the cluster, underlying error is: %v", err)
				break
			} else {
				return err
			}
		}
	}

	logger.Success("created %d managed nodegroup(s) in cluster %q", len(cfg.ManagedNodeGroups), cfg.Metadata.Name)
	return nil
}

func ShowDevicePluginMessageForNodeGroup(nodeGroup *api.NodeGroup, installNeuronPlugin, installNvidiaPlugin bool) {
	if api.HasInstanceType(nodeGroup, utils.IsInferentiaInstanceType) {
		if installNeuronPlugin {
			logger.Info("as you are using the EKS-Optimized Accelerated AMI with an inf1 instance type, the AWS Neuron Kubernetes device plugin was automatically installed.")
			logger.Info("\t to skip installing it, use --install-neuron-plugin=false.")
		} else {
			// if neuron instance type, give instructions
			logger.Info("as you are using the EKS-Optimized Accelerated AMI with an inf1 instance type, you will need to install the AWS Neuron Kubernetes device plugin.")
			logger.Info("\t see the following page for instructions: https://github.com/aws/aws-neuron-sdk/blob/master/docs/neuron-container-tools/tutorial-k8s.md")
		}
	} else if api.HasInstanceType(nodeGroup, utils.IsGPUInstanceType) {
		if installNvidiaPlugin {
			logger.Info("as you are using the EKS-Optimized Accelerated AMI with a GPU-enabled instance type, the Nvidia Kubernetes device plugin was automatically installed.")
			logger.Info("\t to skip installing it, use --install-nvidia-plugin=false.")
		} else {
			// if GPU instance type, give instructions
			logger.Info("as you are using a GPU optimized instance type you will need to install NVIDIA Kubernetes device plugin.")
			logger.Info("\t see the following page for instructions: https://github.com/NVIDIA/k8s-device-plugin")
		}
	}
}

func checkVersion(ctl *eks.ClusterProvider, meta *api.ClusterMeta) error {
	switch meta.Version {
	case "auto":
		break
	case "":
		meta.Version = "auto"
	case "default":
		meta.Version = api.DefaultVersion
		logger.Info("will use default version (%s) for new nodegroup(s)", meta.Version)
	case "latest":
		meta.Version = api.LatestVersion
		logger.Info("will use latest version (%s) for new nodegroup(s)", meta.Version)
	default:
		if !api.IsSupportedVersion(meta.Version) {
			if api.IsDeprecatedVersion(meta.Version) {
				return fmt.Errorf("invalid version, %s is no longer supported, supported values: auto, default, latest, %s\nsee also: https://docs.aws.amazon.com/eks/latest/userguide/kubernetes-versions.html", meta.Version, strings.Join(api.SupportedVersions(), ", "))
			}
			return fmt.Errorf("invalid version %s, supported values: auto, default, latest, %s", meta.Version, strings.Join(api.SupportedVersions(), ", "))
		}
	}

	if v := ctl.ControlPlaneVersion(); v == "" {
		return fmt.Errorf("unable to get control plane version")
	} else if meta.Version == "auto" {
		meta.Version = v
		logger.Info("will use version %s for new nodegroup(s) based on control plane version", meta.Version)
	} else if meta.Version != v {
		hint := "--version=auto"
		logger.Warning("will use version %s for new nodegroup(s), while control plane version is %s; to automatically inherit the version use %q", meta.Version, v, hint)
	}

	return nil
}

func checkARMSupport(ctl *eks.ClusterProvider, clientSet kubernetes.Interface, cfg *api.ClusterConfig) error {
	rawClient, err := ctl.NewRawClient(cfg)
	if err != nil {
		return err
	}

	kubernetesVersion, err := rawClient.ServerVersion()
	if err != nil {
		return err
	}
	if api.ClusterHasInstanceType(cfg, utils.IsARMInstanceType) {
		upToDate, err := defaultaddons.DoAddonsSupportMultiArch(clientSet, rawClient, kubernetesVersion, ctl.Provider.Region())
		if err != nil {
			return err
		}
		if !upToDate {
			logger.Critical("to create an ARM nodegroup kube-proxy, coredns and aws-node addons should be up to date. " +
				"Please use `eksctl utils update-coredns`, `eksctl utils update-kube-proxy` and `eksctl utils update-aws-node` before proceeding.")
			return errors.New("expected default addons up to date")
		}
	}
	return nil
}
