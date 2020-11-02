package golang

import (
	"errors"
	"fmt"
	"github.com/zyra/gots/pkg/parser/reader"
	"go/ast"
	"go/constant"
)

// Const options
type Const struct {
	reader.Constant
}

func evalBinaryExpr(expr *ast.BinaryExpr) (constant.Value, error) {
	xLit, ok := expr.X.(*ast.BasicLit)
	if !ok {
		return constant.MakeUnknown(), errors.New("left operand is not BasicLit")
	}

	yLit, ok := expr.Y.(*ast.BasicLit)
	if !ok {
		return constant.MakeUnknown(), errors.New("right operand is not BasicLit")
	}

	x := evalBasicLit(xLit)
	y := evalBasicLit(yLit)
	return constant.BinaryOp(x, expr.Op, y), nil
}

func evalBasicLit(expr *ast.BasicLit) constant.Value {
	return constant.MakeFromLiteral(expr.Value, expr.Kind, 0)
}

func ConstFromValueSpec(spec *ast.ValueSpec) (*Const, error) {
	cName := spec.Names[0].Name

	if !ast.IsExported(cName) {
		return nil, ErrNotExported
	}

	c := &Const{
		Constant: reader.Constant{
			Name: cName,
		},
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

		case *ast.BinaryExpr:
			vv, err := evalBinaryExpr(val)
			if err != nil {
				return nil, err
			}
			c.Value = vv.ExactString()

		case *ast.Ident:
			c.Value = constant.Make(val.Name).String()
			fmt.Println("done")
		}

		return nil, fmt.Errorf("%s doesn't have a value", c.Name)
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
