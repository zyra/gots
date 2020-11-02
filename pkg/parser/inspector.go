package parser

import "go/ast"

type Inspector struct {
	inspect func(ast.Node) (bool, error)
	err     error
}

func (f *Inspector) Error() error {
	return f.err
}

func (f *Inspector) Visit(node ast.Node) ast.Visitor {
	if ok, err := f.inspect(node); err != nil {
		f.err = err
		return nil
	} else if ok {
		return f
	}
	return nil
}
