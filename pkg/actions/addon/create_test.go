package addon_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"k8s.io/client-go/kubernetes"

	"github.com/aws/aws-sdk-go-v2/aws"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/kris-nova/logger"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/weaveworks/eksctl/pkg/actions/addon"
	"github.com/weaveworks/eksctl/pkg/actions/addon/fakes"
	addonmocks "github.com/weaveworks/eksctl/pkg/actions/addon/mocks"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/eks/mocksv2"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

type createAddonEntry struct {
	addon       api.Addon
	withOIDC    bool
	waitTimeout time.Duration

	mockClusterConfig func(clusterConfig *api.ClusterConfig)
	mockIAM           func(mockIAMRoleCreator *addonmocks.IAMRoleCreator)
	mockEKS           func(provider *mockprovider.MockProvider)
	mockCFN           func(stackManager *fakes.FakeStackManager)
	mockK8s           bool

	validateCreateAddonInput   func(input *awseks.CreateAddonInput)
	validateCustomLoggerOutput func(output string)
	validateCFNCalls           func(stackManager *fakes.FakeStackManager)
	expectedErr                string
}

var _ = Describe("Create", func() {
	var (
		manager            *addon.Manager
		oidc               *iamoidc.OpenIDConnectManager
		fakeStackManager   *fakes.FakeStackManager
		mockProvider       *mockprovider.MockProvider
		createAddonInput   *awseks.CreateAddonInput
		mockIAMRoleCreator *addonmocks.IAMRoleCreator
		fakeRawClient      *testutils.FakeRawClient
		err                error
		genericErr         = fmt.Errorf("ERR")
	)

	mockDescribeAddon := func(mockEKS *mocksv2.EKS, err error) {
		if err == nil {
			err = &ekstypes.ResourceNotFoundException{
				Message: aws.String(genericErr.Error()),
			}
		}
		mockEKS.
			On("DescribeAddon", mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(2))
				Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeAddonInput{}))
			}).
			Return(nil, err).
			Once()
	}

	mockDescribeAddonVersions := func(mockEKS *mocksv2.EKS, err error) {
		var output *awseks.DescribeAddonVersionsOutput
		if err == nil {
			output = &awseks.DescribeAddonVersionsOutput{
				Addons: []ekstypes.AddonInfo{
					{
						AddonVersions: []ekstypes.AddonVersionInfo{
							{
								AddonVersion:           aws.String("v1.0.0-eksbuild.1"),
								RequiresIamPermissions: false,
							},
							{
								AddonVersion:           aws.String("v1.7.5-eksbuild.1"),
								RequiresIamPermissions: true,
								Compatibilities: []ekstypes.Compatibility{
									{
										ClusterVersion: aws.String(api.DefaultVersion),
										DefaultVersion: true,
									},
								},
							},
							{
								AddonVersion: aws.String("v1.7.5-eksbuild.2"),
							},
							{
								AddonVersion: aws.String("v1.7.7-eksbuild.2"),
							},
							{
								AddonVersion: aws.String("v1.7.6"),
							},
						},
					},
				},
			}
		}
		mockEKS.
			On("DescribeAddonVersions", mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(2))
				Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeAddonVersionsInput{}))
			}).
			Return(output, err).
			Once()
	}

	mockDescribeAddonConfiguration := func(mockEKS *mocksv2.EKS, serviceAccountNames []string, err error) {
		podIDConfig := []ekstypes.AddonPodIdentityConfiguration{}
		for _, sa := range serviceAccountNames {
			podIDConfig = append(podIDConfig, ekstypes.AddonPodIdentityConfiguration{
				ServiceAccount:             &sa,
				RecommendedManagedPolicies: []string{"arn:aws:iam::111122223333:policy/" + sa},
			})
		}
		mockEKS.
			On("DescribeAddonConfiguration", mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(2))
				Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeAddonConfigurationInput{}))
			}).
			Return(&awseks.DescribeAddonConfigurationOutput{
				PodIdentityConfiguration: podIDConfig,
			}, err).
			Once()
	}

	mockCreateAddon := func(mockEKS *mocksv2.EKS, err error) {
		mockEKS.
			On("CreateAddon", mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(2))
				Expect(args[1]).To(BeAssignableToTypeOf(&awseks.CreateAddonInput{}))
				createAddonInput = args[1].(*awseks.CreateAddonInput)
			}).
			Return(&awseks.CreateAddonOutput{
				Addon: &ekstypes.Addon{},
			}, err).
			Once()
	}

	DescribeTable("Create addon", func(e createAddonEntry) {
		if e.addon.Name == "" {
			e.addon.Name = "my-addon"
		}

		clusterConfig := api.NewClusterConfig()
		if e.mockClusterConfig != nil {
			e.mockClusterConfig(clusterConfig)
		}

		mockIAMRoleCreator = new(addonmocks.IAMRoleCreator)
		if e.mockIAM != nil {
			e.mockIAM(mockIAMRoleCreator)
		}

		fakeStackManager = new(fakes.FakeStackManager)
		if e.mockCFN != nil {
			e.mockCFN(fakeStackManager)
		}

		mockProvider = mockprovider.NewMockProvider()
		if e.mockEKS != nil {
			e.mockEKS(mockProvider)
		}

		fakeRawClient = testutils.NewFakeRawClient()
		if e.mockK8s {
			fakeRawClient.AssumeObjectsMissing = true
			sampleAddons := testutils.LoadSamples("testdata/aws-node.json")
			for _, item := range sampleAddons {
				rc, err := fakeRawClient.NewRawResource(item)
				Expect(err).NotTo(HaveOccurred())
				_, err = rc.CreateOrReplace(false)
				Expect(err).NotTo(HaveOccurred())
			}

			ct := fakeRawClient.Collection
			Expect(ct.Updated()).To(BeEmpty())
			Expect(ct.Created()).NotTo(BeEmpty())
			Expect(ct.CreatedItems()).To(HaveLen(10))
		}

		output := &bytes.Buffer{}
		if e.validateCustomLoggerOutput != nil {
			defer func() {
				logger.Writer = os.Stdout
			}()
			logger.Writer = output
		}

		oidc, err = iamoidc.NewOpenIDConnectManager(nil, "111122223333", "https://oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E", "aws", nil)
		Expect(err).NotTo(HaveOccurred())

		manager, err = addon.New(clusterConfig, mockProvider.EKS(), fakeStackManager, e.withOIDC, oidc, func() (kubernetes.Interface, error) {
			return fakeRawClient.ClientSet(), nil
		})
		Expect(err).NotTo(HaveOccurred())

		err = manager.Create(context.Background(), &e.addon, mockIAMRoleCreator, e.waitTimeout)
		if e.expectedErr != "" {
			Expect(err).To(MatchError(ContainSubstring(e.expectedErr)))
			return
		}
		Expect(err).ToNot(HaveOccurred())

		if e.validateCreateAddonInput != nil {
			e.validateCreateAddonInput(createAddonInput)
		}

		if e.validateCustomLoggerOutput != nil {
			e.validateCustomLoggerOutput(output.String())
		}

		if e.validateCFNCalls != nil {
			e.validateCFNCalls(fakeStackManager)
		}
	},
		Entry("[API Error] fails to describe addon", createAddonEntry{
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(provider.MockEKS(), genericErr)
			},
			expectedErr: "failed to describe addon",
		}),

		Entry("[API Error] fails to describe addon versions", createAddonEntry{
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(provider.MockEKS(), nil)
				mockDescribeAddonVersions(provider.MockEKS(), genericErr)
			},
			expectedErr: "failed to describe addon versions",
		}),

		Entry("[API Error] fails to describe addon configuration", createAddonEntry{
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(provider.MockEKS(), nil)
				mockDescribeAddonVersions(provider.MockEKS(), nil)
				mockDescribeAddonConfiguration(provider.MockEKS(), []string{}, genericErr)
			},
			expectedErr: "failed to describe configuration for \"my-addon\" addon",
		}),

		Entry("[API Error] fails to create addon", createAddonEntry{
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(provider.MockEKS(), nil)
				mockDescribeAddonVersions(provider.MockEKS(), nil)
				mockDescribeAddonConfiguration(provider.MockEKS(), []string{}, nil)
				mockCreateAddon(provider.MockEKS(), genericErr)
			},
			expectedErr: "failed to create \"my-addon\" addon",
		}),

		Entry("[API Error] fails to create IAM role for podID", createAddonEntry{
			addon: api.Addon{
				PodIdentityAssociations: &[]api.PodIdentityAssociation{
					{
						ServiceAccountName:   "sa1",
						PermissionPolicyARNs: []string{"arn:aws:iam::111122223333:policy/sa1"},
					},
				},
			},
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(mockProvider.MockEKS(), nil)
				mockDescribeAddonVersions(provider.MockEKS(), nil)
				mockDescribeAddonConfiguration(mockProvider.MockEKS(), []string{"sa1"}, nil)
			},
			mockIAM: func(mockIAMRoleCreator *addonmocks.IAMRoleCreator) {
				mockIAMRoleCreator.
					On("Create", mock.Anything, mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(3))
						Expect(args[1]).To(BeAssignableToTypeOf(&api.PodIdentityAssociation{}))
						Expect(args[1].(*api.PodIdentityAssociation).ServiceAccountName).To(Equal("sa1"))
						Expect(args[1].(*api.PodIdentityAssociation).PermissionPolicyARNs).To(ConsistOf("arn:aws:iam::111122223333:policy/sa1"))
					}).
					Return("", genericErr).
					Once()
			},
			expectedErr: genericErr.Error(),
		}),

		Entry("[API Error] fails to create IAM role for service account", createAddonEntry{
			addon: api.Addon{
				AttachPolicyARNs: []string{"arn:aws:iam::111122223333:policy/policy-name-1"},
			},
			withOIDC: true,
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(mockProvider.MockEKS(), nil)
				mockDescribeAddonVersions(provider.MockEKS(), nil)
				mockDescribeAddonConfiguration(mockProvider.MockEKS(), []string{}, nil)
			},
			mockCFN: func(stackManager *fakes.FakeStackManager) {
				stackManager.CreateStackStub = func(ctx context.Context, s string, rsr builder.ResourceSetReader, m1, m2 map[string]string, c chan error) error {
					go func() {
						c <- nil
					}()
					Expect(rsr).To(BeAssignableToTypeOf(&builder.IAMRoleResourceSet{}))
					output, err := rsr.(*builder.IAMRoleResourceSet).RenderJSON()
					Expect(err).NotTo(HaveOccurred())
					Expect(string(output)).To(ContainSubstring("arn:aws:iam::111122223333:policy/policy-name-1"))
					return genericErr
				}
				stackManager.CreateStackReturns(genericErr)
			},
			validateCFNCalls: func(stackManager *fakes.FakeStackManager) {
				Expect(stackManager.CreateStackCallCount()).To(Equal(1))
			},
			expectedErr: genericErr.Error(),
		}),

		Entry("[Addon already exists] addon is in CREATE_FAILED state", createAddonEntry{
			addon: api.Addon{
				Version: "1.0.0",
			},
			mockEKS: func(provider *mockprovider.MockProvider) {
				provider.MockEKS().
					On("DescribeAddon", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(2))
						Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeAddonInput{}))
					}).
					Return(&awseks.DescribeAddonOutput{
						Addon: &ekstypes.Addon{
							AddonName: aws.String("my-addon"),
							Status:    ekstypes.AddonStatusCreateFailed,
						},
					}, nil).
					Once()

				mockDescribeAddonVersions(provider.MockEKS(), nil)
				mockCreateAddon(provider.MockEKS(), nil)
			},
			validateCustomLoggerOutput: func(output string) {
				Expect(output).NotTo(ContainSubstring("addon is already present on the cluster, as an EKS managed addon, skipping creation"))
			},
		}),

		Entry("[Addon already exists] addon is NOT in CREATE_FAILED state", createAddonEntry{
			mockEKS: func(provider *mockprovider.MockProvider) {
				provider.MockEKS().
					On("DescribeAddon", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(2))
						Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeAddonInput{}))
					}).
					Return(&awseks.DescribeAddonOutput{
						Addon: &ekstypes.Addon{
							AddonName: aws.String("my-addon"),
							Status:    ekstypes.AddonStatusActive,
						},
					}, nil).
					Once()
			},
			validateCustomLoggerOutput: func(output string) {
				Expect(output).To(ContainSubstring("addon is already present on the cluster, as an EKS managed addon, skipping creation"))
			},
		}),

		Entry("[Resolve version] no version found", createAddonEntry{
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(provider.MockEKS(), nil)
				provider.MockEKS().
					On("DescribeAddonVersions", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(2))
						Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeAddonVersionsInput{}))
					}).
					Return(&awseks.DescribeAddonVersionsOutput{
						Addons: []ekstypes.AddonInfo{{AddonVersions: []ekstypes.AddonVersionInfo{}}},
					}, nil).
					Once()
			},
			expectedErr: "no versions available for \"my-addon\"",
		}),

		Entry("[Resolve version] invalid version found", createAddonEntry{
			addon: api.Addon{
				Version: "latest",
			},
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(provider.MockEKS(), nil)
				provider.MockEKS().
					On("DescribeAddonVersions", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(2))
						Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeAddonVersionsInput{}))
					}).
					Return(&awseks.DescribeAddonVersionsOutput{
						Addons: []ekstypes.AddonInfo{{AddonVersions: []ekstypes.AddonVersionInfo{
							{
								AddonVersion: aws.String("totally not semver"),
							},
						}}},
					}, nil).
					Once()
			},
			expectedErr: "failed to parse version \"totally not semver\":",
		}),

		Entry("[Resolve version] missing", createAddonEntry{
			addon: api.Addon{
				Version: "1.100.0",
			},
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(provider.MockEKS(), nil)
				mockDescribeAddonVersions(provider.MockEKS(), nil)
			},
			expectedErr: "no version(s) found matching \"1.100.0\" for \"my-addon\"",
		}),

		Entry("[Resolve version] latest", createAddonEntry{
			addon: api.Addon{
				Version: "latest",
			},
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(provider.MockEKS(), nil)
				mockDescribeAddonVersions(provider.MockEKS(), nil)
				mockCreateAddon(provider.MockEKS(), nil)
			},
			validateCreateAddonInput: func(input *awseks.CreateAddonInput) {
				Expect(*input.AddonVersion).To(Equal("v1.7.7-eksbuild.2"))
			},
		}),

		Entry("[Resolve version] numeric value", createAddonEntry{
			addon: api.Addon{
				Version: "1.7.7",
			},
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(provider.MockEKS(), nil)
				mockDescribeAddonVersions(provider.MockEKS(), nil)
				mockCreateAddon(provider.MockEKS(), nil)
			},
			validateCreateAddonInput: func(input *awseks.CreateAddonInput) {
				Expect(*input.AddonVersion).To(Equal("v1.7.7-eksbuild.2"))
			},
		}),

		Entry("[Resolve version] alphanumeric value", createAddonEntry{
			addon: api.Addon{
				Version: "v1.7.5",
			},
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(provider.MockEKS(), nil)
				mockDescribeAddonVersions(provider.MockEKS(), nil)
				mockCreateAddon(provider.MockEKS(), nil)
			},
			validateCreateAddonInput: func(input *awseks.CreateAddonInput) {
				Expect(*input.AddonVersion).To(Equal("v1.7.5-eksbuild.2"))
			},
		}),

		Entry("[ResolveConflicts] explicitly set to overwrite", createAddonEntry{
			addon: api.Addon{
				Version:          "1.0.0",
				ResolveConflicts: ekstypes.ResolveConflictsOverwrite,
			},
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(provider.MockEKS(), nil)
				mockDescribeAddonVersions(provider.MockEKS(), nil)
				mockCreateAddon(provider.MockEKS(), nil)
			},
			validateCreateAddonInput: func(input *awseks.CreateAddonInput) {
				Expect(input.ResolveConflicts).To(Equal(ekstypes.ResolveConflictsOverwrite))
			},
		}),

		Entry("[ResolveConflicts] implicitly set to overwrite by using `--force` flag", createAddonEntry{
			addon: api.Addon{
				Version: "1.0.0",
				Force:   true,
			},
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(provider.MockEKS(), nil)
				mockDescribeAddonVersions(provider.MockEKS(), nil)
				mockCreateAddon(provider.MockEKS(), nil)
			},
			validateCreateAddonInput: func(input *awseks.CreateAddonInput) {
				Expect(input.ResolveConflicts).To(Equal(ekstypes.ResolveConflictsOverwrite))
			},
		}),

		Entry("[ConfigurationValues] are set", createAddonEntry{
			addon: api.Addon{
				Version:             "1.0.0",
				ConfigurationValues: "{\"replicaCount\":3}",
			},
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(provider.MockEKS(), nil)
				mockDescribeAddonVersions(provider.MockEKS(), nil)
				mockCreateAddon(provider.MockEKS(), nil)
			},
			validateCreateAddonInput: func(input *awseks.CreateAddonInput) {
				Expect(*input.ConfigurationValues).To(Equal("{\"replicaCount\":3}"))
			},
		}),

		Entry("[Tags] are set", createAddonEntry{
			addon: api.Addon{
				Version: "1.0.0",
				Tags:    map[string]string{"foo": "bar", "fox": "brown"},
			},
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(provider.MockEKS(), nil)
				mockDescribeAddonVersions(provider.MockEKS(), nil)
				mockCreateAddon(provider.MockEKS(), nil)
			},
			validateCreateAddonInput: func(input *awseks.CreateAddonInput) {
				Expect(input.Tags["foo"]).To(Equal("bar"))
				Expect(input.Tags["fox"]).To(Equal("brown"))
			},
		}),

		Entry("[Wait is true] addon creation succeeds", createAddonEntry{
			addon: api.Addon{
				Version: "1.0.0",
			},
			waitTimeout: time.Nanosecond,
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(provider.MockEKS(), nil)
				// addon becomes active after creation
				mockProvider.MockEKS().
					On("DescribeAddon", mock.Anything, mock.Anything, mock.Anything).
					Return(&awseks.DescribeAddonOutput{
						Addon: &ekstypes.Addon{
							AddonName: aws.String("my-addon"),
							Status:    ekstypes.AddonStatusActive,
						},
					}, nil).
					Once()
				mockDescribeAddonVersions(provider.MockEKS(), nil)
				mockCreateAddon(provider.MockEKS(), nil)
			},
		}),

		Entry("[Wait is true] addon creation fails", createAddonEntry{
			addon: api.Addon{
				Version: "1.0.0",
			},
			waitTimeout: time.Nanosecond,
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(provider.MockEKS(), nil)
				// addon becomes degraded after creation
				mockProvider.MockEKS().
					On("DescribeAddon", mock.Anything, mock.Anything, mock.Anything).
					Return(&awseks.DescribeAddonOutput{
						Addon: &ekstypes.Addon{
							AddonName: aws.String("my-addon"),
							Status:    ekstypes.AddonStatusDegraded,
						},
					}, nil)
				mockDescribeAddonVersions(provider.MockEKS(), nil)
				mockCreateAddon(provider.MockEKS(), nil)
			},
			expectedErr: "addon status transitioned to \"DEGRADED\"",
		}),

		Entry("[Cluster without nodegroups] should not wait for CoreDNS to become active", createAddonEntry{
			addon: api.Addon{
				Name:    api.CoreDNSAddon,
				Version: "1.0.0",
			},
			waitTimeout: time.Nanosecond,
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(provider.MockEKS(), nil)
				// addon becomes degraded after creation
				mockProvider.MockEKS().
					On("DescribeAddon", mock.Anything, mock.Anything, mock.Anything).
					Return(&awseks.DescribeAddonOutput{
						Addon: &ekstypes.Addon{
							AddonName: aws.String("my-addon"),
							Status:    ekstypes.AddonStatusDegraded,
						},
					}, nil)
				mockDescribeAddonVersions(provider.MockEKS(), nil)
				mockCreateAddon(provider.MockEKS(), nil)
			},
		}),

		Entry("[Cluster without nodegroups] should not wait for EBS CSI driver to become active", createAddonEntry{
			addon: api.Addon{
				Name:    api.AWSEBSCSIDriverAddon,
				Version: "1.0.0",
			},
			waitTimeout: time.Nanosecond,
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(provider.MockEKS(), nil)
				// addon becomes degraded after creation
				mockProvider.MockEKS().
					On("DescribeAddon", mock.Anything, mock.Anything, mock.Anything).
					Return(&awseks.DescribeAddonOutput{
						Addon: &ekstypes.Addon{
							AddonName: aws.String("my-addon"),
							Status:    ekstypes.AddonStatusDegraded,
						},
					}, nil)
				mockDescribeAddonVersions(provider.MockEKS(), nil)
				mockCreateAddon(provider.MockEKS(), nil)
			},
		}),

		Entry("[Cluster without nodegroups] should not wait for EFS CSI driver to become active", createAddonEntry{
			addon: api.Addon{
				Name:    api.AWSEFSCSIDriverAddon,
				Version: "1.0.0",
			},
			waitTimeout: time.Nanosecond,
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(provider.MockEKS(), nil)
				// addon becomes degraded after creation
				mockProvider.MockEKS().
					On("DescribeAddon", mock.Anything, mock.Anything, mock.Anything).
					Return(&awseks.DescribeAddonOutput{
						Addon: &ekstypes.Addon{
							AddonName: aws.String("my-addon"),
							Status:    ekstypes.AddonStatusDegraded,
						},
					}, nil)
				mockDescribeAddonVersions(provider.MockEKS(), nil)
				mockCreateAddon(provider.MockEKS(), nil)
			},
		}),

		Entry("[RequiresIAMPermissions] podIDs set explicitly and NOT supportsPodIDs", createAddonEntry{
			addon: api.Addon{
				PodIdentityAssociations: &[]api.PodIdentityAssociation{
					{
						ServiceAccountName: "sa1",
						RoleARN:            "arn:aws:iam::111122223333:role/role-name-1",
					},
				},
			},
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(mockProvider.MockEKS(), nil)
				mockDescribeAddonVersions(provider.MockEKS(), nil)
				mockDescribeAddonConfiguration(mockProvider.MockEKS(), []string{}, nil)
			},
			expectedErr: "\"my-addon\" addon does not support pod identity associations; use IRSA config (`addon.serviceAccountRoleARN`, `addon.attachPolicyARNs`, `addon.attachPolicy` or `addon.wellKnownPolicies`) instead",
		}),

		Entry("[RequiresIAMPermissions] podIDs set explicitly and supportsPodIDs", createAddonEntry{
			addon: api.Addon{
				PodIdentityAssociations: &[]api.PodIdentityAssociation{
					{
						ServiceAccountName: "sa1",
						RoleARN:            "arn:aws:iam::111122223333:role/role-name-1",
					},
					{
						ServiceAccountName: "sa2",
						RoleARN:            "arn:aws:iam::111122223333:role/role-name-2",
					},
				},
			},
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(mockProvider.MockEKS(), nil)
				mockDescribeAddonVersions(provider.MockEKS(), nil)
				mockDescribeAddonConfiguration(mockProvider.MockEKS(), []string{"sa1", "sa2"}, nil)
				mockCreateAddon(mockProvider.MockEKS(), nil)
			},
			validateCreateAddonInput: func(input *awseks.CreateAddonInput) {
				Expect(input.PodIdentityAssociations).To(HaveLen(2))
				Expect(*input.PodIdentityAssociations[0].ServiceAccount).To(Equal("sa1"))
				Expect(*input.PodIdentityAssociations[0].RoleArn).To(Equal("arn:aws:iam::111122223333:role/role-name-1"))
				Expect(*input.PodIdentityAssociations[1].ServiceAccount).To(Equal("sa2"))
				Expect(*input.PodIdentityAssociations[1].RoleArn).To(Equal("arn:aws:iam::111122223333:role/role-name-2"))
			},
			validateCustomLoggerOutput: func(output string) {
				Expect(output).To(ContainSubstring("pod identity associations are set for \"my-addon\" addon; will use these to configure required IAM permissions"))
			},
		}),

		Entry("[RequiresIAMPermissions] `autoApplyPodIdentityAssociations: true` and NOT supportsPodIDs", createAddonEntry{
			mockClusterConfig: func(clusterConfig *api.ClusterConfig) {
				clusterConfig.AddonsConfig.AutoApplyPodIdentityAssociations = true
			},
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(mockProvider.MockEKS(), nil)
				mockDescribeAddonVersions(provider.MockEKS(), nil)
				mockDescribeAddonConfiguration(mockProvider.MockEKS(), []string{}, nil)
				mockCreateAddon(mockProvider.MockEKS(), nil)
			},
			validateCreateAddonInput: func(input *awseks.CreateAddonInput) {
				Expect(input.PodIdentityAssociations).To(HaveLen(0))
			},
			validateCustomLoggerOutput: func(output string) {
				Expect(output).To(ContainSubstring("IAM permissions are required for \"my-addon\" addon; " +
					"the recommended way to provide IAM permissions for \"my-addon\" addon is via IRSA; " +
					"after addon creation is completed, add all recommended policies to the config file, " +
					"under `addon.ServiceAccountRoleARN`, `addon.AttachPolicyARNs`, `addon.AttachPolicy` or `addon.WellKnownPolicies`, " +
					"and run `eksctl update addon`"))
			},
		}),

		Entry("[RequiresIAMPermissions] `autoApplyPodIdentityAssociations: true` and supportsPodIDs", createAddonEntry{
			mockClusterConfig: func(clusterConfig *api.ClusterConfig) {
				clusterConfig.AddonsConfig.AutoApplyPodIdentityAssociations = true
			},
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(mockProvider.MockEKS(), nil)
				mockDescribeAddonVersions(provider.MockEKS(), nil)
				mockDescribeAddonConfiguration(mockProvider.MockEKS(), []string{"sa1"}, nil)
				mockCreateAddon(mockProvider.MockEKS(), nil)
			},
			mockIAM: func(mockIAMRoleCreator *addonmocks.IAMRoleCreator) {
				mockIAMRoleCreator.
					On("Create", mock.Anything, mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(3))
						Expect(args[1]).To(BeAssignableToTypeOf(&api.PodIdentityAssociation{}))
						Expect(args[1].(*api.PodIdentityAssociation).ServiceAccountName).To(Equal("sa1"))
						Expect(args[1].(*api.PodIdentityAssociation).PermissionPolicyARNs).To(ConsistOf("arn:aws:iam::111122223333:policy/sa1"))
					}).
					Return("arn:aws:iam::111122223333:role/sa1", nil).
					Once()
			},
			validateCreateAddonInput: func(input *awseks.CreateAddonInput) {
				Expect(input.PodIdentityAssociations).To(HaveLen(1))
				Expect(*input.PodIdentityAssociations[0].ServiceAccount).To(Equal("sa1"))
				Expect(*input.PodIdentityAssociations[0].RoleArn).To(Equal("arn:aws:iam::111122223333:role/sa1"))
			},
			validateCustomLoggerOutput: func(output string) {
				Expect(output).To(ContainSubstring("\"addonsConfig.autoApplyPodIdentityAssociations\" is set to true; will lookup recommended pod identity configuration for \"my-addon\" addon"))
			},
		}),

		Entry("[RequiresIAMPermissions] `autoApplyPodIdentityAssociations: true` and supportsPodIDs (vpc-cni && ipv6)", createAddonEntry{
			addon: api.Addon{
				Name: api.VPCCNIAddon,
			},
			mockK8s: true,
			mockClusterConfig: func(clusterConfig *api.ClusterConfig) {
				clusterConfig.AddonsConfig.AutoApplyPodIdentityAssociations = true
				clusterConfig.KubernetesNetworkConfig = &api.KubernetesNetworkConfig{
					IPFamily: "IPv6",
				}
			},
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(mockProvider.MockEKS(), nil)
				mockDescribeAddonVersions(provider.MockEKS(), nil)
				mockDescribeAddonConfiguration(mockProvider.MockEKS(), []string{"sa1"}, nil)
				mockCreateAddon(mockProvider.MockEKS(), nil)
			},
			mockIAM: func(mockIAMRoleCreator *addonmocks.IAMRoleCreator) {
				mockIAMRoleCreator.
					On("Create", mock.Anything, mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(3))
						Expect(args[1]).To(BeAssignableToTypeOf(&api.PodIdentityAssociation{}))
						Expect(args[1].(*api.PodIdentityAssociation).ServiceAccountName).To(Equal("aws-node"))
						Expect(args[1].(*api.PodIdentityAssociation).PermissionPolicy).NotTo(BeEmpty())
					}).
					Return("arn:aws:iam::111122223333:role/aws-node", nil).
					Once()
			},
			validateCreateAddonInput: func(input *awseks.CreateAddonInput) {
				Expect(input.PodIdentityAssociations).To(HaveLen(1))
				Expect(*input.PodIdentityAssociations[0].ServiceAccount).To(Equal("aws-node"))
				Expect(*input.PodIdentityAssociations[0].RoleArn).To(Equal("arn:aws:iam::111122223333:role/aws-node"))
			},
			validateCustomLoggerOutput: func(output string) {
				Expect(output).To(ContainSubstring("\"addonsConfig.autoApplyPodIdentityAssociations\" is set to true; will lookup recommended pod identity configuration for \"vpc-cni\" addon"))
			},
		}),

		Entry("[RequiresIAMPermissions] podIDs already exist on cluster", createAddonEntry{
			addon: api.Addon{
				PodIdentityAssociations: &[]api.PodIdentityAssociation{
					{
						ServiceAccountName: "sa1",
						RoleARN:            "arn:aws:iam::111122223333:role/role-name-1",
					},
					{
						ServiceAccountName: "sa2",
						RoleARN:            "arn:aws:iam::111122223333:role/role-name-2",
					},
				},
			},
			mockClusterConfig: func(clusterConfig *api.ClusterConfig) {
				clusterConfig.AddonsConfig.AutoApplyPodIdentityAssociations = true
			},
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(mockProvider.MockEKS(), nil)
				mockDescribeAddonVersions(provider.MockEKS(), nil)
				mockDescribeAddonConfiguration(mockProvider.MockEKS(), []string{"sa1", "sa2"}, nil)
				mockCreateAddon(mockProvider.MockEKS(), &ekstypes.ResourceInUseException{})
			},
			mockCFN: func(stackManager *fakes.FakeStackManager) {
				stackManager.DeleteStackBySpecReturns(nil, nil)
			},
			validateCFNCalls: func(stackManager *fakes.FakeStackManager) {
				Expect(stackManager.DeleteStackBySpecCallCount()).To(Equal(2))
			},
			expectedErr: "creating addon: one or more service accounts corresponding to \"my-addon\" addon is already associated with a different IAM role; " +
				"please delete all pre-existing pod identity associations corresponding to \"sa1\",\"sa2\" service account(s) in the addon's namespace, then re-try creating the addon",
		}),

		Entry("[RequiresIAMPermissions] IRSA set explicitly and NOT supportsPodIDs", createAddonEntry{
			addon: api.Addon{
				ServiceAccountRoleARN: "arn:aws:iam::111122223333:role/role-name-1",
			},
			withOIDC: true,
			mockClusterConfig: func(clusterConfig *api.ClusterConfig) {
				clusterConfig.AddonsConfig.AutoApplyPodIdentityAssociations = true
			},
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(mockProvider.MockEKS(), nil)
				mockDescribeAddonVersions(provider.MockEKS(), nil)
				mockDescribeAddonConfiguration(mockProvider.MockEKS(), []string{}, nil)
				mockCreateAddon(mockProvider.MockEKS(), nil)
			},
			validateCreateAddonInput: func(input *awseks.CreateAddonInput) {
				Expect(input.PodIdentityAssociations).To(HaveLen(0))
				Expect(input.ServiceAccountRoleArn).NotTo(BeNil())
				Expect(*input.ServiceAccountRoleArn).To(Equal("arn:aws:iam::111122223333:role/role-name-1"))
			},
			validateCustomLoggerOutput: func(output string) {
				Expect(output).NotTo(ContainSubstring("the recommended way to provide IAM permissions"))
			},
		}),

		Entry("[RequiresIAMPermissions] IRSA set explicitly and supportsPodIDs", createAddonEntry{
			addon: api.Addon{
				AttachPolicy: api.InlineDocument{
					"foo": "policy-bar",
				},
			},
			withOIDC: true,
			mockClusterConfig: func(clusterConfig *api.ClusterConfig) {
				clusterConfig.AddonsConfig.AutoApplyPodIdentityAssociations = true
			},
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(mockProvider.MockEKS(), nil)
				mockDescribeAddonVersions(provider.MockEKS(), nil)
				mockDescribeAddonConfiguration(mockProvider.MockEKS(), []string{"sa1"}, nil)
				mockCreateAddon(mockProvider.MockEKS(), nil)
			},
			mockCFN: func(stackManager *fakes.FakeStackManager) {
				stackManager.CreateStackStub = func(ctx context.Context, s string, rsr builder.ResourceSetReader, m1, m2 map[string]string, c chan error) error {
					go func() {
						c <- nil
					}()
					Expect(rsr).To(BeAssignableToTypeOf(&builder.IAMRoleResourceSet{}))
					output, err := rsr.(*builder.IAMRoleResourceSet).RenderJSON()
					Expect(err).NotTo(HaveOccurred())
					Expect(string(output)).To(ContainSubstring("policy-bar"))
					rsr.(*builder.IAMRoleResourceSet).OutputRole = "arn:aws:iam::111122223333:role/role-name-1"
					return nil
				}
			},
			validateCreateAddonInput: func(input *awseks.CreateAddonInput) {
				Expect(input.PodIdentityAssociations).To(HaveLen(0))
				Expect(input.ServiceAccountRoleArn).NotTo(BeNil())
				Expect(*input.ServiceAccountRoleArn).To(Equal("arn:aws:iam::111122223333:role/role-name-1"))
			},
			validateCFNCalls: func(stackManager *fakes.FakeStackManager) {
				Expect(stackManager.CreateStackCallCount()).To(Equal(1))
			},
			validateCustomLoggerOutput: func(output string) {
				Expect(output).To(ContainSubstring("the recommended way to provide IAM permissions for \"my-addon\" addon is via pod identity associations; " +
					"after addon creation is completed, run `eksctl utils migrate-to-pod-identity`"))
			},
		}),

		Entry("[RequiresIAMPermissions] IRSA set implicitly (vpc-cni && ipv4)", createAddonEntry{
			addon: api.Addon{
				Name: api.VPCCNIAddon,
			},
			withOIDC: true,
			mockK8s:  true,
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(mockProvider.MockEKS(), nil)
				mockDescribeAddonVersions(provider.MockEKS(), nil)
				mockDescribeAddonConfiguration(mockProvider.MockEKS(), []string{"aws-node"}, nil)
				mockCreateAddon(mockProvider.MockEKS(), nil)
			},
			mockCFN: func(stackManager *fakes.FakeStackManager) {
				stackManager.CreateStackStub = func(ctx context.Context, s string, rsr builder.ResourceSetReader, m1, m2 map[string]string, c chan error) error {
					go func() {
						c <- nil
					}()
					Expect(rsr).To(BeAssignableToTypeOf(&builder.IAMRoleResourceSet{}))
					output, err := rsr.(*builder.IAMRoleResourceSet).RenderJSON()
					Expect(err).NotTo(HaveOccurred())
					Expect(string(output)).To(ContainSubstring("arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"))
					Expect(string(output)).To(ContainSubstring(":sub\":\"system:serviceaccount:kube-system:aws-node"))
					rsr.(*builder.IAMRoleResourceSet).OutputRole = "arn:aws:iam::111122223333:role/role-name-1"
					return nil
				}
			},
			validateCreateAddonInput: func(input *awseks.CreateAddonInput) {
				Expect(input.PodIdentityAssociations).To(HaveLen(0))
				Expect(input.ServiceAccountRoleArn).NotTo(BeNil())
				Expect(*input.ServiceAccountRoleArn).To(Equal("arn:aws:iam::111122223333:role/role-name-1"))
			},
			validateCFNCalls: func(stackManager *fakes.FakeStackManager) {
				Expect(stackManager.CreateStackCallCount()).To(Equal(1))
			},
			validateCustomLoggerOutput: func(output string) {
				Expect(output).To(ContainSubstring("the recommended way to provide IAM permissions for \"vpc-cni\" addon is via pod identity associations; " +
					"after addon creation is completed, run `eksctl utils migrate-to-pod-identity`"))
			},
		}),

		Entry("[RequiresIAMPermissions] IRSA set implicitly (vpc-cni && ipv6)", createAddonEntry{
			addon: api.Addon{
				Name: api.VPCCNIAddon,
			},
			withOIDC: true,
			mockK8s:  true,
			mockClusterConfig: func(clusterConfig *api.ClusterConfig) {
				clusterConfig.KubernetesNetworkConfig = &api.KubernetesNetworkConfig{
					IPFamily: "IPv6",
				}
			},
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(mockProvider.MockEKS(), nil)
				mockDescribeAddonVersions(provider.MockEKS(), nil)
				mockDescribeAddonConfiguration(mockProvider.MockEKS(), []string{"aws-node"}, nil)
				mockCreateAddon(mockProvider.MockEKS(), nil)
			},
			mockCFN: func(stackManager *fakes.FakeStackManager) {
				stackManager.CreateStackStub = func(ctx context.Context, s string, rsr builder.ResourceSetReader, m1, m2 map[string]string, c chan error) error {
					go func() {
						c <- nil
					}()
					Expect(rsr).To(BeAssignableToTypeOf(&builder.IAMRoleResourceSet{}))
					output, err := rsr.(*builder.IAMRoleResourceSet).RenderJSON()
					Expect(err).NotTo(HaveOccurred())
					Expect(string(output)).To(ContainSubstring("AssignIpv6Addresses"))
					rsr.(*builder.IAMRoleResourceSet).OutputRole = "arn:aws:iam::111122223333:role/role-name-1"
					return nil
				}
			},
			validateCreateAddonInput: func(input *awseks.CreateAddonInput) {
				Expect(input.PodIdentityAssociations).To(HaveLen(0))
				Expect(input.ServiceAccountRoleArn).NotTo(BeNil())
				Expect(*input.ServiceAccountRoleArn).To(Equal("arn:aws:iam::111122223333:role/role-name-1"))
			},
			validateCFNCalls: func(stackManager *fakes.FakeStackManager) {
				Expect(stackManager.CreateStackCallCount()).To(Equal(1))
			},
			validateCustomLoggerOutput: func(output string) {
				Expect(output).To(ContainSubstring("the recommended way to provide IAM permissions for \"vpc-cni\" addon is via pod identity associations; " +
					"after addon creation is completed, run `eksctl utils migrate-to-pod-identity`"))
			},
		}),

		Entry("[RequiresIAMPermissions] IRSA set implicitly and supportsPodIDs (aws-ebs-csi-driver)", createAddonEntry{
			addon: api.Addon{
				Name: api.AWSEBSCSIDriverAddon,
			},
			withOIDC: true,
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(mockProvider.MockEKS(), nil)
				mockDescribeAddonVersions(provider.MockEKS(), nil)
				mockDescribeAddonConfiguration(mockProvider.MockEKS(), []string{"sa1"}, nil)
				mockCreateAddon(mockProvider.MockEKS(), nil)
			},
			mockCFN: func(stackManager *fakes.FakeStackManager) {
				stackManager.CreateStackStub = func(ctx context.Context, s string, rsr builder.ResourceSetReader, m1, m2 map[string]string, c chan error) error {
					go func() {
						c <- nil
					}()
					Expect(rsr).To(BeAssignableToTypeOf(&builder.IAMRoleResourceSet{}))
					output, err := rsr.(*builder.IAMRoleResourceSet).RenderJSON()
					Expect(err).NotTo(HaveOccurred())
					Expect(string(output)).To(ContainSubstring("PolicyEBSCSIController"))
					rsr.(*builder.IAMRoleResourceSet).OutputRole = "arn:aws:iam::111122223333:role/role-name-1"
					return nil
				}
			},
			validateCreateAddonInput: func(input *awseks.CreateAddonInput) {
				Expect(input.PodIdentityAssociations).To(HaveLen(0))
				Expect(input.ServiceAccountRoleArn).NotTo(BeNil())
				Expect(*input.ServiceAccountRoleArn).To(Equal("arn:aws:iam::111122223333:role/role-name-1"))
			},
			validateCFNCalls: func(stackManager *fakes.FakeStackManager) {
				Expect(stackManager.CreateStackCallCount()).To(Equal(1))
			},
			validateCustomLoggerOutput: func(output string) {
				Expect(output).To(ContainSubstring("the recommended way to provide IAM permissions for \"aws-ebs-csi-driver\" addon is via pod identity associations; " +
					"after addon creation is completed, run `eksctl utils migrate-to-pod-identity`"))
			},
		}),

		Entry("[RequiresIAMPermissions] IRSA set implicitly and not supportsPodIDs (aws-efs-csi-driver)", createAddonEntry{
			addon: api.Addon{
				Name: api.AWSEFSCSIDriverAddon,
			},
			withOIDC: true,
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(mockProvider.MockEKS(), nil)
				mockDescribeAddonVersions(provider.MockEKS(), nil)
				mockDescribeAddonConfiguration(mockProvider.MockEKS(), []string{}, nil)
				mockCreateAddon(mockProvider.MockEKS(), nil)
			},
			mockCFN: func(stackManager *fakes.FakeStackManager) {
				stackManager.CreateStackStub = func(ctx context.Context, s string, rsr builder.ResourceSetReader, m1, m2 map[string]string, c chan error) error {
					go func() {
						c <- nil
					}()
					Expect(rsr).To(BeAssignableToTypeOf(&builder.IAMRoleResourceSet{}))
					output, err := rsr.(*builder.IAMRoleResourceSet).RenderJSON()
					Expect(err).NotTo(HaveOccurred())
					Expect(string(output)).To(ContainSubstring("PolicyEFSCSIController"))
					rsr.(*builder.IAMRoleResourceSet).OutputRole = "arn:aws:iam::111122223333:role/role-name-1"
					return nil
				}
			},
			validateCreateAddonInput: func(input *awseks.CreateAddonInput) {
				Expect(input.PodIdentityAssociations).To(HaveLen(0))
				Expect(input.ServiceAccountRoleArn).NotTo(BeNil())
				Expect(*input.ServiceAccountRoleArn).To(Equal("arn:aws:iam::111122223333:role/role-name-1"))
			},
			validateCFNCalls: func(stackManager *fakes.FakeStackManager) {
				Expect(stackManager.CreateStackCallCount()).To(Equal(1))
			},
			validateCustomLoggerOutput: func(output string) {
				Expect(output).NotTo(ContainSubstring("the recommended way to provide IAM permissions"))
			},
		}),

		Entry("[RequiresIAMPermissions] OIDC is disabled and IRSA set explicitly", createAddonEntry{
			addon: api.Addon{
				ServiceAccountRoleARN: "arn:aws:iam::111122223333:role/role-name-1",
			},
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(mockProvider.MockEKS(), nil)
				mockDescribeAddonVersions(provider.MockEKS(), nil)
				mockDescribeAddonConfiguration(mockProvider.MockEKS(), []string{"sa1"}, nil)
				mockCreateAddon(mockProvider.MockEKS(), nil)
			},
			validateCreateAddonInput: func(input *awseks.CreateAddonInput) {
				Expect(input.PodIdentityAssociations).To(HaveLen(0))
				Expect(input.ServiceAccountRoleArn).To(BeNil())
			},
			validateCustomLoggerOutput: func(output string) {
				Expect(output).To(ContainSubstring("IRSA config is set for \"my-addon\" addon, " +
					"but since OIDC is disabled on the cluster, eksctl cannot configure the requested permissions; " +
					"the recommended way to provide IAM permissions for \"my-addon\" addon is via pod identity associations; " +
					"after addon creation is completed, add all recommended policies to the config file, under `addon.PodIdentityAssociations`, " +
					"and run `eksctl update addon`"))
			},
		}),

		Entry("[RequiresIAMPermissions] OIDC is disabled and IRSA set implicitly", createAddonEntry{
			addon: api.Addon{
				Name: api.AWSEFSCSIDriverAddon,
			},
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(mockProvider.MockEKS(), nil)
				mockDescribeAddonVersions(provider.MockEKS(), nil)
				mockDescribeAddonConfiguration(mockProvider.MockEKS(), []string{}, nil)
				mockCreateAddon(mockProvider.MockEKS(), nil)
			},
			validateCreateAddonInput: func(input *awseks.CreateAddonInput) {
				Expect(input.PodIdentityAssociations).To(HaveLen(0))
				Expect(input.ServiceAccountRoleArn).To(BeNil())
			},
			validateCustomLoggerOutput: func(output string) {
				Expect(output).To(ContainSubstring("recommended policies were found for \"aws-efs-csi-driver\" addon, " +
					"but since OIDC is disabled on the cluster, eksctl cannot configure the requested permissions; " +
					"users are responsible for attaching the policies to all nodegroup roles"))
			},
		}),

		Entry("[RequiresIAMPermissions] neither IRSA nor podIDs are being set and NOT supportsPodIDs", createAddonEntry{
			addon: api.Addon{},
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(mockProvider.MockEKS(), nil)
				mockDescribeAddonVersions(provider.MockEKS(), nil)
				mockDescribeAddonConfiguration(mockProvider.MockEKS(), []string{}, nil)
				mockCreateAddon(mockProvider.MockEKS(), nil)
			},
			validateCreateAddonInput: func(input *awseks.CreateAddonInput) {
				Expect(input.PodIdentityAssociations).To(HaveLen(0))
				Expect(input.ServiceAccountRoleArn).To(BeNil())
			},
			validateCustomLoggerOutput: func(output string) {
				Expect(output).To(ContainSubstring("IAM permissions are required for \"my-addon\" addon; " +
					"the recommended way to provide IAM permissions for \"my-addon\" addon is via IRSA; " +
					"after addon creation is completed, add all recommended policies to the config file, " +
					"under `addon.ServiceAccountRoleARN`, `addon.AttachPolicyARNs`, `addon.AttachPolicy` or `addon.WellKnownPolicies`, " +
					"and run `eksctl update addon`"))
			},
		}),

		Entry("[RequiresIAMPermissions] neither IRSA nor podIDs are being set and supportsPodIDs", createAddonEntry{
			addon: api.Addon{},
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(mockProvider.MockEKS(), nil)
				mockDescribeAddonVersions(provider.MockEKS(), nil)
				mockDescribeAddonConfiguration(mockProvider.MockEKS(), []string{"sa1"}, nil)
				mockCreateAddon(mockProvider.MockEKS(), nil)
			},
			validateCreateAddonInput: func(input *awseks.CreateAddonInput) {
				Expect(input.PodIdentityAssociations).To(HaveLen(0))
				Expect(input.ServiceAccountRoleArn).To(BeNil())
			},
			validateCustomLoggerOutput: func(output string) {
				Expect(output).To(ContainSubstring("IAM permissions are required for \"my-addon\" addon; " +
					"the recommended way to provide IAM permissions for \"my-addon\" addon is via pod identity associations; " +
					"after addon creation is completed, add all recommended policies to the config file, " +
					"under `addon.PodIdentityAssociations`, and run `eksctl update addon`"))
			},
		}),

		Entry("[RequiresIAMPermissions is false] podIDs set", createAddonEntry{
			addon: api.Addon{
				Version: "1.0.0",
				PodIdentityAssociations: &[]api.PodIdentityAssociation{
					{
						ServiceAccountName: "sa1",
						RoleARN:            "arn:aws:iam::111122223333:role/role-name-1",
					},
				},
			},
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(mockProvider.MockEKS(), nil)
				mockDescribeAddonVersions(provider.MockEKS(), nil)
				mockDescribeAddonConfiguration(mockProvider.MockEKS(), []string{}, nil)
				mockCreateAddon(mockProvider.MockEKS(), nil)
			},
			validateCreateAddonInput: func(input *awseks.CreateAddonInput) {
				Expect(input.PodIdentityAssociations).To(HaveLen(0))
			},
			validateCustomLoggerOutput: func(output string) {
				Expect(output).To(ContainSubstring("IAM permissions are not required for \"my-addon\" addon; " +
					"any IRSA configuration or pod identity associations will be ignored"))
			},
		}),

		Entry("[RequiresIAMPermissions is false] IRSA set", createAddonEntry{
			addon: api.Addon{
				Version:               "1.0.0",
				ServiceAccountRoleARN: "arn:aws:iam::111122223333:role/role-name-1",
			},
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockDescribeAddon(mockProvider.MockEKS(), nil)
				mockDescribeAddonVersions(provider.MockEKS(), nil)
				mockDescribeAddonConfiguration(mockProvider.MockEKS(), []string{}, nil)
				mockCreateAddon(mockProvider.MockEKS(), nil)
			},
			validateCreateAddonInput: func(input *awseks.CreateAddonInput) {
				Expect(input.ServiceAccountRoleArn).To(BeNil())
			},
			validateCustomLoggerOutput: func(output string) {
				Expect(output).To(ContainSubstring("IAM permissions are not required for \"my-addon\" addon; " +
					"any IRSA configuration or pod identity associations will be ignored"))
			},
		}),
	)
})
