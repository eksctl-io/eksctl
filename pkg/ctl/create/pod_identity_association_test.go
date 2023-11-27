package create

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("create pod identity association", func() {

	var (
		defaultArgs = []string{
			"--cluster", "test-cluster",
			"--namespace", "test-namespace",
			"--service-account-name", "test-sa-name",
		}
		configFile = "../../../examples/01-simple-cluster.yaml"
	)

	type createPodIdentityAssociationEntry struct {
		args        []string
		expectedErr string
	}

	DescribeTable("unsupported arguments", func(e createPodIdentityAssociationEntry) {
		cmd := newDefaultCmd(append([]string{"podidentityassociation"}, e.args...)...)
		_, err := cmd.execute()
		Expect(err).To(MatchError(ContainSubstring(e.expectedErr)))
	},
		Entry("missing required flag --cluster", createPodIdentityAssociationEntry{
			expectedErr: "--cluster must be set",
		}),
		Entry("missing required flag --namespace", createPodIdentityAssociationEntry{
			args:        []string{"--cluster", "test-cluster"},
			expectedErr: "--namespace is required",
		}),
		Entry("missing required flag --service-account-name", createPodIdentityAssociationEntry{
			args:        []string{"--cluster", "test-cluster", "--namespace", "test-namespace"},
			expectedErr: "--service-account-name is required",
		}),
		Entry("setting --cluster and --config-file at the same time", createPodIdentityAssociationEntry{
			args:        []string{"--cluster", "test-cluster", "--config-file", configFile},
			expectedErr: "cannot use --cluster when --config-file/-f is set",
		}),
		Entry("setting --namespace and --config-file at the same time", createPodIdentityAssociationEntry{
			args:        []string{"--namespace", "test-namespace", "--config-file", configFile},
			expectedErr: "cannot use --namespace when --config-file/-f is set",
		}),
		Entry("setting --service-account-name and --config-file at the same time", createPodIdentityAssociationEntry{
			args:        []string{"--service-account-name", "test-sa-name", "--config-file", configFile},
			expectedErr: "cannot use --service-account-name when --config-file/-f is set",
		}),
		Entry("missing all --role-arn, --permission-policy-arns and --well-known-policies", createPodIdentityAssociationEntry{
			args:        defaultArgs,
			expectedErr: "at least one of the following flags must be specified: --role-arn, --permission-policy-arns, --well-known-policies",
		}),
		Entry("setting --permissions-policy-arns and --role-arn at the same time", createPodIdentityAssociationEntry{
			args:        append(defaultArgs, "--role-arn", "test-role", "--permission-policy-arns=test-policy"),
			expectedErr: "--permission-policy-arns cannot be specified when --role-arn is set",
		}),
		Entry("setting --well-known-policies and --role-arn at the same time", createPodIdentityAssociationEntry{
			args:        append(defaultArgs, "--role-arn", "test-role", "--well-known-policies=autoScaler,externalDNS"),
			expectedErr: "--well-known-policies cannot be specified when --role-arn is set",
		}),
		Entry("invalid --well-known-policies value", createPodIdentityAssociationEntry{
			args:        append(defaultArgs, "--well-known-policies=invalid"),
			expectedErr: "invalid wellKnownPolicy",
		}),
	)
})
