package utils

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/kubicorn/kubicorn/pkg/namer"
)

const (
	randNodeGroupNameLength     = 8
	randNodeGroupNameComponents = "abcdef0123456789"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

// NodeGroupName generates a name string when a and b are empty strings.
// If either a or b are non-empty, it returns whichever is non-empty.
// If neither a nor b are empty, it returns empty name, to indicate
// ambiguous usage.
// It uses a different naming scheme from ClusterName, so that users can
// easily distinguish a cluster name from nodegroup name.
func NodeGroupName(a, b string) string {
	return UseNameOrGenerate(a, b, func() string {
		name := make([]byte, randNodeGroupNameLength)
		for i := 0; i < randNodeGroupNameLength; i++ {
			name[i] = randNodeGroupNameComponents[r.Intn(len(randNodeGroupNameComponents))]
		}
		return fmt.Sprintf("ng-%s", string(name))
	})
}

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
