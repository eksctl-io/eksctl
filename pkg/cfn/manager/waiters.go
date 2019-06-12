package manager

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/request"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/utils/waiters"
)

const (
	stackStatus     = "Stacks[].StackStatus"
	changesetStatus = "Status"
)

// cloudformation.WaitUntilStackCreateComplete doesn't detect in-progress status early enough,
// so this is custom version that is more suitable for our use, as there is no way to add any
// custom acceptors

func (c *StackCollection) waitWithAcceptors(i *Stack, acceptors []request.WaiterAcceptor) error {
	msg := fmt.Sprintf("waiting for CloudFormation stack %q", *i.StackName)

	newRequest := func() *request.Request {
		input := &cfn.DescribeStacksInput{
			StackName: i.StackName,
		}
		if api.IsSetAndNonEmptyString(i.StackId) {
			input.StackName = i.StackId
		}
		req, _ := c.provider.CloudFormation().DescribeStacksRequest(input)
		return req
	}

	troubleshoot := func(desiredStatus string) {
		s, err := c.DescribeStack(i)
		if err != nil {
			logger.Debug("describeErr=%v", err)
		} else {
			logger.Critical("unexpected status %q while %s", *s.StackStatus, msg)
			c.troubleshootStackFailureCause(i, desiredStatus)
		}
	}

	return waiters.Wait(*i.StackName, msg, acceptors, newRequest, c.provider.WaitTimeout(), troubleshoot)
}

func (c *StackCollection) waitWithAcceptorsChangeSet(i *Stack, changesetName string, acceptors []request.WaiterAcceptor) error {
	msg := fmt.Sprintf("waiting for CloudFormation changeset %q for stack %q", changesetName, *i.StackName)

	newRequest := func() *request.Request {
		input := &cfn.DescribeChangeSetInput{
			StackName:     i.StackName,
			ChangeSetName: &changesetName,
		}
		req, _ := c.provider.CloudFormation().DescribeChangeSetRequest(input)
		return req
	}

	troubleshoot := func(desiredStatus string) {
		s, err := c.DescribeStackChangeSet(i, changesetName)
		if err != nil {
			logger.Debug("describeChangeSetErr=%v", err)
		} else {
			logger.Critical("unexpected status %q while %s, reason: %s", *s.Status, msg, *s.StatusReason)
		}
	}

	return waiters.Wait(*i.StackName, msg, acceptors, newRequest, c.provider.WaitTimeout(), troubleshoot)
}

func (c *StackCollection) troubleshootStackFailureCause(i *Stack, desiredStatus string) {
	logger.Info("fetching stack events in attempt to troubleshoot the root cause of the failure")
	events, err := c.DescribeStackEvents(i)
	if err != nil {
		logger.Critical("cannot fetch stack events: %v", err)
		return
	}
	for _, e := range events {
		msg := fmt.Sprintf("%s/%s: %s", *e.ResourceType, *e.LogicalResourceId, *e.ResourceStatus)
		if e.ResourceStatusReason != nil {
			msg = fmt.Sprintf("%s – %#v", msg, *e.ResourceStatusReason)
		}
		switch desiredStatus {
		case cfn.StackStatusCreateComplete:
			switch *e.ResourceStatus {
			case cfn.ResourceStatusCreateFailed:
				logger.Critical(msg)
			case cfn.ResourceStatusDeleteInProgress:
				logger.Warning(msg)
			default:
				logger.Debug(msg) // only output this when verbose logging is enabled
			}
		case cfn.StackStatusDeleteComplete:
			switch *e.ResourceStatus {
			case cfn.ResourceStatusDeleteFailed:
				logger.Critical(msg)
			case cfn.ResourceStatusDeleteSkipped:
				logger.Warning(msg)
			default:
				logger.Info(msg)
			}
		default:
			logger.Info(msg)
		}
	}
}

