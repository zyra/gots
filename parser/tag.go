package parser

import (
	"errors"
	"go/ast"
	"go/token"
	"regexp"
	"strings"
)

type tag struct {
	Name, Type string
	Optional   bool
}

var jsonTagRgx = regexp.MustCompile(`(?i)json:"([a-z0-9_-]+),?(omitempty)?"`)
var gotsTagRgx = regexp.MustCompile(`(?i)gots:"([a-z0-9_,:\[\]]+)"`)
var gotsInnerTagRgx = regexp.MustCompile(`(?i)(name|type|optional):?([a-z0-9_\[\]]+)?`)

var errJsonTagNotPresent = errors.New("json tag not present")
var errJsonIgnored = errors.New("field is ignored")

func parseTags(val string) (*tag, error) {
	m := jsonTagRgx.FindAllStringSubmatch(val, -1)

	if len(m) == 0 {
		return nil, errJsonTagNotPresent
	}

	sm := m[0]

	if sm[1] == "-" {
		return nil, errJsonIgnored
	}

	t := &tag{
		Name:     sm[1],
		Optional: sm[2] == "omitempty",
	}

	m = gotsTagRgx.FindAllStringSubmatch(val, -1)

	if len(m) == 1 && len(m[0]) > 0 {
		sm = m[0]
		m = gotsInnerTagRgx.FindAllStringSubmatch(sm[1], -1)

		for i := range m {
			switch m[i][1] {
			case "type":
				t.Type = m[i][2]
			case "name":
				t.Name = m[i][2]
			case "optional":
				t.Optional = true
			}
		}
	}

	return t, nil
}

func parseTypeFromKind(t token.Token) string {
	switch t {
	case token.INT, token.FLOAT:
		return "number"
	case token.STRING:
		return "string"
	default:
		panic("???")
	}
}

func parseTypeFromName(n string) string {
	switch n {
	case "string":
		return "string"
	case "uint8", "uint16", "uint32", "uint64", "uint", "int8", "int16", "int32", "int64", "int", "float32", "float64":
		return "number"
	case "bool":
		return "boolean"
	default:
		return n
	}
}

func parseType(t ast.Expr) string {
	switch t.(type) {
	case *ast.Ident:
		return parseTypeFromName(t.(*ast.Ident).Name)
	case *ast.StarExpr:
		return parseType(t.(*ast.StarExpr).X)
	case *ast.MapType:
		tt := t.(*ast.MapType)
		k, v := parseType(tt.Key), parseType(tt.Value)
		return strings.Join([]string{"{[key:", k, "]:", v, "}"}, "")
	case *ast.ArrayType:
		return strings.Join([]string{parseType(t.(*ast.ArrayType).Elt), "[]"}, "")
	case *ast.SelectorExpr, *ast.InterfaceType:
		return "any"
	default:
		return ""
	}
}
