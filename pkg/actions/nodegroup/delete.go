package nodegroup

import (
	"context"
	"fmt"

	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

// StackHelper is a helper for managing nodegroup stacks.
//
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_stack_helper.go . StackHelper
type StackHelper interface {
	NewTasksToDeleteNodeGroups(stacks []manager.NodeGroupStack, shouldDelete func(_ string) bool, wait bool, cleanup func(chan error, string) error) (*tasks.TaskTree, error)
	ListNodeGroupStacksWithStatuses(ctx context.Context) ([]manager.NodeGroupStack, error)
	NewTaskToDeleteUnownedNodeGroup(ctx context.Context, clusterName, nodegroup string, nodeGroupDeleter manager.NodeGroupDeleter, waitCondition *manager.DeleteWaitCondition) tasks.Task
}

// AuthConfigMapUpdater updates the aws-auth ConfigMap.
//
//counterfeiter:generate -o fakes/fake_auth_configmap_updater.go . AuthConfigMapUpdater
type AuthConfigMapUpdater interface {
	// RemoveNodeGroup removes the specified nodegroup.
	RemoveNodeGroup(*api.NodeGroup) error
}

// A Deleter deletes nodegroups.
type Deleter struct {
	StackHelper          StackHelper
	NodeGroupDeleter     manager.NodeGroupDeleter
	ClusterName          string
	AuthConfigMapUpdater AuthConfigMapUpdater
}

// DeleteOptions represents the options for deleting nodegroups.
type DeleteOptions struct {
	Wait                bool
	Plan                bool
	UpdateAuthConfigMap bool
}

// Delete deletes the specified nodegroups.
func (d *Deleter) Delete(ctx context.Context, nodeGroups []*api.NodeGroup, managedNodeGroups []*api.ManagedNodeGroup, options DeleteOptions) error {
	nodeGroupsWithStacks := map[string]struct{}{}
	for _, n := range nodeGroups {
		nodeGroupsWithStacks[n.NameString()] = struct{}{}
	}

	stacks, err := d.StackHelper.ListNodeGroupStacksWithStatuses(ctx)
	if err != nil {
		return err
	}

	taskTree := &tasks.TaskTree{
		Parallel: true,
		PlanMode: options.Plan,
	}
	for _, n := range managedNodeGroups {
		if findStack(stacks, n.Name) != nil {
			nodeGroupsWithStacks[n.NameString()] = struct{}{}
		} else {
			taskTree.Append(d.StackHelper.NewTaskToDeleteUnownedNodeGroup(ctx, d.ClusterName, n.Name, d.NodeGroupDeleter, nil))
		}
	}

	var deleteTasks tasks.Task
	if len(stacks) > 0 {
		shouldDelete := func(ngName string) bool {
			_, ok := nodeGroupsWithStacks[ngName]
			return ok
		}
		deleteTasks, err = d.StackHelper.NewTasksToDeleteNodeGroups(stacks, shouldDelete, options.Wait, nil)
		if err != nil {
			return err
		}
	}
	if authTask := d.updateAuthConfigMapTask(nodeGroups, stacks, options); authTask != nil {
		if deleteTasks != nil {
			var subTasks tasks.TaskTree
			subTasks.Append(
				deleteTasks,
				authTask,
			)
			deleteTasks = &subTasks
		} else {
			deleteTasks = authTask
		}
	}

	if deleteTasks != nil {
		taskTree.Append(deleteTasks)
	}

	logger.Info(taskTree.Describe())
	if errs := taskTree.DoAllSync(); len(errs) > 0 {
		return handleErrors(errs, "nodegroup(s)")
	}
	return nil
}

func (d *Deleter) updateAuthConfigMapTask(nodeGroups []*api.NodeGroup, stacks []manager.NodeGroupStack, options DeleteOptions) tasks.Task {
	if !options.UpdateAuthConfigMap {
		return nil
	}
	var nodeGroupsWithoutAccessEntry []*api.NodeGroup
	for _, ng := range nodeGroups {
		if stack := findStack(stacks, ng.NameString()); stack == nil || !stack.UsesAccessEntry {
			nodeGroupsWithoutAccessEntry = append(nodeGroupsWithoutAccessEntry, ng)
		}
	}
	if len(nodeGroupsWithoutAccessEntry) == 0 {
		return nil
	}

	return &tasks.GenericTask{
		Description: "update auth ConfigMap",
		Doer: func() error {
			cmdutils.LogIntendedAction(options.Plan, "delete %d nodegroups from auth ConfigMap in cluster %q", len(nodeGroupsWithoutAccessEntry), d.ClusterName)
			if options.Plan {
				return nil
			}
			for _, ng := range nodeGroupsWithoutAccessEntry {
				if err := d.AuthConfigMapUpdater.RemoveNodeGroup(ng); err != nil {
					logger.Warning(err.Error())
				}
			}
			return nil
		},
	}
}

func handleErrors(errs []error, subject string) error {
	logger.Info("%d error(s) occurred while deleting %s", len(errs), subject)
	for _, err := range errs {
		logger.Critical("%s\n", err.Error())
	}
	return fmt.Errorf("failed to delete %s", subject)
}
