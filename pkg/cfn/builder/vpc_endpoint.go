package builder

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"

	"github.com/pkg/errors"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"

	gfnec2 "github.com/weaveworks/goformation/v4/cloudformation/ec2"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"
)

// A VPCEndpointResourceSet represents the resources required for VPC endpoints
type VPCEndpointResourceSet struct {
	provider        provider
	rs              *resourceSet
	vpc             *gfnt.Value
	clusterConfig   *api.ClusterConfig
	subnets         []subnetResource
	clusterSharedSG *gfnt.Value
}

type provider interface {
	EC2() ec2iface.EC2API
	Region() string
}

// NewVPCEndpointResourceSet creates a new VPCEndpointResourceSet
func NewVPCEndpointResourceSet(provider provider, rs *resourceSet, clusterConfig *api.ClusterConfig, vpc *gfnt.Value, subnets []subnetResource, clusterSharedSG *gfnt.Value) *VPCEndpointResourceSet {
	return &VPCEndpointResourceSet{
		provider:        provider,
		rs:              rs,
		clusterConfig:   clusterConfig,
		vpc:             vpc,
		subnets:         subnets,
		clusterSharedSG: clusterSharedSG,
	}
}

// VPCEndpointServiceDetails holds the details for a VPC endpoint service
type VPCEndpointServiceDetails struct {
	ServiceName       string
	Service           string
	EndpointType      string
	AvailabilityZones []string
}

// AddResources adds resources for VPC endpoints
func (e *VPCEndpointResourceSet) AddResources() error {
	endpointServices := append(api.RequiredEndpointServices(), e.clusterConfig.PrivateCluster.AdditionalEndpointServices...)
	if e.clusterConfig.HasClusterCloudWatchLogging() && !e.hasEndpoint(api.EndpointServiceCloudWatch) {
		endpointServices = append(endpointServices, api.EndpointServiceCloudWatch)
	}
	endpointServiceDetails, err := BuildVPCEndpointServices(e.provider.EC2(), e.provider.Region(), endpointServices)
	if err != nil {
		return errors.Wrap(err, "error building endpoint service details")
	}

	for _, endpointDetail := range endpointServiceDetails {
		endpoint := &gfnec2.VPCEndpoint{
			ServiceName:     gfnt.NewString(endpointDetail.ServiceName),
			VpcId:           e.vpc,
			VpcEndpointType: gfnt.NewString(endpointDetail.EndpointType),
		}

		if endpointDetail.EndpointType == ec2.VpcEndpointTypeGateway {
			endpoint.RouteTableIds = gfnt.NewSlice(e.routeTableIDs()...)
		} else {
			endpoint.SubnetIds = gfnt.NewSlice(e.subnetsForAZs(endpointDetail.AvailabilityZones)...)
			endpoint.PrivateDnsEnabled = gfnt.NewBoolean(true)
			endpoint.SecurityGroupIds = gfnt.NewSlice(e.clusterSharedSG)
		}

		resourceName := fmt.Sprintf("VPCEndpoint%s", strings.ToUpper(
			strings.ReplaceAll(endpointDetail.Service, ".", ""),
		))

		// TODO attach policy document
		e.rs.newResource(resourceName, endpoint)
	}
	return nil
}

func (e *VPCEndpointResourceSet) subnetsForAZs(azs []string) []*gfnt.Value {
	var subnetRefs []*gfnt.Value
	for _, az := range azs {
		for _, subnet := range e.subnets {
			if subnet.AvailabilityZone == az {
				subnetRefs = append(subnetRefs, subnet.Subnet)
			}
		}

	}
	return subnetRefs
}

func (e *VPCEndpointResourceSet) routeTableIDs() []*gfnt.Value {
	var routeTableIDs []*gfnt.Value
	for _, subnet := range e.subnets {
		routeTableIDs = append(routeTableIDs, subnet.RouteTable)
	}
	return routeTableIDs
}

func (e *VPCEndpointResourceSet) hasEndpoint(endpoint string) bool {
	for _, ae := range e.clusterConfig.PrivateCluster.AdditionalEndpointServices {
		if ae == endpoint {
			return true
		}
	}
	return false
}

// BuildVPCEndpointServices builds a slice of VPCEndpointServiceDetails for the specified endpoint names
func BuildVPCEndpointServices(ec2API ec2iface.EC2API, region string, endpoints []string) ([]VPCEndpointServiceDetails, error) {
	serviceNames := make([]string, len(endpoints))
	serviceRegionPrefix := fmt.Sprintf("com.amazonaws.%s.", region)
	for i, endpoint := range endpoints {
		serviceNames[i] = serviceRegionPrefix + endpoint
	}

	var serviceDetails []*ec2.ServiceDetail
	var nextToken *string

	for {
		output, err := ec2API.DescribeVpcEndpointServices(&ec2.DescribeVpcEndpointServicesInput{
			ServiceNames: aws.StringSlice(serviceNames),
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("service-name"),
					Values: aws.StringSlice(serviceNames),
				},
			},
			NextToken: nextToken,
		})

		if err != nil {
			return nil, errors.Wrap(err, "error describing VPC endpoint services")
		}
		serviceDetails = append(serviceDetails, output.ServiceDetails...)

		if nextToken = output.NextToken; nextToken == nil {
			break
		}
	}

	ret := make([]VPCEndpointServiceDetails, len(serviceDetails))

	for i, sd := range serviceDetails {
		if len(sd.ServiceType) > 1 {
			return nil, errors.Errorf("endpoint service %q with multiple service types isn't supported", *sd.ServiceName)
		}

		ret[i] = VPCEndpointServiceDetails{
			ServiceName:       *sd.ServiceName,
			Service:           strings.TrimPrefix(*sd.ServiceName, serviceRegionPrefix),
			EndpointType:      *sd.ServiceType[0].ServiceType,
			AvailabilityZones: aws.StringValueSlice(sd.AvailabilityZones),
		}
	}

	return ret, nil
}
