package get

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

var _ = Describe("get", func() {
	Describe("labels", func() {
		It("with cluster name as flag", func() {
			count := 0
			cmd := newMockEmptyCmd("labels", "--cluster", "clusterName")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				getLabelsWithRunFunc(cmd, func(cmd *cmdutils.Cmd, nodeGroupName string) error {
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
			cmd := newMockEmptyCmd("labels", "--cluster", "clusterName", "--nodegroup", "nodeGroup")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				getLabelsWithRunFunc(cmd, func(cmd *cmdutils.Cmd, nodeGroupName string) error {
					Expect(cmd.ClusterConfig.Metadata.Name).To(Equal("clusterName"))
					Expect(nodeGroupName).To(Equal("nodeGroup"))
					count++
					return nil
				})
			})
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
		})

		It("with one name argument", func() {
			cmd := newMockDefaultCmd("labels", "clusterName", "--cluster", "dummyName")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("name argument is not supported"))
		})

		It("missing required flag --cluster", func() {
			cmd := newMockDefaultCmd( "labels")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--cluster must be set"))
		})

		It("missing required flag --cluster, but with --nodegroup", func() {
			cmd := newMockDefaultCmd( "labels", "--nodegroup", "dummyNodeGroup")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--cluster must be set"))
		})

		It("setting name argument", func() {
			cmd := newMockDefaultCmd("labels", "--cluster", "dummy", "dummyName")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("name argument is not supported"))
		})
	})
})
