package accessentry

import (
	"context"
	"errors"
	"fmt"
	"strings"

	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/bxcodec/faker/support/slice"
	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

// A StackManager manages CloudFormation stacks for access entries.
type StackManager interface {
	ListAccessEntryStackNames(ctx context.Context, clusterName string) ([]string, error)
	DescribeStack(ctx context.Context, s *cfntypes.Stack) (*cfntypes.Stack, error)
	DeleteStackBySpec(ctx context.Context, s *cfntypes.Stack) (*cfntypes.Stack, error)
}

type Remover struct {
	clusterName  string
	stackManager StackManager
	eksAPI       awsapi.EKS
}

func NewRemover(clusterConfig *api.ClusterConfig, stackManager StackManager, eksAPI awsapi.EKS) *Remover {
	return &Remover{
		clusterName:  clusterConfig.Metadata.Name,
		stackManager: stackManager,
		eksAPI:       eksAPI,
	}
}

func (aer *Remover) Delete(ctx context.Context, entries []api.AccessEntry) error {
	tasks := &tasks.TaskTree{Parallel: true}

	// may replace this call with something dedicated to figure out authenticationMode
	_, err := aer.eksAPI.ListAccessEntries(ctx, &eks.ListAccessEntriesInput{
		ClusterName: &aer.clusterName,
	})
	if err != nil {
		var invalidRequestErr *ekstypes.InvalidRequestException
		if errors.As(err, &invalidRequestErr) && strings.Contains(err.Error(), AuthenticationModeErr) {
			return ErrDisabledAccessEntryAPI
		}
		return fmt.Errorf("calling EKS API to list access entries: %w", err)
	}

	stacks, err := aer.stackManager.ListAccessEntryStackNames(ctx, aer.clusterName)
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
			stackManager: aer.stackManager,
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
