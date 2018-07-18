package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/kubicorn/kubicorn/pkg/logger"

	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/utils"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
)

func createCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create resource(s)",
		Run: func(c *cobra.Command, _ []string) {
			c.Help()
		},
	}

	cmd.AddCommand(createClusterCmd())

	return cmd
}

const (
	EKS_REGION_US_WEST_2   = "us-west-2"
	EKS_REGION_US_EAST_1   = "us-east-1"
	DEFAULT_EKS_REGION     = EKS_REGION_US_WEST_2
	DEFAULT_NODE_COUNT     = 2
	DEFAULT_NODE_TYPE      = "m5.large"
	DEFAULT_SSH_PUBLIC_KEY = "~/.ssh/id_rsa.pub"
)

var (
	writeKubeconfig    bool
	kubeconfigPath     string
	autoKubeconfigPath bool
	setContext         bool
	availabilityZones  []string
)

func createClusterCmd() *cobra.Command {
	cfg := &eks.ClusterConfig{}

	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Create a cluster",
		Run: func(_ *cobra.Command, args []string) {
			if err := doCreateCluster(cfg, getNameArg(args)); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	fs := cmd.Flags()

	exampleClusterName := utils.ClusterName("", "")

	fs.StringVarP(&cfg.ClusterName, "name", "n", "", fmt.Sprintf("EKS cluster name (generated if unspecified, e.g. %q)", exampleClusterName))

	fs.StringVarP(&cfg.Region, "region", "r", DEFAULT_EKS_REGION, "AWS region")
	fs.StringVarP(&cfg.Profile, "profile", "p", "", "AWS creditials profile to use (overrides the AWS_PROFILE environment variable)")

	fs.StringVarP(&cfg.NodeType, "node-type", "t", DEFAULT_NODE_TYPE, "node instance type")
	fs.IntVarP(&cfg.Nodes, "nodes", "N", DEFAULT_NODE_COUNT, "total number of nodes (for a static ASG)")

	// TODO: https://github.com/weaveworks/eksctl/issues/28
	fs.IntVarP(&cfg.MinNodes, "nodes-min", "m", 0, "minimum nodes in ASG")
	fs.IntVarP(&cfg.MaxNodes, "nodes-max", "M", 0, "maximum nodes in ASG")

	fs.StringSliceVar(&availabilityZones, "zones", nil, "(auto-select if unspecified)")

	fs.StringVar(&cfg.SSHPublicKeyPath, "ssh-public-key", DEFAULT_SSH_PUBLIC_KEY, "SSH public key to use for nodes (import from local path, or use existing EC2 key pair)")

	fs.BoolVar(&writeKubeconfig, "write-kubeconfig", true, "toggle writing of kubeconfig")
	fs.BoolVar(&autoKubeconfigPath, "auto-kubeconfig", false, fmt.Sprintf("save kubconfig file by cluster name, e.g. %q", kubeconfig.AutoPath(exampleClusterName)))
	fs.StringVar(&kubeconfigPath, "kubeconfig", kubeconfig.DefaultPath, "path to write kubeconfig (incompatible with --auto-kubeconfig)")
	fs.BoolVar(&setContext, "set-kubeconfig-context", true, "if true then current-context will be set in kubeconfig; if a context is already set then it will be overwritten")

	fs.DurationVar(&cfg.WaitTimeout, "aws-api-timeout", eks.DefaultWaitTimeout, "")
	fs.MarkHidden("aws-api-timeout") // TODO deprecate in 0.2.0
	fs.DurationVar(&cfg.WaitTimeout, "timeout", eks.DefaultWaitTimeout, "max wait time in any polling operations")

	fs.BoolVar(&cfg.Addons.WithIAM.PolicyAmazonEC2ContainerRegistryPowerUser, "full-ecr-access", false, "enable full access to ECR")

	fs.BoolVar(&cfg.CtrlPlaneCfn, "ctrl-plane-cloudformation", true, "Create the EKS control plane using a cloudformation resource.  Disabling this flag falls back to the original behavior.")

	return cmd
}

func doCreateCluster(cfg *eks.ClusterConfig, name string) error {
	ctl := eks.New(cfg)

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

	if cfg.SSHPublicKeyPath == "" {
		return fmt.Errorf("--ssh-public-key must be non-empty string")
	}

	if cfg.Region != EKS_REGION_US_WEST_2 && cfg.Region != EKS_REGION_US_EAST_1 {
		return fmt.Errorf("--region=%s is not supported only %s and %s are supported", cfg.Region, EKS_REGION_US_WEST_2, EKS_REGION_US_EAST_1)
	}

	if err := ctl.SetAvailabilityZones(availabilityZones); err != nil {
		return err
	}

	if err := ctl.LoadSSHPublicKey(); err != nil {
		return err
	}

	logger.Debug("cfg = %#v", cfg)

	logger.Info("creating EKS cluster %q in %q region", cfg.ClusterName, cfg.Region)

	{ // core action
		taskErr := make(chan error)
		// create each of the core cloudformation stacks
		go ctl.CreateCluster(taskErr)
		// read any errors (it only gets non-nil errors)
		for err := range taskErr {
			logger.Info("an error has occurred and cluster hasn't beend created properly")
			logger.Info("to cleanup resources, run 'eksctl delete cluster --region=%s --name=%s'", cfg.Region, cfg.ClusterName)
			return err
		}
	}

	logger.Success("all EKS cluster %q resources has been created", cfg.ClusterName)

	// obtain cluster credentials, write kubeconfig

	{ // post-creation action
		clientConfigBase, err := ctl.NewClientConfig()
		if err != nil {
			return err
		}

		if writeKubeconfig {
			config := clientConfigBase.WithExecHeptioAuthenticator()
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

		if err := cfg.WaitForControlPlane(clientSet); err != nil {
			return err
		}

		// authorise nodes to join
		if err := cfg.CreateDefaultNodeGroupAuthConfigMap(clientSet); err != nil {
			return err
		}

		// wait for nodes to join
		if err := cfg.WaitForNodes(clientSet); err != nil {
			return err
		}

		// check kubectl version, and offer install instructions if missing or old
		// also check heptio-authenticator
		// TODO: https://github.com/weaveworks/eksctl/issues/30
		env, err := ctl.GetCredentialsEnv()
		if err != nil {
			return err
		}
		if err := utils.CheckAllCommands(kubeconfigPath, setContext, clientConfigBase.ContextName, env); err != nil {
			logger.Critical(err.Error())
			logger.Info("cluster should be functional despite missing (or misconfigured) client binaries")
		}
	}

	logger.Success("EKS cluster %q in %q region is ready", cfg.ClusterName, cfg.Region)

	return nil
}
