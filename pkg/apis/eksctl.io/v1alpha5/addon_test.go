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
				expectedErr: "cannot set pod identity associtations for \"eks-pod-identity-agent\" addon",
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
