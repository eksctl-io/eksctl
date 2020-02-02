package set

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

var _ = Describe("set", func() {
	Describe("labels", func() {
		It("with valid details", func() {
			count := 0
			cmd := newMockEmptyCmd("labels", "--cluster", "clusterName", "--labels", "testLabel=testValue")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				setLabelsWithRunFunc(cmd, func(cmd *cmdutils.Cmd, options labelOptions) error {
					Expect(cmd.ClusterConfig.Metadata.Name).To(Equal("clusterName"))
					Expect(options.labels).To(Equal(map[string]string{
						"testLabel": "testValue",
					}))
					count++
					return nil
				})
			})
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(count).To(Equal(1))
		})

		It("with invalid label format", func() {
			cmd := newMockDefaultCmd("labels", "--cluster", "clusterName", "--labels", "testLabel")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid argument \"testLabel\" for \"-l, --labels\" flag: testLabel must be formatted as key=value"))
		})

		It("with cluster name and node group flags", func() {
			count := 0
			cmd := newMockEmptyCmd("labels", "--cluster", "clusterName", "--nodegroup", "nodeGroup", "--labels", "testLabel=testValue")
			cmdutils.AddResourceCmd(cmdutils.NewGrouping(), cmd.parentCmd, func(cmd *cmdutils.Cmd) {
				setLabelsWithRunFunc(cmd, func(cmd *cmdutils.Cmd, options labelOptions) error {
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

		It("with one name argument", func() {
			cmd := newMockDefaultCmd("labels", "clusterName", "--cluster", "dummyName", "--labels", "testLabel=testValue")
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
			cmd := newMockDefaultCmd("labels", "--cluster", "dummy", "dummyName", "--labels", "testLabel=testValue")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("name argument is not supported"))
		})
	})
})
