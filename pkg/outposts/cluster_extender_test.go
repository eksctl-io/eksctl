package outposts_test

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	awsoutposts "github.com/aws/aws-sdk-go-v2/service/outposts"
	outpoststypes "github.com/aws/aws-sdk-go-v2/service/outposts/types"

	"github.com/stretchr/testify/mock"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/outposts"
	"github.com/weaveworks/eksctl/pkg/outposts/fakes"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
	"github.com/weaveworks/eksctl/pkg/utils/ipnet"
)

type extendedClusterEntry struct {
	cluster         outposts.ClusterToExtend
	clusterVPC      *api.ClusterVPC
	mockProvider    func(*mockprovider.MockProvider)
	externalSubnets []ec2types.Subnet

	expectedSubnets                                *api.ClusterSubnets
	expectedAppendNewClusterStackResourceCallCount int
	expectedErr                                    string
}

func newFakeCluster(controlPlaneOnOutposts bool, nodeGroupOutpostARN string) outposts.ClusterToExtend {
	return &fakes.FakeClusterToExtend{
		IsControlPlaneOnOutpostsStub: func() bool {
			return controlPlaneOnOutposts
		},
		FindNodeGroupOutpostARNStub: func() (string, bool) {
			if nodeGroupOutpostARN != "" {
				return nodeGroupOutpostARN, true
			}
			return "", false
		},
	}
}

