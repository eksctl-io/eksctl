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
func (p *Waiter) Wait(ctx context.Context) error {
	for attempts := 1; ; attempts++ {
		done, err := p.wait(ctx, p.NextDelay(attempts))
		if err != nil {
			return err
		}
		if done {
			return nil
		}
	}
}

func (p *Waiter) wait(ctx context.Context, d time.Duration) (bool, error) {
	waitTimer := time.NewTimer(d)
	select {
	case <-waitTimer.C:
		return p.Operation()

	case <-ctx.Done():
		waitTimer.Stop()
		return false, ctx.Err()
	}
}
