package utils_test

import (
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
