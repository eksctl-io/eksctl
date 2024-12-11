package builder_test

import (
	"net"
	"testing"

	. "github.com/onsi/gomega"

	gfnt "goformation/v4/cloudformation/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/utils/ipnet"
)

func TestCfnBuilder(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var (
	azA, azB, azC                            = "us-west-2a", "us-west-2b", "us-west-2c"
	azAFormatted, azBFormatted, azCFormatted = "USWEST2A", "USWEST2B", "USWEST2C"
	privateSubnet1, privateSubnet2           = "subnet-0ade11bad78dced9f", "subnet-0f98135715dfcf55a"
	publicSubnet1, publicSubnet2             = "subnet-0ade11bad78dced9e", "subnet-0f98135715dfcf55f"
	privateSubnetRef1, privateSubnetRef2     = "SubnetPrivateUSWEST2A", "SubnetPrivateUSWEST2B"
	publicSubnetRef1, publicSubnetRef2       = "SubnetPublicUSWEST2A", "SubnetPublicUSWEST2B"
	vpcResourceKey, igwKey, gaKey            = "VPC", "InternetGateway", "VPCGatewayAttachment"
	pubRouteTable, pubSubnetRoute            = "PublicRouteTable", "PublicSubnetRoute"
	privRouteTableA, privRouteTableB         = "PrivateRouteTableUSWEST2A", "PrivateRouteTableUSWEST2B"
	rtaPublicA, rtaPublicB                   = "RouteTableAssociationPublicUSWEST2A", "RouteTableAssociationPublicUSWEST2B"
	rtaPrivateA, rtaPrivateB                 = "RouteTableAssociationPrivateUSWEST2A", "RouteTableAssociationPrivateUSWEST2B"
)

func vpcConfig() *api.ClusterVPC {
	disable := api.ClusterDisableNAT
	return &api.ClusterVPC{
		NAT: &api.ClusterNAT{
			Gateway: &disable,
		},
		ClusterEndpoints: api.ClusterEndpointAccessDefaults(),
		Subnets: &api.ClusterSubnets{
			Public: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
				azB: {
					ID: publicSubnet2,
					CIDR: &ipnet.IPNet{
						IPNet: net.IPNet{
							IP:   []byte{192, 168, 0, 0},
							Mask: []byte{255, 255, 224, 0},
						},
					},
				},
				azA: {
					ID: publicSubnet1,
					CIDR: &ipnet.IPNet{
						IPNet: net.IPNet{
							IP:   []byte{192, 168, 32, 0},
							Mask: []byte{255, 255, 224, 0},
						},
					},
				},
			}),
			Private: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
				azB: {
					ID: privateSubnet2,
					CIDR: &ipnet.IPNet{
						IPNet: net.IPNet{
							IP:   []byte{192, 168, 96, 0},
							Mask: []byte{255, 255, 224, 0},
						},
					},
				},
				azA: {
					ID: privateSubnet1,
					CIDR: &ipnet.IPNet{
						IPNet: net.IPNet{
							IP:   []byte{192, 168, 128, 0},
							Mask: []byte{255, 255, 224, 0},
						},
					},
				},
			}),
		},
	}
}

func isRefTo(obj interface{}, value string) bool {
	Expect(obj).NotTo(BeEmpty())
	o, ok := obj.(map[string]interface{})
	Expect(ok).To(BeTrue())
	Expect(o).To(HaveKey(gfnt.Ref))
	return o[gfnt.Ref] == value
}
