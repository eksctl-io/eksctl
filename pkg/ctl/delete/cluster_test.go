package delete

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

const (
	clusterName = "clusterName"
)

var _ = Describe("delete cluster", func() {
	DescribeTable("should be called with no extra flags",
		func(args ...string) {
			cmd := newMockEmptyCmd(args...)
			count := 0
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				deleteClusterWithRunFunc(cmd, func(cmd *cmdutils.Cmd, force bool, disableNodegroupEviction bool) error {
					Expect(cmd.ClusterConfig.Metadata.Name).To(Equal(clusterName))
					Expect(force).To(Equal(false))
					Expect(disableNodegroupEviction).To(Equal(false))
					count++
					return nil
				})
			})
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
		},
		Entry("with only valid cluster name", "cluster", "--name", clusterName),
	)

	DescribeTable("should be called with force flag",
		func(args ...string) {
			cmd := newMockEmptyCmd(args...)
			count := 0
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				deleteClusterWithRunFunc(cmd, func(cmd *cmdutils.Cmd, force bool, disableNodegroupEviction bool) error {
					Expect(cmd.ClusterConfig.Metadata.Name).To(Equal(clusterName))
					Expect(force).To(Equal(true))
					Expect(disableNodegroupEviction).To(Equal(false))
					count++
					return nil
				})
			})
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
		},
		Entry("with only valid cluster name", "cluster", "--name", clusterName, "--force"),
	)

	DescribeTable("should be called with disable nodegroup eviction flag",
		func(args ...string) {
			cmd := newMockEmptyCmd(args...)
			count := 0
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				deleteClusterWithRunFunc(cmd, func(cmd *cmdutils.Cmd, force bool, disableNodegroupEviction bool) error {
					Expect(cmd.ClusterConfig.Metadata.Name).To(Equal(clusterName))
					Expect(force).To(Equal(false))
					Expect(disableNodegroupEviction).To(Equal(true))
					count++
					return nil
				})
			})
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
		},
		Entry("with only valid cluster name", "cluster", "--name", clusterName, "--disable-nodegroup-eviction"),
	)

	DescribeTable("should be called with both force and disable nodegroup eviction flags",
		func(args ...string) {
			cmd := newMockEmptyCmd(args...)
			count := 0
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				deleteClusterWithRunFunc(cmd, func(cmd *cmdutils.Cmd, force bool, disableNodegroupEviction bool) error {
					Expect(cmd.ClusterConfig.Metadata.Name).To(Equal(clusterName))
					Expect(force).To(Equal(true))
					Expect(disableNodegroupEviction).To(Equal(true))
					count++
					return nil
				})
			})
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
		},
		Entry("with only valid cluster name", "cluster", "--name", clusterName, "--force", "--disable-nodegroup-eviction"),
	)
})
