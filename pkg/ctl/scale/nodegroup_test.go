package scale

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

var _ = Describe("scale", func() {
	Describe("scale nodegroup", func() {
		It("requires the cluster's name, and if missing, prints an error and the usage", func() {
			cmd := newMockScaleNodegroupCmd("nodegroup")
			out, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("One of [--cluster --use-kubeconfig-context] must be set"))
			Expect(out).To(ContainSubstring("Error: One of [--cluster --use-kubeconfig-context] must be set"))
			Expect(out).To(ContainSubstring("Usage:"))
		})
	})
})

func newMockScaleNodegroupCmd(args ...string) *mockScaleNodegroupCmd {
	mockCmd := &mockScaleNodegroupCmd{}
	grouping := cmdutils.NewGrouping()
	parentCmd := cmdutils.NewVerbCmd("scale", "", "")
	cmdutils.AddResourceCmd(grouping, parentCmd, func(cmd *cmdutils.Cmd) {
		scaleNodeGroupCmd(cmd)
	})
	parentCmd.SetArgs(args)
	mockCmd.parentCmd = parentCmd
	return mockCmd
}

type mockScaleNodegroupCmd struct {
	parentCmd *cobra.Command
	cmd       *cmdutils.Cmd
}

func (c mockScaleNodegroupCmd) execute() (string, error) {
	buf := new(bytes.Buffer)
	c.parentCmd.SetOutput(buf)
	err := c.parentCmd.Execute()
	return buf.String(), err
}
