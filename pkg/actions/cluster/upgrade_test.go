package cluster

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var _ = Describe("upgrade cluster", func() {
	type upgradeCase struct {
		givenVersion           string
		eksVersion             string
		expectedUpgradeVersion string
		expectedUpgrade        bool
		expectedErrorText      string
	}

	DescribeTable("checks the specified version",
		func(c upgradeCase) {
			clusterMeta := api.ClusterMeta{
				Version: c.givenVersion,
			}
			result, err := requiresVersionUpgrade(&clusterMeta, c.eksVersion)

			if c.expectedErrorText != "" {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(c.expectedErrorText))
			} else {
				Expect(clusterMeta.Version).To(Equal(c.expectedUpgradeVersion))
				Expect(result).To(Equal(c.expectedUpgrade))
			}
		},

		Entry("upgrades by default when the version is not specified", upgradeCase{
			givenVersion:           "",
			eksVersion:             "1.16",
			expectedUpgradeVersion: "1.17",
			expectedUpgrade:        true,
		}),

		Entry("upgrades by default when the version is auto", upgradeCase{
			givenVersion:           "auto",
			eksVersion:             "1.16",
			expectedUpgradeVersion: "1.17",
			expectedUpgrade:        true,
		}),

		Entry("does not upgrade or fail when the cluster is already in the last version", upgradeCase{
			givenVersion:           "",
			eksVersion:             "1.19",
			expectedUpgradeVersion: "1.19",
			expectedUpgrade:        false,
		}),

		Entry("upgrades to the next version when specified", upgradeCase{
			givenVersion:           "1.16",
			eksVersion:             "1.15",
			expectedUpgradeVersion: "1.16",
			expectedUpgrade:        true,
		}),

		Entry("does not upgrade when the current version is specified", upgradeCase{
			givenVersion:           "1.15",
			eksVersion:             "1.15",
			expectedUpgradeVersion: "1.15",
			expectedUpgrade:        false,
		}),

		Entry("fails when the upgrade jumps more than one kubernetes version", upgradeCase{
			givenVersion:      "1.16",
			eksVersion:        "1.14",
			expectedErrorText: "upgrading more than one version at a time is not supported",
		}),

		Entry("fails when the given version is lower than the current one", upgradeCase{
			givenVersion:      "1.14",
			eksVersion:        "1.15",
			expectedErrorText: "cannot upgrade to a lower version. Found given target version \"1.14\", current cluster version \"1.15\"",
		}),

		Entry("fails when the version is deprecated", upgradeCase{
			givenVersion:      "1.12",
			eksVersion:        "1.12",
			expectedErrorText: "control plane version \"1.12\" has been deprecated",
		}),

		Entry("fails when the version is still not supported", upgradeCase{
			givenVersion:      "1.20",
			eksVersion:        "1.19",
			expectedErrorText: "control plane version \"1.20\" is not known to this version of eksctl",
		}),
	)
})
