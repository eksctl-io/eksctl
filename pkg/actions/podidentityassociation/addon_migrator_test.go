package podidentityassociation_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/smithy-go"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/stretchr/testify/mock"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/actions/podidentityassociation"
	"github.com/weaveworks/eksctl/pkg/actions/podidentityassociation/mocks"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/eks/mocksv2"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

var _ = Describe("Addon Migration", func() {
	type addonMocks struct {
		eksAddonsAPI   *mocksv2.EKS
		iamRoleGetter  *mocksv2.IAM
		stackDescriber *mocks.StackDeleter
		roleMigrator   *mocks.RoleMigrator
	}
	type migrateEntry struct {
		mockCalls func(m addonMocks)

		expectedTasks bool
		expectedErr   string
	}

	const clusterName = "cluster"

	mockAddonCalls := func(eksAddonsAPI *mocksv2.EKS) {
		eksAddonsAPI.On("ListAddons", mock.Anything, &eks.ListAddonsInput{
			ClusterName: aws.String(clusterName),
		}, mock.Anything).Return(&eks.ListAddonsOutput{
			Addons: []string{"vpc-cni"},
		}, nil)
		eksAddonsAPI.On("DescribeAddon", mock.Anything, &eks.DescribeAddonInput{
			AddonName:   aws.String("vpc-cni"),
			ClusterName: aws.String(clusterName),
		}).Return(&eks.DescribeAddonOutput{
			Addon: &ekstypes.Addon{
				AddonName:             aws.String("vpc-cni"),
				AddonVersion:          aws.String("v1"),
				ServiceAccountRoleArn: aws.String("arn:aws:iam::000:role/role-1"),
				ClusterName:           aws.String(clusterName),
				ConfigurationValues:   aws.String("{}"),
			},
		}, nil)
	}

	mockGetRole := func(iamRoleGetter *mocksv2.IAM, includeServiceAccountSubject bool) {
		strEquals := map[string]string{
			"oidc.eks.us-west-2.amazonaws.com/id/00:aud": "sts.amazonaws.com",
		}
		if includeServiceAccountSubject {
			strEquals["oidc.eks.us-west-2.amazonaws.com/id/00:sub"] = "system:serviceaccount:kube-system:aws-node"
		}

		val, err := json.Marshal(strEquals)
		Expect(err).NotTo(HaveOccurred())
		iamRoleGetter.On("GetRole", mock.Anything, &iam.GetRoleInput{
			RoleName: aws.String("role-1"),
		}).Return(&iam.GetRoleOutput{
			Role: &iamtypes.Role{
				AssumeRolePolicyDocument: aws.String(fmt.Sprintf(`
{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Effect": "Allow",
			"Principal": {
				"Federated": "arn:aws:iam::00:oidc-provider/oidc.eks.eu-north-1.amazonaws.com/id/00"
			},
			"Action": "sts:AssumeRoleWithWebIdentity",
			"Condition": {
				"StringEquals": %s
			}
		}
	]
}`, val)),
			},
		}, nil)
	}

	DescribeTable("update pod identity associations", func(e migrateEntry) {
		var (
			provider       = mockprovider.NewMockProvider()
			stackDescriber mocks.StackDeleter
			roleMigrator   mocks.RoleMigrator
		)

		addonMocks := addonMocks{
			eksAddonsAPI:   provider.MockEKS(),
			iamRoleGetter:  provider.MockIAM(),
			stackDescriber: &stackDescriber,
			roleMigrator:   &roleMigrator,
		}
		e.mockCalls(addonMocks)
		addonMapper, err := podidentityassociation.CreateAddonServiceAccountRoleMapper(context.Background(), clusterName, provider.MockEKS())
		Expect(err).NotTo(HaveOccurred())
		addonMigrator := &podidentityassociation.AddonMigrator{
			ClusterName:                   clusterName,
			AddonServiceAccountRoleMapper: addonMapper,
			IAMRoleGetter:                 provider.MockIAM(),
			StackDescriber:                &stackDescriber,
			EKSAddonsAPI:                  provider.MockEKS(),
			RoleMigrator:                  &roleMigrator,
		}
		taskTree, err := addonMigrator.Migrate(context.Background())
		if e.expectedErr != "" {
			Expect(err).To(MatchError(e.expectedErr))
		} else {
			Expect(err).NotTo(HaveOccurred())
		}
		if e.expectedTasks {
			Expect(taskTree.Tasks).NotTo(BeEmpty(), "expected tasks to be non-empty")
		} else {
			Expect(taskTree.Tasks).To(BeEmpty(), "expected tasks to be empty")
		}
		for _, err := range taskTree.DoAllSync() {
			Expect(err).NotTo(HaveOccurred())
		}
		for _, asserter := range []interface {
			AssertExpectations(t mock.TestingT) bool
		}{
			addonMocks.eksAddonsAPI,
			addonMocks.iamRoleGetter,
			addonMocks.stackDescriber,
			addonMocks.roleMigrator,
		} {
			asserter.AssertExpectations(GinkgoT())
		}
	},
		Entry("migrating an addon with unowned IAM resources", migrateEntry{
			mockCalls: func(m addonMocks) {
				mockAddonCalls(m.eksAddonsAPI)
				m.eksAddonsAPI.On("DescribeAddonConfiguration", mock.Anything, &eks.DescribeAddonConfigurationInput{
					AddonName:    aws.String("vpc-cni"),
					AddonVersion: aws.String("v1"),
				}).Return(&eks.DescribeAddonConfigurationOutput{
					PodIdentityConfiguration: []ekstypes.AddonPodIdentityConfiguration{
						{
							ServiceAccount: aws.String("aws-node"),
						},
					},
				}, nil)
				mockGetRole(m.iamRoleGetter, true)

				m.stackDescriber.On("DescribeStack", mock.Anything, &manager.Stack{
					StackName: aws.String(fmt.Sprintf("eksctl-%s-addon-%s", clusterName, "vpc-cni")),
				}).Return(nil, &smithy.OperationError{
					Err: errors.New("ValidationError"),
				})
				m.roleMigrator.On("UpdateTrustPolicyForUnownedRoleTask", mock.Anything, "role-1", true).Return(&tasks.GenericTask{
					Description: `update trust policy for unowned role "role-1"`,
					Doer: func() error {
						return nil
					},
				})

				m.eksAddonsAPI.On("UpdateAddon", mock.Anything, &eks.UpdateAddonInput{
					AddonName:           aws.String("vpc-cni"),
					AddonVersion:        aws.String("v1"),
					ClusterName:         aws.String(clusterName),
					ConfigurationValues: aws.String("{}"),
					PodIdentityAssociations: []ekstypes.AddonPodIdentityAssociations{
						{
							RoleArn:        aws.String("arn:aws:iam::000:role/role-1"),
							ServiceAccount: aws.String("aws-node"),
						},
					},
				}).Return(&eks.UpdateAddonOutput{}, nil)
			},

			expectedTasks: true,
		}),

		Entry("migrating an addon with owned IAM resources", migrateEntry{
			mockCalls: func(m addonMocks) {
				mockAddonCalls(m.eksAddonsAPI)
				m.eksAddonsAPI.On("DescribeAddonConfiguration", mock.Anything, &eks.DescribeAddonConfigurationInput{
					AddonName:    aws.String("vpc-cni"),
					AddonVersion: aws.String("v1"),
				}).Return(&eks.DescribeAddonConfigurationOutput{
					PodIdentityConfiguration: []ekstypes.AddonPodIdentityConfiguration{
						{
							ServiceAccount: aws.String("aws-node"),
						},
					},
				}, nil)
				mockGetRole(m.iamRoleGetter, true)

				m.stackDescriber.On("DescribeStack", mock.Anything, &manager.Stack{
					StackName: aws.String(fmt.Sprintf("eksctl-%s-addon-%s", clusterName, "vpc-cni")),
				}).Return(&manager.Stack{
					StackName: aws.String(fmt.Sprintf("eksctl-%s-addon-%s", clusterName, "vpc-cni")),
					Tags: []cfntypes.Tag{
						{
							Key:   aws.String("alpha.eksctl.io/addon-name"),
							Value: aws.String("vpc-cni"),
						},
					},
					Capabilities: []cfntypes.Capability{cfntypes.CapabilityCapabilityIam},
				}, nil)
				m.roleMigrator.On("UpdateTrustPolicyForOwnedRoleTask", mock.Anything, "role-1", "aws-node", podidentityassociation.IRSAv1StackSummary{
					Name: fmt.Sprintf("eksctl-%s-addon-%s", clusterName, "vpc-cni"),
					Tags: map[string]string{
						"alpha.eksctl.io/addon-name": "vpc-cni",
					},
					Capabilities: []string{string(cfntypes.CapabilityCapabilityIam)},
				}, true).Return(&tasks.GenericTask{
					Description: `update trust policy for owned role "role-1"`,
					Doer: func() error {
						return nil
					},
				})

				m.eksAddonsAPI.On("UpdateAddon", mock.Anything, &eks.UpdateAddonInput{
					AddonName:           aws.String("vpc-cni"),
					AddonVersion:        aws.String("v1"),
					ClusterName:         aws.String(clusterName),
					ConfigurationValues: aws.String("{}"),
					PodIdentityAssociations: []ekstypes.AddonPodIdentityAssociations{
						{
							RoleArn:        aws.String("arn:aws:iam::000:role/role-1"),
							ServiceAccount: aws.String("aws-node"),
						},
					},
				}).Return(&eks.UpdateAddonOutput{}, nil)
			},

			expectedTasks: true,
		}),

		Entry("addon that does not support pod identity is skipped", migrateEntry{
			mockCalls: func(m addonMocks) {
				mockAddonCalls(m.eksAddonsAPI)
				m.eksAddonsAPI.On("DescribeAddonConfiguration", mock.Anything, &eks.DescribeAddonConfigurationInput{
					AddonName:    aws.String("vpc-cni"),
					AddonVersion: aws.String("v1"),
				}).Return(&eks.DescribeAddonConfigurationOutput{}, nil)
			},
		}),

		Entry("addon without service account in IAM role policy uses service account from addon configuration", migrateEntry{
			mockCalls: func(m addonMocks) {
				mockAddonCalls(m.eksAddonsAPI)
				m.eksAddonsAPI.On("DescribeAddonConfiguration", mock.Anything, &eks.DescribeAddonConfigurationInput{
					AddonName:    aws.String("vpc-cni"),
					AddonVersion: aws.String("v1"),
				}).Return(&eks.DescribeAddonConfigurationOutput{
					PodIdentityConfiguration: []ekstypes.AddonPodIdentityConfiguration{
						{
							ServiceAccount: aws.String("aws-node"),
						},
					},
				}, nil)
				mockGetRole(m.iamRoleGetter, false)

				m.stackDescriber.On("DescribeStack", mock.Anything, &manager.Stack{
					StackName: aws.String(fmt.Sprintf("eksctl-%s-addon-%s", clusterName, "vpc-cni")),
				}).Return(&manager.Stack{
					StackName: aws.String(fmt.Sprintf("eksctl-%s-addon-%s", clusterName, "vpc-cni")),
					Tags: []cfntypes.Tag{
						{
							Key:   aws.String("alpha.eksctl.io/addon-name"),
							Value: aws.String("vpc-cni"),
						},
					},
					Capabilities: []cfntypes.Capability{cfntypes.CapabilityCapabilityIam},
				}, nil)
				m.roleMigrator.On("UpdateTrustPolicyForOwnedRoleTask", mock.Anything, "role-1", "aws-node", podidentityassociation.IRSAv1StackSummary{
					Name: fmt.Sprintf("eksctl-%s-addon-%s", clusterName, "vpc-cni"),
					Tags: map[string]string{
						"alpha.eksctl.io/addon-name": "vpc-cni",
					},
					Capabilities: []string{string(cfntypes.CapabilityCapabilityIam)},
				}, true).Return(&tasks.GenericTask{
					Description: `update trust policy for owned role "role-1"`,
					Doer: func() error {
						return nil
					},
				})

				m.eksAddonsAPI.On("UpdateAddon", mock.Anything, &eks.UpdateAddonInput{
					AddonName:           aws.String("vpc-cni"),
					AddonVersion:        aws.String("v1"),
					ClusterName:         aws.String(clusterName),
					ConfigurationValues: aws.String("{}"),
					PodIdentityAssociations: []ekstypes.AddonPodIdentityAssociations{
						{
							RoleArn:        aws.String("arn:aws:iam::000:role/role-1"),
							ServiceAccount: aws.String("aws-node"),
						},
					},
				}).Return(&eks.UpdateAddonOutput{}, nil)
			},

			expectedTasks: true,
		}),

		Entry("addon without service account in IAM role policy and multiple service accounts in addon configuration is skipped", migrateEntry{
			mockCalls: func(m addonMocks) {
				mockAddonCalls(m.eksAddonsAPI)
				m.eksAddonsAPI.On("DescribeAddonConfiguration", mock.Anything, &eks.DescribeAddonConfigurationInput{
					AddonName:    aws.String("vpc-cni"),
					AddonVersion: aws.String("v1"),
				}).Return(&eks.DescribeAddonConfigurationOutput{
					PodIdentityConfiguration: []ekstypes.AddonPodIdentityConfiguration{
						{
							ServiceAccount: aws.String("aws-node"),
						},
						{
							ServiceAccount: aws.String("ipam"),
						},
					},
				}, nil)
				mockGetRole(m.iamRoleGetter, false)
			},

			expectedTasks: false,
		}),

		Entry("addon with Statement.Condition missing in IAM role is skipped", migrateEntry{
			mockCalls: func(m addonMocks) {
				mockAddonCalls(m.eksAddonsAPI)
				m.eksAddonsAPI.On("DescribeAddonConfiguration", mock.Anything, &eks.DescribeAddonConfigurationInput{
					AddonName:    aws.String("vpc-cni"),
					AddonVersion: aws.String("v1"),
				}).Return(&eks.DescribeAddonConfigurationOutput{
					PodIdentityConfiguration: []ekstypes.AddonPodIdentityConfiguration{
						{
							ServiceAccount: aws.String("aws-node"),
						},
						{
							ServiceAccount: aws.String("ipam"),
						},
					},
				}, nil)

				m.iamRoleGetter.On("GetRole", mock.Anything, &iam.GetRoleInput{
					RoleName: aws.String("role-1"),
				}).Return(&iam.GetRoleOutput{
					Role: &iamtypes.Role{
						AssumeRolePolicyDocument: aws.String(`
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "beta.pods.eks.aws.internal"
      },
      "Action": [
        "sts:AssumeRole",
        "sts:TagSession"
      ]
    }
  ]
}
`),
					},
				}, nil)
			},

			expectedTasks: false,
		}),
	)
})
