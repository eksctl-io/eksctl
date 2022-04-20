package v1alpha5_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var _ = Describe("Addon", func() {
	Describe("Validate", func() {
		When("name is not set", func() {
			It("errors", func() {
				err := v1alpha5.Addon{}.Validate()
				Expect(err).To(MatchError("name required"))
			})
		})

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
