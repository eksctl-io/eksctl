package eks

import (
	"fmt"

	"github.com/kubicorn/kubicorn/pkg/logger"
	"k8s.io/kops/pkg/util/subnet"

	"github.com/weaveworks/eksctl/pkg/eks/api"
)

// SetSubnets defines CIDRs for each of the subnets,
// it must be called after SetAvailabilityZones
func (c *ClusterProvider) SetSubnets() error {
	var err error

	c.Spec.VPC.Subnets = map[api.SubnetTopology]map[string]api.Network{
		api.SubnetTopologyPublic:  map[string]api.Network{},
		api.SubnetTopologyPrivate: map[string]api.Network{},
	}

	zoneCIDRs, err := subnet.SplitInto8(c.Spec.VPC.CIDR)
	if err != nil {
		return err
	}

	logger.Debug("VPC CIDR (%s) was divided into 8 subnets %v", c.Spec.VPC.CIDR.String(), zoneCIDRs)

	zonesTotal := len(c.Spec.AvailabilityZones)
	if 2*zonesTotal > len(zoneCIDRs) {
		return fmt.Errorf("insuffience number of subnets (have %d, but need %d) for %d availability zones", len(zoneCIDRs), 2*zonesTotal, zonesTotal)
	}

	for i, zone := range c.Spec.AvailabilityZones {
		public := zoneCIDRs[i]
		private := zoneCIDRs[i+zonesTotal]
		c.Spec.VPC.Subnets[api.SubnetTopologyPublic][zone] = api.Network{
			CIDR: public,
		}
		c.Spec.VPC.Subnets[api.SubnetTopologyPrivate][zone] = api.Network{
			CIDR: private,
		}
		logger.Info("subnets for %s - public:%s private:%s", zone, public.String(), private.String())
	}

	return nil
}
