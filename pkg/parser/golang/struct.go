package golang

import (
	"errors"
	"fmt"
	"github.com/zyra/gots/pkg/parser/reader"
	"github.com/zyra/gots/pkg/parser/tag"
	"go/ast"
)

func ParseStruct(spec *ast.TypeSpec) (*reader.Interface, error) {
	s, ok := spec.Type.(*ast.StructType)
	if !ok {
		return nil, errors.New("invalid TypeSpec")
	}

	itName := spec.Name.Name
	if !ast.IsExported(itName) {
		return nil, ErrNotExported
	}

	st := reader.Interface{}
	st.Name = itName

	if s.Fields == nil || s.Fields.NumFields() == 0 {
		return &st, nil
	}

	props := make([]*reader.Property, 0)
	for i := range s.Fields.List {
		f := s.Fields.List[i]
		var t *tag.Tag
		var err error
		if f.Tag != nil {
			t, err = tag.ParseTag(f.Tag.Value)
		}

		propType := TypeFromExpr(f.Type)

		prop := reader.Property{}

		if err != nil {
			if err == tag.ErrJsonIgnored || err == tag.ErrJsonTagNotPresent || err == tag.ErrPropertyIgnored {
				continue
			}

			if err == tag.ErrPropertyInlined {
				prop.Inline = true
			} else {
				return nil, fmt.Errorf("failed to parse tag: %v", err)
			}
		} else if t != nil {
			if t.Type != "" {
				propType.Name = t.Type
				propType.From = ""
			}

			prop.Name = t.Name
			prop.Optional = t.Optional
		} else {
			prop.Name = spec.Name.Name
		}

		prop.Type = propType
		props = append(props, &prop)
	}

	st.Properties = props

	return &st, nil
}
