package unset

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

var _ = Describe("unset", func() {
	Describe("labels", func() {
		It("with valid details", func() {
			count := 0
			cmd := newMockEmptyCmd("labels", "--cluster", "clusterName", "--labels", "testLabel")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				unsetLabelsWithRunFunc(cmd, func(cmd *cmdutils.Cmd, nodeGroupName string, removeLabels []string) error {
					Expect(cmd.ClusterConfig.Metadata.Name).To(Equal("clusterName"))
					Expect(removeLabels).To(ConsistOf("testLabel"))
					count++
					return nil
				})
			})
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
		})

		It("with multiple labels", func() {
			count := 0
			cmd := newMockEmptyCmd("labels", "--cluster", "clusterName", "--labels", "testLabel,testAnotherLabel")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				unsetLabelsWithRunFunc(cmd, func(cmd *cmdutils.Cmd, nodeGroupName string, removeLabels []string) error {
					Expect(cmd.ClusterConfig.Metadata.Name).To(Equal("clusterName"))
					Expect(removeLabels).To(ConsistOf("testLabel", "testAnotherLabel"))
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
			cmd := newMockEmptyCmd("labels", "--cluster", "clusterName", "--nodegroup", "nodeGroup", "--labels", "testLabel")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				unsetLabelsWithRunFunc(cmd, func(cmd *cmdutils.Cmd, nodeGroupName string, removeLabels []string) error {
					Expect(cmd.ClusterConfig.Metadata.Name).To(Equal("clusterName"))
					Expect(nodeGroupName).To(Equal("nodeGroup"))
					Expect(removeLabels).To(ConsistOf("testLabel"))
					count++
					return nil
				})
			})
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
		})

		It("with one name argument", func() {
			cmd := newMockDefaultCmd("labels", "clusterName", "--cluster", "dummyName", "--labels", "testLabel")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("name argument is not supported"))
		})

		It("missing required flag --labels", func() {
			cmd := newMockDefaultCmd("labels")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("required flag(s) \"labels\" not set"))
		})

		It("missing required flag --labels", func() {
			cmd := newMockDefaultCmd("labels", "--cluster", "dummyName")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("required flag(s) \"labels\" not set"))
		})

		It("setting name argument", func() {
			cmd := newMockDefaultCmd("labels", "--cluster", "dummy", "dummyName", "--labels", "testLabel")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("name argument is not supported"))
		})
	})
})
