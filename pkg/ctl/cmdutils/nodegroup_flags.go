package cmdutils

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/weaveworks/eksctl/pkg/ami"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// AddCommonCreateNodeGroupFlags adds common flags for creating a nodegroup
func AddCommonCreateNodeGroupFlags(fs *pflag.FlagSet, cmd *Cmd, ng *api.NodeGroup) {
	fs.StringVarP(&ng.InstanceType, "node-type", "t", api.DefaultNodeType, "node instance type")

	desiredCapacity := fs.IntP("nodes", "N", api.DefaultNodeCount, "total number of nodes (for a static ASG)")
	minSize := fs.IntP("nodes-min", "m", api.DefaultNodeCount, "minimum nodes in ASG")
	maxSize := fs.IntP("nodes-max", "M", api.DefaultNodeCount, "maximum nodes in ASG")

	AddPreRun(cmd.CobraCommand, func(cobraCmd *cobra.Command, args []string) {
		if f := cobraCmd.Flag("nodes"); f.Changed {
			ng.DesiredCapacity = desiredCapacity
		}
		if f := cobraCmd.Flag("nodes-min"); f.Changed {
			ng.MinSize = minSize
		}
		if f := cobraCmd.Flag("nodes-max"); f.Changed {
			ng.MaxSize = maxSize
		}
	})

	fs.IntVar(ng.VolumeSize, "node-volume-size", *ng.VolumeSize, "node volume size in GB")
	fs.StringVar(ng.VolumeType, "node-volume-type", *ng.VolumeType, fmt.Sprintf("node volume type (valid options: %s)", strings.Join(api.SupportedNodeVolumeTypes(), ", ")))

	fs.IntVar(&ng.MaxPodsPerNode, "max-pods-per-node", 0, "maximum number of pods per node (set automatically if unspecified)")

	ng.SSH.Allow = fs.Bool("ssh-access", *ng.SSH.Allow, "control SSH access for nodes. Uses ~/.ssh/id_rsa.pub as default key path if enabled")
	ng.SSH.PublicKeyPath = fs.String("ssh-public-key", "", "SSH public key to use for nodes (import from local path, or use existing EC2 key pair)")

	fs.StringVar(&ng.AMI, "node-ami", ami.ResolverStatic, "Advanced use cases only. If 'static' is supplied (default) then eksctl will use static AMIs; if 'auto' is supplied then eksctl will automatically set the AMI based on version/region/instance type; if any other value is supplied it will override the AMI to use for the nodes. Use with extreme care.")
	fs.StringVar(&ng.AMIFamily, "node-ami-family", api.DefaultNodeImageFamily, "Advanced use cases only. If 'AmazonLinux2' is supplied (default), then eksctl will use the official AWS EKS AMIs (Amazon Linux 2); if 'Ubuntu1804' is supplied, then eksctl will use the official Canonical EKS AMIs (Ubuntu 18.04).")

	fs.BoolVarP(&ng.PrivateNetworking, "node-private-networking", "P", false, "whether to make nodegroup networking private")

	fs.StringSliceVar(&ng.SecurityGroups.AttachIDs, "node-security-groups", []string{}, "Attach additional security groups to nodes, so that it can be used to allow extra ingress/egress access from/to pods")

	fs.StringToStringVar(&ng.Labels, "node-labels", nil, `Extra labels to add when registering the nodes in the nodegroup, e.g. "partition=backend,nodeclass=hugememory"`)
	fs.StringSliceVar(&ng.AvailabilityZones, "node-zones", nil, "(inherited from the cluster if unspecified)")
}

// AddCommonCreateNodeGroupIAMAddonsFlags adds flags to set ng.IAM.WithAddonPolicies
func AddCommonCreateNodeGroupIAMAddonsFlags(fs *pflag.FlagSet, ng *api.NodeGroup) {
	ng.IAM.WithAddonPolicies.AutoScaler = new(bool)
	ng.IAM.WithAddonPolicies.ExternalDNS = new(bool)
	ng.IAM.WithAddonPolicies.ImageBuilder = new(bool)
	ng.IAM.WithAddonPolicies.AppMesh = new(bool)
	ng.IAM.WithAddonPolicies.ALBIngress = new(bool)
	ng.IAM.WithAddonPolicies.XRay = new(bool)
	ng.IAM.WithAddonPolicies.CloudWatch = new(bool)
	fs.BoolVar(ng.IAM.WithAddonPolicies.AutoScaler, "asg-access", false, "enable IAM policy for cluster-autoscaler")
	fs.BoolVar(ng.IAM.WithAddonPolicies.ExternalDNS, "external-dns-access", false, "enable IAM policy for external-dns")
	fs.BoolVar(ng.IAM.WithAddonPolicies.ImageBuilder, "full-ecr-access", false, "enable full access to ECR")
	fs.BoolVar(ng.IAM.WithAddonPolicies.AppMesh, "appmesh-access", false, "enable full access to AppMesh")
	fs.BoolVar(ng.IAM.WithAddonPolicies.ALBIngress, "alb-ingress-access", false, "enable full access for alb-ingress-controller")
}

// AddNodeGroupFilterFlags add common `--include` and `--exclude` flags for filtering nodegroups
func AddNodeGroupFilterFlags(fs *pflag.FlagSet, includeGlobs, excludeGlobs *[]string) {
	fs.StringSliceVar(includeGlobs, "only", nil, "")
	_ = fs.MarkDeprecated("only", "use --include")

	fs.StringSliceVar(includeGlobs, "include", nil,
		"nodegroups to include (list of globs), e.g.: 'ng-team-?,prod-*'")

	fs.StringSliceVar(excludeGlobs, "exclude", nil,
		"nodegroups to exclude (list of globs), e.g.: 'ng-team-?,prod-*'")
}
