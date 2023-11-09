package delete

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("delete access entry", func() {

	type deleteAccessEntryTest struct {
		args        []string
		expectedErr string
	}

	DescribeTable("unsupported arguments", func(e deleteAccessEntryTest) {
		cmd := newDefaultCmd(append([]string{"accessentry"}, e.args...)...)
		_, err := cmd.execute()
		Expect(err).To(MatchError(ContainSubstring(e.expectedErr)))
	},
		Entry("missing required flag --cluster", deleteAccessEntryTest{
			expectedErr: "Error: --cluster must be set",
		}),
		Entry("missing required flag --principal-arn", deleteAccessEntryTest{
			expectedErr: "Error: --principal-arn must be set",
			args:        []string{"--cluster", "test"},
		}),
		Entry("setting invalid value for --principal-arn", deleteAccessEntryTest{
			expectedErr: "Error: invalid argument \"invalid\" for \"--principal-arn\" flag",
			args:        []string{"--cluster", "test", "--principal-arn", "invalid"},
		}),
		Entry("setting --principal-arn and --config-file at the same time", deleteAccessEntryTest{
			expectedErr: "Error: cannot use --principal-arn when --config-file/-f is set",
			args:        []string{"--principal-arn", "arn:aws:iam::111122223333:user/my-user-name", "--config-file", "../../../examples/01-simple-cluster.yaml"},
		}),
		Entry("setting --cluster and --config-file at the same time", deleteAccessEntryTest{
			expectedErr: "Error: cannot use --cluster when --config-file/-f is set",
			args:        []string{"--cluster", "test", "--config-file", "../../../examples/01-simple-cluster.yaml"},
		}),
	)
})
