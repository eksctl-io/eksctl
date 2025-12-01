package capability

import (
	"context"
	"errors"
	"fmt"
	"strings"

	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"

	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

type RemoverInterface interface {
	Delete(ctx context.Context, capabilities []Summary) error
	DeleteTasks(ctx context.Context, capabilities []Summary) *tasks.TaskTree
}

type StackRemover interface {
	DescribeStack(ctx context.Context, s *cfntypes.Stack) (*cfntypes.Stack, error)
	DeleteStackSync(ctx context.Context, s *cfntypes.Stack) error
	ListCapabilityStacks(ctx context.Context) ([]*cfntypes.Stack, error)
	ListCapabilitiesIAMStacks(ctx context.Context) ([]*cfntypes.Stack, error)
}

type Remover struct {
	clusterName  string
	stackRemover StackRemover
}

func NewRemover(clusterName string, stackRemover StackRemover) *Remover {
	return &Remover{
		clusterName:  clusterName,
		stackRemover: stackRemover,
	}
}

func (r *Remover) Delete(ctx context.Context, capabilities []Summary) error {
	// Create parallel tasks for capability deletion and waiting
	capabilityTasks, err := r.DeleteTasks(ctx, capabilities)

	if err != nil {
		return err
	}

	// Execute capability deletions in parallel
	if errs := capabilityTasks.DoAllSync(); len(errs) > 0 {
		var allErrs []string
		for _, err := range errs {
			allErrs = append(allErrs, err.Error())
		}
		return errors.New(strings.Join(allErrs, "\n"))
	}

	return nil
}

func (r *Remover) DeleteTasks(ctx context.Context, capabilities []Summary) (*tasks.TaskTree, error) {

	// Get CapabilityIAMRole stacks
	capabilityStacks, err := r.stackRemover.ListCapabilityStacks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stacks for capabilities: %w", err)
	}

	// Get CapabilityIAMRole stacks
	iamStacks, err := r.stackRemover.ListCapabilitiesIAMStacks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch IAM role stacks for capabilities: %w", err)
	}

	// Create parallel tasks for capability deletion and waiting
	capabilityTasks := &tasks.TaskTree{
		Parallel: true,
	}

	for _, cap := range capabilities {
		capabilityTasks.Append(&tasks.GenericTask{
			Description: fmt.Sprintf("delete and wait for capability %s", cap.Name),
			Doer: func() error {
				if err := r.deleteCapabilityStack(ctx, cap.Name, capabilityStacks); err != nil {
					return err
				}

				return r.deleteCapabilityIAMRoleStack(ctx, cap.Name, iamStacks)
			},
		})
	}

	return capabilityTasks, nil
}

// deleteCapabilityStack deletes capability stacks filtered by capability name
func (r *Remover) deleteCapabilityStack(ctx context.Context, capabilityName string, stacks []*cfntypes.Stack) error {
	for _, stackItem := range stacks {
		if manager.GetCapabilityNameFromStack(stackItem) != capabilityName {
			continue
		}
		stack, err := r.stackRemover.DescribeStack(ctx, &cfntypes.Stack{StackName: stackItem.StackName})
		if err != nil {
			// the stack should not be missing as we retrieved its name previously
			return fmt.Errorf("failed to describe stack for capability %s: %w", capabilityName, err)
		}

		if err := r.stackRemover.DeleteStackSync(ctx, stack); err != nil {
			return fmt.Errorf("deleting capability IAM role stack of capbility %s: %w", *stack.StackName, err)
		}
	}

	return nil
}

// deleteCapabilityIAMTasksFiltered deletes capability IAM stacks filtered by capability name
func (r *Remover) deleteCapabilityIAMRoleStack(ctx context.Context, capabilityName string, stacks []*cfntypes.Stack) error {
	for _, stackItem := range stacks {
		if manager.GetCapabilityNameFromIAMStack(stackItem) != capabilityName {
			continue
		}
		stack, err := r.stackRemover.DescribeStack(ctx, &cfntypes.Stack{StackName: stackItem.StackName})
		if err != nil {
			// the stack should not be missing as we retrieved its name previously
			return fmt.Errorf("failed to describe stack for IAM role of capability %s: %w", capabilityName, err)
		}

		if err := r.stackRemover.DeleteStackSync(ctx, stack); err != nil {
			return fmt.Errorf("deleting capability IAM role stack of capbility %s: %w", *stack.StackName, err)
		}
	}

	return nil
}
