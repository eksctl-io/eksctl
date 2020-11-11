package definition

import (
	"go/ast"

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
			return importer.PackageInfo{
				Pkg: &ast.Object{
					Data: &ast.Scope{
						Objects: map[string]*ast.Object{},
					},
					Kind: ast.Typ,
				},
			}, nil
		}
		dg := Generator{Strict: false, Importer: dummy}
		commentMeta, err := dg.handleComment("Struct", comment, def)
		Expect(err).ToNot(HaveOccurred())
		Expect(commentMeta.NoDerive).To(BeTrue())
		Expect(def.Description).To(Equal("holds some info"))
		Expect(def.Type).To(Equal("string"))
	})
	It("interprets +required", func() {
		def := &Definition{}
		comment := `Struct holds some info
+required`
		dummy := func(path string) (importer.PackageInfo, error) {
			return importer.PackageInfo{}, nil
		}
		dg := Generator{Strict: false, Importer: dummy}
		commentMeta, err := dg.handleComment("Struct", comment, def)
		Expect(err).ToNot(HaveOccurred())
		Expect(commentMeta.Required).To(BeTrue())
		Expect(def.Description).To(Equal("holds some info"))
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
