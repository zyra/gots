package golang

import (
	"go/ast"
)

// Type alias options
type TypeAlias struct {
	// Type name
	Name string `json:"name"`

	// Aliased type data
	Type *Type `json:"type"`
}

// Parse Go type alias
func ParseTypeAlias(spec *ast.TypeSpec) *TypeAlias {
	return &TypeAlias{
		Name: spec.Name.Name,
		Type: TypeFromExpr(spec.Type),
	}
}
