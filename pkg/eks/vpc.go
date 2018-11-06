package eks

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/weaveworks/eksctl/pkg/eks/api"
	"k8s.io/kops/pkg/util/subnet"
)

// SetSubnets defines CIDRs for each of the subnets,
// it must be called after SetAvailabilityZones
func (c *ClusterProvider) SetSubnets() error {
	var err error

	vpc := c.Spec.VPC
	vpc.Subnets = map[api.SubnetTopology]map[string]api.Network{
		api.SubnetTopologyPublic:  map[string]api.Network{},
		api.SubnetTopologyPrivate: map[string]api.Network{},
	}
	prefix, _ := c.Spec.VPC.CIDR.Mask.Size()
	if (prefix < 16) || (prefix > 24) {
		return fmt.Errorf("VPC CIDR prefix must be betwee /16 and /24")
	}
	zoneCIDRs, err := subnet.SplitInto8(c.Spec.VPC.CIDR)
	if err != nil {
		return err
	}

	logger.Debug("VPC CIDR (%s) was divided into 8 subnets %v", vpc.CIDR.String(), zoneCIDRs)

	zonesTotal := len(c.Spec.AvailabilityZones)
	if 2*zonesTotal > len(zoneCIDRs) {
		return fmt.Errorf("insufficient number of subnets (have %d, but need %d) for %d availability zones", len(zoneCIDRs), 2*zonesTotal, zonesTotal)
	}

	for i, zone := range c.Spec.AvailabilityZones {
		public := zoneCIDRs[i]
		private := zoneCIDRs[i+zonesTotal]
		vpc.Subnets[api.SubnetTopologyPublic][zone] = api.Network{
			CIDR: public,
		}
		vpc.Subnets[api.SubnetTopologyPrivate][zone] = api.Network{
			CIDR: private,
		}
		logger.Info("subnets for %s - public:%s private:%s", zone, public.String(), private.String())
	}

	return nil
}

// UseSubnets imports
func (c *ClusterProvider) UseSubnets(topology api.SubnetTopology, subnetIDs []string) error {
	if len(subnetIDs) == 0 {
		return nil
	}
	input := &ec2.DescribeSubnetsInput{
		SubnetIds: aws.StringSlice(subnetIDs),
	}
	output, err := c.Provider.EC2().DescribeSubnets(input)
	if err != nil {
		return err
	}

	for _, subnet := range output.Subnets {
		if c.Spec.VPC.ID == "" {
			c.Spec.VPC.ID = *subnet.VpcId
		} else if c.Spec.VPC.ID != *subnet.VpcId {
			return fmt.Errorf("given subnets (%s) are not in the same VPC", strings.Join(subnetIDs, ", "))
		}

		c.Spec.ImportSubnet(topology, *subnet.AvailabilityZone, *subnet.SubnetId)
		c.Spec.AppendAvailabilityZone(*subnet.AvailabilityZone)
	}

	return nil
}
