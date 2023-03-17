package get

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("get addon", func() {

	type getAddonEntry struct {
		args        []string
		expectedErr string
	}

	DescribeTable("unsupported arguments", func(e getAddonEntry) {
		cmd := newMockCmd(append([]string{"addon"}, e.args...)...)
		_, err := cmd.execute()
		Expect(err).To(MatchError(ContainSubstring(e.expectedErr)))
	},
		Entry("missing required flag --cluster", getAddonEntry{
			expectedErr: "Error: --cluster must be set",
		}),
		Entry("unsupported name argument", getAddonEntry{
			expectedErr: "Error: name argument is not supported",
			args:        []string{"--cluster", "test", "kube-proxy"},
		}),
		Entry("setting --name and --config-file at the same time", getAddonEntry{
			expectedErr: "Error: cannot use --name when --config-file/-f is set",
			args:        []string{"--name", "kube-proxy", "--config-file", "../../../examples/01-simple-cluster.yaml"},
		}),
		Entry("setting --cluster and --config-file at the same time", getAddonEntry{
			expectedErr: "Error: cannot use --cluster when --config-file/-f is set",
			args:        []string{"--cluster", "test", "--config-file", "../../../examples/01-simple-cluster.yaml"},
		}),
	)
})
