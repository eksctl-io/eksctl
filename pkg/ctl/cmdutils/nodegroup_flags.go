package cmdutils

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// AddCommonCreateNodeGroupFlags adds common flags for creating a nodegroup
func AddCommonCreateNodeGroupFlags(fs *pflag.FlagSet, cmd *Cmd, ng *api.NodeGroup, mngOptions *CreateManagedNGOptions) {
	fs.StringVarP(&ng.InstanceType, "node-type", "t", "", "node instance type")

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
	ng.SSH.EnableSSM = fs.Bool("enable-ssm", false, "Enable AWS Systems Manager (SSM)")

	fs.StringVar(&ng.AMI, "node-ami", "", "'auto-ssm', 'auto', 'static' (deprecated) or an AMI id (advanced use)")
	fs.StringVar(&ng.AMIFamily, "node-ami-family", api.DefaultNodeImageFamily, "'AmazonLinux2' for the Amazon EKS optimized AMI, or use 'Ubuntu2004' or 'Ubuntu1804' for the official Canonical EKS AMIs")

	fs.BoolVarP(&ng.PrivateNetworking, "node-private-networking", "P", false, "whether to make nodegroup networking private")

	fs.StringSliceVar(&ng.SecurityGroups.AttachIDs, "node-security-groups", []string{}, "attach additional security groups to nodes")

	AddStringToStringVarPFlag(fs, &ng.Labels, "node-labels", "", nil, "extra labels to add when registering the nodes in the nodegroup")
	fs.StringSliceVar(&ng.AvailabilityZones, "node-zones", nil, "(inherited from the cluster if unspecified)")

	fs.StringVar(&ng.InstancePrefix, "instance-prefix", "", "add a prefix value in front of the instance's name")
	fs.StringVar(&ng.InstanceName, "instance-name", "", "overrides the default instance's name")

	fs.BoolVar(ng.DisablePodIMDS, "disable-pod-imds", false, "Blocks IMDS requests from non host networking pods")

	fs.BoolVarP(&mngOptions.Managed, "managed", "", false, "Create EKS-managed nodegroup")
	fs.BoolVar(&mngOptions.Spot, "spot", false, "Create a spot nodegroup (managed nodegroups only)")
	fs.StringSliceVar(&mngOptions.InstanceTypes, "instance-types", nil, "Comma-separated list of instance types (e.g., --instance-types=c3.large,c4.large,c5.large")
}

func incompatibleManagedNodesFlags() []string {
	return []string{
		"max-pods-per-node",
		"node-ami",
		"node-security-groups",
	}
}

// AddCommonCreateNodeGroupIAMAddonsFlags adds flags to set ng.IAM.WithAddonPolicies
func AddCommonCreateNodeGroupIAMAddonsFlags(fs *pflag.FlagSet, ng *api.NodeGroup) {
	ng.IAM.WithAddonPolicies.AutoScaler = new(bool)
	ng.IAM.WithAddonPolicies.ExternalDNS = new(bool)
	ng.IAM.WithAddonPolicies.ImageBuilder = new(bool)
	ng.IAM.WithAddonPolicies.AppMesh = new(bool)
	ng.IAM.WithAddonPolicies.AppMeshPreview = new(bool)
	ng.IAM.WithAddonPolicies.AWSLoadBalancerController = new(bool)
	ng.IAM.WithAddonPolicies.XRay = new(bool)
	ng.IAM.WithAddonPolicies.CloudWatch = new(bool)
	fs.BoolVar(ng.IAM.WithAddonPolicies.AutoScaler, "asg-access", false, "enable IAM policy for cluster-autoscaler")
	fs.BoolVar(ng.IAM.WithAddonPolicies.ExternalDNS, "external-dns-access", false, "enable IAM policy for external-dns")
	fs.BoolVar(ng.IAM.WithAddonPolicies.ImageBuilder, "full-ecr-access", false, "enable full access to ECR")
	fs.BoolVar(ng.IAM.WithAddonPolicies.AppMesh, "appmesh-access", false, "enable full access to AppMesh")
	fs.BoolVar(ng.IAM.WithAddonPolicies.AppMeshPreview, "appmesh-preview-access", false, "enable full access to AppMesh Preview")
	fs.BoolVar(ng.IAM.WithAddonPolicies.AWSLoadBalancerController, "alb-ingress-access", false, "enable full access for alb-ingress-controller")
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
