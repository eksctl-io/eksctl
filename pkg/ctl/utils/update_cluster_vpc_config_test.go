package utils_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

type updateClusterVPCEntry struct {
	args        []string
	expectedErr string
}

var _ = DescribeTable("invalid usage of update-cluster-vpc-config", func(e updateClusterVPCEntry) {
	cmd := newMockCmd(append([]string{"update-cluster-vpc-config"}, e.args...)...)
	_, err := cmd.execute()
	Expect(err).To(MatchError(ContainSubstring(e.expectedErr)))
},
	Entry("missing --cluster option", updateClusterVPCEntry{
		expectedErr: "--cluster must be set",
	}),

	Entry("missing a required parameter", updateClusterVPCEntry{
		args:        []string{"--cluster", "test"},
		expectedErr: "at least one of these options must be specified: --private-access, --public-access, --public-access-cidrs, --control-plane-subnet-ids, --control-plane-security-group-ids",
	}),
)

var _ = Describe("update-cluster-vpc-config with a config file", func() {
	writeConfigFile := func(vpc string) string {
		path := filepath.Join(GinkgoT().TempDir(), "config.yaml")
		contents := `apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig
metadata:
  name: test
  region: us-west-2
vpc:
` + vpc
		Expect(os.WriteFile(path, []byte(contents), 0600)).To(Succeed())
		return path
	}

	// load runs the loader for `eksctl utils update-cluster-vpc-config` directly, without
	// going through the command handler that talks to AWS.
	load := func(configFile string) error {
		cfg := api.NewClusterConfig()
		cmd := &cmdutils.Cmd{
			ClusterConfig:     cfg,
			ClusterConfigFile: configFile,
			CobraCommand: &cobra.Command{
				Use: "update-cluster-vpc-config",
				Run: func(_ *cobra.Command, _ []string) {},
			},
		}
		return cmdutils.NewUpdateClusterVPCLoader(cmd, cmdutils.UpdateClusterVPCOptions{}).Load()
	}

	When("vpc.controlPlaneOnPrivateSubnets is set", func() {
		It("does not return a validation error", func() {
			err := load(writeConfigFile("  controlPlaneOnPrivateSubnets: true\n"))
			Expect(err).NotTo(HaveOccurred())
		})
	})

	When("vpc.controlPlaneOnPrivateSubnets is explicitly false", func() {
		It("does not return a validation error", func() {
			err := load(writeConfigFile("  controlPlaneOnPrivateSubnets: false\n  controlPlaneSubnetIDs: [subnet-1234, subnet-5678]\n"))
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
