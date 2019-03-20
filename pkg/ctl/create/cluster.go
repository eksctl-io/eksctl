package create

import (
	"fmt"
	"os"
	"strings"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/kops"
	"github.com/weaveworks/eksctl/pkg/printers"
	"github.com/weaveworks/eksctl/pkg/utils"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

const (
	defaultSSHPublicKey = "~/.ssh/id_rsa.pub"
)

var (
	writeKubeconfig    bool
	kubeconfigPath     string
	autoKubeconfigPath bool
	setContext         bool
	availabilityZones  []string

	clusterConfigFile = ""

	kopsClusterNameForVPC string
	subnets               map[api.SubnetTopology]*[]string
	addonsStorageClass    bool
	withoutNodeGroup      bool
)

func createClusterCmd(g *cmdutils.Grouping) *cobra.Command {
	p := &api.ProviderConfig{}
	cfg := api.NewClusterConfig()
	ng := cfg.NewNodeGroup()

	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Create a cluster",
		Run: func(cmd *cobra.Command, args []string) {
			if err := doCreateCluster(p, cfg, cmdutils.GetNameArg(args), cmd); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	group := g.New(cmd)

	exampleClusterName := utils.ClusterName("", "")
	exampleNodeGroupName := utils.NodeGroupName("", "")

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVarP(&cfg.Metadata.Name, "name", "n", "", fmt.Sprintf("EKS cluster name (generated if unspecified, e.g. %q)", exampleClusterName))
		fs.StringToStringVarP(&cfg.Metadata.Tags, "tags", "", map[string]string{}, `A list of KV pairs used to tag the AWS resources (e.g. "Owner=John Doe,Team=Some Team")`)
		cmdutils.AddRegionFlag(fs, p)
		fs.StringSliceVar(&availabilityZones, "zones", nil, "(auto-select if unspecified)")
		cmdutils.AddVersionFlag(fs, cfg.Metadata, "")
		cmdutils.AddConfigFileFlag(&clusterConfigFile, fs)
	})

	group.InFlagSet("Initial nodegroup", func(fs *pflag.FlagSet) {
		fs.StringVar(&ng.Name, "nodegroup-name", "", fmt.Sprintf("name of the nodegroup (generated if unspecified, e.g. %q)", exampleNodeGroupName))
		fs.BoolVar(&withoutNodeGroup, "without-nodegroup", false, "if set, initial nodegroup will not be created")
		cmdutils.AddCommonCreateNodeGroupFlags(cmd, fs, p, cfg, ng)
	})

	group.InFlagSet("Cluster and nodegroup add-ons", func(fs *pflag.FlagSet) {
		cmdutils.AddCommonCreateNodeGroupIAMAddonsFlags(fs, ng)
		fs.BoolVar(&addonsStorageClass, "storage-class", true, "if true (default) then a default StorageClass of type gp2 provisioned by EBS will be created")
	})

	group.InFlagSet("VPC networking", func(fs *pflag.FlagSet) {
		fs.IPNetVar(&cfg.VPC.CIDR.IPNet, "vpc-cidr", cfg.VPC.CIDR.IPNet, "global CIDR to use for VPC")
		subnets = map[api.SubnetTopology]*[]string{
			api.SubnetTopologyPrivate: fs.StringSlice("vpc-private-subnets", nil, "re-use private subnets of an existing VPC"),
			api.SubnetTopologyPublic:  fs.StringSlice("vpc-public-subnets", nil, "re-use public subnets of an existing VPC"),
		}
		fs.StringVar(&kopsClusterNameForVPC, "vpc-from-kops-cluster", "", "re-use VPC from a given kops cluster")
	})

	cmdutils.AddCommonFlagsForAWS(group, p, true)

	group.InFlagSet("Output kubeconfig", func(fs *pflag.FlagSet) {
		cmdutils.AddCommonFlagsForKubeconfig(fs, &kubeconfigPath, &setContext, &autoKubeconfigPath, exampleClusterName)
		fs.BoolVar(&writeKubeconfig, "write-kubeconfig", true, "toggle writing of kubeconfig")
	})

	group.AddTo(cmd)

	return cmd
}

