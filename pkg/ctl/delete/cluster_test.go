package delete

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

const (
	clusterName = "clusterName"
)

var _ = Describe("delete cluster", func() {
	DescribeTable("should be called to delete the cluster",
		func(forceExpected bool, disableNodegroupEvictionExpected bool, args ...string) {
			cmd := newMockEmptyCmd(args...)
			count := 0
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				deleteClusterWithRunFunc(cmd, func(cmd *cmdutils.Cmd, force bool, disableNodegroupEviction bool, podEvictionWaitPeriod time.Duration, parallel int) error {
					Expect(cmd.ClusterConfig.Metadata.Name).To(Equal(clusterName))
					Expect(force).To(Equal(forceExpected))
					Expect(disableNodegroupEviction).To(Equal(disableNodegroupEvictionExpected))
					count++
					return nil
				})
			})
			_, err := cmd.execute()
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		},
		Entry("with only valid cluster name", false, false, "cluster", "--name", clusterName),
		Entry("with valid cluster name and force flag", true, false, "cluster", "--name", clusterName, "--force"),
		Entry("with valid cluster name and disableNodeGroupEviction flag", false, true, "cluster", "--name", clusterName, "--disable-nodegroup-eviction"),
		Entry("with valid cluster name, force & disableNodeGroupEviction flags", true, true, "cluster", "--name", clusterName, "--force", "--disable-nodegroup-eviction"),
	)
})
