package waiter

import (
	"context"
	"time"
)

// A Waiter keeps retrying the specified operation until it returns true, or an error.
type Waiter struct {
	// NextDelay is the amount of time to wait before the next retry given the number of attempts.
	NextDelay NextDelay

	// Operation is the function to invoke.
	Operation func() (bool, error)
}

// Wait waits for the specified operation to complete.
func (w *Waiter) Wait(ctx context.Context) error {
	for attempts := 1; ; attempts++ {
		done, err := w.wait(ctx, w.NextDelay(attempts))
		if err != nil {
			return err
		}
		if done {
			return nil
		}
	}
}

// WaitWithTimeout is a wrapper around Wait that takes a timeout value instead of a Context,
// and returns a DeadlineExceeded error when the timeout expires.
// It exists to allow interfacing with code that is not using contexts yet.
func (w *Waiter) WaitWithTimeout(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return w.Wait(ctx)
}

func (w *Waiter) wait(ctx context.Context, d time.Duration) (bool, error) {
	waitTimer := time.NewTimer(d)
	select {
	case <-waitTimer.C:
		return w.Operation()

	case <-ctx.Done():
		waitTimer.Stop()
		return false, ctx.Err()
	}
}
