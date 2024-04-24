package manager_test

import (
	"context"
	"fmt"

	"github.com/weaveworks/eksctl/pkg/cfn/manager/mocks"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	bootstrapfakes "github.com/weaveworks/eksctl/pkg/nodebootstrap/fakes"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("Unmanaged Nodegroup Task", func() {

	type resourceSetOptions struct {
		nodeGroupName              string
		disableAccessEntryResource bool
	}

	type ngEntry struct {
		nodeGroups []*api.NodeGroup
		mockCalls  func(*mockprovider.MockProvider, *mocks.NodeGroupStackManager)

		expectedResourceSetOptions []resourceSetOptions

		expectedErr string
	}

	newNodeGroups := func(instanceRoleARNs ...string) []*api.NodeGroup {
		var nodeGroups []*api.NodeGroup
		for i, roleARN := range instanceRoleARNs {
			ng := api.NewNodeGroup()
			ng.Name = fmt.Sprintf("ng-%d", i+1)
			ng.IAM.InstanceRoleARN = roleARN
			nodeGroups = append(nodeGroups, ng)
		}
		return nodeGroups
	}

	const clusterName = "cluster"

	makeAccessEntryInput := func(roleName string) *eks.CreateAccessEntryInput {
		return &eks.CreateAccessEntryInput{
			ClusterName:  aws.String(clusterName),
			PrincipalArn: aws.String(roleName),
			Type:         aws.String("EC2_LINUX"),
			Tags: map[string]string{
				api.ClusterNameLabel: clusterName,
			},
		}
	}

	DescribeTable("Create", func(e ngEntry) {
		clusterConfig := api.NewClusterConfig()
		clusterConfig.Metadata.Name = clusterName
		clusterConfig.NodeGroups = e.nodeGroups

		var (
			provider             = mockprovider.NewMockProvider()
			stackManager         mocks.NodeGroupStackManager
			resourceSetArgs      []resourceSetOptions
			nodeGroupResourceSet mocks.NodeGroupResourceSet
		)
		nodeGroupResourceSet.On("AddAllResources", mock.Anything).Return(nil).Times(len(clusterConfig.NodeGroups))

		t := &manager.UnmanagedNodeGroupTask{
			ClusterConfig: clusterConfig,
			NodeGroups:    clusterConfig.NodeGroups,
			CreateNodeGroupResourceSet: func(options builder.NodeGroupOptions) manager.NodeGroupResourceSet {
				resourceSetArgs = append(resourceSetArgs, resourceSetOptions{
					nodeGroupName:              options.NodeGroup.Name,
					disableAccessEntryResource: options.DisableAccessEntryResource,
				})
				return &nodeGroupResourceSet
			},
			NewBootstrapper: func(_ *api.ClusterConfig, _ *api.NodeGroup) (nodebootstrap.Bootstrapper, error) {
				var bootstrapper bootstrapfakes.FakeBootstrapper
				bootstrapper.UserDataReturns("", nil)
				return &bootstrapper, nil
			},
			EKSAPI:       provider.EKS(),
			StackManager: &stackManager,
		}

		stackManager.On("CreateStack", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.MatchedBy(func(errs chan error) bool {
			close(errs)
			return true
		})).Return(nil)

		if e.mockCalls != nil {
			e.mockCalls(provider, &stackManager)
		}

		taskTree := t.Create(context.Background(), manager.CreateNodeGroupOptions{
			ForceAddCNIPolicy:          false,
			SkipEgressRules:            false,
			DisableAccessEntryCreation: false,
		})
		errs := taskTree.DoAllSync()
		if e.expectedErr != "" {
			Expect(errs).To(ContainElement(MatchError(ContainSubstring(e.expectedErr))))
			return
		}

		Expect(errs).To(HaveLen(0))
		Expect(e.expectedResourceSetOptions).To(ConsistOf(resourceSetArgs))
		mock.AssertExpectationsForObjects(GinkgoT(), provider.MockEKS(), &stackManager, &nodeGroupResourceSet)
	},
		Entry("nodegroups with no instanceRoleARN set", ngEntry{
			nodeGroups: newNodeGroups("", ""),
			expectedResourceSetOptions: []resourceSetOptions{
				{
					nodeGroupName:              "ng-1",
					disableAccessEntryResource: false,
				},
				{
					nodeGroupName:              "ng-2",
					disableAccessEntryResource: false,
				},
			},
		}),

		Entry("some nodegroups with instanceRoleARN set", ngEntry{
			nodeGroups: newNodeGroups("", "role-1", ""),
			mockCalls: func(provider *mockprovider.MockProvider, stackManager *mocks.NodeGroupStackManager) {
				provider.MockEKS().On("CreateAccessEntry", mock.Anything, makeAccessEntryInput("role-1")).Return(&eks.CreateAccessEntryOutput{}, nil)
			},
			expectedResourceSetOptions: []resourceSetOptions{
				{
					nodeGroupName:              "ng-1",
					disableAccessEntryResource: false,
				},
				{
					nodeGroupName:              "ng-2",
					disableAccessEntryResource: true,
				},
				{
					nodeGroupName:              "ng-3",
					disableAccessEntryResource: false,
				},
			},
		}),

		Entry("all nodegroups with instanceRoleARN set", ngEntry{
			nodeGroups: newNodeGroups("role-1", "role-2"),
			mockCalls: func(provider *mockprovider.MockProvider, stackManager *mocks.NodeGroupStackManager) {
				for _, role := range []string{"role-1", "role-2"} {
					provider.MockEKS().On("CreateAccessEntry", mock.Anything, makeAccessEntryInput(role)).Return(&eks.CreateAccessEntryOutput{}, nil).Once()
				}

			},
			expectedResourceSetOptions: []resourceSetOptions{
				{
					nodeGroupName:              "ng-1",
					disableAccessEntryResource: true,
				},
				{
					nodeGroupName:              "ng-2",
					disableAccessEntryResource: true,
				},
			},
		}),

		Entry("all nodegroups with the same instanceRoleARN", ngEntry{
			nodeGroups: newNodeGroups("role-3", "role-3"),
			mockCalls: func(provider *mockprovider.MockProvider, stackManager *mocks.NodeGroupStackManager) {
				provider.MockEKS().On("CreateAccessEntry", mock.Anything, makeAccessEntryInput("role-3")).Return(&eks.CreateAccessEntryOutput{}, nil).Once()
				provider.MockEKS().On("CreateAccessEntry", mock.Anything, makeAccessEntryInput("role-3")).Return(nil, &ekstypes.ResourceInUseException{
					ClusterName: aws.String(clusterName),
				}).Once()

			},
			expectedResourceSetOptions: []resourceSetOptions{
				{
					nodeGroupName:              "ng-1",
					disableAccessEntryResource: true,
				},
				{
					nodeGroupName:              "ng-2",
					disableAccessEntryResource: true,
				},
			},
		}),

		Entry("single nodegroup with a pre-existing access entry", ngEntry{
			nodeGroups: newNodeGroups("role-3"),
			mockCalls: func(provider *mockprovider.MockProvider, stackManager *mocks.NodeGroupStackManager) {
				provider.MockEKS().On("CreateAccessEntry", mock.Anything, makeAccessEntryInput("role-3")).Return(&eks.CreateAccessEntryOutput{}, &ekstypes.ResourceInUseException{
					ClusterName: aws.String(clusterName),
				}).Once()

			},
			expectedResourceSetOptions: []resourceSetOptions{
				{
					nodeGroupName:              "ng-1",
					disableAccessEntryResource: true,
				},
			},
		}),
	)

})
