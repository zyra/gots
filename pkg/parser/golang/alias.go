package golang

import (
	"github.com/zyra/gots/pkg/parser/reader"
	"go/ast"
)

// Type alias options
type TypeAlias struct {
	reader.TypeAlias
}

// Parse Go type alias
func ParseTypeAlias(spec *ast.TypeSpec) *TypeAlias {
	t := TypeFromExpr(spec.Type)
	return &TypeAlias{
		TypeAlias: reader.TypeAlias{
			Name:        spec.Name.Name,
			AliasedType: t,
		},
	}
}
