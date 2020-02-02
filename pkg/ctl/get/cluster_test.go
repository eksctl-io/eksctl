package get

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

var _ = Describe("get", func() {
	Describe("cluster", func() {
		It("without cluster name", func() {
			count := 0
			cmd := newMockEmptyCmd("cluster")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				getClusterWithRunFunc(cmd, func(cmd *cmdutils.Cmd, params *getCmdParams, listAllRegions bool) error {
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
				getClusterWithRunFunc(cmd, func(cmd *cmdutils.Cmd, params *getCmdParams, listAllRegions bool) error {
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
				getClusterWithRunFunc(cmd, func(cmd *cmdutils.Cmd, params *getCmdParams, listAllRegions bool) error {
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

		It("with all-regions flags", func() {
			count := 0
			cmd := newMockEmptyCmd("cluster", "--all-regions")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				getClusterWithRunFunc(cmd, func(cmd *cmdutils.Cmd, params *getCmdParams, listAllRegions bool) error {
					count++
					return nil
				})
			})
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
		})

		It("with both cluster name argument and --all-regions flag", func() {
			cmd := newMockDefaultCmd("cluster", "clusterName", "--all-regions")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--all-regions is for listing all clusters, it must be used without cluster name flag/argument"))
		})

		It("with both cluster name and --all-regions flags", func() {
			cmd := newMockDefaultCmd("cluster", "--name", "clusterName", "--all-regions")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--all-regions is for listing all clusters, it must be used without cluster name flag/argument"))
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
				getClusterWithRunFunc(cmd, func(cmd *cmdutils.Cmd, params *getCmdParams, listAllRegions bool) error {
					return fmt.Errorf("unable to fetch the cluster")
				})
			})
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unable to fetch the cluster"))
		})
	})
})
