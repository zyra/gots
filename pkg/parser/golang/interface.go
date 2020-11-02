package golang

import (
	"errors"
	"go/ast"
)

// Interface options
type Interface struct {
	// Interface name
	Name string `json:"name"`
}

// Parse Go interface from ast.TypeSpec
func ParseInterface(spec *ast.TypeSpec) (*Interface, error) {
	_, ok := spec.Type.(*ast.InterfaceType)
	if !ok {
		return nil, errors.New("invalid TypeSpec")
	}
	name := spec.Name.Name
	return &Interface{Name: name}, nil
}
