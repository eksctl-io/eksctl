package delete

import (
	"bytes"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/fargate"
)

var _ = Describe("delete", func() {
	Describe("delete fargateprofile", func() {
		It("requires the cluster's name, and if missing, prints an error and the usage", func() {
			cmd := newMockDeleteFargateProfileCmd("fargateprofile")
			out, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Error: --cluster must be set"))
			Expect(out).To(ContainSubstring("Usage:"))
		})

		It("requires a profile name, and if missing, prints an error and the usage", func() {
			cmd := newMockDeleteFargateProfileCmd("fargateprofile", "--cluster", "foo")
			out, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Error: invalid Fargate profile: empty name"))
			Expect(out).To(ContainSubstring("Usage:"))
		})

		It("requires a profile name, which can be provided as an argument", func() {
			cmd := newMockDeleteFargateProfileCmd("fargateprofile", "--cluster", "foo", "fp-default")
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(cmd.cmd.ClusterConfig.Metadata.Name).To(Equal("foo"))
			Expect(cmd.options.ProfileName).To(Equal("fp-default"))
		})

		It("requires a profile name, which can be provided via the --name flag", func() {
			cmd := newMockDeleteFargateProfileCmd("fargateprofile", "--cluster", "foo", "--name", "fp-default")
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(cmd.cmd.ClusterConfig.Metadata.Name).To(Equal("foo"))
			Expect(cmd.options.ProfileName).To(Equal("fp-default"))
		})

		It("supports the cluster name to be provided by a ClusterConfig file, but still requires a profile name, and if missing, prints an error and the usage", func() {
			cmd := newMockDeleteFargateProfileCmd("fargateprofile", "-f", "../../../examples/01-simple-cluster.yaml")
			out, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Error: invalid Fargate profile: empty name"))
			Expect(out).To(ContainSubstring("Usage:"))
		})

		It("supports the cluster name to be provided by a ClusterConfig file, and requires a profile name provided via the --name flag", func() {
			cmd := newMockDeleteFargateProfileCmd("fargateprofile", "-f", "../../../examples/01-simple-cluster.yaml", "--name", "fp-default")
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(cmd.cmd.ClusterConfig.Metadata.Name).To(Equal("cluster-1"))
			Expect(cmd.options.ProfileName).To(Equal("fp-default"))
		})
	})
})

func newMockDeleteFargateProfileCmd(args ...string) *mockDeleteFargateProfileCmd {
	mockCmd := &mockDeleteFargateProfileCmd{}
	grouping := cmdutils.NewGrouping()
	parentCmd := cmdutils.NewVerbCmd("delete", "", "")
	cmdutils.AddResourceCmd(grouping, parentCmd, func(cmd *cmdutils.Cmd) {
		deleteFargateProfileWithRunFunc(cmd, func(cmd *cmdutils.Cmd, options *fargate.Options) error {
			mockCmd.cmd = cmd
			mockCmd.options = options
			return nil // no-op, to only test input aggregation & validation.
		})
	})
	parentCmd.SetArgs(args)
	mockCmd.parentCmd = parentCmd
	return mockCmd
}

type mockDeleteFargateProfileCmd struct {
	parentCmd *cobra.Command
	cmd       *cmdutils.Cmd
	options   *fargate.Options
}

func (c mockDeleteFargateProfileCmd) execute() (string, error) {
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
