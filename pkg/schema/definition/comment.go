package definition

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/russross/blackfriday/v2"
	"github.com/weaveworks/eksctl/pkg/schema/importer"
)

var (
	regexpDefaults      = regexp.MustCompile("(.*)Defaults to `(.*)`")
	regexpExample       = regexp.MustCompile("(.*)For example: `(.*)`")
	typeOverridePattern = regexp.MustCompile("(.*)Schema type is `([a-zA-Z]+)`")
	pTags               = regexp.MustCompile("(<p>)|(</p>)")

	// patterns for enum-type values
	enumValuePattern     = "^[ \t]*`(?P<name>[^`]+)`([ \t]*\\(default\\))?(?:(: .*)|,)?$"
	regexpEnumDefinition = regexp.MustCompile("(?m).*Valid [a-z]+ are:((\\n" + enumValuePattern + ")*)")
	regexpEnumValues     = regexp.MustCompile("(?m)" + enumValuePattern)
)

func getTypeName(rawName string) string {
	splits := strings.Split(rawName, ".")
	return splits[len(splits)-1]
}

// findLiteralValue does a lookup for constant values
func findLiteralValue(importer importer.Importer, name string) (string, error) {
	obj, ok := importer.FindPkgObj(name)
	if !ok {
		return "", errors.New("Not in package")
	}
	valueSpec, ok := obj.Decl.(*ast.ValueSpec)
	if !ok {
		return "", errors.New("obj must refer to a value")
	}
	basicLit, ok := valueSpec.Values[0].(*ast.BasicLit)
	if !ok {
		return "", errors.New("obj must have a literal value")
	}
	switch basicLit.Kind {
	case token.STRING, token.CHAR:
		str, err := strconv.Unquote(basicLit.Value)
		if err != nil {
			panic("Couldn't unquote basic literal of type STRING or CHAR")
		}
		return str, nil
	default:
		return basicLit.Value, nil
	}
}

// handleComment interprets as much as it can from the comment and saves this
// information in the Definition
func (dg *Generator) handleComment(rawName, comment string, def *Definition) (bool, error) {
	var noDerive bool
	name := getTypeName(rawName)
	if dg.Strict && name != "" {
		if !strings.HasPrefix(comment, name+" ") {
			return noDerive, errors.Errorf("comment should start with field name on field %s", name)
		}
	}

	// process enums before stripping out newlines
	if m := regexpEnumDefinition.FindStringSubmatch(comment); m != nil {
		enums := make([]string, 0)
		if n := regexpEnumValues.FindAllStringSubmatch(m[1], -1); n != nil {
			for _, matches := range n {
				rawVal := matches[1]
				isDefault := matches[2] != ""
				var enumVal string
				if literal, err := strconv.Unquote(rawVal); err == nil {
					enumVal = literal
				} else {
					val, err := findLiteralValue(dg.Importer, rawVal)
					if err != nil {
						return noDerive, errors.Wrapf(err, "couldn't resolve %s in package", rawVal)
					}
					enumVal = val
					comment = strings.ReplaceAll(comment, rawVal, fmt.Sprintf(`"%s"`, val))
				}
				if isDefault {
					def.Default = enumVal
				}
				enums = append(enums, enumVal)
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

	if dg.Strict && name != "" {
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
