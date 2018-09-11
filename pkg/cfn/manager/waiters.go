package manager

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws/request"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"

	"github.com/kubicorn/kubicorn/pkg/logger"

	"github.com/weaveworks/eksctl/pkg/cfn/builder"
)

// cloudformation.WaitUntilStackCreateComplete doesn't detect in-progress status early enough,
// so this is custom version that is more suitable for our use, as there is no way to add any
// custom acceptors

func makeStatusAcceptor(status string) request.WaiterAcceptor {
	return request.WaiterAcceptor{
		Matcher:  request.PathAllWaiterMatch,
		Argument: "Stacks[].StackStatus",
		Expected: status,
		State:    request.FailureWaiterState,
	}
}

func makeAcceptors(successStatus string, failureStates []string, extraAcceptors ...request.WaiterAcceptor) []request.WaiterAcceptor {
	acceptors := []request.WaiterAcceptor{makeStatusAcceptor(successStatus)}
	acceptors[0].State = request.SuccessWaiterState

	for _, s := range failureStates {
		acceptors = append(acceptors, makeStatusAcceptor(s))
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

	ctx, cancel := context.WithTimeout(context.Background(), c.spec.WaitTimeout)
	defer cancel()

	startTime := time.Now()

	w := request.Waiter{
		Name:        strings.Join([]string{"wait", *i.StackName, desiredStatus}, "_"),
		MaxAttempts: 1024, // we use context deadline instead
		Delay:       makeWaiterDelay(),
		Acceptors:   acceptors,
		NewRequest: func(_ []request.Option) (*request.Request, error) {
			input := &cfn.DescribeStacksInput{
				StackName: i.StackId,
			}
			logger.Debug(msg)
			req, _ := c.cfn.DescribeStacksRequest(input)
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
		}
		return errors.Wrap(waitErr, msg)
	}

	logger.Debug("done after %s of %s", time.Since(startTime), msg)

	return nil
}

func (c *StackCollection) doWaitUntilStackIsCreated(i *Stack) error {
	return c.waitWithAcceptors(i,
		makeAcceptors(
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
