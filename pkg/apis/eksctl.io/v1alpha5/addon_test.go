package v1alpha5_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var _ = Describe("Addon", func() {
	Describe("Validating configuration", func() {
		When("name is not set", func() {
			It("errors", func() {
				err := api.Addon{}.Validate()
				Expect(err).To(MatchError(ContainSubstring("name is required")))
			})
		})

		DescribeTable("when configurationValues is in invalid format",
			func(configurationValues string) {
				err := api.Addon{
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
				err := api.Addon{
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

		DescribeTable("when namespace config is valid",
			func(namespaceConfig *api.AddonNamespaceConfig) {
				err := api.Addon{
					Name:            "name",
					Version:         "version",
					NamespaceConfig: namespaceConfig,
				}.Validate()
				Expect(err).NotTo(HaveOccurred())
			},
			Entry("nil namespace config", (*api.AddonNamespaceConfig)(nil)),
			Entry("empty namespace config", &api.AddonNamespaceConfig{}),
			Entry("empty namespace string", &api.AddonNamespaceConfig{Namespace: ""}),
			Entry("valid namespace name", &api.AddonNamespaceConfig{Namespace: "kube-system"}),
			Entry("valid namespace with numbers", &api.AddonNamespaceConfig{Namespace: "app-v2"}),
			Entry("valid single character namespace", &api.AddonNamespaceConfig{Namespace: "a"}),
			Entry("valid namespace with hyphens", &api.AddonNamespaceConfig{Namespace: "my-app-namespace"}),
			Entry("valid namespace starting with letter", &api.AddonNamespaceConfig{Namespace: "system"}),
			Entry("valid namespace ending with number", &api.AddonNamespaceConfig{Namespace: "app1"}),
			Entry("valid namespace with multiple hyphens", &api.AddonNamespaceConfig{Namespace: "my-long-app-namespace"}),
			Entry("valid two character namespace", &api.AddonNamespaceConfig{Namespace: "ab"}),
			Entry("valid namespace with alternating letters and numbers", &api.AddonNamespaceConfig{Namespace: "a1b2c3"}),
			Entry("valid namespace at maximum length", &api.AddonNamespaceConfig{Namespace: "a12345678901234567890123456789012345678901234567890123456789012"}),
		)

		DescribeTable("when namespace config is invalid",
			func(namespaceConfig *api.AddonNamespaceConfig, expectedError string) {
				err := api.Addon{
					Name:            "name",
					Version:         "version",
					NamespaceConfig: namespaceConfig,
				}.Validate()
				Expect(err).To(MatchError(ContainSubstring(expectedError)))
			},
			Entry("namespace starting with hyphen", &api.AddonNamespaceConfig{Namespace: "-invalid"}, "is not a valid Kubernetes namespace name"),
			Entry("namespace ending with hyphen", &api.AddonNamespaceConfig{Namespace: "invalid-"}, "is not a valid Kubernetes namespace name"),
			Entry("namespace with uppercase letters", &api.AddonNamespaceConfig{Namespace: "Invalid"}, "is not a valid Kubernetes namespace name"),
			Entry("namespace with special characters", &api.AddonNamespaceConfig{Namespace: "invalid_namespace"}, "is not a valid Kubernetes namespace name"),
			Entry("namespace with dots", &api.AddonNamespaceConfig{Namespace: "invalid.namespace"}, "is not a valid Kubernetes namespace name"),
			Entry("namespace starting with number followed by hyphen", &api.AddonNamespaceConfig{Namespace: "1-invalid"}, "is not a valid Kubernetes namespace name"),
			Entry("namespace too long", &api.AddonNamespaceConfig{Namespace: "this-is-a-very-long-namespace-name-that-exceeds-the-maximum-length-allowed-for-kubernetes-namespaces"}, "is too long"),
			Entry("namespace starting with number", &api.AddonNamespaceConfig{Namespace: "1invalid"}, "is not a valid Kubernetes namespace name"),
			Entry("namespace with spaces", &api.AddonNamespaceConfig{Namespace: "invalid namespace"}, "is not a valid Kubernetes namespace name"),
			Entry("namespace with forward slash", &api.AddonNamespaceConfig{Namespace: "invalid/namespace"}, "is not a valid Kubernetes namespace name"),
			Entry("namespace with backslash", &api.AddonNamespaceConfig{Namespace: "invalid\\namespace"}, "is not a valid Kubernetes namespace name"),
			Entry("namespace with colon", &api.AddonNamespaceConfig{Namespace: "invalid:namespace"}, "is not a valid Kubernetes namespace name"),
			Entry("namespace with semicolon", &api.AddonNamespaceConfig{Namespace: "invalid;namespace"}, "is not a valid Kubernetes namespace name"),
			Entry("namespace with at symbol", &api.AddonNamespaceConfig{Namespace: "invalid@namespace"}, "is not a valid Kubernetes namespace name"),
			Entry("namespace with hash symbol", &api.AddonNamespaceConfig{Namespace: "invalid#namespace"}, "is not a valid Kubernetes namespace name"),
			Entry("namespace with percent symbol", &api.AddonNamespaceConfig{Namespace: "invalid%namespace"}, "is not a valid Kubernetes namespace name"),
			Entry("namespace with ampersand", &api.AddonNamespaceConfig{Namespace: "invalid&namespace"}, "is not a valid Kubernetes namespace name"),
			Entry("namespace with asterisk", &api.AddonNamespaceConfig{Namespace: "invalid*namespace"}, "is not a valid Kubernetes namespace name"),
			Entry("namespace with plus sign", &api.AddonNamespaceConfig{Namespace: "invalid+namespace"}, "is not a valid Kubernetes namespace name"),
			Entry("namespace with equals sign", &api.AddonNamespaceConfig{Namespace: "invalid=namespace"}, "is not a valid Kubernetes namespace name"),
			Entry("namespace with question mark", &api.AddonNamespaceConfig{Namespace: "invalid?namespace"}, "is not a valid Kubernetes namespace name"),
			Entry("namespace with exclamation mark", &api.AddonNamespaceConfig{Namespace: "invalid!namespace"}, "is not a valid Kubernetes namespace name"),
			Entry("namespace with parentheses", &api.AddonNamespaceConfig{Namespace: "invalid(namespace)"}, "is not a valid Kubernetes namespace name"),
			Entry("namespace with brackets", &api.AddonNamespaceConfig{Namespace: "invalid[namespace]"}, "is not a valid Kubernetes namespace name"),
			Entry("namespace with braces", &api.AddonNamespaceConfig{Namespace: "invalid{namespace}"}, "is not a valid Kubernetes namespace name"),
			Entry("namespace with pipe symbol", &api.AddonNamespaceConfig{Namespace: "invalid|namespace"}, "is not a valid Kubernetes namespace name"),
			Entry("namespace with tilde", &api.AddonNamespaceConfig{Namespace: "invalid~namespace"}, "is not a valid Kubernetes namespace name"),
			Entry("namespace with backtick", &api.AddonNamespaceConfig{Namespace: "invalid`namespace"}, "is not a valid Kubernetes namespace name"),
			Entry("namespace with single quote", &api.AddonNamespaceConfig{Namespace: "invalid'namespace"}, "is not a valid Kubernetes namespace name"),
			Entry("namespace with double quote", &api.AddonNamespaceConfig{Namespace: "invalid\"namespace"}, "is not a valid Kubernetes namespace name"),
			Entry("namespace with comma", &api.AddonNamespaceConfig{Namespace: "invalid,namespace"}, "is not a valid Kubernetes namespace name"),
			Entry("namespace with less than", &api.AddonNamespaceConfig{Namespace: "invalid<namespace"}, "is not a valid Kubernetes namespace name"),
			Entry("namespace with greater than", &api.AddonNamespaceConfig{Namespace: "invalid>namespace"}, "is not a valid Kubernetes namespace name"),
			Entry("namespace with tab character", &api.AddonNamespaceConfig{Namespace: "invalid\tnamespace"}, "is not a valid Kubernetes namespace name"),
			Entry("namespace with newline character", &api.AddonNamespaceConfig{Namespace: "invalid\nnamespace"}, "is not a valid Kubernetes namespace name"),
			Entry("namespace with carriage return", &api.AddonNamespaceConfig{Namespace: "invalid\rnamespace"}, "is not a valid Kubernetes namespace name"),
			Entry("namespace exactly 64 characters", &api.AddonNamespaceConfig{Namespace: "a1234567890123456789012345678901234567890123456789012345678901234"}, "is too long"),
			Entry("namespace much longer than limit", &api.AddonNamespaceConfig{Namespace: "this-is-an-extremely-long-namespace-name-that-definitely-exceeds-the-sixty-three-character-limit-for-kubernetes-namespace-names-and-should-fail-validation"}, "is too long"),
		)

		Describe("namespace config edge cases", func() {
			When("namespace config is nil", func() {
				It("should not error", func() {
					err := api.Addon{
						Name:            "test-addon",
						Version:         "v1.0.0",
						NamespaceConfig: nil,
					}.Validate()
					Expect(err).NotTo(HaveOccurred())
				})
			})

			When("namespace config has empty namespace string", func() {
				It("should not error", func() {
					err := api.Addon{
						Name:            "test-addon",
						Version:         "v1.0.0",
						NamespaceConfig: &api.AddonNamespaceConfig{Namespace: ""},
					}.Validate()
					Expect(err).NotTo(HaveOccurred())
				})
			})

			When("namespace config is provided with other valid addon fields", func() {
				It("should validate successfully", func() {
					err := api.Addon{
						Name:            "test-addon",
						Version:         "v1.0.0",
						NamespaceConfig: &api.AddonNamespaceConfig{Namespace: "custom-namespace"},
						Tags: map[string]string{
							"Environment": "test",
						},
						ConfigurationValues: `{"replicaCount": 2}`,
					}.Validate()
					Expect(err).NotTo(HaveOccurred())
				})
			})

			When("namespace config validation fails", func() {
				It("should include addon name in error message", func() {
					err := api.Addon{
						Name:            "my-test-addon",
						Version:         "v1.0.0",
						NamespaceConfig: &api.AddonNamespaceConfig{Namespace: "Invalid-Namespace"},
					}.Validate()
					Expect(err).To(MatchError(ContainSubstring("invalid configuration for \"my-test-addon\" addon")))
					Expect(err).To(MatchError(ContainSubstring("is not a valid Kubernetes namespace name")))
				})
			})

			When("namespace config has whitespace", func() {
				It("should fail validation", func() {
					err := api.Addon{
						Name:            "test-addon",
						Version:         "v1.0.0",
						NamespaceConfig: &api.AddonNamespaceConfig{Namespace: " invalid-namespace "},
					}.Validate()
					Expect(err).To(MatchError(ContainSubstring("is not a valid Kubernetes namespace name")))
				})
			})

			When("namespace config has leading whitespace", func() {
				It("should fail validation", func() {
					err := api.Addon{
						Name:            "test-addon",
						Version:         "v1.0.0",
						NamespaceConfig: &api.AddonNamespaceConfig{Namespace: " valid-namespace"},
					}.Validate()
					Expect(err).To(MatchError(ContainSubstring("is not a valid Kubernetes namespace name")))
				})
			})

			When("namespace config has trailing whitespace", func() {
				It("should fail validation", func() {
					err := api.Addon{
						Name:            "test-addon",
						Version:         "v1.0.0",
						NamespaceConfig: &api.AddonNamespaceConfig{Namespace: "valid-namespace "},
					}.Validate()
					Expect(err).To(MatchError(ContainSubstring("is not a valid Kubernetes namespace name")))
				})
			})
		})

		When("specifying more than one of serviceAccountRoleARN, attachPolicyARNs, attachPolicy, wellKnownPolicies", func() {
			It("errors", func() {
				err := api.Addon{
					Name:                  "name",
					Version:               "version",
					ServiceAccountRoleARN: "foo",
					AttachPolicyARNs:      []string{"arn"},
				}.Validate()
				Expect(err).To(MatchError(ContainSubstring("at most one of wellKnownPolicies, serviceAccountRoleARN, attachPolicyARNs and attachPolicy can be specified")))

				err = api.Addon{
					Name:    "name",
					Version: "version",
					AttachPolicy: map[string]interface{}{
						"foo": "bar",
					},
					AttachPolicyARNs: []string{"arn"},
				}.Validate()
				Expect(err).To(MatchError(ContainSubstring("at most one of wellKnownPolicies, serviceAccountRoleARN, attachPolicyARNs and attachPolicy can be specified")))

				err = api.Addon{
					Name:                  "name",
					Version:               "version",
					ServiceAccountRoleARN: "foo",
					AttachPolicy: map[string]interface{}{
						"foo": "bar",
					},
				}.Validate()
				Expect(err).To(MatchError(ContainSubstring("at most one of wellKnownPolicies, serviceAccountRoleARN, attachPolicyARNs and attachPolicy can be specified")))

				err = api.Addon{
					Name:    "name",
					Version: "version",
					WellKnownPolicies: api.WellKnownPolicies{
						AutoScaler: true,
					},
					AttachPolicy: map[string]interface{}{
						"foo": "bar",
					},
				}.Validate()
				Expect(err).To(MatchError(ContainSubstring("at most one of wellKnownPolicies, serviceAccountRoleARN, attachPolicyARNs and attachPolicy can be specified")))
			})
		})

		type addonWithPodIDEntry struct {
			addon       api.Addon
			expectedErr string
		}
		DescribeTable("pod identity associations", func(e addonWithPodIDEntry) {
			err := e.addon.Validate()
			if e.expectedErr != "" {
				Expect(err).To(MatchError(ContainSubstring(e.expectedErr)))
			} else {
				Expect(err).NotTo(HaveOccurred())
			}
		},
			Entry("setting podIDs for eks-pod-identity-agent addon", addonWithPodIDEntry{
				addon: api.Addon{
					Name:                    api.PodIdentityAgentAddon,
					PodIdentityAssociations: &[]api.PodIdentityAssociation{{}},
				},
				expectedErr: "cannot set pod identity associations for \"eks-pod-identity-agent\" addon",
			}),
			Entry("namespace is not set", addonWithPodIDEntry{
				addon: api.Addon{
					Name:                    "name",
					PodIdentityAssociations: &[]api.PodIdentityAssociation{{}},
				},
				expectedErr: "podIdentityAssociations[0].namespace must be set",
			}),
			Entry("service account name is not set", addonWithPodIDEntry{
				addon: api.Addon{
					Name: "name",
					PodIdentityAssociations: &[]api.PodIdentityAssociation{{
						Namespace: "kube-system",
					}},
				},
				expectedErr: "podIdentityAssociations[0].serviceAccountName must be set",
			}),
			Entry("no IAM role or policies are set", addonWithPodIDEntry{
				addon: api.Addon{
					Name: "name",
					PodIdentityAssociations: &[]api.PodIdentityAssociation{
						{
							Namespace:          "kube-system",
							ServiceAccountName: "aws-node",
						},
					},
				},
				expectedErr: fmt.Sprintf("at least one of the following must be specified: %[1]s.roleARN, %[1]s.permissionPolicy, %[1]s.permissionPolicyARNs, %[1]s.wellKnownPolicies", "podIdentityAssociations[0]"),
			}),
			Entry("IAM role and permissionPolicy are set at the same time", addonWithPodIDEntry{
				addon: api.Addon{
					Name: "name",
					PodIdentityAssociations: &[]api.PodIdentityAssociation{
						{
							Namespace:          "kube-system",
							ServiceAccountName: "aws-node",
							RoleARN:            "arn:aws:iam::111122223333:role/role-name-1",
							PermissionPolicy: api.InlineDocument{
								"Version": "2012-10-17",
								"Statement": []map[string]interface{}{
									{
										"Effect": "Allow",
										"Action": []string{
											"autoscaling:DescribeAutoScalingGroups",
											"autoscaling:DescribeAutoScalingInstances",
										},
										"Resource": "*",
									},
								},
							},
						},
					},
				},
				expectedErr: fmt.Sprintf("%[1]s.permissionPolicy cannot be specified when %[1]s.roleARN is set", "podIdentityAssociations[0]"),
			}),
			Entry("IAM role and permissionPolicyARNs are set at the same time", addonWithPodIDEntry{
				addon: api.Addon{
					Name: "name",
					PodIdentityAssociations: &[]api.PodIdentityAssociation{
						{
							Namespace:            "kube-system",
							ServiceAccountName:   "aws-node",
							RoleARN:              "arn:aws:iam::111122223333:role/role-name-1",
							PermissionPolicyARNs: []string{"arn:aws:iam::111122223333:policy/policy-name-1"},
						},
					},
				},
				expectedErr: fmt.Sprintf("%[1]s.permissionPolicyARNs cannot be specified when %[1]s.roleARN is set", "podIdentityAssociations[0]"),
			}),
			Entry("IAM role and permissionPolicyARNs are set at the same time", addonWithPodIDEntry{
				addon: api.Addon{
					Name: "name",
					PodIdentityAssociations: &[]api.PodIdentityAssociation{
						{
							Namespace:          "kube-system",
							ServiceAccountName: "aws-node",
							RoleARN:            "arn:aws:iam::111122223333:role/role-name-1",
							WellKnownPolicies: api.WellKnownPolicies{
								EBSCSIController: true,
							},
						},
					},
				},
				expectedErr: fmt.Sprintf("%[1]s.wellKnownPolicies cannot be specified when %[1]s.roleARN is set", "podIdentityAssociations[0]"),
			}),
			Entry("podIDs and ServiceAccountRoleARN are set at the same time", addonWithPodIDEntry{
				addon: api.Addon{
					Name:                  "name",
					ServiceAccountRoleARN: "arn:aws:iam::111122223333:role/role-name-1",
					PodIdentityAssociations: &[]api.PodIdentityAssociation{
						{
							Namespace:          "kube-system",
							ServiceAccountName: "aws-node",
							RoleARN:            "arn:aws:iam::111122223333:role/role-name-1",
						},
					},
				},
				expectedErr: "cannot set IRSA config (`addon.ServiceAccountRoleARN`, `addon.AttachPolicyARNs`, `addon.AttachPolicy`, `addon.WellKnownPolicies`) and pod identity associations at the same time",
			}),
			Entry("podIDs and AttachPolicyARNs are set at the same time", addonWithPodIDEntry{
				addon: api.Addon{
					Name:             "name",
					AttachPolicyARNs: []string{"arn:aws:iam::111122223333:policy/policy-name-1"},
					PodIdentityAssociations: &[]api.PodIdentityAssociation{
						{
							Namespace:          "kube-system",
							ServiceAccountName: "aws-node",
							RoleARN:            "arn:aws:iam::111122223333:role/role-name-1",
						},
					},
				},
				expectedErr: "cannot set IRSA config (`addon.ServiceAccountRoleARN`, `addon.AttachPolicyARNs`, `addon.AttachPolicy`, `addon.WellKnownPolicies`) and pod identity associations at the same time",
			}),
			Entry("podIDs and AttachPolicy are set at the same time", addonWithPodIDEntry{
				addon: api.Addon{
					Name: "name",
					AttachPolicy: api.InlineDocument{
						"Version": "2012-10-17",
						"Statement": []map[string]interface{}{
							{
								"Effect": "Allow",
								"Action": []string{
									"autoscaling:DescribeAutoScalingGroups",
									"autoscaling:DescribeAutoScalingInstances",
								},
								"Resource": "*",
							},
						},
					},
					PodIdentityAssociations: &[]api.PodIdentityAssociation{
						{
							Namespace:          "kube-system",
							ServiceAccountName: "aws-node",
							RoleARN:            "arn:aws:iam::111122223333:role/role-name-1",
						},
					},
				},
				expectedErr: "cannot set IRSA config (`addon.ServiceAccountRoleARN`, `addon.AttachPolicyARNs`, `addon.AttachPolicy`, `addon.WellKnownPolicies`) and pod identity associations at the same time",
			}),
			Entry("podIDs and WellKnownPolicies are set at the same time", addonWithPodIDEntry{
				addon: api.Addon{
					Name: "name",
					WellKnownPolicies: api.WellKnownPolicies{
						EBSCSIController: true,
					},
					PodIdentityAssociations: &[]api.PodIdentityAssociation{
						{
							Namespace:          "kube-system",
							ServiceAccountName: "aws-node",
							RoleARN:            "arn:aws:iam::111122223333:role/role-name-1",
						},
					},
				},
				expectedErr: "cannot set IRSA config (`addon.ServiceAccountRoleARN`, `addon.AttachPolicyARNs`, `addon.AttachPolicy`, `addon.WellKnownPolicies`) and pod identity associations at the same time",
			}),
		)
	})
})
