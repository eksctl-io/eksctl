package create

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/aws/amazon-ec2-instance-selector/v2/pkg/selector"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/runtime"
	kubeclient "k8s.io/client-go/kubernetes"
	clientcmdlatest "k8s.io/client-go/tools/clientcmd/api/latest"

	"github.com/weaveworks/eksctl/pkg/accessentry"
	accessentryactions "github.com/weaveworks/eksctl/pkg/actions/accessentry"
	"github.com/weaveworks/eksctl/pkg/actions/addon"
	"github.com/weaveworks/eksctl/pkg/actions/flux"
	"github.com/weaveworks/eksctl/pkg/actions/karpenter"
	"github.com/weaveworks/eksctl/pkg/actions/podidentityassociation"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils/filter"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/kops"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/outposts"
	"github.com/weaveworks/eksctl/pkg/printers"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
	"github.com/weaveworks/eksctl/pkg/utils/names"
	"github.com/weaveworks/eksctl/pkg/utils/nodes"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

const (
	vpcControllerInfoMessage = "you no longer need to install the VPC resource controller on Linux worker nodes to run " +
		"Windows workloads in EKS clusters created after Oct 22, 2021. You can enable Windows IP address management on the EKS control plane via " +
		"a ConﬁgMap setting (see https://docs.aws.amazon.com/eks/latest/userguide/windows-support.html for details). eksctl will automatically patch the ConfigMap to enable " +
		"Windows IP address management when a Windows nodegroup is created. For existing clusters, you can enable it manually " +
		"and run `eksctl utils install-vpc-controllers` with the --delete ﬂag to remove the worker node installation of the VPC resource controller"
)

var (
	once                     sync.Once
	createKarpenterInstaller = karpenter.NewInstaller
)

func createClusterCmd(cmd *cmdutils.Cmd) {
	createClusterCmdWithRunFunc(cmd, func(cmd *cmdutils.Cmd, ngFilter *filter.NodeGroupFilter, params *cmdutils.CreateClusterCmdParams) error {
		ctl, err := cmd.NewCtl()
		if err != nil {
			return err
		}
		return doCreateCluster(cmd, ngFilter, params, ctl, func(clusterName string, stackCreator accessentryactions.StackCreator) accessentryactions.CreatorInterface {
			return &accessentryactions.Creator{
				ClusterName:  clusterName,
				StackCreator: stackCreator,
			}
		})
	})
}

func checkClusterVersion(cfg *api.ClusterConfig) error {
	switch cfg.Metadata.Version {
	case "auto":
		cfg.Metadata.Version = api.DefaultVersion
	case "latest":
		cfg.Metadata.Version = api.LatestVersion
	}

	if err := api.ValidateClusterVersion(cfg); err != nil {
		return err
	}
	if cfg.Metadata.Version == "" {
		cfg.Metadata.Version = api.DefaultVersion
	}
	return nil
}

