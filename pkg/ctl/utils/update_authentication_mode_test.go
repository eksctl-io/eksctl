package utils_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("update authentication mode", func() {

	type updateAuthenticationModeEntry struct {
		args        []string
		expectedErr string
	}

	DescribeTable("unsupported arguments", func(e updateAuthenticationModeEntry) {
		cmd := newMockCmd(append([]string{"update-authentication-mode"}, e.args...)...)
		_, err := cmd.execute()
		Expect(err).To(MatchError(ContainSubstring(e.expectedErr)))
	},
		Entry("missing required flag --cluster", updateAuthenticationModeEntry{
			expectedErr: "Error: --cluster must be set",
		}),
		Entry("missing required flag --authentication-mode", updateAuthenticationModeEntry{
			args:        []string{"--cluster", "test"},
			expectedErr: "Error: --authentication-mode must be set",
		}),
		Entry("unsupported name argument", updateAuthenticationModeEntry{
			expectedErr: "Error: name argument is not supported",
			args:        []string{"--cluster", "test", "CONFIG_MAP"},
		}),
		Entry("setting --cluster and --config-file at the same time", updateAuthenticationModeEntry{
			expectedErr: "Error: cannot use --cluster when --config-file/-f is set",
			args:        []string{"--cluster", "test", "--config-file", "../../../examples/01-simple-cluster.yaml"},
		}),
		Entry("setting --authentication-mode and --config-file at the same time", updateAuthenticationModeEntry{
			expectedErr: "Error: cannot use --authentication-mode when --config-file/-f is set",
			args:        []string{"--authentication-mode", "CONFIG_MAP", "--config-file", "../../../examples/01-simple-cluster.yaml"},
		}),
	)
})
