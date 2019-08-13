package random_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/weaveworks/eksctl/pkg/utils/random"
)

func TestString(t *testing.T) {
	assert.Equal(t, "", String(-1))
	assert.Equal(t, "", String(0))
	assert.NotEqual(t, String(10), String(10))

	randomString := String(10)
	assert.Len(t, randomString, 10)
	for _, c := range randomString {
		assert.Contains(t, "abcdefghijklmnopqrstuvwxyz0123456789", string(c))
	}
}
