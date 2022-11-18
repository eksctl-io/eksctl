package create

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cftypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go-v2/service/outposts"
	outpoststypes "github.com/aws/aws-sdk-go-v2/service/outposts/types"
	"github.com/aws/aws-sdk-go/aws/credentials"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	kubefake "k8s.io/client-go/kubernetes/fake"
	restclient "k8s.io/client-go/rest"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils/filter"
	"github.com/weaveworks/eksctl/pkg/ctl/ctltest"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/eks/fakes"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("create cluster", func() {
	Describe("un-managed node group", func() {
		It("understands ssh access arguments correctly", func() {
			commandArgs := []string{"cluster", "--managed=false", "--ssh-access=false", "--ssh-public-key=dummy-key"}
			cmd := newMockEmptyCmd(commandArgs...)
			count := 0
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				createClusterCmdWithRunFunc(cmd, func(cmd *cmdutils.Cmd, ngFilter *filter.NodeGroupFilter, params *cmdutils.CreateClusterCmdParams) error {
					Expect(*cmd.ClusterConfig.NodeGroups[0].SSH.Allow).To(BeFalse())
					count++
					return nil
				})
			})
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
		})
		DescribeTable("create cluster successfully",
			func(args ...string) {
				commandArgs := append([]string{"cluster"}, args...)
				cmd := newMockEmptyCmd(commandArgs...)
				count := 0
				cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
					createClusterCmdWithRunFunc(cmd, func(cmd *cmdutils.Cmd, ngFilter *filter.NodeGroupFilter, params *cmdutils.CreateClusterCmdParams) error {
						Expect(cmd.ClusterConfig.Metadata.Name).NotTo(BeNil())
						count++
						return nil
					})
				})
				_, err := cmd.execute()
				Expect(err).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			},
			Entry("without cluster name", ""),
			Entry("with cluster name as flag", "--name", "clusterName"),
			Entry("with cluster name as argument", "clusterName"),
			Entry("with cluster name with hyphen as flag", "--name", "my-cluster-name-is-fine10"),
			Entry("with cluster name with hyphen as argument", "my-Cluster-name-is-fine10"),
			// vpc networking flags
			Entry("with vpc-cidr flag", "--vpc-cidr", "10.0.0.0/20"),
			Entry("with vpc-private-subnets flag", "--vpc-private-subnets", "10.0.0.0/24"),
			Entry("with vpc-public-subnets flag", "--vpc-public-subnets", "10.0.0.0/24"),
			Entry("with vpc-from-kops-cluster flag", "--vpc-from-kops-cluster", "dummy-kops-cluster"),
			Entry("with vpc-nat-mode flag", "--vpc-nat-mode", "Single"),
			// kubeconfig flags
			Entry("with write-kubeconfig flag", "--write-kubeconfig"),
			Entry("with kubeconfig flag", "--kubeconfig", "~/.kube"),
			Entry("with authenticator-role-arn flag", "--authenticator-role-arn", "arn::dummy::123/role"),
			Entry("with auto-kubeconfig flag", "--auto-kubeconfig"),
			// common node group flags
			Entry("with node-type flag", "--node-type", "m5.large"),
			Entry("with nodes flag", "--nodes", "2"),
			Entry("with nodes-min flag", "--nodes-min", "2"),
			Entry("with nodes-max flag", "--nodes-max", "2"),
			Entry("with node-volume-size flag", "--node-volume-size", "2"),
			Entry("with node-volume-type flag", "--node-volume-type", "gp2"),
			Entry("with max-pods-per-node flag", "--max-pods-per-node", "20"),
			Entry("with ssh-access flag", "--ssh-access", "true"),
			Entry("with ssh-public-key flag", "--ssh-public-key", "dummy-public-key"),
			Entry("with enable-ssm flag", "--enable-ssm"),
			Entry("with node-ami flag", "--node-ami", "ami-dummy-123"),
			Entry("with node-ami-family flag", "--node-ami-family", "AmazonLinux2"),
			Entry("with node-private-networking flag", "--node-private-networking", "true"),
			Entry("with node-security-groups flag", "--node-security-groups", "sg-123"),
			Entry("with node-labels flag", "--node-labels", "partition=backend,nodeclass=hugememory"),
			Entry("with node-zones flag", "--node-zones", "zone1,zone2,zone3", "--zones", "zone1,zone2,zone3,zone4"),
			// commons node group IAM flags
			Entry("with asg-access flag", "--asg-access", "true"),
			Entry("with external-dns-access flag", "--external-dns-access", "true"),
			Entry("with full-ecr-access flag", "--full-ecr-access", "true"),
			Entry("with appmesh-access flag", "--appmesh-access", "true"),
			Entry("with alb-ingress-access flag", "--alb-ingress-access", "true"),
			Entry("with managed flag unset", "--managed", "false"),
		)

		DescribeTable("invalid flags or arguments",
			func(c invalidParamsCase) {
				commandArgs := append([]string{"cluster", "--managed=false"}, c.args...)
				cmd := newDefaultCmd(commandArgs...)
				_, err := cmd.execute()
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(ContainSubstring(c.error)))
			},
			Entry("with cluster name as argument and flag", invalidParamsCase{
				args:  []string{"clusterName", "--name", "clusterName"},
				error: "--name=clusterName and argument clusterName cannot be used at the same time",
			}),
			Entry("with invalid flags", invalidParamsCase{
				args:  []string{"cluster", "--invalid", "dummy"},
				error: "unknown flag: --invalid",
			}),
		)
	})

	Describe("managed node group", func() {
		DescribeTable("create cluster successfully",
			func(args ...string) {
				commandArgs := append([]string{"cluster"}, args...)
				cmd := newMockEmptyCmd(commandArgs...)
				count := 0
				cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
					createClusterCmdWithRunFunc(cmd, func(cmd *cmdutils.Cmd, ngFilter *filter.NodeGroupFilter, params *cmdutils.CreateClusterCmdParams) error {
						Expect(cmd.ClusterConfig.Metadata.Name).NotTo(BeNil())
						count++
						return nil
					})
				})
				_, err := cmd.execute()
				Expect(err).To(Not(HaveOccurred()))
				Expect(count).To(Equal(1))
			},
			Entry("without cluster name", ""),
			Entry("with cluster name as flag", "--name", "clusterName"),
			Entry("with cluster name as argument", "clusterName"),
			Entry("with cluster name with hyphen as flag", "--name", "my-cluster-name-is-fine10"),
			Entry("with cluster name with hyphen as argument", "my-Cluster-name-is-fine10"),
			// vpc networking flags
			Entry("with vpc-cidr flag", "--vpc-cidr", "10.0.0.0/20"),
			Entry("with vpc-private-subnets flag", "--vpc-private-subnets", "10.0.0.0/24"),
			Entry("with vpc-public-subnets flag", "--vpc-public-subnets", "10.0.0.0/24"),
			Entry("with vpc-from-kops-cluster flag", "--vpc-from-kops-cluster", "dummy-kops-cluster"),
			Entry("with vpc-nat-mode flag", "--vpc-nat-mode", "Single"),
			// kubeconfig flags
			Entry("with write-kubeconfig flag", "--write-kubeconfig"),
			Entry("with kubeconfig flag", "--kubeconfig", "~/.kube"),
			Entry("with authenticator-role-arn flag", "--authenticator-role-arn", "arn::dummy::123/role"),
			Entry("with auto-kubeconfig flag", "--auto-kubeconfig"),
			// common node group flags
			Entry("with node-type flag", "--node-type", "m5.large"),
			Entry("with nodes flag", "--nodes", "2"),
			Entry("with nodes-min flag", "--nodes-min", "2"),
			Entry("with nodes-max flag", "--nodes-max", "2"),
			Entry("with node-volume-size flag", "--node-volume-size", "2"),
			Entry("with ssh-access flag", "--ssh-access", "true"),
			Entry("with ssh-public-key flag", "--ssh-public-key", "dummy-public-key"),
			Entry("with enable-ssm flag", "--enable-ssm"),
			Entry("with node-ami-family flag", "--node-ami-family", "AmazonLinux2"),
			Entry("with node-private-networking flag", "--node-private-networking", "true"),
			Entry("with node-labels flag", "--node-labels", "partition=backend,nodeclass=hugememory"),
			Entry("with node-zones flag", "--node-zones", "zone1,zone2,zone3", "--zones", "zone1,zone2,zone3,zone4"),
			Entry("with asg-access flag", "--asg-access", "true"),
			Entry("with external-dns-access flag", "--external-dns-access", "true"),
			Entry("with full-ecr-access flag", "--full-ecr-access", "true"),
			Entry("with appmesh-access flag", "--appmesh-access", "true"),
			Entry("with alb-ingress-access flag", "--alb-ingress-access", "true"),
		)

		DescribeTable("invalid flags or arguments",
			func(c invalidParamsCase) {
				commandArgs := append([]string{"cluster"}, c.args...)
				cmd := newDefaultCmd(commandArgs...)
				_, err := cmd.execute()
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(ContainSubstring(c.error)))
			},
			Entry("with cluster name as argument and flag", invalidParamsCase{
				args:  []string{"clusterName", "--name", "clusterName"},
				error: "--name=clusterName and argument clusterName cannot be used at the same time",
			}),
			Entry("with invalid flags", invalidParamsCase{
				args:  []string{"cluster", "--invalid", "dummy"},
				error: "unknown flag: --invalid",
			}),
			Entry("with enableSsm disabled", invalidParamsCase{
				args:  []string{"--name=test", "--enable-ssm=false"},
				error: "SSM agent is now built into EKS AMIs and cannot be disabled",
			}),
			Entry("with node zones without zones", invalidParamsCase{
				args:  []string{"--zones=zone1,zone2", "--node-zones=zone3"},
				error: "validation for --zones and --node-zones failed: node-zones [zone3] must be a subset of zones [zone1 zone2]; \"zone3\" was not found in zones",
			}),
		)
	})

	type createClusterEntry struct {
		updateClusterConfig func(*api.ClusterConfig)
		updateClusterParams func(*cmdutils.CreateClusterCmdParams)
		updateMocks         func(*mockprovider.MockProvider)
		mockOutposts        bool

		expectedErr string
	}

	Describe("[Outposts] cluster version issues", func() {
		Describe("version not set", func() {
			It("should return an error", func() {
				cfg := &api.ClusterConfig{
					TypeMeta: api.ClusterConfigTypeMeta(),
					Metadata: &api.ClusterMeta{
						Name:    "cluster-1",
						Region:  "us-west-2",
						Version: "",
					},
					Outpost: &api.Outpost{
						ControlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
					},
				}

				cmd := newDefaultCmd("cluster", "--config-file", ctltest.CreateConfigFile(cfg))
				_, err := cmd.execute()
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(ContainSubstring("cluster version must be explicitly set to 1.21 for Outposts clusters as only version 1.21 is currently supported")))
			})
		})

		Describe("version set to unsupported version", func() {
			It("should return an error", func() {
				cfg := &api.ClusterConfig{
					TypeMeta: api.ClusterConfigTypeMeta(),
					Metadata: &api.ClusterMeta{
						Name:    "cluster-1",
						Region:  "us-west-2",
						Version: "1.20",
					},
					Outpost: &api.Outpost{
						ControlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
					},
				}

				cmd := newDefaultCmd("cluster", "--config-file", ctltest.CreateConfigFile(cfg))
				_, err := cmd.execute()
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(ContainSubstring("only version 1.21 is supported on Outposts")))
			})
		})
	})

	DescribeTable("doCreateCluster", func(ce createClusterEntry) {
		p := mockprovider.NewMockProvider()
		defaultProviderMocks(p, defaultOutput, ce.mockOutposts)
		if ce.updateMocks != nil {
			ce.updateMocks(p)
		}
		fk := &fakes.FakeKubeProvider{}
		clientset := kubefake.NewSimpleClientset()
		client, err := kubernetes.NewRawClient(clientset, &restclient.Config{})
		Expect(err).NotTo(HaveOccurred())
		fk.NewRawClientReturns(client, nil)
		fk.ServerVersionReturns("1.22", nil)
		msp := &mockSessionProvider{}
		ctl := &eks.ClusterProvider{
			AWSProvider: p,
			Status: &eks.ProviderStatus{
				ClusterInfo: &eks.ClusterInfo{
					Cluster: testutils.NewFakeCluster("my-cluster", ""),
				},
				SessionCreds: msp,
			},
			KubeProvider: fk,
		}

		clusterConfig := api.NewClusterConfig()
		clusterConfig.Metadata.Name = "my-cluster"
		clusterConfig.VPC.ClusterEndpoints = api.ClusterEndpointAccessDefaults()

		if ce.updateClusterConfig != nil {
			ce.updateClusterConfig(clusterConfig)
		}
		cmd := &cmdutils.Cmd{
			ClusterConfig: clusterConfig,
			ProviderConfig: api.ProviderConfig{
				WaitTimeout: time.Second * 1,
			},
		}
		filter := filter.NewNodeGroupFilter()
		params := &cmdutils.CreateClusterCmdParams{
			Subnets: map[api.SubnetTopology]*[]string{
				api.SubnetTopologyPrivate: {},
				api.SubnetTopologyPublic:  {},
			},
		}
		if ce.updateClusterParams != nil {
			ce.updateClusterParams(params)
		}
		filter.SetExcludeAll(params.WithoutNodeGroup)
		err = doCreateCluster(cmd, filter, params, ctl)
		if ce.expectedErr != "" {
			Expect(err).To(MatchError(ContainSubstring(ce.expectedErr)))
			return
		}

		Expect(err).NotTo(HaveOccurred())
	},
		Entry("standard cluster", createClusterEntry{}),

		Entry("[Outposts] control plane on Outposts with valid config", createClusterEntry{
			updateClusterConfig: func(c *api.ClusterConfig) {
				c.Metadata.Version = api.Version1_21
				c.Outpost = &api.Outpost{
					ControlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
				}
			},
			mockOutposts: true,
		}),

		Entry("[Outposts] unavailable instance type specified for the control plane", createClusterEntry{
			updateClusterConfig: func(c *api.ClusterConfig) {
				c.Metadata.Version = api.Version1_21
				c.Outpost = &api.Outpost{
					ControlPlaneOutpostARN:   "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
					ControlPlaneInstanceType: "t2.medium",
				}
			},
			mockOutposts: true,

			expectedErr: `instance type "t2.medium" does not exist in Outpost "arn:aws:outposts:us-west-2:1234:outpost/op-1234"`,
		}),

		Entry("[Outposts] available instance type specified for the control plane", createClusterEntry{
			updateClusterConfig: func(c *api.ClusterConfig) {
				c.Metadata.Version = api.Version1_21
				c.Outpost = &api.Outpost{
					ControlPlaneOutpostARN:   "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
					ControlPlaneInstanceType: "m5.xlarge",
				}
			},
			mockOutposts: true,
		}),

		Entry("[Outposts] nodegroups specified when the VPC will be created by eksctl", createClusterEntry{
			updateClusterConfig: func(c *api.ClusterConfig) {
				c.Metadata.Version = api.Version1_21
				c.Outpost = &api.Outpost{
					ControlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
				}
				c.NodeGroups = []*api.NodeGroup{
					api.NewNodeGroup(),
				}
			},
			mockOutposts: true,

			expectedErr: "cannot create nodegroups on Outposts when the VPC is created by eksctl as it will not have connectivity to the API server; please rerun the command with `--without-nodegroup` and run `eksctl create nodegroup` after associating the VPC with a local gateway and ensuring connectivity to the API server",
		}),

		Entry("[Outposts] nodegroups specified on Outposts but the control plane is not on Outposts", createClusterEntry{
			updateClusterConfig: func(c *api.ClusterConfig) {
				c.Metadata.Version = api.Version1_21
				ng := api.NewNodeGroup()
				ng.OutpostARN = "arn:aws:outposts:us-west-2:1234:outpost/op-1234"
				c.NodeGroups = []*api.NodeGroup{ng}
			},
			mockOutposts: true,

			expectedErr: "creating nodegroups on Outposts when the control plane is not on Outposts is not supported during cluster creation; " +
				"either create the nodegroups after cluster creation or consider creating the control plane on Outposts",
		}),

		Entry("[Outposts] nodegroups specified when the VPC will be created by eksctl, but with --without-nodegroup", createClusterEntry{
			updateClusterConfig: func(c *api.ClusterConfig) {
				c.Metadata.Version = api.Version1_21
				c.Outpost = &api.Outpost{
					ControlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
				}
				c.NodeGroups = []*api.NodeGroup{
					api.NewNodeGroup(),
				}
			},
			updateClusterParams: func(params *cmdutils.CreateClusterCmdParams) {
				params.WithoutNodeGroup = true
			},
			mockOutposts: true,
		}),

		Entry("[Outposts] specified Outpost does not exist", createClusterEntry{
			updateClusterConfig: func(c *api.ClusterConfig) {
				c.Metadata.Version = api.Version1_21
				c.Outpost = &api.Outpost{
					ControlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
				}
			},
			updateMocks: func(provider *mockprovider.MockProvider) {
				provider.MockOutposts().On("GetOutpost", mock.Anything, &outposts.GetOutpostInput{
					OutpostId: aws.String("arn:aws:outposts:us-west-2:1234:outpost/op-1234"),
				}).Return(nil, &outpoststypes.NotFoundException{Message: aws.String("not found")})
			},
			expectedErr: "error getting Outpost details: NotFoundException: not found",
		}),
	)
})

