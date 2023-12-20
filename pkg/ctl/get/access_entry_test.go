package get

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("get access entry", func() {

	type getAccessEntryTest struct {
		args        []string
		expectedErr string
	}

	DescribeTable("unsupported arguments", func(e getAccessEntryTest) {
		cmd := newMockCmd(append([]string{"accessentry"}, e.args...)...)
		_, err := cmd.execute()
		Expect(err).To(MatchError(ContainSubstring(e.expectedErr)))
	},
		Entry("missing required flag --cluster", getAccessEntryTest{
			expectedErr: "Error: --cluster must be set",
		}),
		Entry("unsupported name argument", getAccessEntryTest{
			expectedErr: "Error: name argument is not supported",
			args:        []string{"--cluster", "test", "entry-name"},
		}),
		Entry("setting invalid value for --principal-arn", getAccessEntryTest{
			expectedErr: "Error: invalid argument \"invalid\" for \"--principal-arn\" flag",
			args:        []string{"--cluster", "test", "--principal-arn", "invalid"},
		}),
		Entry("setting --name and --config-file at the same time", getAccessEntryTest{
			expectedErr: "Error: cannot use --principal-arn when --config-file/-f is set",
			args:        []string{"--principal-arn", "arn:aws:iam::123456:role/testing-role", "--config-file", "../../../examples/01-simple-cluster.yaml"},
		}),
		Entry("setting --cluster and --config-file at the same time", getAccessEntryTest{
			expectedErr: "Error: cannot use --cluster when --config-file/-f is set",
			args:        []string{"--cluster", "test", "--config-file", "../../../examples/01-simple-cluster.yaml"},
		}),
	)
})
