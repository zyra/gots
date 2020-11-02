package tsgo

import (
	"github.com/zyra/gots/pkg/parser/godef"
	"github.com/zyra/gots/pkg/parser/ts"
)

func TypeAlias(in *godef.TypeAlias) *ts.TypeAlias {
	t := NewType(in.Type)
	return &ts.TypeAlias{
		Name: in.Name,
		Type: &t,
	}
}
