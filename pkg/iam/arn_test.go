package iam_test

import (
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/weaveworks/eksctl/pkg/iam"
)

var _ = Describe("iam", func() {
	Describe("ARN", func() {
		var role string
		var bytesRole []byte
		var jsonRole string

		BeforeEach(func() {
			var err error
			role = "arn:aws:iam::123456:role/testing"
			bytesRole, err = json.Marshal(role)
			Expect(err).ToNot(HaveOccurred())
			jsonRole = string(bytesRole)
		})

		It("marshals to string", func() {
			arn, err := Parse(role)
			Expect(err).ToNot(HaveOccurred())

			out, err := json.Marshal(arn)
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(MatchJSON(jsonRole))
		})
		It("unmarshals from string", func() {
			var arn ARN
			err := json.Unmarshal(bytesRole, &arn)
			Expect(err).ToNot(HaveOccurred())
			Expect(arn.String()).To(Equal(role))
		})
	})
})
