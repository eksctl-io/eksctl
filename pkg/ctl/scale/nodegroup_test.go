package scale

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

var _ = Describe("scale", func() {
	Describe("nodegroup", func() {
		It("with valid details", func() {
			count := 0
			cmd := newMockEmptyCmd("nodegroup", "--cluster", "clusterName", "--name", "nodeGroup", "--nodes", "2")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				scaleNodeGroupWithRunFunc(cmd, func(cmd *cmdutils.Cmd, ng *v1alpha5.NodeGroup) error {
					Expect(cmd.ClusterConfig.Metadata.Name).To(Equal("clusterName"))
					Expect(ng.Name).To(Equal("nodeGroup"))
					Expect(*ng.DesiredCapacity).To(Equal(2))
					count++
					return nil
				})
			})
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
		})

		It("missing required flag --cluster", func() {
			cmd := newMockDefaultCmd("nodegroup")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--cluster must be set"))
		})

		It("setting --name and argument at the same time", func() {
			cmd := newMockDefaultCmd("nodegroup", "ng", "--cluster", "dummy", "--name", "ng")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--name=ng and argument ng cannot be used at the same time"))
		})

		It("missing required nodes flag --nodes", func() {
			cmd := newMockDefaultCmd("nodegroup", "ng", "--cluster", "dummy")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("number of nodes must be 0 or greater. Use the --nodes/-N flag"))
		})

		It("invalid flag", func() {
			cmd := newMockDefaultCmd("nodegroup", "--invalid", "dummy")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unknown flag: --invalid"))
		})
	})
})
