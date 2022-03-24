package manager

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/kris-nova/logger"

	"github.com/weaveworks/eksctl/pkg/cfn/builder"
)

const (
	stackStatus     = "Stacks[].StackStatus"
	changesetStatus = "Status"
)

func (c *StackCollection) troubleshootStackFailureCause(ctx context.Context, i *Stack, desiredStatus string) {
	logger.Info("fetching stack events in attempt to troubleshoot the root cause of the failure")
	events, err := c.DescribeStackEvents(ctx, i)
	if err != nil {
		logger.Critical("cannot fetch stack events: %v", err)
		return
	}
	for _, e := range events {
		msg := fmt.Sprintf("%s/%s: %s", *e.ResourceType, *e.LogicalResourceId, e.ResourceStatus)
		if e.ResourceStatusReason != nil {
			msg = fmt.Sprintf("%s – %#v", msg, *e.ResourceStatusReason)
		}
		switch desiredStatus {
		case string(types.StackStatusCreateComplete):
			switch string(e.ResourceStatus) {
			case string(types.ResourceStatusCreateFailed):
				logger.Critical(msg)
			case string(types.ResourceStatusDeleteInProgress):
				logger.Warning(msg)
			default:
				logger.Debug(msg) // only output this when verbose logging is enabled
			}
		case string(types.StackStatusDeleteComplete):
			switch string(e.ResourceStatus) {
			case string(types.ResourceStatusDeleteFailed):
				logger.Critical(msg)
			case string(types.ResourceStatusDeleteSkipped):
				logger.Warning(msg)
			default:
				logger.Debug(msg) // only output this when verbose logging is enabled
			}
		default:
			logger.Info(msg)
		}
	}
}

type noChangeError struct {
	msg string
}

func (e *noChangeError) Error() string {
	return e.msg
}

// DoWaitUntilStackIsCreated blocks until the given stack's
// creation has completed.
func (c *StackCollection) DoWaitUntilStackIsCreated(i *Stack) error {
	return nil
	// return c.waitWithAcceptors(i,
	// 	waiters.MakeAcceptors(
	// 		stackStatus,
	// 		types.StackStatusCreateComplete,
	// 		[]string{
	// 			string(types.StackStatusCreateFailed),
	// 			string(types.StackStatusRollbackInProgress),
	// 			string(types.StackStatusRollbackFailed),
	// 			string(types.StackStatusRollbackComplete),
	// 			string(types.StackStatusDeleteInProgress),
	// 			string(types.StackStatusDeleteFailed),
	// 			string(types.StackStatusDeleteComplete),
	// 		},
	// 		request.WaiterAcceptor{
	// 			State:    request.FailureWaiterState,
	// 			Matcher:  request.ErrorWaiterMatch,
	// 			Expected: "ValidationError",
	// 		},
	// 	),
	// )
}

func (c *StackCollection) waitUntilStackIsCreated(ctx context.Context, i *Stack, stack builder.ResourceSet, errs chan error) {
	// defer close(errs)

	// if err := c.DoWaitUntilStackIsCreated(ctx, i); err != nil {
	// 	errs <- err
	// 	return
	// }
	// s, err := c.DescribeStack(ctx, i)
	// if err != nil {
	// 	errs <- err
	// 	return
	// }
	// if err := stack.GetAllOutputs(*s); err != nil {
	// 	errs <- errors.Wrapf(err, "getting stack %q outputs", *i.StackName)
	// 	return
	// }
	// errs <- nil
	// }

	// func (c *StackCollection) doWaitUntilStackIsDeleted(i *Stack) error {
	// return c.waitWithAcceptors(i,
	// 	waiters.MakeAcceptors(
	// 		stackStatus,
	// 		types.StackStatusDeleteComplete,
	// 		[]string{
	// 			string(types.StackStatusDeleteFailed),
	// 			string(types.StackStatusCreateInProgress),
	// 			string(types.StackStatusCreateFailed),
	// 			string(types.StackStatusCreateComplete),
	// 			string(types.StackStatusRollbackInProgress),
	// 			string(types.StackStatusRollbackFailed),
	// 			string(types.StackStatusRollbackComplete),
	// 			string(types.StackStatusUpdateInProgress),
	// 			string(types.StackStatusUpdateCompleteCleanupInProgress),
	// 			string(types.StackStatusUpdateComplete),
	// 			string(types.StackStatusUpdateRollbackInProgress),
	// 			string(types.StackStatusUpdateRollbackFailed),
	// 			string(types.StackStatusUpdateRollbackCompleteCleanupInProgress),
	// 			string(types.StackStatusUpdateRollbackComplete),
	// 			string(types.StackStatusReviewInProgress),
	// 		},
	// 		// ValidationError is expected as success, although
	// 		// we use stack ARN, so should normally see actual
	// 		// StackStatusDeleteComplete
	// 		request.WaiterAcceptor{
	// 			State:    request.SuccessWaiterState,
	// 			Matcher:  request.ErrorWaiterMatch,
	// 			Expected: "ValidationError",
	// 		},
	// 	),
	// )
}

