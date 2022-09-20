package kops

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"k8s.io/kops/upup/pkg/fi/cloudup/awsup"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/vpc"
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

func (k *Wrapper) topologyOf(s ec2types.Subnet) api.SubnetTopology {
	for _, t := range s.Tags {
		if *t.Key == "SubnetType" && *t.Value == "Private" {
			return api.SubnetTopologyPrivate
		}
	}
	return api.SubnetTopologyPublic // "Utility", "Public" or unspecified
}

// UseVPC finds VPC and subnets that give kops cluster uses and add those to EKS cluster config
func (k *Wrapper) UseVPC(ctx context.Context, ec2API awsapi.EC2, spec *api.ClusterConfig) error {
	spec.VPC.CIDR = nil // ensure to reset the CIDR

	clusterTag := fmt.Sprintf("kubernetes.io/cluster/%s", k.clusterName)
	output, err := ec2API.DescribeSubnets(ctx, &ec2.DescribeSubnetsInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String(fmt.Sprintf("tag:%s", clusterTag)),
				Values: []string{"owned"},
			},
		},
	})
	if err != nil {
		return err
	}

	var (
		publicSubnets  []ec2types.Subnet
		privateSubnets []ec2types.Subnet
	)

	for _, subnet := range output.Subnets {
		switch k.topologyOf(subnet) {
		case api.SubnetTopologyPublic:
			publicSubnets = append(publicSubnets, subnet)
		case api.SubnetTopologyPrivate:
			privateSubnets = append(privateSubnets, subnet)
		}
	}

	if err := vpc.ImportSubnets(ctx, ec2API, spec, spec.VPC.Subnets.Private, privateSubnets, nil); err != nil {
		return err
	}

	if err := vpc.ImportSubnets(ctx, ec2API, spec, spec.VPC.Subnets.Public, publicSubnets, nil); err != nil {
		return err
	}

	logger.Debug("subnets = %#v", spec.VPC.Subnets)
	if err := spec.HasSufficientSubnets(); err != nil {
		return errors.Wrapf(err, "using VPC from kops cluster %q", k.clusterName)
	}
	return nil
}