var _ = Describe("Cluster Extender", func() {
	const outpostARN = "arn:aws:outposts:us-west-2:1234:outpost/op-1234"

	DescribeTable("Extend cluster with Outpost subnets", func(e extendedClusterEntry) {
		provider := mockprovider.NewMockProvider()
		if e.mockProvider != nil {
			e.mockProvider(provider)
		}
		stackUpdater := &fakes.FakeStackUpdater{
			AppendNewClusterStackResourceStub: func(_ context.Context, _ bool, _ bool) (bool, error) {
				return true, nil
			},
		}

		clusterExtender := &outposts.ClusterExtender{
			StackUpdater: stackUpdater,
			EC2API:       provider.EC2(),
			OutpostsAPI:  provider.Outposts(),
		}
		if e.clusterVPC != nil {
			mockDescribeSubnets(provider, e.clusterVPC.Subnets, e.externalSubnets)
		}

		err := clusterExtender.ExtendWithOutpostSubnetsIfRequired(context.Background(), e.cluster, e.clusterVPC)
		if e.expectedErr != "" {
			Expect(err).To(MatchError(ContainSubstring(e.expectedErr)))
			return
		}

		Expect(err).NotTo(HaveOccurred())
		if e.clusterVPC != nil {
			Expect(e.clusterVPC.Subnets).To(Equal(e.expectedSubnets))
		}
		Expect(stackUpdater.AppendNewClusterStackResourceCallCount()).To(Equal(e.expectedAppendNewClusterStackResourceCallCount))
	},
		Entry("control plane on Outposts should result in a no-op", extendedClusterEntry{
			cluster: newFakeCluster(true, ""),
		}),

		Entry("nodegroup not on Outposts should result in a no-op", extendedClusterEntry{
			cluster: newFakeCluster(false, ""),
		}),

		Entry("nodegroup on a different Outpost than subnets", extendedClusterEntry{
			cluster: newFakeCluster(false, "arn:aws:outposts:us-west-2:1234:outpost/op-5678"),

			clusterVPC: &api.ClusterVPC{
				Network: api.Network{
					ID:   "vpc-1234",
					CIDR: ipnet.MustParseCIDR("192.168.0.0/16"),
				},
				Subnets: &api.ClusterSubnets{
					Public: api.AZSubnetMapping{
						"us-west-2a": api.AZSubnetSpec{
							ID:   "subnet-1",
							AZ:   "us-west-2a",
							CIDR: ipnet.MustParseCIDR("192.168.0.0/19"),
						},
						"us-west-2b": api.AZSubnetSpec{
							ID:   "subnet-2",
							AZ:   "us-west-2b",
							CIDR: ipnet.MustParseCIDR("192.168.32.0/19"),
						},
						"outpost-us-west-2o-1": api.AZSubnetSpec{
							AZ:         "us-west-2o",
							CIDR:       ipnet.MustParseCIDR("192.168.128.0/19"),
							OutpostARN: outpostARN,
						},
					},
					Private: api.AZSubnetMapping{
						"us-west-2a": api.AZSubnetSpec{
							ID:   "subnet-3",
							AZ:   "us-west-2a",
							CIDR: ipnet.MustParseCIDR("192.168.64.0/19"),
						},
						"us-west-2b": api.AZSubnetSpec{
							ID:   "subnet-4",
							AZ:   "us-west-2b",
							CIDR: ipnet.MustParseCIDR("192.168.96.0/19"),
						},
						"outpost-us-west-2o-1": api.AZSubnetSpec{
							AZ:         "us-west-2o",
							CIDR:       ipnet.MustParseCIDR("192.168.160.0/19"),
							OutpostARN: outpostARN,
						},
					},
				},
			},

			expectedErr: fmt.Sprintf(`cannot extend a cluster with two different Outposts; found subnets on Outpost %q but nodegroup is using %q`, outpostARN, "arn:aws:outposts:us-west-2:1234:outpost/op-5678"),
		}),

		Entry("subnets already exist on Outposts", extendedClusterEntry{
			cluster: newFakeCluster(false, outpostARN),

			clusterVPC: &api.ClusterVPC{
				Network: api.Network{
					ID:   "vpc-1234",
					CIDR: ipnet.MustParseCIDR("192.168.0.0/16"),
				},
				Subnets: &api.ClusterSubnets{
					Public: api.AZSubnetMapping{
						"us-west-2a": api.AZSubnetSpec{
							ID:   "subnet-1",
							AZ:   "us-west-2a",
							CIDR: ipnet.MustParseCIDR("192.168.0.0/19"),
						},
						"us-west-2b": api.AZSubnetSpec{
							ID:   "subnet-2",
							AZ:   "us-west-2b",
							CIDR: ipnet.MustParseCIDR("192.168.32.0/19"),
						},
						"outpost-us-west-2o-1": api.AZSubnetSpec{
							AZ:         "us-west-2o",
							CIDR:       ipnet.MustParseCIDR("192.168.128.0/19"),
							OutpostARN: outpostARN,
						},
					},
					Private: api.AZSubnetMapping{
						"us-west-2a": api.AZSubnetSpec{
							ID:   "subnet-3",
							AZ:   "us-west-2a",
							CIDR: ipnet.MustParseCIDR("192.168.64.0/19"),
						},
						"us-west-2b": api.AZSubnetSpec{
							ID:   "subnet-4",
							AZ:   "us-west-2b",
							CIDR: ipnet.MustParseCIDR("192.168.96.0/19"),
						},
						"outpost-us-west-2o-1": api.AZSubnetSpec{
							AZ:         "us-west-2o",
							CIDR:       ipnet.MustParseCIDR("192.168.160.0/19"),
							OutpostARN: outpostARN,
						},
					},
				},
			},

			expectedSubnets: &api.ClusterSubnets{
				Public: api.AZSubnetMapping{
					"us-west-2a": api.AZSubnetSpec{
						ID:   "subnet-1",
						AZ:   "us-west-2a",
						CIDR: ipnet.MustParseCIDR("192.168.0.0/19"),
					},
					"us-west-2b": api.AZSubnetSpec{
						ID:   "subnet-2",
						AZ:   "us-west-2b",
						CIDR: ipnet.MustParseCIDR("192.168.32.0/19"),
					},
					"outpost-us-west-2o-1": api.AZSubnetSpec{
						AZ:         "us-west-2o",
						CIDR:       ipnet.MustParseCIDR("192.168.128.0/19"),
						OutpostARN: outpostARN,
					},
				},
				Private: api.AZSubnetMapping{
					"us-west-2a": api.AZSubnetSpec{
						ID:   "subnet-3",
						AZ:   "us-west-2a",
						CIDR: ipnet.MustParseCIDR("192.168.64.0/19"),
					},
					"us-west-2b": api.AZSubnetSpec{
						ID:   "subnet-4",
						AZ:   "us-west-2b",
						CIDR: ipnet.MustParseCIDR("192.168.96.0/19"),
					},
					"outpost-us-west-2o-1": api.AZSubnetSpec{
						AZ:         "us-west-2o",
						CIDR:       ipnet.MustParseCIDR("192.168.160.0/19"),
						OutpostARN: outpostARN,
					},
				},
			},
		}),

		Entry("VPC CIDR block exhausted", extendedClusterEntry{
			cluster: newFakeCluster(false, outpostARN),

			clusterVPC: &api.ClusterVPC{
				Network: api.Network{
					ID:   "vpc-1234",
					CIDR: ipnet.MustParseCIDR("192.168.0.0/16"),
				},
				Subnets: &api.ClusterSubnets{
					Public: api.AZSubnetMapping{
						"us-west-2a": api.AZSubnetSpec{
							ID:   "subnet-1",
							AZ:   "us-west-2a",
							CIDR: ipnet.MustParseCIDR("192.168.0.0/20"),
						},
						"us-west-2b": api.AZSubnetSpec{
							ID:   "subnet-2",
							AZ:   "us-west-2b",
							CIDR: ipnet.MustParseCIDR("192.168.16.0/20"),
						},
						"us-west-2c": api.AZSubnetSpec{
							ID:   "subnet-3",
							AZ:   "us-west-2b",
							CIDR: ipnet.MustParseCIDR("192.168.32.0/20"),
						},
						"us-west-2d": api.AZSubnetSpec{
							ID:   "subnet-4",
							AZ:   "us-west-2b",
							CIDR: ipnet.MustParseCIDR("192.168.48.0/20"),
						},
						"us-west-2e": api.AZSubnetSpec{
							ID:   "subnet-5",
							AZ:   "us-west-2b",
							CIDR: ipnet.MustParseCIDR("192.168.64.0/20"),
						},
						"us-west-2f": api.AZSubnetSpec{
							ID:   "subnet-6",
							AZ:   "us-west-2f",
							CIDR: ipnet.MustParseCIDR("192.168.80.0/20"),
						},
						"us-west-2g": api.AZSubnetSpec{
							ID:   "subnet-7",
							AZ:   "us-west-2g",
							CIDR: ipnet.MustParseCIDR("192.168.96.0/20"),
						},
						"us-west-2h": api.AZSubnetSpec{
							ID:   "subnet-8",
							AZ:   "us-west-2g",
							CIDR: ipnet.MustParseCIDR("192.168.112.0/20"),
						},
					},
					Private: api.AZSubnetMapping{
						"us-west-2a": api.AZSubnetSpec{
							ID:   "subnet-9",
							AZ:   "us-west-2a",
							CIDR: ipnet.MustParseCIDR("192.168.128.0/20"),
						},
						"us-west-2b": api.AZSubnetSpec{
							ID:   "subnet-10",
							AZ:   "us-west-2b",
							CIDR: ipnet.MustParseCIDR("192.168.144.0/20"),
						},
						"us-west-2c": api.AZSubnetSpec{
							ID:   "subnet-11",
							AZ:   "us-west-2c",
							CIDR: ipnet.MustParseCIDR("192.168.160.0/20"),
						},
						"us-west-2d": api.AZSubnetSpec{
							ID:   "subnet-12",
							AZ:   "us-west-2d",
							CIDR: ipnet.MustParseCIDR("192.168.176.0/20"),
						},
						"us-west-2e": api.AZSubnetSpec{
							ID:   "subnet-13",
							AZ:   "us-west-2e",
							CIDR: ipnet.MustParseCIDR("192.168.192.0/20"),
						},
						"us-west-2f": api.AZSubnetSpec{
							ID:   "subnet-14",
							AZ:   "us-west-2f",
							CIDR: ipnet.MustParseCIDR("192.168.208.0/20"),
						},
						"us-west-2g": api.AZSubnetSpec{
							ID:   "subnet-15",
							AZ:   "us-west-2g",
							CIDR: ipnet.MustParseCIDR("192.168.224.0/20"),
						},
						"us-west-2h": api.AZSubnetSpec{
							ID:   "subnet-16",
							AZ:   "us-west-2h",
							CIDR: ipnet.MustParseCIDR("192.168.240.0/20"),
						},
					},
				},
			},
			mockProvider: func(provider *mockprovider.MockProvider) {
				mockOutposts(provider, outpostARN)
			},
			expectedErr: "VPC cannot be extended with more subnets: expected to find at least two free CIDRs in VPC; got 0",
		}),

		Entry("external subnet exists that overlaps with new CIDRs", extendedClusterEntry{
			cluster: newFakeCluster(false, outpostARN),

			clusterVPC: &api.ClusterVPC{
				Network: api.Network{
					ID:   "vpc-1234",
					CIDR: ipnet.MustParseCIDR("192.168.0.0/16"),
				},
				Subnets: &api.ClusterSubnets{
					Public: api.AZSubnetMapping{
						"us-west-2a": api.AZSubnetSpec{
							ID:   "subnet-1",
							AZ:   "us-west-2a",
							CIDR: ipnet.MustParseCIDR("192.168.0.0/19"),
						},
						"us-west-2b": api.AZSubnetSpec{
							ID:   "subnet-2",
							AZ:   "us-west-2b",
							CIDR: ipnet.MustParseCIDR("192.168.32.0/19"),
						},
					},
					Private: api.AZSubnetMapping{
						"us-west-2a": api.AZSubnetSpec{
							ID:   "subnet-3",
							AZ:   "us-west-2a",
							CIDR: ipnet.MustParseCIDR("192.168.64.0/19"),
						},
						"us-west-2b": api.AZSubnetSpec{
							ID:   "subnet-4",
							AZ:   "us-west-2b",
							CIDR: ipnet.MustParseCIDR("192.168.96.0/19"),
						},
					},
				},
			},

			mockProvider: func(provider *mockprovider.MockProvider) {
				provider.MockOutposts().On("GetOutpost", mock.Anything, &awsoutposts.GetOutpostInput{
					OutpostId: aws.String(outpostARN),
				}).Return(&awsoutposts.GetOutpostOutput{
					Outpost: &outpoststypes.Outpost{
						OutpostId:        aws.String(outpostARN),
						AvailabilityZone: aws.String("us-west-2o"),
					},
				}, nil)
			},

			externalSubnets: []ec2types.Subnet{
				{
					SubnetId:         aws.String("subnet-external"),
					AvailabilityZone: aws.String("us-west-2a"),
					CidrBlock:        aws.String("192.168.160.0/19"),
				},
			},

			expectedErr: `cannot create subnets on Outpost; subnet CIDR "192.168.160.0/19" (ID: subnet-external) created outside of eksctl overlaps with new CIDR "192.168.160.0/19"`,
		}),

		Entry("VPC is updated with new CIDRs for Outpost subnets for /19 block", extendedClusterEntry{
			cluster: &fakes.FakeClusterToExtend{
				IsControlPlaneOnOutpostsStub: func() bool {
					return false
				},
				FindNodeGroupOutpostARNStub: func() (string, bool) {
					return outpostARN, true
				},
			},

			clusterVPC: &api.ClusterVPC{
				Network: api.Network{
					ID:   "vpc-1234",
					CIDR: ipnet.MustParseCIDR("192.168.0.0/16"),
				},
				Subnets: &api.ClusterSubnets{
					Public: api.AZSubnetMapping{
						"us-west-2a": api.AZSubnetSpec{
							ID:   "subnet-1",
							AZ:   "us-west-2a",
							CIDR: ipnet.MustParseCIDR("192.168.0.0/19"),
						},
						"us-west-2b": api.AZSubnetSpec{
							ID:   "subnet-2",
							AZ:   "us-west-2b",
							CIDR: ipnet.MustParseCIDR("192.168.32.0/19"),
						},
					},
					Private: api.AZSubnetMapping{
						"us-west-2a": api.AZSubnetSpec{
							ID:   "subnet-3",
							AZ:   "us-west-2a",
							CIDR: ipnet.MustParseCIDR("192.168.64.0/19"),
						},
						"us-west-2b": api.AZSubnetSpec{
							ID:   "subnet-4",
							AZ:   "us-west-2b",
							CIDR: ipnet.MustParseCIDR("192.168.96.0/19"),
						},
					},
				},
			},

			mockProvider: func(provider *mockprovider.MockProvider) {
				provider.MockOutposts().On("GetOutpost", mock.Anything, &awsoutposts.GetOutpostInput{
					OutpostId: aws.String(outpostARN),
				}).Return(&awsoutposts.GetOutpostOutput{
					Outpost: &outpoststypes.Outpost{
						OutpostId:        aws.String(outpostARN),
						AvailabilityZone: aws.String("us-west-2o"),
					},
				}, nil)
			},

			expectedSubnets: &api.ClusterSubnets{
				Public: api.AZSubnetMapping{
					"us-west-2a": api.AZSubnetSpec{
						ID:   "subnet-1",
						AZ:   "us-west-2a",
						CIDR: ipnet.MustParseCIDR("192.168.0.0/19"),
					},
					"us-west-2b": api.AZSubnetSpec{
						ID:   "subnet-2",
						AZ:   "us-west-2b",
						CIDR: ipnet.MustParseCIDR("192.168.32.0/19"),
					},
					"outpost-us-west-2o-1": api.AZSubnetSpec{
						AZ:         "us-west-2o",
						CIDR:       ipnet.MustParseCIDR("192.168.128.0/19"),
						CIDRIndex:  5,
						OutpostARN: outpostARN,
					},
				},
				Private: api.AZSubnetMapping{
					"us-west-2a": api.AZSubnetSpec{
						ID:   "subnet-3",
						AZ:   "us-west-2a",
						CIDR: ipnet.MustParseCIDR("192.168.64.0/19"),
					},
					"us-west-2b": api.AZSubnetSpec{
						ID:   "subnet-4",
						AZ:   "us-west-2b",
						CIDR: ipnet.MustParseCIDR("192.168.96.0/19"),
					},
					"outpost-us-west-2o-1": api.AZSubnetSpec{
						AZ:         "us-west-2o",
						CIDR:       ipnet.MustParseCIDR("192.168.160.0/19"),
						CIDRIndex:  6,
						OutpostARN: outpostARN,
					},
				},
			},

			expectedAppendNewClusterStackResourceCallCount: 1,
		}),

		Entry("VPC is updated with new CIDRs for Outpost subnets for /20 block", extendedClusterEntry{
			cluster: newFakeCluster(false, outpostARN),

			clusterVPC: &api.ClusterVPC{
				Network: api.Network{
					ID:   "vpc-1234",
					CIDR: ipnet.MustParseCIDR("192.168.0.0/16"),
				},
				Subnets: &api.ClusterSubnets{
					Public: api.AZSubnetMapping{
						"us-west-2a": api.AZSubnetSpec{
							ID:   "subnet-1",
							AZ:   "us-west-2a",
							CIDR: ipnet.MustParseCIDR("192.168.0.0/20"),
						},
						"us-west-2b": api.AZSubnetSpec{
							ID:   "subnet-2",
							AZ:   "us-west-2b",
							CIDR: ipnet.MustParseCIDR("192.168.16.0/20"),
						},
						"us-west-2c": api.AZSubnetSpec{
							ID:   "subnet-3",
							AZ:   "us-west-2c",
							CIDR: ipnet.MustParseCIDR("192.168.32.0/20"),
						},
						"us-west-2d": api.AZSubnetSpec{
							ID:   "subnet-4",
							AZ:   "us-west-2d",
							CIDR: ipnet.MustParseCIDR("192.168.48.0/20"),
						},
						"us-west-2e": api.AZSubnetSpec{
							ID:   "subnet-5",
							AZ:   "us-west-2e",
							CIDR: ipnet.MustParseCIDR("192.168.64.0/20"),
						},
					},
					Private: api.AZSubnetMapping{
						"us-west-2a": api.AZSubnetSpec{
							ID:   "subnet-6",
							AZ:   "us-west-2a",
							CIDR: ipnet.MustParseCIDR("192.168.80.0/20"),
						},
						"us-west-2b": api.AZSubnetSpec{
							ID:   "subnet-7",
							AZ:   "us-west-2b",
							CIDR: ipnet.MustParseCIDR("192.168.96.0/20"),
						},
						"us-west-2c": api.AZSubnetSpec{
							ID:   "subnet-8",
							AZ:   "us-west-2c",
							CIDR: ipnet.MustParseCIDR("192.168.112.0/20"),
						},
						"us-west-2d": api.AZSubnetSpec{
							ID:   "subnet-9",
							AZ:   "us-west-2d",
							CIDR: ipnet.MustParseCIDR("192.168.128.0/20"),
						},
						"us-west-2e": api.AZSubnetSpec{
							ID:   "subnet-10",
							AZ:   "us-west-2e",
							CIDR: ipnet.MustParseCIDR("192.168.144.0/20"),
						},
					},
				},
			},

			mockProvider: func(provider *mockprovider.MockProvider) {
				mockOutposts(provider, outpostARN)
			},

			expectedSubnets: &api.ClusterSubnets{
				Public: api.AZSubnetMapping{
					"us-west-2a": api.AZSubnetSpec{
						ID:   "subnet-1",
						AZ:   "us-west-2a",
						CIDR: ipnet.MustParseCIDR("192.168.0.0/20"),
					},
					"us-west-2b": api.AZSubnetSpec{
						ID:   "subnet-2",
						AZ:   "us-west-2b",
						CIDR: ipnet.MustParseCIDR("192.168.16.0/20"),
					},
					"us-west-2c": api.AZSubnetSpec{
						ID:   "subnet-3",
						AZ:   "us-west-2c",
						CIDR: ipnet.MustParseCIDR("192.168.32.0/20"),
					},
					"us-west-2d": api.AZSubnetSpec{
						ID:   "subnet-4",
						AZ:   "us-west-2d",
						CIDR: ipnet.MustParseCIDR("192.168.48.0/20"),
					},
					"us-west-2e": api.AZSubnetSpec{
						ID:   "subnet-5",
						AZ:   "us-west-2e",
						CIDR: ipnet.MustParseCIDR("192.168.64.0/20"),
					},
					"outpost-us-west-2o-1": api.AZSubnetSpec{
						AZ:         "us-west-2o",
						CIDR:       ipnet.MustParseCIDR("192.168.160.0/20"),
						CIDRIndex:  11,
						OutpostARN: outpostARN,
					},
				},
				Private: api.AZSubnetMapping{
					"us-west-2a": api.AZSubnetSpec{
						ID:   "subnet-6",
						AZ:   "us-west-2a",
						CIDR: ipnet.MustParseCIDR("192.168.80.0/20"),
					},
					"us-west-2b": api.AZSubnetSpec{
						ID:   "subnet-7",
						AZ:   "us-west-2b",
						CIDR: ipnet.MustParseCIDR("192.168.96.0/20"),
					},
					"us-west-2c": api.AZSubnetSpec{
						ID:   "subnet-8",
						AZ:   "us-west-2c",
						CIDR: ipnet.MustParseCIDR("192.168.112.0/20"),
					},
					"us-west-2d": api.AZSubnetSpec{
						ID:   "subnet-9",
						AZ:   "us-west-2d",
						CIDR: ipnet.MustParseCIDR("192.168.128.0/20"),
					},
					"us-west-2e": api.AZSubnetSpec{
						ID:   "subnet-10",
						AZ:   "us-west-2e",
						CIDR: ipnet.MustParseCIDR("192.168.144.0/20"),
					},
					"outpost-us-west-2o-1": api.AZSubnetSpec{
						AZ:         "us-west-2o",
						CIDR:       ipnet.MustParseCIDR("192.168.176.0/20"),
						CIDRIndex:  12,
						OutpostARN: outpostARN,
					},
				},
			},

			expectedAppendNewClusterStackResourceCallCount: 1,
		}),

		Entry("external subnet exists that overlaps in a VPC with /20 subnets", extendedClusterEntry{
			cluster: newFakeCluster(false, outpostARN),

			clusterVPC: &api.ClusterVPC{
				Network: api.Network{
					ID:   "vpc-1234",
					CIDR: ipnet.MustParseCIDR("192.168.0.0/16"),
				},
				Subnets: &api.ClusterSubnets{
					Public: api.AZSubnetMapping{
						"us-west-2a": api.AZSubnetSpec{
							ID:   "subnet-1",
							AZ:   "us-west-2a",
							CIDR: ipnet.MustParseCIDR("192.168.0.0/20"),
						},
						"us-west-2b": api.AZSubnetSpec{
							ID:   "subnet-2",
							AZ:   "us-west-2b",
							CIDR: ipnet.MustParseCIDR("192.168.16.0/20"),
						},
						"us-west-2c": api.AZSubnetSpec{
							ID:   "subnet-3",
							AZ:   "us-west-2b",
							CIDR: ipnet.MustParseCIDR("192.168.32.0/20"),
						},
						"us-west-2d": api.AZSubnetSpec{
							ID:   "subnet-4",
							AZ:   "us-west-2d",
							CIDR: ipnet.MustParseCIDR("192.168.48.0/20"),
						},
						"us-west-2e": api.AZSubnetSpec{
							ID:   "subnet-5",
							AZ:   "us-west-2b",
							CIDR: ipnet.MustParseCIDR("192.168.64.0/20"),
						},
					},
					Private: api.AZSubnetMapping{
						"us-west-2a": api.AZSubnetSpec{
							ID:   "subnet-6",
							AZ:   "us-west-2a",
							CIDR: ipnet.MustParseCIDR("192.168.80.0/20"),
						},
						"us-west-2b": api.AZSubnetSpec{
							ID:   "subnet-7",
							AZ:   "us-west-2b",
							CIDR: ipnet.MustParseCIDR("192.168.96.0/20"),
						},
						"us-west-2c": api.AZSubnetSpec{
							ID:   "subnet-8",
							AZ:   "us-west-2c",
							CIDR: ipnet.MustParseCIDR("192.168.112.0/20"),
						},
						"us-west-2d": api.AZSubnetSpec{
							ID:   "subnet-9",
							AZ:   "us-west-2d",
							CIDR: ipnet.MustParseCIDR("192.168.128.0/20"),
						},
						"us-west-2e": api.AZSubnetSpec{
							ID:   "subnet-10",
							AZ:   "us-west-2e",
							CIDR: ipnet.MustParseCIDR("192.168.144.0/20"),
						},
					},
				},
			},

			mockProvider: func(provider *mockprovider.MockProvider) {
				mockOutposts(provider, outpostARN)
			},

			externalSubnets: []ec2types.Subnet{
				{
					SubnetId:         aws.String("subnet-external"),
					AvailabilityZone: aws.String("us-west-2b"),
					CidrBlock:        aws.String("192.168.168.0/22"),
				},
			},

			expectedErr: `cannot create subnets on Outpost; subnet CIDR "192.168.168.0/22" (ID: subnet-external) created outside of eksctl overlaps with new CIDR "192.168.160.0/20"`,
		}),
	)
})

