package completion

import (
	"bytes"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

var _ = Describe("completion", func() {
	It("with bash", func() {
		cmd := newMockCmd("bash")
		out, err := cmd.execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(out).To(ContainSubstring("bash completion for eksctl"))
	})

	It("with zsh", func() {
		cmd := newMockCmd("zsh")
		out, err := cmd.execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(out).To(ContainSubstring("#compdef _eksctl eksctl"))
	})

	It("with fish", func() {
		cmd := newMockCmd("fish")
		out, _ := cmd.execute()
		Expect(out).To(ContainSubstring("fish completion for eksctl"))
	})

	It("with invalid shell", func() {
		cmd := newMockCmd("invalid-shell")
		out, err := cmd.execute()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("unknown command \"invalid-shell\" for \"completion\""))
		Expect(out).To(ContainSubstring("usage"))
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
	cmd       *cmdutils.Cmd
}

func (c mockVerbCmd) execute() (string, error) {
	buf := new(bytes.Buffer)
	c.parentCmd.SetOut(buf)
	err := c.parentCmd.Execute()
	return buf.String(), err
}
