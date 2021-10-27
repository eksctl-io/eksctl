package v1alpha5

import (
	// For go:embed
	_ "embed"
)

//go:generate go run ../../../../cmd/schema assets/schema.json

//go:embed assets/schema.json
var SchemaJson string
