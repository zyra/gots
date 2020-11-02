package tsgo

import (
	"github.com/zyra/gots/pkg/parser/godef"
	"github.com/zyra/gots/pkg/parser/ts"
)

func EnumValue(in *godef.Const) *ts.EnumValue {
	return &ts.EnumValue{
		Name:  in.Name,
		Value: in.Value,
	}
}

func Enum(in *godef.TypeAlias, values []*godef.Const) *ts.Enum {
	evals := make([]*ts.EnumValue, len(values))
	for i := range values {
		evals[i] = EnumValue(values[i])
	}
	return &ts.Enum{
		Name:   in.Name,
		Values: evals,
	}
}
