package retry

import "time"

// Policy defines an "iterator"-style of interface for retry policies.
// The client should repeatedly call Done() and Duration() such that if Done()
// returns false, then Duration() is called, and the value returned should be
// the time to wait for before retrying an operation defined by the client.
// Once Done() returns true, the client should stop using the retry policy and
// "give up" in an appropriate manner for its use-case.
// Reset() resets the current policy's state, so that it can be re-used.
// Clone() clones the current policy, so that it can be used in a different
// thread/go-routine.
type Policy interface {
	Done() bool
	Duration() time.Duration
	Reset()
	Clone() Policy
}
