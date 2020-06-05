package definition

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
	blackfriday "github.com/russross/blackfriday/v2"
)

var (
	regexpDefaults      = regexp.MustCompile("(.*)Defaults to `(.*)`")
	regexpExample       = regexp.MustCompile("(.*)For example: `(.*)`")
	typeOverridePattern = regexp.MustCompile("(.*)Schema type is `([a-zA-Z]+)`")
	pTags               = regexp.MustCompile("(<p>)|(</p>)")

	// patterns for enum-type values
	enumValuePattern     = "^[ \t]*`(?P<name>[^`]+)`([ \t]*\\(default\\))?: .*$"
	regexpEnumDefinition = regexp.MustCompile("(?m).*Valid [a-z]+ are((\\n" + enumValuePattern + ")*)")
	regexpEnumValues     = regexp.MustCompile("(?m)" + enumValuePattern)
)

func getTypeName(rawName string) string {
	splits := strings.Split(rawName, ".")
	return splits[len(splits)-1]
}

// HandleComment interprets as much as it can from the comment and saves this
// information in the Definition
func HandleComment(rawName, comment string, def *Definition, strict bool) (bool, error) {
	var noDerive bool
	name := getTypeName(rawName)
	if strict && name != "" {
		if !strings.HasPrefix(comment, name+" ") {
			return noDerive, errors.Errorf("comment should start with field name on field %s", name)
		}
	}

	// process enums before stripping out newlines
	if m := regexpEnumDefinition.FindStringSubmatch(comment); m != nil {
		enums := make([]string, 0)
		if n := regexpEnumValues.FindAllStringSubmatch(m[1], -1); n != nil {
			for _, matches := range n {
				enums = append(enums, matches[1])
			}
			def.Enum = enums
		}
	}

	// Remove kubernetes-style annotations from comments
	description := strings.TrimSpace(
		strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ReplaceAll(comment, "+required", ""),
				"+optional", "",
			), "\n", " ",
		),
	)

	// Extract default value
	if m := regexpDefaults.FindStringSubmatch(description); m != nil {
		description = strings.TrimSpace(m[1])
		def.Default = m[2]
	}

	// Extract schema type, disabling derivation
	if m := typeOverridePattern.FindStringSubmatch(description); m != nil {
		description = strings.TrimSpace(m[1])
		noDerive = true
		def.Type = m[2]
	}

	// Extract example
	if m := regexpExample.FindStringSubmatch(description); m != nil {
		description = strings.TrimSpace(m[1])
		def.Examples = []string{m[2]}
	}

	// Remove type prefix
	description = regexp.MustCompile("^"+name+" (\\*.*\\* )?((is (the )?)|(are (the )?)|(lists ))?").ReplaceAllString(description, "$1")

	if strict && name != "" {
		if description == "" {
			return noDerive, errors.Errorf("no description on field %s", name)
		}
		if !strings.HasSuffix(description, ".") {
			return noDerive, errors.Errorf("description should end with a dot on field %s", name)
		}
	}
	def.Description = description

	// Convert to HTML
	html := string(blackfriday.Run([]byte(description), blackfriday.WithNoExtensions()))
	def.HTMLDescription = strings.TrimSpace(pTags.ReplaceAllString(html, ""))
	return noDerive, nil
}
