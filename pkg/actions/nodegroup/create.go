package nodegroup

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	defaultaddons "github.com/weaveworks/eksctl/pkg/addons/default"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils/filter"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/printers"
	instanceutils "github.com/weaveworks/eksctl/pkg/utils/instance"
	"github.com/weaveworks/eksctl/pkg/utils/nodes"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
	"github.com/weaveworks/eksctl/pkg/vpc"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

// CreateOpts controls specific steps of node group creation
type CreateOpts struct {
	UpdateAuthConfigMap       bool
	InstallNeuronDevicePlugin bool
	InstallNvidiaDevicePlugin bool
	DryRun                    bool
	SkipOutdatedAddonsCheck   bool
	ConfigFileProvided        bool
}

// Create creates a new nodegroup with the given options.
func (m *Manager) Create(ctx context.Context, options CreateOpts, nodegroupFilter filter.NodegroupFilter) error {
	cfg := m.cfg
	meta := cfg.Metadata
	ctl := m.ctl

	if err := checkVersion(ctl, meta); err != nil {
		return err
	}

	if err := m.checkARMSupport(ctx, ctl, m.clientSet, cfg, options.SkipOutdatedAddonsCheck); err != nil {
		return err
	}

	isOwnedCluster := true
	if err := ctl.LoadClusterIntoSpecFromStack(ctx, cfg, m.stackManager); err != nil {
		switch e := err.(type) {
		case *manager.StackNotFoundErr:
			logger.Warning("%s, will attempt to create nodegroup(s) on non eksctl-managed cluster", e.Error())
			if err := loadVPCFromConfig(ctx, ctl.AWSProvider, cfg); err != nil {
				return errors.Wrapf(err, "loading VPC spec for cluster %q", meta.Name)
			}

			isOwnedCluster = false
		default:
			return errors.Wrapf(e, "getting existing configuration for cluster %q", meta.Name)
		}
	}

	nodePools := nodes.ToNodePools(cfg)

	nodeGroupService := eks.NewNodeGroupService(ctl.AWSProvider, m.instanceSelector)
	if err := nodeGroupService.ExpandInstanceSelectorOptions(nodePools, cfg.AvailabilityZones); err != nil {
		return err
	}

	if !options.DryRun {
		if err := nodeGroupService.Normalize(ctx, nodePools, cfg.Metadata); err != nil {
			return err
		}
	}

	printer := printers.NewJSONPrinter()
	if err := printer.LogObj(logger.Debug, "cfg.json = \\\n%s\n", cfg); err != nil {
		return err
	}

	if isOwnedCluster {
		if err := ctl.ValidateClusterForCompatibility(ctx, cfg, m.stackManager); err != nil {
			return errors.Wrap(err, "cluster compatibility check failed")
		}
	}

	if err := vpc.ValidateLegacySubnetsForNodeGroups(ctx, cfg, ctl.AWSProvider); err != nil {
		return err
	}

	if err := nodegroupFilter.SetOnlyLocal(ctx, m.ctl.AWSProvider.EKS(), m.stackManager, cfg); err != nil {
		return err
	}

	logFiltered := cmdutils.ApplyFilter(cfg, nodegroupFilter)
	logFiltered()
	logMsg := func(resource string, count int) {
		logger.Info("will create a CloudFormation stack for each of %d %s in cluster %q", count, resource, meta.Name)
	}
	if len(m.cfg.NodeGroups) > 0 {
		logMsg("nodegroups", len(cfg.NodeGroups))
	}

	if len(m.cfg.ManagedNodeGroups) > 0 {
		logMsg("managed nodegroups", len(cfg.ManagedNodeGroups))
	}

	if options.DryRun {
		clusterConfigCopy := cfg.DeepCopy()
		// Set filtered nodegroups
		clusterConfigCopy.NodeGroups = cfg.NodeGroups
		clusterConfigCopy.ManagedNodeGroups = cfg.ManagedNodeGroups
		if options.ConfigFileProvided {
			return cmdutils.PrintDryRunConfig(clusterConfigCopy, os.Stdout)
		}
		return cmdutils.PrintNodeGroupDryRunConfig(clusterConfigCopy, os.Stdout)
	}

	if err := m.nodeCreationTasks(ctx, isOwnedCluster); err != nil {
		return err
	}

	if err := m.postNodeCreationTasks(ctx, m.clientSet, options); err != nil {
		return err
	}

	if err := eks.ValidateExistingNodeGroupsForCompatibility(ctx, cfg, m.stackManager); err != nil {
		logger.Critical("failed checking nodegroups", err.Error())
	}

	return nil
}

