package iam

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("iam", func() {
	Describe("ARN", func() {
		It("determines if it is a user", func() {
			role := "arn:aws:iam::123456:role/testing"
			arn, err := Parse(role)
			Expect(err).ToNot(HaveOccurred())

			Expect(arn.IsUser()).To(BeFalse())
			Expect(arn.IsRole()).To(BeTrue())
		})
		It("determines if it is a role", func() {
			user := "arn:aws:iam::123456:user/testing"
			arn, err := Parse(user)
			Expect(err).ToNot(HaveOccurred())

			Expect(arn.IsUser()).To(BeTrue())
			Expect(arn.IsRole()).To(BeFalse())
		})
	})
})
