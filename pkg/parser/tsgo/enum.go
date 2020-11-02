package tsgo

import (
	"github.com/zyra/gots/pkg/parser/golang"
	"github.com/zyra/gots/pkg/parser/ts"
)

func EnumValue(in *golang.Const) *ts.EnumValue {
	return &ts.EnumValue{
		Name:  in.Name,
		Value: in.Value,
	}
}

func Enum(in *golang.TypeAlias, values []*golang.Const) *ts.Enum {
	evals := make([]*ts.EnumValue, len(values))
	for i := range values {
		evals[i] = EnumValue(values[i])
	}
	return &ts.Enum{
		Name:   in.Name,
		Values: evals,
	}
}
