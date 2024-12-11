package goformation_test

import (
	// Note that this is a fork of the main repo: github.com/xeipuuv/gojsonschema
	// CloudFormation uses nested schema def references, which is currently broken
	// in the main repo: https://github.com/xeipuuv/gojsonschema/pull/146
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xeipuuv/gojsonschema"
)

var _ = Describe("Goformation-generated JSON schemas", func() {

	Context("with a valid CloudFormation template", func() {

		pwd, _ := os.Getwd()
		schemaLoader := NewReferenceLoader("file://" + pwd + "/schema/cloudformation.schema.json")

		It("should successfully validate the CloudFormation template", func() {
			documentLoader := NewReferenceLoader("file://" + pwd + "/test/json/valid-template.json")
			result, err := Validate(schemaLoader, documentLoader)

			Expect(err).To(BeNil())
			Expect(result).ShouldNot(BeNil())
			Expect(result.Errors()).Should(BeEmpty())
			Expect(result.Valid()).Should(BeTrue())
		})

		It("should successfully validate a template with resource attributes", func() {
			documentLoader := NewReferenceLoader("file://" + pwd + "/test/json/valid-template-resource-attributes.json")
			result, err := Validate(schemaLoader, documentLoader)

			Expect(err).To(BeNil())
			Expect(result).ShouldNot(BeNil())
			Expect(result.Errors()).Should(BeEmpty())
			Expect(result.Valid()).Should(BeTrue())
		})
	})

	Context("with a valid SAM template", func() {

		pwd, _ := os.Getwd()
		schemaLoader := NewReferenceLoader("file://" + pwd + "/schema/sam.schema.json")
		documentLoader := NewReferenceLoader("file://" + pwd + "/test/json/valid-sam-template.json")

		result, err := Validate(schemaLoader, documentLoader)
		It("should successfully validate the SAM template", func() {
			Expect(err).To(BeNil())
			Expect(result).ShouldNot(BeNil())
			Expect(result.Errors()).Should(BeEmpty())
			Expect(result.Valid()).Should(BeTrue())
		})
	})

	Context("with a valid template, but which contains attributes not yet supported by the generated schema", func() {

		pwd, _ := os.Getwd()
		schemaLoader := NewReferenceLoader("file://" + pwd + "/schema/cloudformation.schema.json")

		It("should report that a template with intrinsic functions is invalid", func() {
			documentLoader := NewReferenceLoader("file://" + pwd + "/test/json/valid-template-with-fns.json")
			result, err := Validate(schemaLoader, documentLoader)

			Expect(err).To(BeNil())
			Expect(result).ShouldNot(BeNil())
			Expect(result.Valid()).Should(BeFalse())
		})
	})

	Context("with an invalid CloudFormation template", func() {

		pwd, _ := os.Getwd()
		schemaLoader := NewReferenceLoader("file://" + pwd + "/schema/cloudformation.schema.json")

		It("should report that a template missing required root properties is invalid", func() {
			documentLoader := NewReferenceLoader("file://" + pwd + "/test/json/invalid-template-missing-resources.json")
			result, err := Validate(schemaLoader, documentLoader)
			Expect(err).To(BeNil())
			Expect(result).ShouldNot(BeNil())
			Expect(result.Valid()).Should(BeFalse())
		})

		It("should report that a template missing required resource properties is invalid", func() {
			documentLoader := NewReferenceLoader("file://" + pwd + "/test/json/invalid-template-missing-resource-properties.json")
			result, err := Validate(schemaLoader, documentLoader)
			Expect(err).To(BeNil())
			Expect(result).ShouldNot(BeNil())
			Expect(result.Valid()).Should(BeFalse())
		})

		It("should report that a template with missing single required resource property is invalid", func() {
			documentLoader := NewReferenceLoader("file://" + pwd + "/test/json/invalid-template-empty-resource-properties.json")
			result, err := Validate(schemaLoader, documentLoader)
			Expect(err).To(BeNil())
			Expect(result).ShouldNot(BeNil())
			Expect(result.Valid()).Should(BeFalse())
		})

		It("should report that a template with a non-alphanumeric resource key is invalid", func() {
			documentLoader := NewReferenceLoader("file://" + pwd + "/test/json/invalid-template-invalid-resource-name.json")
			result, err := Validate(schemaLoader, documentLoader)
			Expect(err).To(BeNil())
			Expect(result).ShouldNot(BeNil())
			Expect(result.Valid()).Should(BeFalse())
		})

		It("should report that a template with unknown root properties is invalid", func() {
			documentLoader := NewReferenceLoader("file://" + pwd + "/test/json/invalid-template-additional-property.json")
			result, err := Validate(schemaLoader, documentLoader)
			Expect(err).To(BeNil())
			Expect(result).ShouldNot(BeNil())
			Expect(result.Valid()).Should(BeFalse())
		})

		It("should report that a template with unknown resource type is invalid", func() {
			documentLoader := NewReferenceLoader("file://" + pwd + "/test/json/invalid-template-unknown-resource-type.json")
			result, err := Validate(schemaLoader, documentLoader)
			Expect(err).To(BeNil())
			Expect(result).ShouldNot(BeNil())
			Expect(result.Valid()).Should(BeFalse())
		})

		It("should report that a template with unknown resource type is invalid", func() {
			documentLoader := NewReferenceLoader("file://" + pwd + "/test/json/invalid-template-missing-resource-type.json")
			result, err := Validate(schemaLoader, documentLoader)
			Expect(err).To(BeNil())
			Expect(result).ShouldNot(BeNil())
			Expect(result.Valid()).Should(BeFalse())
		})

		It("should report that a template with unknown resource properties is invalid", func() {
			documentLoader := NewReferenceLoader("file://" + pwd + "/test/json/invalid-template-additional-resource-property.json")
			result, err := Validate(schemaLoader, documentLoader)
			Expect(err).To(BeNil())
			Expect(result).ShouldNot(BeNil())
			Expect(result.Valid()).Should(BeFalse())
		})

		It("should report that a template with an unknown subproperty property is invalid", func() {
			documentLoader := NewReferenceLoader("file://" + pwd + "/test/json/invalid-template-unknown-subproperty-property.json")
			result, err := Validate(schemaLoader, documentLoader)
			Expect(err).To(BeNil())
			Expect(result).ShouldNot(BeNil())
			Expect(result.Valid()).Should(BeFalse())
		})

		It("should report that a template with a wrong subproperty type is invalid", func() {
			documentLoader := NewReferenceLoader("file://" + pwd + "/test/json/invalid-template-invalid-resource-subproperty.json")
			result, err := Validate(schemaLoader, documentLoader)
			Expect(err).To(BeNil())
			Expect(result).ShouldNot(BeNil())
			Expect(result.Valid()).Should(BeFalse())
		})

		It("should report that a template with a subproperty that is missing a required property is invalid", func() {
			documentLoader := NewReferenceLoader("file://" + pwd + "/test/json/invalid-template-subproperty-missing-property.json")
			result, err := Validate(schemaLoader, documentLoader)
			Expect(err).To(BeNil())
			Expect(result).ShouldNot(BeNil())
			Expect(result.Valid()).Should(BeFalse())
		})
	})

	Context("with an invalid SAM template", func() {

		pwd, _ := os.Getwd()
		schemaLoader := NewReferenceLoader("file://" + pwd + "/schema/sam.schema.json")

		It("should report that a template missing the Transform property is invalid", func() {
			documentLoader := NewReferenceLoader("file://" + pwd + "/test/json/invalid-sam-template-no-transform.json")
			result, err := Validate(schemaLoader, documentLoader)
			Expect(err).To(BeNil())
			Expect(result).ShouldNot(BeNil())
			Expect(result.Valid()).Should(BeFalse())
		})

		It("should report that a template with a non-SAM Transform value is invalid", func() {
			documentLoader := NewReferenceLoader("file://" + pwd + "/test/json/invalid-sam-template-wrong-transform.json")
			result, err := Validate(schemaLoader, documentLoader)
			Expect(err).To(BeNil())
			Expect(result).ShouldNot(BeNil())
			Expect(result.Valid()).Should(BeFalse())
		})
	})
})
