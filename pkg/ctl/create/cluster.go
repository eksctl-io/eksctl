package create

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/eks/api"
	"github.com/weaveworks/eksctl/pkg/kops"
	"github.com/weaveworks/eksctl/pkg/utils"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

var (
	writeKubeconfig    bool
	kubeconfigPath     string
	autoKubeconfigPath bool
	setContext         bool
	availabilityZones  []string

	kopsClusterNameForVPC string
	subnets               map[api.SubnetTopology]*[]string
)

func createClusterCmd() *cobra.Command {
	p := &api.ProviderConfig{}
	cfg := api.NewClusterConfig()
	ng := cfg.NewNodeGroup()

	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Create a cluster",
		Run: func(_ *cobra.Command, args []string) {
			if err := doCreateCluster(p, cfg, ng, cmdutils.GetNameArg(args)); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	fs := cmd.Flags()

	exampleClusterName := utils.ClusterName("", "")

	fs.StringVarP(&cfg.Metadata.Name, "name", "n", "", fmt.Sprintf("EKS cluster name (generated if unspecified, e.g. %q)", exampleClusterName))

	addCommonCreateFlags(fs, p, cfg, ng)

	fs.BoolVar(&writeKubeconfig, "write-kubeconfig", true, "toggle writing of kubeconfig")
	cmdutils.AddCommonFlagsForKubeconfig(fs, &kubeconfigPath, &setContext, &autoKubeconfigPath, exampleClusterName)

	//TODO: is this needed for nodegroup creation????
	fs.StringSliceVar(&availabilityZones, "zones", nil, "(auto-select if unspecified)")

	fs.BoolVar(&cfg.Addons.WithIAM.PolicyAmazonEC2ContainerRegistryPowerUser, "full-ecr-access", false, "enable full access to ECR")
	fs.BoolVar(&cfg.Addons.WithIAM.PolicyAutoScaling, "asg-access", false, "enable iam policy dependency for cluster-autoscaler")
	fs.BoolVar(&cfg.Addons.Storage, "storage-class", true, "if true (default) then a default StorageClass of type gp2 provisioned by EBS will be created")

	fs.StringVar(&kopsClusterNameForVPC, "vpc-from-kops-cluster", "", "re-use VPC from a given kops cluster")

	fs.IPNetVar(cfg.VPC.CIDR, "vpc-cidr", api.DefaultCIDR(), "global CIDR to use for VPC")

	subnets = map[api.SubnetTopology]*[]string{
		api.SubnetTopologyPrivate: fs.StringSlice("vpc-private-subnets", nil, "re-use private subnets of an existing VPC"),
		api.SubnetTopologyPublic:  fs.StringSlice("vpc-public-subnets", nil, "re-use public subnets of an existing VPC"),
	}

	groupFlagsInUsage(cmd)

	return cmd
}

func groupFlagsInUsage(cmd *cobra.Command) {
	// Group flags by their categories determined by name prefixes
	groupToPatterns := map[string][]string{
		"Node":       {"node", "storage-class", "ssh", "max-pods-per-node", "full-ecr-access", "asg-access"},
		"Networking": {"vpc", "zones",},
		"Stack":      {"region", "tags"},
		"Other":      {},
	}
	groups := []string{}
	for k := range groupToPatterns {
		groups = append(groups, k)
	}
	groupToFlagSet := make(map[string]*pflag.FlagSet)
	for _, g := range groups {
		groupToFlagSet[g] = pflag.NewFlagSet(g, /* Unused. Can be anythng. */ pflag.ContinueOnError)
	}
	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		for _, g := range groups {
			for _, p := range groupToPatterns[g] {
				if strings.HasPrefix(f.Name, p) {
					groupToFlagSet[g].AddFlag(f)
					return
				}
			}
		}
		groupToFlagSet["Other"].AddFlag(f)
	})

	// The usage template is based on the one bundled into cobra
	// https://github.com/spf13/cobra/blob/1e58aa3361fd650121dceeedc399e7189c05674a/command.go#L397
	origFlagUsages := `

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}`

	altFlagUsages := ``
	for _, g := range groups {
		set := groupToFlagSet[g]
		altFlagUsages += fmt.Sprintf(`

%s Flags:
%s`, g, strings.TrimRightFunc(set.FlagUsages(), unicode.IsSpace))
	}

	cmd.SetUsageTemplate(strings.Replace(cmd.UsageTemplate(), origFlagUsages, altFlagUsages, 1))
}

