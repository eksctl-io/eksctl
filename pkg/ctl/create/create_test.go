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

func newDefaultCmd(args ...string) *mockCmd {
	cmd := Command(cmdutils.NewGrouping())
	cmd.SetArgs(args)
	return &mockCmd{
		parentCmd: cmd,
	}
}

func newMockEmptyCmd(args ...string) *mockCmd {
	cmd := cmdutils.NewVerbCmd("create", "Create resource(s)", "")
	cmd.SetArgs(args)
	return &mockCmd{
		parentCmd: cmd,
	}
}

func newMockCmdWithRunFunc(verb string, runFunc func(cmd *cmdutils.Cmd), args ...string) *mockCmd {
	grouping := cmdutils.NewGrouping()
	parentCmd := cmdutils.NewVerbCmd(verb, "", "")

	var mc mockCmd
	cmdutils.AddResourceCmd(grouping, parentCmd, func(cmd *cmdutils.Cmd) {
		mc.cmd = cmd
		runFunc(cmd)
	})
	parentCmd.SetArgs(args)
	mc.parentCmd = parentCmd
	return &mc
}

type mockCmd struct {
	parentCmd *cobra.Command
	cmd       *cmdutils.Cmd
}

func (c *mockCmd) execute() (string, error) {
	var (
		stdOut bytes.Buffer
		stdErr bytes.Buffer
	)
	c.parentCmd.SetOut(&stdOut)
	c.parentCmd.SetErr(&stdErr)
	err := c.parentCmd.Execute()
	if err != nil {
		err = errors.New(stdErr.String())
	}
	return stdOut.String(), err
}
