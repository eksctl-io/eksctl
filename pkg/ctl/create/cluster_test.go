package create

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils/filter"
	"github.com/weaveworks/eksctl/pkg/eks"
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
			Entry("with --name option with invalid characters that are rejected by cloudformation", invalidParamsCase{
				args:  []string{"test-k8_cluster01"},
				error: "validation for test-k8_cluster01 failed, name must satisfy regular expression pattern: [a-zA-Z][-a-zA-Z0-9]*",
			}),
			Entry("with cluster name argument with invalid characters that are rejected by cloudformation", invalidParamsCase{
				args:  []string{"--name", "eksctl-testing-k_8_cluster01"},
				error: "validation for eksctl-testing-k_8_cluster01 failed, name must satisfy regular expression pattern: [a-zA-Z][-a-zA-Z0-9]*",
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
			Entry("with --name option with invalid characters that are rejected by cloudformation", invalidParamsCase{
				args:  []string{"test-k8_cluster01"},
				error: "validation for test-k8_cluster01 failed, name must satisfy regular expression pattern: [a-zA-Z][-a-zA-Z0-9]*",
			}),
			Entry("with cluster name argument with invalid characters that are rejected by cloudformation", invalidParamsCase{
				args:  []string{"--name", "eksctl-testing-k_8_cluster01"},
				error: "validation for eksctl-testing-k_8_cluster01 failed, name must satisfy regular expression pattern: [a-zA-Z][-a-zA-Z0-9]*",
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
	FDescribe("createClusterCmd", func() {
		Context("create cluster", func() {
			BeforeEach(func() {
			})
			AfterEach(func() {
			})
			It("can create a cluster", func() {

				// Context deadline exceeded because
				// waitTimeout
				cfg := api.NewClusterConfig()
				cfg.Metadata.Name = "gb-test-cluster-1"
				cfg.VPC.ClusterEndpoints = api.ClusterEndpointAccessDefaults()
				cfg.Metadata.Version = "1.22"
				cmd := &cmdutils.Cmd{
					ClusterConfig: cfg,
					ProviderConfig: api.ProviderConfig{
						WaitTimeout: time.Minute * 1,
					},
				}
				ngFilter := &filter.NodeGroupFilter{}
				params := &cmdutils.CreateClusterCmdParams{
					Subnets: map[api.SubnetTopology]*[]string{
						api.SubnetTopologyPrivate: {},
						api.SubnetTopologyPublic:  {},
					},
				}
				p := mockprovider.NewMockProvider()
				ctl := &eks.ClusterProvider{
					AWSProvider: p,
					Status: &eks.ProviderStatus{
						ClusterInfo: &eks.ClusterInfo{
							Cluster: testutils.NewFakeCluster("my-cluster", ""),
						},
					},
				}
				Expect(doCreateCluster(cmd, ngFilter, params, ctl)).To(Succeed())
			})
		})
	})
})

// responseHandler returns the required response based on the body and Request type in the body.
// TODO: Extend this to allow to return an error response for a specific
// response.
func responseHandler(body string) (string, error) {
	getContent := func(file string) string {
		content, err := os.ReadFile(filepath.Join("testdata", file))
		Expect(err).NotTo(HaveOccurred())
		return string(content)
	}
	switch {
	case strings.Contains(body, "Action=DescribeAvailabilityZones"):
		return getContent("describe_availability_zones.xml"), nil
	case strings.Contains(body, "Action=GetCallerIdentity"):
		return getContent("caller_identity_success_response.xml"), nil
	case strings.Contains(body, "Action=CreateStack"):
		return getContent("create_stack_success_response.xml"), nil
	case strings.Contains(body, "Action=DescribeStacks"):
		return getContent("describe_stacks_response.xml"), nil
	case strings.Contains(body, "Action=DescribeSubnets"):
		return getContent("describe_subnet_response_empty.xml"), nil
	case strings.Contains(body, "Action=DescribeVpcs"):
		return getContent("describe_vpcs_response.xml"), nil
	case strings.Contains(body, "Action=ListStacks"):
		return getContent("list_stacks_response.xml"), nil
	default:
		fmt.Println("Body: ", body)
		return "", errors.New("unknown request in body")
	}
}
