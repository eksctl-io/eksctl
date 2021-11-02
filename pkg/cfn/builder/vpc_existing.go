package builder

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"
)

type ExistingVPCResourceSet struct {
	rs            *resourceSet
	clusterConfig *api.ClusterConfig
	ec2API        ec2iface.EC2API
	vpcID         *gfnt.Value
	subnetDetails *SubnetDetails
}

// NewExistingVPCResourceSet creates and returns a new VPCResourceSet
func NewExistingVPCResourceSet(rs *resourceSet, clusterConfig *api.ClusterConfig, ec2API ec2iface.EC2API) *ExistingVPCResourceSet {
	return &ExistingVPCResourceSet{
		rs:            rs,
		clusterConfig: clusterConfig,
		ec2API:        ec2API,
		vpcID:         gfnt.NewString(clusterConfig.VPC.ID),
		subnetDetails: &SubnetDetails{},
	}
}

func (v *ExistingVPCResourceSet) CreateTemplate() (*gfnt.Value, *SubnetDetails, error) {
	if err := v.importExistingResources(); err != nil {
		return nil, nil, errors.Wrap(err, "error importing VPC resources")
	}
	return v.vpcID, v.subnetDetails, nil
}

func (v *ExistingVPCResourceSet) importExistingResources() error {
	if subnets := v.clusterConfig.VPC.Subnets.Private; subnets != nil {
		var (
			subnetRoutes map[string]string
			err          error
		)
		if v.isFullyPrivate() {
			subnetRoutes, err = importRouteTables(v.ec2API, v.clusterConfig.VPC.Subnets.Private)
			if err != nil {
				return err
			}
		}

		subnetResources, err := makeSubnetResources(subnets, subnetRoutes)
		if err != nil {
			return err
		}
		v.subnetDetails.Private = subnetResources
	}

	if subnets := v.clusterConfig.VPC.Subnets.Public; subnets != nil {
		subnetResources, err := makeSubnetResources(subnets, nil)
		if err != nil {
			return err
		}
		v.subnetDetails.Public = subnetResources
	}

	return nil
}

func makeSubnetResources(subnets map[string]api.AZSubnetSpec, subnetRoutes map[string]string) ([]SubnetResource, error) {
	subnetResources := make([]SubnetResource, len(subnets))
	i := 0
	for _, network := range subnets {
		az := network.AZ
		sr := SubnetResource{
			AvailabilityZone: az,
			Subnet:           gfnt.NewString(network.ID),
		}

		if subnetRoutes != nil {
			rt, ok := subnetRoutes[network.ID]
			if !ok {
				return nil, errors.Errorf("failed to find an explicit route table associated with subnet %q; "+
					"eksctl does not modify the main route table if a subnet is not associated with an explicit route table", network.ID)
			}
			sr.RouteTable = gfnt.NewString(rt)
		}
		subnetResources[i] = sr
		i++
	}
	return subnetResources, nil
}

func importRouteTables(ec2API ec2iface.EC2API, subnets map[string]api.AZSubnetSpec) (map[string]string, error) {
	var subnetIDs []string
	for _, subnet := range subnets {
		subnetIDs = append(subnetIDs, subnet.ID)
	}

	var routeTables []*ec2.RouteTable
	var nextToken *string

	for {
		output, err := ec2API.DescribeRouteTables(&ec2.DescribeRouteTablesInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("association.subnet-id"),
					Values: aws.StringSlice(subnetIDs),
				},
			},
			NextToken: nextToken,
		})

		if err != nil {
			return nil, errors.Wrap(err, "error describing route tables")
		}

		routeTables = append(routeTables, output.RouteTables...)

		if nextToken = output.NextToken; nextToken == nil {
			break
		}
	}

	subnetRoutes := make(map[string]string)
	for _, rt := range routeTables {
		for _, rta := range rt.Associations {
			if rta.Main != nil && *rta.Main {
				return nil, errors.New("subnets must be associated with a non-main route table; eksctl does not modify the main route table")
			}
			subnetRoutes[*rta.SubnetId] = *rt.RouteTableId
		}
	}
	return subnetRoutes, nil
}

func (v *ExistingVPCResourceSet) isFullyPrivate() bool {
	return v.clusterConfig.PrivateCluster.Enabled
}
