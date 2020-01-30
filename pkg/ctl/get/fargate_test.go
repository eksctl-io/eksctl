package get

import (
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"testing"
)

func TestSuite(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("get", func() {
	Describe("get fargateprofile", func() {
		It("requires the cluster's name, and if missing, prints an error and the usage", func() {
			cmd := newMockGetFargateProfileCmd("fargateprofile")
			out, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--cluster must be set"))
			Expect(out).To(ContainSubstring("Error: --cluster must be set"))
			Expect(out).To(ContainSubstring("Usage:"))
		})

		It("requires the cluster's name, and does not have any profile name filter by default", func() {
			cmd := newMockGetFargateProfileCmd("fargateprofile", "--cluster", "foo")
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(cmd.cmd.ClusterConfig.Metadata.Name).To(Equal("foo"))
			Expect(cmd.options.ProfileName).To(Equal(""))
		})

		It("optionally accepts a profile name, which can be provided as an argument", func() {
			cmd := newMockGetFargateProfileCmd("fargateprofile", "--cluster", "foo", "fp-default")
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(cmd.cmd.ClusterConfig.Metadata.Name).To(Equal("foo"))
			Expect(cmd.options.ProfileName).To(Equal("fp-default"))
		})

		It("optionally accepts a profile name, which can be provided via the --name flag", func() {
			cmd := newMockGetFargateProfileCmd("fargateprofile", "--cluster", "foo", "--name", "fp-default")
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(cmd.cmd.ClusterConfig.Metadata.Name).To(Equal("foo"))
			Expect(cmd.options.ProfileName).To(Equal("fp-default"))
		})

		It("supports the cluster name to be provided by a ClusterConfig file", func() {
			cmd := newMockGetFargateProfileCmd("fargateprofile", "-f", "../../../examples/01-simple-cluster.yaml")
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(cmd.cmd.ClusterConfig.Metadata.Name).To(Equal("cluster-1"))
			Expect(cmd.options.ProfileName).To(Equal(""))
		})

		It("supports the cluster name to be provided by a ClusterConfig file, and an optional profile name provided via the --name flag", func() {
			cmd := newMockGetFargateProfileCmd("fargateprofile", "-f", "../../../examples/01-simple-cluster.yaml", "--name", "fp-default")
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(cmd.cmd.ClusterConfig.Metadata.Name).To(Equal("cluster-1"))
			Expect(cmd.options.ProfileName).To(Equal("fp-default"))
		})
	})
})

func newMockGetFargateProfileCmd(args ...string) *mockGetFargateProfileCmd {
	mockCmd := &mockGetFargateProfileCmd{}
	grouping := cmdutils.NewGrouping()
	parentCmd := cmdutils.NewVerbCmd("get", "", "")
	cmdutils.AddResourceCmd(grouping, parentCmd, func(cmd *cmdutils.Cmd) {
		getFargateProfileWithRunFunc(cmd, func(cmd *cmdutils.Cmd, options *options) error {
			mockCmd.cmd = cmd
			mockCmd.options = options
			return nil // no-op, to only test input aggregation & validation.
		})
	})
	parentCmd.SetArgs(args)
	mockCmd.parentCmd = parentCmd
	return mockCmd
}

type mockGetFargateProfileCmd struct {
	parentCmd *cobra.Command
	cmd       *cmdutils.Cmd
	options   *options
}

func (c mockGetFargateProfileCmd) execute() (string, error) {
	buf := new(bytes.Buffer)
	c.parentCmd.SetOut(buf)
	err := c.parentCmd.Execute()
	return buf.String(), err
}