var defaultOutput = []cftypes.Output{
	{
		OutputKey:   aws.String("ClusterSecurityGroupId"),
		OutputValue: aws.String("csg-1234"),
	},
	{
		OutputKey:   aws.String("SecurityGroup"),
		OutputValue: aws.String("sg-1"),
	},
	{
		OutputKey:   aws.String("VPC"),
		OutputValue: aws.String("vpc-1"),
	},
	{
		OutputKey:   aws.String("SharedNodeSecurityGroup"),
		OutputValue: aws.String("sg-1"),
	},
	{
		OutputKey:   aws.String("FeatureNATMode"),
		OutputValue: aws.String("Single"),
	},
	{
		OutputKey:   aws.String("SubnetsPrivate"),
		OutputValue: aws.String("sub-priv-1 sub-priv-2 sub-priv-3"),
	},
	{
		OutputKey:   aws.String("SubnetsPublic"),
		OutputValue: aws.String("sub-pub-1 sub-pub-2 sub-pub-3"),
	},
	{
		OutputKey:   aws.String("ServiceRoleARN"),
		OutputValue: aws.String("arn:aws:iam::123456:role/amazingrole-1"),
	},
	{
		OutputKey:   aws.String("ARN"),
		OutputValue: aws.String("arn:aws:iam::123456:role/amazingrole-2"),
	},
	{
		OutputKey:   aws.String("CertificateAuthorityData"),
		OutputValue: aws.String("dGVzdAo="),
	},
	{
		OutputKey:   aws.String("ClusterStackName"),
		OutputValue: aws.String("eksctl-my-cluster-cluster"),
	},
	{
		OutputKey:   aws.String("Endpoint"),
		OutputValue: aws.String("https://endpoint.com"),
	},
}

