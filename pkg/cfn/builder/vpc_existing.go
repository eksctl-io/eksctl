package builder

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	awsec2 "github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/vpc"
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
	out, err := v.ec2API.DescribeVpcs(&awsec2.DescribeVpcsInput{
		VpcIds: aws.StringSlice([]string{v.clusterConfig.VPC.ID}),
	})

	if err != nil {
		return nil, nil, fmt.Errorf("failed to describe VPC %q: %w", v.clusterConfig.VPC.ID, err)
	}

	if len(out.Vpcs) == 0 {
		return nil, nil, fmt.Errorf("VPC %q does not exist", v.clusterConfig.VPC.ID)
	}

	if v.clusterConfig.IPv6Enabled() {
		if err := v.checkIPv6CidrBlockAssociated(out); err != nil {
			return nil, nil, err
		}
	}
	if err := v.importExistingResources(); err != nil {
		return nil, nil, errors.Wrap(err, "error importing VPC resources")
	}

	v.addOutputs()
	return v.vpcID, v.subnetDetails, nil
}

// addOutputs adds VPC resource outputs
func (v *ExistingVPCResourceSet) addOutputs() {
	v.rs.defineOutput(outputs.ClusterVPC, v.vpcID, true, func(val string) error {
		v.clusterConfig.VPC.ID = val
		return nil
	})

	if v.clusterConfig.VPC.NAT != nil {
		v.rs.defineOutputWithoutCollector(outputs.ClusterFeatureNATMode, v.clusterConfig.VPC.NAT.Gateway, false)
	}

	addSubnetOutput := func(subnetRefs []*gfnt.Value, topology api.SubnetTopology, outputName string) {
		v.rs.defineJoinedOutput(outputName, subnetRefs, true, func(value string) error {
			return vpc.ImportSubnetsFromIDList(v.ec2API, v.clusterConfig, topology, strings.Split(value, ","))
		})
	}

	if subnetAZs := v.subnetDetails.PrivateSubnetRefs(); len(subnetAZs) > 0 {
		addSubnetOutput(subnetAZs, api.SubnetTopologyPrivate, outputs.ClusterSubnetsPrivate)
	}

	if subnetAZs := v.subnetDetails.PublicSubnetRefs(); len(subnetAZs) > 0 {
		addSubnetOutput(subnetAZs, api.SubnetTopologyPublic, outputs.ClusterSubnetsPublic)
	}

	if v.isFullyPrivate() {
		v.rs.defineOutputWithoutCollector(outputs.ClusterFullyPrivate, true, true)
	}
}

func (v *ExistingVPCResourceSet) checkIPv6CidrBlockAssociated(describeVPCOutput *awsec2.DescribeVpcsOutput) error {
	if len(describeVPCOutput.Vpcs[0].Ipv6CidrBlockAssociationSet) == 0 {
		return fmt.Errorf("VPC %q does not have any associated IPv6 CIDR blocks", v.clusterConfig.VPC.ID)
	}
	return nil
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
	var subnetResources []SubnetResource
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
		subnetResources = append(subnetResources, sr)
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

// RenderJSON returns the rendered JSON
func (v *ExistingVPCResourceSet) RenderJSON() ([]byte, error) {
	return v.rs.renderJSON()
}
