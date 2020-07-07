package builder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/vpc"
	gfn "github.com/weaveworks/goformation/cloudformation"
)

type vpcResourceSetCase struct {
	clusterConfig  *api.ClusterConfig
	expectedFile   string
	createProvider func() api.ClusterProvider
}

var _ = Describe("VPC Endpoint Builder", func() {

	DescribeTable("Add resources", func(vc vpcResourceSetCase) {
		api.SetClusterConfigDefaults(vc.clusterConfig)

		if len(vc.clusterConfig.AvailabilityZones) == 0 {
			vc.clusterConfig.AvailabilityZones = []string{"us-west-2a", "us-west-2b", "us-west-2c", "us-west-2d"}
		}
		if vc.clusterConfig.VPC.ID == "" {
			Expect(vpc.SetSubnets(vc.clusterConfig.VPC, vc.clusterConfig.AvailabilityZones)).To(Succeed())
		}

		var provider api.ClusterProvider
		if vc.createProvider != nil {
			provider = vc.createProvider()
		} else {
			provider = mockprovider.NewMockProvider()
		}

		rs := newResourceSet()
		vpcResourceSet := NewVPCResourceSet(rs, vc.clusterConfig, provider)
		vpcResource, err := vpcResourceSet.AddResources()
		Expect(err).ToNot(HaveOccurred())

		if vc.clusterConfig.PrivateCluster.Enabled {
			vpcEndpointResourceSet := NewVPCEndpointResourceSet(provider, rs, vc.clusterConfig, vpcResource.VPC, vpcResource.SubnetDetails.Private, gfn.NewString("sg-test"))
			Expect(vpcEndpointResourceSet.AddResources()).To(Succeed())
			s3Endpoint := rs.template.Resources["VPCEndpointS3"].(*gfn.AWSEC2VPCEndpoint)
			sort.Slice(s3Endpoint.RouteTableIds, func(i, j int) bool {
				return s3Endpoint.RouteTableIds[i].String() < s3Endpoint.RouteTableIds[j].String()
			})
		} else if vc.clusterConfig.VPC.ID != "" {
			Expect(rs.template.Resources).To(BeEmpty())
			return
		}

		resourceJSON, err := rs.template.JSON()
		Expect(err).ToNot(HaveOccurred())

		expectedJSON, err := ioutil.ReadFile("testdata/" + vc.expectedFile)
		Expect(err).ToNot(HaveOccurred())
		Expect(resourceJSON).To(MatchJSON(expectedJSON))
	},
		Entry("Standard cluster", vpcResourceSetCase{
			clusterConfig: api.NewClusterConfig(),
			expectedFile:  "vpc_public.json",
		}),
		Entry("Private cluster", vpcResourceSetCase{
			clusterConfig: &api.ClusterConfig{
				VPC: api.NewClusterVPC(),
				PrivateCluster: &api.PrivateCluster{
					Enabled: true,
				},
			},
			createProvider: func() api.ClusterProvider {
				provider := mockprovider.NewMockProvider()
				mockDescribeVPCEndpoints(provider)
				return provider
			},
			expectedFile: "vpc_private.json",
		}),
		Entry("Non-private cluster with a user-supplied VPC", vpcResourceSetCase{
			clusterConfig: &api.ClusterConfig{
				AvailabilityZones: []string{"us-west-2a", "us-west-2b"},
				VPC: &api.ClusterVPC{
					Network: api.Network{
						ID: "vpc-custom",
					},
					Subnets: &api.ClusterSubnets{
						Private: map[string]api.Network{
							"us-west-2a": {
								ID: "subnet-custom1",
							},
							"us-west-2b": {
								ID: "subnet-custom2",
							},
						},
						Public: map[string]api.Network{},
					},
				},
			},
		}),
		Entry("Private cluster with a user-supplied VPC", vpcResourceSetCase{
			clusterConfig: &api.ClusterConfig{
				PrivateCluster: &api.PrivateCluster{
					Enabled: true,
				},
				VPC: &api.ClusterVPC{
					Network: api.Network{
						ID: "vpc-custom",
					},
					Subnets: &api.ClusterSubnets{
						Private: map[string]api.Network{
							"us-west-2a": {
								ID: "subnet-custom1",
							},
							"us-west-2b": {
								ID: "subnet-custom2",
							},
						},
					},
				},
			},
			createProvider: func() api.ClusterProvider {
				provider := mockprovider.NewMockProvider()
				mockDescribeVPCEndpoints(provider)
				mockDescribeRouteTables(provider, []string{"subnet-custom1", "subnet-custom2"})
				return provider
			},
			expectedFile: "custom_vpc_private_endpoint.json",
		}),
	)
})

