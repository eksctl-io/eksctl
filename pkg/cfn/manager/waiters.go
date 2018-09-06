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

func (c *StackCollection) waitWithAcceptors(name string, acceptors []request.WaiterAcceptor) error {
	desiredStatus := fmt.Sprintf("%v", acceptors[0].Expected)
	msg := fmt.Sprintf("waiting for CloudFormation stack %q to reach %q status", name, desiredStatus)

	ctx, cancel := context.WithTimeout(context.Background(), c.spec.WaitTimeout)
	defer cancel()

	startTime := time.Now()

	w := request.Waiter{
		Name:        strings.Join([]string{"wait", name, desiredStatus}, "_"),
		MaxAttempts: 1024, // we use context deadline instead
		Delay:       makeWaiterDelay(),
		Acceptors:   acceptors,
		NewRequest: func(_ []request.Option) (*request.Request, error) {
			input := &cfn.DescribeStacksInput{
				StackName: &name,
			}
			logger.Debug(msg)
			req, _ := c.cfn.DescribeStacksRequest(input)
			req.SetContext(ctx)
			return req, nil
		},
	}

	logger.Debug("start %s", msg)

	if waitErr := w.WaitWithContext(ctx); waitErr != nil {
		s, err := c.describeStack(name)
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

func (c *StackCollection) doWaitUntilStackIsCreated(name string) error {
	return c.waitWithAcceptors(name,
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

func (c *StackCollection) waitUntilStackIsCreated(name string, stack builder.ResourceSet, errs chan error) {
	defer close(errs)

	if err := c.doWaitUntilStackIsCreated(name); err != nil {
		errs <- err
		return
	}
	s, err := c.describeStack(name)
	if err != nil {
		errs <- err
		return
	}
	if err := stack.GetAllOutputs(*s); err != nil {
		errs <- errors.Wrapf(err, "getting stack %q outputs", name)
		return
	}
	errs <- nil
}

func (c *StackCollection) doWaitUntilStackIsDeleted(name string) error {
	return c.waitWithAcceptors(name,
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
			// ValidationError really means stack not found, that's what
			// you get from DescribeStack call, you never get actual stack;
			// you do get stack with StackStatusDeleteComplete on ListStacks,
			// but that returns pages and pages, so we don't actually want
			// to worry about that, in fact it also returns ones that had
			// the same name but were deted a long time ago
			request.WaiterAcceptor{
				State:    request.SuccessWaiterState,
				Matcher:  request.ErrorWaiterMatch,
				Expected: "ValidationError",
			},
		),
	)
}