func defaultProviderMocks(p *mockprovider.MockProvider, output []cftypes.Output, controlPlaneOnOutposts bool) {
	p.MockEC2().On("DescribeAvailabilityZones", mock.Anything, &ec2.DescribeAvailabilityZonesInput{
		Filters: []ec2types.Filter{{
			Name:   aws.String("region-name"),
			Values: []string{"us-west-2"},
		}, {
			Name:   aws.String("state"),
			Values: []string{string(ec2types.AvailabilityZoneStateAvailable)},
		}, {
			Name:   aws.String("zone-type"),
			Values: []string{string(ec2types.LocationTypeAvailabilityZone)},
		}},
	}).Return(&ec2.DescribeAvailabilityZonesOutput{
		AvailabilityZones: []ec2types.AvailabilityZone{
			{
				GroupName: aws.String("name"),
				ZoneName:  aws.String("us-west-2-1b"),
				ZoneId:    aws.String("id"),
			},
			{
				GroupName: aws.String("name"),
				ZoneName:  aws.String("us-west-2-1a"),
				ZoneId:    aws.String("id"),
			}},
	}, nil)
	p.MockCloudFormation().On("ListStacks", mock.Anything, mock.Anything).Return(&cloudformation.ListStacksOutput{
		StackSummaries: []cftypes.StackSummary{
			{
				StackName:   aws.String("eksctl-my-cluster-cluster"),
				StackStatus: "CREATE_COMPLETE",
			},
		},
	}, nil)
	p.MockCloudFormation().On("DescribeStacks", mock.Anything, mock.Anything).Return(&cloudformation.DescribeStacksOutput{
		Stacks: []cftypes.Stack{
			{
				StackName:   aws.String("eksctl-my-cluster-cluster"),
				StackStatus: "CREATE_COMPLETE",
				Tags: []cftypes.Tag{
					{
						Key:   aws.String(api.ClusterNameTag),
						Value: aws.String("eksctl-my-cluster-cluster"),
					},
				},
				Outputs: output,
			},
		},
	}, nil)

	const outpostID = "arn:aws:outposts:us-west-2:1234:outpost/op-1234"
	var outpostConfig *ekstypes.OutpostConfigResponse
	if controlPlaneOnOutposts {
		mockOutposts(p, outpostID)
		outpostConfig = &ekstypes.OutpostConfigResponse{
			OutpostArns:              []string{outpostID},
			ControlPlaneInstanceType: aws.String("m5.xlarge"),
		}
	}
	p.MockEKS().On("DescribeCluster", mock.Anything, mock.Anything).Return(&awseks.DescribeClusterOutput{
		Cluster: &ekstypes.Cluster{
			CertificateAuthority: &ekstypes.Certificate{
				Data: aws.String("dGVzdAo="),
			},
			Endpoint:                aws.String("endpoint"),
			Arn:                     aws.String("arn"),
			KubernetesNetworkConfig: nil,
			Logging:                 nil,
			Name:                    aws.String("my-cluster"),
			PlatformVersion:         aws.String("1.22"),
			ResourcesVpcConfig: &ekstypes.VpcConfigResponse{
				ClusterSecurityGroupId: aws.String("csg-1234"),
				EndpointPublicAccess:   true,
				PublicAccessCidrs:      []string{"1.2.3.4/24", "1.2.3.4/12"},
				SecurityGroupIds:       []string{"sg-1", "sg-2"},
				SubnetIds:              []string{"sub-1", "sub-2"},
				VpcId:                  aws.String("vpc-1"),
			},
			Status: "CREATE_COMPLETE",
			Tags: map[string]string{
				api.ClusterNameTag: "eksctl-my-cluster-cluster",
			},
			Version:       aws.String("1.22"),
			OutpostConfig: outpostConfig,
		},
	}, nil)

	p.MockEC2().On("DescribeImages", mock.Anything, mock.Anything).
		Return(&ec2.DescribeImagesOutput{
			Images: []ec2types.Image{
				{
					ImageId:        aws.String("ami-123"),
					State:          ec2types.ImageStateAvailable,
					OwnerId:        aws.String("123"),
					RootDeviceType: ec2types.DeviceTypeEbs,
					RootDeviceName: aws.String("/dev/sda1"),
					BlockDeviceMappings: []ec2types.BlockDeviceMapping{
						{
							DeviceName: aws.String("/dev/sda1"),
							Ebs: &ec2types.EbsBlockDevice{
								Encrypted: aws.Bool(false),
							},
						},
					},
				},
			},
		}, nil)
	p.MockEC2().On("DescribeSubnets", mock.Anything, mock.Anything).Return(&ec2.DescribeSubnetsOutput{
		Subnets: []ec2types.Subnet{},
	}, nil)
	p.MockEC2().On("DescribeVpcs", mock.Anything, mock.Anything).Return(&ec2.DescribeVpcsOutput{
		Vpcs: []ec2types.Vpc{
			{
				VpcId:     aws.String("vpc-1"),
				CidrBlock: aws.String("192.168.0.0/16"),
			},
		},
	}, nil)
	p.MockCloudFormation().On("CreateStack", mock.Anything, mock.Anything).Return(&cloudformation.CreateStackOutput{
		StackId: aws.String("eksctl-my-cluster-cluster"),
	}, nil)
}

