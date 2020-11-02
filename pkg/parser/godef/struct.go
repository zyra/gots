package godef

import (
	"errors"
	"fmt"
	"github.com/zyra/gots/pkg/parser/tag"
	"go/ast"
)

// Property options
type Property struct {
	// Property name
	Name string `json:"name"`

	// Property type
	Type *Type `json:"typeOpts"`

	// Whether the property is optional
	Optional bool `json:"optional"`

	// Property is inlined
	Inline bool `json:"inline"`

	// Inlined struct
	// Populated at a later stage
	InlinedStruct *Struct `json:"inlinedStruct"`
}

// Struct options
type Struct struct {
	// Type name
	Name string `json:"name"`

	// Struct properties
	Properties []*Property `json:"properties"`
}

func ParseStruct(spec *ast.TypeSpec) (*Struct, error) {
	s, ok := spec.Type.(*ast.StructType)
	if !ok {
		return nil, errors.New("invalid TypeSpec")
	}

	itName := spec.Name.Name
	if !ast.IsExported(itName) {
		return nil, ErrNotExported
	}

	var props []*Property
	if s.Fields != nil {
		nf := s.Fields.NumFields()
		if nf != 0 {
			props = make([]*Property, 0, nf)
			for i := range s.Fields.List {
				f := s.Fields.List[i]
				t, err := tag.ParseTag(f.Tag.Value)

				propType := TypeFromExpr(f.Type)

				prop := Property{}

				if err != nil {
					if err == tag.ErrJsonIgnored || err == tag.ErrJsonTagNotPresent || err == tag.ErrPropertyIgnored {
						continue
					}

					if err == tag.ErrPropertyInlined {
						prop.Inline = true
					} else {
						return nil, fmt.Errorf("failed to parse tag: %v", err)
					}
				} else {
					if t.Type != "" {
						propType.Name = t.Type
						propType.From = nil
					}

					prop.Name = t.Name
					prop.Optional = t.Optional
				}

				prop.Type = propType
				props = append(props, &prop)
			}
		}
	}

	st := Struct{
		Name:       itName,
		Properties: props,
	}

	return &st, nil
}
