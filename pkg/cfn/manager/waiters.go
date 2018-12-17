package manager

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/request"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
)

const (
	stackStatus     = "Stacks[].StackStatus"
	changesetStatus = "Status"
)

// cloudformation.WaitUntilStackCreateComplete doesn't detect in-progress status early enough,
// so this is custom version that is more suitable for our use, as there is no way to add any
// custom acceptors

func makeStatusAcceptor(status string, statusPath string) request.WaiterAcceptor {
	return request.WaiterAcceptor{
		Matcher:  request.PathAllWaiterMatch,
		Argument: statusPath,
		Expected: status,
		State:    request.FailureWaiterState,
	}
}

func makeAcceptors(statusPath string, successStatus string, failureStates []string, extraAcceptors ...request.WaiterAcceptor) []request.WaiterAcceptor {
	acceptors := []request.WaiterAcceptor{makeStatusAcceptor(successStatus, statusPath)}
	acceptors[0].State = request.SuccessWaiterState

	for _, s := range failureStates {
		acceptors = append(acceptors, makeStatusAcceptor(s, statusPath))
	}

	acceptors = append(acceptors, extraAcceptors...)

	return acceptors
}

// makeWaiterDelay returns delay ranging between 15s and 20s
func makeWaiterDelay() request.WaiterDelay {
	const (
		base        = 15 * time.Second
		offsetSteps = 200
		offsetMax   = 5000
		stepMult    = offsetMax / offsetSteps
	)

	offsets := rand.Perm(offsetSteps)

	return func(attempt int) time.Duration {
		s := rand.Intn(offsetSteps)
		d := stepMult * offsets[s]

		offset := time.Duration(d) * time.Millisecond

		return base + offset
	}
}

func (c *StackCollection) waitWithAcceptors(i *Stack, acceptors []request.WaiterAcceptor) error {
	desiredStatus := fmt.Sprintf("%v", acceptors[0].Expected)
	msg := fmt.Sprintf("waiting for CloudFormation stack %q to reach %q status", *i.StackName, desiredStatus)

	ctx, cancel := context.WithTimeout(context.Background(), c.provider.WaitTimeout())
	defer cancel()

	startTime := time.Now()

	w := request.Waiter{
		Name:        strings.Join([]string{"wait", *i.StackName, desiredStatus}, "_"),
		MaxAttempts: 1024, // we use context deadline instead
		Delay:       makeWaiterDelay(),
		Acceptors:   acceptors,
		NewRequest: func(_ []request.Option) (*request.Request, error) {
			input := &cfn.DescribeStacksInput{
				StackName: i.StackName,
			}
			if i.StackId != nil && *i.StackId != "" {
				input.StackName = i.StackId
			}
			logger.Debug(msg)
			req, _ := c.provider.CloudFormation().DescribeStacksRequest(input)
			req.SetContext(ctx)
			return req, nil
		},
	}

	logger.Debug("start %s", msg)

	if waitErr := w.WaitWithContext(ctx); waitErr != nil {
		s, err := c.describeStack(i)
		if err != nil {
			logger.Debug("describeErr=%v", err)
		} else {
			logger.Critical("unexpected status %q while %s", *s.StackStatus, msg)
			c.troubleshootStackFailureCause(i, desiredStatus)
		}
		return errors.Wrap(waitErr, msg)
	}

	logger.Debug("done after %s of %s", time.Since(startTime), msg)

	return nil
}

func (c *StackCollection) waitWithAcceptorsChangeSet(i *Stack, changesetName *string, acceptors []request.WaiterAcceptor) error {
	desiredStatus := fmt.Sprintf("%v", acceptors[0].Expected)
	msg := fmt.Sprintf("waiting for CloudFormation changeset %q for stack %q to reach %q status", *changesetName, *i.StackName, desiredStatus)
	ctx, cancel := context.WithTimeout(context.Background(), c.provider.WaitTimeout())
	defer cancel()
	startTime := time.Now()
	w := request.Waiter{
		Name:        strings.Join([]string{"waitCS", *i.StackName, *changesetName, desiredStatus}, "_"),
		MaxAttempts: 1024, // we use context deadline instead
		Delay:       makeWaiterDelay(),
		Acceptors:   acceptors,
		NewRequest: func(_ []request.Option) (*request.Request, error) {
			input := &cfn.DescribeChangeSetInput{
				StackName:     i.StackName,
				ChangeSetName: changesetName,
			}
			logger.Debug(msg)
			req, _ := c.provider.CloudFormation().DescribeChangeSetRequest(input)
			req.SetContext(ctx)
			return req, nil
		},
	}
	logger.Debug("start %s", msg)
	if waitErr := w.WaitWithContext(ctx); waitErr != nil {
		s, err := c.describeStackChangeSet(i, changesetName)
		if err != nil {
			logger.Debug("describeChangeSetErr=%v", err)
		} else {
			logger.Critical("unexpected status %q while %s, reason %s", *s.Status, msg, *s.StatusReason)
		}
		return errors.Wrap(waitErr, msg)
	}
	logger.Debug("done after %s of %s", time.Since(startTime), msg)
	return nil
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
				logger.Info(msg)
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

func (c *StackCollection) doWaitUntilStackIsCreated(i *Stack) error {
	return c.waitWithAcceptors(i,
		makeAcceptors(
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

	if err := c.doWaitUntilStackIsCreated(i); err != nil {
		errs <- err
		return
	}
	s, err := c.describeStack(i)
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
		makeAcceptors(
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

func (c *StackCollection) doWaitUntilStackIsUpdated(i *Stack) error {
	return c.waitWithAcceptors(i,
		makeAcceptors(
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

func (c *StackCollection) doWaitUntilChangeSetIsCreated(i *Stack, changesetName *string) error {
	return c.waitWithAcceptorsChangeSet(i, changesetName,
		makeAcceptors(
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
