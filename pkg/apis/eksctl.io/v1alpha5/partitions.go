package v1alpha5

import "fmt"

// Partitions.
const (
	PartitionAWS   = "aws"
	PartitionChina = "aws-cn"
	PartitionUSGov = "aws-us-gov"
	PartitionISO   = "aws-iso"
	PartitionISOB  = "aws-iso-b"
)

// partition is an AWS partition.
type partition struct {
	name                        string
	serviceMappings             map[string]string
	regions                     []string
	endpointServiceDomainPrefix string
}

type partitions []partition

var standardServiceMappings = map[string]string{
	"EC2":            "ec2.amazonaws.com",
	"EKS":            "eks.amazonaws.com",
	"EKSFargatePods": "eks-fargate-pods.amazonaws.com",
}

const standardPartitionServiceDomainPrefix = "com.amazonaws"

var awsPartition = partition{
	name:                        PartitionAWS,
	serviceMappings:             standardServiceMappings,
	endpointServiceDomainPrefix: standardPartitionServiceDomainPrefix,
}

// Partitions is a list of supported AWS partitions.
var Partitions = partitions{
	awsPartition,
	{
		name:                        PartitionUSGov,
		serviceMappings:             standardServiceMappings,
		regions:                     []string{RegionUSGovEast1, RegionUSGovWest1},
		endpointServiceDomainPrefix: standardPartitionServiceDomainPrefix,
	},
	{
		name: PartitionChina,
		serviceMappings: map[string]string{
			"EC2":            "ec2.amazonaws.com.cn",
			"EKS":            "eks.amazonaws.com",
			"EKSFargatePods": "eks-fargate-pods.amazonaws.com",
		},
		regions:                     []string{RegionCNNorth1, RegionCNNorthwest1},
		endpointServiceDomainPrefix: fmt.Sprintf("cn.%s", standardPartitionServiceDomainPrefix),
	},
	{
		name: PartitionISO,
		serviceMappings: map[string]string{
			"EC2":            "ec2.c2s.ic.gov",
			"EKS":            "eks.amazonaws.com",
			"EKSFargatePods": "eks-fargate-pods.amazonaws.com",
		},
		regions:                     []string{RegionUSISOEast1},
		endpointServiceDomainPrefix: "gov.ic.c2s",
	},
	{
		name: PartitionISOB,
		serviceMappings: map[string]string{
			"EC2":            "ec2.sc2s.sgov.gov",
			"EKS":            "eks.amazonaws.com",
			"EKSFargatePods": "eks-fargate-pods.amazonaws.com",
		},
		regions:                     []string{RegionUSISOBEast1},
		endpointServiceDomainPrefix: "gov.sgov.sc2s",
	},
}

// ForRegion returns the partition a region belongs to.
func (p partitions) ForRegion(region string) string {
	for _, pt := range p {
		for _, r := range pt.regions {
			if r == region {
				return pt.name
			}
		}
	}
	return PartitionAWS
}

// GetEndpointServiceDomainPrefix returns the domain prefix for the endpoint service.
func (p partitions) GetEndpointServiceDomainPrefix(endpointService EndpointService, region string) string {
	for _, pt := range p {
		for _, r := range pt.regions {
			if r == region {
				switch pt.name {
				case PartitionChina:
					if endpointService.RequiresChinaPrefix {
						return pt.endpointServiceDomainPrefix
					}
					return standardPartitionServiceDomainPrefix
				case PartitionISO, PartitionISOB:
					if endpointService.RequiresISOPrefix {
						return pt.endpointServiceDomainPrefix
					}
					return standardPartitionServiceDomainPrefix
				default:
					return pt.endpointServiceDomainPrefix
				}
			}
		}
	}
	return awsPartition.endpointServiceDomainPrefix
}

// IsSupported returns true if the partition is supported.
func (p partitions) IsSupported(partition string) bool {
	for _, pt := range p {
		if pt.name == partition {
			return true
		}
	}
	return false
}

// ServicePrincipalPartitionMappings returns the service principal partition mappings for all supported partitions.
func (p partitions) ServicePrincipalPartitionMappings() map[string]map[string]string {
	ret := make(map[string]map[string]string, len(p))
	for _, pt := range p {
		ret[pt.name] = pt.serviceMappings
	}
	return ret
}
