package nodegroup_test

import (
	"context"
	"errors"
	"sync/atomic"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"
	"github.com/weaveworks/eksctl/pkg/actions/nodegroup/fakes"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	managerfakes "github.com/weaveworks/eksctl/pkg/cfn/manager/fakes"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
	taskfakes "github.com/weaveworks/eksctl/pkg/utils/tasks/fakes"
)

var _ = Describe("Delete", func() {
	type deleteErrors struct {
		nodeGroup        bool
		unownedNodeGroup bool
		authNodeGroup    bool
	}

	type callsCount struct {
		newTasksToDeleteNodeGroups      int
		newTaskToDeleteUnownedNodeGroup int
		deleteNgTaskRuns                int
		deleteUnownedNgTaskRuns         int
		authRemoveNodeGroups            []string
	}

	type deleteNgTest struct {
		nodeGroupNames        []string
		managedNodeGroupNames []string
		updateAuthConfigMap   bool
		stacks                []manager.NodeGroupStack
		errors                deleteErrors

		expectedErr        bool
		expectedCallsCount callsCount
	}

	assertRemoveNodeGroupCalls := func(authConfigMapUpdater *fakes.FakeAuthConfigMapUpdater, ngNames ...string) {
		Expect(authConfigMapUpdater.RemoveNodeGroupCallCount()).To(Equal(len(ngNames)))
		for i, ngName := range ngNames {
			ng := authConfigMapUpdater.RemoveNodeGroupArgsForCall(i)
			Expect(ng.Name).To(Equal(ngName))
		}
	}

	failedErr := errors.New("failed")

	DescribeTable("nodegroups", func(dt deleteNgTest) {
		var (
			nodeGroups        []*api.NodeGroup
			managedNodeGroups []*api.ManagedNodeGroup
		)
		for _, ngName := range dt.nodeGroupNames {
			nodeGroups = append(nodeGroups, &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Name: ngName,
				},
			})
		}
		for _, mngName := range dt.managedNodeGroupNames {
			managedNodeGroups = append(managedNodeGroups, &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Name: mngName,
				},
			})
		}
		var (
			deleteNgTasks    []tasks.Task
			deleteNgTaskRuns atomic.Uint32
		)
		for _, stack := range dt.stacks {
			stack := stack
			deleteNgTasks = append(deleteNgTasks, &taskfakes.FakeTask{
				DescribeStub: func() string {
					return "delete " + stack.NodeGroupName
				},
				DoStub: func(errCh chan error) error {
					defer close(errCh)
					deleteNgTaskRuns.Add(1)
					if dt.errors.nodeGroup {
						return failedErr
					}
					return nil
				},
			})
		}
		var stackHelper fakes.FakeStackHelper
		stackHelper.ListNodeGroupStacksWithStatusesReturns(dt.stacks, nil)
		stackHelper.NewTasksToDeleteNodeGroupsReturns(&tasks.TaskTree{
			Parallel: true,
			Tasks:    deleteNgTasks,
		}, nil)

		var deleteUnownedNgTaskRuns atomic.Uint32
		stackHelper.NewTaskToDeleteUnownedNodeGroupReturns(&taskfakes.FakeTask{
			DescribeStub: func() string {
				return "delete managed nodegroup"
			},
			DoStub: func(errCh chan error) error {
				defer close(errCh)
				deleteUnownedNgTaskRuns.Add(1)
				if dt.errors.unownedNodeGroup {
					return failedErr
				}
				return nil
			},
		})

		var authConfigMapUpdater fakes.FakeAuthConfigMapUpdater
		authConfigMapUpdater.RemoveNodeGroupStub = func(_ *api.NodeGroup) error {
			if dt.errors.authNodeGroup {
				return failedErr
			}
			return nil
		}
		ngDeleter := &nodegroup.Deleter{
			StackHelper:          &stackHelper,
			NodeGroupDeleter:     &managerfakes.FakeNodeGroupDeleter{},
			ClusterName:          "cluster",
			AuthConfigMapUpdater: &authConfigMapUpdater,
		}

		err := ngDeleter.Delete(context.Background(), nodeGroups, managedNodeGroups, nodegroup.DeleteOptions{
			Wait:                true,
			Plan:                false,
			UpdateAuthConfigMap: dt.updateAuthConfigMap,
		})
		if dt.expectedErr {
			Expect(err).To(MatchError(ContainSubstring("failed to delete nodegroup(s)")))
		} else {
			Expect(err).NotTo(HaveOccurred())
		}
		Expect(stackHelper.ListNodeGroupStacksWithStatusesCallCount()).To(Equal(1))
		Expect(deleteNgTaskRuns.Load()).To(Equal(uint32(dt.expectedCallsCount.deleteNgTaskRuns)))
		Expect(deleteUnownedNgTaskRuns.Load()).To(Equal(uint32(dt.expectedCallsCount.deleteUnownedNgTaskRuns)))
		Expect(stackHelper.NewTasksToDeleteNodeGroupsCallCount()).To(Equal(dt.expectedCallsCount.newTasksToDeleteNodeGroups))
		Expect(stackHelper.NewTaskToDeleteUnownedNodeGroupCallCount()).To(Equal(dt.expectedCallsCount.newTaskToDeleteUnownedNodeGroup))
		assertRemoveNodeGroupCalls(&authConfigMapUpdater, dt.expectedCallsCount.authRemoveNodeGroups...)
	},
		Entry("delete self-managed nodegroups", deleteNgTest{
			nodeGroupNames: []string{"ng1", "ng2", "ng3", "ng4"},
			stacks: []manager.NodeGroupStack{
				{
					NodeGroupName: "ng1",
					Type:          api.NodeGroupTypeUnmanaged,
				},
				{
					NodeGroupName:   "ng2",
					Type:            api.NodeGroupTypeUnmanaged,
					UsesAccessEntry: true,
				},
				{
					NodeGroupName: "ng3",
					Type:          api.NodeGroupTypeUnmanaged,
				},
				{
					NodeGroupName:   "ng4",
					Type:            api.NodeGroupTypeUnmanaged,
					UsesAccessEntry: true,
				},
			},
			updateAuthConfigMap: true,

			expectedCallsCount: callsCount{
				newTasksToDeleteNodeGroups: 1,
				deleteNgTaskRuns:           4,
				authRemoveNodeGroups:       []string{"ng1", "ng3"},
			},
		}),

		Entry("delete self-managed and managed nodegroups", deleteNgTest{
			nodeGroupNames:        []string{"ng1", "ng2"},
			managedNodeGroupNames: []string{"mng1", "mng2", "mng3", "mng4"},
			stacks: []manager.NodeGroupStack{
				{
					NodeGroupName: "ng1",
					Type:          api.NodeGroupTypeUnmanaged,
				},
				{
					NodeGroupName:   "ng2",
					Type:            api.NodeGroupTypeUnmanaged,
					UsesAccessEntry: true,
				},
				{
					NodeGroupName: "mng1",
					Type:          api.NodeGroupTypeManaged,
				},
				{
					NodeGroupName: "mng3",
					Type:          api.NodeGroupTypeManaged,
				},
			},
			updateAuthConfigMap: true,

			expectedCallsCount: callsCount{
				newTasksToDeleteNodeGroups:      1,
				newTaskToDeleteUnownedNodeGroup: 2,
				deleteNgTaskRuns:                4,
				deleteUnownedNgTaskRuns:         2,
				authRemoveNodeGroups:            []string{"ng1"},
			},
		}),

		Entry("delete unowned managed nodegroups", deleteNgTest{
			managedNodeGroupNames: []string{"mng1", "mng2"},
			updateAuthConfigMap:   true,

			expectedCallsCount: callsCount{
				newTaskToDeleteUnownedNodeGroup: 2,
				deleteUnownedNgTaskRuns:         2,
			},
		}),

		Entry("delete self-managed nodegroups and skip updating aws-auth ConfigMap", deleteNgTest{
			nodeGroupNames: []string{"ng1"},
			stacks: []manager.NodeGroupStack{
				{
					NodeGroupName: "ng1",
					Type:          api.NodeGroupTypeUnmanaged,
				},
			},

			expectedCallsCount: callsCount{
				newTasksToDeleteNodeGroups: 1,
				deleteNgTaskRuns:           1,
			},
		}),

		Entry("deleting self-managed nodegroups returns an error", deleteNgTest{
			nodeGroupNames: []string{"ng1"},
			stacks: []manager.NodeGroupStack{
				{
					NodeGroupName: "ng1",
					Type:          api.NodeGroupTypeUnmanaged,
				},
			},
			errors: deleteErrors{
				nodeGroup: true,
			},

			expectedCallsCount: callsCount{
				newTasksToDeleteNodeGroups: 1,
				deleteNgTaskRuns:           1,
			},
			expectedErr: true,
		}),

		Entry("deleting unowned nodegroups returns an error", deleteNgTest{
			managedNodeGroupNames: []string{"mng1"},
			stacks:                []manager.NodeGroupStack{},
			errors: deleteErrors{
				unownedNodeGroup: true,
			},

			expectedCallsCount: callsCount{
				newTaskToDeleteUnownedNodeGroup: 1,
				deleteUnownedNgTaskRuns:         1,
			},
			expectedErr: true,
		}),

		Entry("ignore error if removing nodegroup from aws-auth ConfigMap returns an error", deleteNgTest{
			nodeGroupNames:        []string{"ng1"},
			managedNodeGroupNames: []string{"mng1"},
			stacks: []manager.NodeGroupStack{
				{
					NodeGroupName: "ng1",
					Type:          api.NodeGroupTypeUnmanaged,
				},
				{
					NodeGroupName: "mng1",
					Type:          api.NodeGroupTypeManaged,
				},
			},
			errors: deleteErrors{
				authNodeGroup: true,
			},
			updateAuthConfigMap: true,
			expectedCallsCount: callsCount{
				newTasksToDeleteNodeGroups: 1,
				deleteNgTaskRuns:           2,
				authRemoveNodeGroups:       []string{"ng1"},
			},
		}),
	)
})
