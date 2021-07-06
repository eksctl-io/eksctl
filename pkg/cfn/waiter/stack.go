package waiter

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
)

// NextDelay returns the amount of time to wait before the next retry given the number of attempts.
type NextDelay func(attempts int) time.Duration

// WaitForStack waits for the cluster stack to reach a success or failure state, and returns the stack.
func WaitForStack(ctx context.Context, cfnAPI cloudformationiface.CloudFormationAPI, stackID, stackName string, nextDelay NextDelay) (*cloudformation.Stack, error) {
	var lastStack *cloudformation.Stack
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

func describeStackStatus(ctx context.Context, cfnAPI cloudformationiface.CloudFormationAPI, stackID, stackName string) (*cloudformation.Stack, bool, error) {
	logger.Info("waiting for CloudFormation stack %q", stackName)
	req, output := cfnAPI.DescribeStacksRequest(&cloudformation.DescribeStacksInput{
		StackName: aws.String(stackID),
	})
	req.SetContext(ctx)
	if err := req.Send(); err != nil {
		return nil, false, errors.Wrap(err, "error describing stack")
	}
	if len(output.Stacks) != 1 {
		return nil, false, errors.Errorf("expected a single stack; got %d", len(output.Stacks))
	}

	switch stack := output.Stacks[0]; *stack.StackStatus {
	case cloudformation.StackStatusCreateComplete,
		cloudformation.StackStatusUpdateComplete:
		return stack, true, nil

	case cloudformation.StackStatusCreateFailed,
		cloudformation.StackStatusRollbackInProgress,
		cloudformation.StackStatusRollbackFailed,
		cloudformation.StackStatusRollbackComplete,
		cloudformation.StackStatusDeleteInProgress,
		cloudformation.StackStatusDeleteFailed,
		cloudformation.StackStatusDeleteComplete:
		return stack, false, errors.New("ResourceNotReady: failed waiting for successful resource state")

	default:
		return stack, false, nil
	}
}