func createClusterCmdWithRunFunc(cmd *cmdutils.Cmd, runFunc func(cmd *cmdutils.Cmd, ngFilter *filter.NodeGroupFilter, params *cmdutils.CreateClusterCmdParams) error) {
	cfg := api.NewClusterConfig()
	ng := api.NewNodeGroup()
	cmd.ClusterConfig = cfg

	params := &cmdutils.CreateClusterCmdParams{}

	cmd.SetDescription("cluster", "Create a cluster", "")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		ngFilter := filter.NewNodeGroupFilter()
		if err := cmdutils.NewCreateClusterLoader(cmd, ngFilter, ng, params).Load(); err != nil {
			return err
		}
		err := checkClusterVersion(cmd.ClusterConfig)
		if err != nil {
			return err
		}
		return runFunc(cmd, ngFilter, params)
	}

	exampleClusterName := names.ForCluster("", "")
	exampleNodeGroupName := names.ForNodeGroup("", "")

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVarP(&cfg.Metadata.Name, "name", "n", "", fmt.Sprintf("EKS cluster name (generated if unspecified, e.g. %q)", exampleClusterName))
		cmdutils.AddStringToStringVarPFlag(fs, &cfg.Metadata.Tags, "tags", "", map[string]string{}, "Used to tag the AWS resources")
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		fs.BoolVar(cfg.IAM.WithOIDC, "with-oidc", false, "Enable the IAM OIDC provider")
		fs.StringSliceVar(&params.AvailabilityZones, "zones", nil, "(auto-select if unspecified)")
		cmdutils.AddVersionFlag(fs, cfg.Metadata, "")
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
		fs.BoolVarP(&params.InstallWindowsVPCController, "install-vpc-controllers", "", false, "Install VPC controller that's required for Windows workloads")
		fs.BoolVarP(&params.Fargate, "fargate", "", false, "Create a Fargate profile scheduling pods in the default and kube-system namespaces onto Fargate")
		fs.BoolVarP(&params.DryRun, "dry-run", "", false, "Dry-run mode that skips cluster creation and outputs a ClusterConfig")

		_ = fs.MarkDeprecated("install-vpc-controllers", vpcControllerInfoMessage)
	})

	cmd.FlagSetGroup.InFlagSet("Initial nodegroup", func(fs *pflag.FlagSet) {
		fs.StringVar(&ng.Name, "nodegroup-name", "", fmt.Sprintf("name of the nodegroup (generated if unspecified, e.g. %q)", exampleNodeGroupName))
		fs.BoolVar(&params.WithoutNodeGroup, "without-nodegroup", false, "if set, initial nodegroup will not be created")
		cmdutils.AddCommonCreateNodeGroupFlags(fs, cmd, ng, &params.CreateManagedNGOptions)
	})

	cmd.FlagSetGroup.InFlagSet("Cluster and nodegroup add-ons", func(fs *pflag.FlagSet) {
		cmdutils.AddCommonCreateNodeGroupAddonsFlags(fs, ng, &params.CreateNGOptions)
	})

	cmd.FlagSetGroup.InFlagSet("VPC networking", func(fs *pflag.FlagSet) {
		fs.IPNetVar(&cfg.VPC.CIDR.IPNet, "vpc-cidr", cfg.VPC.CIDR.IPNet, "global CIDR to use for VPC")
		params.Subnets = map[api.SubnetTopology]*[]string{
			api.SubnetTopologyPrivate: fs.StringSlice("vpc-private-subnets", nil, "re-use private subnets of an existing VPC; the subnets must exist in availability zones and not other types of zones"),
			api.SubnetTopologyPublic:  fs.StringSlice("vpc-public-subnets", nil, "re-use public subnets of an existing VPC; the subnets must exist in availability zones and not other types of zones"),
		}
		fs.StringVar(&params.KopsClusterNameForVPC, "vpc-from-kops-cluster", "", "re-use VPC from a given kops cluster")
		fs.StringVar(cfg.VPC.NAT.Gateway, "vpc-nat-mode", api.ClusterSingleNAT, "VPC NAT mode, valid options: HighlyAvailable, Single, Disable")
	})

	cmdutils.AddInstanceSelectorOptions(cmd.FlagSetGroup, ng)

	cmdutils.AddCommonFlagsForAWS(cmd, &cmd.ProviderConfig, true)

	cmd.FlagSetGroup.InFlagSet("Output kubeconfig", func(fs *pflag.FlagSet) {
		cmdutils.AddCommonFlagsForKubeconfig(fs, &params.KubeconfigPath, &params.AuthenticatorRoleARN, &params.SetContext, &params.AutoKubeconfigPath, exampleClusterName)
		fs.BoolVar(&params.WriteKubeconfig, "write-kubeconfig", true, "toggle writing of kubeconfig")
	})
}

