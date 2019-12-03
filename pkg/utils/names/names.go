package names

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

// ForCluster generates a name string when a and b are empty strings.
// If either a or b are non-empty, it returns whichever is non-empty.
// If neither a nor b are empty, it returns empty name, to indicate
// ambiguous usage.
func ForCluster(a, b string) string {
	return useNameOrGenerate(a, b, func() string {
		return fmt.Sprintf("%s-%d", namer.RandomName(), time.Now().Unix())
	})
}

// ForNodeGroup generates a name string when a and b are empty strings.
// If either a or b are non-empty, it returns whichever is non-empty.
// If neither a nor b are empty, it returns empty name, to indicate
// ambiguous usage.
// It uses a different naming scheme from ClusterName, so that users can
// easily distinguish a cluster name from nodegroup name.
func ForNodeGroup(a, b string) string {
	return useNameOrGenerate(a, b, func() string {
		return fmt.Sprintf("ng-%s", RandomName(randNodeGroupNameLength, randNodeGroupNameComponents))
	})
}

// ForFargateProfile returns the provided name if non-empty, or else generates
// a random name matching: fp-[abcdef0123456789]{8}
func ForFargateProfile(name string) string {
	if name != "" {
		return name
	}
	length := 8                 // Length of automatically generated Fargate profile names.
	chars := "abcdef0123456789" // Characters used to generate Fargate profile names.
	return fmt.Sprintf("fp-%s", RandomName(length, chars))
}

// useNameOrGenerate picks one of the provided strings or generates a
// new one using the provided generate function
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

// RandomName generates a string of the provided length by randomly selecting
// characters in the provided set.
func RandomName(length int, chars string) string {
	randomName := make([]byte, length)
	for i := 0; i < length; i++ {
		randomName[i] = chars[r.Intn(len(chars))]
	}
	return string(randomName)
}
