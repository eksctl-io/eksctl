package v1alpha5

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/utils/ipnet"
)

type endpointAccessCase struct {
	vpc       *ClusterVPC
	endpoints *ClusterEndpoints
	Public    *bool
	Private   *bool
	Expected  bool
}

type subnetCase struct {
	subnets            AZSubnetMapping
	localSubnetsConfig AZSubnetMapping
	az                 string
	subnetID           string
	cidr               string
	err                string
	expected           AZSubnetMapping
}

var _ = Describe("VPC Configuration", func() {
	DescribeTable("Subnet import",
		func(e subnetCase) {
			ec2Subnet := ec2types.Subnet{
				AvailabilityZone: aws.String(e.az),
				SubnetId:         aws.String(e.subnetID),
			}
			if len(e.cidr) > 0 {
				ec2Subnet.CidrBlock = aws.String(e.cidr)
			}
			err := ImportSubnet(e.subnets, e.localSubnetsConfig, &ec2Subnet, func(subnet *ec2types.Subnet) string {
				return *subnet.AvailabilityZone
			})
			if e.err != "" {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(e.err))
			} else {
				Expect(err).NotTo(HaveOccurred())
				Expect(e.subnets).To(Equal(e.expected))
			}
		},
		Entry("No subnets", subnetCase{
			subnets:            NewAZSubnetMapping(),
			localSubnetsConfig: NewAZSubnetMapping(),
			az:                 "us-east-1a",
			subnetID:           "subnet-1",
			cidr:               "192.168.0.0/16",
			expected: AZSubnetMappingFromMap(map[string]AZSubnetSpec{
				"us-east-1a": {
					AZ:   "us-east-1a",
					ID:   "subnet-1",
					CIDR: ipnet.MustParseCIDR("192.168.0.0/16"),
				},
			}),
		}),
		Entry("Existing subnets", subnetCase{
			subnets: AZSubnetMappingFromMap(map[string]AZSubnetSpec{
				"us-east-1a": {},
			}),
			localSubnetsConfig: AZSubnetMappingFromMap(map[string]AZSubnetSpec{
				"us-east-1a": {},
			}),
			az:       "us-east-1a",
			subnetID: "subnet-1",
			cidr:     "192.168.0.0/16",
			expected: AZSubnetMappingFromMap(map[string]AZSubnetSpec{
				"us-east-1a": {
					AZ:   "us-east-1a",
					ID:   "subnet-1",
					CIDR: ipnet.MustParseCIDR("192.168.0.0/16"),
				},
			}),
		}),
		Entry("Existing subnets w/o IPv4 CIDR", subnetCase{
			subnets: AZSubnetMappingFromMap(map[string]AZSubnetSpec{
				"us-east-1a": {},
			}),
			localSubnetsConfig: AZSubnetMappingFromMap(map[string]AZSubnetSpec{
				"us-east-1a": {},
			}),
			az:       "us-east-1a",
			subnetID: "subnet-1",
			expected: AZSubnetMappingFromMap(map[string]AZSubnetSpec{
				"us-east-1a": {
					AZ: "us-east-1a",
					ID: "subnet-1",
				},
			}),
		}),
		Entry("Existing subnet with ID", subnetCase{
			subnets: AZSubnetMappingFromMap(map[string]AZSubnetSpec{
				"us-east-1a": {
					ID: "subnet-1",
				},
			}),
			localSubnetsConfig: AZSubnetMappingFromMap(map[string]AZSubnetSpec{
				"us-east-1a": {
					ID: "subnet-1",
				},
			}),
			az:       "us-east-1a",
			subnetID: "subnet-1",
			cidr:     "192.168.0.0/16",
			expected: AZSubnetMappingFromMap(map[string]AZSubnetSpec{
				"us-east-1a": {
					AZ:   "us-east-1a",
					ID:   "subnet-1",
					CIDR: ipnet.MustParseCIDR("192.168.0.0/16"),
				},
			}),
		}),
		Entry("ID only subnet", subnetCase{
			subnets: AZSubnetMappingFromMap(map[string]AZSubnetSpec{
				"main-subnet": {
					ID: "subnet-1",
				},
			}),
			localSubnetsConfig: AZSubnetMappingFromMap(map[string]AZSubnetSpec{
				"main-subnet": {
					ID: "subnet-1",
				},
			}),
			az:       "us-east-1a",
			subnetID: "subnet-1",
			cidr:     "192.168.0.0/24",
			expected: AZSubnetMappingFromMap(map[string]AZSubnetSpec{
				"main-subnet": {
					AZ:   "us-east-1a",
					ID:   "subnet-1",
					CIDR: ipnet.MustParseCIDR("192.168.0.0/24"),
				},
			}),
		}),
		Entry("Conflicting existing subnets by ID", subnetCase{
			subnets: AZSubnetMappingFromMap(map[string]AZSubnetSpec{
				"us-east-1a": {
					ID: "subnet-2",
				},
			}),
			localSubnetsConfig: AZSubnetMappingFromMap(map[string]AZSubnetSpec{
				"us-east-1a": {
					ID: "subnet-2",
				},
			}),
			az:       "us-east-1a",
			subnetID: "subnet-1",
			cidr:     "192.168.0.0/16",
			err:      "mismatch found between local and remote VPC config: subnet ID",
		}),
		Entry("Conflicting existing subnets by CIDR", subnetCase{
			subnets: AZSubnetMappingFromMap(map[string]AZSubnetSpec{
				"us-east-1a": {
					ID:   "subnet-1",
					CIDR: ipnet.MustParseCIDR("192.168.1.0/24"),
				},
			}),
			localSubnetsConfig: AZSubnetMappingFromMap(map[string]AZSubnetSpec{
				"us-east-1a": {
					ID:   "subnet-1",
					CIDR: ipnet.MustParseCIDR("192.168.1.0/24"),
				},
			}),
			az:       "us-east-1a",
			subnetID: "subnet-1",
			cidr:     "192.168.0.0/16",
			err:      "mismatch found between local and remote VPC config: subnet CIDR",
		}),
		Entry("Named subnets placeholder", subnetCase{
			subnets: AZSubnetMappingFromMap(map[string]AZSubnetSpec{
				"main-subnet": {
					AZ: "us-east-1a",
				},
			}),
			localSubnetsConfig: AZSubnetMappingFromMap(map[string]AZSubnetSpec{
				"main-subnet": {
					AZ: "us-east-1a",
				},
			}),
			az:       "us-east-1a",
			subnetID: "subnet-1",
			cidr:     "192.168.0.0/16",
			expected: AZSubnetMappingFromMap(map[string]AZSubnetSpec{
				"main-subnet": {
					AZ:   "us-east-1a",
					ID:   "subnet-1",
					CIDR: ipnet.MustParseCIDR("192.168.0.0/16"),
				},
			}),
		}),
		Entry("Ambiguous ID list", subnetCase{
			subnets: AZSubnetMappingFromMap(map[string]AZSubnetSpec{
				"main-subnet": {
					ID: "subnet-1",
				},
				"other-subnet": {
					ID: "subnet-1",
				},
			}),
			localSubnetsConfig: AZSubnetMappingFromMap(map[string]AZSubnetSpec{
				"main-subnet": {
					ID: "subnet-1",
				},
				"other-subnet": {
					ID: "subnet-1",
				},
			}),
			az:       "us-east-1a",
			subnetID: "subnet-1",
			cidr:     "192.168.0.0/16",
			err:      "mismatch found between local and remote VPC config: unable to unambiguously import subnet by ID",
		}),
		Entry("Ambiguous CIDR+AZ list", subnetCase{
			subnets: AZSubnetMappingFromMap(map[string]AZSubnetSpec{
				"main-subnet": {
					AZ: "us-east-1a",
				},
				"other-subnet": {
					AZ: "us-east-1a",
				},
			}),
			localSubnetsConfig: AZSubnetMappingFromMap(map[string]AZSubnetSpec{
				"main-subnet": {
					AZ: "us-east-1a",
				},
				"other-subnet": {
					AZ: "us-east-1a",
				},
			}),
			az:       "us-east-1a",
			subnetID: "subnet-1",
			cidr:     "192.168.0.0/16",
			err:      "mismatch found between local and remote VPC config: unable to unambiguously import subnet by <AZ,CIDR> pair",
		}),
		Entry("CIDR+AZ differentiated list", subnetCase{
			subnets: AZSubnetMappingFromMap(map[string]AZSubnetSpec{
				"main-subnet": {
					AZ:   "us-east-1a",
					CIDR: ipnet.MustParseCIDR("192.168.0.0/24"),
				},
				"other-subnet": {
					AZ:   "us-east-1a",
					CIDR: ipnet.MustParseCIDR("192.168.1.0/24"),
				},
			}),
			localSubnetsConfig: AZSubnetMappingFromMap(map[string]AZSubnetSpec{
				"main-subnet": {
					AZ:   "us-east-1a",
					CIDR: ipnet.MustParseCIDR("192.168.0.0/24"),
				},
				"other-subnet": {
					AZ:   "us-east-1a",
					CIDR: ipnet.MustParseCIDR("192.168.1.0/24"),
				},
			}),
			az:       "us-east-1a",
			subnetID: "subnet-1",
			cidr:     "192.168.0.0/24",
			expected: AZSubnetMappingFromMap(map[string]AZSubnetSpec{
				"main-subnet": {
					AZ:   "us-east-1a",
					ID:   "subnet-1",
					CIDR: ipnet.MustParseCIDR("192.168.0.0/24"),
				},
				"other-subnet": {
					AZ:   "us-east-1a",
					CIDR: ipnet.MustParseCIDR("192.168.1.0/24"),
				},
			}),
		}),
		Entry("ID disambiguating list", subnetCase{
			subnets: AZSubnetMappingFromMap(map[string]AZSubnetSpec{
				"main-subnet": {
					AZ: "us-east-1a",
					ID: "subnet-1",
				},
				"other-subnet": {
					AZ: "us-east-1a",
				},
			}),
			localSubnetsConfig: AZSubnetMappingFromMap(map[string]AZSubnetSpec{
				"main-subnet": {
					AZ: "us-east-1a",
					ID: "subnet-1",
				},
				"other-subnet": {
					AZ: "us-east-1a",
				},
			}),
			az:       "us-east-1a",
			subnetID: "subnet-1",
			cidr:     "192.168.0.0/24",
			expected: AZSubnetMappingFromMap(map[string]AZSubnetSpec{
				"main-subnet": {
					AZ:   "us-east-1a",
					ID:   "subnet-1",
					CIDR: ipnet.MustParseCIDR("192.168.0.0/24"),
				},
				"other-subnet": {
					AZ: "us-east-1a",
				},
			}),
		}),
		Entry("Two subnets in same AZ, without VPC config provided", subnetCase{
			subnets: AZSubnetMappingFromMap(map[string]AZSubnetSpec{
				"us-east-1a": {
					ID:   "subnet-1",
					CIDR: ipnet.MustParseCIDR("192.168.0.0/24"),
				},
			}),
			localSubnetsConfig: NewAZSubnetMapping(),
			az:                 "us-east-1a",
			subnetID:           "subnet-2",
			cidr:               "192.168.32.0/26",
			expected: AZSubnetMappingFromMap(map[string]AZSubnetSpec{
				"us-east-1a": {
					ID:   "subnet-2",
					CIDR: ipnet.MustParseCIDR("192.168.32.0/26"),
				},
			}),
		}),
	)
	DescribeTable("Can determine if VPC config in config file has cluster endpoints",
		func(e endpointAccessCase) {
			cc := &ClusterConfig{}
			Expect(cc.HasClusterEndpointAccess()).Should(BeTrue())
			cc.VPC = e.vpc
			Expect(cc.HasClusterEndpointAccess()).Should(BeTrue())
			if cc.VPC != nil {
				cc.VPC.ClusterEndpoints = e.endpoints
			}
			if e.Public != nil && cc.VPC.ClusterEndpoints != nil {
				cc.VPC.ClusterEndpoints.PublicAccess = e.Public
			}
			if e.Private != nil && cc.VPC.ClusterEndpoints != nil {
				cc.VPC.ClusterEndpoints.PrivateAccess = e.Private
			}
			if e.Public != nil && e.Private != nil {
				Expect(cc.HasClusterEndpointAccess()).Should(Equal(e.Expected))
			}
		},
		Entry("No VPC", endpointAccessCase{
			vpc:       nil,
			endpoints: nil,
			Public:    nil,
			Private:   nil,
			Expected:  true,
		}),
		Entry("Has Empty VPC", endpointAccessCase{
			vpc:       &ClusterVPC{},
			endpoints: nil,
			Public:    nil,
			Private:   nil,
			Expected:  true,
		}),
		Entry("Public=true, Private=true", endpointAccessCase{
			vpc:       &ClusterVPC{},
			endpoints: &ClusterEndpoints{},
			Public:    Enabled(),
			Private:   Enabled(),
			Expected:  true,
		}),
		Entry("Public=true, Private=false", endpointAccessCase{
			vpc:       &ClusterVPC{},
			endpoints: &ClusterEndpoints{},
			Public:    Enabled(),
			Private:   Disabled(),
			Expected:  true,
		}),
		Entry("Public=false, Private=true", endpointAccessCase{
			vpc:       &ClusterVPC{},
			endpoints: &ClusterEndpoints{},
			Public:    Disabled(),
			Private:   Enabled(),
			Expected:  true,
		}),
		Entry("Public=false, Private=false", endpointAccessCase{
			vpc:       &ClusterVPC{},
			endpoints: &ClusterEndpoints{},
			Public:    Disabled(),
			Private:   Disabled(),
			Expected:  false,
		}),
	)
})
