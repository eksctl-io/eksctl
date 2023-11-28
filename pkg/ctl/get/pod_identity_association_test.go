package get

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("get pod identity association", func() {

	type getPodIdentityAssociationEntry struct {
		args        []string
		expectedErr string
	}

	DescribeTable("unsupported arguments", func(e getPodIdentityAssociationEntry) {
		cmd := newMockCmd(append([]string{"podidentityassociation"}, e.args...)...)
		_, err := cmd.execute()
		Expect(err).To(MatchError(ContainSubstring(e.expectedErr)))
	},
		Entry("missing required flag --cluster", getPodIdentityAssociationEntry{
			expectedErr: "--cluster must be set",
		}),
		Entry("using --service-account-name without --namespace", getPodIdentityAssociationEntry{
			args:        []string{"--cluster", "test-cluster", "--service-account-name", "test-sa-name"},
			expectedErr: "--namespace must be set in order to specify --service-account-name",
		}),
	)
})
