package golang

import (
	"github.com/stretchr/testify/assert"
	"go/ast"
	"go/token"
	"testing"
)

func TestTypeFromToken(t *testing.T) {
	a := assert.New(t)
	tk := TypeFromToken(token.INT)
	a.Equal("int", tk.Name)

	tk = TypeFromToken(token.FLOAT)
	a.Equal("float", tk.Name)
}

func TestTypeFromName(t *testing.T) {
	a := assert.New(t)

	tk := TypeFromName("Image")
	a.Equal("Image", tk.Name)
}

func TestTypeFromExpr(t *testing.T) {
	a := assert.New(t)

	var expr ast.Expr

	expr = &ast.Ident{Name: "Image"}
	res := TypeFromExpr(expr)
	a.Equal("Image", res.Name)

	expr = &ast.StarExpr{
		X: &ast.Ident{Name: "Image"},
	}
	res = TypeFromExpr(expr)
	a.Equal("Image", res.Name)
	a.True(res.Pointer)

	expr = &ast.MapType{
		Key:   &ast.Ident{Name: "string"},
		Value: &ast.Ident{Name: "bool"},
	}
	res = TypeFromExpr(expr)
	a.True(res.Map)
	a.NotNil(res.MapKey)
	a.NotNil(res.MapValue)
	a.Equal("string", res.MapKey.Name)
	a.Equal("bool", res.MapValue.Name)

	expr = &ast.ArrayType{
		Elt: &ast.Ident{Name: "string"},
	}
	res = TypeFromExpr(expr)
	a.True(res.Array)
	a.Equal("string", res.Name)

	expr = &ast.SelectorExpr{
		X: &ast.Ident{
			Name: "ObjectID",
		},
		Sel: &ast.Ident{
			Name: "primitive",
		},
	}
	res = TypeFromExpr(expr)
	a.Equal("primitive", res.From)
	a.Equal("ObjectID", res.Name)

	expr = &ast.InterfaceType{}
	res = TypeFromExpr(expr)
	a.True(res.Generic)
}
