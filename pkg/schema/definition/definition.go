package definition

import (
	"go/ast"
	"go/token"

	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/schema/importer"
)

const (
	// DefPrefix is the JSON Schema prefix required in the definition map
	DefPrefix = "#/definitions/"
)

// Generator can create definitions from Exprs
type Generator struct {
	Strict      bool
	Definitions map[string]*Definition
	PkgScope    *ast.Scope
	Importer    importer.Importer
}

// newStructDefinition handles saving definitions for refs in the map
func (dg *Generator) newStructDefinition(name string, typeSpec ast.Expr, comment string) *Definition {
	def := Definition{}
	noDerive, err := HandleComment(name, comment, &def, dg.Strict)
	if err != nil {
		panic(err)
	}
	if noDerive {
		return &def
	}
	structType, ok := typeSpec.(*ast.StructType)
	if !ok {
		panic(errors.Errorf("Cannot handle non-struct TypeSpec %s", name))
	}
	for _, field := range structType.Fields.List {
		tag := GetFieldTag(field)
		fieldName := JSONPropName(tag)

		if def.Properties == nil {
			def.Properties = make(map[string]*Definition)
		}

		var required []string
		var preferredOrder []string
		var properties map[string]*Definition
		if len(field.Names) == 0 {
			// We have to handle an embedded field, get its definition
			// and deconstruct it into this def
			ref := dg.newPropertyRef("", field.Type, field.Doc.Text(), true)
			properties = ref.Properties
			preferredOrder = ref.PreferredOrder
			required = ref.Required
		} else {
			if fieldName == "" {
				// private field
				continue
			}
			required = []string{}
			if IsRequired(tag) {
				required = []string{fieldName}
			}

			preferredOrder = []string{fieldName}
			properties = map[string]*Definition{
				fieldName: dg.newPropertyRef(field.Names[0].Name, field.Type, field.Doc.Text(), false),
			}
			def.AdditionalProperties = false
		}

		def.PreferredOrder = append(def.PreferredOrder, preferredOrder...)
		for k, v := range properties {
			def.Properties[k] = v
		}
		def.Required = append(def.Required, required...)
	}
	return &def
}

// newPropertyRef creates a new JSON schema Definition
func (dg *Generator) newPropertyRef(referenceName string, t ast.Expr, comment string, inline bool) *Definition {
	var def *Definition

	var refTypeName string
	var refTypeSpec *ast.TypeSpec

	switch tt := t.(type) {
	case *ast.Ident:
		typeName := tt.Name
		if obj, ok := dg.PkgScope.Objects[typeName]; ok {
			// If we have a declared type behind our ident, add it
			refTypeName, refTypeSpec = typeName, obj.Decl.(*ast.TypeSpec)
		}
		def = &Definition{}
		setTypeOrRef(def, typeName)
		setDefaultForNonPointerType(def, typeName)

	case *ast.StarExpr:
		def = dg.newPropertyRef(referenceName, tt.X, comment, inline)
		def.Default = nil

	case *ast.SelectorExpr:
		var err error
		refTypeName, refTypeSpec, err = dg.Importer.FindImportedTypeSpec(tt)
		if err != nil {
			panic(errors.Wrapf(err, "Couldn't import type from identifier"))
		}
		def = &Definition{}
		setTypeOrRef(def, refTypeName)

	case *ast.ArrayType:
		def = &Definition{
			Type:  "array",
			Items: dg.newPropertyRef("", tt.Elt, "", inline),
		}

	case *ast.MapType:
		def = &Definition{
			Type:                 "object",
			Default:              "{}",
			AdditionalProperties: dg.newPropertyRef("", tt.Value, "", inline),
		}

	case *ast.StructType:
		return dg.newStructDefinition(referenceName, t, comment)

	case *ast.InterfaceType:
		// Only `interface{}` is supported
		def = &Definition{}

	default:
		panic(errors.Errorf("Unexpected type %v for %s", t, referenceName))
	}

	// Add a new definition if necessary
	if refTypeSpec != nil {
		structDef := dg.newPropertyRef(refTypeName, refTypeSpec.Type, refTypeSpec.Doc.Text(), inline)
		// If we're inlining this, we want the struct definition, not the ref
		// and we also don't need to save it in our definitions
		if inline {
			return structDef
		}
		dg.Definitions[refTypeName] = structDef
	}

	_, err := HandleComment(referenceName, comment, def, dg.Strict)
	if err != nil {
		panic(err)
	}

	return def
}

// CollectDefinitionsFromStruct gets a complete definition for the root object
func (dg *Generator) CollectDefinitionsFromStruct(root string) {
	rootIdent := ast.Ident{
		NamePos: token.NoPos,
		Name:    root,
	}
	_ = dg.newPropertyRef(root, &rootIdent, "", false)
}
