package create

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

var _ = Describe("create", func() {
	Describe("create fargateprofile", func() {
		It("requires the cluster's name, and if missing, prints an error and the usage", func() {
			cmd := newMockCreateFargateProfileCmd("fargateprofile")
			out, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--cluster must be set"))
			Expect(out).To(ContainSubstring("Error: --cluster must be set"))
			Expect(out).To(ContainSubstring("Usage:"))
		})

		It("requires a Kubernetes namespace to be provided, and if missing, prints an error and the usage", func() {
			cmd := newMockCreateFargateProfileCmd("fargateprofile", "--cluster", "foo")
			out, err := cmd.execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("invalid Fargate profile: empty selector namespace"))
			Expect(out).To(ContainSubstring("Error: invalid Fargate profile: empty selector namespace"))
			Expect(out).To(ContainSubstring("Usage:"))
		})

		It("requires at least a Kubernetes namespace to be provided, in which case it generates the profile name", func() {
			cmd := newMockCreateFargateProfileCmd("fargateprofile", "--cluster", "foo", "--namespace", "default")
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(cmd.cmd.ClusterConfig.Metadata.Name).To(Equal("foo"))
			profiles := cmd.cmd.ClusterConfig.FargateProfiles
			Expect(profiles).To(HaveLen(1))
			profile := profiles[0]
			Expect(profile.Name).To(MatchRegexp("fp-[abcdef0123456789]{8}"))
			Expect(profile.Selectors).To(HaveLen(1))
			selector := profile.Selectors[0]
			Expect(selector.Namespace).To(Equal("default"))
		})

		It("the profile name can be provided as an argument", func() {
			cmd := newMockCreateFargateProfileCmd("fargateprofile", "--cluster", "foo", "--namespace", "default", "fp-default")
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(cmd.cmd.ClusterConfig.Metadata.Name).To(Equal("foo"))
			profiles := cmd.cmd.ClusterConfig.FargateProfiles
			Expect(profiles).To(HaveLen(1))
			profile := profiles[0]
			Expect(profile.Name).To(Equal("fp-default"))
			Expect(profile.Selectors).To(HaveLen(1))
			selector := profile.Selectors[0]
			Expect(selector.Namespace).To(Equal("default"))
		})

		It("the profile name can be provided via the --name flag", func() {
			cmd := newMockCreateFargateProfileCmd("fargateprofile", "--cluster", "foo", "--namespace", "default", "--name", "fp-default")
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(cmd.cmd.ClusterConfig.Metadata.Name).To(Equal("foo"))
			profiles := cmd.cmd.ClusterConfig.FargateProfiles
			Expect(profiles).To(HaveLen(1))
			profile := profiles[0]
			Expect(profile.Name).To(Equal("fp-default"))
			Expect(profile.Selectors).To(HaveLen(1))
			selector := profile.Selectors[0]
			Expect(selector.Namespace).To(Equal("default"))
		})

		It("the fargate profile with tags", func() {
			cmd := newMockCreateFargateProfileCmd("fargateprofile", "--cluster", "foo", "--tags", "env=dev,name=fp-default", "--namespace", "default", "fp-default")
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(cmd.cmd.ClusterConfig.Metadata.Name).To(Equal("foo"))
			profiles := cmd.cmd.ClusterConfig.FargateProfiles
			Expect(profiles).To(HaveLen(1))
			profile := profiles[0]
			Expect(profile.Name).To(Equal("fp-default"))
			Expect(profile.Selectors).To(HaveLen(1))
			selector := profile.Selectors[0]
			Expect(selector.Namespace).To(Equal("default"))
			Expect(profile.Tags).To(HaveKeyWithValue("env", "dev"))
			Expect(profile.Tags).To(HaveKeyWithValue("name", "fp-default"))
		})

		It("supports all arguments to be provided by a ClusterConfig file", func() {
			cmd := newMockCreateFargateProfileCmd("fargateprofile", "-f", "../../../examples/16-fargate-profile.yaml")
			_, err := cmd.execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(cmd.cmd.ClusterConfig.Metadata.Name).To(Equal("cluster-16"))
			profiles := cmd.cmd.ClusterConfig.FargateProfiles
			Expect(profiles).To(HaveLen(2))
			Expect(profiles[0].Name).To(Equal("fp-default"))
			Expect(profiles[0].Selectors).To(HaveLen(2))
			Expect(profiles[0].Selectors[0].Namespace).To(Equal("default"))
			Expect(profiles[0].Selectors[0].Labels).To(BeEmpty())
			Expect(profiles[0].Selectors[1].Namespace).To(Equal("kube-system"))
			Expect(profiles[0].Selectors[1].Labels).To(BeEmpty())
			Expect(profiles[1].Name).To(Equal("fp-dev"))
			Expect(profiles[1].Selectors).To(HaveLen(1))
			Expect(profiles[1].Selectors[0].Namespace).To(Equal("dev"))
			Expect(profiles[1].Selectors[0].Labels).To(HaveLen(2))
			Expect(profiles[1].Selectors[0].Labels).To(HaveKeyWithValue("env", "dev"))
			Expect(profiles[1].Selectors[0].Labels).To(HaveKeyWithValue("checks", "passed"))
			Expect(profiles[1].Tags).To(HaveKeyWithValue("env", "dev"))
			Expect(profiles[1].Tags).To(HaveKeyWithValue("name", "fp-dev"))
		})
	})
})

func newMockCreateFargateProfileCmd(args ...string) *mockCreateFargateProfileCmd {
	mockCmd := &mockCreateFargateProfileCmd{}
	grouping := cmdutils.NewGrouping()
	parentCmd := cmdutils.NewVerbCmd("create", "", "")
	cmdutils.AddResourceCmd(grouping, parentCmd, func(cmd *cmdutils.Cmd) {
		createFargateProfileWithRunFunc(cmd, func(cmd *cmdutils.Cmd) error {
			mockCmd.cmd = cmd
			return nil // no-op, to only test input aggregation & validation.
		})
	})
	parentCmd.SetArgs(args)
	mockCmd.parentCmd = parentCmd
	return mockCmd
}

type mockCreateFargateProfileCmd struct {
	parentCmd *cobra.Command
	cmd       *cmdutils.Cmd
}

func (c mockCreateFargateProfileCmd) execute() (string, error) {
	buf := new(bytes.Buffer)
	c.parentCmd.SetOut(buf)
	err := c.parentCmd.Execute()
	return buf.String(), err
}
