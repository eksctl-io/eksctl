package nodegroup

import (
	"fmt"

	awseks "github.com/aws/aws-sdk-go/service/eks"

	"github.com/aws/aws-sdk-go/service/eks/eksiface"

	"github.com/weaveworks/eksctl/pkg/eks"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"

	"github.com/kris-nova/logger"
)

func (ng *NodeGroupManager) Delete(nodeGroups []*api.NodeGroup, managedNodeGroups []*api.ManagedNodeGroup, wait, plan bool) error {
	plan = true
	var nodesWithStacks []eks.KubeNodeGroup
	var nodesWithoutStacksDeleteTask []*DeleteUnownedNodegroupTask

	for _, n := range nodeGroups {
		nodesWithStacks = append(nodesWithStacks, n)
	}

	for _, n := range managedNodeGroups {
		hasStacks, err := ng.hasStacks(n)
		if err != nil {
			return err
		}

		if hasStacks {
			nodesWithStacks = append(nodesWithStacks, n)
		} else {
			nodesWithoutStacksDeleteTask = append(nodesWithoutStacksDeleteTask, &DeleteUnownedNodegroupTask{
				cluster:   ng.cfg.Metadata.Name,
				nodegroup: n.Name,
				eksAPI:    ng.ctl.Provider.EKS(),
				info:      fmt.Sprintf("delete unowned nodegroup %q", n.Name),
			})
		}
	}

	shouldDelete := func(ngName string) bool {
		for _, n := range nodesWithStacks {
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

	for _, t := range nodesWithoutStacksDeleteTask {
		tasks.Append(t)
	}

	tasks.PlanMode = plan
	logger.Info(tasks.Describe())
	if errs := tasks.DoAllSync(); len(errs) > 0 {
		return handleErrors(errs, "nodegroup(s)")
	}
	return nil
}

func (ng *NodeGroupManager) hasStacks(n *api.ManagedNodeGroup) (bool, error) {
	stacks, err := ng.manager.ListNodeGroupStacks()
	if err != nil {
		return false, err
	}
	for _, stack := range stacks {
		if stack.NodeGroupName == n.Name {
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
