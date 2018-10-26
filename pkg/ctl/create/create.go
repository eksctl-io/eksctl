package create

import (
	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/ami"
	"github.com/weaveworks/eksctl/pkg/eks/api"
)

const (
	defaultNodeType     = "m5.large"
	defaultSSHPublicKey = "~/.ssh/id_rsa.pub"
)

// Command will create the `create` commands
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create resource(s)",
		Run: func(c *cobra.Command, _ []string) {
			if err := c.Help(); err != nil {
				logger.Debug("ignoring error %q", err.Error())
			}
		},
	}

	cmd.AddCommand(createClusterCmd())
	cmd.AddCommand(createNodeGroupCmd())

	return cmd
}

func addCommonCreateFlags(fs *pflag.FlagSet, cfg *api.ClusterConfig, ng *api.NodeGroup) {
	fs.StringVarP(&cfg.Region, "region", "r", "", "AWS region")
	fs.StringVarP(&cfg.Profile, "profile", "p", "", "AWS credentials profile to use (overrides the AWS_PROFILE environment variable)")
	fs.StringToStringVarP(&cfg.Tags, "tags", "", map[string]string{}, `A list of KV pairs used to tag the AWS resources (e.g. "Owner=John Doe,Team=Some Team")`)

	fs.StringVarP(&ng.InstanceType, "node-type", "t", defaultNodeType, "node instance type")
	fs.IntVarP(&ng.DesiredCapacity, "nodes", "N", api.DefaultNodeCount, "total number of nodes (for a static ASG)")

	// TODO: https://github.com/weaveworks/eksctl/issues/28
	fs.IntVarP(&ng.MinSize, "nodes-min", "m", 0, "minimum nodes in ASG")
	fs.IntVarP(&ng.MaxSize, "nodes-max", "M", 0, "maximum nodes in ASG")

	fs.StringVar(&ng.AMI, "node-ami", ami.ResolverStatic, "Advanced use cases only. If 'static' is supplied (default) then eksctl will use static AMIs; if 'auto' is supplied then eksctl will automatically set the AMI based on region/instance type; if any other value is supplied it will override the AMI to use for the nodes. Use with extreme care.")

	fs.IntVarP(&ng.VolumeSize, "node-volume-size", "", 0, "Node volume size (in GB)")
	fs.IntVar(&ng.MaxPodsPerNode, "max-pods-per-node", 0, "maximum number of pods per node (set automatically if unspecified)")

	fs.BoolVar(&ng.AllowSSH, "ssh-access", false, "control SSH access for nodes")
	fs.StringVar(&ng.SSHPublicKeyPath, "ssh-public-key", defaultSSHPublicKey, "SSH public key to use for nodes (import from local path, or use existing EC2 key pair)")

	fs.DurationVar(&cfg.WaitTimeout, "aws-api-timeout", api.DefaultWaitTimeout, "")
	// TODO deprecate in 0.2.0
	if err := fs.MarkHidden("aws-api-timeout"); err != nil {
		logger.Debug("ignoring error %q", err.Error())
	}
	fs.DurationVar(&cfg.WaitTimeout, "timeout", api.DefaultWaitTimeout, "max wait time in any polling operations")

}
