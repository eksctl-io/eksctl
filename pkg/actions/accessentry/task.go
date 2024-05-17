package accessentry

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"

	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_stack_creator.go . StackCreator
type StackCreator interface {
	CreateStack(ctx context.Context, stackName string, resourceSet builder.ResourceSetReader, tags, parameters map[string]string, errs chan error) error
	TroubleshootStackFailureCause(ctx context.Context, stack *cfntypes.Stack, desiredStatus cfntypes.StackStatus)
}

//counterfeiter:generate -o fakes/fake_stack_remover.go . StackRemover
type StackRemover interface {
	DescribeStack(ctx context.Context, s *cfntypes.Stack) (*cfntypes.Stack, error)
	DeleteStackBySpecSync(ctx context.Context, s *cfntypes.Stack, errs chan error) error
	ListAccessEntryStackNames(ctx context.Context, clusterName string) ([]string, error)
}

type accessEntryTask struct {
	info         string
	clusterName  string
	accessEntry  api.AccessEntry
	stackCreator StackCreator
	ctx          context.Context
}

func (t *accessEntryTask) Describe() string { return t.info }

func (t *accessEntryTask) Do(errorCh chan error) error {
	defer close(errorCh)
	rs := builder.NewAccessEntryResourceSet(t.clusterName, t.accessEntry)
	if err := rs.AddAllResources(); err != nil {
		return err
	}
	principalARN := t.accessEntry.PrincipalARN.String()
	logger.Info("creating access entry for principal ARN %q", principalARN)
	stackErrCh := make(chan error)
	stackName := MakeStackName(t.clusterName, t.accessEntry)
	if err := t.stackCreator.CreateStack(t.ctx, stackName, rs, nil, nil, stackErrCh); err != nil {
		return err
	}
	select {
	case err := <-stackErrCh:
		if err != nil {
			return t.troubleshootFailure(stackName, err)
		}
		logger.Info("created access entry for principal ARN %q", principalARN)
		return nil
	case <-t.ctx.Done():
		return fmt.Errorf("timed out waiting for access entry %q: %w", principalARN, t.ctx.Err())
	}
}

func (t *accessEntryTask) troubleshootFailure(stackName string, err error) error {
	stack := &cfntypes.Stack{
		StackName: aws.String(stackName),
	}
	t.stackCreator.TroubleshootStackFailureCause(t.ctx, stack, cfntypes.StackStatusCreateComplete)
	if strings.Contains(err.Error(), "waiter state transitioned to Failure") {
		return fmt.Errorf("failed to create access entry for principal ARN %q", t.accessEntry.PrincipalARN.String())
	}
	return err
}

type deleteUnownedAccessEntryTask struct {
	info         string
	clusterName  string
	principalARN api.ARN
	eksAPI       awsapi.EKS
	ctx          context.Context
}

func (t *deleteUnownedAccessEntryTask) Describe() string {
	return t.info
}

func (t *deleteUnownedAccessEntryTask) Do(errorCh chan error) error {
	defer close(errorCh)

	if _, err := t.eksAPI.DeleteAccessEntry(t.ctx, &awseks.DeleteAccessEntryInput{
		ClusterName:  &t.clusterName,
		PrincipalArn: aws.String(t.principalARN.String()),
	}); err != nil {
		return fmt.Errorf("deleting access entry with principalARN %s: %w", t.principalARN, err)
	}
	logger.Info("started deleting access entry with principalARN %q", t.principalARN)

	return nil
}

type deleteOwnedAccessEntryTask struct {
	info         string
	stackName    string
	stackRemover StackRemover
	principalARN api.ARN
	ctx          context.Context
}

func (t *deleteOwnedAccessEntryTask) Describe() string {
	return t.info
}

func (t *deleteOwnedAccessEntryTask) Do(errorCh chan error) error {
	stack, err := t.stackRemover.DescribeStack(t.ctx, &cfntypes.Stack{StackName: &t.stackName})
	if err != nil {
		// the stack should not be missing as we retrieved its name at a previous point
		return fmt.Errorf("failed to describe stack for access entry with principal ARN %s: %w", t.principalARN.String(), err)
	}

	if err := t.stackRemover.DeleteStackBySpecSync(t.ctx, stack, errorCh); err != nil {
		return fmt.Errorf("deleting access entry with principalARN %s: %w", t.principalARN.String(), err)
	}

	return nil
}

func runAllTasks(taskTree *tasks.TaskTree) error {
	logger.Info(taskTree.Describe())
	if errs := taskTree.DoAllSync(); len(errs) > 0 {
		var allErrs []string
		for _, err := range errs {
			allErrs = append(allErrs, err.Error())
		}
		return fmt.Errorf(strings.Join(allErrs, "\n"))
	}
	completedAction := func() string {
		if taskTree.PlanMode {
			return "skipped"
		}
		return "completed successfully"
	}
	logger.Info("all tasks were %s", completedAction())
	return nil
}
