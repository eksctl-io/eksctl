package definition

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HandleComment", func() {
	It("interprets type override", func() {
		def := &Definition{}
		comment := `Struct holds some info
		Schema type is ` + "`string`"
		noderive, err := HandleComment("Struct", comment, def, false)
		Expect(err).ToNot(HaveOccurred())
		Expect(noderive).To(BeTrue())
		Expect(def.Description).To(Equal("holds some info"))
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
