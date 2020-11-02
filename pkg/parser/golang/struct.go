package golang

import (
	"errors"
	"fmt"
	"github.com/zyra/gots/pkg/parser/reader"
	"github.com/zyra/gots/pkg/parser/tag"
	"go/ast"
)

// Property options
type Property = reader.Property

// Struct options
type Struct struct {
	reader.Interface
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

	st := Struct{}
	st.Name = itName

	if s.Fields == nil || s.Fields.NumFields() == 0 {
		return &st, nil
	}

	props := make([]*Property, 0)
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
				propType.From = ""
			}

			prop.Name = t.Name
			prop.Optional = t.Optional
		}

		prop.Type = propType
		props = append(props, &prop)
	}

	st.Properties = props

	return &st, nil
}
