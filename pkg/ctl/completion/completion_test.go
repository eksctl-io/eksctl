package completion

import (
	"bytes"
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

var _ = Describe("completion", func() {
	It("with bash", func() {
		cmd := newMockCmd("bash")
		out, err := cmd.execute()
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(ContainSubstring("bash completion for eksctl"))
	})

	It("with zsh", func() {
		cmd := newMockCmd("zsh")
		out, err := cmd.execute()
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(ContainSubstring("#compdef _eksctl eksctl"))
	})

	It("with fish", func() {
		cmd := newMockCmd("fish")
		out, _ := cmd.execute()
		Expect(out).To(ContainSubstring("fish completion for eksctl"))
	})

	It("with powershell", func() {
		cmd := newMockCmd("powershell")
		out, _ := cmd.execute()
		Expect(out).To(ContainSubstring("Register-ArgumentCompleter -CommandName 'eksctl'"))
	})

	It("with invalid shell", func() {
		cmd := newMockCmd("invalid-shell")
		_, err := cmd.execute()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Error: unknown command \"invalid-shell\" for \"completion\""))
		Expect(err.Error()).To(ContainSubstring("usage"))
	})

	It("with no shell", func() {
		cmd := newMockCmd("")
		out, _ := cmd.execute()
		Expect(out).To(ContainSubstring("Usage"))
	})
})

func newMockCmd(args ...string) *mockVerbCmd {
	cmd := Command(&cobra.Command{
		Use:   "eksctl [command]",
		Short: "The official CLI for Amazon EKS",
		Run: func(c *cobra.Command, _ []string) {
			if err := c.Help(); err != nil {
				fmt.Println(err.Error())
			}
		},
		SilenceUsage: true,
	})
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
