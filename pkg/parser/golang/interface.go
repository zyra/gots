package golang

import (
	"errors"
	"github.com/zyra/gots/pkg/parser/reader"
	"go/ast"
)

// Parse Go interface from ast.TypeSpec
func ParseInterface(spec *ast.TypeSpec) (*reader.Interface, error) {
	_, ok := spec.Type.(*ast.InterfaceType)
	if !ok {
		return nil, errors.New("invalid TypeSpec")
	}
	name := spec.Name.Name
	return &reader.Interface{Name: name}, nil
}
