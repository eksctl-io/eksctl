package definition

import "strings"

// interpretReference takes a literal identifier or selector and gives us a pkg
// and identifier
func interpretReference(ref string) (string, string) {
	splits := strings.Split(ref, ".")
	var pkg string
	if len(splits) > 1 {
		pkg = strings.Join(splits[:len(splits)-1], "")
	}
	return pkg, splits[len(splits)-1]
}
