package cmdutils

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/weaveworks/eksctl/pkg/ami"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
)

const (
	defaultNodeType     = "m5.large"
	defaultSSHPublicKey = "~/.ssh/id_rsa.pub"
)

// AddCommonCreateNodeGroupFlags adds common flags for creating a node group
func AddCommonCreateNodeGroupFlags(cmd *cobra.Command, fs *pflag.FlagSet, p *api.ProviderConfig, cfg *api.ClusterConfig, ng *api.NodeGroup) {
	var desiredCapacity int
	var minSize int
	var maxSize int

	fs.StringVarP(&ng.InstanceType, "node-type", "t", defaultNodeType, "node instance type")
	fs.IntVarP(&desiredCapacity, "nodes", "N", api.DefaultNodeCount, "total number of nodes (for a static ASG)")

	// TODO: https://github.com/weaveworks/eksctl/issues/28
	fs.IntVarP(&minSize, "nodes-min", "m", api.DefaultNodeCount, "minimum nodes in ASG")
	fs.IntVarP(&maxSize, "nodes-max", "M", api.DefaultNodeCount, "maximum nodes in ASG")

	cmd.PreRun = func(cmd *cobra.Command, args[] string) {
		if f := cmd.Flag("nodes"); f.Changed {
			ng.DesiredCapacity = &desiredCapacity
		}
		if f := cmd.Flag("nodes-min"); f.Changed {
			ng.MinSize= &minSize
		}
		if f := cmd.Flag("nodes-max"); f.Changed {
			ng.MaxSize= &maxSize
		}
	}

	fs.IntVar(&ng.VolumeSize, "node-volume-size", ng.VolumeSize, "node volume size in GB")
	fs.StringVar(&ng.VolumeType, "node-volume-type", ng.VolumeType, fmt.Sprintf("node volume type (valid options: %s)", strings.Join(api.SupportedNodeVolumeTypes(), ", ")))

	fs.IntVar(&ng.MaxPodsPerNode, "max-pods-per-node", 0, "maximum number of pods per node (set automatically if unspecified)")

	fs.BoolVar(&ng.AllowSSH, "ssh-access", false, "control SSH access for nodes")
	fs.StringVar(&ng.SSHPublicKeyPath, "ssh-public-key", defaultSSHPublicKey, "SSH public key to use for nodes (import from local path, or use existing EC2 key pair)")

	fs.StringVar(&ng.AMI, "node-ami", ami.ResolverStatic, "Advanced use cases only. If 'static' is supplied (default) then eksctl will use static AMIs; if 'auto' is supplied then eksctl will automatically set the AMI based on version/region/instance type; if any other value is supplied it will override the AMI to use for the nodes. Use with extreme care.")
	fs.StringVar(&ng.AMIFamily, "node-ami-family", ami.ImageFamilyAmazonLinux2, "Advanced use cases only. If 'AmazonLinux2' is supplied (default), then eksctl will use the offical AWS EKS AMIs (Amazon Linux 2); if 'Ubuntu1804' is supplied, then eksctl will use the offical Canonical EKS AMIs (Ubuntu 18.04).")

	fs.BoolVarP(&ng.PrivateNetworking, "node-private-networking", "P", false, "whether to make nodegroup networking private")

	fs.StringSliceVar(&ng.SecurityGroups.AttachIDs, "node-security-groups", []string{}, "Attach additional security groups to nodes, so that it can be used to allow extra ingress/egress access from/to pods")

	fs.Var(&ng.Labels, "node-labels", `Extra labels to add when registering the nodes in the nodegroup, e.g. "partition=backend,nodeclass=hugememory"`)
	fs.StringSliceVar(&ng.AvailabilityZones, "node-zones", nil, "(inherited from the cluster if unspecified)")
}

// AddCommonCreateNodeGroupIAMAddonsFlags adds flags to set ng.IAM.WithAddonPolicies
func AddCommonCreateNodeGroupIAMAddonsFlags(fs *pflag.FlagSet, ng *api.NodeGroup) {
	fs.StringSliceVar(&ng.IAM.AttachPolicyARNs, "temp-node-role-policies", []string{}, "Advanced use cases only. "+
		"All the IAM policies to be associated to the node's instance role. "+
		"Beware that you MUST include the policies for EKS and CNI related AWS API Access, like `arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy` and `arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy` that are used by default when this flag is omitted.")
	fs.MarkHidden("temp-node-role-policies")
	fs.StringVar(&ng.IAM.InstanceRoleName, "temp-node-role-name", "", "Advanced use cases only. Specify the exact name of the node's instance role for easier integration with K8S-IAM integrations like kube2iam. See https://github.com/weaveworks/eksctl/issues/398 for more information.")
	fs.MarkHidden("temp-node-role-name")

	ng.IAM.WithAddonPolicies.AutoScaler = new(bool)
	ng.IAM.WithAddonPolicies.ExternalDNS = new(bool)
	ng.IAM.WithAddonPolicies.ImageBuilder = new(bool)
	fs.BoolVar(ng.IAM.WithAddonPolicies.AutoScaler, "asg-access", false, "enable IAM policy for cluster-autoscaler")
	fs.BoolVar(ng.IAM.WithAddonPolicies.ExternalDNS, "external-dns-access", false, "enable IAM policy for external-dns")
	fs.BoolVar(ng.IAM.WithAddonPolicies.ImageBuilder, "full-ecr-access", false, "enable full access to ECR")
}
