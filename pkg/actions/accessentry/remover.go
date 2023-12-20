package accessentry

import (
	"context"
	"fmt"

	"github.com/bxcodec/faker/support/slice"
	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

type Remover struct {
	clusterName  string
	stackRemover StackRemover
	eksAPI       awsapi.EKS
}

func NewRemover(clusterName string, stackRemover StackRemover, eksAPI awsapi.EKS) *Remover {
	return &Remover{
		clusterName:  clusterName,
		stackRemover: stackRemover,
		eksAPI:       eksAPI,
	}
}

func (aer *Remover) Delete(ctx context.Context, accessEntries []api.AccessEntry) error {
	tasks, err := aer.DeleteTasks(ctx, accessEntries)
	if err != nil {
		return err
	}

	logger.Info(tasks.Describe())
	if errs := tasks.DoAllSync(); len(errs) > 0 {
		return handleErrors(errs, "accessentry(ies)")
	}
	return nil
}

func (aer Remover) DeleteTasks(ctx context.Context, accessEntries []api.AccessEntry) (*tasks.TaskTree, error) {
	stacks, err := aer.stackRemover.ListAccessEntryStackNames(ctx, aer.clusterName)
	if err != nil {
		return nil, fmt.Errorf("listing access entry stacks: %w", err)
	}

	tasks := &tasks.TaskTree{
		Parallel: true,
	}

	// this is true during cluster deletion, when no access entry is given as user input.
	if len(accessEntries) == 0 {
		for _, s := range stacks {
			tasks.Append(&deleteOwnedAccessEntryTask{
				info:         fmt.Sprintf("delete access entry stack %q", s),
				stackName:    s,
				stackRemover: aer.stackRemover,
				principalARN: api.ARN{Partition: "unknown"},
				ctx:          ctx,
			})
		}
		return tasks, nil
	}

	for _, e := range accessEntries {
		stackName := MakeStackName(aer.clusterName, e)
		if !slice.Contains(stacks, stackName) {
			tasks.Append(&deleteUnownedAccessEntryTask{
				info:         fmt.Sprintf("delete access entry for principal ARN %s", e.PrincipalARN.String()),
				clusterName:  aer.clusterName,
				principalARN: e.PrincipalARN,
				eksAPI:       aer.eksAPI,
				ctx:          ctx,
			})
			continue
		}
		tasks.Append(&deleteOwnedAccessEntryTask{
			info:         fmt.Sprintf("delete access entry for principal ARN %s", e.PrincipalARN.String()),
			stackName:    stackName,
			stackRemover: aer.stackRemover,
			principalARN: e.PrincipalARN,
			ctx:          ctx,
		})
	}

	return tasks, nil
}

func handleErrors(errs []error, subject string) error {
	logger.Info("%d error(s) occurred while deleting %s", len(errs), subject)
	for _, err := range errs {
		logger.Critical("%s\n", err.Error())
	}
	return fmt.Errorf("failed to delete %s", subject)
}
