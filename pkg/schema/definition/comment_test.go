package definition

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HandleJSONSchemaComment", func() {
	It("works", func() {
		def := &Definition{}
		comment := `Comment about struct
		+jsonschema noderive
		+jsonschema { "type": "string" }`
		noderive, remaining, err := HandleJSONSchemaComment(comment, def)
		Expect(err).ToNot(HaveOccurred())
		Expect(noderive).To(BeTrue())
		Expect(remaining).To(Equal("Comment about struct"))
		Expect(def.Type).To(Equal("string"))
	})
})

var _ = Describe("GetTypeName", func() {
	It("handles imported types", func() {
		Expect(getTypeName("some/pkg.Thing")).To(Equal("Thing"))
	})
	It("handles in scope types", func() {
		Expect(getTypeName("Thing")).To(Equal("Thing"))
	})
})
