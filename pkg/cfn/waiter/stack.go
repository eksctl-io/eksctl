package waiter

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/kris-nova/logger"

	"github.com/weaveworks/eksctl/pkg/awsapi"
)

var ClusterCreationNextDelay = func(attempts int) time.Duration {
	// Wait 30s for the first two requests, and 1m for subsequent requests.
	if attempts <= 2 {
		return 30 * time.Second
	}
	return 1 * time.Minute
}

// NextDelay returns the amount of time to wait before the next retry given the number of attempts.
type NextDelay func(attempts int) time.Duration

// WaitForStack waits for the cluster stack to reach a success or failure state, and returns the stack.
func WaitForStack(ctx context.Context, cfnAPI awsapi.CloudFormation, stackID, stackName string, nextDelay NextDelay) (*types.Stack, error) {
	var lastStack *types.Stack
	waiter := &Waiter{
		NextDelay: nextDelay,
		Operation: func() (bool, error) {
			var (
				err     error
				success bool
			)
			lastStack, success, err = describeStackStatus(context.Background(), cfnAPI, stackID, stackName)
			return success, err
		},
	}

	if err := waiter.Wait(ctx); err != nil {
		return nil, err
	}

	if lastStack == nil {
		return nil, errors.New("unexpected nil value for stack")
	}

	return lastStack, nil
}

func describeStackStatus(ctx context.Context, cfnAPI awsapi.CloudFormation, stackID, stackName string) (*types.Stack, bool, error) {
	logger.Info("waiting for CloudFormation stack %q", stackName)
	output, err := cfnAPI.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{
		StackName: aws.String(stackID),
	})
	if err != nil {
		return nil, false, fmt.Errorf("error describing stack: %w", err)
	}
	if len(output.Stacks) != 1 {
		return nil, false, fmt.Errorf("expected a single stack; got %d", len(output.Stacks))
	}

	switch stack := output.Stacks[0]; stack.StackStatus {
	case types.StackStatusCreateComplete,
		types.StackStatusUpdateComplete:
		return &stack, true, nil

	case types.StackStatusCreateFailed,
		types.StackStatusRollbackInProgress,
		types.StackStatusRollbackFailed,
		types.StackStatusRollbackComplete,
		types.StackStatusDeleteInProgress,
		types.StackStatusDeleteFailed,
		types.StackStatusDeleteComplete:
		return &stack, false, errors.New("ResourceNotReady: failed waiting for successful resource state")

	default:
		return &stack, false, nil
	}
}
