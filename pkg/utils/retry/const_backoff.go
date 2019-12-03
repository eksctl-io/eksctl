package retry

import (
	"time"
)

// ConstantBackoff defines a retry policy in which in we retry up to the
// provided maximum number of retries (MaxRetries). The duration is constant
// and always the one provided, i.e. Time converted to TimeUnit.
type ConstantBackoff struct {
	retry      int
	MaxRetries int
	Time       int
	TimeUnit   time.Duration
}

// Done implements retry.Policy#Done() bool.
func (b ConstantBackoff) Done() bool {
	return b.retry == b.MaxRetries
}

// Duration implements retry.Policy#Duration() time.Duration.
func (b *ConstantBackoff) Duration() time.Duration {
	b.retry++
	return time.Duration(b.Time) * b.TimeUnit
}

// Reset implements retry.Policy#Reset().
func (b *ConstantBackoff) Reset() {
	b.retry = 0
}

// Clone implements retry.Policy#Clone() retry.Policy.
func (b ConstantBackoff) Clone() Policy {
	return &ConstantBackoff{
		MaxRetries: b.MaxRetries,
		Time:       b.Time,
		TimeUnit:   b.TimeUnit,
	}
}
