package waiters

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
)

// Wait for something with a name to reach status that is expressed by acceptors using newRequest
// until we hit waitTimeout, on unexpected status troubleshoot will be called with the desired
// status as an argument, so that it can find what migth have gone wrong
func Wait(name, msg string, acceptors []request.WaiterAcceptor, newRequest func() *request.Request, waitTimeout time.Duration, troubleshoot func(string) error) error {
	desiredStatus := fmt.Sprintf("%v", acceptors[0].Expected)
	name = strings.Join([]string{"wait", name, desiredStatus}, "_")

	ctx, cancel := context.WithTimeout(context.Background(), waitTimeout)
	defer cancel()
	startTime := time.Now()
	w := makeWaiter(ctx, name, msg, acceptors, newRequest)
	logger.Info("start %s", msg)
	if waitErr := w.WaitWithContext(ctx); waitErr != nil {
		if troubleshoot != nil {
			if wrappedErr := troubleshoot(desiredStatus); wrappedErr != nil {
				return wrappedErr
			}
		}
		return errors.Wrap(waitErr, msg)
	}
	logger.Info("done after %s of %s", time.Since(startTime), msg)
	return nil
}

func makeWaiter(ctx context.Context, name, msg string, acceptors []request.WaiterAcceptor, newRequest func() *request.Request) request.Waiter {
	return request.Waiter{
		Name:        name,
		MaxAttempts: 1024, // we use context deadline instead
		Delay:       makeWaiterDelay(),
		Acceptors:   acceptors,
		NewRequest: func(_ []request.Option) (*request.Request, error) {
			logger.Info(msg)
			req := newRequest()
			req.SetContext(ctx)
			return req, nil
		},
	}
}

// MakeAcceptors constructs a slice of request acceptors
func MakeAcceptors(statusPath string, successStatus string, failureStates []string, extraAcceptors ...request.WaiterAcceptor) []request.WaiterAcceptor {
	acceptors := []request.WaiterAcceptor{makeStatusAcceptor(successStatus, statusPath)}
	acceptors[0].State = request.SuccessWaiterState

	for _, s := range failureStates {
		acceptors = append(acceptors, makeStatusAcceptor(s, statusPath))
	}

	acceptors = append(acceptors, extraAcceptors...)

	return acceptors
}

func makeStatusAcceptor(status string, statusPath string) request.WaiterAcceptor {
	return request.WaiterAcceptor{
		Matcher:  request.PathAllWaiterMatch,
		Argument: statusPath,
		Expected: status,
		State:    request.FailureWaiterState,
	}
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