func doCreateCluster(cmd *cmdutils.Cmd, ngFilter *filter.NodeGroupFilter, params *cmdutils.CreateClusterCmdParams, ctl *eks.ClusterProvider,
	makeAccessEntryCreator func(clusterName string, creator accessentryactions.StackCreator) accessentryactions.CreatorInterface) error {
	var err error
	cfg := cmd.ClusterConfig
	meta := cmd.ClusterConfig.Metadata

	if meta.Name != "" && api.IsInvalidNameArg(meta.Name) {
		return api.ErrInvalidName(meta.Name)
	}
	printer := printers.NewJSONPrinter()

	if params.DryRun {
		originalWriter := logger.Writer
		logger.Writer = io.Discard
		defer func() {
			logger.Writer = originalWriter
		}()
	}

	//prevent logging multiple times
	once.Do(func() {
		cmdutils.LogRegionAndVersionInfo(meta)
	})

	if err := cfg.ValidatePrivateCluster(); err != nil {
		return err
	}

	if err := cfg.ValidateClusterEndpointConfig(); err != nil {
		return err
	}

	// If it's a private-only cluster warn the user.
	if api.PrivateOnly(cfg.VPC.ClusterEndpoints) && !cfg.IsControlPlaneOnOutposts() {
		logger.Warning(api.ErrClusterEndpointPrivateOnly.Error())
	}

	// if using a custom shared node security group, warn that the rules are managed by default
	if cfg.VPC.SharedNodeSecurityGroup != "" && api.IsEnabled(cfg.VPC.ManageSharedNodeSecurityGroupRules) {
		logger.Warning("security group rules may be added by eksctl; see vpc.manageSharedNodeSecurityGroupRules to disable this behavior")
	}

	if params.AutoKubeconfigPath {
		if params.KubeconfigPath != kubeconfig.DefaultPath() {
			return fmt.Errorf("--kubeconfig and --auto-kubeconfig %s", cmdutils.IncompatibleFlags)
		}
		params.KubeconfigPath = kubeconfig.AutoPath(meta.Name)
	}

	ctx := context.Background()

	if checkSubnetsGivenAsFlags(params) {
		// undo defaulting and reset it, as it's not set via config file;
		// default value here causes errors as vpc.ImportVPC doesn't
		// treat remote state as authority over local state
		cfg.VPC.CIDR = nil
		if cfg.VPC.Subnets == nil {
			cfg.VPC.Subnets = &api.ClusterSubnets{
				Private: api.NewAZSubnetMapping(),
				Public:  api.NewAZSubnetMapping(),
			}
		}

		// load subnets from local map created from flags, into the config
		importSubnets := func(subnetMapping api.AZSubnetMapping, subnetIDs *[]string) error {
			if subnetIDs != nil {
				if err := vpc.ImportSubnetsFromIDList(ctx, ctl.AWSProvider.EC2(), cfg, subnetMapping, *subnetIDs); err != nil {
					return err
				}
			}
			return nil
		}
		if err := importSubnets(cfg.VPC.Subnets.Public, params.Subnets[api.SubnetTopologyPublic]); err != nil {
			return err
		}
		if err := importSubnets(cfg.VPC.Subnets.Private, params.Subnets[api.SubnetTopologyPrivate]); err != nil {
			return err
		}
		if params.DryRun {
			cfg.AvailabilityZones = nil
		}
	}
	logFiltered := cmdutils.ApplyFilter(cfg, ngFilter)
	kubeNodeGroups := cmdutils.ToKubeNodeGroups(cfg.NodeGroups, cfg.ManagedNodeGroups)

	// Check if flux binary exists early in the process, so it doesn't fail at the end when the cluster
	// has already been created with a missing flux binary error which should have been caught earlier.
	// Note: we aren't running PreFlight here, we just check for the binary.
	if cfg.HasGitOpsFluxConfigured() {
		if _, err := exec.LookPath("flux"); err != nil {
			return fmt.Errorf("flux binary is required when gitops configuration is set: %w", err)
		}
	}

	if params.InstallWindowsVPCController {
		if !eks.SupportsWindowsWorkloads(kubeNodeGroups) {
			return errors.New("running Windows workloads requires having both Windows and Linux (AmazonLinux2) node groups")
		}
		logger.Warning(vpcControllerInfoMessage)
	} else {
		eks.LogWindowsCompatibility(kubeNodeGroups, cfg.Metadata)
	}

	if err := accessentry.ValidateAPIServerAccess(cfg); err != nil {
		return err
	}

	var outpostsService *outposts.Service

	if cfg.IsControlPlaneOnOutposts() {
		outpostsService = &outposts.Service{
			OutpostsAPI: ctl.AWSProvider.Outposts(),
			EC2API:      ctl.AWSProvider.EC2(),
			OutpostID:   cfg.Outpost.ControlPlaneOutpostARN,
		}
		outpost, err := outpostsService.GetOutpost(ctx)
		if err != nil {
			return fmt.Errorf("error getting Outpost details: %w", err)
		}
		if !params.DryRun {
			cfg.AvailabilityZones = []string{aws.ToString(outpost.AvailabilityZone)}
		}
		if cfg.Outpost.HasPlacementGroup() {
			if err := outpostsService.ValidatePlacementGroup(ctx, cfg.Outpost.ControlPlanePlacement); err != nil {
				return err
			}
		}

		if err := outpostsService.SetOrValidateOutpostInstanceType(ctx, cfg.Outpost); err != nil {
			return fmt.Errorf("error setting or validating instance type for the control plane: %w", err)
		}

		if !cfg.HasAnySubnets() && len(kubeNodeGroups) > 0 {
			return errors.New("cannot create nodegroups on Outposts when the VPC is created by eksctl as it will not have connectivity to the API server; please rerun the command with `--without-nodegroup` and run `eksctl create nodegroup` after associating the VPC with a local gateway and ensuring connectivity to the API server")
		}
	} else if _, hasNodeGroupsOnOutposts := cfg.FindNodeGroupOutpostARN(); hasNodeGroupsOnOutposts {
		return errors.New("creating nodegroups on Outposts when the control plane is not on Outposts is not supported during cluster creation; " +
			"either create the nodegroups after cluster creation or consider creating the control plane on Outposts")
	}

	if err := createOrImportVPC(ctx, cmd, cfg, params, ctl); err != nil {
		return err
	}

	instanceSelector, err := selector.New(ctx, ctl.AWSProvider.AWSConfig())
	if err != nil {
		return err
	}
	nodeGroupService := eks.NewNodeGroupService(ctl.AWSProvider, instanceSelector, outpostsService)
	nodePools := nodes.ToNodePools(cfg)
	if err := nodeGroupService.ExpandInstanceSelectorOptions(nodePools, cfg.AvailabilityZones); err != nil {
		return err
	}

	if params.DryRun {
		return cmdutils.PrintDryRunConfig(cfg, cmd.CobraCommand.OutOrStdout())
	}

	if err := nodeGroupService.Normalize(ctx, nodePools, cfg); err != nil {
		return err
	}

	logger.Info("using Kubernetes version %s", meta.Version)
	logger.Info("creating %s", cfg.LogString())

	// TODO dry-run mode should provide a way to render config with all defaults set
	// we should also make a call to resolve the AMI and write the result, similarly
	// the body of the SSH key can be read

	if err := printer.LogObj(logger.Debug, "cfg.json = \\\n%s\n", cfg); err != nil {
		return err
	}

	stackManager := ctl.NewStackManager(cfg)
	if cmd.ClusterConfigFile == "" {
		logMsg := func(resource string) {
			logger.Info("will create 2 separate CloudFormation stacks for cluster itself and the initial %s", resource)
		}
		if len(cfg.NodeGroups) == 1 {
			logMsg("nodegroup")
		} else if len(cfg.ManagedNodeGroups) == 1 {
			logMsg("managed nodegroup")
		}
	} else {
		logMsg := func(resource string, count int) {
			logger.Info("will create a CloudFormation stack for cluster itself and %d %s stack(s)", count, resource)
		}
		logFiltered()

		logMsg("nodegroup", len(cfg.NodeGroups))
		logMsg("managed nodegroup", len(cfg.ManagedNodeGroups))
	}

	logger.Info("if you encounter any issues, check CloudFormation console or try 'eksctl utils describe-stacks --region=%s --cluster=%s'", meta.Region, meta.Name)

	eks.LogEnabledFeatures(cfg)
	postClusterCreationTasks := ctl.CreateExtraClusterConfigTasks(ctx, cfg)

	var preNodegroupAddons, postNodegroupAddons *tasks.TaskTree
	if len(cfg.Addons) > 0 {
		preNodegroupAddons, postNodegroupAddons = addon.CreateAddonTasks(ctx, cfg, ctl, true, cmd.ProviderConfig.WaitTimeout)
		postClusterCreationTasks.Append(preNodegroupAddons)
	}

	taskTree := stackManager.NewTasksToCreateCluster(ctx, cfg.NodeGroups, cfg.ManagedNodeGroups, cfg.AccessConfig, makeAccessEntryCreator(cfg.Metadata.Name, stackManager), postClusterCreationTasks)

	logger.Info(taskTree.Describe())
	if errs := taskTree.DoAllSync(); len(errs) > 0 {
		logger.Warning("%d error(s) occurred and cluster hasn't been created properly, you may wish to check CloudFormation console", len(errs))
		logger.Info("to cleanup resources, run 'eksctl delete cluster --region=%s --name=%s'", meta.Region, meta.Name)
		for _, err := range errs {
			ufe := &api.UnsupportedFeatureError{}
			if errors.As(err, &ufe) {
				logger.Critical(ufe.Message)
			}
			logger.Critical("%s\n", err.Error())
		}
		return fmt.Errorf("failed to create cluster %q", meta.Name)
	}

	logger.Info("waiting for the control plane to become ready")

	// obtain cluster credentials, write kubeconfig

	{ // post-creation action
		var kubeconfigContextName string

		if params.WriteKubeconfig {
			kubectlConfig := kubeconfig.NewForKubectl(cfg, eks.GetUsername(ctl.Status.IAMRoleARN), params.AuthenticatorRoleARN, ctl.AWSProvider.Profile().Name)
			kubeconfigContextName = kubectlConfig.CurrentContext

			params.KubeconfigPath, err = kubeconfig.Write(params.KubeconfigPath, *kubectlConfig, params.SetContext)
			if err != nil {
				logger.Warning("unable to write kubeconfig %s, please retry with 'eksctl utils write-kubeconfig -n %s': %v", params.KubeconfigPath, meta.Name, err)
			} else {
				logger.Success("saved kubeconfig as %q", params.KubeconfigPath)
			}
		} else {
			params.KubeconfigPath = ""
		}

		ngTasks := ctl.ClusterTasksForNodeGroups(cfg, params.InstallNeuronDevicePlugin, params.InstallNvidiaDevicePlugin)

		logger.Info(ngTasks.Describe())
		if errs := ngTasks.DoAllSync(); len(errs) > 0 {
			logger.Warning("%d error(s) occurred and post actions have failed, you may wish to check CloudFormation console", len(errs))
			logger.Info("to cleanup resources, run 'eksctl delete cluster --region=%s --name=%s'", meta.Region, meta.Name)
			for _, err := range errs {
				logger.Critical("%s\n", err.Error())
			}
			return fmt.Errorf("failed to create cluster %q", meta.Name)
		}
		logger.Success("all EKS cluster resources for %q have been created", meta.Name)

		makeClientSet := clientSetCreator(ctl, cfg)
		{
			clientSet, err := makeClientSet()
			if err != nil {
				if !api.IsDisabled(cfg.AccessConfig.BootstrapClusterCreatorAdminPermissions) {
					return err
				}
				if !accessentry.IsEnabled(cfg.AccessConfig.AuthenticationMode) {
					return err
				}
				if len(cfg.NodeGroups) > 0 {
					logger.Warning("not waiting for self-managed nodes to become ready as API server is not accessible; run `kubectl get nodes` to ensure the nodes are ready: %v", err)
				}
			} else {
				ngCtx, cancel := context.WithTimeout(ctx, cmd.ProviderConfig.WaitTimeout)
				defer cancel()

				// authorize self-managed nodes to join the cluster via aws-auth configmap
				// only if EKS access entries are disabled
				if cfg.AccessConfig.AuthenticationMode == ekstypes.AuthenticationModeConfigMap {
					if err := eks.UpdateAuthConfigMap(cfg.NodeGroups, clientSet); err != nil {
						return err
					}
				}

				for _, ng := range cfg.NodeGroups {
					if err := eks.WaitForNodes(ngCtx, clientSet, ng); err != nil {
						return err
					}
				}
				logger.Success("created %d nodegroup(s) in cluster %q", len(cfg.NodeGroups), cfg.Metadata.Name)

				for _, ng := range cfg.ManagedNodeGroups {
					if err := eks.WaitForNodes(ngCtx, clientSet, ng); err != nil {
						return err
					}
				}
				logger.Success("created %d managed nodegroup(s) in cluster %q", len(cfg.ManagedNodeGroups), cfg.Metadata.Name)
			}
		}
		if postNodegroupAddons != nil && postNodegroupAddons.Len() > 0 {
			if errs := postNodegroupAddons.DoAllSync(); len(errs) > 0 {
				logger.Warning("%d error(s) occurred while creating addons", len(errs))
				for _, err := range errs {
					logger.Critical("%s\n", err.Error())
				}
				return errors.New("failed to create addons")
			}
		}

		if len(cfg.IAM.PodIdentityAssociations) > 0 {
			if err := podidentityassociation.NewCreator(
				cfg.Metadata.Name,
				stackManager,
				ctl.AWSProvider.EKS(),
			).CreatePodIdentityAssociations(ctx, cfg.IAM.PodIdentityAssociations); err != nil {
				return err
			}
		}

		// After we have the cluster config and all the nodes are done, we install Karpenter if necessary.
		if cfg.Karpenter != nil {
			config := kubeconfig.NewForKubectl(cfg, eks.GetUsername(ctl.Status.IAMRoleARN), params.AuthenticatorRoleARN, ctl.AWSProvider.Profile().Name)
			kubeConfigBytes, err := runtime.Encode(clientcmdlatest.Codec, config)
			if err != nil {
				return errors.Wrap(err, "generating kubeconfig")
			}
			clientSet, err := makeClientSet()
			if err != nil {
				return fmt.Errorf("error installing Karpenter: %w", err)
			}
			if err := installKarpenter(ctx, ctl, cfg, stackManager, clientSet, kubernetes.NewRESTClientGetter("karpenter", string(kubeConfigBytes))); err != nil {
				return err
			}
		}

		if cfg.HasGitOpsFluxConfigured() {
			clientSet, err := makeClientSet()
			if err != nil {
				return fmt.Errorf("error installing Flux: %w", err)
			}
			installer, err := flux.New(clientSet, cfg.GitOps)
			logger.Info("gitops configuration detected, setting installer to Flux v2")
			if err != nil {
				return errors.Wrapf(err, "could not initialise Flux installer")
			}

			if err := installer.Run(); err != nil {
				return err
			}

			//TODO why was it returning early before? I want to remove this line :thinking:
			return nil
		}

		env, err := ctl.GetCredentialsEnv(ctx)
		if err != nil {
			return err
		}
		if err := kubeconfig.CheckAllCommands(params.KubeconfigPath, params.SetContext, kubeconfigContextName, env); err != nil {
			logger.Critical("%s\n", err.Error())
			logger.Info("cluster should be functional despite missing (or misconfigured) client binaries")
		}

		if cfg.IsFullyPrivate() && !cfg.IsControlPlaneOnOutposts() {
			// disable public access
			logger.Info("disabling public endpoint access for the cluster")
			cfg.VPC.ClusterEndpoints.PublicAccess = api.Disabled()
			if err := ctl.UpdateClusterConfigForEndpoints(ctx, cfg); err != nil {
				return errors.Wrap(err, "error disabling public endpoint access for the cluster")
			}
			logger.Info("fully private cluster %q has been created. For subsequent operations, eksctl must be run from within the cluster's VPC, a peered VPC or some other means like AWS Direct Connect", cfg.Metadata.Name)
		}
	}

	logger.Success("%s is ready", meta.LogString())

	return printer.LogObj(logger.Debug, "cfg.json = \\\n%s\n", cfg)
}

