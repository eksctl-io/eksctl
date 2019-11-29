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

			It("fails when the profile's name starts with eks-", func() {
				options := fargate.CreateOptions{
					Options: fargate.Options{
						ProfileName: "eks-foo",
					},
					ProfileSelectorNamespace: "default",
				}
				err := options.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid Fargate profile: name should NOT start with \"eks-\""))
			})

			It("passes on randomly generated input", func() {
				options := fargate.CreateOptions{}
				err := faker.FakeData(&options)
				Expect(err).To(Not(HaveOccurred()))
				err = options.Validate()
				Expect(err).To(Not(HaveOccurred()))
			})
		})

		Describe("ToFargateProfile", func() {
			It("creates a FargateProfile DTO object from this CreateOptions object's values", func() {
				options := fargate.CreateOptions{
					Options: fargate.Options{
						ProfileName: "default",
					},
					ProfileSelectorNamespace: "development",
					ProfileSelectorLabels: map[string]string{
						"env": "dev",
					},
				}
				profile := options.ToFargateProfile()
				Expect(profile).To(Not(BeNil()))
				Expect(profile.Validate()).To(Not(HaveOccurred()))
				Expect(profile.Name).To(Equal("default"))
				Expect(profile.Selectors).To(HaveLen(1))
				selector := profile.Selectors[0]
				Expect(selector.Namespace).To(Equal("development"))
				Expect(selector.Labels).To(HaveLen(1))
				Expect(selector.Labels).To(HaveKeyWithValue("env", "dev"))
			})
		})
	})
})
