package definition

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

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

// parseAsValue takes a string and parses it for use as a default, example or
// enum variant
func parserAsValue(v string) (interface{}, error) {
	expr, err := parser.ParseExpr(v)
	if err != nil {
		// return as string
		return v, nil
	}
	switch lit := expr.(type) {
	case *ast.BasicLit:
		switch lit.Kind {
		case token.STRING, token.CHAR:
			str, err := strconv.Unquote(lit.Value)
			if err != nil {
				panic("Couldn't unquote basic literal of type STRING or CHAR")
			}
			return str, nil
		case token.INT:
			return strconv.Atoi(v)
		case token.FLOAT:
			return strconv.ParseFloat(v, 64)
		default:
			return nil, errors.Errorf("unsupported literal kind %s", lit.Kind)
		}
	case *ast.Ident:
		switch lit.Name {
		// Go, where you can redefine `true := false`
		case "true", "false":
			return strconv.ParseBool(lit.Name)
		default:
			return nil, errors.Errorf("can't understand literal %v", lit)
		}
	default:
		return nil, errors.Errorf("can't handle %s (type %T)", v, expr)
	}
}