// installKarpenter prepares the environment for Karpenter, by creating the following resources:
// - iam roles and profiles
// - service account
// - identity mapping
// then proceeds with installing Karpenter using Helm.
func installKarpenter(ctx context.Context, ctl *eks.ClusterProvider, cfg *api.ClusterConfig, stackManager manager.StackManager, clientSet kubeclient.Interface, restClientGetter *kubernetes.SimpleRESTClientGetter) error {
	installer, err := createKarpenterInstaller(ctx, cfg, ctl, stackManager, clientSet, restClientGetter)
	if err != nil {
		return fmt.Errorf("failed to create installer: %w", err)
	}
	if err := installer.Create(ctx); err != nil {
		return fmt.Errorf("failed to install Karpenter: %w", err)
	}

	return nil
}

func createOrImportVPC(ctx context.Context, cmd *cmdutils.Cmd, cfg *api.ClusterConfig, params *cmdutils.CreateClusterCmdParams, ctl *eks.ClusterProvider) error {
	customNetworkingNotice := "custom VPC/subnets will be used; if resulting cluster doesn't function as expected, make sure to review the configuration of VPC/subnets"

	subnetsGiven := cfg.HasAnySubnets() // this will be false when neither flags nor config has any subnets
	if !subnetsGiven && params.KopsClusterNameForVPC == "" {
		if !cfg.IsControlPlaneOnOutposts() {
			userProvidedAZs, err := eks.SetAvailabilityZones(ctx, cfg, params.AvailabilityZones, ctl.AWSProvider.EC2(), ctl.AWSProvider.Region())
			if err != nil {
				return err
			}

			// If the availability zones were provided at random, we already did this check.
			if userProvidedAZs {
				if err := eks.CheckInstanceAvailability(ctx, cfg, ctl.AWSProvider.EC2()); err != nil {
					return err
				}
			}

			if len(cfg.LocalZones) > 0 {
				if err := eks.ValidateLocalZones(ctx, ctl.AWSProvider.EC2(), cfg.LocalZones, ctl.AWSProvider.Region()); err != nil {
					return err
				}
			}

			// Skip setting subnets
			// The default subnet config set by SetSubnets will fail validation on a subsequent run of `create cluster`
			// because those fields indicate usage of pre-existing VPC and subnets
			// default: create dedicated VPC
			if params.DryRun {
				return nil
			}
		}
		return vpc.SetSubnets(cfg.VPC, cfg.AvailabilityZones, cfg.LocalZones)
	}

	if params.KopsClusterNameForVPC != "" {
		if cfg.IsControlPlaneOnOutposts() {
			return errors.New("cannot specify --vpc-from-kops-cluster when creating a cluster on Outposts")
		}
		// import VPC from a given kops cluster
		if len(params.AvailabilityZones) != 0 {
			return fmt.Errorf("--vpc-from-kops-cluster and --zones %s", cmdutils.IncompatibleFlags)
		}
		if cmd.CobraCommand.Flag("vpc-cidr").Changed {
			return fmt.Errorf("--vpc-from-kops-cluster and --vpc-cidr %s", cmdutils.IncompatibleFlags)
		}

		if subnetsGiven {
			return fmt.Errorf("--vpc-from-kops-cluster and --vpc-private-subnets/--vpc-public-subnets %s", cmdutils.IncompatibleFlags)
		}

		kw, err := kops.NewWrapper(cmd.ProviderConfig.Region, params.KopsClusterNameForVPC)
		if err != nil {
			return err
		}

		if params.DryRun {
			return nil
		}

		if err := kw.UseVPC(ctx, ctl.AWSProvider.EC2(), cfg); err != nil {
			return err
		}

		if err := cfg.CanUseForPrivateNodeGroups(); err != nil {
			return err
		}

		logger.Success("using %s from kops cluster %q", cfg.SubnetInfo(), params.KopsClusterNameForVPC)
		logger.Warning(customNetworkingNotice)
		return nil
	}

	// use subnets as specified by --vpc-{private,public}-subnets flags

	if len(params.AvailabilityZones) != 0 {
		return fmt.Errorf("--vpc-private-subnets/--vpc-public-subnets and --zones %s", cmdutils.IncompatibleFlags)
	}
	if cmd.CobraCommand.Flag("vpc-cidr").Changed {
		return fmt.Errorf("--vpc-private-subnets/--vpc-public-subnets and --vpc-cidr %s", cmdutils.IncompatibleFlags)
	}

	if params.DryRun {
		if cfg.VPC.NAT != nil {
			disableNAT := api.ClusterDisableNAT
			cfg.VPC.NAT = &api.ClusterNAT{
				Gateway: &disableNAT,
			}
		}
		return nil
	}

	if err := vpc.ImportSubnetsFromSpec(ctx, ctl.AWSProvider, cfg); err != nil {
		return err
	}

	if err := cfg.HasSufficientSubnets(); err != nil {
		logger.Critical("unable to use given %s", cfg.SubnetInfo())
		return err
	}

	if err := cfg.CanUseForPrivateNodeGroups(); err != nil {
		return err
	}

	logger.Success("using existing %s", cfg.SubnetInfo())
	logger.Warning(customNetworkingNotice)
	return nil
}

func clientSetCreator(ctl *eks.ClusterProvider, cfg *api.ClusterConfig) func() (kubernetes.Interface, error) {
	var (
		err       error
		clientSet kubernetes.Interface
	)
	return func() (kubernetes.Interface, error) {
		if clientSet != nil || err != nil {
			return clientSet, err
		}
		clientSet, err = ctl.NewStdClientSet(cfg)
		return clientSet, err
	}
}

func checkSubnetsGivenAsFlags(params *cmdutils.CreateClusterCmdParams) bool {
	return len(*params.Subnets[api.SubnetTopologyPrivate])+len(*params.Subnets[api.SubnetTopologyPublic]) != 0
}
