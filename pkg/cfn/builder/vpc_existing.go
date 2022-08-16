package builder

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/pkg/errors"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

type ExistingVPCResourceSet struct {
	rs            *resourceSet
	clusterConfig *api.ClusterConfig
	ec2API        awsapi.EC2
	vpcID         *gfnt.Value
	subnetDetails *SubnetDetails
}

// NewExistingVPCResourceSet creates and returns a new VPCResourceSet
func NewExistingVPCResourceSet(rs *resourceSet, clusterConfig *api.ClusterConfig, ec2API awsapi.EC2) *ExistingVPCResourceSet {
	return &ExistingVPCResourceSet{
		rs:            rs,
		clusterConfig: clusterConfig,
		ec2API:        ec2API,
		vpcID:         gfnt.NewString(clusterConfig.VPC.ID),
		subnetDetails: &SubnetDetails{
			controlPlaneOnOutposts: clusterConfig.IsControlPlaneOnOutposts(),
		},
	}
}

func (v *ExistingVPCResourceSet) CreateTemplate(ctx context.Context) (*gfnt.Value, *SubnetDetails, error) {
	out, err := v.ec2API.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{
		VpcIds: []string{v.clusterConfig.VPC.ID},
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
	if err := v.importExistingResources(ctx); err != nil {
		return nil, nil, errors.Wrap(err, "error importing VPC resources")
	}

	v.addOutputs(ctx)
	return v.vpcID, v.subnetDetails, nil
}

// addOutputs adds VPC resource outputs
func (v *ExistingVPCResourceSet) addOutputs(ctx context.Context) {
	v.rs.defineOutput(outputs.ClusterVPC, v.vpcID, true, func(val string) error {
		v.clusterConfig.VPC.ID = val
		return nil
	})

	if v.clusterConfig.VPC.NAT != nil {
		v.rs.defineOutputWithoutCollector(outputs.ClusterFeatureNATMode, v.clusterConfig.VPC.NAT.Gateway, false)
	}

	addSubnetOutput := func(subnetRefs []*gfnt.Value, subnetMapping api.AZSubnetMapping, outputName string) {
		v.rs.defineJoinedOutput(outputName, subnetRefs, true, func(value string) error {
			return vpc.ImportSubnetsFromIDList(ctx, v.ec2API, v.clusterConfig, subnetMapping, strings.Split(value, ","))
		})
	}

	if subnetAZs := v.subnetDetails.PrivateSubnetRefs(); len(subnetAZs) > 0 {
		addSubnetOutput(subnetAZs, v.clusterConfig.VPC.Subnets.Private, outputs.ClusterSubnetsPrivate)
	}

	if subnetAZs := v.subnetDetails.PublicSubnetRefs(); len(subnetAZs) > 0 {
		addSubnetOutput(subnetAZs, v.clusterConfig.VPC.Subnets.Public, outputs.ClusterSubnetsPublic)
	}

	if v.clusterConfig.IsFullyPrivate() {
		v.rs.defineOutputWithoutCollector(outputs.ClusterFullyPrivate, true, true)
	}
}

func (v *ExistingVPCResourceSet) checkIPv6CidrBlockAssociated(describeVPCOutput *ec2.DescribeVpcsOutput) error {
	if len(describeVPCOutput.Vpcs[0].Ipv6CidrBlockAssociationSet) == 0 {
		return fmt.Errorf("VPC %q does not have any associated IPv6 CIDR blocks", v.clusterConfig.VPC.ID)
	}
	return nil
}

func (v *ExistingVPCResourceSet) importExistingResources(ctx context.Context) error {
	if subnets := v.clusterConfig.VPC.Subnets.Private; subnets != nil {
		var (
			subnetRoutes map[string]string
			err          error
		)
		if v.clusterConfig.IsFullyPrivate() {
			subnetRoutes, err = importRouteTables(ctx, v.ec2API, v.clusterConfig.VPC.Subnets.Private)
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

func importRouteTables(ctx context.Context, ec2API awsapi.EC2, subnets map[string]api.AZSubnetSpec) (map[string]string, error) {
	var subnetIDs []string
	for _, subnet := range subnets {
		subnetIDs = append(subnetIDs, subnet.ID)
	}

	var routeTables []ec2types.RouteTable

	paginator := ec2.NewDescribeRouteTablesPaginator(ec2API, &ec2.DescribeRouteTablesInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("association.subnet-id"),
				Values: subnetIDs,
			},
		},
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "error describing route tables")
		}

		routeTables = append(routeTables, output.RouteTables...)
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

// RenderJSON returns the rendered JSON
func (v *ExistingVPCResourceSet) RenderJSON() ([]byte, error) {
	return v.rs.renderJSON()
}
