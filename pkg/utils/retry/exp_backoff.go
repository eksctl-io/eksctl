package retry

import (
	"math"
	"time"
)

// ExponentialBackoff defines a retry policy in which in we exponentially
// retry up to the provided maximum number of retries (MaxRetries).
type ExponentialBackoff struct {
	retry      int
	MaxRetries int
	TimeUnit   time.Duration
}

// Done implements retry.Policy#Done() bool.
func (b ExponentialBackoff) Done() bool {
	return b.retry == b.MaxRetries
}

// Duration implements retry.Policy#Duration() time.Duration.
func (b *ExponentialBackoff) Duration() time.Duration {
	duration := time.Duration(pow(2, b.retry)) * b.TimeUnit
	b.retry++
	return duration
}

// Reset implements retry.Policy#Reset().
func (b *ExponentialBackoff) Reset() {
	b.retry = 0
}

// Clone implements retry.Policy#Clone() retry.Policy.
func (b ExponentialBackoff) Clone() Policy {
	return &ExponentialBackoff{
		MaxRetries: b.MaxRetries,
		TimeUnit:   b.TimeUnit,
	}
}

// TimingOutExponentialBackoff defines a retry policy in which we exponentially
// retry up to the provided maximum duration (Timeout).
type TimingOutExponentialBackoff struct {
	retry          int
	totalTimeSoFar time.Duration
	Timeout        time.Duration
	TimeUnit       time.Duration
}

// Done implements retry.Policy#Done() bool.
func (b TimingOutExponentialBackoff) Done() bool {
	return b.totalTimeSoFar == b.Timeout
}

// Duration implements retry.Policy#Duration() time.Duration.
func (b *TimingOutExponentialBackoff) Duration() time.Duration {
	duration := time.Duration(pow(2, b.retry)) * b.TimeUnit
	b.retry++
	// Cap duration so that the configured timeout is never exceeded:
	if b.totalTimeSoFar+duration > b.Timeout {
		duration = b.Timeout - b.totalTimeSoFar
	}
	b.totalTimeSoFar += duration
	return duration
}

// Reset implements retry.Policy#Reset().
func (b *TimingOutExponentialBackoff) Reset() {
	b.retry = 0
	b.totalTimeSoFar = 0
}

// Clone implements retry.Policy#Clone() retry.Policy.
func (b TimingOutExponentialBackoff) Clone() Policy {
	return &TimingOutExponentialBackoff{
		Timeout:  b.Timeout,
		TimeUnit: b.TimeUnit,
	}
}

func pow(x, y int) int32 {
	return int32(math.Pow(float64(x), float64(y)))
}