func (m *Manager) nodeCreationTasks(ctx context.Context, isOwnedCluster bool) error {
	cfg := m.cfg
	meta := cfg.Metadata

	taskTree := &tasks.TaskTree{
		Parallel: false,
	}

	if isOwnedCluster {
		taskTree.Append(m.stackManager.NewClusterCompatTask(ctx))
	}

	awsNodeUsesIRSA, err := eks.DoesAWSNodeUseIRSA(ctx, m.ctl.AWSProvider, m.clientSet)
	if err != nil {
		return errors.Wrap(err, "couldn't check aws-node for annotation")
	}

	if !awsNodeUsesIRSA && api.IsEnabled(cfg.IAM.WithOIDC) {
		logger.Debug("cluster has withOIDC enabled but is not using IRSA for CNI, will add CNI policy to node role")
	}

	var vpcImporter vpc.Importer
	if isOwnedCluster {
		vpcImporter = vpc.NewStackConfigImporter(m.stackManager.MakeClusterStackName())
	} else {
		vpcImporter = vpc.NewSpecConfigImporter(*m.ctl.Status.ClusterInfo.Cluster.ResourcesVpcConfig.ClusterSecurityGroupId, cfg.VPC)
	}

	allNodeGroupTasks := &tasks.TaskTree{
		Parallel: true,
	}
	nodeGroupTasks := m.stackManager.NewUnmanagedNodeGroupTask(ctx, cfg.NodeGroups, !awsNodeUsesIRSA, vpcImporter)
	if nodeGroupTasks.Len() > 0 {
		allNodeGroupTasks.Append(nodeGroupTasks)
	}
	managedTasks := m.stackManager.NewManagedNodeGroupTask(ctx, cfg.ManagedNodeGroups, !awsNodeUsesIRSA, vpcImporter)
	if managedTasks.Len() > 0 {
		allNodeGroupTasks.Append(managedTasks)
	}

	taskTree.Append(allNodeGroupTasks)
	return eks.DoAllNodegroupStackTasks(taskTree, meta.Region, meta.Name)
}

func (m *Manager) postNodeCreationTasks(ctx context.Context, clientSet kubernetes.Interface, options CreateOpts) error {
	tasks := m.ctl.ClusterTasksForNodeGroups(m.cfg, options.InstallNeuronDevicePlugin, options.InstallNvidiaDevicePlugin)
	logger.Info(tasks.Describe())
	errs := tasks.DoAllSync()
	if len(errs) > 0 {
		logger.Info("%d error(s) occurred and nodegroups haven't been created properly, you may wish to check CloudFormation console", len(errs))
		logger.Info("to cleanup resources, run 'eksctl delete nodegroup --region=%s --cluster=%s --name=<name>' for each of the failed nodegroups", m.cfg.Metadata.Region, m.cfg.Metadata.Name)
		for _, err := range errs {
			if err != nil {
				logger.Critical("%s\n", err.Error())
			}
		}
		return fmt.Errorf("failed to create nodegroups for cluster %q", m.cfg.Metadata.Name)
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, m.ctl.AWSProvider.WaitTimeout())
	defer cancel()

	if options.UpdateAuthConfigMap {
		if err := eks.UpdateAuthConfigMap(timeoutCtx, m.cfg.NodeGroups, clientSet); err != nil {
			return err
		}
	}
	logger.Success("created %d nodegroup(s) in cluster %q", len(m.cfg.NodeGroups), m.cfg.Metadata.Name)
	for _, ng := range m.cfg.ManagedNodeGroups {
		if err := eks.WaitForNodes(timeoutCtx, clientSet, ng); err != nil {
			if m.cfg.PrivateCluster.Enabled {
				logger.Info("error waiting for nodes to join the cluster; this command was likely run from outside the cluster's VPC as the API server is not reachable, nodegroup(s) should still be able to join the cluster, underlying error is: %v", err)
				break
			} else {
				return err
			}
		}
	}

	logger.Success("created %d managed nodegroup(s) in cluster %q", len(m.cfg.ManagedNodeGroups), m.cfg.Metadata.Name)
	return nil
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

func (m *Manager) checkARMSupport(ctx context.Context, ctl *eks.ClusterProvider, clientSet kubernetes.Interface, cfg *api.ClusterConfig, skipOutdatedAddonsCheck bool) error {
	kubeProvider := m.ctl
	rawClient, err := kubeProvider.NewRawClient(cfg)
	if err != nil {
		return err
	}

	kubernetesVersion, err := kubeProvider.ServerVersion(rawClient)
	if err != nil {
		return err
	}
	if api.ClusterHasInstanceType(cfg, instanceutils.IsARMInstanceType) {
		upToDate, err := defaultaddons.DoAddonsSupportMultiArch(ctx, ctl.AWSProvider.EKS(), rawClient, kubernetesVersion, ctl.AWSProvider.Region())
		if err != nil {
			return err
		}
		if !skipOutdatedAddonsCheck && !upToDate {
			logger.Critical("to create an ARM nodegroup kube-proxy, coredns and aws-node addons should be up to date. " +
				"Please use `eksctl utils update-coredns`, `eksctl utils update-kube-proxy` and `eksctl utils update-aws-node` before proceeding.\n" +
				"To ignore this check and proceed with the nodegroup creation, please run again with --skip-outdated-addons-check=true.")
			return errors.New("expected default addons up to date")
		}
	}
	return nil
}

func loadVPCFromConfig(ctx context.Context, provider api.ClusterProvider, cfg *api.ClusterConfig) error {
	if cfg.VPC == nil || cfg.VPC.Subnets == nil || cfg.VPC.SecurityGroup == "" || cfg.VPC.ID == "" {
		return errors.New("VPC configuration required for creating nodegroups on clusters not owned by eksctl: vpc.subnets, vpc.id, vpc.securityGroup")
	}

	if err := vpc.ImportSubnetsFromSpec(ctx, provider, cfg); err != nil {
		return err
	}

	if err := cfg.HasSufficientSubnets(); err != nil {
		logger.Critical("unable to use given %s", cfg.SubnetInfo())
		return err
	}

	return cfg.CanUseForPrivateNodeGroups()
}
