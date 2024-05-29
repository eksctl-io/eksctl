package podidentityassociation_test

import (
	"context"
	"fmt"

	awseks "github.com/aws/aws-sdk-go-v2/service/eks"

	"k8s.io/apimachinery/pkg/runtime"
	kubeclientfakes "k8s.io/client-go/kubernetes/fake"
	kubeclienttesting "k8s.io/client-go/testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/weaveworks/eksctl/pkg/actions/podidentityassociation"
	"github.com/weaveworks/eksctl/pkg/actions/podidentityassociation/fakes"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

type createPodIdentityAssociationEntry struct {
	toBeCreated              []api.PodIdentityAssociation
	mockEKS                  func(provider *mockprovider.MockProvider)
	mockCFN                  func(stackCreator *fakes.FakeStackCreator)
	mockK8s                  func(clientSet *kubeclientfakes.Clientset)
	expectedCreateStackCalls int
	expectedErr              string
}

var _ = Describe("Create", func() {
	var (
		creator          *podidentityassociation.Creator
		fakeStackCreator *fakes.FakeStackCreator
		fakeClientSet    *kubeclientfakes.Clientset
		mockProvider     *mockprovider.MockProvider

		clusterName         = "test-cluster"
		namespace           = "test-namespace"
		serviceAccountName1 = "test-service-account-name-1"
		serviceAccountName2 = "test-service-account-name-2"
		genericErr          = fmt.Errorf("ERR")
		roleARN             = "arn:aws:iam::111122223333:role/TestRole"
	)

	DescribeTable("Create", func(e createPodIdentityAssociationEntry) {
		fakeStackCreator = new(fakes.FakeStackCreator)
		if e.mockCFN != nil {
			e.mockCFN(fakeStackCreator)
		}

		fakeClientSet = kubeclientfakes.NewSimpleClientset()
		if e.mockK8s != nil {
			e.mockK8s(fakeClientSet)
		}

		mockProvider = mockprovider.NewMockProvider()
		if e.mockEKS != nil {
			e.mockEKS(mockProvider)
		}

		creator = podidentityassociation.NewCreator(clusterName, fakeStackCreator, mockProvider.MockEKS(), fakeClientSet)

		err := creator.CreatePodIdentityAssociations(context.Background(), e.toBeCreated)
		if e.expectedErr != "" {
			Expect(err).To(MatchError(ContainSubstring(e.expectedErr)))
			return
		}
		Expect(err).ToNot(HaveOccurred())
		Expect(fakeStackCreator.CreateStackCallCount()).To(Equal(e.expectedCreateStackCalls))
	},
		Entry("returns an error if creating the IAM role fails", createPodIdentityAssociationEntry{
			toBeCreated: []api.PodIdentityAssociation{
				{
					Namespace:          namespace,
					ServiceAccountName: serviceAccountName1,
				},
			},
			mockCFN: func(stackCreator *fakes.FakeStackCreator) {
				stackCreator.CreateStackStub = func(ctx context.Context, s string, rsr builder.ResourceSetReader, m1, m2 map[string]string, c chan error) error {
					defer close(c)
					Expect(s).To(Equal(podidentityassociation.MakeStackName(
						clusterName,
						namespace,
						serviceAccountName1,
					)))
					return genericErr
				}
			},
			expectedErr: "creating IAM role for pod identity association",
		}),

		Entry("returns an error if creating the service account fails", createPodIdentityAssociationEntry{
			toBeCreated: []api.PodIdentityAssociation{
				{
					Namespace:            namespace,
					ServiceAccountName:   serviceAccountName1,
					RoleARN:              roleARN,
					CreateServiceAccount: true,
				},
			},
			mockK8s: func(clientSet *kubeclientfakes.Clientset) {
				clientSet.PrependReactor("get", "namespaces", func(action kubeclienttesting.Action) (bool, runtime.Object, error) {
					return true, nil, genericErr
				})
			},
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockProvider.MockEKS().
					On("CreatePodIdentityAssociation", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(2))
						Expect(args[1]).To(BeAssignableToTypeOf(&awseks.CreatePodIdentityAssociationInput{}))
					}).
					Return(&awseks.CreatePodIdentityAssociationOutput{}, nil).
					Once()
			},
			expectedErr: "failed to create service account",
		}),

		Entry("returns an error if creating the association fails", createPodIdentityAssociationEntry{
			toBeCreated: []api.PodIdentityAssociation{
				{
					Namespace:          namespace,
					ServiceAccountName: serviceAccountName1,
					RoleARN:            roleARN,
				},
			},
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockProvider.MockEKS().
					On("CreatePodIdentityAssociation", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(2))
						Expect(args[1]).To(BeAssignableToTypeOf(&awseks.CreatePodIdentityAssociationInput{}))
					}).
					Return(nil, genericErr).
					Once()
			},
			expectedErr: "creating pod identity association",
		}),

		Entry("creates all expected roles and associations successfully", createPodIdentityAssociationEntry{
			toBeCreated: []api.PodIdentityAssociation{
				{
					Namespace:          namespace,
					ServiceAccountName: serviceAccountName1,
					RoleARN:            roleARN,
				},
				{
					Namespace:          namespace,
					ServiceAccountName: serviceAccountName2,
				},
			},
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockProvider.MockEKS().
					On("CreatePodIdentityAssociation", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(2))
						Expect(args[1]).To(BeAssignableToTypeOf(&awseks.CreatePodIdentityAssociationInput{}))
					}).
					Return(&awseks.CreatePodIdentityAssociationOutput{}, nil).
					Twice()
			},
			mockCFN: func(stackCreator *fakes.FakeStackCreator) {
				stackCreator.CreateStackStub = func(ctx context.Context, s string, rsr builder.ResourceSetReader, m1, m2 map[string]string, c chan error) error {
					defer close(c)
					return nil
				}
			},
			expectedCreateStackCalls: 1,
		}),
	)
})
