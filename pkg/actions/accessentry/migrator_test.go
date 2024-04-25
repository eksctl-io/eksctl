package accessentry_test

import (
	"bytes"
	"context"
	"encoding/base32"
	"fmt"
	"os"
	"strings"

	"github.com/kris-nova/logger"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"gopkg.in/yaml.v2"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/aws/aws-sdk-go-v2/aws"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	awsiam "github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"

	"github.com/weaveworks/eksctl/pkg/actions/accessentry"
	"github.com/weaveworks/eksctl/pkg/actions/accessentry/fakes"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/iam"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

type migrateToAccessEntryEntry struct {
	curAuthMode                ekstypes.AuthenticationMode
	tgAuthMode                 ekstypes.AuthenticationMode
	mockIAM                    func(provider *mockprovider.MockProvider)
	mockK8s                    func(clientSet *fake.Clientset)
	mockAccessEntries          func(getter *fakes.FakeGetterInterface)
	validateCustomLoggerOutput func(output string)
	options                    accessentry.MigrationOptions
	expectedErr                string
}

var _ = Describe("Migrate Access Entry", func() {

	var (
		migrator      *accessentry.Migrator
		mockProvider  *mockprovider.MockProvider
		fakeClientset *fake.Clientset
		fakeAECreator accessentry.CreatorInterface
		fakeAEGetter  accessentry.GetterInterface
		clusterName   = "test-cluster"
		genericErr    = fmt.Errorf("ERR")
	)

	DescribeTable("Migrate", func(ae migrateToAccessEntryEntry) {
		var s fakes.FakeStackCreator
		s.CreateStackStub = func(ctx context.Context, stackName string, r builder.ResourceSetReader, tags map[string]string, parameters map[string]string, errorCh chan error) error {
			defer close(errorCh)
			prefix := fmt.Sprintf("eksctl-%s-accessentry-", clusterName)
			idx := strings.Index(stackName, prefix)
			if idx < 0 {
				return fmt.Errorf("expected stack name to have prefix %q", prefix)
			}
			suffix := stackName[idx+len(prefix):]
			_, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(suffix)
			if err != nil {
				return fmt.Errorf("expected stack name to have a base32-encoded suffix: %w", err)
			}
			return nil
		}

		mockProvider = mockprovider.NewMockProvider()
		if ae.mockIAM != nil {
			ae.mockIAM(mockProvider)
		}

		fakeClientset = fake.NewSimpleClientset()
		if ae.mockK8s != nil {
			ae.mockK8s(fakeClientset)
		}

		fakeAECreator = &accessentry.Creator{ClusterName: clusterName}
		fakeAEGetter = &fakes.FakeGetterInterface{}
		if ae.mockAccessEntries != nil {
			ae.mockAccessEntries(fakeAEGetter.(*fakes.FakeGetterInterface))
		}

		output := &bytes.Buffer{}
		if ae.validateCustomLoggerOutput != nil {
			defer func() {
				logger.Writer = os.Stdout
			}()
			logger.Writer = output
		}

		migrator = accessentry.NewMigrator(
			clusterName,
			mockProvider.MockEKS(),
			mockProvider.MockIAM(),
			fakeClientset,
			fakeAECreator,
			fakeAEGetter,
			ae.curAuthMode,
			ae.tgAuthMode,
		)

		err := migrator.MigrateToAccessEntry(context.Background(), ae.options)

		if ae.expectedErr != "" {
			Expect(err).To(MatchError(ContainSubstring(ae.expectedErr)))
			return
		}

		Expect(err).ToNot(HaveOccurred())

		if ae.validateCustomLoggerOutput != nil {
			ae.validateCustomLoggerOutput(output.String())
		}
	},
		Entry("[Validation Error] target authentication mode is CONFIG_MAP", migrateToAccessEntryEntry{
			tgAuthMode:  ekstypes.AuthenticationModeConfigMap,
			expectedErr: "target authentication mode is invalid",
		}),

		Entry("[Validation Error] current authentication mode is API", migrateToAccessEntryEntry{
			curAuthMode: ekstypes.AuthenticationModeApi,
			tgAuthMode:  ekstypes.AuthenticationModeApi,
			validateCustomLoggerOutput: func(output string) {
				Expect(output).To(ContainSubstring(fmt.Sprintf("cluster authentication mode is already %s; there is no need to migrate to access entries", ekstypes.AuthenticationModeApi)))
			},
		}),

		Entry("[API Error] getting access entries fails", migrateToAccessEntryEntry{
			curAuthMode: ekstypes.AuthenticationModeApiAndConfigMap,
			tgAuthMode:  ekstypes.AuthenticationModeApi,
			mockAccessEntries: func(getter *fakes.FakeGetterInterface) {
				getter.GetReturns(nil, genericErr)
			},
			expectedErr: "fetching existing access entries",
		}),

		Entry("[API Error] getting role fails", migrateToAccessEntryEntry{
			curAuthMode: ekstypes.AuthenticationModeApiAndConfigMap,
			tgAuthMode:  ekstypes.AuthenticationModeApi,
			mockAccessEntries: func(getter *fakes.FakeGetterInterface) {
				getter.GetReturns([]accessentry.Summary{}, nil)
			},
			mockIAM: func(provider *mockprovider.MockProvider) {
				provider.MockIAM().
					On("GetRole", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(2))
						Expect(args[1]).To(BeAssignableToTypeOf(&awsiam.GetRoleInput{}))
					}).
					Return(nil, &iamtypes.NoSuchEntityException{}).
					Once()
			},
			mockK8s: func(clientSet *fake.Clientset) {
				roles := []iam.RoleIdentity{{RoleARN: "arn:aws:iam::111122223333:role/test"}}
				rolesBytes, err := yaml.Marshal(roles)
				Expect(err).NotTo(HaveOccurred())

				_, err = clientSet.CoreV1().ConfigMaps(authconfigmap.ObjectNamespace).Create(context.Background(), &v1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name: authconfigmap.ObjectName,
					},
					Data: map[string]string{
						"mapRoles": string(rolesBytes),
					},
				}, metav1.CreateOptions{})
				Expect(err).NotTo(HaveOccurred())
			},
			expectedErr: fmt.Sprintf("role %q does not exists", "test"),
		}),

		Entry("[API Error] getting user fails", migrateToAccessEntryEntry{
			curAuthMode: ekstypes.AuthenticationModeApiAndConfigMap,
			tgAuthMode:  ekstypes.AuthenticationModeApi,
			mockAccessEntries: func(getter *fakes.FakeGetterInterface) {
				getter.GetReturns([]accessentry.Summary{}, nil)
			},
			mockIAM: func(provider *mockprovider.MockProvider) {
				provider.MockIAM().
					On("GetUser", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(2))
						Expect(args[1]).To(BeAssignableToTypeOf(&awsiam.GetUserInput{}))
					}).
					Return(nil, &iamtypes.NoSuchEntityException{}).
					Once()
			},
			mockK8s: func(clientSet *fake.Clientset) {
				users := []iam.UserIdentity{{UserARN: "arn:aws:iam::111122223333:user/test"}}
				usersBytes, err := yaml.Marshal(users)
				Expect(err).NotTo(HaveOccurred())

				_, err = clientSet.CoreV1().ConfigMaps(authconfigmap.ObjectNamespace).Create(context.Background(), &v1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name: authconfigmap.ObjectName,
					},
					Data: map[string]string{
						"mapUsers": string(usersBytes),
					},
				}, metav1.CreateOptions{})
				Expect(err).NotTo(HaveOccurred())
			},
			expectedErr: fmt.Sprintf("user %q does not exists", "test"),
		}),

		Entry("[TaskTree] should not switch to API mode if service-linked role iamidentitymapping is found", migrateToAccessEntryEntry{
			curAuthMode: ekstypes.AuthenticationModeApiAndConfigMap,
			tgAuthMode:  ekstypes.AuthenticationModeApi,
			mockAccessEntries: func(getter *fakes.FakeGetterInterface) {
				getter.GetReturns([]accessentry.Summary{}, nil)
			},
			mockIAM: func(provider *mockprovider.MockProvider) {
				provider.MockIAM().
					On("GetRole", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(2))
						Expect(args[1]).To(BeAssignableToTypeOf(&awsiam.GetRoleInput{}))
					}).
					Return(&awsiam.GetRoleOutput{
						Role: &iamtypes.Role{
							Arn: aws.String("arn:aws:iam::111122223333:role/aws-service-role/test"),
						},
					}, nil).
					Once()
			},
			mockK8s: func(clientSet *fake.Clientset) {
				roles := []iam.RoleIdentity{{RoleARN: "arn:aws:iam::111122223333:role/aws-service-role/test"}}
				rolesBytes, err := yaml.Marshal(roles)
				Expect(err).NotTo(HaveOccurred())

				_, err = clientSet.CoreV1().ConfigMaps(authconfigmap.ObjectNamespace).Create(context.Background(), &v1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name: authconfigmap.ObjectName,
					},
					Data: map[string]string{
						"mapRoles": string(rolesBytes),
					},
				}, metav1.CreateOptions{})
				Expect(err).NotTo(HaveOccurred())
			},
			validateCustomLoggerOutput: func(output string) {
				Expect(output).To(ContainSubstring("found service-linked role iamidentitymapping"))
				Expect(output).NotTo(ContainSubstring("update authentication mode from API_AND_CONFIG_MAP to API"))
				Expect(output).NotTo(ContainSubstring("delete aws-auth configMap when authentication mode is API"))
			},
		}),

		Entry("[TaskTree] should not switch to API mode if iamidentitymapping with non-master `system:` group is found", migrateToAccessEntryEntry{
			curAuthMode: ekstypes.AuthenticationModeApiAndConfigMap,
			tgAuthMode:  ekstypes.AuthenticationModeApi,
			mockAccessEntries: func(getter *fakes.FakeGetterInterface) {
				getter.GetReturns([]accessentry.Summary{}, nil)
			},
			mockIAM: func(provider *mockprovider.MockProvider) {
				provider.MockIAM().
					On("GetRole", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(2))
						Expect(args[1]).To(BeAssignableToTypeOf(&awsiam.GetRoleInput{}))
					}).
					Return(&awsiam.GetRoleOutput{
						Role: &iamtypes.Role{
							Arn: aws.String("arn:aws:iam::111122223333:role/test"),
						},
					}, nil).
					Once()
			},
			mockK8s: func(clientSet *fake.Clientset) {
				roles := []iam.RoleIdentity{
					{
						RoleARN: "arn:aws:iam::111122223333:role/test",
						KubernetesIdentity: iam.KubernetesIdentity{
							KubernetesGroups: []string{"system:tests"},
						},
					},
				}
				rolesBytes, err := yaml.Marshal(roles)
				Expect(err).NotTo(HaveOccurred())

				_, err = clientSet.CoreV1().ConfigMaps(authconfigmap.ObjectNamespace).Create(context.Background(), &v1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name: authconfigmap.ObjectName,
					},
					Data: map[string]string{
						"mapRoles": string(rolesBytes),
					},
				}, metav1.CreateOptions{})
				Expect(err).NotTo(HaveOccurred())
			},
			validateCustomLoggerOutput: func(output string) {
				Expect(output).To(ContainSubstring("at least one group name associated with %q starts with \"system:\"", "arn:aws:iam::111122223333:role/test"))
				Expect(output).NotTo(ContainSubstring("update authentication mode from API_AND_CONFIG_MAP to API"))
				Expect(output).NotTo(ContainSubstring("delete aws-auth configMap when authentication mode is API"))
			},
		}),

		Entry("[TaskTree] should not switch to API mode if account iamidentitymapping is found", migrateToAccessEntryEntry{
			curAuthMode: ekstypes.AuthenticationModeApiAndConfigMap,
			tgAuthMode:  ekstypes.AuthenticationModeApi,
			mockAccessEntries: func(getter *fakes.FakeGetterInterface) {
				getter.GetReturns([]accessentry.Summary{}, nil)
			},
			mockK8s: func(clientSet *fake.Clientset) {
				accounts := []string{"test-account"}
				accountsBytes, err := yaml.Marshal(accounts)
				Expect(err).NotTo(HaveOccurred())

				_, err = clientSet.CoreV1().ConfigMaps(authconfigmap.ObjectNamespace).Create(context.Background(), &v1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name: authconfigmap.ObjectName,
					},
					Data: map[string]string{
						"mapAccounts": string(accountsBytes),
					},
				}, metav1.CreateOptions{})
				Expect(err).NotTo(HaveOccurred())
			},
			validateCustomLoggerOutput: func(output string) {
				Expect(output).To(ContainSubstring("found account iamidentitymapping"))
				Expect(output).NotTo(ContainSubstring("update authentication mode from API_AND_CONFIG_MAP to API"))
				Expect(output).NotTo(ContainSubstring("delete aws-auth configMap when authentication mode is API"))
			},
		}),

		Entry("[TaskTree] should contain all expected tasks", migrateToAccessEntryEntry{
			curAuthMode: ekstypes.AuthenticationModeConfigMap,
			tgAuthMode:  ekstypes.AuthenticationModeApi,
			mockAccessEntries: func(getter *fakes.FakeGetterInterface) {
				getter.GetReturns([]accessentry.Summary{
					{
						PrincipalARN: "arn:aws:iam::111122223333:role/eksctl-test-cluster-nodegroup-NodeInstanceRole-1",
					},
				}, nil)
			},
			mockIAM: func(provider *mockprovider.MockProvider) {
				provider.MockIAM().
					On("GetRole", mock.Anything, &awsiam.GetRoleInput{
						RoleName: aws.String("eksctl-test-cluster-nodegroup-NodeInstanceRole-1"),
					}).
					Return(&awsiam.GetRoleOutput{
						Role: &iamtypes.Role{
							Arn: aws.String("arn:aws:iam::111122223333:role/eksctl-test-cluster-nodegroup-NodeInstanceRole-1"),
						},
					}, nil)
				provider.MockIAM().
					On("GetRole", mock.Anything, &awsiam.GetRoleInput{
						RoleName: aws.String("eksctl-test-cluster-nodegroup-NodeInstanceRole-2"),
					}).
					Return(&awsiam.GetRoleOutput{
						Role: &iamtypes.Role{
							Arn: aws.String("arn:aws:iam::111122223333:role/eksctl-test-cluster-nodegroup-NodeInstanceRole-2"),
						},
					}, nil)
				provider.MockIAM().
					On("GetUser", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(2))
						Expect(args[1]).To(BeAssignableToTypeOf(&awsiam.GetUserInput{}))
					}).
					Return(&awsiam.GetUserOutput{
						User: &iamtypes.User{
							Arn: aws.String("arn:aws:iam::111122223333:user/admin"),
						},
					}, nil).
					Once()
			},
			mockK8s: func(clientSet *fake.Clientset) {
				roles := []iam.RoleIdentity{
					{
						RoleARN: "arn:aws:iam::111122223333:role/eksctl-test-cluster-nodegroup-NodeInstanceRole-1",
						KubernetesIdentity: iam.KubernetesIdentity{
							KubernetesUsername: "system:node:{{EC2PrivateDNSName}}",
							KubernetesGroups:   []string{"system:nodes", "system:bootstrappers"},
						},
					},
					{
						RoleARN: "arn:aws:iam::111122223333:role/eksctl-test-cluster-nodegroup-NodeInstanceRole-2",
						KubernetesIdentity: iam.KubernetesIdentity{
							KubernetesUsername: "system:node:{{EC2PrivateDNSName}}",
							KubernetesGroups:   []string{"system:nodes", "system:bootstrappers"},
						},
					},
				}
				users := []iam.UserIdentity{
					{
						UserARN: "arn:aws:iam::111122223333:user/admin",
						KubernetesIdentity: iam.KubernetesIdentity{
							KubernetesUsername: "admin",
							KubernetesGroups:   []string{"system:masters"},
						},
					},
				}
				rolesBytes, err := yaml.Marshal(roles)
				Expect(err).NotTo(HaveOccurred())
				usersBytes, err := yaml.Marshal(users)
				Expect(err).NotTo(HaveOccurred())

				_, err = clientSet.CoreV1().ConfigMaps(authconfigmap.ObjectNamespace).Create(context.Background(), &v1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name: authconfigmap.ObjectName,
					},
					Data: map[string]string{
						"mapRoles": string(rolesBytes),
						"mapUsers": string(usersBytes),
					},
				}, metav1.CreateOptions{})
				Expect(err).NotTo(HaveOccurred())
			},
			validateCustomLoggerOutput: func(output string) {
				Expect(output).To(ContainSubstring("create access entry for principal ARN arn:aws:iam::111122223333:user/admin"))
				Expect(output).To(ContainSubstring("create access entry for principal ARN arn:aws:iam::111122223333:role/eksctl-test-cluster-nodegroup-NodeInstanceRole-2"))
				Expect(output).To(ContainSubstring("update authentication mode from CONFIG_MAP to API_AND_CONFIG_MAP"))
				Expect(output).To(ContainSubstring("update authentication mode from API_AND_CONFIG_MAP to API"))
				// filter out existing access entries
				Expect(output).NotTo(ContainSubstring("create access entry for principal ARN arn:aws:iam::111122223333:role/eksctl-test-cluster-nodegroup-NodeInstanceRole-1"))
			},
		}),
	)
})