// DoWaitUntilStackIsCreated blocks until the given stack's
// creation has completed.
func (c *StackCollection) DoWaitUntilStackIsCreated(i *Stack) error {
	return c.waitWithAcceptors(i,
		waiters.MakeAcceptors(
			stackStatus,
			cfn.StackStatusCreateComplete,
			[]string{
				cfn.StackStatusCreateFailed,
				cfn.StackStatusRollbackInProgress,
				cfn.StackStatusRollbackFailed,
				cfn.StackStatusRollbackComplete,
				cfn.StackStatusDeleteInProgress,
				cfn.StackStatusDeleteFailed,
				cfn.StackStatusDeleteComplete,
			},
			request.WaiterAcceptor{
				State:    request.FailureWaiterState,
				Matcher:  request.ErrorWaiterMatch,
				Expected: "ValidationError",
			},
		),
	)
}

func (c *StackCollection) waitUntilStackIsCreated(i *Stack, stack builder.ResourceSet, errs chan error) {
	defer close(errs)

	if err := c.DoWaitUntilStackIsCreated(i); err != nil {
		errs <- err
		return
	}
	s, err := c.DescribeStack(i)
	if err != nil {
		errs <- err
		return
	}
	if err := stack.GetAllOutputs(*s); err != nil {
		errs <- errors.Wrapf(err, "getting stack %q outputs", *i.StackName)
		return
	}
	errs <- nil
}

func (c *StackCollection) doWaitUntilStackIsDeleted(i *Stack) error {
	return c.waitWithAcceptors(i,
		waiters.MakeAcceptors(
			stackStatus,
			cfn.StackStatusDeleteComplete,
			[]string{
				cfn.StackStatusDeleteFailed,
				cfn.StackStatusCreateInProgress,
				cfn.StackStatusCreateFailed,
				cfn.StackStatusCreateComplete,
				cfn.StackStatusRollbackInProgress,
				cfn.StackStatusRollbackFailed,
				cfn.StackStatusRollbackComplete,
				cfn.StackStatusUpdateInProgress,
				cfn.StackStatusUpdateCompleteCleanupInProgress,
				cfn.StackStatusUpdateComplete,
				cfn.StackStatusUpdateRollbackInProgress,
				cfn.StackStatusUpdateRollbackFailed,
				cfn.StackStatusUpdateRollbackCompleteCleanupInProgress,
				cfn.StackStatusUpdateRollbackComplete,
				cfn.StackStatusReviewInProgress,
			},
			// ValidationError is expected as success, although
			// we use stack ARN, so should normally see actual
			// StackStatusDeleteComplete
			request.WaiterAcceptor{
				State:    request.SuccessWaiterState,
				Matcher:  request.ErrorWaiterMatch,
				Expected: "ValidationError",
			},
		),
	)
}

func (c *StackCollection) waitUntilStackIsDeleted(i *Stack, errs chan error) {
	defer close(errs)

	if err := c.doWaitUntilStackIsDeleted(i); err != nil {
		errs <- err
		return
	}
	errs <- nil
}

func (c *StackCollection) doWaitUntilStackIsUpdated(i *Stack) error {
	return c.waitWithAcceptors(i,
		waiters.MakeAcceptors(
			stackStatus,
			cfn.StackStatusUpdateComplete,
			[]string{
				cfn.StackStatusUpdateRollbackComplete,
				cfn.StackStatusUpdateRollbackFailed,
				cfn.StackStatusUpdateRollbackInProgress,
				cfn.StackStatusRollbackInProgress,
				cfn.StackStatusRollbackFailed,
				cfn.StackStatusRollbackComplete,
				cfn.StackStatusDeleteInProgress,
				cfn.StackStatusDeleteFailed,
				cfn.StackStatusDeleteComplete,
			},
			request.WaiterAcceptor{
				State:    request.FailureWaiterState,
				Matcher:  request.ErrorWaiterMatch,
				Expected: "ValidationError",
			},
		),
	)
}

func (c *StackCollection) doWaitUntilChangeSetIsCreated(i *Stack, changesetName string) error {
	return c.waitWithAcceptorsChangeSet(i, changesetName,
		waiters.MakeAcceptors(
			changesetStatus,
			cfn.ChangeSetStatusCreateComplete,
			[]string{
				cfn.ChangeSetStatusDeleteComplete,
				cfn.ChangeSetStatusFailed,
			},
			request.WaiterAcceptor{
				State:    request.FailureWaiterState,
				Matcher:  request.ErrorWaiterMatch,
				Expected: "ValidationError",
			},
		),
	)
}
