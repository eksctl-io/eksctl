package waiter

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/kris-nova/logger"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/aws/smithy-go/middleware"
	smithytime "github.com/aws/smithy-go/time"
	smithywaiter "github.com/aws/smithy-go/waiter"

	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/utils/apierrors"
)

type UpdateWaiterOptions struct {
	// Set of options to modify how an operation is invoked. These apply to all
	// operations invoked for this client. Use functional options on operation call to
	// modify this list for per operation behaviour.
	APIOptions []func(*middleware.Stack) error

	// MinDelay is the minimum amount of time to delay between retries. If unset,
	// NodegroupDeletedWaiter will use default minimum delay of 30 seconds. Note that
	// MinDelay must resolve to a value lesser than or equal to the MaxDelay.
	MinDelay time.Duration

	// MaxDelay is the maximum amount of time to delay between retries. If unset or set
	// to zero, NodegroupDeletedWaiter will use default max delay of 120 seconds. Note
	// that MaxDelay must resolve to value greater than or equal to the MinDelay.
	MaxDelay time.Duration

	// RetryAttemptLogMessage is the message to log before attempting the next request.
	RetryAttemptLogMessage string

	// Retryable is function that can be used to override the service defined
	// waiter-behavior based on operation output, or returned error. This function is
	// used by the waiter to decide if a state is retryable or a terminal state. By
	// default service-modeled logic will populate this option. This option can thus be
	// used to define a custom waiter state with fall-back to service-modeled waiter
	// state mutators.The function returns an error in case of a failure state. In case
	// of retry state, this function returns a bool value of true and nil error, while
	// in case of success it returns a bool value of false and nil error.
	Retryable func(context.Context, *eks.DescribeUpdateInput, *eks.DescribeUpdateOutput, error) (bool, error)
}

type UpdateWaiter struct {
	client  awsapi.EKS
	options UpdateWaiterOptions
}

type UpdateFailedError struct {
	Status      string
	UpdateError string
}

func (u *UpdateFailedError) Error() string {
	return fmt.Sprintf("update failed with status %q: %s", u.Status, u.UpdateError)
}

// NewUpdateWaiter constructs an UpdateWaiter.
// It provides an interface similar to the waiter types in the AWS SDK.
func NewUpdateWaiter(client awsapi.EKS, optFns ...func(options *UpdateWaiterOptions)) *UpdateWaiter {
	options := UpdateWaiterOptions{
		MinDelay: 30 * time.Second,
		MaxDelay: 120 * time.Second,
	}
	options.Retryable = func(ctx context.Context, input *eks.DescribeUpdateInput, output *eks.DescribeUpdateOutput, err error) (bool, error) {
		if err != nil {
			if apierrors.IsRetriableError(err) {
				return true, nil
			}
			return false, err
		}

		switch output.Update.Status {
		case ekstypes.UpdateStatusSuccessful:
			return false, nil
		case ekstypes.UpdateStatusFailed, ekstypes.UpdateStatusCancelled:
			return false, &UpdateFailedError{
				Status:      string(output.Update.Status),
				UpdateError: fmt.Sprintf("update errors:\n%s", aggregateErrors(output.Update.Errors)),
			}
		default:
			return true, nil
		}
	}

	for _, fn := range optFns {
		fn(&options)
	}
	return &UpdateWaiter{
		client:  client,
		options: options,
	}
}

// Wait calls the waiter function for UpdateWaiter waiter. The maxWaitDur is the
// maximum wait duration the waiter will wait. The maxWaitDur is required and must
// be greater than zero.
func (w *UpdateWaiter) Wait(ctx context.Context, params *eks.DescribeUpdateInput, maxWaitDur time.Duration, optFns ...func(waiter *UpdateWaiterOptions)) error {
	_, err := w.WaitForOutput(ctx, params, maxWaitDur, optFns...)
	return err
}

// WaitForOutput calls the waiter function for UpdateWaiter waiter and returns
// the output of the successful operation. The maxWaitDur is the maximum wait
// duration the waiter will wait. The maxWaitDur is required and must be greater
// than zero.
func (w *UpdateWaiter) WaitForOutput(ctx context.Context, params *eks.DescribeUpdateInput, maxWaitDur time.Duration, optFns ...func(options *UpdateWaiterOptions)) (*eks.DescribeUpdateOutput, error) {
	if maxWaitDur <= 0 {
		return nil, errors.New("maximum wait time for waiter must be greater than zero")
	}

	options := w.options
	for _, fn := range optFns {
		fn(&options)
	}

	if options.MaxDelay <= 0 {
		options.MaxDelay = 120 * time.Second
	}

	if options.MinDelay > options.MaxDelay {
		return nil, fmt.Errorf("minimum waiter delay %v must be lesser than or equal to maximum waiter delay of %v", options.MinDelay, options.MaxDelay)
	}

	ctx, cancelFn := context.WithTimeout(ctx, maxWaitDur)
	defer cancelFn()

	remainingTime := maxWaitDur
	startTime := time.Now()

	var attempt int64
	for {

		attempt++
		apiOptions := options.APIOptions
		start := time.Now()

		logger.Debug(options.RetryAttemptLogMessage)

		out, err := w.client.DescribeUpdate(ctx, params, func(o *eks.Options) {
			o.APIOptions = append(o.APIOptions, apiOptions...)
		})

		retryable, err := options.Retryable(ctx, params, out, err)
		if err != nil {
			return nil, err
		}
		if !retryable {
			logger.Debug("done after %s of %s", time.Since(startTime), options.RetryAttemptLogMessage)
			return out, nil
		}

		remainingTime -= time.Since(start)
		if remainingTime < options.MinDelay || remainingTime <= 0 {
			break
		}

		// compute exponential backoff between waiter retries
		delay, err := smithywaiter.ComputeDelay(
			attempt, options.MinDelay, options.MaxDelay, remainingTime,
		)
		if err != nil {
			return nil, fmt.Errorf("error computing waiter delay, %w", err)
		}

		remainingTime -= delay
		// sleep for the delay amount before invoking a request
		if err := smithytime.SleepWithContext(ctx, delay); err != nil {
			return nil, fmt.Errorf("request cancelled while waiting, %w", err)
		}
	}
	return nil, errors.New("exceeded max wait time for UpdateWaiter")
}

func aggregateErrors(errorDetails []ekstypes.ErrorDetail) string {
	var aggregatedErrors []string
	for _, err := range errorDetails {
		msg := fmt.Sprintf("%s; errorCode: %s, resourceIDs: %v", *err.ErrorMessage, err.ErrorCode, err.ResourceIds)
		aggregatedErrors = append(aggregatedErrors, fmt.Sprintf("- %s", msg))
	}
	return strings.Join(aggregatedErrors, "\n")
}
