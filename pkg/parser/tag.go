package parser

import (
	"errors"
	"go/ast"
	"go/token"
	"regexp"
	"strings"
)

type Tag struct {
	// Tag name
	Name string

	// Tag type
	Type string

	// Whether the tagged property is optional
	Optional bool
}

var jsonTagRgx = regexp.MustCompile(`(?i)json:"([a-z0-9_-]+),?(omitempty)?"`)
var gotsTagRgx = regexp.MustCompile(`(?i)gots:"([a-z0-9_,:\[\]]+)"`)
var gotsInnerTagRgx = regexp.MustCompile(`(?i)(name|type|optional):?([a-z0-9_\[\]]+)?`)

var errJsonTagNotPresent = errors.New("json Tag not present")
var errJsonIgnored = errors.New("field is ignored")

func ParseTag(val string) (*Tag, error) {
	m := jsonTagRgx.FindAllStringSubmatch(val, -1)

	if len(m) == 0 {
		return nil, errJsonTagNotPresent
	}

	sm := m[0]

	if sm[1] == "-" {
		return nil, errJsonIgnored
	}

	t := &Tag{
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

type TagType string

const (
	TagTypeString  TagType = "string"
	TagTypeBool    TagType = "bool"
	TagTypeNumber  TagType = "number"
	TagTypeUnknown TagType = "unknown"
)

func (t TagType) String() string {
	return string(t)
}

func TagTypeFromToken(t token.Token) TagType {
	switch t {
	case token.INT, token.FLOAT:
		return TagTypeNumber
	case token.STRING:
		return TagTypeString
	default:
		return TagTypeUnknown
	}
}

func ParseTypeFromToken(t token.Token) string {
	switch t {
	case token.INT, token.FLOAT:
		return "number"
	case token.STRING:
		return "string"
	default:
		panic("???")
	}
}

func TagTypeFromName(n string) TagType {
	switch n {
	case "string":
		return TagTypeString
	case "uint8", "uint16", "uint32", "uint64", "uint", "int8", "int16", "int32", "int64", "int", "float32", "float64":
		return TagTypeNumber
	case "bool":
		return TagTypeBool
	default:
		return TagTypeUnknown
	}
}

func ParseTagTypeFromName(n string) string {
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

func ParseTagType(t ast.Expr) string {
	switch v := t.(type) {
	case *ast.Ident:
		return TagTypeFromName(v.Name).String()
	case *ast.StarExpr:
		return ParseTagType(v.X)
	case *ast.MapType:
		key, value := ParseTagType(v.Key), ParseTagType(v.Value)
		return strings.Join([]string{"{[key:", key, "]:", value, "}"}, "")
	case *ast.ArrayType:
		return strings.Join([]string{ParseTagType(v.Elt), "[]"}, "")
	case *ast.SelectorExpr, *ast.InterfaceType:
		return "any"
	default:
		return ""
	}
}
