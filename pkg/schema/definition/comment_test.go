package definition

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleJSONSchemaComment(t *testing.T) {
	def := &Definition{}
	comment := `Comment about struct
	+jsonschema noderive
	+jsonschema { "type": "string" }`
	noderive, remaining, err := HandleJSONSchemaComment(comment, def)
	assert.Nil(t, err)
	assert.True(t, noderive)
	assert.Equal(t, "string", def.Type)
	assert.Equal(t, "Comment about struct", remaining)
}

func TestGetTypeName(t *testing.T) {
	assert.Equal(t, "Thing", getTypeName("some/pkg.Thing"))
	assert.Equal(t, "Thing", getTypeName("Thing"))
}
