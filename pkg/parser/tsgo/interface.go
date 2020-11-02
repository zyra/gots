package tsgo

import (
	"github.com/zyra/gots/pkg/parser/golang"
	"github.com/zyra/gots/pkg/parser/ts"
)

func Interface(in *golang.Interface) *ts.Interface {
	it := ts.Interface{
		Name: in.Name,
	}
	return &it
}

func Property(in *golang.Property) *ts.Property {
	p := ts.Property{
		Name:     in.Name,
		Type:     NewType(in.Type),
		Optional: in.Optional,
	}
	return &p
}

func InterfaceFromStruct(in *golang.Struct) *ts.Interface {
	it := ts.Interface{Name: in.Name}

	if len(it.Properties) == 0 {
		return &it
	}

	it.Properties = make([]*ts.Property, 0, len(in.Properties))
	for _, p := range in.Properties {
		if p.Inline {

		}
		it.Properties = append(it.Properties, Property(p))
	}

	return &it
}
