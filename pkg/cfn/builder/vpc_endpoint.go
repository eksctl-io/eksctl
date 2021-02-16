package builder

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/pkg/errors"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"

	gfnec2 "github.com/weaveworks/goformation/v4/cloudformation/ec2"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"
)

// A VPCEndpointResourceSet represents the resources required for VPC endpoints
type VPCEndpointResourceSet struct {
	ec2API          ec2iface.EC2API
	region          string
	rs              *resourceSet
	vpc             *gfnt.Value
	clusterConfig   *api.ClusterConfig
	subnets         []subnetResource
	clusterSharedSG *gfnt.Value
}

// NewVPCEndpointResourceSet creates a new VPCEndpointResourceSet
func NewVPCEndpointResourceSet(ec2API ec2iface.EC2API, region string, rs *resourceSet, clusterConfig *api.ClusterConfig, vpc *gfnt.Value, subnets []subnetResource, clusterSharedSG *gfnt.Value) *VPCEndpointResourceSet {
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

// VPCEndpointServiceDetails holds the details for a VPC endpoint service
type VPCEndpointServiceDetails struct {
	ServiceName         string
	ServiceReadableName string
	EndpointType        string
	AvailabilityZones   []string
}

// AddResources adds resources for VPC endpoints
func (e *VPCEndpointResourceSet) AddResources() error {
	endpointServices := append(api.RequiredEndpointServices(), e.clusterConfig.PrivateCluster.AdditionalEndpointServices...)
	if e.clusterConfig.HasClusterCloudWatchLogging() && !e.hasEndpoint(api.EndpointServiceCloudWatch) {
		endpointServices = append(endpointServices, api.EndpointServiceCloudWatch)
	}
	endpointServiceDetails, err := buildVPCEndpointServices(e.ec2API, e.region, endpointServices)
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

func (e *VPCEndpointResourceSet) hasEndpoint(endpoint string) bool {
	for _, ae := range e.clusterConfig.PrivateCluster.AdditionalEndpointServices {
		if ae == endpoint {
			return true
		}
	}
	return false
}

var chinaPartitionServiceHasChinaPrefix = map[string]bool{
	api.EndpointServiceEC2:            true,
	api.EndpointServiceECRAPI:         true,
	api.EndpointServiceECRDKR:         true,
	api.EndpointServiceS3:             false,
	api.EndpointServiceSTS:            true,
	api.EndpointServiceCloudFormation: true,
	api.EndpointServiceAutoscaling:    true,
	api.EndpointServiceCloudWatch:     false,
}

// buildVPCEndpointServices builds a slice of VPCEndpointServiceDetails for the specified endpoint names
func buildVPCEndpointServices(ec2API ec2iface.EC2API, region string, endpoints []string) ([]VPCEndpointServiceDetails, error) {
	serviceNames := make([]string, len(endpoints))
	serviceDomain := fmt.Sprintf("com.amazonaws.%s", region)
	for i, endpoint := range endpoints {
		serviceName, err := makeServiceName(region, endpoint)
		if err != nil {
			return nil, err
		}
		serviceNames[i] = serviceName
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
			if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "InvalidServiceName" {
				return nil, &api.UnsupportedFeatureError{
					Message: fmt.Sprintf("fully-private clusters are not supported in region %q, please retry with a different region", region),
					Err:     err,
				}
			}
			return nil, errors.Wrap(err, "error describing VPC endpoint services")
		}

		serviceDetails = append(serviceDetails, output.ServiceDetails...)
		if nextToken = output.NextToken; nextToken == nil {
			break
		}
	}

	var ret []VPCEndpointServiceDetails
	s3EndpointName, err := makeServiceName(region, api.EndpointServiceS3)
	if err != nil {
		return nil, err
	}

	for _, sd := range serviceDetails {
		if len(sd.ServiceType) > 1 {
			return nil, errors.Errorf("endpoint service %q with multiple service types isn't supported", *sd.ServiceName)
		}

		endpointType := *sd.ServiceType[0].ServiceType
		if !serviceEndpointTypeExpected(*sd.ServiceName, endpointType, s3EndpointName) {
			continue
		}

		// Trim the domain (potentially with a partition-specific part) from the `ServiceName`
		parts := strings.Split(*sd.ServiceName, fmt.Sprintf("%s.", serviceDomain))
		if len(parts) != 2 {
			return nil, errors.Errorf("error parsing service name %s %s", *sd.ServiceName, serviceDomain)
		}
		readableName := parts[1]

		ret = append(ret, VPCEndpointServiceDetails{
			ServiceName:         *sd.ServiceName,
			ServiceReadableName: readableName,
			EndpointType:        endpointType,
			AvailabilityZones:   aws.StringValueSlice(sd.AvailabilityZones),
		})
	}

	return ret, nil
}

// serviceEndpointTypeExpected returns true if the endpoint service is expected to use the specified endpoint type
func serviceEndpointTypeExpected(serviceName, endpointType, s3EndpointName string) bool {
	if serviceName == s3EndpointName {
		return endpointType == ec2.VpcEndpointTypeGateway
	}
	return endpointType == ec2.VpcEndpointTypeInterface
}

func makeServiceName(region, endpoint string) (string, error) {
	serviceName := fmt.Sprintf("com.amazonaws.%s.%s", region, endpoint)
	hasChinaPrefix, ok := chinaPartitionServiceHasChinaPrefix[endpoint]
	if !ok {
		return "", errors.Errorf("couldn't determine endpoint domain for service %s", endpoint)
	}
	if api.Partition(region) == api.PartitionChina && hasChinaPrefix {
		serviceName = "cn." + serviceName
	}
	return serviceName, nil
}
