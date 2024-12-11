package builder

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/smithy-go"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"

	gfnec2 "goformation/v4/cloudformation/ec2"
	gfnt "goformation/v4/cloudformation/types"
)

// A VPCEndpointResourceSet holds the resources required for VPC endpoints.
type VPCEndpointResourceSet struct {
	ec2API          awsapi.EC2
	region          string
	rs              *resourceSet
	vpc             *gfnt.Value
	clusterConfig   *api.ClusterConfig
	subnets         []SubnetResource
	clusterSharedSG *gfnt.Value
}

// NewVPCEndpointResourceSet creates a new VPCEndpointResourceSet.
func NewVPCEndpointResourceSet(ec2API awsapi.EC2, region string, rs *resourceSet, clusterConfig *api.ClusterConfig, vpc *gfnt.Value, subnets []SubnetResource, clusterSharedSG *gfnt.Value) *VPCEndpointResourceSet {
	return &VPCEndpointResourceSet{
		ec2API:          ec2API,
		region:          region,
		rs:              rs,
		clusterConfig:   clusterConfig,
		vpc:             vpc,
		subnets:         subnets,
		clusterSharedSG: clusterSharedSG,
	}
}

// VPCEndpointServiceDetails holds the details for a VPC endpoint service.
type VPCEndpointServiceDetails struct {
	ServiceName         string
	ServiceReadableName string
	EndpointType        string
	AvailabilityZones   []string
}

// AddResources adds resources for VPC endpoints.
func (e *VPCEndpointResourceSet) AddResources(ctx context.Context) error {
	additionalServices, err := api.MapOptionalEndpointServices(e.clusterConfig.PrivateCluster.AdditionalEndpointServices, e.clusterConfig.HasClusterCloudWatchLogging())
	if err != nil {
		return err
	}
	endpointServices := append(api.RequiredEndpointServices(e.clusterConfig.IsControlPlaneOnOutposts()), additionalServices...)
	endpointServiceDetails, err := e.buildVPCEndpointServices(ctx, endpointServices)
	if err != nil {
		return fmt.Errorf("error building endpoint service details: %w", err)
	}

	for _, endpointDetail := range endpointServiceDetails {
		endpoint := &gfnec2.VPCEndpoint{
			ServiceName:     gfnt.NewString(endpointDetail.ServiceName),
			VpcId:           e.vpc,
			VpcEndpointType: gfnt.NewString(endpointDetail.EndpointType),
		}

		if endpointDetail.EndpointType == string(ec2types.VpcEndpointTypeGateway) {
			endpoint.RouteTableIds = gfnt.NewSlice(e.routeTableIDs()...)
		} else {
			endpoint.SubnetIds = gfnt.NewSlice(e.subnetsForAZs(endpointDetail.AvailabilityZones)...)
			endpoint.PrivateDnsEnabled = gfnt.NewBoolean(true)
			endpoint.SecurityGroupIds = gfnt.NewSlice(e.clusterSharedSG)
		}

		resourceName := fmt.Sprintf("VPCEndpoint%s", strings.ToUpper(
			strings.ReplaceAll(endpointDetail.ServiceReadableName, ".", ""),
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
	m := make(map[string]bool)
	for _, subnet := range e.subnets {
		routeTableStr := subnet.RouteTable.String()

		if !m[routeTableStr] {
			m[routeTableStr] = true
			routeTableIDs = append(routeTableIDs, subnet.RouteTable)
		}
	}
	return routeTableIDs
}

// buildVPCEndpointServices builds a slice of VPCEndpointServiceDetails for the specified endpoint names.
func (e *VPCEndpointResourceSet) buildVPCEndpointServices(ctx context.Context, endpointServices []api.EndpointService) ([]VPCEndpointServiceDetails, error) {
	serviceNames := make([]string, len(endpointServices))
	for i, endpoint := range endpointServices {
		serviceNames[i] = makeServiceName(endpoint, e.region)
	}

	var (
		serviceDetails []ec2types.ServiceDetail
		nextToken      *string
	)

	for {
		output, err := e.ec2API.DescribeVpcEndpointServices(ctx, &ec2.DescribeVpcEndpointServicesInput{
			ServiceNames: serviceNames,
			Filters: []ec2types.Filter{
				{
					Name:   aws.String("service-name"),
					Values: serviceNames,
				},
			},
			NextToken: nextToken,
		})

		if err != nil {
			var ae smithy.APIError
			if errors.As(err, &ae) && ae.ErrorCode() == "InvalidServiceName" {
				return nil, &api.UnsupportedFeatureError{
					Message: fmt.Sprintf("fully-private clusters are not supported in region %q, please retry with a different region", e.region),
					Err:     err,
				}
			}
			return nil, fmt.Errorf("error describing VPC endpoint services: %w", err)
		}

		serviceDetails = append(serviceDetails, output.ServiceDetails...)
		if nextToken = output.NextToken; nextToken == nil {
			break
		}
	}

	var ret []VPCEndpointServiceDetails
	s3ServiceName := makeServiceName(api.EndpointServiceS3, e.region)
	for _, sd := range serviceDetails {
		if len(sd.ServiceType) > 1 {
			return nil, fmt.Errorf("endpoint service %q with multiple service types isn't supported", *sd.ServiceName)
		}
		if len(sd.ServiceType) == 0 {
			return nil, fmt.Errorf("expected to find a service type for endpoint service %q", *sd.ServiceName)
		}

		endpointType := sd.ServiceType[0].ServiceType
		if !serviceEndpointTypeExpected(*sd.ServiceName, endpointType, s3ServiceName) {
			continue
		}

		readableName, err := makeReadableName(*sd.ServiceName, e.region)
		if err != nil {
			return nil, err
		}

		ret = append(ret, VPCEndpointServiceDetails{
			ServiceName:         *sd.ServiceName,
			ServiceReadableName: readableName,
			EndpointType:        string(endpointType),
			AvailabilityZones:   sd.AvailabilityZones,
		})
	}

	return ret, nil
}

func makeReadableName(serviceName, region string) (string, error) {
	search := fmt.Sprintf(".%s.", region)
	idx := strings.Index(serviceName, search)
	if idx == -1 {
		return "", fmt.Errorf("unexpected format for endpoint service name: %q", serviceName)
	}
	return serviceName[idx+len(search)-1:], nil
}

// serviceEndpointTypeExpected returns true if the endpoint service is expected to use the specified endpoint type.
func serviceEndpointTypeExpected(serviceName string, endpointType ec2types.ServiceType, s3ServiceName string) bool {
	if serviceName == s3ServiceName {
		return endpointType == ec2types.ServiceTypeGateway
	}
	return endpointType == ec2types.ServiceTypeInterface
}

func makeServiceName(endpointService api.EndpointService, region string) string {
	serviceDomainPrefix := api.Partitions.GetEndpointServiceDomainPrefix(endpointService, region)
	return fmt.Sprintf("%s.%s.%s", serviceDomainPrefix, region, endpointService.Name)
}
