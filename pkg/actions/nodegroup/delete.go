package nodegroup

import (
	"fmt"

	awseks "github.com/aws/aws-sdk-go/service/eks"

	"github.com/aws/aws-sdk-go/service/eks/eksiface"

	"github.com/weaveworks/eksctl/pkg/eks"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"

	"github.com/kris-nova/logger"
)

func (ng *Manager) Delete(nodeGroups []*api.NodeGroup, managedNodeGroups []*api.ManagedNodeGroup, wait, plan bool) error {
	var nodeGroupsWithStacks []eks.KubeNodeGroup
	var nodeGroupsWithoutStacksDeleteTasks []*DeleteUnownedNodegroupTask

	for _, n := range nodeGroups {
		nodeGroupsWithStacks = append(nodeGroupsWithStacks, n)
	}

	for _, n := range managedNodeGroups {
		hasStacks, err := ng.hasStacks(n.Name)
		if err != nil {
			return err
		}

		if hasStacks {
			nodeGroupsWithStacks = append(nodeGroupsWithStacks, n)
		} else {
			nodeGroupsWithoutStacksDeleteTasks = append(nodeGroupsWithoutStacksDeleteTasks, &DeleteUnownedNodegroupTask{
				cluster:   ng.cfg.Metadata.Name,
				nodegroup: n.Name,
				eksAPI:    ng.ctl.Provider.EKS(),
				info:      fmt.Sprintf("delete unowned nodegroup %q", n.Name),
			})
		}
	}

	shouldDelete := func(ngName string) bool {
		for _, n := range nodeGroupsWithStacks {
			if n.NameString() == ngName {
				return true
			}
		}
		return false
	}

	tasks, err := ng.manager.NewTasksToDeleteNodeGroups(shouldDelete, wait, nil)
	if err != nil {
		return err
	}

	for _, t := range nodeGroupsWithoutStacksDeleteTasks {
		tasks.Append(t)
	}

	tasks.PlanMode = plan
	logger.Info(tasks.Describe())
	if errs := tasks.DoAllSync(); len(errs) > 0 {
		return handleErrors(errs, "nodegroup(s)")
	}
	return nil
}

func (ng *Manager) hasStacks(name string) (bool, error) {
	stacks, err := ng.manager.ListNodeGroupStacks()
	if err != nil {
		return false, err
	}
	for _, stack := range stacks {
		if stack.NodeGroupName == name {
			return true, nil
		}
	}
	return false, nil
}

type DeleteUnownedNodegroupTask struct {
	info      string
	eksAPI    eksiface.EKSAPI
	cluster   string
	nodegroup string
}

func (d *DeleteUnownedNodegroupTask) Describe() string {
	return d.info
}

func (d *DeleteUnownedNodegroupTask) Do(errorchan chan error) error {
	out, err := d.eksAPI.DeleteNodegroup(&awseks.DeleteNodegroupInput{
		ClusterName:   &d.cluster,
		NodegroupName: &d.nodegroup,
	})
	go func() {
		errorchan <- err
	}()

	if out != nil {
		logger.Debug("Delete nodegroup %q output: %s", d.nodegroup, out.String())
	}
	return err
}

func handleErrors(errs []error, subject string) error {
	logger.Info("%d error(s) occurred while deleting %s", len(errs), subject)
	for _, err := range errs {
		logger.Critical("%s\n", err.Error())
	}
	return fmt.Errorf("failed to delete %s", subject)
}
