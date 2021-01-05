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
	enumValuePattern         = "^[ \t]*`(?P<name>[^`]+)`([ \t]*\\(default\\))?(?:(: .*)|,)?$"
	regexpItemEnumReference  = regexp.MustCompile("(?m).*Valid (entries|items) are `(.*)` [a-z]+")
	regexpItemEnumDefinition = regexp.MustCompile("(?m).*Valid (?:entries|items) are:((\\n" + enumValuePattern + ")*)")
	regexpEnumReference      = regexp.MustCompile("(?m).*Valid ([a-z]+) are `(.*)` [a-z]+")
	regexpEnumDefinition     = regexp.MustCompile("(?m).*Valid [a-z]+ are:((\\n" + enumValuePattern + ")*)")
	regexpEnumValues         = regexp.MustCompile("(?m)" + enumValuePattern)
)

// findLiteralFromString does a lookup for constant values
func findLiteralFromString(i importer.Importer, path, name string) (string, error) {
	obj, ok := i.SearchPackageForObj(path, name)
	if !ok {
		return "", errors.New("Not in package")
	}
	valueSpec, ok := obj.Decl.(*ast.ValueSpec)
	if !ok {
		return "", errors.New("obj must refer to a value")
	}
	return findLiteralValue(i, valueSpec.Values[0])
}

// findLiteralValue gives us a string literal from an expression
// it handles identifiers, selectors and literal strings
func findLiteralValue(i importer.Importer, value ast.Expr) (string, error) {
	switch v := value.(type) {
	case *ast.BasicLit:
		switch v.Kind {
		case token.STRING, token.CHAR:
			str, err := strconv.Unquote(v.Value)
			if err != nil {
				panic("Couldn't unquote basic literal of type STRING or CHAR")
			}
			return str, nil
		default:
			return v.Value, nil
		}
	case *ast.SelectorExpr:
		importPath, typeName, err := importer.ImportPathFromSelector(v)
		if err != nil {
			return "", err
		}
		return findLiteralFromString(i, importPath, typeName)
	case *ast.Ident:
		return findLiteralFromString(i, "", v.Name)
	default:
		return "", errors.Errorf("obj must have a literal value, instead has %T", v)
	}
}

// interpretReferencedVariantComment interprets the comments for _one_ variant
func interpretReferencedVariantComment(importer importer.Importer, valueSpec *ast.ValueSpec) (def bool, value string, comment string) {
	name := valueSpec.Names[0].Name
	value, err := findLiteralFromString(importer, "", name)
	if err != nil {
		panic(errors.Wrap(err, "Couldn't find a value from an Ident"))
	}
	enum := value
	if strings.Contains(valueSpec.Doc.Text(), "(default)") {
		def = true
	}
	// Replace with the real literal
	generatedComment := strings.ReplaceAll(valueSpec.Doc.Text(), valueSpec.Names[0].Name, "")
	// valueSpec comments have a trailing newline
	generatedComment = strings.TrimSpace(generatedComment)
	return def, enum, generatedComment
}

// findReferencedVariants searches for a GenDecl marked with this name
func findReferencedVariants(importer importer.Importer, variantDeclName string) (importer.Variants, error) {
	pkgName, variantName := interpretReference(variantDeclName)

	pkgInfo, err := importer(pkgName)
	if err != nil {
		return nil, err
	}
	variants, ok := pkgInfo.Variants[variantName]
	if !ok {
		return nil, errors.Errorf("Couldn't find %s in root package", variantName)
	}
	return variants, nil
}

// enumCommentInformation holds information interpreted from a comment string
type enumCommentInformation struct {
	Enum             []string
	Default          *string
	EnumComment      string
	RemainingComment string
}

