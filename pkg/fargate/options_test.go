package fargate_test

import (
	"github.com/bxcodec/faker"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/fargate"
)

var _ = Describe("fargate", func() {
	Describe("Options", func() {
		Describe("Validate", func() {
			It("fails when profile name is empty", func() {
				options := fargate.Options{}
				err := options.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid Fargate profile: empty name"))
			})
			It("passes when profile name is not empty", func() {
				options := fargate.Options{}
				err := faker.FakeData(&options)
				Expect(err).To(Not(HaveOccurred()))
				err = options.Validate()
				Expect(err).To(Not(HaveOccurred()))
			})
		})
	})
})
