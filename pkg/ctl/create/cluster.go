package create

import (
	"fmt"
	"os"

	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/ami"
	"github.com/weaveworks/eksctl/pkg/ctl"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/eks/api"
	"github.com/weaveworks/eksctl/pkg/kops"
	"github.com/weaveworks/eksctl/pkg/utils"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
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
)

func createClusterCmd() *cobra.Command {
	cfg := api.NewClusterConfig()
	ng := cfg.NewNodeGroup()

	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Create a cluster",
		Run: func(_ *cobra.Command, args []string) {
			if err := doCreateCluster(cfg, ng, ctl.GetNameArg(args)); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	fs := cmd.Flags()

	exampleClusterName := utils.ClusterName("", "")

	fs.StringVarP(&cfg.ClusterName, "name", "n", "", fmt.Sprintf("EKS cluster name (generated if unspecified, e.g. %q)", exampleClusterName))

	fs.StringVarP(&cfg.Region, "region", "r", "", "AWS region")
	fs.StringVarP(&cfg.Profile, "profile", "p", "", "AWS credentials profile to use (overrides the AWS_PROFILE environment variable)")
	fs.StringToStringVarP(&cfg.Tags, "tags", "", map[string]string{}, `A list of KV pairs used to tag the AWS resources (e.g. "Owner=John Doe,Team=Some Team")`)

	fs.StringVarP(&ng.InstanceType, "node-type", "t", defaultNodeType, "node instance type")
	fs.IntVarP(&ng.DesiredCapacity, "nodes", "N", api.DefaultNodeCount, "total number of nodes (desired capacity of AGS)")

	// TODO: https://github.com/weaveworks/eksctl/issues/28
	fs.IntVarP(&ng.MinSize, "nodes-min", "m", 0, "minimum nodes in ASG (leave unset for a static nodegroup)")
	fs.IntVarP(&ng.MaxSize, "nodes-max", "M", 0, "maximum nodes in ASG (leave unset for a static nodegroup)")

	fs.IntVarP(&ng.VolumeSize, "node-volume-size", "", 0, "Node volume size (in GB)")
	fs.IntVar(&ng.MaxPodsPerNode, "max-pods-per-node", 0, "maximum number of pods per node (set automatically if unspecified)")
	fs.StringSliceVar(&availabilityZones, "zones", nil, "(auto-select if unspecified)")

	fs.BoolVar(&ng.AllowSSH, "ssh-access", false, "control SSH access for nodes")
	fs.StringVar(&ng.SSHPublicKeyPath, "ssh-public-key", defaultSSHPublicKey, "SSH public key to use for nodes (import from local path, or use existing EC2 key pair)")

	fs.BoolVar(&writeKubeconfig, "write-kubeconfig", true, "toggle writing of kubeconfig")
	fs.BoolVar(&autoKubeconfigPath, "auto-kubeconfig", false, fmt.Sprintf("save kubeconfig file by cluster name, e.g. %q", kubeconfig.AutoPath(exampleClusterName)))
	fs.StringVar(&kubeconfigPath, "kubeconfig", kubeconfig.DefaultPath, "path to write kubeconfig (incompatible with --auto-kubeconfig)")
	fs.BoolVar(&setContext, "set-kubeconfig-context", true, "if true then current-context will be set in kubeconfig; if a context is already set then it will be overwritten")

	fs.DurationVar(&cfg.WaitTimeout, "aws-api-timeout", api.DefaultWaitTimeout, "")
	// TODO deprecate in 0.2.0
	if err := fs.MarkHidden("aws-api-timeout"); err != nil {
		logger.Debug("ignoring error %q", err.Error())
	}
	fs.DurationVar(&cfg.WaitTimeout, "timeout", api.DefaultWaitTimeout, "max wait time in any polling operations")

	fs.BoolVar(&cfg.Addons.WithIAM.PolicyAmazonEC2ContainerRegistryPowerUser, "full-ecr-access", false, "enable full access to ECR")
	fs.BoolVar(&cfg.Addons.WithIAM.PolicyAutoScaling, "asg-access", false, "enable iam policy dependency for cluster-autoscaler")
	fs.BoolVar(&cfg.Addons.Storage, "storage-class", true, "if true (default) then a default StorageClass of type gp2 provisioned by EBS will be created")

	fs.StringVar(&ng.AMI, "node-ami", ami.ResolverStatic, "Advanced use cases only. If 'static' is supplied (default) then eksctl will use static AMIs; if 'auto' is supplied then eksctl will automatically set the AMI based on region/instance type; if any other value is supplied it will override the AMI to use for the nodes. Use with extreme care.")
	fs.StringVar(&ng.AMIFamily, "node-ami-family", ami.ImageFamilyAmazonLinux2, "Advanced use cases only. If 'AmazonLinux2' is supplied (default), then eksctl will use the offical AWS EKS AMIs (Amazon Linux 2); if 'Ubuntu1804' is supplied, then eksctl will use the offical Canonical EKS AMIs (Ubuntu 18.04).")

	fs.StringVar(&kopsClusterNameForVPC, "vpc-from-kops-cluster", "", "re-use VPC from a given kops cluster")

	fs.IPNetVar(cfg.VPC.CIDR, "vpc-cidr", api.DefaultCIDR(), "global CIDR to use for VPC")

	fs.BoolVarP(&ng.PrivateNetworking, "node-private-networking", "P", false, "whether to make initial nodegroup networking private")

	return cmd
}

func doCreateCluster(cfg *api.ClusterConfig, ng *api.NodeGroup, name string) error {
	ctl := eks.New(cfg)

	if !api.IsSupportedRegion(cfg.Region) {
		return fmt.Errorf("%s is not supported. Only %v are supported", cfg.Region, api.SupportedRegions)
	}

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if utils.ClusterName(cfg.ClusterName, name) == "" {
		return fmt.Errorf("--name=%s and argument %s cannot be used at the same time", cfg.ClusterName, name)
	}
	cfg.ClusterName = utils.ClusterName(cfg.ClusterName, name)

	if autoKubeconfigPath {
		if kubeconfigPath != kubeconfig.DefaultPath {
			return fmt.Errorf("--kubeconfig and --auto-kubeconfig cannot be used at the same time")
		}
		kubeconfigPath = kubeconfig.AutoPath(cfg.ClusterName)
	}

	if ng.SSHPublicKeyPath == "" {
		return fmt.Errorf("--ssh-public-key must be non-empty string")
	}

	if kopsClusterNameForVPC != "" {
		if len(availabilityZones) != 0 {
			return fmt.Errorf("--vpc-from-kops-cluster and --zones cannot be used at the same time")
		}
		kw, err := kops.NewWrapper(cfg.Region, kopsClusterNameForVPC)
		if err != nil {
			return err
		}

		if err := kw.UseVPC(cfg); err != nil {
			return err
		}
		logger.Success("using VPC (%s) and subnets (%v) from kops cluster %q", cfg.VPC.ID, cfg.SubnetIDs(api.SubnetTopologyPublic), kopsClusterNameForVPC)
	} else {
		// kw.UseVPC() sets AZs based on subenets used
		if err := ctl.SetAvailabilityZones(availabilityZones); err != nil {
			return err
		}
		if err := ctl.SetSubnets(); err != nil {
			return err
		}
	}

	if err := ctl.EnsureAMI(ng); err != nil {
		return err
	}

	if err := ctl.LoadSSHPublicKey(ng); err != nil {
		return err
	}

	logger.Debug("cfg = %#v", cfg)

	logger.Info("creating EKS cluster %q in %q region", cfg.ClusterName, cfg.Region)

	{ // core action
		stackManager := ctl.NewStackManager()
		logger.Info("will create 2 separate CloudFormation stacks for cluster itself and the initial nodegroup")
		logger.Info("if you encounter any issues, check CloudFormation console or try 'eksctl utils describe-stacks --region=%s --name=%s'", cfg.Region, cfg.ClusterName)
		errs := stackManager.CreateClusterWithNodeGroups()
		// read any errors (it only gets non-nil errors)
		if len(errs) > 0 {
			logger.Info("%d error(s) occurred and cluster hasn't been created properly, you may wish to check CloudFormation console", len(errs))
			logger.Info("to cleanup resources, run 'eksctl delete cluster --region=%s --name=%s'", cfg.Region, cfg.ClusterName)
			for _, err := range errs {
				logger.Critical("%s\n", err.Error())
			}
			return fmt.Errorf("failed to create cluster %q", cfg.ClusterName)
		}
	}

	logger.Success("all EKS cluster resource for %q had been created", cfg.ClusterName)

	// obtain cluster credentials, write kubeconfig

	{ // post-creation action
		clientConfigBase, err := ctl.NewClientConfig()
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

		if err = ctl.WaitForControlPlane(clientSet); err != nil {
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

	logger.Success("EKS cluster %q in %q region is ready", cfg.ClusterName, cfg.Region)

	return nil
}
