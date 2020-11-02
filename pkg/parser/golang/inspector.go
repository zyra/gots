package golang

import "go/ast"

type InspectFn func(ast.Node) (bool, error)

type Inspector struct {
	inspect InspectFn
	err     error
}

func NewInspector(fn InspectFn) *Inspector {
	return &Inspector{inspect: fn}
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
