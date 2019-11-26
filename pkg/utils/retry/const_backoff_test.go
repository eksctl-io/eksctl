package retry_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/utils/retry"
)

var _ = Describe("retry", func() {
	Describe("ConstantBackoff", func() {
		It("generates a sequence of constant durations", func() {
			policy := retry.ConstantBackoff{
				MaxRetries: 5,
				Time:       10,
				TimeUnit:   time.Second,
			}
			Expect(policy.Done()).To(BeFalse())
			Expect(policy.Duration()).To(Equal(10 * time.Second))
			Expect(policy.Done()).To(BeFalse())
			Expect(policy.Duration()).To(Equal(10 * time.Second))
			Expect(policy.Done()).To(BeFalse())
			Expect(policy.Duration()).To(Equal(10 * time.Second))
			Expect(policy.Done()).To(BeFalse())
			Expect(policy.Duration()).To(Equal(10 * time.Second))
			Expect(policy.Done()).To(BeFalse())
			Expect(policy.Duration()).To(Equal(10 * time.Second))
			Expect(policy.Done()).To(BeTrue())
		})

		Describe("Reset", func() {
			It("resets the current policy so it can be re-used", func() {
				policy := retry.ConstantBackoff{
					MaxRetries: 1,
					Time:       1,
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
				policy := retry.ConstantBackoff{
					MaxRetries: 1,
					Time:       1,
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
})
