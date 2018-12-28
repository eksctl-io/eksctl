package utils

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/kubicorn/kubicorn/pkg/namer"
)

func useNameOrGenerate(a, b string, generate func() string) string {
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
	return useNameOrGenerate(a, b, func() string {
		return fmt.Sprintf("%s-%d", namer.RandomName(), time.Now().Unix())
	})
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

const (
	randNodeGroupNameLength     = 8
	randNodeGroupNameComponents = "abcdef0123456789"
)

// NodeGroupName generates a name string when a and b are empty strings.
// If either a or b are non-empty, it returns whichever is non-empty.
// If neither a nor b are empty, it returns empty name, to indicate
// ambiguous usage.
// It uses a different naming scheme from ClusterName, so that users can
// easily distinguish a cluster name from nodegroup name.
func NodeGroupName(a, b string) string {
	return useNameOrGenerate(a, b, func() string {
		name := make([]byte, randNodeGroupNameLength)
		for i := 0; i < randNodeGroupNameLength; i++ {
			name[i] = randNodeGroupNameComponents[r.Intn(len(randNodeGroupNameComponents))]
		}
		return fmt.Sprintf("ng-%s", string(name))
	})
}
