package definition

import (
	"go/ast"
	"reflect"
	"strings"
)

// GetFieldTag gets the StructTag for a field
func GetFieldTag(field *ast.Field) reflect.StructTag {
	if field.Tag == nil {
		return ""
	}
	tag := strings.Replace(field.Tag.Value, "`", "", -1)
	return reflect.StructTag(tag)
}

// JSONPropName returns the name for marshaling to/from json
func JSONPropName(tag reflect.StructTag) string {
	jsonField := tag.Get("json")

	return strings.Split(jsonField, ",")[0]
}

// IsRequired checks whether the field has been marked required
func IsRequired(tag reflect.StructTag) bool {
	jsonField := tag.Get("jsonschema")
	return strings.Contains(jsonField, "required")
}
