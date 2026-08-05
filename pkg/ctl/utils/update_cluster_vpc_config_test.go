package utils_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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

	When("vpc.controlPlaneOnPrivateSubnets is set", func() {
		It("returns an error instead of silently ignoring it", func() {
			cmd := newMockCmd("update-cluster-vpc-config", "-f", writeConfigFile("  controlPlaneOnPrivateSubnets: true\n"))
			_, err := cmd.execute()
			Expect(err).To(MatchError(ContainSubstring("vpc.controlPlaneOnPrivateSubnets is only supported when creating a cluster")))
		})
	})

	When("vpc.controlPlaneOnPrivateSubnets is explicitly false", func() {
		It("does not return a validation error", func() {
			cmd := newMockCmd("update-cluster-vpc-config", "-f", writeConfigFile("  controlPlaneOnPrivateSubnets: false\n  controlPlaneSubnetIDs: [subnet-1234, subnet-5678]\n"))
			_, err := cmd.execute()
			// The command proceeds past validation and fails when reaching AWS.
			if err != nil {
				Expect(err.Error()).NotTo(ContainSubstring("controlPlaneOnPrivateSubnets"))
			}
		})
	})
})
