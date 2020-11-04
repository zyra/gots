package golang

import (
	"github.com/zyra/gots/pkg/parser/reader"
	"go/ast"
	"go/token"
	"strings"
)

// Package import options
type Import struct {
	// Package path
	//
	// This might not be populated immediately
	// since we need access to the file imports to find the path
	Package string `json:"package"`

	// Name of import
	//
	// Might be alias or real package name
	Name string `json:"alias"`
}

// Returns a type with the token's string representation as the name
func TypeFromToken(t token.Token) *reader.Type {
	return &reader.Type{
		Name: strings.ToLower(t.String()),
	}
}

// Returns a Type with the given name
func TypeFromName(n string) *reader.Type {
	return &reader.Type{Name: n}
}

// Parses ast.Expr and returns a Type
func TypeFromExpr(t ast.Expr) *reader.Type {
	switch v := t.(type) {
	case *ast.Ident:
		if v.Name == "byte" {
			return &reader.Type{Generic: true}
		}

		return TypeFromName(v.Name)

	case *ast.StarExpr:
		t := TypeFromExpr(v.X)
		t.Pointer = true
		return t

	case *ast.MapType:
		k, vv := TypeFromExpr(v.Key), TypeFromExpr(v.Value)
		return &reader.Type{
			Map:      true,
			MapKey:   k,
			MapValue: vv,
		}

	case *ast.ArrayType:
		if v, ok := v.Elt.(*ast.Ident); ok && v.Name == "byte" {
			return &reader.Type{Generic: true}
		}

		return &reader.Type{
			Name:  TypeFromExpr(v.Elt).Name,
			Array: true,
		}

	case *ast.SelectorExpr:
		return &reader.Type{
			Name: v.Sel.Name,
			From: v.X.(*ast.Ident).Name,
		}

	case *ast.InterfaceType:
		return &reader.Type{Generic: true}

	default:
		return &reader.Type{Generic: true}
	}
}
