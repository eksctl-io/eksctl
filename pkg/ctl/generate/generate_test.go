package generate

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

var _ = Describe("generate", func() {
	Describe("invalid-resource", func() {
		It("with no flag", func() {
			cmd := newMockCmd("invalid-resource")
			out, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unknown command \"invalid-resource\" for \"generate\""))
			Expect(out).To(ContainSubstring("usage"))
		})
		It("with invalid-resource and some flag", func() {
			cmd := newMockCmd("invalid-resource", "--invalid-flag", "foo")
			out, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unknown command \"invalid-resource\" for \"generate\""))
			Expect(out).To(ContainSubstring("usage"))
		})
		It("with invalid-resource and additional argument", func() {
			cmd := newMockCmd("invalid-resource", "foo")
			out, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unknown command \"invalid-resource\" for \"generate\""))
			Expect(out).To(ContainSubstring("usage"))
		})
	})
})

func newMockCmd(args ...string) *mockVerbCmd {
	cmd := Command(cmdutils.NewGrouping())
	cmd.SetArgs(args)
	return &mockVerbCmd{
		parentCmd: cmd,
	}
}

type mockVerbCmd struct {
	parentCmd *cobra.Command
}

func (c mockVerbCmd) execute() (string, error) {
	buf := new(bytes.Buffer)
	c.parentCmd.SetOut(buf)
	err := c.parentCmd.Execute()
	return buf.String(), err
}
