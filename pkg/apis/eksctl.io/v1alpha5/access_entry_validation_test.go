package v1alpha5_test

import (
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

type accessEntryTest struct {
	authenticationMode ekstypes.AuthenticationMode
	accessEntries      []api.AccessEntry
	expectedErr        string
}

var _ = DescribeTable("Access Entry validation", func(aet accessEntryTest) {
	clusterConfig := api.NewClusterConfig()
	clusterConfig.AccessConfig = &api.AccessConfig{
		AccessEntries:      aet.accessEntries,
		AuthenticationMode: aet.authenticationMode,
	}
	err := api.ValidateClusterConfig(clusterConfig)
	if aet.expectedErr != "" {
		Expect(err).To(MatchError(ContainSubstring(aet.expectedErr)))
	} else {
		Expect(err).NotTo(HaveOccurred())
	}
},
	Entry("access entries specified when authentication mode is set to CONFIG_MAP", accessEntryTest{
		authenticationMode: ekstypes.AuthenticationModeConfigMap,
		accessEntries: []api.AccessEntry{
			{
				PrincipalARN: api.MustParseARN("arn:aws:iam::111122223333:role/role-1"),
				AccessPolicies: []api.AccessPolicy{
					{
						PolicyARN: api.MustParseARN("arn:aws:eks::aws:cluster-access-policy/AmazonEKSViewPolicy"),
						AccessScope: api.AccessScope{
							Type: ekstypes.AccessScopeTypeCluster,
						},
					},
				},
			},
		},

		expectedErr: "accessConfig.authenticationMode must be set to either API_AND_CONFIG_MAP or API to use access entries",
	}),

	Entry("empty principal ARN", accessEntryTest{
		authenticationMode: ekstypes.AuthenticationModeApiAndConfigMap,
		accessEntries: []api.AccessEntry{
			{},
		},

		expectedErr: "accessEntries[0].principalARN must be set to a valid AWS ARN",
	}),

	Entry("empty policy ARN", accessEntryTest{
		authenticationMode: ekstypes.AuthenticationModeApiAndConfigMap,
		accessEntries: []api.AccessEntry{
			{
				PrincipalARN: api.MustParseARN("arn:aws:iam::111122223333:role/role-1"),
				AccessPolicies: []api.AccessPolicy{
					{},
				},
			},
		},

		expectedErr: "accessEntries[0].policyARN must be set to a valid AWS ARN",
	}),

	Entry("empty accessScope.type", accessEntryTest{
		authenticationMode: ekstypes.AuthenticationModeApi,
		accessEntries: []api.AccessEntry{
			{
				PrincipalARN: api.MustParseARN("arn:aws:iam::111122223333:role/role-1"),
				AccessPolicies: []api.AccessPolicy{
					{
						PolicyARN: api.MustParseARN("arn:aws:eks::aws:cluster-access-policy/AmazonEKSViewPolicy"),
					},
				},
			},
		},

		expectedErr: `accessEntries[0].accessScope.type must be set to either "namespace" or "cluster"`,
	}),

	Entry("invalid type", accessEntryTest{
		authenticationMode: ekstypes.AuthenticationModeApiAndConfigMap,
		accessEntries: []api.AccessEntry{
			{
				PrincipalARN: api.MustParseARN("arn:aws:iam::111122223333:role/role-1"),
				Type:         "ec2_linux",
			},
		},

		expectedErr: `invalid access entry type "ec2_linux" for accessEntries[0]`,
	}),

	Entry("kubernetesGroups set for non-standard access entry type", accessEntryTest{
		authenticationMode: ekstypes.AuthenticationModeApiAndConfigMap,
		accessEntries: []api.AccessEntry{
			{
				PrincipalARN:     api.MustParseARN("arn:aws:iam::111122223333:role/role-1"),
				Type:             "FARGATE_LINUX",
				KubernetesGroups: []string{"dummy"},
			},
		},

		expectedErr: `cannot specify accessEntries[0].kubernetesGroups nor accessEntries[0].kubernetesUsername when type is set to FARGATE_LINUX`,
	}),

	Entry("kubernetesUsername set for non-standard access entry type", accessEntryTest{
		authenticationMode: ekstypes.AuthenticationModeApiAndConfigMap,
		accessEntries: []api.AccessEntry{
			{
				PrincipalARN:       api.MustParseARN("arn:aws:iam::111122223333:role/role-1"),
				Type:               "FARGATE_LINUX",
				KubernetesUsername: "dummy",
			},
		},

		expectedErr: `cannot specify accessEntries[0].kubernetesGroups nor accessEntries[0].kubernetesUsername when type is set to FARGATE_LINUX`,
	}),

	Entry("accessPolicies set for non-standard access entry type", accessEntryTest{
		authenticationMode: ekstypes.AuthenticationModeApiAndConfigMap,
		accessEntries: []api.AccessEntry{
			{
				PrincipalARN: api.MustParseARN("arn:aws:iam::111122223333:role/role-1"),
				Type:         "FARGATE_LINUX",
				AccessPolicies: []api.AccessPolicy{
					{
						PolicyARN: api.MustParseARN("arn:aws:eks::aws:cluster-access-policy/AmazonEKSViewPolicy"),
						AccessScope: api.AccessScope{
							Type:       ekstypes.AccessScopeTypeNamespace,
							Namespaces: []string{"default"},
						},
					},
				},
			},
		},

		expectedErr: `cannot specify accessEntries[0].accessPolicies when type is set to FARGATE_LINUX`,
	}),

	Entry("invalid accessScope.type", accessEntryTest{
		authenticationMode: ekstypes.AuthenticationModeApiAndConfigMap,
		accessEntries: []api.AccessEntry{
			{
				PrincipalARN: api.MustParseARN("arn:aws:iam::111122223333:role/role-1"),
				AccessPolicies: []api.AccessPolicy{
					{
						PolicyARN: api.MustParseARN("arn:aws:eks::aws:cluster-access-policy/AmazonEKSViewPolicy"),
						AccessScope: api.AccessScope{
							Type: "resource",
						},
					},
				},
			},
		},

		expectedErr: `invalid access scope type "resource" for accessEntries[0]`,
	}),

	Entry("namespaces set for cluster access scope", accessEntryTest{
		authenticationMode: ekstypes.AuthenticationModeApi,
		accessEntries: []api.AccessEntry{
			{
				PrincipalARN: api.MustParseARN("arn:aws:iam::111122223333:role/role-1"),
				AccessPolicies: []api.AccessPolicy{
					{
						PolicyARN: api.MustParseARN("arn:aws:eks::aws:cluster-access-policy/AmazonEKSViewPolicy"),
						AccessScope: api.AccessScope{
							Type:       ekstypes.AccessScopeTypeCluster,
							Namespaces: []string{"kube-system"},
						},
					},
				},
			},
		},

		expectedErr: "cannot specify accessEntries[0].accessScope.namespaces when accessScope is set to cluster",
	}),

	Entry("namespaces set for cluster access scope", accessEntryTest{
		authenticationMode: ekstypes.AuthenticationModeApiAndConfigMap,
		accessEntries: []api.AccessEntry{
			{
				PrincipalARN: api.MustParseARN("arn:aws:iam::111122223333:role/role-1"),
				AccessPolicies: []api.AccessPolicy{
					{
						PolicyARN: api.MustParseARN("arn:aws:eks::aws:cluster-access-policy/AmazonEKSViewPolicy"),
						AccessScope: api.AccessScope{
							Type: ekstypes.AccessScopeTypeNamespace,
						},
					},
				},
			},
		},

		expectedErr: "at least one namespace must be specified when accessScope is set to namespace: (accessEntries[0])",
	}),

	Entry("duplicate principal ARN", accessEntryTest{
		authenticationMode: ekstypes.AuthenticationModeApi,
		accessEntries: []api.AccessEntry{
			{
				PrincipalARN: api.MustParseARN("arn:aws:iam::111122223333:role/role-1"),
				AccessPolicies: []api.AccessPolicy{
					{
						PolicyARN: api.MustParseARN("arn:aws:eks::aws:cluster-access-policy/AmazonEKSViewPolicy"),
						AccessScope: api.AccessScope{
							Type: ekstypes.AccessScopeTypeCluster,
						},
					},
				},
			},
			{
				PrincipalARN: api.MustParseARN("arn:aws:iam::111122223333:role/role-2"),
				AccessPolicies: []api.AccessPolicy{
					{
						PolicyARN: api.MustParseARN("arn:aws:eks::aws:cluster-access-policy/AmazonEKSViewPolicy"),
						AccessScope: api.AccessScope{
							Type: ekstypes.AccessScopeTypeCluster,
						},
					},
				},
			},
			{
				PrincipalARN: api.MustParseARN("arn:aws:iam::111122223333:role/role-1"),
				AccessPolicies: []api.AccessPolicy{
					{
						PolicyARN: api.MustParseARN("arn:aws:eks::aws:cluster-access-policy/AmazonEKSViewPolicy"),
						AccessScope: api.AccessScope{
							Type:       ekstypes.AccessScopeTypeNamespace,
							Namespaces: []string{"default"},
						},
					},
				},
			},
		},

		expectedErr: `duplicate access entry accessEntries[2] with principal ARN "arn:aws:iam::111122223333:role/role-1"`,
	}),

	Entry("valid access entries", accessEntryTest{
		authenticationMode: ekstypes.AuthenticationModeApiAndConfigMap,
		accessEntries: []api.AccessEntry{
			{
				PrincipalARN: api.MustParseARN("arn:aws:iam::111122223333:role/role-1"),
				AccessPolicies: []api.AccessPolicy{
					{
						PolicyARN: api.MustParseARN("arn:aws:eks::aws:cluster-access-policy/AmazonEKSViewPolicy"),
						AccessScope: api.AccessScope{
							Type: ekstypes.AccessScopeTypeCluster,
						},
					},
				},
			},
			{
				PrincipalARN: api.MustParseARN("arn:aws:iam::111122223333:role/role-2"),
				AccessPolicies: []api.AccessPolicy{
					{
						PolicyARN: api.MustParseARN("arn:aws:eks::aws:cluster-access-policy/AmazonEKSViewPolicy"),
						AccessScope: api.AccessScope{
							Type: ekstypes.AccessScopeTypeCluster,
						},
					},
				},
			},
			{
				PrincipalARN: api.MustParseARN("arn:aws:iam::111122223333:role/role-3"),
			},
			{
				PrincipalARN: api.MustParseARN("arn:aws:iam::111122223333:role/role-4"),
				Type:         "EC2_LINUX",
			},
			{
				PrincipalARN:       api.MustParseARN("arn:aws:iam::111122223333:role/role-5"),
				Type:               "STANDARD",
				KubernetesGroups:   []string{"dummy", "dummy"},
				KubernetesUsername: "dummy",
			},
			{
				PrincipalARN: api.MustParseARN("arn:aws:iam::111122223333:role/role-6"),
				AccessPolicies: []api.AccessPolicy{
					{
						PolicyARN: api.MustParseARN("arn:aws:eks::aws:cluster-access-policy/AmazonEKSViewPolicy"),
						AccessScope: api.AccessScope{
							Type:       ekstypes.AccessScopeTypeNamespace,
							Namespaces: []string{"default"},
						},
					},
				},
			},
		},
	}),
)
