package definition

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/schema/importer"
)

// patterns for enum-type values
var (
	enumValuePattern     = "^[ \t]*`(?P<name>[^`]+)`([ \t]*\\(default\\))?(?:(: .*)|,)?$"
	regexpEnumReference  = regexp.MustCompile("(?m).*Valid [a-z]+ are `(.*)` [a-z]+")
	regexpEnumDefinition = regexp.MustCompile("(?m).*Valid [a-z]+ are:((\\n" + enumValuePattern + ")*)")
	regexpEnumValues     = regexp.MustCompile("(?m)" + enumValuePattern)
)

// findLiteralValue does a lookup for constant values
// at the moment it treats everything as a string
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

// enumCommentInformation holds interpreted information
type enumCommentInformation struct {
	Enum             []string
	Default          string
	EnumComment      string
	RemainingComment string
}

func interpretReferencedVariantComment(importer importer.Importer, valueSpec *ast.ValueSpec) (def bool, value string, comment string) {
	name := valueSpec.Names[0].Name
	value, err := findLiteralValue(importer, name)
	if err != nil {
		panic("Couldn't find a value from an Ident, impossible!")
	}
	enum := value
	if strings.Contains(valueSpec.Doc.Text(), "(default)") {
		def = true
	}
	// valueSpec comments have a trailing newline
	trimmed := strings.TrimSuffix(valueSpec.Doc.Text(), "\n")
	// Replace with the real literal
	generatedComment := strings.ReplaceAll(trimmed, valueSpec.Names[0].Name, fmt.Sprintf("`\"%s\"`", value))
	return def, enum, generatedComment
}

// interpretEnumComments handles interpreting enum information from comments
func interpretEnumComments(importer importer.Importer, comment string) (*enumCommentInformation, error) {
	// process enums before stripping out newlines
	var def string
	var enum = []string{}
	if m := regexpEnumReference.FindStringSubmatch(comment); m != nil {
		comment = strings.ReplaceAll(comment, m[0], "")
		// We need to generate a comment here
		var variantComments = []string{}
		pkgName, variantName := interpretReference(m[1])

		pkgInfo, err := importer(pkgName)
		if err != nil {
			return nil, err
		}
		variants, ok := pkgInfo.Variants[variantName]
		if !ok {
			return nil, errors.Errorf("Couldn't find %s in root package", variantName)
		}
		for _, valueSpec := range variants {
			isDefault, value, variantComment := interpretReferencedVariantComment(importer, valueSpec)

			enum = append(enum, value)
			if isDefault {
				def = value
			}
			variantComments = append(variantComments, variantComment)
		}
		joinedVariantComments := strings.Join(variantComments, ", ")
		enumComment := strings.Join(append([]string{"Valid variants are:"}, joinedVariantComments), " ") + "."
		return &enumCommentInformation{
			Enum:             enum,
			Default:          def,
			EnumComment:      enumComment,
			RemainingComment: comment,
		}, nil
	} else if m := regexpEnumDefinition.FindStringSubmatch(comment); m != nil {
		if n := regexpEnumValues.FindAllStringSubmatch(m[1], -1); n != nil {
			for _, matches := range n {
				rawVal := matches[1]
				isDefault := matches[2] != ""
				var enumVal string
				if literal, err := strconv.Unquote(rawVal); err == nil {
					enumVal = literal
				} else {
					val, err := findLiteralValue(importer, rawVal)
					if err != nil {
						return nil, errors.Wrapf(err, "couldn't resolve %s in package", rawVal)
					}
					enumVal = val
					comment = strings.ReplaceAll(comment, rawVal, fmt.Sprintf(`"%s"`, val))
				}
				if isDefault {
					def = enumVal
				}
				enum = append(enum, enumVal)
			}
		}
		return &enumCommentInformation{
			Enum:    enum,
			Default: def,
			// For now the generated enum comment is empty
			// because we leave the comment as is
			EnumComment:      "",
			RemainingComment: comment,
		}, nil
	}
	return nil, nil
}
