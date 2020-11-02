package golang

import (
	"fmt"
	"go/ast"
)

// Const options
type Const struct {
	// Constant name
	Name string `json:"name"`

	// Constant type data
	Type *Type `json:"type"`

	// Constant value
	Value string `json:"value"`
}

func ConstFromValueSpec(spec *ast.ValueSpec) (*Const, error) {
	cName := spec.Names[0].Name

	if !ast.IsExported(cName) {
		return nil, ErrNotExported
	}

	c := &Const{
		Name: cName,
	}

	if len(spec.Values) == 0 {
		return nil, fmt.Errorf("%s doesn't have a value", c.Name)
	}

	if spec.Type != nil {
		c.Type = TypeFromExpr(spec.Type)
		if val, ok := spec.Values[0].(*ast.BasicLit); ok {
			c.Value = val.Value
			return c, nil
		}
		return nil, fmt.Errorf("%s doesn't have a value", c.Name)
	}

	switch v := spec.Values[0].(type) {
	case *ast.CallExpr:
		if val, ok := v.Args[0].(*ast.BasicLit); ok {
			c.Type = TypeFromToken(val.Kind)
			c.Value = val.Value
			return c, nil
		}
		return nil, fmt.Errorf("unhandled const value type: %t", spec.Values[0])

	case *ast.BasicLit:
		c.Type = TypeFromToken(v.Kind)
		c.Value = v.Value
		return c, nil

	case *ast.Ident:
		c.Value = v.Name
		return c, nil

	default:
		return nil, fmt.Errorf("unhandled const value type: %t", spec.Values[0])
	}
}
