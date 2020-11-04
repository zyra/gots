package golang

import (
	"fmt"
	"github.com/zyra/gots/pkg/parser/reader"
	"go/ast"
	"go/constant"
)

func ConstFromValueSpec(spec *ast.ValueSpec) (*reader.Constant, error) {
	cName := spec.Names[0].Name

	c := &reader.Constant{
		Name: cName,
	}

	if len(spec.Values) == 0 {
		return c, nil
	}

	if spec.Type != nil {
		c.Type = TypeFromExpr(spec.Type)
		switch val := spec.Values[0].(type) {
		case *ast.BasicLit:
			c.Value = val.Value
			return c, nil
		case *ast.Ident:
			if val.Name == "iota" {
				c.Value = "0"
			} else {
				c.Value = constant.Make(val.Name).String()
			}
			return c, nil
		case *ast.BinaryExpr:
			getVal := func(v interface{}) constant.Value {
				switch xV := v.(type) {
				case *ast.BasicLit:
					return constant.MakeFromLiteral(xV.Value, xV.Kind, 0)

				case *ast.Ident:
					if xV.Name == "iota" {
						// Not always accurate, but should work for most enum cases
						return constant.MakeInt64(0)
					}
				}
				return constant.MakeUnknown()
			}

			xVal := getVal(val.X)
			yVal := getVal(val.Y)

			if xVal.Kind() == constant.Unknown || yVal.Kind() == constant.Unknown {
				return nil, fmt.Errorf("%s uses an unsupported binary op", c.Name)
			}

			bOpVal := constant.BinaryOp(xVal, val.Op, yVal)

			if bOpVal.Kind() == constant.Unknown {
				return nil, fmt.Errorf("failed to calculate binary op for %s", c.Name)
			}

			c.Value = bOpVal.String()
			return c, nil
		}

		return nil, fmt.Errorf("%s doesn't have a value or is not supported", c.Name)
	}

	if len(spec.Values) != 1 {
		panic("spec had more than 1 value")
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
