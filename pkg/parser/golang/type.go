package golang

import (
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

// Type options (for variables and properties)
type Type struct {
	// Base type name (e.g string, ObjectID...etc)
	Name string `json:"name"`

	// opts if the type is external
	From *Import `json:"importOpts"`

	// Whether this type is an array
	Array bool `json:"array"`

	// Whether this type is a map
	Map bool `json:"map"`

	// Type used for map keys
	MapKey *Type `json:"mapKey"`

	// Type used for map values
	MapValue *Type `json:"mapValue"`

	// Whether this type is a pointer
	Pointer bool `json:"pointer"`

	// Generic type (i.e interface{})
	Generic bool `json:"generic"`
}

// Return a string representation of the type
func (t *Type) String() string {
	if t.Generic {
		return "interface{}"
	}

	if t.Array {
		return strings.Join([]string{t.Name, "[]"}, "")
	}
	return t.Name
}

// Returns a type with the token's string representation as the name
func TypeFromToken(t token.Token) *Type {
	return &Type{
		Name: t.String(),
	}
}

// Returns a Type with the given name
func TypeFromName(n string) *Type {
	return &Type{Name: n}
}

// Parses ast.Expr and returns a Type
func TypeFromExpr(t ast.Expr) *Type {
	switch v := t.(type) {
	case *ast.Ident:
		return TypeFromName(v.Name)

	case *ast.StarExpr:
		t := TypeFromExpr(v.X)
		t.Pointer = true
		return t

	case *ast.MapType:
		return &Type{
			Map:      true,
			MapKey:   TypeFromExpr(v.Key),
			MapValue: TypeFromExpr(v.Value),
		}

	case *ast.ArrayType:
		return &Type{
			Name:  TypeFromExpr(v.Elt).Name,
			Array: true,
		}

	case *ast.SelectorExpr:
		return &Type{
			Name: v.Sel.Name,
			From: &Import{
				Name: v.X.(*ast.Ident).Name,
			},
		}

	case *ast.InterfaceType:
		return &Type{
			Generic: true,
		}

	default:
		return &Type{
			Generic: true,
		}
	}
}
