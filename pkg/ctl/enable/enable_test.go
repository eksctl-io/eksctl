package enable

import (
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

var _ = Describe("enable", func() {
	Describe("invalid-resource", func() {
		It("with no flag", func() {
			cmd := newMockDefaultCmd("invalid-resource")
			out, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unknown command \"invalid-resource\" for \"enable\""))
			Expect(out).To(ContainSubstring("usage"))
		})
		It("with invalid-resource and some flag", func() {
			cmd := newMockDefaultCmd("invalid-resource", "--invalid-flag", "foo")
			out, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unknown command \"invalid-resource\" for \"enable\""))
			Expect(out).To(ContainSubstring("usage"))
		})
		It("with invalid-resource and additional argument", func() {
			cmd := newMockDefaultCmd("invalid-resource", "foo")
			out, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unknown command \"invalid-resource\" for \"enable\""))
			Expect(out).To(ContainSubstring("usage"))
		})
	})
})

func newMockDefaultCmd(args ...string) *mockVerbCmd {
	cmd := Command(cmdutils.NewGrouping())
	cmd.SetArgs(args)
	return &mockVerbCmd{
		parentCmd: cmd,
	}
}

func newMockEmptyCmd(args ...string) *mockVerbCmd {
	cmd := cmdutils.NewVerbCmd("enable", "Get resource(s)", "")
	cmd.SetArgs(args)
	return &mockVerbCmd{
		parentCmd: cmd,
	}
}

type mockVerbCmd struct {
	parentCmd *cobra.Command
	cmd       *cmdutils.Cmd
}

func (c mockVerbCmd) execute() (string, error) {
	buf := new(bytes.Buffer)
	c.parentCmd.SetOut(buf)
	err := c.parentCmd.Execute()
	return buf.String(), err
}
