package parser

import (
	"go/ast"
	"go/token"
	"strings"
)

func TypeFromToken(t token.Token) *Type {
	switch t {
	case token.INT, token.FLOAT:
		return &Type{
			Name: "number",
		}
	case token.STRING:
		return &Type{Name: "string"}
	default:
		return nil
	}
}

func ParseTypeFromToken(t token.Token) *Type {
	switch t {
	case token.INT, token.FLOAT:
		return &Type{Name: "number"}
	case token.STRING:
		return &Type{Name: "string"}
	default:
		return nil
	}
}

func TypeFromName(n string) *Type {
	switch n {
	case "string":
		return &Type{Name: "string"}
	case "uint8", "uint16", "uint32", "uint64", "uint", "int8", "int16", "int32", "int64", "int", "float32", "float64":
		return &Type{Name: "number"}
	case "bool":
		return &Type{Name: "boolean"}
	default:
		return &Type{Name: n}
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

func ParseType(t ast.Expr) *Type {
	switch v := t.(type) {
	case *ast.Ident:
		t := TypeFromName(v.Name)
		return t
	case *ast.StarExpr:
		return ParseType(v.X)
	case *ast.MapType:
		key, value := ParseType(v.Key), ParseType(v.Value)
		return &Type{
			Name: strings.Join([]string{"{[key:", key.Name, "]:", value.Name, "}"}, ""),
		}
	case *ast.ArrayType:
		return &Type{
			Name:  ParseType(v.Elt).Name,
			Array: true,
		}
	case *ast.SelectorExpr:
		return &Type{
			Name: v.Sel.Name,
			Import: &Import{
				Alias: v.X.(*ast.Ident).Name,
			},
		}
	case *ast.InterfaceType:
		return &Type{
			Name: "any",
		}
	default:
		return &Type{
			Name: "any",
		}
	}
}
