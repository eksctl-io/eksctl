package create

import (
	"fmt"
	"strings"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/eks"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/kops"
	"github.com/weaveworks/eksctl/pkg/printers"
	"github.com/weaveworks/eksctl/pkg/utils"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

type createClusterCmdParams struct {
	writeKubeconfig             bool
	kubeconfigPath              string
	autoKubeconfigPath          bool
	authenticatorRoleARN        string
	setContext                  bool
	availabilityZones           []string
	installWindowsVPCController bool
	kopsClusterNameForVPC       string
	subnets                     map[api.SubnetTopology]*[]string
	withoutNodeGroup            bool
}

func createClusterCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	ng := api.NewNodeGroup()
	cmd.ClusterConfig = cfg

	params := &createClusterCmdParams{}

	cmd.SetDescription("cluster", "Create a cluster", "")

	cmd.SetRunFuncWithNameArg(func() error {
		return doCreateCluster(cmd, ng, params)
	})

	exampleClusterName := cmdutils.ClusterName("", "")
	exampleNodeGroupName := cmdutils.NodeGroupName("", "")

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVarP(&cfg.Metadata.Name, "name", "n", "", fmt.Sprintf("EKS cluster name (generated if unspecified, e.g. %q)", exampleClusterName))
		fs.StringToStringVarP(&cfg.Metadata.Tags, "tags", "", map[string]string{}, `A list of KV pairs used to tag the AWS resources (e.g. "Owner=John Doe,Team=Some Team")`)
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		fs.StringSliceVar(&params.availabilityZones, "zones", nil, "(auto-select if unspecified)")
		cmdutils.AddVersionFlag(fs, cfg.Metadata, "")
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
		fs.BoolVarP(&params.installWindowsVPCController, "install-vpc-controllers", "", false, "Install VPC controller that's required for Windows workloads")
	})

	cmd.FlagSetGroup.InFlagSet("Initial nodegroup", func(fs *pflag.FlagSet) {
		fs.StringVar(&ng.Name, "nodegroup-name", "", fmt.Sprintf("name of the nodegroup (generated if unspecified, e.g. %q)", exampleNodeGroupName))
		fs.BoolVar(&params.withoutNodeGroup, "without-nodegroup", false, "if set, initial nodegroup will not be created")
		cmdutils.AddCommonCreateNodeGroupFlags(fs, cmd, ng)
	})

	cmd.FlagSetGroup.InFlagSet("Cluster and nodegroup add-ons", func(fs *pflag.FlagSet) {
		cmdutils.AddCommonCreateNodeGroupIAMAddonsFlags(fs, ng)
	})

	cmd.FlagSetGroup.InFlagSet("VPC networking", func(fs *pflag.FlagSet) {
		fs.IPNetVar(&cfg.VPC.CIDR.IPNet, "vpc-cidr", cfg.VPC.CIDR.IPNet, "global CIDR to use for VPC")
		params.subnets = map[api.SubnetTopology]*[]string{
			api.SubnetTopologyPrivate: fs.StringSlice("vpc-private-subnets", nil, "re-use private subnets of an existing VPC"),
			api.SubnetTopologyPublic:  fs.StringSlice("vpc-public-subnets", nil, "re-use public subnets of an existing VPC"),
		}
		fs.StringVar(&params.kopsClusterNameForVPC, "vpc-from-kops-cluster", "", "re-use VPC from a given kops cluster")
		fs.StringVar(cfg.VPC.NAT.Gateway, "vpc-nat-mode", api.ClusterSingleNAT, "VPC NAT mode, valid options: HighlyAvailable, Single, Disable")
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, true)

	cmd.FlagSetGroup.InFlagSet("Output kubeconfig", func(fs *pflag.FlagSet) {
		cmdutils.AddCommonFlagsForKubeconfig(fs, &params.kubeconfigPath, &params.authenticatorRoleARN, &params.setContext, &params.autoKubeconfigPath, exampleClusterName)
		fs.BoolVar(&params.writeKubeconfig, "write-kubeconfig", true, "toggle writing of kubeconfig")
	})
}

