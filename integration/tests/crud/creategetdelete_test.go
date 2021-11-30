//go:build integration
// +build integration

package crud

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/client-go/kubernetes"

	"github.com/aws/aws-sdk-go/aws"
	awsec2 "github.com/aws/aws-sdk-go/service/ec2"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	harness "github.com/dlespiau/kube-test-harness"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/yaml"

	. "github.com/weaveworks/eksctl/integration/matchers"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	clusterutils "github.com/weaveworks/eksctl/integration/utilities/cluster"
	"github.com/weaveworks/eksctl/integration/utilities/kube"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/iam"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/utils/file"
)

var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	if err := api.Register(); err != nil {
		panic(errors.Wrap(err, "unexpected error registering API scheme"))
	}
	params = tests.NewParams("crud")

}

func TestCRUD(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("(Integration) Create, Get, Scale & Delete", func() {

	const (
		mngNG1 = "mng-1"
		mngNG2 = "mng-2"

		unmNG1 = "unm-1"
		unmNG2 = "unm-2"
	)

	commonTimeout := 10 * time.Minute
	makeClusterConfig := func() *api.ClusterConfig {
		clusterConfig := api.NewClusterConfig()
		clusterConfig.Metadata.Name = params.ClusterName
		clusterConfig.Metadata.Region = params.Region
		clusterConfig.Metadata.Version = params.Version
		return clusterConfig
	}

	BeforeSuite(func() {
		params.KubeconfigTemp = false
		if params.KubeconfigPath == "" {
			wd, _ := os.Getwd()
			f, _ := os.CreateTemp(wd, "kubeconfig-")
			params.KubeconfigPath = f.Name()
			params.KubeconfigTemp = true
		}

		if params.SkipCreate {
			fmt.Fprintf(GinkgoWriter, "will use existing cluster %s", params.ClusterName)
			if !file.Exists(params.KubeconfigPath) {
				// Generate the Kubernetes configuration that eksctl create
				// would have generated otherwise:
				cmd := params.EksctlUtilsCmd.WithArgs(
					"write-kubeconfig",
					"--verbose", "4",
					"--cluster", params.ClusterName,
					"--kubeconfig", params.KubeconfigPath,
				)
				Expect(cmd).To(RunSuccessfully())
			}
			return
		}

		fmt.Fprintf(GinkgoWriter, "Using kubeconfig: %s\n", params.KubeconfigPath)

		cmd := params.EksctlCreateCmd.WithArgs(
			"cluster",
			"--verbose", "4",
			"--name", params.ClusterName,
			"--tags", "alpha.eksctl.io/description=eksctl integration test",
			"--nodegroup-name", mngNG1,
			"--node-labels", "ng-name="+mngNG1,
			"--nodes", "1",
			"--version", params.Version,
			"--kubeconfig", params.KubeconfigPath,
			"--zones", "us-west-2b,us-west-2c",
		)
		Expect(cmd).To(RunSuccessfully())
	})

	AfterSuite(func() {
		params.DeleteClusters()
		gexec.KillAndWait()
		if params.KubeconfigTemp {
			os.Remove(params.KubeconfigPath)
		}
		os.RemoveAll(params.TestDirectory)
	})

	Describe("cluster with 1 node", func() {
		It("should have created an EKS cluster and two CloudFormation stacks", func() {
			awsSession := NewSession(params.Region)

			Expect(awsSession).To(HaveExistingCluster(params.ClusterName, awseks.ClusterStatusActive, params.Version))

			Expect(awsSession).To(HaveExistingStack(fmt.Sprintf("eksctl-%s-cluster", params.ClusterName)))
			Expect(awsSession).To(HaveExistingStack(fmt.Sprintf("eksctl-%s-nodegroup-%s", params.ClusterName, mngNG1)))
		})

		It("should have created a valid kubectl config file", func() {
			config, err := clientcmd.LoadFromFile(params.KubeconfigPath)
			Expect(err).ShouldNot(HaveOccurred())

			err = clientcmd.ConfirmUsable(*config, "")
			Expect(err).ShouldNot(HaveOccurred())

			Expect(config.CurrentContext).To(ContainSubstring("eksctl"))
			Expect(config.CurrentContext).To(ContainSubstring(params.ClusterName))
			Expect(config.CurrentContext).To(ContainSubstring(params.Region))
		})

		Context("and listing clusters", func() {
			It("should return the previously created cluster", func() {
				cmd := params.EksctlGetCmd.WithArgs("clusters", "--all-regions")
				Expect(cmd).To(RunSuccessfullyWithOutputString(ContainSubstring(params.ClusterName)))
			})
		})

		Context("and describe the stack for the cluster", func() {
			It("should describe the cluster's stack", func() {
				cmd := params.EksctlUtilsCmd.WithArgs("describe-stacks", "--cluster", params.ClusterName, "-o", "yaml")
				session := cmd.Run()
				Expect(session.ExitCode()).To(BeZero())
				var stacks []*cloudformation.Stack
				Expect(yaml.Unmarshal(session.Out.Contents(), &stacks)).To(Succeed())
				Expect(stacks).To(HaveLen(2))
				nodegroupStack := stacks[0]
				clusterStack := stacks[1]
				Expect(aws.StringValue(clusterStack.StackName)).To(ContainSubstring(params.ClusterName))
				Expect(aws.StringValue(nodegroupStack.StackName)).To(ContainSubstring(params.ClusterName))
				Expect(aws.StringValue(clusterStack.Description)).To(Equal("EKS cluster (dedicated VPC: true, dedicated IAM: true) [created and managed by eksctl]"))
				Expect(aws.StringValue(nodegroupStack.Description)).To(Equal("EKS Managed Nodes (SSH access: false) [created by eksctl]"))
			})
		})

		Context("toggling kubernetes API access", func() {
			var (
				clientSet *kubernetes.Clientset
			)
			BeforeEach(func() {
				cfg := &api.ClusterConfig{
					Metadata: &api.ClusterMeta{
						Name:   params.ClusterName,
						Region: params.Region,
					},
				}
				ctl, err := eks.New(&api.ProviderConfig{Region: params.Region}, cfg)
				Expect(err).NotTo(HaveOccurred())
				err = ctl.RefreshClusterStatus(cfg)
				Expect(err).ShouldNot(HaveOccurred())
				clientSet, err = ctl.NewStdClientSet(cfg)
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("should be publicly accessible by default", func() {
				_, err := clientSet.CoreV1().ServiceAccounts(metav1.NamespaceDefault).List(context.TODO(), metav1.ListOptions{})
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("should be able to disable public access", func() {
				cmd := params.EksctlUtilsCmd.WithArgs(
					"set-public-access-cidrs",
					"--cluster", params.ClusterName,
					"1.1.1.1/32,2.2.2.0/24",
					"--approve",
				)
				Expect(cmd).To(RunSuccessfully())

				_, err := clientSet.CoreV1().ServiceAccounts(metav1.NamespaceDefault).List(context.TODO(), metav1.ListOptions{})
				Expect(err).Should(HaveOccurred())
			})

			It("should be able to re-enable public access", func() {
				cmd := params.EksctlUtilsCmd.WithArgs(
					"set-public-access-cidrs",
					"--cluster", params.ClusterName,
					"0.0.0.0/0",
					"--approve",
				)
				Expect(cmd).To(RunSuccessfully())

				_, err := clientSet.CoreV1().ServiceAccounts(metav1.NamespaceDefault).List(context.TODO(), metav1.ListOptions{})
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("and create a new nodegroup with taints and maxPods", func() {
			It("should have taints and maxPods set", func() {
				By("creating a new nodegroup with taints and maxPods set")
				cmd := params.EksctlCreateCmd.
					WithArgs(
						"nodegroup",
						"--config-file", "-",
						"--verbose", "4",
					).
					WithoutArg("--region", params.Region).
					WithStdin(clusterutils.ReaderFromFile(params.ClusterName, params.Region, "testdata/taints-max-pods.yaml"))
				Expect(cmd).To(RunSuccessfully())

				config, err := clientcmd.BuildConfigFromFlags("", params.KubeconfigPath)
				Expect(err).NotTo(HaveOccurred())
				clientset, err := kubernetes.NewForConfig(config)
				Expect(err).NotTo(HaveOccurred())

				By("asserting that both formats for taints are supported")
				var (
					nodeListN1 = tests.ListNodes(clientset, unmNG1)
					nodeListN2 = tests.ListNodes(clientset, unmNG2)
				)

				tests.AssertNodeTaints(nodeListN1, []corev1.Taint{
					{
						Key:    "key1",
						Value:  "val1",
						Effect: "NoSchedule",
					},
					{
						Key:    "key2",
						Effect: "NoExecute",
					},
				})

				tests.AssertNodeTaints(nodeListN2, []corev1.Taint{
					{
						Key:    "key1",
						Value:  "value1",
						Effect: "NoSchedule",
					},
					{
						Key:    "key2",
						Effect: "NoExecute",
					},
				})

				By("asserting that maxPods is set correctly")
				expectedMaxPods := 123
				for _, node := range nodeListN1.Items {
					maxPods, _ := node.Status.Allocatable.Pods().AsInt64()
					Expect(maxPods).To(Equal(int64(expectedMaxPods)))
				}

			})
		})

		Context("can add a nodegroup into a new subnet", func() {
			var (
				subnet        *awsec2.Subnet
				nodegroupName string
			)
			BeforeEach(func() {
				nodegroupName = "test-extra-nodegroup"
			})
			AfterEach(func() {
				cmd := params.EksctlDeleteCmd.WithArgs(
					"nodegroup",
					"--verbose", "4",
					"--cluster", params.ClusterName,
					"--wait",
					nodegroupName,
				)
				Expect(cmd).To(RunSuccessfully())
				awsSession := NewSession(params.Region)
				ec2 := awsec2.New(awsSession)
				output, err := ec2.DeleteSubnet(&awsec2.DeleteSubnetInput{
					SubnetId: subnet.SubnetId,
				})
				Expect(err).NotTo(HaveOccurred(), output.GoString())

			})
			It("creates a new nodegroup", func() {
				cfg := &api.ClusterConfig{
					Metadata: &api.ClusterMeta{
						Name:   params.ClusterName,
						Region: params.Region,
					},
				}
				ctl, err := eks.New(&api.ProviderConfig{Region: params.Region}, cfg)
				Expect(err).NotTo(HaveOccurred())
				cl, err := ctl.GetCluster(params.ClusterName)
				Expect(err).NotTo(HaveOccurred())
				awsSession := NewSession(params.Region)
				ec2 := awsec2.New(awsSession)
				existingSubnets, err := ec2.DescribeSubnets(&awsec2.DescribeSubnetsInput{
					SubnetIds: cl.ResourcesVpcConfig.SubnetIds,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(len(existingSubnets.Subnets) > 0).To(BeTrue())
				s := existingSubnets.Subnets[0]

				cidr := *s.CidrBlock
				var (
					i1, i2, i3, i4, ic int
				)
				fmt.Sscanf(cidr, "%d.%d.%d.%d/%d", &i1, &i2, &i3, &i4, &ic)
				cidr = fmt.Sprintf("%d.%d.%s.%d/%d", i1, i2, "255", i4, ic)

				var tags []*awsec2.Tag

				// filter aws: tags
				for _, t := range s.Tags {
					if !strings.HasPrefix(*t.Key, "aws:") {
						tags = append(tags, t)
					}
				}
				output, err := ec2.CreateSubnet(&awsec2.CreateSubnetInput{
					AvailabilityZone: aws.String("us-west-2a"),
					CidrBlock:        aws.String(cidr),
					TagSpecifications: []*awsec2.TagSpecification{
						{
							ResourceType: aws.String(awsec2.ResourceTypeSubnet),
							Tags:         tags,
						},
					},
					VpcId: s.VpcId,
				})
				Expect(err).NotTo(HaveOccurred())
				moutput, err := ec2.ModifySubnetAttribute(&awsec2.ModifySubnetAttributeInput{
					MapPublicIpOnLaunch: &awsec2.AttributeBooleanValue{
						Value: aws.Bool(true),
					},
					SubnetId: output.Subnet.SubnetId,
				})
				Expect(err).NotTo(HaveOccurred(), moutput.GoString())
				subnet = output.Subnet

				routeTables, err := ec2.DescribeRouteTables(&awsec2.DescribeRouteTablesInput{
					Filters: []*awsec2.Filter{
						{
							Name:   aws.String("association.subnet-id"),
							Values: aws.StringSlice([]string{*s.SubnetId}),
						},
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(len(routeTables.RouteTables) > 0).To(BeTrue(), fmt.Sprintf("route table ended up being empty: %+v", routeTables))
				routput, err := ec2.AssociateRouteTable(&awsec2.AssociateRouteTableInput{
					RouteTableId: routeTables.RouteTables[0].RouteTableId,
					SubnetId:     subnet.SubnetId,
				})
				Expect(err).NotTo(HaveOccurred(), routput)

				// create a new subnet in that given vpc and zone.
				cmd := params.EksctlCreateCmd.WithArgs(
					"nodegroup",
					"--timeout=45m",
					"--cluster", params.ClusterName,
					"--nodes", "1",
					"--node-type", "p2.xlarge",
					"--subnet-ids", *subnet.SubnetId,
					nodegroupName,
				)
				Expect(cmd).To(RunSuccessfully())
			})
		})

		Context("and creating a nodegroup with containerd runtime", func() {
			var (
				nodegroupName string
			)
			BeforeEach(func() {
				nodegroupName = "test-containerd"
			})
			AfterEach(func() {
				cmd := params.EksctlDeleteCmd.WithArgs(
					"nodegroup",
					"--verbose", "4",
					"--cluster", params.ClusterName,
					"--wait",
					nodegroupName,
				)
				Expect(cmd).To(RunSuccessfully())
			})
			It("should create the nodegroup without problems", func() {
				clusterConfig := makeClusterConfig()
				clusterConfig.Metadata.Name = params.ClusterName
				clusterConfig.NodeGroups = []*api.NodeGroup{
					{
						NodeGroupBase: &api.NodeGroupBase{
							Name:         "test-containerd",
							AMIFamily:    api.NodeImageFamilyAmazonLinux2,
							InstanceType: "p2.xlarge",
						},
						ContainerRuntime: aws.String(api.ContainerRuntimeContainerD),
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
				Expect(cmd).To(RunSuccessfully())
			})
		})

		When("scaling nodegroup(s)", func() {
			It("should scale a single nodegroup", func() {
				By("passing the name of the nodegroup as a flag")
				cmd := params.EksctlScaleNodeGroupCmd.WithArgs(
					"--cluster", params.ClusterName,
					"--nodes-min", "4",
					"--nodes", "4",
					"--nodes-max", "4",
					"--name", mngNG1,
				)
				Expect(cmd).To(RunSuccessfully())

				getMngNgCmd := params.EksctlGetCmd.WithArgs(
					"nodegroup",
					"--cluster", params.ClusterName,
					"--name", mngNG1,
					"-o", "yaml",
				)
				Expect(getMngNgCmd).To(RunSuccessfullyWithOutputString(ContainSubstring("MaxSize: 4")))
				Expect(getMngNgCmd).To(RunSuccessfullyWithOutputString(ContainSubstring("MinSize: 4")))
				Expect(getMngNgCmd).To(RunSuccessfullyWithOutputString(ContainSubstring("DesiredCapacity: 4")))
			})

			It("should scale all nodegroups", func() {
				By("scaling all nodegroups in the config file to the desired capacity, max size, and min size")
				cmd := params.EksctlScaleNodeGroupCmd.WithArgs(
					"--config-file", "-",
				).
					WithoutArg("--region", params.Region).
					WithStdin(clusterutils.ReaderFromFile(params.ClusterName, params.Region, "testdata/scale-nodegroups.yaml"))
				Expect(cmd).To(RunSuccessfully())

				getMngNgCmd := params.EksctlGetCmd.WithArgs(
					"nodegroup",
					"--cluster", params.ClusterName,
					"--name", mngNG1,
					"-o", "yaml",
				)
				Expect(getMngNgCmd).To(RunSuccessfullyWithOutputString(ContainSubstring("MaxSize: 5")))
				Expect(getMngNgCmd).To(RunSuccessfullyWithOutputString(ContainSubstring("MinSize: 5")))
				Expect(getMngNgCmd).To(RunSuccessfullyWithOutputString(ContainSubstring("DesiredCapacity: 5")))

				getUnmNgCmd := params.EksctlGetCmd.WithArgs(
					"nodegroup",
					"--cluster", params.ClusterName,
					"--name", unmNG1,
					"-o", "yaml",
				)
				Expect(getUnmNgCmd).To(RunSuccessfullyWithOutputString(ContainSubstring("MaxSize: 5")))
				Expect(getUnmNgCmd).To(RunSuccessfullyWithOutputString(ContainSubstring("MinSize: 5")))
				Expect(getUnmNgCmd).To(RunSuccessfullyWithOutputString(ContainSubstring("DesiredCapacity: 5")))
			})
		})

		Context("and add a second (GPU) nodegroup", func() {
			It("should not return an error", func() {
				cmd := params.EksctlCreateCmd.WithArgs(
					"nodegroup",
					"--timeout=45m",
					"--cluster", params.ClusterName,
					"--nodes", "1",
					"--node-type", "p2.xlarge",
					"--node-private-networking",
					"--node-zones", "us-west-2b,us-west-2c",
					mngNG2,
				)
				Expect(cmd).To(RunSuccessfully())
			})

			It("should be able to list nodegroups", func() {
				cmd := params.EksctlGetCmd.WithArgs(
					"nodegroup",
					"-o", "json",
					"--cluster", params.ClusterName,
					mngNG1,
				)
				Expect(cmd).To(RunSuccessfullyWithOutputString(BeNodeGroupsWithNamesWhich(
					HaveLen(1),
					ContainElement(mngNG1),
					Not(ContainElement(mngNG2)),
				)))
				Expect(cmd).To(RunSuccessfullyWithOutputString(ContainSubstring(params.Version)))

				cmd = params.EksctlGetCmd.WithArgs(
					"nodegroup",
					"-o", "json",
					"--cluster", params.ClusterName,
					mngNG2,
				)
				Expect(cmd).To(RunSuccessfullyWithOutputString(BeNodeGroupsWithNamesWhich(
					HaveLen(1),
					ContainElement(mngNG2),
					Not(ContainElement(mngNG1)),
				)))
				Expect(cmd).To(RunSuccessfullyWithOutputString(ContainSubstring(params.Version)))

				cmd = params.EksctlGetCmd.WithArgs(
					"nodegroup",
					"-o", "json",
					"--cluster", params.ClusterName,
				)
				Expect(cmd).To(RunSuccessfullyWithOutputString(BeNodeGroupsWithNamesWhich(
					HaveLen(4),
					ContainElement(mngNG1),
					ContainElement(mngNG2),
					ContainElement(unmNG1),
					ContainElement(unmNG2),
				)))
				Expect(cmd).To(RunSuccessfullyWithOutputString(ContainSubstring(params.Version)))
			})

			Context("toggle CloudWatch logging", func() {
				var (
					cfg *api.ClusterConfig
					ctl *eks.ClusterProvider
				)

				BeforeEach(func() {
					cfg = &api.ClusterConfig{
						Metadata: &api.ClusterMeta{
							Name:   params.ClusterName,
							Region: params.Region,
						},
					}
					var err error
					ctl, err = eks.New(&api.ProviderConfig{Region: params.Region}, cfg)
					Expect(err).NotTo(HaveOccurred())
				})

				It("should have all types disabled by default", func() {
					enabled, disable, err := ctl.GetCurrentClusterConfigForLogging(cfg)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(enabled.List()).To(HaveLen(0))
					Expect(disable.List()).To(HaveLen(5))
				})

				It("should plan to enable two of the types using flags", func() {
					cmd := params.EksctlUtilsCmd.WithArgs(
						"update-cluster-logging",
						"--cluster", params.ClusterName,
						"--enable-types", "api,controllerManager",
					)
					Expect(cmd).To(RunSuccessfully())

					enabled, disable, err := ctl.GetCurrentClusterConfigForLogging(cfg)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(enabled.List()).To(HaveLen(0))
					Expect(disable.List()).To(HaveLen(5))
				})

				It("should enable two of the types using flags", func() {
					cmd := params.EksctlUtilsCmd.WithArgs(
						"update-cluster-logging",
						"--cluster", params.ClusterName,
						"--approve",
						"--enable-types", "api,controllerManager",
					)
					Expect(cmd).To(RunSuccessfully())

					enabled, disable, err := ctl.GetCurrentClusterConfigForLogging(cfg)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(enabled.List()).To(HaveLen(2))
					Expect(disable.List()).To(HaveLen(3))
					Expect(enabled.List()).To(ConsistOf("api", "controllerManager"))
					Expect(disable.List()).To(ConsistOf("audit", "authenticator", "scheduler"))
				})

				It("should enable all of the types with --enable-types=all", func() {
					cmd := params.EksctlUtilsCmd.WithArgs(
						"update-cluster-logging",
						"--cluster", params.ClusterName,
						"--approve",
						"--enable-types", "all",
					)
					Expect(cmd).To(RunSuccessfully())

					enabled, disable, err := ctl.GetCurrentClusterConfigForLogging(cfg)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(enabled.List()).To(HaveLen(5))
					Expect(disable.List()).To(HaveLen(0))
				})

				It("should enable all but one type", func() {
					cmd := params.EksctlUtilsCmd.WithArgs(
						"update-cluster-logging",
						"--cluster", params.ClusterName,
						"--approve",
						"--enable-types", "all",
						"--disable-types", "controllerManager",
					)
					Expect(cmd).To(RunSuccessfully())

					enabled, disable, err := ctl.GetCurrentClusterConfigForLogging(cfg)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(enabled.List()).To(HaveLen(4))
					Expect(disable.List()).To(HaveLen(1))
					Expect(enabled.List()).To(ConsistOf("api", "audit", "authenticator", "scheduler"))
					Expect(disable.List()).To(ConsistOf("controllerManager"))
				})

				It("should disable all but one type", func() {
					cmd := params.EksctlUtilsCmd.WithArgs(
						"update-cluster-logging",
						"--cluster", params.ClusterName,
						"--approve",
						"--disable-types", "all",
						"--enable-types", "controllerManager",
					)
					Expect(cmd).To(RunSuccessfully())

					enabled, disable, err := ctl.GetCurrentClusterConfigForLogging(cfg)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(disable.List()).To(HaveLen(4))
					Expect(enabled.List()).To(HaveLen(1))
					Expect(disable.List()).To(ConsistOf("api", "audit", "authenticator", "scheduler"))
					Expect(enabled.List()).To(ConsistOf("controllerManager"))
				})

				It("should disable all of the types with --disable-types=all", func() {
					cmd := params.EksctlUtilsCmd.WithArgs(
						"update-cluster-logging",
						"--cluster", params.ClusterName,
						"--approve",
						"--disable-types", "all",
					)
					Expect(cmd).To(RunSuccessfully())

					enabled, disable, err := ctl.GetCurrentClusterConfigForLogging(cfg)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(enabled.List()).To(HaveLen(0))
					Expect(disable.List()).To(HaveLen(5))
					Expect(disable.HasAll(api.SupportedCloudWatchClusterLogTypes()...)).To(BeTrue())
				})
			})

			Context("create, update, and delete iamserviceaccounts", func() {
				var (
					cfg  *api.ClusterConfig
					ctl  *eks.ClusterProvider
					oidc *iamoidc.OpenIDConnectManager
					err  error
				)

				BeforeEach(func() {
					cfg = &api.ClusterConfig{
						Metadata: &api.ClusterMeta{
							Name:   params.ClusterName,
							Region: params.Region,
						},
					}
					ctl, err = eks.New(&api.ProviderConfig{Region: params.Region}, cfg)
					Expect(err).NotTo(HaveOccurred())
					err = ctl.RefreshClusterStatus(cfg)
					Expect(err).ShouldNot(HaveOccurred())
					oidc, err = ctl.NewOpenIDConnectManager(cfg)
					Expect(err).ShouldNot(HaveOccurred())
				})

				It("should enable OIDC, create two iamserviceaccounts and update the policies", func() {
					{
						exists, err := oidc.CheckProviderExists()
						Expect(err).ShouldNot(HaveOccurred())
						Expect(exists).To(BeFalse())
					}

					setupCmd := params.EksctlUtilsCmd.WithArgs(
						"associate-iam-oidc-provider",
						"--cluster", params.ClusterName,
						"--approve",
					)
					Expect(setupCmd).To(RunSuccessfully())

					{
						exists, err := oidc.CheckProviderExists()
						Expect(err).ShouldNot(HaveOccurred())
						Expect(exists).To(BeTrue())
					}

					cmds := []Cmd{
						params.EksctlCreateCmd.WithArgs(
							"iamserviceaccount",
							"--cluster", params.ClusterName,
							"--name", "app-cache-access",
							"--namespace", "app1",
							"--attach-policy-arn", "arn:aws:iam::aws:policy/AmazonDynamoDBReadOnlyAccess",
							"--attach-policy-arn", "arn:aws:iam::aws:policy/AmazonElastiCacheFullAccess",
							"--approve",
						),
						params.EksctlCreateCmd.WithArgs(
							"iamserviceaccount",
							"--cluster", params.ClusterName,
							"--name", "s3-read-only",
							"--attach-policy-arn", "arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess",
							"--approve",
						),
					}

					Expect(cmds).To(RunSuccessfully())

					awsSession := NewSession(params.Region)

					stackNamePrefix := fmt.Sprintf("eksctl-%s-addon-iamserviceaccount-", params.ClusterName)

					Expect(awsSession).To(HaveExistingStack(stackNamePrefix + "default-s3-read-only"))
					Expect(awsSession).To(HaveExistingStack(stackNamePrefix + "app1-app-cache-access"))

					clientSet, err := ctl.NewStdClientSet(cfg)
					Expect(err).ShouldNot(HaveOccurred())

					{
						sa, err := clientSet.CoreV1().ServiceAccounts(metav1.NamespaceDefault).Get(context.TODO(), "s3-read-only", metav1.GetOptions{})
						Expect(err).ShouldNot(HaveOccurred())

						Expect(sa.Annotations).To(HaveLen(1))
						Expect(sa.Annotations).To(HaveKey(api.AnnotationEKSRoleARN))
						Expect(sa.Annotations[api.AnnotationEKSRoleARN]).To(MatchRegexp("^arn:aws:iam::.*:role/eksctl-" + truncate(params.ClusterName) + ".*$"))
					}

					{
						sa, err := clientSet.CoreV1().ServiceAccounts("app1").Get(context.TODO(), "app-cache-access", metav1.GetOptions{})
						Expect(err).ShouldNot(HaveOccurred())

						Expect(sa.Annotations).To(HaveLen(1))
						Expect(sa.Annotations).To(HaveKey(api.AnnotationEKSRoleARN))
						Expect(sa.Annotations[api.AnnotationEKSRoleARN]).To(MatchRegexp("^arn:aws:iam::.*:role/eksctl-" + truncate(params.ClusterName) + ".*$"))
					}

					cmds = []Cmd{
						params.EksctlUpdateCmd.WithArgs(
							"iamserviceaccount",
							"--cluster", params.ClusterName,
							"--name", "app-cache-access",
							"--namespace", "app1",
							"--attach-policy-arn", "arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess",
							"--approve",
						),
					}

					Expect(cmds).To(RunSuccessfully())
				})

				It("should list both iamserviceaccounts", func() {
					cmd := params.EksctlGetCmd.WithArgs(
						"iamserviceaccount",
						"--cluster", params.ClusterName,
					)

					Expect(cmd).To(RunSuccessfullyWithOutputString(MatchRegexp(
						`(?m:^NAMESPACE\s+NAME\s+ROLE\sARN$)` +
							`|(?m:^app1\s+app-cache-access\s+arn:aws:iam::.*$)` +
							`|(?m:^default\s+s3-read-only\s+arn:aws:iam::.*$)`,
					)))
				})

				It("delete both iamserviceaccounts", func() {
					cmds := []Cmd{
						params.EksctlDeleteCmd.WithArgs(
							"iamserviceaccount",
							"--cluster", params.ClusterName,
							"--name", "s3-read-only",
							"--wait",
						),
						params.EksctlDeleteCmd.WithArgs(
							"iamserviceaccount",
							"--cluster", params.ClusterName,
							"--name", "app-cache-access",
							"--namespace", "app1",
							"--wait",
						),
					}
					Expect(cmds).To(RunSuccessfully())

					awsSession := NewSession(params.Region)

					stackNamePrefix := fmt.Sprintf("eksctl-%s-addon-iamserviceaccount-", params.ClusterName)

					Expect(awsSession).NotTo(HaveExistingStack(stackNamePrefix + "default-s3-read-only"))
					Expect(awsSession).NotTo(HaveExistingStack(stackNamePrefix + "app1-app-cache-access"))
				})
			})

			Context("create test workloads", func() {
				var (
					err  error
					test *harness.Test
				)

				BeforeEach(func() {
					test, err = kube.NewTest(params.KubeconfigPath)
					Expect(err).ShouldNot(HaveOccurred())
				})

				AfterEach(func() {
					test.Close()
					Eventually(func() int {
						return len(test.ListPods(test.Namespace, metav1.ListOptions{}).Items)
					}, "3m", "1s").Should(BeZero())
				})

				It("should deploy podinfo service to the cluster and access it via proxy", func() {
					d := test.CreateDeploymentFromFile(test.Namespace, "../../data/podinfo.yaml")
					test.WaitForDeploymentReady(d, commonTimeout)

					pods := test.ListPodsFromDeployment(d)
					Expect(len(pods.Items)).To(Equal(2))

					// For each pod of the Deployment, check we receive a sensible response to a
					// GET request on /version.
					for _, pod := range pods.Items {
						Expect(pod.Namespace).To(Equal(test.Namespace))

						req := test.PodProxyGet(&pod, "", "/version")
						fmt.Fprintf(GinkgoWriter, "url = %#v", req.URL())

						var js interface{}
						test.PodProxyGetJSON(&pod, "", "/version", &js)

						Expect(js.(map[string]interface{})).To(HaveKeyWithValue("version", "1.5.1"))
					}
				})

				It("should have functional DNS", func() {
					d := test.CreateDaemonSetFromFile(test.Namespace, "../../data/test-dns.yaml")

					test.WaitForDaemonSetReady(d, commonTimeout)

					{
						ds, err := test.GetDaemonSet(test.Namespace, d.Name)
						Expect(err).ShouldNot(HaveOccurred())
						fmt.Fprintf(GinkgoWriter, "ds.Status = %#v", ds.Status)
					}
				})

				It("should have access to HTTP(S) sites", func() {
					d := test.CreateDaemonSetFromFile(test.Namespace, "../../data/test-http.yaml")

					test.WaitForDaemonSetReady(d, commonTimeout)

					{
						ds, err := test.GetDaemonSet(test.Namespace, d.Name)
						Expect(err).ShouldNot(HaveOccurred())
						fmt.Fprintf(GinkgoWriter, "ds.Status = %#v", ds.Status)
					}
				})

				It("should be able to run pods with an iamserviceaccount", func() {
					createCmd := params.EksctlCreateCmd.WithArgs(
						"iamserviceaccount",
						"--cluster", params.ClusterName,
						"--name", "s3-reader",
						"--namespace", test.Namespace,
						"--attach-policy-arn", "arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess",
						"--approve",
					)

					Expect(createCmd).To(RunSuccessfully())

					d := test.CreateDeploymentFromFile(test.Namespace, "../../data/iamserviceaccount-checker.yaml")
					test.WaitForDeploymentReady(d, commonTimeout)

					pods := test.ListPodsFromDeployment(d)
					Expect(len(pods.Items)).To(Equal(2))

					// For each pod of the Deployment, check we get expected environment variables
					// via a GET request on /env.
					type sessionObject struct {
						AssumedRoleUser struct {
							AssumedRoleID, Arn string
						}
						Audience, Provider, SubjectFromWebIdentityToken string
						Credentials                                     struct {
							SecretAccessKey, SessionToken, Expiration, AccessKeyID string
						}
					}

					for _, pod := range pods.Items {
						Expect(pod.Namespace).To(Equal(test.Namespace))

						so := sessionObject{}

						var js []string
						test.PodProxyGetJSON(&pod, "", "/env", &js)

						Expect(js).To(ContainElement(HavePrefix("AWS_ROLE_ARN=arn:aws:iam::")))
						Expect(js).To(ContainElement("AWS_WEB_IDENTITY_TOKEN_FILE=/var/run/secrets/eks.amazonaws.com/serviceaccount/token"))
						Expect(js).To(ContainElement(HavePrefix("AWS_SESSION_OBJECT=")))

						for _, envVar := range js {
							if strings.HasPrefix(envVar, "AWS_SESSION_OBJECT=") {
								err := json.Unmarshal([]byte(strings.TrimPrefix(envVar, "AWS_SESSION_OBJECT=")), &so)
								Expect(err).ShouldNot(HaveOccurred())
							}
						}

						Expect(so.AssumedRoleUser.AssumedRoleID).To(HaveSuffix(":integration-test"))

						Expect(so.AssumedRoleUser.Arn).To(MatchRegexp("^arn:aws:sts::.*:assumed-role/eksctl-" + truncate(params.ClusterName) + "-.*/integration-test$"))

						Expect(so.Audience).To(Equal("sts.amazonaws.com"))

						Expect(so.Provider).To(MatchRegexp("^arn:aws:iam::.*:oidc-provider/oidc.eks." + params.Region + ".amazonaws.com/id/.*$"))

						Expect(so.SubjectFromWebIdentityToken).To(Equal("system:serviceaccount:" + test.Namespace + ":s3-reader"))

						Expect(so.Credentials.SecretAccessKey).NotTo(BeEmpty())
						Expect(so.Credentials.SessionToken).NotTo(BeEmpty())
						Expect(so.Credentials.Expiration).NotTo(BeEmpty())
						Expect(so.Credentials.AccessKeyID).NotTo(BeEmpty())
					}

					deleteCmd := params.EksctlDeleteCmd.WithArgs(
						"iamserviceaccount",
						"--cluster", params.ClusterName,
						"--name", "s3-reader",
						"--namespace", test.Namespace,
					)

					Expect(deleteCmd).To(RunSuccessfully())
				})
			})

			Context("and manipulating iam identity mappings", func() {
				var (
					expR0, expR1, expU0 string
					role0, role1        iam.Identity
					user0               iam.Identity
					admin               = "admin"
					alice               = "alice"
				)

				BeforeEach(func() {
					roleCanonicalArn := "arn:aws:iam::123456:role/eksctl-testing-XYZ"
					var err error
					role0 = iam.RoleIdentity{
						RoleARN: roleCanonicalArn,
						KubernetesIdentity: iam.KubernetesIdentity{
							KubernetesUsername: admin,
							KubernetesGroups:   []string{"system:masters", "system:nodes"},
						},
					}
					role1 = iam.RoleIdentity{
						RoleARN: roleCanonicalArn,
						KubernetesIdentity: iam.KubernetesIdentity{
							KubernetesGroups: []string{"system:something"},
						},
					}

					userCanonicalArn := "arn:aws:iam::123456:user/alice"

					user0 = iam.UserIdentity{
						UserARN: userCanonicalArn,
						KubernetesIdentity: iam.KubernetesIdentity{
							KubernetesUsername: alice,
							KubernetesGroups:   []string{"system:masters", "cryptographers"},
						},
					}

					bs, err := yaml.Marshal([]iam.Identity{role0})
					Expect(err).ShouldNot(HaveOccurred())
					expR0 = string(bs)

					bs, err = yaml.Marshal([]iam.Identity{role1})
					Expect(err).ShouldNot(HaveOccurred())
					expR1 = string(bs)

					bs, err = yaml.Marshal([]iam.Identity{user0})
					Expect(err).ShouldNot(HaveOccurred())
					expU0 = string(bs)
				})

				It("fails getting unknown role mapping", func() {
					cmd := params.EksctlGetCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", params.ClusterName,
						"--arn", "arn:aws:iam::123456:role/idontexist",
						"-o", "yaml",
					)
					Expect(cmd).NotTo(RunSuccessfully())
				})
				It("fails getting unknown user mapping", func() {
					cmd := params.EksctlGetCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", params.ClusterName,
						"--arn", "arn:aws:iam::123456:user/bob",
						"-o", "yaml",
					)
					Expect(cmd).NotTo(RunSuccessfully())
				})
				It("creates role mapping", func() {
					create := params.EksctlCreateCmd.WithArgs(
						"iamidentitymapping",
						"--name", params.ClusterName,
						"--arn", role0.ARN(),
						"--username", role0.Username(),
						"--group", role0.Groups()[0],
						"--group", role0.Groups()[1],
					)
					Expect(create).To(RunSuccessfully())

					get := params.EksctlGetCmd.WithArgs(
						"iamidentitymapping",
						"--name", params.ClusterName,
						"--arn", role0.ARN(),
						"-o", "yaml",
					)
					Expect(get).To(RunSuccessfullyWithOutputString(MatchYAML(expR0)))
				})
				It("creates user mapping", func() {
					create := params.EksctlCreateCmd.WithArgs(
						"iamidentitymapping",
						"--name", params.ClusterName,
						"--arn", user0.ARN(),
						"--username", user0.Username(),
						"--group", user0.Groups()[0],
						"--group", user0.Groups()[1],
					)
					Expect(create).To(RunSuccessfully())

					get := params.EksctlGetCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", params.ClusterName,
						"--arn", user0.ARN(),
						"-o", "yaml",
					)
					Expect(get).To(RunSuccessfullyWithOutputString(MatchYAML(expU0)))
				})
				It("creates a duplicate role mapping", func() {
					createRole := params.EksctlCreateCmd.WithArgs(
						"iamidentitymapping",
						"--name", params.ClusterName,
						"--arn", role0.ARN(),
						"--username", role0.Username(),
						"--group", role0.Groups()[0],
						"--group", role0.Groups()[1],
					)
					Expect(createRole).To(RunSuccessfully())

					get := params.EksctlGetCmd.WithArgs(
						"iamidentitymapping",
						"--name", params.ClusterName,
						"--arn", role0.ARN(),
						"-o", "yaml",
					)
					Expect(get).To(RunSuccessfullyWithOutputString(MatchYAML(expR0 + expR0)))
				})
				It("creates a duplicate user mapping", func() {
					createCmd := params.EksctlCreateCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", params.ClusterName,
						"--arn", user0.ARN(),
						"--username", user0.Username(),
						"--group", user0.Groups()[0],
						"--group", user0.Groups()[1],
					)
					Expect(createCmd).To(RunSuccessfully())

					getCmd := params.EksctlGetCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", params.ClusterName,
						"--arn", user0.ARN(),
						"-o", "yaml",
					)
					Expect(getCmd).To(RunSuccessfullyWithOutputString(MatchYAML(expU0 + expU0)))
				})
				It("creates a duplicate role mapping with different identity", func() {
					createCmd := params.EksctlCreateCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", params.ClusterName,
						"--arn", role1.ARN(),
						"--group", role1.Groups()[0],
					)
					Expect(createCmd).To(RunSuccessfully())

					getCmd := params.EksctlGetCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", params.ClusterName,
						"--arn", role1.ARN(),
						"-o", "yaml",
					)
					Expect(getCmd).To(RunSuccessfullyWithOutputString(MatchYAML(expR0 + expR0 + expR1)))
				})
				It("deletes a single role mapping fifo", func() {
					deleteCmd := params.EksctlDeleteCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", params.ClusterName,
						"--arn", role1.ARN(),
					)
					Expect(deleteCmd).To(RunSuccessfully())

					getCmd := params.EksctlGetCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", params.ClusterName,
						"--arn", role1.ARN(),
						"-o", "yaml",
					)
					Expect(getCmd).To(RunSuccessfullyWithOutputString(MatchYAML(expR0 + expR1)))
				})
				It("fails when deleting unknown mapping", func() {
					deleteCmd := params.EksctlDeleteCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", params.ClusterName,
						"--arn", "arn:aws:iam::123456:role/idontexist",
					)
					Expect(deleteCmd).NotTo(RunSuccessfully())
				})
				It("deletes duplicate role mappings with --all", func() {
					deleteCmd := params.EksctlDeleteCmd.WithArgs(
						"iamidentitymapping",
						"--name", params.ClusterName,
						"--arn", role1.ARN(),
						"--all",
					)
					Expect(deleteCmd).To(RunSuccessfully())

					getCmd := params.EksctlGetCmd.WithArgs(
						"iamidentitymapping",
						"--name", params.ClusterName,
						"--arn", role1.ARN(),
						"-o", "yaml",
					)
					Expect(getCmd).NotTo(RunSuccessfully())
				})
				It("deletes duplicate user mappings with --all", func() {
					deleteCmd := params.EksctlDeleteCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", params.ClusterName,
						"--arn", user0.ARN(),
						"--all",
					)
					Expect(deleteCmd).To(RunSuccessfully())

					getCmd := params.EksctlGetCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", params.ClusterName,
						"--arn", user0.ARN(),
						"-o", "yaml",
					)
					Expect(getCmd).NotTo(RunSuccessfully())
				})
			})

			Context("and delete the second nodegroup", func() {
				It("should not return an error", func() {
					cmd := params.EksctlDeleteCmd.WithArgs(
						"nodegroup",
						"--verbose", "4",
						"--cluster", params.ClusterName,
						mngNG2,
					)
					Expect(cmd).To(RunSuccessfully())
				})
			})
		})

		Context("and scale the initial nodegroup back to 1 node", func() {
			It("should not return an error", func() {
				cmd := params.EksctlScaleNodeGroupCmd.WithArgs(
					"--cluster", params.ClusterName,
					"--nodes-min", "1",
					"--nodes", "1",
					"--nodes-max", "1",
					"--name", mngNG1,
				)
				Expect(cmd).To(RunSuccessfully())
			})
		})

		Context("and drain the initial nodegroup", func() {
			It("should not return an error", func() {
				cmd := params.EksctlDrainNodeGroupCmd.WithArgs(
					"--cluster", params.ClusterName,
					"--name", mngNG1,
				)
				Expect(cmd).To(RunSuccessfully())
			})
		})

		Context("and deleting the cluster", func() {
			It("should not return an error", func() {
				if params.SkipDelete {
					Skip("will not delete cluster " + params.ClusterName)
				}

				cmd := params.EksctlDeleteClusterCmd.WithArgs(
					"--name", params.ClusterName,
				)
				Expect(cmd).To(RunSuccessfully())
			})
		})
	})
})

func truncate(clusterName string) string {
	// CloudFormation seems to truncate long cluster names at 37 characters:
	if len(clusterName) > 37 {
		return clusterName[:37]
	}
	return clusterName
}
