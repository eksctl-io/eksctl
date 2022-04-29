package scale

import (
	"bytes"
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo/v2"

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
				_, err := cmd.execute()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(c.error.Error()))
				Expect(err.Error()).To(ContainSubstring(c.output))
			},
			Entry("missing required flag --cluster", invalidParamsCase{
				args:   []string{"invalid-resource"},
				error:  fmt.Errorf("Error: unknown command \"invalid-resource\" for \"scale\""),
				output: "usage",
			}),
			Entry("with invalid-resource and some flag", invalidParamsCase{
				args:   []string{"invalid-resource", "--invalid-flag", "foo"},
				error:  fmt.Errorf("Error: unknown command \"invalid-resource\" for \"scale\""),
				output: "usage",
			}),
			Entry("with invalid-resource and additional argument", invalidParamsCase{
				args:   []string{"invalid-resource", "foo"},
				error:  fmt.Errorf("Error: unknown command \"invalid-resource\" for \"scale\""),
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