func mockDescribeVPCEndpoints(provider *mockprovider.MockProvider) {
	serviceDetailsJSON := `
{
    "ServiceNames": [
        "com.amazonaws.us-west-2.ec2", 
        "com.amazonaws.us-west-2.ecr.api", 
        "com.amazonaws.us-west-2.ecr.dkr", 
        "com.amazonaws.us-west-2.s3", 
        "com.amazonaws.us-west-2.sts"
    ], 
    "ServiceDetails": [
        {
            "ServiceType": [
                {
                    "ServiceType": "Interface"
                }
            ], 
            "Tags": [], 
            "ManagesVpcEndpoints": false, 
            "PrivateDnsName": "ec2.us-west-2.amazonaws.com", 
            "ServiceName": "com.amazonaws.us-west-2.ec2", 
            "VpcEndpointPolicySupported": true, 
            "ServiceId": "vpce-svc-0ee6723c76642b3d8", 
            "Owner": "amazon", 
            "AvailabilityZones": [
                "us-west-2a", 
                "us-west-2b", 
                "us-west-2c"
            ], 
            "AcceptanceRequired": false, 
            "BaseEndpointDnsNames": [
                "ec2.us-west-2.vpce.amazonaws.com"
            ]
        }, 
        {
            "ServiceType": [
                {
                    "ServiceType": "Interface"
                }
            ], 
            "Tags": [], 
            "ManagesVpcEndpoints": false, 
            "PrivateDnsName": "api.ecr.us-west-2.amazonaws.com", 
            "ServiceName": "com.amazonaws.us-west-2.ecr.api", 
            "VpcEndpointPolicySupported": true, 
            "ServiceId": "vpce-svc-07d1f428f072fd172", 
            "Owner": "amazon", 
            "AvailabilityZones": [
                "us-west-2a", 
                "us-west-2b", 
                "us-west-2c", 
                "us-west-2d"
            ], 
            "AcceptanceRequired": false, 
            "BaseEndpointDnsNames": [
                "api.ecr.us-west-2.vpce.amazonaws.com"
            ]
        }, 
        {
            "ServiceType": [
                {
                    "ServiceType": "Interface"
                }
            ], 
            "Tags": [], 
            "ManagesVpcEndpoints": false, 
            "PrivateDnsName": "*.dkr.ecr.us-west-2.amazonaws.com", 
            "ServiceName": "com.amazonaws.us-west-2.ecr.dkr", 
            "VpcEndpointPolicySupported": true, 
            "ServiceId": "vpce-svc-09d74a28015a69002", 
            "Owner": "amazon", 
            "AvailabilityZones": [
                "us-west-2a", 
                "us-west-2b", 
                "us-west-2c", 
                "us-west-2d"
            ], 
            "AcceptanceRequired": false, 
            "BaseEndpointDnsNames": [
                "dkr.ecr.us-west-2.vpce.amazonaws.com"
            ]
        }, 
        {
            "ServiceType": [
                {
                    "ServiceType": "Gateway"
                }
            ], 
            "Tags": [], 
            "ManagesVpcEndpoints": false, 
            "AcceptanceRequired": false, 
            "ServiceName": "com.amazonaws.us-west-2.s3", 
            "VpcEndpointPolicySupported": true, 
            "ServiceId": "vpce-svc-0001be97e1865c74e", 
            "Owner": "amazon", 
            "AvailabilityZones": [
                "us-west-2a", 
                "us-west-2b", 
                "us-west-2c", 
                "us-west-2d"
            ], 
            "BaseEndpointDnsNames": [
                "s3.us-west-2.amazonaws.com"
            ]
        }, 
        {
            "ServiceType": [
                {
                    "ServiceType": "Interface"
                }
            ], 
            "Tags": [], 
            "ManagesVpcEndpoints": false, 
            "PrivateDnsName": "sts.us-west-2.amazonaws.com", 
            "ServiceName": "com.amazonaws.us-west-2.sts", 
            "VpcEndpointPolicySupported": true, 
            "ServiceId": "vpce-svc-06681ce20e9a3e8c4", 
            "Owner": "amazon", 
            "AvailabilityZones": [
                "us-west-2a", 
                "us-west-2b", 
                "us-west-2c"
            ], 
            "AcceptanceRequired": false, 
            "BaseEndpointDnsNames": [
                "sts.us-west-2.vpce.amazonaws.com"
            ]
        }
    ]
}
`

	var output *ec2.DescribeVpcEndpointServicesOutput
	Expect(json.Unmarshal([]byte(serviceDetailsJSON), &output)).To(Succeed())

	provider.MockEC2().On("DescribeVpcEndpointServices", mock.MatchedBy(func(e *ec2.DescribeVpcEndpointServicesInput) bool {
		return len(e.ServiceNames) == 5
	})).Return(output, nil)
}

func mockDescribeRouteTables(provider *mockprovider.MockProvider, subnetIDs []string) {
	output := &ec2.DescribeRouteTablesOutput{
		RouteTables: make([]*ec2.RouteTable, len(subnetIDs)),
	}

	for i, subnetID := range subnetIDs {
		rtID := aws.String(fmt.Sprintf("rtb-custom-%d", i+1))
		output.RouteTables[i] = &ec2.RouteTable{
			VpcId:        aws.String("vpc-custom"),
			RouteTableId: rtID,
			Associations: []*ec2.RouteTableAssociation{
				{
					RouteTableId:            rtID,
					SubnetId:                aws.String(subnetID),
					RouteTableAssociationId: aws.String("rtbassoc-custom"),
					Main:                    aws.Bool(false),
				},
			},
		}
	}

	provider.MockEC2().On("DescribeRouteTables", mock.MatchedBy(func(input *ec2.DescribeRouteTablesInput) bool {
		return len(input.Filters) > 0
	})).Return(output, nil)
}
