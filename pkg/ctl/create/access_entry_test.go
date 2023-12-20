package create

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

var _ = Describe("create access entry", func() {
	type accessEntryTest struct {
		args        []string
		expectedErr string
	}

	DescribeTable("invalid arguments", func(aet accessEntryTest) {
		args := append([]string{"accessentry"}, aet.args...)
		cmd := newMockCmdWithRunFunc("create", func(cmd *cmdutils.Cmd) {
			createAccessEntryCmdWithRunFunc(cmd, func(cmd *cmdutils.Cmd) error {
				return nil
			})
		}, args...)
		_, err := cmd.execute()
		Expect(err).To(MatchError(ContainSubstring(aet.expectedErr)))
	},
		Entry("--cluster not supplied", accessEntryTest{
			expectedErr: "--cluster must be set",
		}),

		Entry("--principal-arn not supplied", accessEntryTest{
			args:        []string{"--cluster", "test"},
			expectedErr: "--principal-arn is required",
		}),

		Entry("invalid principal-arn", accessEntryTest{
			args:        []string{"--cluster", "test", "--principal-arn", "arn:invalid"},
			expectedErr: `invalid argument "arn:invalid" for "--principal-arn" flag: invalid ARN "arn:invalid"`,
		}),
	)
})
