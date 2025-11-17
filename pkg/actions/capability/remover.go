package capability

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/kris-nova/logger"

	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

type RemoverInterface interface {
	Delete(ctx context.Context, capabilities []CapabilitySummary) error
	DeleteTasks(ctx context.Context, capabilities []CapabilitySummary) *tasks.TaskTree
}

type StackRemover interface {
	DescribeStack(ctx context.Context, s *cfntypes.Stack) (*cfntypes.Stack, error)
	DeleteStackBySpec(ctx context.Context, s *cfntypes.Stack) (*cfntypes.Stack, error)
	DeleteStackBySpecSync(ctx context.Context, s *cfntypes.Stack, errs chan error) error
	GetIAMCapabilitiesStacks(ctx context.Context) ([]*cfntypes.Stack, error)
}

type Remover struct {
	clusterName  string
	stackRemover StackRemover
	eksAPI       awsapi.EKS
	waitTimeout  time.Duration
}

func NewRemover(clusterName string, stackRemover StackRemover, eksAPI awsapi.EKS, waitTimeout time.Duration) *Remover {
	return &Remover{
		clusterName:  clusterName,
		stackRemover: stackRemover,
		eksAPI:       eksAPI,
		waitTimeout:  waitTimeout,
	}
}

func (r *Remover) Delete(ctx context.Context, capabilities []CapabilitySummary) error {
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

func (r *Remover) DeleteTasks(ctx context.Context, capabilities []CapabilitySummary) (*tasks.TaskTree, error) {

	// Get CapabilityIAMRole stacks
	stacks, err := r.stackRemover.GetIAMCapabilitiesStacks(ctx)
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
				if err := r.deleteAndWaitForCapability(ctx, cap.Name); err != nil {
					return err
				}

				return r.deleteCapabilityIAMRoleStack(ctx, cap.Name, stacks)
			},
		})
	}

	return capabilityTasks, nil
}

func (r *Remover) deleteAndWaitForCapability(ctx context.Context, capabilityName string) error {
	// Delete capability via EKS API
	logger.Info("deleting capability %s", capabilityName)
	_, err := r.eksAPI.DeleteCapability(ctx, &awseks.DeleteCapabilityInput{
		ClusterName:    aws.String(r.clusterName),
		CapabilityName: aws.String(capabilityName),
	})
	if err != nil {
		var notFoundErr *ekstypes.ResourceNotFoundException
		if errors.As(err, &notFoundErr) {
			logger.Info("capability %s not found, skipping deletion", capabilityName)
			return nil
		}
		return fmt.Errorf("failed to delete capability %s: %w", capabilityName, err)
	}

	// Wait for capability to be deleted
	return r.waitForCapabilityDeletion(ctx, capabilityName)
}

func (r *Remover) waitForCapabilityDeletion(ctx context.Context, capabilityName string) error {
	timeout := r.waitTimeout
	if timeout == 0 {
		timeout = 15 * time.Minute // Default timeout
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	pollInterval := 15 * time.Second
	if timeout < time.Minute {
		// For short timeouts (like tests), use shorter poll interval
		pollInterval = 100 * time.Millisecond
	}

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	logger.Info("waiting for capability %s to be deleted", capabilityName)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for capability %s to be deleted", capabilityName)
		case <-ticker.C:
			_, err := r.eksAPI.DescribeCapability(ctx, &awseks.DescribeCapabilityInput{
				ClusterName:    aws.String(r.clusterName),
				CapabilityName: &capabilityName,
			})
			if err != nil {
				// Check if capability is not found (successfully deleted)
				var notFoundErr *ekstypes.ResourceNotFoundException
				if errors.As(err, &notFoundErr) {
					logger.Success("capability %s successfully deleted", capabilityName)
					return nil
				}
				return fmt.Errorf("error checking capability %s status: %w", capabilityName, err)
			}
			// Capability still exists, continue waiting
			logger.Debug("capability %s still exists, continuing to wait", capabilityName)
		}
	}
}

// DeleteCapabilityIAMTasksFiltered deletes capability IAM stacks filtered by capability name
func (r *Remover) deleteCapabilityIAMRoleStack(ctx context.Context, capabilityName string, stacks []*cfntypes.Stack) error {
	for _, stackItem := range stacks {
		if manager.GetIAMCapabilityName(stackItem) != capabilityName {
			continue
		}
		stack, err := r.stackRemover.DescribeStack(ctx, &cfntypes.Stack{StackName: stackItem.StackName})
		if err != nil {
			// the stack should not be missing as we retrieved its name previously
			return fmt.Errorf("failed to describe stack for IAM role of capability %s: %w", capabilityName, err)
		}

		errCh := make(chan error, 1)
		if err := r.stackRemover.DeleteStackBySpecSync(ctx, stack, errCh); err != nil {
			return fmt.Errorf("deleting capability IAM role stack of capbility %s: %w", *stack.StackName, err)
		}
	}

	return nil
}
