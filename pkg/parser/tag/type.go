package tag

import (
	"go/ast"
	"go/token"
	"strings"
)

type Type string

const (
	TypeString  Type = "string"
	TypeBool    Type = "bool"
	TypeNumber  Type = "number"
	TypeUnknown Type = "unknown"
)

func (t Type) String() string {
	return string(t)
}

func TypeFromToken(t token.Token) Type {
	switch t {
	case token.INT, token.FLOAT:
		return TypeNumber
	case token.STRING:
		return TypeString
	default:
		return TypeUnknown
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

func TypeFromName(n string) Type {
	switch n {
	case "string":
		return TypeString
	case "uint8", "uint16", "uint32", "uint64", "uint", "int8", "int16", "int32", "int64", "int", "float32", "float64":
		return TypeNumber
	case "bool":
		return TypeBool
	default:
		return TypeUnknown
	}
}

func ParseTypeFromName(n string) string {
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

func ParseType(t ast.Expr) string {
	switch v := t.(type) {
	case *ast.Ident:
		return TypeFromName(v.Name).String()
	case *ast.StarExpr:
		return ParseType(v.X)
	case *ast.MapType:
		key, value := ParseType(v.Key), ParseType(v.Value)
		return strings.Join([]string{"{[key:", key, "]:", value, "}"}, "")
	case *ast.ArrayType:
		return strings.Join([]string{ParseType(v.Elt), "[]"}, "")
	case *ast.SelectorExpr, *ast.InterfaceType:
		return "any"
	default:
		return ""
	}
}