func doCreateCluster(p *api.ProviderConfig, cfg *api.ClusterConfig, ng *api.NodeGroup, nameArg string) error {
	meta := cfg.Metadata
	ctl := eks.New(p, cfg)

	if !ctl.IsSupportedRegion() {
		return cmdutils.ErrUnsupportedRegion(p)
	}
	logger.Info("using region %s", meta.Region)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if utils.ClusterName(meta.Name, nameArg) == "" {
		return cmdutils.ErrNameFlagAndArg(meta.Name, nameArg)
	}
	meta.Name = utils.ClusterName(meta.Name, nameArg)

	if autoKubeconfigPath {
		if kubeconfigPath != kubeconfig.DefaultPath {
			return fmt.Errorf("--kubeconfig and --auto-kubeconfig %s", cmdutils.IncompatibleFlags)
		}
		kubeconfigPath = kubeconfig.AutoPath(meta.Name)
	}

	if ng.SSHPublicKeyPath == "" {
		return fmt.Errorf("--ssh-public-key must be non-empty string")
	}

	createOrImportVPC := func() error {
		subnetsGiven := len(*subnets[api.SubnetTopologyPrivate])+len(*subnets[api.SubnetTopologyPublic]) != 0

		subnetInfo := func() string {
			return fmt.Sprintf("VPC (%s) and subnets (private:%v public:%v)",
				cfg.VPC.ID, cfg.SubnetIDs(api.SubnetTopologyPrivate), cfg.SubnetIDs(api.SubnetTopologyPublic))
		}

		customNetworkingNotice := "custom VPC/subnets will be used; if resulting cluster doesn't function as expected, make sure to review the configuration of VPC/subnets"

		canUseForPrivateNodeGroup := func() error {
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

			if err := canUseForPrivateNodeGroup(); err != nil {
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

		for topology := range subnets {
			if err := vpc.UseSubnets(ctl.Provider, cfg, topology, *subnets[topology]); err != nil {
				return err
			}
		}

		if err := cfg.HasSufficientSubnets(); err != nil {
			logger.Critical("unable to use given %s", subnetInfo())
			return err
		}

		if err := canUseForPrivateNodeGroup(); err != nil {
			return err
		}

		logger.Success("using existing %s", subnetInfo())
		logger.Warning(customNetworkingNotice)
		return nil
	}

	if err := createOrImportVPC(); err != nil {
		return err
	}

	if err := ctl.EnsureAMI(ng); err != nil {
		return err
	}

	if err := ctl.LoadSSHPublicKey(meta.Name, ng); err != nil {
		return err
	}

	logger.Debug("cfg = %#v", cfg)

	logger.Info("creating %s", meta.LogString())

	{ // core action
		stackManager := ctl.NewStackManager(cfg)
		logger.Info("will create 2 separate CloudFormation stacks for cluster itself and the initial nodegroup")
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
		clientConfigBase, err := ctl.NewClientConfig(cfg)
		if err != nil {
			return err
		}

		if writeKubeconfig {
			config := clientConfigBase.WithExecAuthenticator()
			kubeconfigPath, err = kubeconfig.Write(kubeconfigPath, config.Client, setContext)
			if err != nil {
				return errors.Wrap(err, "writing kubeconfig")
			}

			logger.Success("saved kubeconfig as %q", kubeconfigPath)
		} else {
			kubeconfigPath = ""
		}

		// create Kubernetes client
		clientSet, err := clientConfigBase.NewClientSetWithEmbeddedToken()
		if err != nil {
			return err
		}

		if err = ctl.WaitForControlPlane(meta, clientSet); err != nil {
			return err
		}

		// authorise nodes to join
		if err = ctl.CreateNodeGroupAuthConfigMap(clientSet, ng); err != nil {
			return err
		}

		// wait for nodes to join
		if err = ctl.WaitForNodes(clientSet, ng); err != nil {
			return err
		}

		// add default storage class
		if cfg.Addons.Storage {
			if err = ctl.AddDefaultStorageClass(clientSet); err != nil {
				return err
			}
		}

		// check kubectl version, and offer install instructions if missing or old
		// also check heptio-authenticator
		// TODO: https://github.com/weaveworks/eksctl/issues/30
		env, err := ctl.GetCredentialsEnv()
		if err != nil {
			return err
		}
		if err := utils.CheckAllCommands(kubeconfigPath, setContext, clientConfigBase.ContextName, env); err != nil {
			logger.Critical("%s\n", err.Error())
			logger.Info("cluster should be functional despite missing (or misconfigured) client binaries")
		}

		// If GPU instance type, give instructions
		if utils.IsGPUInstanceType(ng.InstanceType) {
			logger.Info("as you are using a GPU optimized instance type you will need to install NVIDIA Kubernetes device plugin.")
			logger.Info("\t see the following page for instructions: https://github.com/NVIDIA/k8s-device-plugin")
		}
	}

	logger.Success("%s is ready", meta.LogString())

	return nil
}
