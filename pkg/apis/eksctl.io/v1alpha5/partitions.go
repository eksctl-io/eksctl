package v1alpha5

import "fmt"

// Partitions.
const (
	PartitionAWS   = "aws"
	PartitionChina = "aws-cn"
	PartitionUSGov = "aws-us-gov"
	PartitionISO   = "aws-iso"
	PartitionISOB  = "aws-iso-b"
	PartitionISOF  = "aws-iso-f"
	PartitionISOE  = "aws-iso-e"
)

// partition is an AWS partition.
type partition struct {
	name                           string
	serviceMappings                map[string]string
	regions                        []string
	endpointServiceDomainPrefix    string
	endpointServiceDomainPrefixAlt string
	v1SDKDNSPrefix                 string
}

type partitions []partition

var standardServiceMappings = map[string]string{
	"EC2":            "ec2.amazonaws.com",
	"EKS":            "eks.amazonaws.com",
	"SSM":            "ssm.amazonaws.com",
	"IRA":            "rolesanywhere.amazonaws.com",
	"EKSFargatePods": "eks-fargate-pods.amazonaws.com",
}

const standardPartitionServiceDomainPrefix = "com.amazonaws"

var awsPartition = partition{
	name:                        PartitionAWS,
	serviceMappings:             standardServiceMappings,
	endpointServiceDomainPrefix: standardPartitionServiceDomainPrefix,
	v1SDKDNSPrefix:              "amazonaws.com",
	regions: []string{
		RegionUSWest1,
		RegionUSWest2,
		RegionUSEast1,
		RegionUSEast2,
		RegionCACentral1,
		RegionCAWest1,
		RegionEUWest1,
		RegionEUWest2,
		RegionEUWest3,
		RegionEUNorth1,
		RegionEUCentral1,
		RegionEUCentral2,
		RegionEUSouth1,
		RegionEUSouth2,
		RegionAPNorthEast1,
		RegionAPNorthEast2,
		RegionAPNorthEast3,
		RegionAPSouthEast1,
		RegionAPSouthEast2,
		RegionAPSouthEast3,
		RegionAPSouthEast4,
		RegionAPSouthEast5,
		RegionAPSouthEast7,
		RegionAPSouth1,
		RegionAPSouth2,
		RegionAPEast1,
		RegionAPEast2,
		RegionMECentral1,
		RegionMESouth1,
		RegionSAEast1,
		RegionAFSouth1,
		RegionILCentral1,
		RegionMXCentral1,
		RegionAPSoutheast6,
	},
}

func (p partition) Name() string {
	return p.name
}

// Partitions is a list of supported AWS partitions.
var Partitions = partitions{
	awsPartition,
	{
		name:                        PartitionUSGov,
		serviceMappings:             standardServiceMappings,
		regions:                     []string{RegionUSGovEast1, RegionUSGovWest1},
		endpointServiceDomainPrefix: standardPartitionServiceDomainPrefix,
		v1SDKDNSPrefix:              "amazonaws.com",
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
		v1SDKDNSPrefix:              "amazonaws.com.cn",
	},
	{
		name: PartitionISO,
		serviceMappings: map[string]string{
			"EC2":            "ec2.c2s.ic.gov",
			"EKS":            "eks.amazonaws.com",
			"EKSFargatePods": "eks-fargate-pods.amazonaws.com",
		},
		regions:                     []string{RegionUSISOEast1, RegionUSISOWest1},
		endpointServiceDomainPrefix: "gov.ic.c2s",
		v1SDKDNSPrefix:              "c2s.ic.gov",
	},
	{
		name: PartitionISOB,
		serviceMappings: map[string]string{
			"EC2":            "ec2.sc2s.sgov.gov",
			"EKS":            "eks.amazonaws.com",
			"EKSFargatePods": "eks-fargate-pods.amazonaws.com",
		},
		regions:                     []string{RegionUSISOBEast1,RegionUSISOBWest1},
		endpointServiceDomainPrefix: "gov.sgov.sc2s",
		v1SDKDNSPrefix:              "sc2s.sgov.gov",
	},
	{
		name: PartitionISOE,
		serviceMappings: map[string]string{
			"EC2":            "ec2.amazonaws.com",
			"EKS":            "eks.amazonaws.com",
			"EKSFargatePods": "eks-fargate-pods.amazonaws.com",
		},
		regions:                        []string{RegionEUISOEWest1},
		endpointServiceDomainPrefix:    standardPartitionServiceDomainPrefix,
		endpointServiceDomainPrefixAlt: "uk.adc-e.cloud",
		v1SDKDNSPrefix:                 "cloud.adc-e.uk",
	},
	{
		name: PartitionISOF,
		serviceMappings: map[string]string{
			"EC2":            "ec2.amazonaws.com",
			"EKS":            "eks.amazonaws.com",
			"EKSFargatePods": "eks-fargate-pods.amazonaws.com",
		},
		regions:                        []string{RegionUSISOFSouth1, RegionUSISOFEast1},
		endpointServiceDomainPrefix:    standardPartitionServiceDomainPrefix,
		endpointServiceDomainPrefixAlt: "gov.ic.hci.csp",
		v1SDKDNSPrefix:                 "csp.hci.ic.gov",
	},
}

func (p partitions) partitionFromRegion(region string) *partition {
	for _, pt := range p {
		for _, r := range pt.regions {
			if r == region {
				return &pt
			}
		}
	}
	return nil
}

// ForRegion returns the partition a region belongs to.
func (p partitions) ForRegion(region string) string {
	pt := p.partitionFromRegion(region)
	if pt == nil {
		return PartitionAWS
	}
	return pt.name
}

func (p partitions) V1SDKDNSPrefixForRegion(region string) (string, error) {
	pt := p.partitionFromRegion(region)
	if pt == nil {
		return "", fmt.Errorf("failed to find DNS suffix for region %s", region)
	}
	return pt.v1SDKDNSPrefix, nil
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
				case PartitionISOE, PartitionISOF:
					if endpointService.RequiresISOPrefix {
						//in these partitions four endpoints have an alternate domain prefix
						switch endpointService.Name {
						case "ebs", "ecr.api", "ecr.dkr", "execute-api":
							return pt.endpointServiceDomainPrefixAlt
						default:
							return pt.endpointServiceDomainPrefix
						}
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
