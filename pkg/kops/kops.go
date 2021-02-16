package kops

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/vpc"
	"github.com/weaveworks/logger"
	"k8s.io/kops/pkg/resources/aws"
	"k8s.io/kops/upup/pkg/fi/cloudup/awsup"
)

// Wrapper for interacting with a kops cluster
type Wrapper struct {
	clusterName string
	cloud       awsup.AWSCloud
}

// NewWrapper constructs a kops wrapper for a given EKS cluster config and kops cluster
func NewWrapper(region, kopsClusterName string) (*Wrapper, error) {
	cloud, err := awsup.NewAWSCloud(region, nil)
	if err != nil {
		return nil, err
	}

	return &Wrapper{kopsClusterName, cloud}, nil
}

func (k *Wrapper) isOwned(t *ec2.Tag) bool {
	return *t.Key == "kubernetes.io/cluster/"+k.clusterName && *t.Value == "owned"
}

func (k *Wrapper) topologyOf(s *ec2.Subnet) api.SubnetTopology {
	for _, t := range s.Tags {
		if *t.Key == "SubnetType" && *t.Value == "Private" {
			return api.SubnetTopologyPrivate
		}
	}
	return api.SubnetTopologyPublic // "Utility", "Public" or unspecified
}

// UseVPC finds VPC and subnets that give kops cluster uses and add those to EKS cluster config
func (k *Wrapper) UseVPC(ec2API ec2iface.EC2API, spec *api.ClusterConfig) error {
	spec.VPC.CIDR = nil // ensure to reset the CIDR

	allSubnets, err := aws.ListSubnets(k.cloud, k.clusterName)
	if err != nil {
		return err
	}

	subnetsByTopology := map[api.SubnetTopology][]*ec2.Subnet{
		api.SubnetTopologyPrivate: {},
		api.SubnetTopologyPublic:  {},
	}

	for _, subnet := range allSubnets {
		if subnet.Type != ec2.ResourceTypeSubnet {
			continue
		}
		subnet := subnet.Obj.(*ec2.Subnet)
		for _, tag := range subnet.Tags {
			if k.isOwned(tag) {
				t := k.topologyOf(subnet)
				subnetsByTopology[t] = append(subnetsByTopology[t], subnet)
			}
		}
	}
	for t, subnets := range subnetsByTopology {
		if err := vpc.ImportSubnets(ec2API, spec, t, subnets); err != nil {
			return err
		}
	}

	logger.Debug("subnets = %#v", spec.VPC.Subnets)
	if err := spec.HasSufficientSubnets(); err != nil {
		return errors.Wrapf(err, "using VPC from kops cluster %q", k.clusterName)
	}
	return nil
}
