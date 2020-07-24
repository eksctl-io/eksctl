package definition

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/russross/blackfriday/v2"
)

var (
	regexpDefaults      = regexp.MustCompile("(.*)Defaults to `(.*)`")
	regexpExample       = regexp.MustCompile("(.*)For example: `(.*)`")
	typeOverridePattern = regexp.MustCompile("(.*)Schema type is `([a-zA-Z]+)`")
	pTags               = regexp.MustCompile("(<p>)|(</p>)")
)

// handleComment interprets as much as it can from the comment and saves this
// information in the Definition
func (dg *Generator) handleComment(rawName, comment string, def *Definition) (bool, error) {
	var noDerive bool
	_, name := interpretReference(rawName)
	if dg.Strict && name != "" {
		if !strings.HasPrefix(comment, name+" ") {
			return noDerive, errors.Errorf("comment should start with field name on field %s", name)
		}
	}

	enumInformation, err := interpretEnumComments(dg.Importer, comment)
	if err != nil {
		return noDerive, err
	}
	var enumComment string
	if enumInformation != nil {
		def.Default = enumInformation.Default
		def.Enum = enumInformation.Enum
		comment = enumInformation.RemainingComment
		enumComment = enumInformation.EnumComment
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
			return noDerive, errors.Wrapf(err, "couldn't parse default value from %v", m[2])
		}
		def.Default = parsedDefault
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
	description = removeTypeNameFromComment(name, description)

	if dg.Strict && name != "" {
		if description == "" {
			return noDerive, errors.Errorf("no description on field %s", name)
		}
		if !strings.HasSuffix(description, ".") {
			return noDerive, errors.Errorf("description should end with a dot on field %s", name)
		}
	}
	def.Description = joinIfNotEmpty(" ", description, enumComment)

	// Convert to HTML
	html := string(blackfriday.Run([]byte(def.Description), blackfriday.WithNoExtensions()))
	def.HTMLDescription = strings.TrimSpace(pTags.ReplaceAllString(html, ""))
	return noDerive, nil
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
