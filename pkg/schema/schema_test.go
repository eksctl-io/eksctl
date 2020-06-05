package schema_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	. "github.com/weaveworks/eksctl/pkg/schema"
	"github.com/weaveworks/eksctl/pkg/schema/definition"
)

const configType = "Config"

var _ = Describe("GenerateSchema", func() {
	var schema Schema
	configDef := func() *definition.Definition {
		return schema.Definitions[configType]
	}
	BeforeSuite(func() {
		var err error
		schema, err = GenerateSchema("../../pkg/schema", "test", configType, false)
		Expect(err).NotTo(HaveOccurred())
	})
	It("handles the top level definition", func() {
		props := []string{"num", "option", "pointeroption", "packageoption", "aliasedint", "unknown", "other"}
		expected := Fields{
			"PreferredOrder":       Equal(props),
			"AdditionalProperties": Equal(false),
			"Description":          Equal("describes some settings for _some_ things"),
			"HTMLDescription":      Equal("describes some settings for <em>some</em> things"),
		}
		Expect(*configDef()).To(MatchFields(IgnoreExtras, expected))
	})
	It("handles primitive types", func() {
		expected := definition.Definition{
			Type:            "integer",
			Description:     "This number describes the number of subthings",
			HTMLDescription: "This number describes the number of subthings",
		}
		Expect(configDef().Properties).To(HaveKey("num"))
		Expect(*configDef().Properties["num"]).To(BeEquivalentTo(expected))
	})
	It("handles type aliases", func() {
		expected := definition.Definition{
			Type:            "integer",
			Description:     "just an int",
			HTMLDescription: "just an int",
		}
		Expect(*schema.Definitions["Alias"]).To(BeEquivalentTo(expected))
		ref := definition.Definition{
			Ref: "#/definitions/Alias",
		}
		Expect(configDef().Properties).To(HaveKey("aliasedint"))
		Expect(*configDef().Properties["aliasedint"]).To(BeEquivalentTo(ref))
	})
	It("handles interface{}", func() {
		expected := definition.Definition{}
		Expect(configDef().Properties).To(HaveKey("unknown"))
		Expect(*configDef().Properties["unknown"]).To(BeEquivalentTo(expected))
	})
	It("handles maps", func() {
		expected := definition.Definition{
			Type:                 "object",
			Default:              "{}",
			AdditionalProperties: &definition.Definition{Type: "string", Default: ""},
		}
		Expect(configDef().Properties).To(HaveKey("other"))
		Expect(*configDef().Properties["other"]).To(BeEquivalentTo(expected))
	})
	It("finds referenced structs", func() {
		When("directly referenced", func() {
			expected := definition.Definition{
				Properties: map[string]*definition.Definition{
					"kind": {
						Type:    "string",
						Default: "",
					},
				},
				PreferredOrder:       []string{"kind"},
				AdditionalProperties: false,
				Description:          "describes a sub configuration of the Config",
				HTMLDescription:      "describes a sub configuration of the Config",
			}
			Expect(*schema.Definitions["DirectType"]).To(BeEquivalentTo(expected))
			ref := definition.Definition{
				Ref: "#/definitions/DirectType",
			}
			Expect(configDef().Properties).To(HaveKey("option"))
			Expect(*configDef().Properties["option"]).To(BeEquivalentTo(ref))
		})
		When("referenced by pointer", func() {
			expected := definition.Definition{
				Properties: map[string]*definition.Definition{
					"kind": {
						Type:    "string",
						Default: "",
					},
				},
				PreferredOrder:       []string{"kind"},
				AdditionalProperties: false,
				Description:          "describes a sub configuration of the Config",
				HTMLDescription:      "describes a sub configuration of the Config",
			}
			Expect(*schema.Definitions["PointerType"]).To(BeEquivalentTo(expected))
			ref := definition.Definition{
				Ref: "#/definitions/PointerType",
			}
			Expect(configDef().Properties).To(HaveKey("pointeroption"))
			Expect(*configDef().Properties["pointeroption"]).To(BeEquivalentTo(ref))
		})
		When("referenced by package", func() {
			expected := definition.Definition{
				Properties: map[string]*definition.Definition{
					"kind": {
						Type:    "string",
						Default: "",
					},
				},
				PreferredOrder:       []string{"kind"},
				AdditionalProperties: false,
				Description:          "describes a sub configuration of the Config",
				HTMLDescription:      "describes a sub configuration of the Config",
			}
			ref := "github.com|weaveworks|eksctl|pkg|schema|test|subpkg.PackageType"
			Expect(*schema.Definitions[ref]).To(BeEquivalentTo(expected))
			refDef := definition.Definition{
				Ref: "#/definitions/" + ref,
			}
			Expect(configDef().Properties).To(HaveKey("packageoption"))
			Expect(*configDef().Properties["packageoption"]).To(BeEquivalentTo(refDef))
		})
	})
})
