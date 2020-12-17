package definition

import (
	"go/parser"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/russross/blackfriday/v2"
)

var (
	regexpDefaults      = regexp.MustCompile("(.*)Defaults to `(.*)`")
	regexpRequired      = regexp.MustCompile("(?m)^Required$")
	regexpPlusRequired  = regexp.MustCompile(`(?m)^\+required$`)
	regexpExample       = regexp.MustCompile("(.*)For example: `(.*)`")
	typeOverridePattern = regexp.MustCompile("(.*)Schema type is `(.*)`")
	oneOfEntry          = "`([^`]+)`,?[ \t]*"
	oneOfEntryPattern   = regexp.MustCompile(oneOfEntry)
	typeOneOfPattern    = regexp.MustCompile("(.*)Schema type is one of ((?:" + oneOfEntry + ")*)")
	pTags               = regexp.MustCompile("(<p>)|(</p>)")
)

type Meta struct {
	Required bool
	NoDerive bool
}

// handleComment interprets as much as it can from the comment and saves this
// information in the Definition
func (dg *Generator) handleComment(rawName, comment string, def *Definition) (Meta, error) {
	var noDerive bool
	var required bool
	_, name := interpretReference(rawName)
	if dg.Strict && name != "" {
		if !strings.HasPrefix(comment, name+" ") {
			return Meta{}, errors.Errorf("comment should start with field name on field %s", name)
		}
	}

	enumInformation, err := handleEnumComments(dg.Importer, def, comment)
	if err != nil {
		return Meta{}, err
	}
	var synthesizedComment string
	if enumInformation != nil {
		comment = enumInformation.RemainingComment
		synthesizedComment = enumInformation.SynthesizedComment
	}

	// Extract requiredness
	if regexpPlusRequired.FindStringSubmatch(comment) != nil || regexpRequired.FindStringSubmatch(comment) != nil {
		required = true
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
		parsedDefault, err := parserAsValue(m[2])
		if err != nil {
			return Meta{}, errors.Wrapf(err, "couldn't parse default value from %v", m[2])
		}
		def.Default = parsedDefault
	}

	// Extract schema type, disabling derivation
	if m := typeOverridePattern.FindStringSubmatch(description); m != nil {
		description = strings.TrimSpace(m[1])
		noDerive = true
		expr, err := parser.ParseExpr(m[2])
		if err != nil {
			return Meta{}, errors.Wrapf(err, "couldn't parse type override %v", m[2])
		}
		overrideDef, _ := dg.newPropertyRef("", expr, "", false)
		*def = *overrideDef
	} else if m := typeOneOfPattern.FindStringSubmatch(description); m != nil {
		description = strings.TrimSpace(m[1])
		n := oneOfEntryPattern.FindAllStringSubmatch(m[2], -1)
		if n == nil {
			panic("error matching oneOf comment with regex")
		}
		defs := []*Definition{}
		for _, t := range n {
			schemaRef := t[1]
			expr, err := parser.ParseExpr(schemaRef)
			if err != nil {
				return Meta{}, errors.Wrapf(err, "couldn't parse `oneOf` type %v", schemaRef)
			}
			def, _ := dg.newPropertyRef("", expr, "", false)
			defs = append(defs, def)
		}
		def.OneOf = defs
	}

	// Extract example
	if m := regexpExample.FindStringSubmatch(description); m != nil {
		description = strings.TrimSpace(m[1])
		def.Examples = []string{m[2]}
	}

	// Remove type prefix
	description = removeTypeNameFromComment(name, description)

	if dg.Strict && name != "" {
		if description == "" {
			return Meta{}, errors.Errorf("no description on field %s", name)
		}
		if !strings.HasSuffix(description, ".") {
			return Meta{}, errors.Errorf("description should end with a dot on field %s", name)
		}
	}
	def.Description = joinIfNotEmpty(" ", description, synthesizedComment)

	// Convert to HTML
	html := string(blackfriday.Run([]byte(def.Description), blackfriday.WithNoExtensions()))
	def.HTMLDescription = strings.TrimSpace(pTags.ReplaceAllString(html, ""))
	return Meta{NoDerive: noDerive, Required: required}, nil
}

func removeTypeNameFromComment(name, description string) string {
	return regexp.MustCompile("^"+name+" (\\*.*\\* )?((is (the )?)|(are (the )?)|(lists ))?").ReplaceAllString(description, "$1")
}

// joinIfNotEmpty is sadly necessary
func joinIfNotEmpty(sep string, elems ...string) string {
	var nonEmptyElems = []string{}
	for _, e := range elems {
		if e != "" {
			nonEmptyElems = append(nonEmptyElems, e)
		}
	}
	return strings.Join(nonEmptyElems, sep)
}