func handleGenDeclReference(importer importer.Importer, comment string, match []string) (*enumCommentInformation, error) {
	var enum = []string{}
	var def *string
	// Drop the reference to the GenDecl from our comment
	entireComment := match[0]
	comment = strings.ReplaceAll(comment, entireComment, "")

	termForVariant := match[1]
	variantName := match[2]
	variants, err := findReferencedVariants(importer, variantName)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't find variant")
	}

	// We synthesize the comments for the enum variants into one description
	// see the tests
	var variantComments = []string{}

	// We take into account duplicated enum variants
	seen := make(map[string]int)
	for _, valueSpec := range variants {
		isDefault, value, variantComment := interpretReferencedVariantComment(importer, valueSpec)

		if indexSeen, valueSeen := seen[value]; !valueSeen {
			// The first time we see a value, append a new enum variant
			variantCommentWithValue := joinIfNotEmpty(" ", fmt.Sprintf("`\"%s\"`", value), variantComment)
			enum = append(enum, value)
			seen[value] = len(enum) - 1
			variantComments = append(variantComments, variantCommentWithValue)
		} else {
			// If we've seen this value before, append the comments
			variantComments[indexSeen] = joinIfNotEmpty(" ", variantComments[indexSeen], variantComment)
		}
		if isDefault {
			def = &value
		}
	}

	// Add some commas, a nice prefix and a period
	joinedVariantComments := strings.Join(variantComments, ", ")
	prefix := fmt.Sprintf("Valid %s are:", termForVariant)
	enumComment := strings.Join(append([]string{prefix}, joinedVariantComments), " ") + "."

	return &enumCommentInformation{
		Enum:    enum,
		Default: def,
		// Our synthesized comment for this enum
		EnumComment:      enumComment,
		RemainingComment: comment,
	}, nil

}

func handleVariantList(importer importer.Importer, comment string, match []string) (*enumCommentInformation, error) {
	var enum = []string{}
	var def *string
	for _, match := range regexpEnumValues.FindAllStringSubmatch(match[1], -1) {
		rawVal := match[1]
		isDefault := match[2] != ""
		var value string
		if literal, err := strconv.Unquote(rawVal); err == nil {
			value = literal
		} else {
			val, err := findLiteralFromString(importer, "", rawVal)
			if err != nil {
				return nil, errors.Wrapf(err, "couldn't resolve %s in package", rawVal)
			}
			value = val
			comment = strings.ReplaceAll(comment, rawVal, fmt.Sprintf(`"%s"`, val))
		}
		if isDefault {
			def = &value
		}
		enum = append(enum, value)
	}
	return &enumCommentInformation{
		Enum:    enum,
		Default: def,
		// Our synthesized comment for this enum
		EnumComment:      "",
		RemainingComment: comment,
	}, nil
}

// CommentInformation holds information interpreted from a comment string
type CommentInformation struct {
	SynthesizedComment string
	RemainingComment   string
}

func applyEnumComment(def *Definition, info *enumCommentInformation) (*CommentInformation, error) {
	if info.Default != nil {
		def.Default = *info.Default
	}
	def.Enum = info.Enum
	return &CommentInformation{
		SynthesizedComment: info.EnumComment,
		RemainingComment:   info.RemainingComment,
	}, nil
}

// handleEnumComments handles interpreting enum information from comments
func handleEnumComments(importer importer.Importer, def *Definition, comment string) (*CommentInformation, error) {
	// If this comments refers to a GenDecl of constants
	if m := regexpItemEnumReference.FindStringSubmatch(comment); m != nil {
		info, err := handleGenDeclReference(importer, comment, m)
		if err != nil {
			return nil, err
		}
		if def.Items == nil {
			return nil, errors.Errorf("can't add enum to items for non-array definition %s", def.Type)
		}
		return applyEnumComment(def.Items, info)
	} else if m := regexpEnumReference.FindStringSubmatch(comment); m != nil {
		info, err := handleGenDeclReference(importer, comment, m)
		if err != nil {
			return nil, err
		}
		return applyEnumComment(def, info)
	} else if m := regexpItemEnumDefinition.FindStringSubmatch(comment); m != nil {
		info, err := handleVariantList(importer, comment, m)
		if err != nil {
			return nil, err
		}
		if def.Items == nil {
			return nil, errors.Errorf("can't add enum to items for non-array definition %s", def.Type)
		}
		return applyEnumComment(def.Items, info)
	} else if m := regexpEnumDefinition.FindStringSubmatch(comment); m != nil {
		info, err := handleVariantList(importer, comment, m)
		if err != nil {
			return nil, err
		}
		return applyEnumComment(def, info)
	}
	// We don't have an enum
	return nil, nil
}