func doCreateCluster(cmd *cmdutils.Cmd, ng *api.NodeGroup, params *createClusterCmdParams) error {
	ngFilter := cmdutils.NewNodeGroupFilter()

	if err := cmdutils.NewCreateClusterLoader(cmd, ngFilter, ng, params.withoutNodeGroup).Load(); err != nil {
		return err
	}

	cfg := cmd.ClusterConfig
	meta := cmd.ClusterConfig.Metadata

	printer := printers.NewJSONPrinter()

	ctl, err := cmd.NewCtl()
	if err != nil {
		return err
	}
	cmdutils.LogRegionAndVersionInfo(meta)

	if cfg.Metadata.Version == "" {
		cfg.Metadata.Version = api.DefaultVersion
	}
	if cfg.Metadata.Version != api.DefaultVersion {
		if !isValidVersion(cfg.Metadata.Version) {
			if isDeprecatedVersion(cfg.Metadata.Version) {
				return fmt.Errorf("invalid version, %s is now deprecated, supported values: %s\nsee also: https://docs.aws.amazon.com/eks/latest/userguide/kubernetes-versions.html", cfg.Metadata.Version, strings.Join(api.SupportedVersions(), ", "))
			}
			return fmt.Errorf("invalid version, supported values: %s", strings.Join(api.SupportedVersions(), ", "))
		}
	}

	if err := cfg.ValidateClusterEndpointConfig(); err != nil {
		return err
	}

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if params.autoKubeconfigPath {
		if params.kubeconfigPath != kubeconfig.DefaultPath {
			return fmt.Errorf("--kubeconfig and --auto-kubeconfig %s", cmdutils.IncompatibleFlags)
		}
		params.kubeconfigPath = kubeconfig.AutoPath(meta.Name)
	}

	if checkSubnetsGivenAsFlags(params) {
		// undo defaulting and reset it, as it's not set via config file;
		// default value here causes errors as vpc.ImportVPC doesn't
		// treat remote state as authority over local state
		cfg.VPC.CIDR = nil
		// load subnets from local map created from flags, into the config
		for topology := range params.subnets {
			if err := vpc.ImportSubnetsFromList(ctl.Provider, cfg, topology, *params.subnets[topology]); err != nil {
				return err
			}
		}
	}
	filteredNodeGroups := ngFilter.FilterMatching(cfg.NodeGroups)

	if err := eks.ValidateWindowsCompatibility(filteredNodeGroups, cfg.Metadata.Version); err != nil {
		return err
	}
	if params.installWindowsVPCController {
		if !eks.SupportsWindowsWorkloads(filteredNodeGroups) {
			return errors.New("running Windows workloads requires having both Windows and Linux (AmazonLinux2) node groups")
		}
	} else {
		eks.LogWindowsCompatibility(filteredNodeGroups, cfg.Metadata)
	}

	subnetsGiven := cfg.HasAnySubnets() // this will be false when neither flags nor config has any subnets

	createOrImportVPC := func() error {

		subnetInfo := func() string {
			return fmt.Sprintf("VPC (%s) and subnets (private:%v public:%v)",
				cfg.VPC.ID, cfg.PrivateSubnetIDs(), cfg.PublicSubnetIDs())
		}

		customNetworkingNotice := "custom VPC/subnets will be used; if resulting cluster doesn't function as expected, make sure to review the configuration of VPC/subnets"

		canUseForPrivateNodeGroups := func(ng *api.NodeGroup) error {
			if ng.PrivateNetworking && !cfg.HasSufficientPrivateSubnets() {
				return errors.New("none or too few private subnets to use with --node-private-networking")
			}
			return nil
		}

		if !subnetsGiven && params.kopsClusterNameForVPC == "" {
			// default: create dedicated VPC
			if err := ctl.SetAvailabilityZones(cfg, params.availabilityZones); err != nil {
				return err
			}
			if err := vpc.SetSubnets(cfg); err != nil {
				return err
			}
			return nil
		}

		if params.kopsClusterNameForVPC != "" {
			// import VPC from a given kops cluster
			if len(params.availabilityZones) != 0 {
				return fmt.Errorf("--vpc-from-kops-cluster and --zones %s", cmdutils.IncompatibleFlags)
			}
			if cmd.CobraCommand.Flag("vpc-cidr").Changed {
				return fmt.Errorf("--vpc-from-kops-cluster and --vpc-cidr %s", cmdutils.IncompatibleFlags)
			}

			if subnetsGiven {
				return fmt.Errorf("--vpc-from-kops-cluster and --vpc-private-subnets/--vpc-public-subnets %s", cmdutils.IncompatibleFlags)
			}

			kw, err := kops.NewWrapper(cmd.ProviderConfig.Region, params.kopsClusterNameForVPC)
			if err != nil {
				return err
			}

			if err := kw.UseVPC(ctl.Provider, cfg); err != nil {
				return err
			}

			for _, ng := range filteredNodeGroups {
				if err := canUseForPrivateNodeGroups(ng); err != nil {
					return err
				}
			}

			logger.Success("using %s from kops cluster %q", subnetInfo(), params.kopsClusterNameForVPC)
			logger.Warning(customNetworkingNotice)
			return nil
		}

		// use subnets as specified by --vpc-{private,public}-subnets flags

		if len(params.availabilityZones) != 0 {
			return fmt.Errorf("--vpc-private-subnets/--vpc-public-subnets and --zones %s", cmdutils.IncompatibleFlags)
		}
		if cmd.CobraCommand.Flag("vpc-cidr").Changed {
			return fmt.Errorf("--vpc-private-subnets/--vpc-public-subnets and --vpc-cidr %s", cmdutils.IncompatibleFlags)
		}

		if err := vpc.ImportAllSubnets(ctl.Provider, cfg); err != nil {
			return err
		}

		if err := cfg.HasSufficientSubnets(); err != nil {
			logger.Critical("unable to use given %s", subnetInfo())
			return err
		}

		for _, ng := range filteredNodeGroups {
			if err := canUseForPrivateNodeGroups(ng); err != nil {
				return err
			}
		}

		logger.Success("using existing %s", subnetInfo())
		logger.Warning(customNetworkingNotice)
		return nil
	}

	if err := createOrImportVPC(); err != nil {
		return err
	}

	for _, ng := range filteredNodeGroups {
		// resolve AMI
		if err := ctl.EnsureAMI(meta.Version, ng); err != nil {
			return err
		}
		logger.Info("nodegroup %q will use %q [%s/%s]", ng.Name, ng.AMI, ng.AMIFamily, cfg.Metadata.Version)

		if err := ctl.SetNodeLabels(ng, meta); err != nil {
			return err
		}

		// load or use SSH key - name includes cluster name and the
		// fingerprint, so if unique keys provided, each will get
		// loaded and used as intended and there is no need to have
		// nodegroup name in the key name
		if err := loadSSHKey(ng, meta.Name, ctl.Provider); err != nil {
			return err
		}
	}

	logger.Info("using Kubernetes version %s", meta.Version)
	logger.Info("creating %s", meta.LogString())

	// TODO dry-run mode should provide a way to render config with all defaults set
	// we should also make a call to resolve the AMI and write the result, similarly
	// the body of the SSH key can be read

	if err := printer.LogObj(logger.Debug, "cfg.json = \\\n%s\n", cfg); err != nil {
		return err
	}

	{ // core action
		stackManager := ctl.NewStackManager(cfg)
		if len(filteredNodeGroups) == 1 && cmd.ClusterConfigFile == "" {
			logger.Info("will create 2 separate CloudFormation stacks for cluster itself and the initial nodegroup")
		} else {
			ngFilter.LogInfo(cfg.NodeGroups)
			logger.Info("will create a CloudFormation stack for cluster itself and %d nodegroup stack(s)", len(filteredNodeGroups))
		}
		logger.Info("if you encounter any issues, check CloudFormation console or try 'eksctl utils describe-stacks --region=%s --cluster=%s'", meta.Region, meta.Name)
		tasks := stackManager.NewTasksToCreateClusterWithNodeGroups(filteredNodeGroups)
		ctl.AppendExtraClusterConfigTasks(cfg, params.installWindowsVPCController, tasks)

		logger.Info(tasks.Describe())
		if errs := tasks.DoAllSync(); len(errs) > 0 {
			logger.Info("%d error(s) occurred and cluster hasn't been created properly, you may wish to check CloudFormation console", len(errs))
			logger.Info("to cleanup resources, run 'eksctl delete cluster --region=%s --name=%s'", meta.Region, meta.Name)
			for _, err := range errs {
				logger.Critical("%s\n", err.Error())
			}
			return fmt.Errorf("failed to create cluster %q", meta.Name)
		}
	}

	logger.Success("all EKS cluster resources for %q have been created", meta.Name)

	// obtain cluster credentials, write kubeconfig

	{ // post-creation action
		var kubeconfigContextName string

		if params.writeKubeconfig {
			kubectlConfig := kubeconfig.NewForKubectl(cfg, ctl.GetUsername(), params.authenticatorRoleARN, ctl.Provider.Profile())
			kubeconfigContextName = kubectlConfig.CurrentContext

			params.kubeconfigPath, err = kubeconfig.Write(params.kubeconfigPath, *kubectlConfig, params.setContext)
			if err != nil {
				return errors.Wrap(err, "writing kubeconfig")
			}
			logger.Success("saved kubeconfig as %q", params.kubeconfigPath)
		} else {
			params.kubeconfigPath = ""
		}

		// create Kubernetes client
		clientSet, err := ctl.NewStdClientSet(cfg)
		if err != nil {
			return err
		}

		if err = ctl.WaitForControlPlane(meta, clientSet); err != nil {
			return err
		}

		for _, ng := range filteredNodeGroups {
			// authorise nodes to join
			if err = authconfigmap.AddNodeGroup(clientSet, ng); err != nil {
				return err
			}

			// wait for nodes to join
			if err = ctl.WaitForNodes(clientSet, ng); err != nil {
				return err
			}

			// if GPU instance type, give instructions
			if utils.IsGPUInstanceType(ng.InstanceType) || (ng.InstancesDistribution != nil && utils.HasGPUInstanceType(ng.InstancesDistribution.InstanceTypes)) {
				logger.Info("as you are using a GPU optimized instance type you will need to install NVIDIA Kubernetes device plugin.")
				logger.Info("\t see the following page for instructions: https://github.com/NVIDIA/k8s-device-plugin")
			}

		}

		// check kubectl version, and offer install instructions if missing or old
		// also check heptio-authenticator
		// TODO: https://github.com/weaveworks/eksctl/issues/30
		env, err := ctl.GetCredentialsEnv()
		if err != nil {
			return err
		}
		if err := utils.CheckAllCommands(params.kubeconfigPath, params.setContext, kubeconfigContextName, env); err != nil {
			logger.Critical("%s\n", err.Error())
			logger.Info("cluster should be functional despite missing (or misconfigured) client binaries")
		}
	}

	logger.Success("%s is ready", meta.LogString())

	if err := printer.LogObj(logger.Debug, "cfg.json = \\\n%s\n", cfg); err != nil {
		return err
	}

	return nil
}
