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

func NewRemover(clusterConfig *api.ClusterConfig, stackRemover StackRemover, eksAPI awsapi.EKS) *Remover {
	return &Remover{
		clusterName:  clusterConfig.Metadata.Name,
		stackRemover: stackRemover,
		eksAPI:       eksAPI,
	}
}

func (aer *Remover) Delete(ctx context.Context, entries []api.AccessEntry) error {
	tasks := &tasks.TaskTree{Parallel: true}

	stacks, err := aer.stackRemover.ListAccessEntryStackNames(ctx, aer.clusterName)
	if err != nil {
		return fmt.Errorf("listing access entry stacks: %w", err)
	}

	for _, e := range entries {
		stackName := MakeStackName(aer.clusterName, e)
		if !slice.Contains(stacks, stackName) {
			tasks.Append(&deleteUnownedAccessEntryTask{
				info:         fmt.Sprintf("delete access entry for principal ARN %s", e.PrincipalARN),
				clusterName:  aer.clusterName,
				principalARN: e.PrincipalARN,
				eksAPI:       aer.eksAPI,
				ctx:          ctx,
			})
			continue
		}
		tasks.Append(&deleteOwnedAccessEntryTask{
			info:         fmt.Sprintf("delete access entry for principal ARN %s", e.PrincipalARN),
			stackName:    stackName,
			stackRemover: aer.stackRemover,
			principalARN: e.PrincipalARN,
			ctx:          ctx,
		})
	}

	logger.Info(tasks.Describe())
	if errs := tasks.DoAllSync(); len(errs) > 0 {
		return handleErrors(errs, "accessentry(ies)")
	}
	return nil
}

func handleErrors(errs []error, subject string) error {
	logger.Info("%d error(s) occurred while deleting %s", len(errs), subject)
	for _, err := range errs {
		logger.Critical("%s\n", err.Error())
	}
	return fmt.Errorf("failed to delete %s", subject)
}