func mockDescribeSubnets(provider *mockprovider.MockProvider, clusterSubnets *api.ClusterSubnets, externalSubnets []ec2types.Subnet) {
	var allSubnets []ec2types.Subnet

	addSubnets := func(subnets api.AZSubnetMapping) {
		for _, s := range subnets {
			allSubnets = append(allSubnets, ec2types.Subnet{
				SubnetId:         aws.String(s.ID),
				CidrBlock:        aws.String(s.CIDR.String()),
				AvailabilityZone: aws.String(s.AZ),
			})
		}
	}

	addSubnets(clusterSubnets.Public)
	addSubnets(clusterSubnets.Private)

	allSubnets = append(allSubnets, externalSubnets...)

	provider.MockEC2().On("DescribeSubnets", mock.Anything, &ec2.DescribeSubnetsInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{"vpc-1234"},
			},
		},
	}, mock.Anything).Return(&ec2.DescribeSubnetsOutput{
		Subnets: allSubnets,
	}, nil)
}

func mockOutposts(provider *mockprovider.MockProvider, outpostARN string) {
	provider.MockOutposts().On("GetOutpost", mock.Anything, &awsoutposts.GetOutpostInput{
		OutpostId: aws.String(outpostARN),
	}).Return(&awsoutposts.GetOutpostOutput{
		Outpost: &outpoststypes.Outpost{
			OutpostId:        aws.String(outpostARN),
			AvailabilityZone: aws.String("us-west-2o"),
		},
	}, nil)
}

// TODO: add tests for find nodegroup clusterconfig outpost etc.
