package nodegroup

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
	"github.com/weaveworks/eksctl/pkg/printers"
	"github.com/weaveworks/eksctl/pkg/utils"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
	"github.com/weaveworks/eksctl/pkg/vpc"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

// Options controls specific steps of node group creation
type CreateOpts struct {
	UpdateAuthConfigMap       bool
	InstallNeuronDevicePlugin bool
	InstallNvidiaDevicePlugin bool
}

func (m *Manager) Create(options CreateOpts, nodegroupFilter filter.NodeGroupFilter) error {
	cfg := m.cfg
	meta := cfg.Metadata
	ctl := m.ctl

	if err := checkVersion(ctl, meta); err != nil {
		return err
	}

	if err := checkARMSupport(ctl, m.clientSet, cfg); err != nil {
		return err
	}

	var isOwnedCluster = true
	if err := ctl.LoadClusterIntoSpecFromStack(cfg, m.stackManager); err != nil {
		switch e := err.(type) {
		case *manager.StackNotFoundErr:
			logger.Warning("%s, will attempt to create nodegroup(s) on non eksctl-managed cluster", e.Error())
			if err := loadVPCFromConfig(ctl.Provider, cfg); err != nil {
				return errors.Wrapf(err, "loading VPC spec for cluster %q", meta.Name)
			}

			isOwnedCluster = false
		default:
			return errors.Wrapf(e, "getting existing configuration for cluster %q", meta.Name)
		}
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

	printer := printers.NewJSONPrinter()
	if err := printer.LogObj(logger.Debug, "cfg.json = \\\n%s\n", cfg); err != nil {
		return err
	}

	if isOwnedCluster {
		if err := ctl.ValidateClusterForCompatibility(cfg, m.stackManager); err != nil {
			return errors.Wrap(err, "cluster compatibility check failed")
		}
	}

	if err := vpc.ValidateLegacySubnetsForNodeGroups(cfg, ctl.Provider); err != nil {
		return err
	}

	{
		if err := nodegroupFilter.SetOnlyLocal(ctl.Provider.EKS(), m.stackManager, cfg); err != nil {
			return err
		}

		logFiltered := cmdutils.ApplyFilter(cfg, &nodegroupFilter)
		logFiltered()
		logMsg := func(resource string, count int) {
			logger.Info("will create a CloudFormation stack for each of %d %s in cluster %q", count, resource, meta.Name)
		}
		if len(cfg.NodeGroups) > 0 {
			logMsg("nodegroups", len(cfg.NodeGroups))
		}

		if len(cfg.ManagedNodeGroups) > 0 {
			logMsg("managed nodegroups", len(cfg.ManagedNodeGroups))
		}

		taskTree := &tasks.TaskTree{
			Parallel: false,
		}

		if supportsManagedNodes && isOwnedCluster {
			taskTree.Append(m.stackManager.NewClusterCompatTask())
		}

		awsNodeUsesIRSA, err := eks.DoesAWSNodeUseIRSA(ctl.Provider, m.clientSet)
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
			vpcImporter = vpc.NewSpecConfigImporter(*ctl.Status.ClusterInfo.Cluster.ResourcesVpcConfig.ClusterSecurityGroupId, cfg.VPC)
		}

		allNodeGroupTasks := &tasks.TaskTree{
			Parallel: true,
		}
		nodeGroupTasks := m.stackManager.NewUnmanagedNodeGroupTask(cfg.NodeGroups, supportsManagedNodes, !awsNodeUsesIRSA, vpcImporter)
		if nodeGroupTasks.Len() > 0 {
			allNodeGroupTasks.Append(nodeGroupTasks)
		}
		managedTasks := m.stackManager.NewManagedNodeGroupTask(cfg.ManagedNodeGroups, !awsNodeUsesIRSA, vpcImporter)
		if managedTasks.Len() > 0 {
			allNodeGroupTasks.Append(managedTasks)
		}

		taskTree.Append(allNodeGroupTasks)
		logger.Info(taskTree.Describe())
		errs := taskTree.DoAllSync()
		if len(errs) > 0 {
			logger.Info("%d error(s) occurred and nodegroups haven't been created properly, you may wish to check CloudFormation console", len(errs))
			logger.Info("to cleanup resources, run 'eksctl delete nodegroup --region=%s --cluster=%s --name=<name>' for each of the failed nodegroup", meta.Region, meta.Name)
			for _, err := range errs {
				if err != nil {
					logger.Critical("%s\n", err.Error())
				}
			}
			return fmt.Errorf("failed to create nodegroups for cluster %q", meta.Name)
		}
	}

	if err := m.postNodeCreationTasks(m.clientSet, options); err != nil {
		return err
	}

	if err := ctl.ValidateExistingNodeGroupsForCompatibility(cfg, m.stackManager); err != nil {
		logger.Critical("failed checking nodegroups", err.Error())
	}

	return nil
}

func (m *Manager) postNodeCreationTasks(clientSet kubernetes.Interface, options CreateOpts) error {
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

	if options.UpdateAuthConfigMap {
		for _, ng := range m.cfg.NodeGroups {
			// authorise nodes to join
			if err := authconfigmap.AddNodeGroup(clientSet, ng); err != nil {
				return err
			}

			// wait for nodes to join
			if err := m.ctl.WaitForNodes(clientSet, ng); err != nil {
				return err
			}
		}
	}
	logger.Success("created %d nodegroup(s) in cluster %q", len(m.cfg.NodeGroups), m.cfg.Metadata.Name)

	for _, ng := range m.cfg.ManagedNodeGroups {
		if err := m.ctl.WaitForNodes(clientSet, ng); err != nil {
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

func loadVPCFromConfig(provider api.ClusterProvider, cfg *api.ClusterConfig) error {
	if cfg.VPC.Subnets == nil || cfg.VPC.SecurityGroup == "" || cfg.VPC.ID == "" {
		return errors.New("VPC configuration required for creating nodegroups on clusters not owned by eksctl: vpc.subnets, vpc.id, vpc.securityGroup")
	}

	if err := vpc.ImportSubnetsFromSpec(provider, cfg); err != nil {
		return err
	}

	if err := cfg.HasSufficientSubnets(); err != nil {
		logger.Critical("unable to use given %s", cfg.SubnetInfo())
		return err
	}

	if err := cfg.CanUseForPrivateNodeGroups(); err != nil {
		return err
	}

	return nil
}
