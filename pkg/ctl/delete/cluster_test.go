package delete

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

var _ = Describe("delete", func() {
	Describe("cluster", func() {
		It("without cluster name", func() {
			cmd := newMockDefaultCmd("cluster")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--name must be set"))
		})

		It("with cluster name as flag", func() {
			count := 0
			cmd := newMockEmptyCmd("cluster", "--name", "clusterName")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				deleteClusterWithRunFunc(cmd, func(cmd *cmdutils.Cmd) error {
					Expect(cmd.ClusterConfig.Metadata.Name).To(Equal("clusterName"))
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
				deleteClusterWithRunFunc(cmd, func(cmd *cmdutils.Cmd) error {
					Expect(cmd.NameArg).To(Equal("clusterName"))
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

		It("with invalid flags", func() {
			cmd := newMockDefaultCmd("cluster", "--invalid", "dummy")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unknown flag: --invalid"))
		})

		It("failed due to any reason", func() {
			cmd := newMockEmptyCmd("cluster", "clusterName")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				deleteClusterWithRunFunc(cmd, func(cmd *cmdutils.Cmd) error {
					return fmt.Errorf("unable to fetch the cluster")
				})
			})
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unable to fetch the cluster"))
		})
	})
})
