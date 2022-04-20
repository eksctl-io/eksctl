package create

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

var _ = Describe("create", func() {
	Describe("invalid-resource", func() {
		It("with no flag", func() {
			cmd := newDefaultCmd("invalid-resource")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Error: unknown command \"invalid-resource\" for \"create\""))
			Expect(err.Error()).To(ContainSubstring("usage"))
		})
		It("with invalid-resource and some flag", func() {
			cmd := newDefaultCmd("invalid-resource", "--invalid-flag", "foo")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Error: unknown command \"invalid-resource\" for \"create\""))
			Expect(err.Error()).To(ContainSubstring("usage"))
		})
		It("with invalid-resource and additional argument", func() {
			cmd := newDefaultCmd("invalid-resource", "foo")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Error: unknown command \"invalid-resource\" for \"create\""))
			Expect(err.Error()).To(ContainSubstring("usage"))
		})
	})
})

type invalidParamsCase struct {
	args  []string
	error string
}

func newDefaultCmd(args ...string) *mockVerbCmd {
	cmd := Command(cmdutils.NewGrouping())
	cmd.SetArgs(args)
	return &mockVerbCmd{
		parentCmd: cmd,
	}
}

func newMockEmptyCmd(args ...string) *mockVerbCmd {
	cmd := cmdutils.NewVerbCmd("create", "Create resource(s)", "")
	cmd.SetArgs(args)
	return &mockVerbCmd{
		parentCmd: cmd,
	}
}

type mockVerbCmd struct {
	parentCmd *cobra.Command
}

func (c mockVerbCmd) execute() (string, error) {
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	c.parentCmd.SetOut(outBuf)
	c.parentCmd.SetErr(errBuf)
	err := c.parentCmd.Execute()
	if err != nil {
		err = errors.New(errBuf.String())
	}
	return outBuf.String(), err
}
