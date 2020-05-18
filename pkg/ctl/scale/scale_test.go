package scale

import (
	"bytes"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

type invalidParamsCase struct {
	args   []string
	error  error
	output string
}

var _ = Describe("generate", func() {
	Describe("invalid-resource", func() {
		DescribeTable("invalid flags or arguments",
			func(c invalidParamsCase) {
				cmd := newDefaultCmd(c.args...)
				out, err := cmd.execute()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(c.error.Error()))
				Expect(out).To(ContainSubstring(c.output))
			},
			Entry("missing required flag --cluster", invalidParamsCase{
				args:   []string{"invalid-resource"},
				error:  fmt.Errorf("unknown command \"invalid-resource\" for \"scale\""),
				output: "usage",
			}),
			Entry("with invalid-resource and some flag", invalidParamsCase{
				args:   []string{"invalid-resource", "--invalid-flag", "foo"},
				error:  fmt.Errorf("unknown command \"invalid-resource\" for \"scale\""),
				output: "usage",
			}),
			Entry("with invalid-resource and additional argument", invalidParamsCase{
				args:   []string{"invalid-resource", "foo"},
				error:  fmt.Errorf("unknown command \"invalid-resource\" for \"scale\""),
				output: "usage",
			}),
		)
	})
})

func newDefaultCmd(args ...string) *mockVerbCmd {
	flagGrouping := cmdutils.NewGrouping()
	cmd := Command(flagGrouping)
	cmd.SetArgs(args)
	return &mockVerbCmd{
		parentCmd: cmd,
	}
}

func newMockEmptyCmd(args ...string) *mockVerbCmd {
	cmd := cmdutils.NewVerbCmd("scale", "Scale resources(s)", "")
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
