package create

import (
	"fmt"

	. "github.com/onsi/ginkgo/extensions/table"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

type invalidParamsCase struct {
	args  []string
	error error
}

var _ = Describe("create cluster", func() {
	Describe("un-managed node group", func() {
		DescribeTable("create cluster successfully",
			func(args ...string) {
				commandArgs := append([]string{"cluster"}, args...)
				cmd := newMockEmptyCmd(commandArgs...)
				count := 0
				cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
					createClusterCmdWithRunFunc(cmd, func(cmd *cmdutils.Cmd, ng *api.NodeGroup, params *cmdutils.CreateClusterCmdParams) error {
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
			Entry("with node-type flag", "--node-type", "m5.large"),
			Entry("with nodes flag", "--nodes", "2"),
			Entry("with nodes-min flag", "--nodes-min", "2"),
			Entry("with nodes-max flag", "--nodes-max", "2"),
			Entry("with node-volume-size flag", "--node-volume-size", "2"),
			Entry("with node-volume-type flag", "--node-volume-type", "gp2"),
			Entry("with max-pods-per-node flag", "--max-pods-per-node", "20"),
			Entry("with ssh-access flag", "--ssh-access", "true"),
			Entry("with ssh-public-key flag", "--ssh-public-key", "dummy-public-key"),
			Entry("with node-ami flag", "--node-ami", "ami-dummy-123"),
			Entry("with node-ami-family flag", "--node-ami-family", "AmazonLinux2"),
			Entry("with node-private-networking flag", "--node-private-networking", "true"),
			Entry("with node-security-groups flag", "--node-security-groups", "sg-123"),
			Entry("with node-labels flag", "--node-labels", "partition=backend,nodeclass=hugememory"),
			Entry("with node-zones flag", "--node-zones", "zone1,zone2,zone3"),
			Entry("with asg-access flag", "--asg-access", "true"),
			Entry("with external-dns-access flag", "--external-dns-access", "true"),
			Entry("with full-ecr-access flag", "--full-ecr-access", "true"),
			Entry("with appmesh-access flag", "--appmesh-access", "true"),
			Entry("with alb-ingress-access flag", "--alb-ingress-access", "true"),
		)

		DescribeTable("invalid flags or arguments",
			func(c invalidParamsCase) {
				commandArgs := append([]string{"cluster"}, c.args...)
				cmd := newMockDefaultCmd(commandArgs...)
				_, err := cmd.execute()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(c.error.Error()))
			},
			Entry("with cluster name as argument and flag", invalidParamsCase{
				args:  []string{"clusterName", "--name", "clusterName"},
				error: fmt.Errorf("--name=clusterName and argument clusterName cannot be used at the same time"),
			}),
			Entry("with invalid flags", invalidParamsCase{
				args:  []string{"cluster", "--invalid", "dummy"},
				error: fmt.Errorf("unknown flag: --invalid"),
			}),
		)
	})

	Describe("managed node group", func() {
		DescribeTable("create cluster successfully",
			func(args ...string) {
				commandArgs := append([]string{"cluster", "--managed"}, args...)
				cmd := newMockEmptyCmd(commandArgs...)
				count := 0
				cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
					createClusterCmdWithRunFunc(cmd, func(cmd *cmdutils.Cmd, ng *api.NodeGroup, params *cmdutils.CreateClusterCmdParams) error {
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
			Entry("with node-type flag", "--node-type", "m5.large"),
			Entry("with nodes flag", "--nodes", "2"),
			Entry("with nodes-min flag", "--nodes-min", "2"),
			Entry("with nodes-max flag", "--nodes-max", "2"),
			Entry("with node-volume-size flag", "--node-volume-size", "2"),
			Entry("with ssh-access flag", "--ssh-access", "true"),
			Entry("with ssh-public-key flag", "--ssh-public-key", "dummy-public-key"),
			Entry("with node-ami-family flag", "--node-ami-family", "AmazonLinux2"),
			Entry("with node-private-networking flag", "--node-private-networking", "true"),
			Entry("with node-labels flag", "--node-labels", "partition=backend,nodeclass=hugememory"),
			Entry("with node-zones flag", "--node-zones", "zone1,zone2,zone3"),
			Entry("with asg-access flag", "--asg-access", "true"),
			Entry("with external-dns-access flag", "--external-dns-access", "true"),
			Entry("with full-ecr-access flag", "--full-ecr-access", "true"),
			Entry("with appmesh-access flag", "--appmesh-access", "true"),
			Entry("with alb-ingress-access flag", "--alb-ingress-access", "true"),
		)

		DescribeTable("with un-supported flags",
			func(args ...string) {
				commandArgs := append([]string{"cluster", "--managed"}, args...)
				cmd := newMockDefaultCmd(commandArgs...)
				_, err := cmd.execute()
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(fmt.Errorf("%s is not supported for Managed Nodegroups (--managed=true)", args[0])))
			},
			Entry("node-volume-type", "--node-volume-type", "gp2"),
			Entry("max-pods-per-node", "--max-pods-per-node", "2"),
			Entry("node-ami", "--node-ami", "ami-dummy-123"),
			Entry("node-security-groups", "--node-security-groups", "sg-123"),
		)

		DescribeTable("invalid flags or arguments",
			func(c invalidParamsCase) {
				commandArgs := append([]string{"cluster", "--managed"}, c.args...)
				cmd := newMockDefaultCmd(commandArgs...)
				_, err := cmd.execute()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(c.error.Error()))
			},
			Entry("with cluster name as argument and flag", invalidParamsCase{
				args:  []string{"clusterName", "--name", "clusterName"},
				error: fmt.Errorf("--name=clusterName and argument clusterName cannot be used at the same time"),
			}),
			Entry("with invalid flags", invalidParamsCase{
				args:  []string{"cluster", "--invalid", "dummy"},
				error: fmt.Errorf("unknown flag: --invalid"),
			}),
		)
	})
})
