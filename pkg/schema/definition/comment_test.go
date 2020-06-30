package definition

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/schema/importer"
)

var _ = Describe("HandleComment", func() {
	It("interprets type override", func() {
		def := &Definition{}
		comment := `Struct holds some info
		Schema type is ` + "`string`"
		dummy := func(path string) (importer.PackageInfo, error) {
			return importer.PackageInfo{}, nil
		}
		dg := Generator{Strict: false, Importer: dummy}
		noderive, err := dg.handleComment("Struct", comment, def)
		Expect(err).ToNot(HaveOccurred())
		Expect(noderive).To(BeTrue())
		Expect(def.Description).To(Equal("holds some info"))
		Expect(def.Type).To(Equal("string"))
	})
})

var _ = Describe("interpretReference", func() {
	It("interprets root name", func() {
		pkg, name := interpretReference("SomeType")
		Expect(name).To(Equal("SomeType"))
		Expect(pkg).To(Equal(""))
	})
	It("interprets pkg name", func() {
		pkg, name := interpretReference("some/pkg.SomeType")
		Expect(name).To(Equal("SomeType"))
		Expect(pkg).To(Equal("some/pkg"))
	})
})
