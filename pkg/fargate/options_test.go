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
			It("fails when the profile's name is empty", func() {
				options := fargate.Options{}
				err := options.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid Fargate profile: empty name"))
			})

			It("passes when the profile's name is not empty", func() {
				options := fargate.Options{
					ProfileName: "default",
				}
				err := options.Validate()
				Expect(err).To(Not(HaveOccurred()))
			})

			It("passes on randomly generated input", func() {
				options := fargate.Options{}
				err := faker.FakeData(&options)
				Expect(err).To(Not(HaveOccurred()))
				err = options.Validate()
				Expect(err).To(Not(HaveOccurred()))
			})
		})
	})

	Describe("CreateOptions", func() {
		Describe("Validate", func() {
			It("fails when the profile's selector namespace is empty", func() {
				options := fargate.CreateOptions{}
				err := options.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid Fargate profile: empty selector namespace"))
			})

			It("passes when the profile's selector namespace is not empty", func() {
				options := fargate.CreateOptions{
					ProfileSelectorNamespace: "default",
				}
				err := options.Validate()
				Expect(err).To(Not(HaveOccurred()))
			})

			It("passes on randomly generated input", func() {
				options := fargate.CreateOptions{}
				err := faker.FakeData(&options)
				Expect(err).To(Not(HaveOccurred()))
				err = options.Validate()
				Expect(err).To(Not(HaveOccurred()))
			})
		})
	})

	Describe("GetOrDefaultProfileName", func() {
		It("returns the provided name if non-empty", func() {
			name := fargate.GetOrDefaultProfileName("my-favourite-name")
			Expect(name).To(Equal("my-favourite-name"))
		})
		It("generates a random name otherwise", func() {
			name := fargate.GetOrDefaultProfileName("")
			Expect(name).To(MatchRegexp("fp-[abcdef0123456789]{8}"))
		})
	})
})
