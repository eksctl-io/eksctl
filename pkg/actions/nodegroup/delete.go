package nodegroup

import (
	"fmt"

	awseks "github.com/aws/aws-sdk-go/service/eks"

	"github.com/aws/aws-sdk-go/service/eks/eksiface"

	"github.com/weaveworks/eksctl/pkg/eks"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"

	"github.com/kris-nova/logger"
)

func (ng *NodeGroup) Delete(nodeGroups []*api.NodeGroup, managedNodeGroups []*api.ManagedNodeGroup, wait, plan bool) error {
	var cloudformationNodegroups []eks.KubeNodeGroup
	var deleteUnownedNodegroupTasks []*DeleteUnownedNodegroupTask

	for _, n := range nodeGroups {
		cloudformationNodegroups = append(cloudformationNodegroups, n)
	}

	for _, n := range managedNodeGroups {
		if n.Unowned {
			deleteUnownedNodegroupTasks = append(deleteUnownedNodegroupTasks, &DeleteUnownedNodegroupTask{
				cluster:   ng.cfg.Metadata.Name,
				nodegroup: n.Name,
				eksAPI:    ng.ctl.Provider.EKS(),
				info:      fmt.Sprintf("delete unowned nodegroup %q", n.Name),
			})
		} else {
			cloudformationNodegroups = append(cloudformationNodegroups, n)
		}
	}

	shouldDelete := func(ngName string) bool {
		for _, n := range cloudformationNodegroups {
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

	for _, t := range deleteUnownedNodegroupTasks {
		tasks.Append(t)
	}

	tasks.PlanMode = plan
	logger.Info(tasks.Describe())
	if errs := tasks.DoAllSync(); len(errs) > 0 {
		return handleErrors(errs, "nodegroup(s)")
	}
	return nil
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
