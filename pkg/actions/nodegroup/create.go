package nodegroup

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/kris-nova/logger"

	"github.com/weaveworks/eksctl/pkg/actions/addon"
	defaultaddons "github.com/weaveworks/eksctl/pkg/addons/default"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils/filter"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/outposts"
	"github.com/weaveworks/eksctl/pkg/printers"
	instanceutils "github.com/weaveworks/eksctl/pkg/utils/instance"
	"github.com/weaveworks/eksctl/pkg/utils/nodes"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

// CreateOpts controls specific steps of node group creation
type CreateOpts struct {
	UpdateAuthConfigMap       *bool
	InstallNeuronDevicePlugin bool
	InstallNvidiaDevicePlugin bool
	DryRunSettings            DryRunSettings
	SkipOutdatedAddonsCheck   bool
	ConfigFileProvided        bool
	Parallelism               int
}

type DryRunSettings struct {
	DryRun    bool
	OutStream io.Writer
}

// Create creates a new nodegroup with the given options.
func (m *Manager) Create(ctx context.Context, options CreateOpts, nodegroupFilter filter.NodegroupFilter) error {
	cfg := m.cfg
	meta := cfg.Metadata
	ctl := m.ctl

	if cfg.IsControlPlaneOnOutposts() && len(cfg.ManagedNodeGroups) > 0 {
		const msg = "Managed Nodegroups are not supported on Outposts"
		if !options.ConfigFileProvided {
			return fmt.Errorf("%s; please rerun the command with --managed=false", msg)
		}
		return errors.New(msg)
	}
	if m.accessEntry.IsAWSAuthDisabled() && options.UpdateAuthConfigMap != nil {
		return errors.New("--update-auth-configmap is not supported when authenticationMode is set to API")
	}

	if cfg.IsAutoModeEnabled() && (len(cfg.ManagedNodeGroups) > 0 || len(cfg.NodeGroups) > 0) {
		addonHelper := &addon.Helper{
			Lister:      ctl.AWSProvider.EKS(),
			ClusterName: cfg.Metadata.Name,
		}
		if err := addonHelper.ValidateNodeGroupCreation(ctx); err != nil {
			return fmt.Errorf("checking if core networking addons are installed: %w", err)
		}
	}

	var (
		isOwnedCluster  = true
		skipEgressRules = false
	)

	clusterStack, err := m.stackManager.DescribeClusterStack(ctx)
	if err != nil {
		var stackNotFoundErr *manager.StackNotFoundErr
		if !errors.As(err, &stackNotFoundErr) {
			return fmt.Errorf("getting existing configuration for cluster %q: %w", meta.Name, err)
		}
		if cfg.IsControlPlaneOnOutposts() {
			return errors.New("Outposts is not supported on non eksctl-managed clusters")
		}
		logger.Warning("%s, will attempt to create nodegroup(s) on non eksctl-managed cluster", err.Error())
		if err := loadVPCFromConfig(ctx, ctl.AWSProvider, cfg); err != nil {
			return fmt.Errorf("loading VPC spec for cluster %q: %w", meta.Name, err)
		}
		isOwnedCluster = false
		if len(cfg.NodeGroups) > 0 {
			skipEgressRules, err = validateSecurityGroup(ctx, ctl.AWSProvider.EC2(), cfg.VPC.SecurityGroup)
			if err != nil {
				return err
			}
		}
	} else if err := ctl.LoadClusterIntoSpecFromStack(ctx, cfg, clusterStack); err != nil {
		return err
	}

	rawClient, err := ctl.NewRawClient(cfg)
	if err != nil {
		var unreachableErr *kubernetes.APIServerUnreachableError
		if errors.As(err, &unreachableErr) {
			const msg = "eksctl requires connectivity to the API server to create nodegroups"
			if cfg.IsControlPlaneOnOutposts() {
				return fmt.Errorf("%s; please ensure the Outpost VPC is associated with your local gateway and you are able to connect to the API server before rerunning the command: %w", msg, err)
			}
			if !m.ctl.Status.ClusterInfo.Cluster.ResourcesVpcConfig.EndpointPublicAccess {
				return fmt.Errorf("%s; please run eksctl from an environment that has access to the API server: %w", msg, err)
			}
			return fmt.Errorf("%s: %w", msg, err)
		}
		return err
	}

	if err := m.checkARMSupport(ctx, rawClient, cfg, options.SkipOutdatedAddonsCheck); err != nil {
		return err
	}

	nodePools := nodes.ToNodePools(cfg)

	nodeGroupService := eks.NewNodeGroupService(ctl.AWSProvider, m.instanceSelector, makeOutpostsService(cfg, ctl.AWSProvider))
	if err := nodeGroupService.ExpandInstanceSelectorOptions(nodePools, cfg.AvailabilityZones); err != nil {
		return err
	}

	if !options.DryRunSettings.DryRun {
		if err := nodeGroupService.Normalize(ctx, nodePools, cfg); err != nil {
			return err
		}
	}

	printer := printers.NewJSONPrinter()
	if err := printer.LogObj(logger.Debug, "cfg.json = \\\n%s\n", cfg); err != nil {
		return err
	}

	if isOwnedCluster {
		if err := ctl.ValidateClusterForCompatibility(ctx, cfg, m.stackManager); err != nil {
			return fmt.Errorf("cluster compatibility check failed: %w", err)
		}
	}

	if err := validateSubnetsAvailability(cfg); err != nil {
		return err
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

	if options.DryRunSettings.DryRun {
		clusterConfigCopy := cfg.DeepCopy()
		// Set filtered nodegroups
		clusterConfigCopy.NodeGroups = cfg.NodeGroups
		clusterConfigCopy.ManagedNodeGroups = cfg.ManagedNodeGroups
		if options.ConfigFileProvided {
			return cmdutils.PrintDryRunConfig(clusterConfigCopy, options.DryRunSettings.OutStream)
		}
		return cmdutils.PrintNodeGroupDryRunConfig(clusterConfigCopy, options.DryRunSettings.OutStream)
	}

	if err := m.nodeCreationTasks(ctx, isOwnedCluster, skipEgressRules, options.UpdateAuthConfigMap, options.Parallelism); err != nil {
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

func makeOutpostsService(clusterConfig *api.ClusterConfig, provider api.ClusterProvider) *outposts.Service {
	var outpostARN string
	if clusterConfig.IsControlPlaneOnOutposts() {
		outpostARN = clusterConfig.Outpost.ControlPlaneOutpostARN
	} else if nodeGroupOutpostARN, found := clusterConfig.FindNodeGroupOutpostARN(); found {
		outpostARN = nodeGroupOutpostARN
	} else {
		return nil
	}

	return &outposts.Service{
		OutpostsAPI: provider.Outposts(),
		EC2API:      provider.EC2(),
		OutpostID:   outpostARN,
	}
}

func (m *Manager) nodeCreationTasks(ctx context.Context, isOwnedCluster, skipEgressRules bool, updateAuthConfigMap *bool, parallelism int) error {
	cfg := m.cfg
	meta := cfg.Metadata

	taskTree := &tasks.TaskTree{
		Parallel: false,
	}

	if isOwnedCluster {
		taskTree.Append(&tasks.GenericTask{
			Doer: func() error {
				if err := m.stackManager.FixClusterCompatibility(ctx); err != nil {
					return err
				}
				hasDedicatedVPC, err := m.stackManager.ClusterHasDedicatedVPC(ctx)
				if err != nil {
					return fmt.Errorf("error checking if cluster has a dedicated VPC: %w", err)
				}
				if !hasDedicatedVPC {
					return nil
				}
				clusterExtender := &outposts.ClusterExtender{
					StackUpdater: m.stackManager,
					EC2API:       m.ctl.AWSProvider.EC2(),
					OutpostsAPI:  m.ctl.AWSProvider.Outposts(),
				}
				if err := clusterExtender.ExtendWithOutpostSubnetsIfRequired(ctx, m.cfg, m.cfg.VPC); err != nil {
					return fmt.Errorf("error extending cluster with Outpost subnets: %w", err)
				}
				return nil
			},
			Description: "fix cluster compatibility",
		})
	}

	awsNodeUsesIRSA, err := eks.DoesAWSNodeUseIRSA(ctx, m.ctl.AWSProvider, m.clientSet)
	if err != nil {
		return fmt.Errorf("couldn't check aws-node for annotation: %w", err)
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
	disableAccessEntryCreation := !m.accessEntry.IsEnabled() || updateAuthConfigMap != nil
	if nodeGroupTasks := m.stackManager.NewUnmanagedNodeGroupTask(ctx, cfg.NodeGroups, !awsNodeUsesIRSA, skipEgressRules,
		disableAccessEntryCreation, vpcImporter, parallelism); nodeGroupTasks.Len() > 0 {
		allNodeGroupTasks.Append(nodeGroupTasks)
	}
	managedTasks := m.stackManager.NewManagedNodeGroupTask(ctx, cfg.ManagedNodeGroups, !awsNodeUsesIRSA, vpcImporter, parallelism)
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

	// authorize self-managed nodes to join the cluster via aws-auth configmap
	// if EKS access entries are disabled OR
	if (!m.accessEntry.IsEnabled() && !api.IsDisabled(options.UpdateAuthConfigMap)) ||
		// if explicitly requested by the user
		api.IsEnabled(options.UpdateAuthConfigMap) {
		if err := eks.UpdateAuthConfigMap(m.cfg.NodeGroups, clientSet); err != nil {
			return err
		}
	}

	// only wait for self-managed nodes to join if either authorization method is being used
	if !api.IsDisabled(options.UpdateAuthConfigMap) {
		for _, ng := range m.cfg.NodeGroups {
			if err := eks.WaitForNodes(timeoutCtx, clientSet, ng); err != nil {
				return err
			}
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

func (m *Manager) checkARMSupport(ctx context.Context, rawClient *kubernetes.RawClient, cfg *api.ClusterConfig, skipOutdatedAddonsCheck bool) error {
	if api.ClusterHasInstanceType(cfg, instanceutils.IsARMInstanceType) {
		upToDate, err := defaultaddons.DoAddonsSupportMultiArch(ctx, rawClient.ClientSet())
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

	if err := vpc.ImportSubnetsFromSpec(ctx, provider.EC2(), cfg); err != nil {
		return err
	}
	if err := cfg.HasSufficientSubnets(); err != nil {
		logger.Critical("unable to use given %s", cfg.SubnetInfo())
		return err
	}
	return cfg.CanUseForPrivateNodeGroups()
}

func validateSecurityGroup(ctx context.Context, ec2API awsapi.EC2, securityGroupID string) (hasDefaultEgressRule bool, err error) {
	paginator := ec2.NewDescribeSecurityGroupRulesPaginator(ec2API, &ec2.DescribeSecurityGroupRulesInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("group-id"),
				Values: []string{securityGroupID},
			},
		},
	})
	var sgRules []ec2types.SecurityGroupRule
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return false, err
		}
		sgRules = append(sgRules, output.SecurityGroupRules...)
	}

	makeError := func(sgRuleID string) error {
		return fmt.Errorf("vpc.securityGroup (%s) has egress rules that were not attached by eksctl; "+
			"vpc.securityGroup should not contain any non-default external egress rules on a cluster not created by eksctl (rule ID: %s)", securityGroupID, sgRuleID)
	}

	isDefaultEgressRule := func(sgRule ec2types.SecurityGroupRule) bool {
		return aws.ToString(sgRule.IpProtocol) == "-1" && aws.ToInt32(sgRule.FromPort) == -1 && aws.ToInt32(sgRule.ToPort) == -1 && aws.ToString(sgRule.CidrIpv4) == "0.0.0.0/0"
	}

	for _, sgRule := range sgRules {
		if !aws.ToBool(sgRule.IsEgress) {
			continue
		}
		if !hasDefaultEgressRule && isDefaultEgressRule(sgRule) {
			hasDefaultEgressRule = true
			continue
		}
		if !strings.HasPrefix(aws.ToString(sgRule.Description), builder.ControlPlaneEgressRuleDescriptionPrefix) {
			return false, makeError(aws.ToString(sgRule.SecurityGroupRuleId))
		}
		matched := false
		for _, egressRule := range builder.ControlPlaneNodeGroupEgressRules {
			if aws.ToString(sgRule.IpProtocol) == egressRule.IPProtocol &&
				aws.ToInt32(sgRule.FromPort) == int32(egressRule.FromPort) &&
				aws.ToInt32(sgRule.ToPort) == int32(egressRule.ToPort) {
				matched = true
				break
			}
		}
		if !matched {
			return false, makeError(aws.ToString(sgRule.SecurityGroupRuleId))
		}
	}
	return hasDefaultEgressRule, nil
}

func validateSubnetsAvailability(spec *api.ClusterConfig) error {
	getAZs := func(subnetMapping api.AZSubnetMapping) map[string]struct{} {
		azs := make(map[string]struct{})
		for _, subnet := range subnetMapping {
			azs[subnet.AZ] = struct{}{}
		}
		return azs
	}
	privateAZs := getAZs(spec.VPC.Subnets.Private)
	publicAZs := getAZs(spec.VPC.Subnets.Public)

	validateSubnetsAvailabilityForNg := func(np api.NodePool) error {
		ng := np.BaseNodeGroup()
		subnetTypeForPrivateNetworking := map[bool]string{
			true:  "private",
			false: "public",
		}
		unavailableSubnetsErr := func(subnetLocation string) error {
			return fmt.Errorf("all %[1]s subnets from %[2]s, that the cluster was originally created on, have been deleted; to create %[1]s nodegroups within %[2]s please manually set valid %[1]s subnets via nodeGroup.SubnetIDs",
				subnetTypeForPrivateNetworking[ng.PrivateNetworking], subnetLocation)
		}

		// don't check private networking compatibility for:
		// self-managed nodegroups on local zones
		if nodeGroup, ok := np.(*api.NodeGroup); (ok && len(nodeGroup.LocalZones) > 0) ||
			// nodegroups on outposts
			(ng.OutpostARN != "" || spec.IsControlPlaneOnOutposts()) ||
			// nodegroups on user specified subnets
			len(ng.Subnets) > 0 {
			return nil
		}
		shouldCheckAcrossAllAZs := true
		for _, az := range ng.AvailabilityZones {
			shouldCheckAcrossAllAZs = false
			if _, ok := privateAZs[az]; !ok && ng.PrivateNetworking {
				return unavailableSubnetsErr(az)
			}
			if _, ok := publicAZs[az]; !ok && !ng.PrivateNetworking {
				return unavailableSubnetsErr(az)
			}
		}
		if shouldCheckAcrossAllAZs {
			if ng.PrivateNetworking && len(privateAZs) == 0 {
				return unavailableSubnetsErr(spec.VPC.ID)
			}
			if !ng.PrivateNetworking && len(publicAZs) == 0 {
				return unavailableSubnetsErr(spec.VPC.ID)
			}
		}
		return nil
	}

	for _, np := range nodes.ToNodePools(spec) {
		if err := validateSubnetsAvailabilityForNg(np); err != nil {
			return err
		}
	}

	return nil
}
