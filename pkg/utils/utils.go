package utils

import (
	"hash/fnv"
	"regexp"
	"strings"

	"github.com/weaveworks/eksctl/pkg/utils/version"
)

var matchFirstCap = regexp.MustCompile("([0-9]+|[A-Z])")

// ToKebabCase turns a CamelCase string into a kebab-case string
func ToKebabCase(str string) string {
	kebab := matchFirstCap.ReplaceAllString(str, "-${1}")
	kebab = strings.TrimPrefix(kebab, "-")
	return strings.ToLower(kebab)
}

// IsMinVersion compares a given version number with a minimum one and returns true if
// version >= minimumVersion
func IsMinVersion(minimumVersion, versionString string) (bool, error) {
	return version.IsMinVersion(minimumVersion, versionString)
}

// CompareVersions compares two version strings with the usual conventions:
// returns 0 if a == b
// returns 1 if a > b
// returns -1 if a < b
func CompareVersions(a, b string) (int, error) {
	return version.CompareVersions(a, b)
}

// FnvHash computes the hash of a string using the Fowler–Noll–Vo hash function
func FnvHash(s string) []byte {
	fnvHash := fnv.New32a()
	fnvHash.Write([]byte(s))
	return fnvHash.Sum(nil)
}

func StringPtr(s string) *string {
	return &s
}

func BoolPtr(b bool) *bool {
	return &b
}

func IntPtr(i int) *int {
	return &i
}
