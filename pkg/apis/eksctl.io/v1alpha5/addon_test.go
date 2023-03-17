package v1alpha5_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var _ = Describe("Addon", func() {
	Describe("Validating configuration", func() {
		When("name is not set", func() {
			It("errors", func() {
				err := v1alpha5.Addon{}.Validate()
				Expect(err).To(MatchError("name is required"))
			})
		})

		DescribeTable("when configurationValues is in invalid format",
			func(configurationValues string) {
				err := v1alpha5.Addon{
					Name:                "name",
					Version:             "version",
					ConfigurationValues: configurationValues,
				}.Validate()
				Expect(err).To(MatchError(ContainSubstring("is not valid, supported format(s) are: JSON and YAML")))
			},
			Entry("non-empty string", "this a string not an object"),
			Entry("invalid yaml", "\"replicaCount: 1"),
		)

		DescribeTable("when configurationValues is in valid format",
			func(configurationValues string) {
				err := v1alpha5.Addon{
					Name:                "name",
					Version:             "version",
					ConfigurationValues: configurationValues,
				}.Validate()
				Expect(err).NotTo(HaveOccurred())
			},
			Entry("empty string", ""),
			Entry("empty json", "{}"),
			Entry("non-empty json", "{\"replicaCount\":3}"),
			Entry("non-empty yaml", "replicaCount: 3"),
		)

		When("specifying more than one of serviceAccountRoleARN, attachPolicyARNs, attachPolicy", func() {
			It("errors", func() {
				err := v1alpha5.Addon{
					Name:                  "name",
					Version:               "version",
					ServiceAccountRoleARN: "foo",
					AttachPolicyARNs:      []string{"arn"},
				}.Validate()
				Expect(err).To(MatchError("at most one of wellKnownPolicies, serviceAccountRoleARN, attachPolicyARNs and attachPolicy can be specified"))

				err = v1alpha5.Addon{
					Name:    "name",
					Version: "version",
					AttachPolicy: map[string]interface{}{
						"foo": "bar",
					},
					AttachPolicyARNs: []string{"arn"},
				}.Validate()
				Expect(err).To(MatchError("at most one of wellKnownPolicies, serviceAccountRoleARN, attachPolicyARNs and attachPolicy can be specified"))

				err = v1alpha5.Addon{
					Name:                  "name",
					Version:               "version",
					ServiceAccountRoleARN: "foo",
					AttachPolicy: map[string]interface{}{
						"foo": "bar",
					},
				}.Validate()
				Expect(err).To(MatchError("at most one of wellKnownPolicies, serviceAccountRoleARN, attachPolicyARNs and attachPolicy can be specified"))

				err = v1alpha5.Addon{
					Name:    "name",
					Version: "version",
					WellKnownPolicies: v1alpha5.WellKnownPolicies{
						AutoScaler: true,
					},
					AttachPolicy: map[string]interface{}{
						"foo": "bar",
					},
				}.Validate()
				Expect(err).To(MatchError("at most one of wellKnownPolicies, serviceAccountRoleARN, attachPolicyARNs and attachPolicy can be specified"))
			})
		})
	})
})
