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
		props := []string{
			"num", "option", "pointeroption", "packageoption", "aliasedint", "unknown", "other", "version", "kind", "kinds", "sumType",
		}
		required := []string{"option"}
		expected := Fields{
			"Required":             Equal(required),
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
			Description:     "describes the number of subthings",
			HTMLDescription: "describes the number of subthings",
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
			AdditionalProperties: &definition.Definition{Type: "string"},
		}
		Expect(configDef().Properties).To(HaveKey("other"))
		Expect(*configDef().Properties["other"]).To(BeEquivalentTo(expected))
	})
	It("handles enums", func() {
		expected := definition.Definition{
			Type:            "string",
			Default:         "X",
			Enum:            []string{"X", "Y", "2.0", "2"},
			Description:     "Determines the version of the main thing. Valid variants are: `\"X\"` (default): This is the right option, `\"Y\"`: Will be deprecated, `\"2.0\"`, `\"2\"`",
			HTMLDescription: "Determines the version of the main thing. Valid variants are: <code>&quot;X&quot;</code> (default): This is the right option, <code>&quot;Y&quot;</code>: Will be deprecated, <code>&quot;2.0&quot;</code>, <code>&quot;2&quot;</code>",
		}
		Expect(configDef().Properties).To(HaveKey("version"))
		Expect(*configDef().Properties["version"]).To(BeEquivalentTo(expected))
	})
	It("handles enums by reference", func() {
		expected := definition.Definition{
			Type:            "string",
			Default:         "SecondKind",
			Enum:            []string{"FirstKind", "SecondKind", "SpecialKind"},
			Description:     "Tells us which kind of config. Valid variants are: `\"FirstKind\"` is legacy, `\"SecondKind\"` should be used (default) and this comment combines with secondKind, `\"SpecialKind\"` is from some other package.",
			HTMLDescription: "Tells us which kind of config. Valid variants are: <code>&quot;FirstKind&quot;</code> is legacy, <code>&quot;SecondKind&quot;</code> should be used (default) and this comment combines with secondKind, <code>&quot;SpecialKind&quot;</code> is from some other package.",
		}
		Expect(configDef().Properties).To(HaveKey("kind"))
		Expect(*configDef().Properties["kind"]).To(BeEquivalentTo(expected))
	})
	It("handles enums by reference in lists", func() {
		expected := definition.Definition{
			Type: "array",
			Items: &definition.Definition{
				Type:    "string",
				Default: "SecondKind",
				Enum:    []string{"FirstKind", "SecondKind", "SpecialKind"},
			},
			Description:     "Valid entries are: `\"FirstKind\"` is legacy, `\"SecondKind\"` should be used (default) and this comment combines with secondKind, `\"SpecialKind\"` is from some other package.",
			HTMLDescription: "Valid entries are: <code>&quot;FirstKind&quot;</code> is legacy, <code>&quot;SecondKind&quot;</code> should be used (default) and this comment combines with secondKind, <code>&quot;SpecialKind&quot;</code> is from some other package.",
		}
		Expect(configDef().Properties).To(HaveKey("kinds"))
		Expect(*configDef().Properties["kinds"]).To(BeEquivalentTo(expected))
	})
	It("finds referenced structs", func() {
		When("directly referenced", func() {
			expected := definition.Definition{
				Properties: map[string]*definition.Definition{
					"kind": {
						Type: "string",
					},
				},
				PreferredOrder:       []string{"kind"},
				AdditionalProperties: false,
				Description:          "describes a sub configuration of the Config",
				HTMLDescription:      "describes a sub configuration of the Config",
			}
			Expect(*schema.Definitions["DirectType"]).To(BeEquivalentTo(expected))
			ref := definition.Definition{
				Ref:             "#/definitions/DirectType",
				Description:     "An option",
				HTMLDescription: "An option",
			}
			Expect(configDef().Properties).To(HaveKey("option"))
			Expect(*configDef().Properties["option"]).To(BeEquivalentTo(ref))
		})
		When("referenced by pointer", func() {
			expected := definition.Definition{
				Properties: map[string]*definition.Definition{
					"kind": {
						Type: "string",
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
						Type: "string",
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
	It("handles oneOf", func() {
		refDef := definition.Definition{
			Ref: "#/definitions/SumType",
		}
		Expect(configDef().Properties).To(HaveKey("sumType"))
		Expect(*configDef().Properties["sumType"]).To(BeEquivalentTo(refDef))
		expected := definition.Definition{
			Properties: map[string]*definition.Definition{
				"type": {
					Type:            "string",
					Enum:            []string{"a", "b"},
					Description:     "Valid variants are: `\"a\"`: type A `\"b\"`: type B",
					HTMLDescription: "Valid variants are: <code>&quot;a&quot;</code>: type A <code>&quot;b&quot;</code>: type B",
				},
			},
			PreferredOrder: []string{"type"},
			OneOf: []*definition.Definition{
				{Ref: "#/definitions/SumTypeA"},
				{Ref: "#/definitions/SumTypeB"},
			},
		}
		Expect(*schema.Definitions["SumType"]).To(BeEquivalentTo(expected))
	})
})
