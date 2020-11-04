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

	props := make([]*reader.Property, 0, s.Fields.NumFields())
	for i := range s.Fields.List {
		f := s.Fields.List[i]

		var t *tag.Tag
		var err error
		if f.Tag != nil {
			t, err = tag.ParseTag(f.Tag.Value)
		}

		propType := TypeFromExpr(f.Type)
		if t != nil && len(t.Type) > 0 {
			propType.Name = t.Type
		}

		inline := false
		optional := false

		if err != nil {
			if err == tag.ErrJsonIgnored || err == tag.ErrJsonTagNotPresent || err == tag.ErrPropertyIgnored {
				continue
			}

			if err == tag.ErrPropertyInlined {
				inline = true
			} else {
				return nil, fmt.Errorf("failed to parse tag: %v", err)
			}
		}

		if t != nil {
			if len(t.Type) > 0 {
				propType.Name = t.Type
				propType.From = ""
			}

			if t.Optional {
				optional = t.Optional
			}
		}

		for _, n := range f.Names {
			if !ast.IsExported(n.Name) {
				continue
			}
			prop := reader.Property{
				Name:     n.Name,
				Type:     propType,
				Optional: optional,
				Inline:   inline,
			}

			if t != nil && len(t.Name) > 0 {
				prop.Name = t.Name
			}

			props = append(props, &prop)
		}
	}

	st.Properties = props

	return &st, nil
}
