package iam

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("iam", func() {
	Describe("Identity", func() {
		It("determines two Identity values are the same", func() {
			a, err := NewIdentity("arn:aws:iam::123456:role/testing", "user", []string{"system:masters", "system:bootstrap"})
			Expect(err).NotTo(HaveOccurred())
			b, err := NewIdentity("arn:aws:iam::123456:role/testing", "user", []string{"system:masters", "system:bootstrap"})
			Expect(err).NotTo(HaveOccurred())

			Expect(CompareIdentity(a, b)).To(BeTrue())
		})
		It("determines that Identity values with same groups in different order are the same", func() {
			a, err := NewIdentity("arn:aws:iam::123456:role/testing", "user", []string{"system:bootstrap", "system:masters"})
			Expect(err).NotTo(HaveOccurred())
			b, err := NewIdentity("arn:aws:iam::123456:role/testing", "user", []string{"system:masters", "system:bootstrap"})
			Expect(err).NotTo(HaveOccurred())

			Expect(CompareIdentity(a, b)).To(BeTrue())
		})
		It("determines that Identity values with different Arn not the same", func() {
			roleA := "arn:aws:iam::123456:role/testing-a"
			roleB := "arn:aws:iam::123456:role/testing-b"
			a, err := NewIdentity(roleA, "user", []string{"system:masters", "system:bootstrap"})
			Expect(err).NotTo(HaveOccurred())
			b, err := NewIdentity(roleB, "user", []string{"system:masters", "system:bootstrap"})
			Expect(err).NotTo(HaveOccurred())

			Expect(CompareIdentity(a, b)).To(BeFalse())
		})
		It("determines that Identity values with different user not the same", func() {
			a, err := NewIdentity("arn:aws:iam::123456:role/testing", "userA", []string{"system:masters", "system:bootstrap"})
			Expect(err).NotTo(HaveOccurred())
			b, err := NewIdentity("arn:aws:iam::123456:role/testing", "userB", []string{"system:masters", "system:bootstrap"})
			Expect(err).NotTo(HaveOccurred())

			Expect(CompareIdentity(a, b)).To(BeFalse())
		})
		It("determines that Identity values with different group lens not the same", func() {
			a, err := NewIdentity("arn:aws:iam::123456:role/testing", "user", []string{"system:masters", "system:bootstrap"})
			Expect(err).NotTo(HaveOccurred())
			b, err := NewIdentity("arn:aws:iam::123456:role/testing", "user", []string{"system:masters"})
			Expect(err).NotTo(HaveOccurred())

			Expect(CompareIdentity(a, b)).To(BeFalse())
		})
		It("determines that Identity values with different groups not the same", func() {
			a, err := NewIdentity("arn:aws:iam::123456:role/testing", "user", []string{"system:masters", "system:bootstrap"})
			Expect(err).NotTo(HaveOccurred())
			b, err := NewIdentity("arn:aws:iam::123456:role/testing", "user", []string{"system:masters", "system:different"})
			Expect(err).NotTo(HaveOccurred())

			Expect(CompareIdentity(a, b)).To(BeFalse())
		})
	})
})
