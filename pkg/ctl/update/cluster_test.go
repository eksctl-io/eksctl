package update

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

var _ = Describe("update", func() {
	Describe("cluster", func() {
		It("without cluster name", func() {
			count := 0
			cmd := newMockEmptyCmd("cluster")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				updateClusterWithRunFunc(cmd, func(cmd *cmdutils.Cmd) error {
					count++
					return nil
				})
			})
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
		})

		It("with cluster name as flag", func() {
			count := 0
			cmd := newMockEmptyCmd("cluster", "--name", "clusterName")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				updateClusterWithRunFunc(cmd, func(cmd *cmdutils.Cmd) error {
					count++
					return nil
				})
			})
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
		})

		It("with cluster name as argument", func() {
			count := 0
			cmd := newMockEmptyCmd("cluster", "clusterName")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				updateClusterWithRunFunc(cmd, func(cmd *cmdutils.Cmd) error {
					count++
					return nil
				})
			})
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
		})

		It("with cluster name as argument and flag", func() {
			cmd := newMockDefaultCmd("cluster", "clusterName", "--name", "clusterName")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--name=clusterName and argument clusterName cannot be used at the same time"))
		})

		It("with config file flag", func() {
			count := 0
			cmd := newMockEmptyCmd("cluster", "--config-file", "../../../examples/01-simple-cluster.yaml")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				updateClusterWithRunFunc(cmd, func(cmd *cmdutils.Cmd) error {
					count++
					return nil
				})
			})
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
		})

		It("with cluster name and deprecated dry-run flag", func() {
			count := 0
			cmd := newMockEmptyCmd("cluster", "clusterName", "--dry-run")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				updateClusterWithRunFunc(cmd, func(cmd *cmdutils.Cmd) error {
					count++
					return nil
				})
			})
			out, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
			Expect(out).To(ContainSubstring("Flag --dry-run has been deprecated, see --approve"))
		})

		It("with cluster name and deprecated wait flag", func() {
			count := 0
			cmd := newMockEmptyCmd("cluster", "clusterName", "--wait")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				updateClusterWithRunFunc(cmd, func(cmd *cmdutils.Cmd) error {
					count++
					return nil
				})
			})
			out, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
			Expect(out).To(ContainSubstring("Flag --wait has been deprecated, --wait is no longer respected; the cluster update always waits to complete"))
		})

		It("with invalid flags", func() {
			cmd := newMockDefaultCmd("cluster", "--invalid", "dummy")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unknown flag: --invalid"))
		})

		It("failed due to any reason", func() {
			cmd := newMockEmptyCmd("cluster")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				updateClusterWithRunFunc(cmd, func(cmd *cmdutils.Cmd) error {
					return fmt.Errorf("unable to fetch the cluster")
				})
			})
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unable to fetch the cluster"))
		})
	})
})

