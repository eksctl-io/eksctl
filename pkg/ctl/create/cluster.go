package create

import (
	"fmt"
	"os"
	"strings"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/ami"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/eks/api"
	"github.com/weaveworks/eksctl/pkg/kops"
	"github.com/weaveworks/eksctl/pkg/utils"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

const (
	defaultNodeType     = "m5.large"
	defaultSSHPublicKey = "~/.ssh/id_rsa.pub"
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

func createClusterCmd(g *cmdutils.Grouping) *cobra.Command {
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

	group := g.New(cmd)

	exampleClusterName := utils.ClusterName("", "")

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVarP(&cfg.Metadata.Name, "name", "n", "", fmt.Sprintf("EKS cluster name (generated if unspecified, e.g. %q)", exampleClusterName))
		fs.StringToStringVarP(&cfg.Metadata.Tags, "tags", "", map[string]string{}, `A list of KV pairs used to tag the AWS resources (e.g. "Owner=John Doe,Team=Some Team")`)
		cmdutils.AddRegionFlag(fs, p)
		fs.StringSliceVar(&availabilityZones, "zones", nil, "(auto-select if unspecified)")
		fs.StringVar(&cfg.Metadata.Version, "version", api.LatestVersion, fmt.Sprintf("Kubernetes version (valid options: %s)", strings.Join(api.SupportedVersions(), ",")))
	})

	group.InFlagSet("Initial nodegroup", func(fs *pflag.FlagSet) {
		fs.IntVarP(&ng.DesiredCapacity, "nodes", "N", api.DefaultNodeCount, "total number of nodes (desired capacity of ASG)")

		// TODO: https://github.com/weaveworks/eksctl/issues/28
		fs.IntVarP(&ng.MinSize, "nodes-min", "m", 0, "minimum nodes in ASG (leave unset for a static nodegroup)")
		fs.IntVarP(&ng.MaxSize, "nodes-max", "M", 0, "maximum nodes in ASG (leave unset for a static nodegroup)")

		fs.StringVarP(&ng.InstanceType, "node-type", "t", defaultNodeType, "node instance type")

		fs.IntVarP(&ng.VolumeSize, "node-volume-size", "", 0, "Node volume size (in GB)")
		fs.IntVar(&ng.MaxPodsPerNode, "max-pods-per-node", 0, "maximum number of pods per node (set automatically if unspecified)")

		fs.StringVar(&ng.AMI, "node-ami", ami.ResolverStatic, "Advanced use cases only. If 'static' is supplied (default) then eksctl will use static AMIs; if 'auto' is supplied then eksctl will automatically set the AMI based on version/region/instance type; if any other value is supplied it will override the AMI to use for the nodes. Use with extreme care.")
		fs.StringVar(&ng.AMIFamily, "node-ami-family", ami.ImageFamilyAmazonLinux2, "Advanced use cases only. If 'AmazonLinux2' is supplied (default), then eksctl will use the offical AWS EKS AMIs (Amazon Linux 2); if 'Ubuntu1804' is supplied, then eksctl will use the offical Canonical EKS AMIs (Ubuntu 18.04).")

		fs.BoolVar(&ng.AllowSSH, "ssh-access", false, "control SSH access for nodes")
		fs.StringVar(&ng.SSHPublicKeyPath, "ssh-public-key", defaultSSHPublicKey, "SSH public key to use for nodes (import from local path, or use existing EC2 key pair)")

		fs.BoolVarP(&ng.PrivateNetworking, "node-private-networking", "P", false, "whether to make initial nodegroup networking private")
	})

	group.InFlagSet("Cluster add-ons", func(fs *pflag.FlagSet) {
		fs.BoolVar(&cfg.Addons.WithIAM.PolicyAutoScaling, "asg-access", false, "enable iam policy dependency for cluster-autoscaler")
		fs.BoolVar(&cfg.Addons.WithIAM.PolicyExternalDNS, "external-dns-access", false, "enable iam policy dependency for external-dns")
		fs.BoolVar(&cfg.Addons.WithIAM.PolicyAmazonEC2ContainerRegistryPowerUser, "full-ecr-access", false, "enable full access to ECR")
		fs.BoolVar(&cfg.Addons.Storage, "storage-class", true, "if true (default) then a default StorageClass of type gp2 provisioned by EBS will be created")
	})

	group.InFlagSet("VPC networking", func(fs *pflag.FlagSet) {
		fs.IPNetVar(cfg.VPC.CIDR, "vpc-cidr", api.DefaultCIDR(), "global CIDR to use for VPC")
		subnets = map[api.SubnetTopology]*[]string{
			api.SubnetTopologyPrivate: fs.StringSlice("vpc-private-subnets", nil, "re-use private subnets of an existing VPC"),
			api.SubnetTopologyPublic:  fs.StringSlice("vpc-public-subnets", nil, "re-use public subnets of an existing VPC"),
		}
		fs.StringVar(&kopsClusterNameForVPC, "vpc-from-kops-cluster", "", "re-use VPC from a given kops cluster")
	})

	cmdutils.AddCommonFlagsForAWS(group, p)

	group.InFlagSet("Output kubeconfig", func(fs *pflag.FlagSet) {
		cmdutils.AddCommonFlagsForKubeconfig(fs, &kubeconfigPath, &setContext, &autoKubeconfigPath, exampleClusterName)
		fs.BoolVar(&writeKubeconfig, "write-kubeconfig", true, "toggle writing of kubeconfig")
	})

	group.AddTo(cmd)

	return cmd
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

	if err := ctl.EnsureAMI(meta.Version, ng); err != nil {
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

		// add default storage class only for version 1.10 clusters
		if cfg.Addons.Storage && meta.Version == "1.10" {
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