func doCreateCluster(p *api.ProviderConfig, cfg *api.ClusterConfig, nameArg string, cmd *cobra.Command) error {
	meta := cfg.Metadata

	printer := printers.NewJSONPrinter()

	if err := api.Register(); err != nil {
		return err
	}

	if clusterConfigFile != "" {
		if err := eks.LoadConfigFromFile(clusterConfigFile, cfg); err != nil {
			return err
		}
		meta = cfg.Metadata

		incompatibleFlags := []string{
			"name",
			"tags",
			"zones",
			"version",
			"region",
			"nodes",
			"nodes-min",
			"nodes-max",
			"node-type",
			"node-volume-size",
			"node-volume-type",
			"max-pods-per-node",
			"node-ami",
			"node-ami-family",
			"ssh-access",
			"ssh-public-key",
			"node-private-networking",
			"node-security-groups",
			"node-labels",
			"node-zones",
			"asg-access",
			"external-dns-access",
			"full-ecr-access",
			"storage-class",
			"vpc-private-subnets",
			"vpc-public-subnets",
			"vpc-cidr",
			"vpc-from-kops-cluster",
		}

		for _, f := range incompatibleFlags {
			if cmd.Flag(f).Changed {
				return fmt.Errorf("cannot use --%s when --config-file/-f is set", f)
			}
		}

		if meta.Name == "" {
			return fmt.Errorf("metadata.name must be set")
		}

		if meta.Region == "" {
			return fmt.Errorf("metadata.region must be set")
		}

		p.Region = meta.Region

		if cfg.VPC == nil {
			cfg.VPC = api.NewClusterVPC()
		}

		if checkSubnetsGiven(cfg) && len(cfg.AvailabilityZones) != 0 {
			return fmt.Errorf("vpc.subnets and availabilityZones cannot be set at the same time")
		}

		skipNodeGroupsIfRequested(cfg)

		if err := checkEachNodeGroup(cfg, newNodeGroupChecker); err != nil {
			return err
		}
	} else {
		// validation and defaulting specific to when --config-file is unused

		// generate cluster name or use either flag or argument
		if utils.ClusterName(meta.Name, nameArg) == "" {
			return cmdutils.ErrNameFlagAndArg(meta.Name, nameArg)
		}
		meta.Name = utils.ClusterName(meta.Name, nameArg)

		if cfg.Status != nil {
			return fmt.Errorf("status fields are read-only")
		}

		skipNodeGroupsIfRequested(cfg)

		err := checkEachNodeGroup(cfg, func(i int, ng *api.NodeGroup) error {
			if ng.AllowSSH && ng.SSHPublicKeyPath == "" {
				return fmt.Errorf("--ssh-public-key must be non-empty string")
			}

			// generate nodegroup name or use flag
			ng.Name = utils.NodeGroupName(ng.Name, "")

			return nil
		})
		if err != nil {
			return err
		}
	}

	ctl := eks.New(p, cfg)

	if !ctl.IsSupportedRegion() {
		return cmdutils.ErrUnsupportedRegion(p)
	}
	logger.Info("using region %s", meta.Region)

	if cfg.Metadata.Version == "" {
		cfg.Metadata.Version = api.LatestVersion
	}
	if cfg.Metadata.Version != api.LatestVersion {
		validVersion := false
		for _, v := range api.SupportedVersions() {
			if cfg.Metadata.Version == v {
				validVersion = true
			}
		}
		if !validVersion {
			return fmt.Errorf("invalid version, supported values: %s", strings.Join(api.SupportedVersions(), ", "))
		}
	}

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if autoKubeconfigPath {
		if kubeconfigPath != kubeconfig.DefaultPath {
			return fmt.Errorf("--kubeconfig and --auto-kubeconfig %s", cmdutils.IncompatibleFlags)
		}
		kubeconfigPath = kubeconfig.AutoPath(meta.Name)
	}

	if checkSubnetsGivenAsFlags() {
		// undo defaulting and reset it, as it's not set via config file;
		// default value here causes errors as vpc.ImportVPC doesn't
		// treat remote state as authority over local state
		cfg.VPC.CIDR = nil
		// load subnets from local map created from flags, into the config
		for topology := range subnets {
			if err := vpc.ImportSubnetsFromList(ctl.Provider, cfg, topology, *subnets[topology]); err != nil {
				return err
			}
		}
	}
	subnetsGiven := checkSubnetsGiven(cfg) // this will be false when neither flags nor config has any subnets

	createOrImportVPC := func() error {

		subnetInfo := func() string {
			return fmt.Sprintf("VPC (%s) and subnets (private:%v public:%v)",
				cfg.VPC.ID, cfg.PrivateSubnetIDs(), cfg.PublicSubnetIDs())
		}

		customNetworkingNotice := "custom VPC/subnets will be used; if resulting cluster doesn't function as expected, make sure to review the configuration of VPC/subnets"

		canUseForPrivateNodeGroups := func(_ int, ng *api.NodeGroup) error {
			if ng.PrivateNetworking && !cfg.HasSufficientPrivateSubnets() {
				return fmt.Errorf("none or too few private subnets to use with --node-private-networking")
			}
			return nil
		}

		if !subnetsGiven && kopsClusterNameForVPC == "" {
			// default: create dedicated VPC
			if err := ctl.SetAvailabilityZones(cfg, availabilityZones); err != nil {
				return err
			}
			if err := vpc.SetSubnets(cfg); err != nil {
				return err
			}
			return nil
		}

		if kopsClusterNameForVPC != "" {
			// import VPC from a given kops cluster
			if len(availabilityZones) != 0 {
				return fmt.Errorf("--vpc-from-kops-cluster and --zones %s", cmdutils.IncompatibleFlags)
			}
			if cmd.Flag("vpc-cidr").Changed {
				return fmt.Errorf("--vpc-from-kops-cluster and --vpc-cidr %s", cmdutils.IncompatibleFlags)
			}

			if subnetsGiven {
				return fmt.Errorf("--vpc-from-kops-cluster and --vpc-private-subnets/--vpc-public-subnets %s", cmdutils.IncompatibleFlags)
			}

			kw, err := kops.NewWrapper(p.Region, kopsClusterNameForVPC)
			if err != nil {
				return err
			}

			if err := kw.UseVPC(ctl.Provider, cfg); err != nil {
				return err
			}

			if err := checkEachNodeGroup(cfg, canUseForPrivateNodeGroups); err != nil {
				return err
			}

			logger.Success("using %s from kops cluster %q", subnetInfo(), kopsClusterNameForVPC)
			logger.Warning(customNetworkingNotice)
			return nil
		}

		// use subnets as specified by --vpc-{private,public}-subnets flags

		if len(availabilityZones) != 0 {
			return fmt.Errorf("--vpc-private-subnets/--vpc-public-subnets and --zones %s", cmdutils.IncompatibleFlags)
		}
		if cmd.Flag("vpc-cidr").Changed {
			return fmt.Errorf("--vpc-private-subnets/--vpc-public-subnets and --vpc-cidr %s", cmdutils.IncompatibleFlags)
		}

		if err := vpc.ImportAllSubnets(ctl.Provider, cfg); err != nil {
			return err
		}

		if err := cfg.HasSufficientSubnets(); err != nil {
			logger.Critical("unable to use given %s", subnetInfo())
			return err
		}

		if err := checkEachNodeGroup(cfg, canUseForPrivateNodeGroups); err != nil {
			return err
		}

		logger.Success("using existing %s", subnetInfo())
		logger.Warning(customNetworkingNotice)
		return nil
	}

	if err := createOrImportVPC(); err != nil {
		return err
	}

	err := checkEachNodeGroup(cfg, func(_ int, ng *api.NodeGroup) error {
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
		if err := ctl.LoadSSHPublicKey(meta.Name, ng); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	logger.Info("creating %s", meta.LogString())

	// TODO dry-run mode should provide a way to render config with all defaults set
	// we should also make a call to resolve the AMI and write the result, similaraly
	// the body of the SSH key can be read

	if err := printer.LogObj(logger.Debug, "cfg.json = \\\n%s\n", cfg); err != nil {
		return err
	}

	{ // core action
		stackManager := ctl.NewStackManager(cfg)
		if len(cfg.NodeGroups) == 1 {
			logger.Info("will create 2 separate CloudFormation stacks for cluster itself and the initial nodegroup")
		} else {
			logger.Info("will create a CloudFormation stack for cluster itself and %d nodegroup stack(s)", len(cfg.NodeGroups))

		}
		logger.Info("if you encounter any issues, check CloudFormation console or try 'eksctl utils describe-stacks --region=%s --name=%s'", meta.Region, meta.Name)
		errs := stackManager.CreateClusterWithNodeGroups()
		// read any errors (it only gets non-nil errors)
		if len(errs) > 0 {
			logger.Info("%d error(s) occurred and cluster hasn't been created properly, you may wish to check CloudFormation console", len(errs))
			logger.Info("to cleanup resources, run 'eksctl delete cluster --region=%s --name=%s'", meta.Region, meta.Name)
			for _, err := range errs {
				logger.Critical("%s\n", err.Error())
			}
			return fmt.Errorf("failed to create cluster %q", meta.Name)
		}
	}

	logger.Success("all EKS cluster resource for %q had been created", meta.Name)

	// obtain cluster credentials, write kubeconfig

	{ // post-creation action
		var kubeconfigContextName string

		if writeKubeconfig {
			client, err := ctl.NewClient(cfg, false)
			if err != nil {
				return err
			}
			kubeconfigContextName = client.ContextName

			kubeconfigPath, err = kubeconfig.Write(kubeconfigPath, *client.Config, setContext)
			if err != nil {
				return errors.Wrap(err, "writing kubeconfig")
			}
			logger.Success("saved kubeconfig as %q", kubeconfigPath)
		} else {
			kubeconfigPath = ""
		}

		// create Kubernetes client
		clientSet, err := ctl.NewStdClientSet(cfg)
		if err != nil {
			return err
		}

		if err = ctl.WaitForControlPlane(meta, clientSet); err != nil {
			return err
		}

		err = checkEachNodeGroup(cfg, func(_ int, ng *api.NodeGroup) error {
			// authorise nodes to join
			if err = authconfigmap.AddNodeGroup(clientSet, ng); err != nil {
				return err
			}

			// wait for nodes to join
			if err = ctl.WaitForNodes(clientSet, ng); err != nil {
				return err
			}

			// if GPU instance type, give instructions
			if utils.IsGPUInstanceType(ng.InstanceType) {
				logger.Info("as you are using a GPU optimized instance type you will need to install NVIDIA Kubernetes device plugin.")
				logger.Info("\t see the following page for instructions: https://github.com/NVIDIA/k8s-device-plugin")
			}

			return nil
		})
		if err != nil {
			return err
		}

		// add default storage class only for version 1.10 clusters
		if meta.Version == "1.10" {
			// --storage-class flag is only for backwards compatibility,
			// we always create the storage class when --config-file is
			// used, as this is 1.10-only
			if addonsStorageClass || clusterConfigFile != "" {
				if err = ctl.AddDefaultStorageClass(clientSet); err != nil {
					return err
				}
			}
		}

		// check kubectl version, and offer install instructions if missing or old
		// also check heptio-authenticator
		// TODO: https://github.com/weaveworks/eksctl/issues/30
		env, err := ctl.GetCredentialsEnv()
		if err != nil {
			return err
		}
		if err := utils.CheckAllCommands(kubeconfigPath, setContext, kubeconfigContextName, env); err != nil {
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
