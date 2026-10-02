//go:build integration

package managed

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	clusterutils "github.com/weaveworks/eksctl/integration/utilities/cluster"
	"github.com/weaveworks/eksctl/integration/utilities/kube"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParamsWithGivenClusterName("managed", "test-cluster")
}

func TestManaged(t *testing.T) {
	testutils.RegisterAndRun(t)
}

const initialAl2Nodegroup = "ng-al2"

var _ = SynchronizedBeforeSuite(func() {
	fmt.Fprintf(GinkgoWriter, "Using kubeconfig: %s\n", params.KubeconfigPath)

	cmd := params.EksctlCreateCmd.WithArgs(
		"cluster",
		"--verbose", "4",
		"--name", params.ClusterName,
		"--tags", "alpha.eksctl.io/description=eksctl integration test",
		"--managed",
		"--nodegroup-name", initialAl2Nodegroup,
		"--node-labels", "ng-name="+initialAl2Nodegroup,
		"--nodes", "2",
		"--instance-types", "t3a.xlarge",
		"--version", params.Version,
		"--kubeconfig", params.KubeconfigPath,
	)
	Expect(cmd).To(RunSuccessfully())
}, func() {})

var _ = Describe("(Integration) Create Managed Nodegroups", func() {

	const (
		updateConfigNodegroup    = "ng-update-config"
		al2023Nodegroup          = "ng-al2023"
		bottlerocketNodegroup    = "ng-bottlerocket"
		bottlerocketGPUNodegroup = "ng-bottlerocket-gpu"
		ubuntuNodegroup          = "ng-ubuntu"
		publicNodeGroup          = "ng-public"
		privateNodeGroup         = "ng-private"
		connTrackingNodeGroup    = "ng-conn-tracking"
	)

	var (
		makeClusterConfig = func() *api.ClusterConfig {
			clusterConfig := api.NewClusterConfig()
			clusterConfig.Metadata.Name = params.ClusterName
			clusterConfig.Metadata.Region = params.Region
			clusterConfig.Metadata.Version = params.Version
			return clusterConfig
		}

		checkNg = func(ngName string) {
			cmd := params.EksctlUtilsCmd.WithArgs(
				"nodegroup-health",
				"--cluster", params.ClusterName,
				"--name", ngName,
			)
			Expect(cmd).To(RunSuccessfullyWithOutputString(ContainSubstring("active")))
		}
	)

	type managedCLIEntry struct {
		createArgs []string
	}

	DescribeTable("Managed CLI features", func(m managedCLIEntry) {
		runAssertions := func(createArgs ...string) {
			cmd := params.EksctlCreateCmd.
				WithArgs(createArgs...).
				WithArgs(m.createArgs...)

			Expect(cmd).To(RunSuccessfully())
		}

		// Run the same assertions for both `create cluster` and `create nodegroup`
		runAssertions("cluster")
		runAssertions(
			"nodegroup",
			"--cluster", params.ClusterName,
		)
	},
		Entry("Windows AMI with dry-run", managedCLIEntry{
			createArgs: []string{
				"--node-ami-family=WindowsServer2019FullContainer",
				"--dry-run",
			},
		}),

		Entry("Bottlerocket with dry-run", managedCLIEntry{
			createArgs: []string{
				"--node-ami-family=Bottlerocket",
				"--instance-prefix=bottle",
				"--instance-name=rocket",
				"--dry-run",
			},
		}),

		Entry("Ubuntu with dry-run", managedCLIEntry{
			createArgs: []string{
				"--node-ami-family=Ubuntu2004",
				"--dry-run",
			},
		}),
	)

	Context("adding new managed nodegroups", func() {
		params.LogStacksEventsOnFailure()

		It("supports a public nodegroup", func() {
			By("creating it")
			cmd := params.EksctlCreateCmd.WithArgs(
				"nodegroup",
				"--cluster", params.ClusterName,
				"--nodes", "4",
				"--managed",
				"--instance-types", "t3a.xlarge",
				publicNodeGroup,
			)
			Expect(cmd).To(RunSuccessfully())

			By("ensuring it is healthy")
			checkNg(publicNodeGroup)

			By("deleting it")
			cmd = params.EksctlDeleteCmd.WithArgs(
				"nodegroup",
				"--verbose", "4",
				"--cluster", params.ClusterName,
				publicNodeGroup,
			)
			Expect(cmd).To(RunSuccessfully())
		})

		It("supports a private nodegroup", func() {
			By("creating it")
			cmd := params.EksctlCreateCmd.WithArgs(
				"nodegroup",
				"--cluster", params.ClusterName,
				"--nodes", "2",
				"--managed",
				"--instance-types", "t3a.xlarge",
				"--node-private-networking",
				privateNodeGroup,
			)
			Expect(cmd).To(RunSuccessfully())

			By("ensuring it is healthy")
			checkNg(privateNodeGroup)

			By("deleting it")
			cmd = params.EksctlDeleteCmd.WithArgs(
				"nodegroup",
				"--verbose", "4",
				"--cluster", params.ClusterName,
				privateNodeGroup,
			)
			Expect(cmd).To(RunSuccessfully())
		})

		It("supports a public multi-AZ nodegroup with connection tracking timeouts", func() {
			// A non-EFA managed nodegroup that sets connectionTracking is the first managed-nodegroup
			// path to put a NetworkInterfaces block in the launch template without also pinning the
			// nodegroup to a single subnet, which efaEnabled does. This exercises that combination on
			// a real cluster: a nodegroup spanning every public subnet, whose nodes must still join
			// and still receive the subnets' auto-assigned public IPv4 addresses.
			connectionTracking := &api.ConnectionTracking{
				TCPEstablishedTimeout: aws.Int(432000),
				UDPStreamTimeout:      aws.Int(180),
				UDPTimeout:            aws.Int(60),
			}
			clusterConfig := makeClusterConfig()
			clusterConfig.ManagedNodeGroups = []*api.ManagedNodeGroup{
				{
					NodeGroupBase: &api.NodeGroupBase{
						Name:               connTrackingNodeGroup,
						InstanceType:       "t3a.xlarge",
						ScalingConfig:      &api.ScalingConfig{DesiredCapacity: aws.Int(2)},
						ConnectionTracking: connectionTracking,
					},
				},
			}
			By("creating it")
			cmd := params.EksctlCreateCmd.
				WithArgs(
					"nodegroup",
					"--config-file", "-",
					"--verbose", "4",
				).
				WithoutArg("--region", params.Region).
				WithStdin(clusterutils.Reader(clusterConfig))
			Expect(cmd).To(RunSuccessfully())
			By("ensuring it is healthy")
			checkNg(connTrackingNodeGroup)
			ctx := context.Background()
			clusterProvider, err := eks.New(ctx, &api.ProviderConfig{Region: params.Region}, clusterConfig)
			Expect(err).NotTo(HaveOccurred())
			ctl := clusterProvider.AWSProvider
			By("checking the nodegroup spans more than one subnet")
			nodegroupOutput, err := ctl.EKS().DescribeNodegroup(ctx, &awseks.DescribeNodegroupInput{
				ClusterName:   &params.ClusterName,
				NodegroupName: aws.String(connTrackingNodeGroup),
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodegroupOutput.Nodegroup.Subnets)).To(BeNumerically(">", 1),
				"the nodegroup should span every public subnet, so the launch template's network interface is used across subnets")
			By("checking the launch template carries the timeouts on its network interface, and no instance-level security groups")
			Expect(nodegroupOutput.Nodegroup.LaunchTemplate).NotTo(BeNil())
			launchTemplateOutput, err := ctl.EC2().DescribeLaunchTemplateVersions(ctx, &ec2.DescribeLaunchTemplateVersionsInput{
				LaunchTemplateId: nodegroupOutput.Nodegroup.LaunchTemplate.Id,
				Versions:         []string{*nodegroupOutput.Nodegroup.LaunchTemplate.Version},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(launchTemplateOutput.LaunchTemplateVersions).To(HaveLen(1))
			launchTemplateData := launchTemplateOutput.LaunchTemplateVersions[0].LaunchTemplateData
			Expect(launchTemplateData.SecurityGroupIds).To(BeEmpty(),
				"security groups must be on the network interface instead, as EKS requires when a launch template specifies one")
			Expect(launchTemplateData.NetworkInterfaces).To(HaveLen(1))
			launchTemplateInterface := launchTemplateData.NetworkInterfaces[0]
			Expect(launchTemplateInterface.Groups).NotTo(BeEmpty())
			Expect(launchTemplateInterface.AssociatePublicIpAddress).To(BeNil(),
				"leaving this unset is what keeps the subnet's auto-assign public IPv4 setting in effect")
			Expect(launchTemplateInterface.ConnectionTrackingSpecification).To(Equal(&ec2types.ConnectionTrackingSpecification{
				TcpEstablishedTimeout: aws.Int32(432000),
				UdpStreamTimeout:      aws.Int32(180),
				UdpTimeout:            aws.Int32(60),
			}))
			By("checking the timeouts reached the primary interface of every node, and that the nodes have public IPs")
			Eventually(func() int {
				return len(tests.NodeInstanceIDs(params.KubeconfigPath, connTrackingNodeGroup))
			}, "5m", "15s").Should(Equal(2), "both nodes should register with the cluster")
			instances := tests.NodeInstances(params.KubeconfigPath, params.Region, connTrackingNodeGroup)
			Expect(instances).To(HaveLen(2))
			for _, instance := range instances {
				Expect(instance.PublicIpAddress).NotTo(BeNil(),
					fmt.Sprintf("instance %s is in a public subnet and should have been assigned a public IPv4 address", aws.ToString(instance.InstanceId)))
				primaryInterface := tests.PrimaryNetworkInterface(instance)
				Expect(primaryInterface).NotTo(BeNil(),
					fmt.Sprintf("instance %s should have an interface at device index 0 on network card 0", aws.ToString(instance.InstanceId)))
				Expect(primaryInterface.ConnectionTrackingConfiguration).To(Equal(&ec2types.ConnectionTrackingSpecificationResponse{
					TcpEstablishedTimeout: aws.Int32(432000),
					UdpStreamTimeout:      aws.Int32(180),
					UdpTimeout:            aws.Int32(60),
				}))
			}
			By("reporting connectionTracking as unchanged on `eksctl update nodegroup`, since it only takes effect at creation")
			clusterConfig.ManagedNodeGroups[0].ConnectionTracking.TCPEstablishedTimeout = aws.Int(3600)
			clusterConfig.ManagedNodeGroups[0].UpdateConfig = &api.NodeGroupUpdateConfig{
				MaxUnavailable: aws.Int(1),
			}
			cmd = params.EksctlUpdateCmd.
				WithArgs(
					"nodegroup",
					"--config-file", "-",
					"--verbose", "4",
				).
				WithoutArg("--region", params.Region).
				WithStdin(clusterutils.Reader(clusterConfig))
			Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
				ContainElement(ContainSubstring(fmt.Sprintf("unchanged fields for nodegroup %s", connTrackingNodeGroup))),
				ContainElement(ContainSubstring("ConnectionTracking")),
			))
			By("deleting it")
			cmd = params.EksctlDeleteCmd.WithArgs(
				"nodegroup",
				"--verbose", "4",
				"--cluster", params.ClusterName,
				connTrackingNodeGroup,
			)
			Expect(cmd).To(RunSuccessfully())
		})
		It("rejects connection tracking on a non-Nitro instance type before creating anything", func() {
			clusterConfig := makeClusterConfig()
			clusterConfig.ManagedNodeGroups = []*api.ManagedNodeGroup{
				{
					NodeGroupBase: &api.NodeGroupBase{
						Name:         "ng-conn-tracking-xen",
						InstanceType: "t2.medium",
						ConnectionTracking: &api.ConnectionTracking{
							TCPEstablishedTimeout: aws.Int(432000),
						},
					},
				},
			}
			cmd := params.EksctlCreateCmd.
				WithArgs(
					"nodegroup",
					"--config-file", "-",
					"--verbose", "4",
				).
				WithoutArg("--region", params.Region).
				WithStdin(clusterutils.Reader(clusterConfig))
			// Run once and assert on both the exit code and the message; each use of a runner matcher
			// would run the command again.
			session := cmd.Run()
			Expect(session.ExitCode()).NotTo(Equal(0), "eksctl should refuse the nodegroup")
			Expect(string(session.Buffer().Contents())).To(
				ContainSubstring("connectionTracking is only supported on Nitro-based instance types; t2.medium cannot be used with it"))
		})
		It("supports a nodegroup with taints", func() {
			taints := []api.NodeGroupTaint{
				{
					Key:    "key1",
					Value:  "value1",
					Effect: "NoSchedule",
				},
				{
					Key:    "key2",
					Effect: "NoSchedule",
				},
				{
					Key:    "key3",
					Value:  "value2",
					Effect: "NoExecute",
				},
			}
			clusterConfig := makeClusterConfig()
			clusterConfig.ManagedNodeGroups = []*api.ManagedNodeGroup{
				{
					NodeGroupBase: &api.NodeGroupBase{
						Name: "taints",
					},
					Taints: taints,
				},
			}

			By("creating it")
			cmd := params.EksctlCreateCmd.
				WithArgs(
					"nodegroup",
					"--config-file", "-",
					"--verbose", "4",
				).
				WithoutArg("--region", params.Region).
				WithStdin(clusterutils.Reader(clusterConfig))
			Expect(cmd).To(RunSuccessfully())

			By("ensuring it is healthy")
			checkNg("taints")

			By("asserting node taints")
			config, err := clientcmd.BuildConfigFromFlags("", params.KubeconfigPath)
			Expect(err).NotTo(HaveOccurred())
			clientset, err := kubernetes.NewForConfig(config)
			Expect(err).NotTo(HaveOccurred())

			mapTaints := func(taints []api.NodeGroupTaint) []corev1.Taint {
				var ret []corev1.Taint
				for _, t := range taints {
					ret = append(ret, corev1.Taint{
						Key:    t.Key,
						Value:  t.Value,
						Effect: t.Effect,
					})
				}
				return ret
			}
			tests.AssertNodeTaints(tests.ListNodes(clientset, "taints"), mapTaints(taints))

			By("deleting it")
			cmd = params.EksctlDeleteCmd.WithArgs(
				"nodegroup",
				"--verbose", "4",
				"--cluster", params.ClusterName,
				"taints",
			)
			Expect(cmd).To(RunSuccessfully())
		})

		It("supports a bottlerocket nodegroup with gpu nodes", func() {
			clusterConfig := makeClusterConfig()
			clusterConfig.ManagedNodeGroups = []*api.ManagedNodeGroup{
				{
					NodeGroupBase: &api.NodeGroupBase{
						Name:         bottlerocketGPUNodegroup,
						VolumeSize:   aws.Int(35),
						AMIFamily:    "Bottlerocket",
						InstanceType: "g4dn.xlarge",
						Bottlerocket: &api.NodeGroupBottlerocket{
							EnableAdminContainer: api.Enabled(),
							Settings: &api.InlineDocument{
								"motd": "Bottlerocket is the future",
								"network": map[string]string{
									"hostname": "custom-bottlerocket-host",
								},
							},
						},
					},
				},
			}

			By("creating it")
			cmd := params.EksctlCreateCmd.
				WithArgs(
					"nodegroup",
					"--config-file", "-",
					"--skip-outdated-addons-check",
					"--verbose", "4",
				).
				WithoutArg("--region", params.Region).
				WithStdin(clusterutils.Reader(clusterConfig))
			Expect(cmd).To(RunSuccessfully())

			By("ensuring it is healthy")
			checkNg(bottlerocketGPUNodegroup)
		})

		It("supports AmazonLinux2023 nodegroups", func() {
			clusterConfig := makeClusterConfig()
			clusterConfig.ManagedNodeGroups = []*api.ManagedNodeGroup{
				{
					NodeGroupBase: &api.NodeGroupBase{
						Name:      al2023Nodegroup,
						AMIFamily: "AmazonLinux2023",
					},
				},
			}

			By("creating it")
			Expect(params.EksctlCreateCmd.
				WithArgs(
					"nodegroup",
					"--config-file", "-",
					"--verbose", "4",
				).
				WithoutArg("--region", params.Region).
				WithStdin(clusterutils.Reader(clusterConfig))).
				To(RunSuccessfully())

			By("ensuring it is healthy")
			checkNg(al2023Nodegroup)
		})

		It("supports bottlerocket and ubuntu nodegroups with additional volumes", func() {
			clusterConfig := makeClusterConfig()
			clusterConfig.ManagedNodeGroups = []*api.ManagedNodeGroup{
				{
					NodeGroupBase: &api.NodeGroupBase{
						Name:         bottlerocketNodegroup,
						VolumeSize:   aws.Int(35),
						AMIFamily:    "Bottlerocket",
						InstanceType: "t3a.xlarge",
						Bottlerocket: &api.NodeGroupBottlerocket{
							EnableAdminContainer: api.Enabled(),
							Settings: &api.InlineDocument{
								"motd": "Bottlerocket is the future",
								"network": map[string]string{
									"hostname": "custom-bottlerocket-host",
								},
							},
						},
						AdditionalVolumes: []*api.VolumeMapping{
							{
								VolumeName: aws.String("/dev/sdb"),
							},
						},
					},
					Taints: []api.NodeGroupTaint{
						{
							Key:    "key2",
							Value:  "value2",
							Effect: "PreferNoSchedule",
						},
					},
				},
				{
					NodeGroupBase: &api.NodeGroupBase{
						Name:         ubuntuNodegroup,
						VolumeSize:   aws.Int(25),
						AMIFamily:    "Ubuntu2404",
						InstanceType: "t3a.xlarge",
					},
				},
			}

			By("creating it")
			cmd := params.EksctlCreateCmd.
				WithArgs(
					"nodegroup",
					"--config-file", "-",
					"--verbose", "4",
				).
				WithoutArg("--region", params.Region).
				WithStdin(clusterutils.Reader(clusterConfig))

			Expect(cmd).To(RunSuccessfully())

			By("ensuring they are healthy")
			checkNg(bottlerocketNodegroup)
			checkNg(ubuntuNodegroup)

			By("asserting node volumes")
			tests.AssertNodeVolumes(params.KubeconfigPath, params.Region, bottlerocketNodegroup, "/dev/sdb")

			By("correctly configuring the bottlerocket nodegroup")
			kubeTest, err := kube.NewTest(params.KubeconfigPath)
			Expect(err).NotTo(HaveOccurred())

			nodeList := kubeTest.ListNodes(metav1.ListOptions{
				LabelSelector: fmt.Sprintf("%s=%s", "eks.amazonaws.com/nodegroup", bottlerocketNodegroup),
			})
			Expect(nodeList.Items).NotTo(BeEmpty())
			for _, node := range nodeList.Items {
				Expect(node.Status.NodeInfo.OSImage).To(ContainSubstring("Bottlerocket"))
				// On k8s 1.26 and greater, --hostname-override flag is passed to kubelet to allow nodes to join cluster.
				// For further reference, please see: https://github.com/bottlerocket-os/bottlerocket/pull/3033
				Expect(node.Labels["kubernetes.io/hostname"]).To(ContainSubstring(fmt.Sprintf("%s.compute.internal", params.Region)))
			}
			kubeTest.Close()
		})

		It("supports a nodegroup with an update config", func() {
			updateConfig := &api.NodeGroupUpdateConfig{
				MaxUnavailable: aws.Int(2),
			}
			clusterConfig := makeClusterConfig()
			clusterConfig.ManagedNodeGroups = []*api.ManagedNodeGroup{
				{
					NodeGroupBase: &api.NodeGroupBase{
						Name: updateConfigNodegroup,
					},
					UpdateConfig: updateConfig,
				},
			}

			By("creating it")
			cmd := params.EksctlCreateCmd.
				WithArgs(
					"nodegroup",
					"--config-file", "-",
					"--verbose", "4",
				).
				WithoutArg("--region", params.Region).
				WithStdin(clusterutils.Reader(clusterConfig))
			Expect(cmd).To(RunSuccessfully())

			By("ensuring it is healthy")
			checkNg(updateConfigNodegroup)

			ctx := context.Background()
			clusterProvider, err := eks.New(ctx, &api.ProviderConfig{Region: params.Region}, clusterConfig)
			Expect(err).NotTo(HaveOccurred())
			ctl := clusterProvider.AWSProvider
			out, err := ctl.EKS().DescribeNodegroup(ctx, &awseks.DescribeNodegroupInput{
				ClusterName:   &params.ClusterName,
				NodegroupName: aws.String(updateConfigNodegroup),
			})
			Expect(err).NotTo(HaveOccurred())
			Eventually(out.Nodegroup.UpdateConfig.MaxUnavailable).Should(Equal(aws.Int32(2)))

			By("updating the nodegroup's UpdateConfig")
			clusterConfig.ManagedNodeGroups[0].Spot = true
			clusterConfig.ManagedNodeGroups[0].UpdateConfig = &api.NodeGroupUpdateConfig{
				MaxUnavailable: aws.Int(1),
			}

			cmd = params.EksctlUpdateCmd.
				WithArgs(
					"nodegroup",
					"--config-file", "-",
					"--verbose", "4",
				).
				WithoutArg("--region", params.Region).
				WithStdin(clusterutils.Reader(clusterConfig))

			msg := fmt.Sprintf("unchanged fields for nodegroup %s: the following fields remain unchanged; they are not supported by `eksctl update nodegroup`: Spot", updateConfigNodegroup)
			Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
				ContainElement(ContainSubstring(msg)),
			))

			out, err = ctl.EKS().DescribeNodegroup(ctx, &awseks.DescribeNodegroupInput{
				ClusterName:   &params.ClusterName,
				NodegroupName: aws.String(updateConfigNodegroup),
			})
			Expect(err).NotTo(HaveOccurred())
			Eventually(out.Nodegroup.UpdateConfig.MaxUnavailable).Should(Equal(aws.Int32(1)))
		})
	})

	Context("scaling the initial nodegroup", func() {
		It("should not return an error", func() {
			cmd := params.EksctlScaleNodeGroupCmd.WithArgs(
				"--cluster", params.ClusterName,
				"--nodes-min", "2",
				"--nodes", "3",
				"--nodes-max", "4",
				"--name", initialAl2Nodegroup,
			)
			Expect(cmd).To(RunSuccessfully())
		})
	})

	Context("eksctl utils update-cluster-vpc-config", Serial, func() {
		makeAWSProvider := func(ctx context.Context, clusterConfig *api.ClusterConfig) api.ClusterProvider {
			clusterProvider, err := eks.New(ctx, &api.ProviderConfig{Region: params.Region}, clusterConfig)
			Expect(err).NotTo(HaveOccurred())
			return clusterProvider.AWSProvider
		}
		getPrivateSubnetIDs := func(ctx context.Context, ec2API awsapi.EC2, vpcID string) []string {
			out, err := ec2API.DescribeSubnets(ctx, &ec2.DescribeSubnetsInput{
				Filters: []ec2types.Filter{
					{
						Name:   aws.String("vpc-id"),
						Values: []string{vpcID},
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			var subnetIDs []string
			for _, s := range out.Subnets {
				if !*s.MapPublicIpOnLaunch {
					subnetIDs = append(subnetIDs, *s.SubnetId)
				}
			}
			return subnetIDs
		}
		It("should update the VPC config", func() {
			clusterConfig := makeClusterConfig()
			ctx := context.Background()
			awsProvider := makeAWSProvider(ctx, clusterConfig)
			cluster, err := awsProvider.EKS().DescribeCluster(ctx, &awseks.DescribeClusterInput{
				Name: aws.String(params.ClusterName),
			})
			Expect(err).NotTo(HaveOccurred(), "error describing cluster")
			clusterSubnetIDs := getPrivateSubnetIDs(ctx, awsProvider.EC2(), *cluster.Cluster.ResourcesVpcConfig.VpcId)
			Expect(len(cluster.Cluster.ResourcesVpcConfig.SecurityGroupIds) > 0).To(BeTrue(), "at least one security group ID must be associated with the cluster")

			clusterVPC := &api.ClusterVPC{
				ClusterEndpoints: &api.ClusterEndpoints{
					PrivateAccess: api.Enabled(),
					PublicAccess:  api.Enabled(),
				},
				PublicAccessCIDRs:            []string{"127.0.0.1/32"},
				ControlPlaneSubnetIDs:        clusterSubnetIDs,
				ControlPlaneSecurityGroupIDs: []string{cluster.Cluster.ResourcesVpcConfig.SecurityGroupIds[0]},
			}
			By("accepting CLI options")
			cmd := params.EksctlUtilsCmd.WithArgs(
				"update-cluster-vpc-config",
				"--cluster", params.ClusterName,
				"--private-access",
				"--public-access",
				"--public-access-cidrs", strings.Join(clusterVPC.PublicAccessCIDRs, ","),
				"--control-plane-subnet-ids", strings.Join(clusterVPC.ControlPlaneSubnetIDs, ","),
				"--control-plane-security-group-ids", strings.Join(clusterVPC.ControlPlaneSecurityGroupIDs, ","),
				"-v4",
				"--approve",
			).
				WithTimeout(45 * time.Minute)
			session := cmd.Run()
			Expect(session.ExitCode()).To(Equal(0))

			formatWithClusterAndRegion := func(format string, values ...any) string {
				return fmt.Sprintf(format, append([]any{params.ClusterName, params.Region}, values...)...)
			}
			Expect(strings.Split(string(session.Buffer().Contents()), "\n")).To(ContainElements(
				ContainSubstring(formatWithClusterAndRegion("control plane subnets and security groups for cluster %q in %q have been updated to: "+
					"controlPlaneSubnetIDs=%v, controlPlaneSecurityGroupIDs=%v", clusterVPC.ControlPlaneSubnetIDs, clusterVPC.ControlPlaneSecurityGroupIDs)),
				ContainSubstring(formatWithClusterAndRegion("Kubernetes API endpoint access for cluster %q in %q has been updated to: privateAccess=%v, publicAccess=%v",
					*clusterVPC.ClusterEndpoints.PrivateAccess, *clusterVPC.ClusterEndpoints.PublicAccess)),
				ContainSubstring(formatWithClusterAndRegion("public access CIDRs for cluster %q in %q have been updated to: %v", clusterVPC.PublicAccessCIDRs)),
			))

			By("accepting a config file")
			clusterConfig.VPC = clusterVPC
			cmd = params.EksctlUtilsCmd.WithArgs(
				"update-cluster-vpc-config",
				"--config-file", "-",
				"-v4",
				"--approve",
			).
				WithoutArg("--region", params.Region).
				WithStdin(clusterutils.Reader(clusterConfig))
			session = cmd.Run()
			Expect(session.ExitCode()).To(Equal(0))
			Expect(strings.Split(string(session.Buffer().Contents()), "\n")).To(ContainElements(
				ContainSubstring(formatWithClusterAndRegion("Kubernetes API endpoint access for cluster %q in %q is already up-to-date")),
				ContainSubstring(formatWithClusterAndRegion("control plane subnet IDs for cluster %q in %q are already up-to-date")),
				ContainSubstring(formatWithClusterAndRegion("control plane security group IDs for cluster %q in %q are already up-to-date")),
			))

			By("resetting public access CIDRs")
			cmd = params.EksctlUtilsCmd.WithArgs(
				"update-cluster-vpc-config",
				"--cluster", params.ClusterName,
				"--public-access-cidrs", "0.0.0.0/0",
				"-v4",
				"--approve",
			)
			Expect(cmd).To(RunSuccessfully())
		})
	})
})

var _ = SynchronizedAfterSuite(func() {}, func() {
	params.DeleteClusters()
})
