package utils

import (
	"fmt"
	"time"

	"github.com/kubicorn/kubicorn/pkg/namer"
)

// UseNameOrGenerate picks one of the provided strings or generates a
// new one using the provided generate function
func UseNameOrGenerate(a, b string, generate func() string) string {
	if a != "" && b != "" {
		return ""
	}
	if a != "" {
		return a
	}
	if b != "" {
		return b
	}
	return generate()
}

// ClusterName generates a name string when a and b are empty strings.
// If either a or b are non-empty, it returns whichever is non-empty.
// If neither a nor b are empty, it returns empty name, to indicate
// ambiguous usage.
func ClusterName(a, b string) string {
	return UseNameOrGenerate(a, b, func() string {
		return fmt.Sprintf("%s-%d", namer.RandomName(), time.Now().Unix())
	})
}
