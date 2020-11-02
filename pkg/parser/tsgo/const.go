package tsgo

import (
	"github.com/zyra/gots/pkg/parser/godef"
	"github.com/zyra/gots/pkg/parser/ts"
)

func NewType(in *godef.Type) ts.Type {
	t := ts.Type{
		Name:    in.Name,
		Array:   in.Array,
		Map:     in.Map,
		Generic: in.Generic,
	}

	if in.From != nil {
		t.From = in.From.Name
	}

	if in.Map {
		k, v := NewType(in.MapKey), NewType(in.MapValue)
		t.MapKey, t.MapValue = &k, &v
	}

	return t
}

func NewConst(in *godef.Const) *ts.Constant {
	return &ts.Constant{
		Name:  in.Name,
		Type:  NewType(in.Type),
		Value: in.Value,
	}
}
