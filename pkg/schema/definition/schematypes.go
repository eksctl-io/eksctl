package definition

//nolint:golint,goconst
func setTypeOrRef(def *Definition, typeName string) {
	switch typeName {
	case "string":
		def.Type = "string"
	case "bool":
		def.Type = "boolean"
	case "int", "int64", "int32":
		def.Type = "integer"
	case "float64":
		def.Type = "number"
	case "byte":
		def.Type = "string" // TODO mediaEncoding
	default:
		def.Ref = DefPrefix + typeName
	}
}

func setDefaultForNonPointerType(def *Definition, typeName string) {
	switch typeName {
	case "bool":
		def.Default = "false"
	case "string":
		def.Default = ""
	}
}
