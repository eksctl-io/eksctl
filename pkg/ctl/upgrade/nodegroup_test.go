package upgrade

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

var _ = Describe("upgrade", func() {
	Describe("nodegroup", func() {
		It("with cluster name as flag", func() {
			count := 0
			cmd := newMockEmptyUpgradeCmd("nodegroup", "--cluster", "clusterName")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				upgradeNodeGroupWithRunFunc(cmd, func(cmd *cmdutils.Cmd, options upgradeOptions) error {
					Expect(cmd.ClusterConfig.Metadata.Name).To(Equal("clusterName"))
					count++
					return nil
				})
			})
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
		})

		It("with cluster name and node group flags", func() {
			count := 0
			cmd := newMockEmptyUpgradeCmd("nodegroup", "--cluster", "clusterName", "--name", "nodeGroup")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				upgradeNodeGroupWithRunFunc(cmd, func(cmd *cmdutils.Cmd, options upgradeOptions) error {
					Expect(cmd.ClusterConfig.Metadata.Name).To(Equal("clusterName"))
					Expect(options.nodeGroupName).To(Equal("nodeGroup"))
					count++
					return nil
				})
			})
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
		})

		It("with cluster name, node group and kube version flags", func() {
			count := 0
			cmd := newMockEmptyUpgradeCmd("nodegroup", "--cluster", "clusterName", "--name", "nodeGroup", "--kubernetes-version", "1.14")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				upgradeNodeGroupWithRunFunc(cmd, func(cmd *cmdutils.Cmd, options upgradeOptions) error {
					Expect(cmd.ClusterConfig.Metadata.Name).To(Equal("clusterName"))
					Expect(options.nodeGroupName).To(Equal("nodeGroup"))
					Expect(options.kubernetesVersion).To(Equal("1.14"))
					count++
					return nil
				})
			})
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
		})

		It("missing required flag --cluster", func() {
			cmd := newMockDefaultUpgradeCmd("nodegroup")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--cluster must be set"))
		})

		It("missing required node group name", func() {
			cmd := newMockDefaultUpgradeCmd("nodegroup", "--cluster", "clusterName")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("name must be set"))
		})

		It("setting --name and argument at the same time", func() {
			cmd := newMockDefaultUpgradeCmd("nodegroup", "ng", "--cluster", "dummy", "--name", "ng")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--name=ng and argument ng cannot be used at the same time"))
		})

		It("invalid flag", func() {
			cmd := newMockDefaultUpgradeCmd("nodegroup", "--invalid", "dummy")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unknown flag: --invalid"))
		})
	})
})
