package retry_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/utils/retry"
)

var _ = Describe("retry", func() {
	Describe("ExponentialBackoff", func() {
		It("generates a sequence of exponentially increasing durations", func() {
			policy := retry.ExponentialBackoff{
				MaxRetries: 5,
				TimeUnit:   time.Second,
			}
			Expect(policy.Done()).To(BeFalse())
			Expect(policy.Duration()).To(Equal(1 * time.Second))
			Expect(policy.Done()).To(BeFalse())
			Expect(policy.Duration()).To(Equal(2 * time.Second))
			Expect(policy.Done()).To(BeFalse())
			Expect(policy.Duration()).To(Equal(4 * time.Second))
			Expect(policy.Done()).To(BeFalse())
			Expect(policy.Duration()).To(Equal(8 * time.Second))
			Expect(policy.Done()).To(BeFalse())
			Expect(policy.Duration()).To(Equal(16 * time.Second))
			Expect(policy.Done()).To(BeTrue())
		})

		Describe("Reset", func() {
			It("resets the current policy so it can be re-used", func() {
				policy := retry.ExponentialBackoff{
					MaxRetries: 1,
					TimeUnit:   time.Second,
				}
				Expect(policy.Done()).To(BeFalse())
				Expect(policy.Duration()).To(Equal(1 * time.Second))
				Expect(policy.Done()).To(BeTrue())
				policy.Reset()
				Expect(policy.Done()).To(BeFalse())
				Expect(policy.Duration()).To(Equal(1 * time.Second))
				Expect(policy.Done()).To(BeTrue())
			})
		})

		Describe("Clone", func() {
			It("clones the current policy so both can be used concurrently", func() {
				policy := retry.ExponentialBackoff{
					MaxRetries: 1,
					TimeUnit:   time.Second,
				}
				Expect(policy.Done()).To(BeFalse())
				Expect(policy.Duration()).To(Equal(1 * time.Second))
				Expect(policy.Done()).To(BeTrue())
				clone := policy.Clone()
				addrPolicy := fmt.Sprintf("%p", &policy)
				addrClone := fmt.Sprintf("%p", &clone)
				Expect(addrPolicy).To(Not(Equal(addrClone)))
				// The original policy is still in a terminal state:
				Expect(policy.Done()).To(BeTrue())
				// The cloned policy can be used:
				Expect(clone.Done()).To(BeFalse())
				Expect(clone.Duration()).To(Equal(1 * time.Second))
				Expect(clone.Done()).To(BeTrue())
			})
		})
	})

	Describe("TimingOutExponentialBackoff", func() {
		It("generates a sequence of exponentially increasing durations, capped by the provided timeout", func() {
			policy := retry.TimingOutExponentialBackoff{
				Timeout:  10 * time.Minute,
				TimeUnit: time.Second,
			}
			Expect(policy.Done()).To(BeFalse())
			Expect(policy.Duration()).To(Equal(1 * time.Second))
			Expect(policy.Done()).To(BeFalse())
			Expect(policy.Duration()).To(Equal(2 * time.Second))
			Expect(policy.Done()).To(BeFalse())
			Expect(policy.Duration()).To(Equal(4 * time.Second))
			Expect(policy.Done()).To(BeFalse())
			Expect(policy.Duration()).To(Equal(8 * time.Second))
			Expect(policy.Done()).To(BeFalse())
			Expect(policy.Duration()).To(Equal(16 * time.Second))
			Expect(policy.Done()).To(BeFalse())
			Expect(policy.Duration()).To(Equal(32 * time.Second))
			Expect(policy.Done()).To(BeFalse())
			Expect(policy.Duration()).To(Equal(64 * time.Second))
			Expect(policy.Done()).To(BeFalse())
			Expect(policy.Duration()).To(Equal(128 * time.Second))
			Expect(policy.Done()).To(BeFalse())
			Expect(policy.Duration()).To(Equal(256 * time.Second))
			Expect(policy.Done()).To(BeFalse())
			// Instead of being 512s in a normal exponential backoff, given the
			// timeout is 10m (600s), and 511s have passed now (2**0 + 2**1 +
			// ... 2**8 = 511), the next duration should be 600 - 511 = 89s:
			Expect(policy.Duration()).To(Equal(89 * time.Second))
			Expect(policy.Done()).To(BeTrue())
		})

		Describe("Reset", func() {
			It("resets the current policy so it can be re-used", func() {
				policy := retry.TimingOutExponentialBackoff{
					Timeout:  1 * time.Second,
					TimeUnit: time.Second,
				}
				Expect(policy.Done()).To(BeFalse())
				Expect(policy.Duration()).To(Equal(1 * time.Second))
				Expect(policy.Done()).To(BeTrue())
				policy.Reset()
				Expect(policy.Done()).To(BeFalse())
				Expect(policy.Duration()).To(Equal(1 * time.Second))
				Expect(policy.Done()).To(BeTrue())
			})
		})

		Describe("Clone", func() {
			It("clones the current policy so both can be used concurrently", func() {
				policy := retry.TimingOutExponentialBackoff{
					Timeout:  1 * time.Second,
					TimeUnit: time.Second,
				}
				Expect(policy.Done()).To(BeFalse())
				Expect(policy.Duration()).To(Equal(1 * time.Second))
				Expect(policy.Done()).To(BeTrue())
				clone := policy.Clone()
				addrPolicy := fmt.Sprintf("%p", &policy)
				addrClone := fmt.Sprintf("%p", &clone)
				Expect(addrPolicy).To(Not(Equal(addrClone)))
				// The original policy is still in a terminal state:
				Expect(policy.Done()).To(BeTrue())
				// The cloned policy can be used:
				Expect(clone.Done()).To(BeFalse())
				Expect(clone.Duration()).To(Equal(1 * time.Second))
				Expect(clone.Done()).To(BeTrue())
			})
		})
	})
})
