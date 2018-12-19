package cmdutils

import (
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/ami"
	"github.com/weaveworks/eksctl/pkg/eks/api"
)

const (
	defaultNodeType     = "m5.large"
	defaultSSHPublicKey = "~/.ssh/id_rsa.pub"
)

func AddCommonCreateNodeGroupFlags(fs *pflag.FlagSet, p *api.ProviderConfig, cfg *api.ClusterConfig, ng *api.NodeGroup) {
	fs.StringVarP(&ng.InstanceType, "node-type", "t", defaultNodeType, "node instance type")
	fs.IntVarP(&ng.DesiredCapacity, "nodes", "N", api.DefaultNodeCount, "total number of nodes (for a static ASG)")

	// TODO: https://github.com/weaveworks/eksctl/issues/28
	fs.IntVarP(&ng.MinSize, "nodes-min", "m", 0, "minimum nodes in ASG")
	fs.IntVarP(&ng.MaxSize, "nodes-max", "M", 0, "maximum nodes in ASG")

	fs.IntVarP(&ng.VolumeSize, "node-volume-size", "", 0, "Node volume size (in GB)")
	fs.IntVar(&ng.MaxPodsPerNode, "max-pods-per-node", 0, "maximum number of pods per node (set automatically if unspecified)")

	fs.BoolVar(&ng.AllowSSH, "ssh-access", false, "control SSH access for nodes")
	fs.StringVar(&ng.SSHPublicKeyPath, "ssh-public-key", defaultSSHPublicKey, "SSH public key to use for nodes (import from local path, or use existing EC2 key pair)")

	fs.StringVar(&ng.AMI, "node-ami", ami.ResolverStatic, "Advanced use cases only. If 'static' is supplied (default) then eksctl will use static AMIs; if 'auto' is supplied then eksctl will automatically set the AMI based on version/region/instance type; if any other value is supplied it will override the AMI to use for the nodes. Use with extreme care.")
	fs.StringVar(&ng.AMIFamily, "node-ami-family", ami.ImageFamilyAmazonLinux2, "Advanced use cases only. If 'AmazonLinux2' is supplied (default), then eksctl will use the offical AWS EKS AMIs (Amazon Linux 2); if 'Ubuntu1804' is supplied, then eksctl will use the offical Canonical EKS AMIs (Ubuntu 18.04).")

	fs.BoolVarP(&ng.PrivateNetworking, "node-private-networking", "P", false, "whether to make nodegroup networking private")

	fs.Var(&ng.Labels, "node-labels", `labels to set for the nodegroup, e.g. "partition=backend,nodeclass=hugememory"`)
	fs.StringSliceVar(&ng.AvailabilityZones, "node-zones", nil, "(iherited from the cluster if unspecified)")

}
