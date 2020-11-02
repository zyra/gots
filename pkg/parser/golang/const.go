package golang

import (
	"fmt"
	"github.com/zyra/gots/pkg/parser/reader"
	"go/ast"
	"go/constant"
)

func ConstFromValueSpec(spec *ast.ValueSpec) (*reader.Constant, error) {
	cName := spec.Names[0].Name

	if !ast.IsExported(cName) {
		return nil, ErrNotExported
	}

	c := &reader.Constant{
		Name: cName,
	}

	if len(spec.Values) == 0 {
		return nil, fmt.Errorf("%s doesn't have a value", c.Name)
	}

	if spec.Type != nil {
		c.Type = TypeFromExpr(spec.Type)
		switch val := spec.Values[0].(type) {
		case *ast.BasicLit:
			c.Value = val.Value
			return c, nil
		case *ast.Ident:
			c.Value = constant.Make(val.Name).String()
		}

		return nil, fmt.Errorf("%s doesn't have a value or is not supported", c.Name)
	}

	switch v := spec.Values[0].(type) {
	case *ast.CallExpr:
		if val, ok := v.Args[0].(*ast.BasicLit); ok {
			t := TypeFromToken(val.Kind)
			c.Type = *t
			c.Value = val.Value
			return c, nil
		}
		return nil, fmt.Errorf("unhandled const value type: %t", spec.Values[0])

	case *ast.BasicLit:
		t := TypeFromToken(v.Kind)
		c.Type = *t
		c.Value = v.Value
		return c, nil

	case *ast.Ident:
		c.Value = v.Name
		return c, nil

	default:
		return nil, fmt.Errorf("unhandled const value type: %t", spec.Values[0])
	}
}