func (c *StackCollection) doWaitUntilStackIsDeleted(i *Stack) error {
	return nil
	// return c.waitWithAcceptors(i,
	// 	waiters.MakeAcceptors(
	// 		stackStatus,
	// 		cfn.StackStatusDeleteComplete,
	// 		[]string{
	// 			cfn.StackStatusDeleteFailed,
	// 			cfn.StackStatusCreateInProgress,
	// 			cfn.StackStatusCreateFailed,
	// 			cfn.StackStatusCreateComplete,
	// 			cfn.StackStatusRollbackInProgress,
	// 			cfn.StackStatusRollbackFailed,
	// 			cfn.StackStatusRollbackComplete,
	// 			cfn.StackStatusUpdateInProgress,
	// 			cfn.StackStatusUpdateCompleteCleanupInProgress,
	// 			cfn.StackStatusUpdateComplete,
	// 			cfn.StackStatusUpdateRollbackInProgress,
	// 			cfn.StackStatusUpdateRollbackFailed,
	// 			cfn.StackStatusUpdateRollbackCompleteCleanupInProgress,
	// 			cfn.StackStatusUpdateRollbackComplete,
	// 			cfn.StackStatusReviewInProgress,
	// 		},
	// 		// ValidationError is expected as success, although
	// 		// we use stack ARN, so should normally see actual
	// 		// StackStatusDeleteComplete
	// 		request.WaiterAcceptor{
	// 			State:    request.SuccessWaiterState,
	// 			Matcher:  request.ErrorWaiterMatch,
	// 			Expected: "ValidationError",
	// 		},
	// 	),
	// )
}

func (c *StackCollection) waitUntilStackIsDeleted(i *Stack, errs chan error) {
	// defer close(errs)

	// if err := c.doWaitUntilStackIsDeleted(i); err != nil {
	// 	errs <- err
	// 	return
	// }
	// errs <- nil
}

func (c *StackCollection) doWaitUntilStackIsUpdated(i *Stack) error {
	return nil
	// return c.waitWithAcceptors(i,
	// 	waiters.MakeAcceptors(
	// 		stackStatus,
	// 		types.StackStatusUpdateComplete,
	// 		[]string{
	// 			string(types.StackStatusUpdateRollbackComplete),
	// 			string(types.StackStatusUpdateRollbackFailed),
	// 			string(types.StackStatusUpdateRollbackInProgress),
	// 			string(types.StackStatusRollbackInProgress),
	// 			string(types.StackStatusRollbackFailed),
	// 			string(types.StackStatusRollbackComplete),
	// 			string(types.StackStatusDeleteInProgress),
	// 			string(types.StackStatusDeleteFailed),
	// 			string(types.StackStatusDeleteComplete),
	// 		},
	// 		request.WaiterAcceptor{
	// 			State:    request.FailureWaiterState,
	// 			Matcher:  request.ErrorWaiterMatch,
	// 			Expected: "ValidationError",
	// 		},
	// 	),
	// )
}

func (c *StackCollection) doWaitUntilChangeSetIsCreated(i *Stack, changesetName string) error {
	return nil
	// return c.waitWithAcceptorsChangeSet(i, changesetName,
	// 	waiters.MakeAcceptors(
	// 		changesetStatus,
	// 		types.ChangeSetStatusCreateComplete,
	// 		[]string{
	// 			string(types.ChangeSetStatusDeleteComplete),
	// 			string(types.ChangeSetStatusFailed),
	// 		},
	// 		request.WaiterAcceptor{
	// 			State:    request.FailureWaiterState,
	// 			Matcher:  request.ErrorWaiterMatch,
	// 			Expected: "ValidationError",
	// 		},
	// 	),
	// )
}
