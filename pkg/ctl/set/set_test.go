package set

import (
	"bytes"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

var _ = Describe("set", func() {
	Describe("invalid-resource", func() {
		It("fails", func() {
			cmd := newMockCmd("invalid-resource")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Error: unknown command \"invalid-resource\" for \"set\""))
			Expect(err.Error()).To(ContainSubstring("usage"))
		})
	})

	Describe("labels", func() {
		It("fails when no flags set", func() {
			cmd := newMockCmd("labels")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Error: --labels must be set"))
		})

		It("fails when cluster flag not set", func() {
			cmd := newMockCmd("labels", "-l", "k=v")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Error: --cluster must be set"))
		})

		It("fails when --nodegroup flag not set", func() {
			cmd := newMockCmd("labels", "--cluster", "dummy", "-l", "k=v")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Error: --nodegroup must be set"))
		})

		It("fails when name argument is used", func() {
			cmd := newMockCmd("labels", "--cluster", "dummy", "--nodegroup", "dummyNodeGroup", "dummyName", "-l", "k=v")
			_, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Error: name argument is not supported"))
		})
	})
})

func newMockCmd(args ...string) *mockVerbCmd {
	flagGrouping := cmdutils.NewGrouping()
	cmd := Command(flagGrouping)
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
