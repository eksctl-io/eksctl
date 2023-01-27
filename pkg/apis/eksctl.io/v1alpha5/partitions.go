package v1alpha5

// Partitions.
const (
	PartitionAWS   = "aws"
	PartitionChina = "aws-cn"
	PartitionUSGov = "aws-us-gov"
	PartitionISO   = "aws-iso"
	PartitionISOB  = "aws-iso-b"
)

type partition struct {
	name            string
	serviceMappings map[string]string
	regions         []string
}

type partitions []partition

var standardServiceMappings = map[string]string{
	"EC2":            "ec2.amazonaws.com",
	"EKS":            "eks.amazonaws.com",
	"EKSFargatePods": "eks-fargate-pods.amazonaws.com",
}

// Partitions is a list of supported AWS partitions.
var Partitions = partitions{
	{
		name:            PartitionAWS,
		serviceMappings: standardServiceMappings,
	},
	{
		name:            PartitionUSGov,
		serviceMappings: standardServiceMappings,
		regions:         []string{RegionUSGovEast1, RegionUSGovWest1},
	},
	{
		name: PartitionChina,
		serviceMappings: map[string]string{
			"EC2":            "ec2.amazonaws.com.cn",
			"EKS":            "eks.amazonaws.com",
			"EKSFargatePods": "eks-fargate-pods.amazonaws.com",
		},
		regions: []string{RegionCNNorth1, RegionCNNorthwest1},
	},
	{
		name: PartitionISO,
		serviceMappings: map[string]string{
			"EC2":            "ec2.c2s.ic.gov",
			"EKS":            "eks.amazonaws.com",
			"EKSFargatePods": "eks-fargate-pods.amazonaws.com",
		},
		regions: []string{RegionUSISOEast1},
	},
	{
		name: PartitionISOB,
		serviceMappings: map[string]string{
			"EC2":            "ec2.sc2s.sgov.gov",
			"EKS":            "eks.amazonaws.com",
			"EKSFargatePods": "eks-fargate-pods.amazonaws.com",
		},
		regions: []string{RegionUSISOBEast1},
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
