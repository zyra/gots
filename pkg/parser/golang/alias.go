package golang

import (
	"github.com/zyra/gots/pkg/parser/reader"
	"go/ast"
)

// Parse Go type alias
func ParseTypeAlias(spec *ast.TypeSpec) *reader.TypeAlias {
	t := TypeFromExpr(spec.Type)
	return &reader.TypeAlias{
		Name:        spec.Name.Name,
		AliasedType: t,
	}
}