func mockOutposts(provider *mockprovider.MockProvider, outpostID string) {
	provider.MockOutposts().On("GetOutpost", mock.Anything, &outposts.GetOutpostInput{
		OutpostId: aws.String(outpostID),
	}).Return(&outposts.GetOutpostOutput{
		Outpost: &outpoststypes.Outpost{
			AvailabilityZone: aws.String("us-west-2a"),
		},
	}, nil)
	provider.MockOutposts().On("GetOutpostInstanceTypes", mock.Anything, &outposts.GetOutpostInstanceTypesInput{
		OutpostId: aws.String(outpostID),
	}).Return(&outposts.GetOutpostInstanceTypesOutput{
		InstanceTypes: []outpoststypes.InstanceTypeItem{
			{
				InstanceType: aws.String("m5.xlarge"),
			},
		},
	}, nil)
	provider.MockEC2().On("DescribeInstanceTypes", mock.Anything, &ec2.DescribeInstanceTypesInput{
		InstanceTypes: []ec2types.InstanceType{"m5.xlarge"},
	}).Return(&ec2.DescribeInstanceTypesOutput{
		InstanceTypes: []ec2types.InstanceTypeInfo{
			{
				InstanceType: "m5.xlarge",
				VCpuInfo: &ec2types.VCpuInfo{
					DefaultVCpus:          aws.Int32(4),
					DefaultCores:          aws.Int32(2),
					DefaultThreadsPerCore: aws.Int32(2),
				},
				MemoryInfo: &ec2types.MemoryInfo{
					SizeInMiB: aws.Int64(16384),
				},
			},
		},
	}, nil)
}

type mockSessionProvider struct {
}

func (m *mockSessionProvider) Get() (credentials.Value, error) {
	return credentials.Value{
		AccessKeyID:     "key-id",
		SecretAccessKey: "secret-access-key",
		SessionToken:    "token",
		ProviderName:    "aws",
	}, nil
}
