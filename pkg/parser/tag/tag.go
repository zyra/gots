package tag

import (
	"errors"
	"fmt"
)

type Tag struct {
	// Tag name
	Name string

	// Tag type
	Type string

	// Whether the tagged property is optional
	Optional bool
}

var ErrPropertyInlined = errors.New("property is inlined")
var ErrPropertyIgnored = errors.New("property is ignored")

func ParseTag(val string) (*Tag, error) {
	t := Tag{}
	j, err := ParseJsonTag(val)
	if err != nil && err != ErrJsonTagNotPresent {
		return nil, fmt.Errorf("failed to parse json tag: %v", err)
	} else if j != nil && j.Inline {
		return nil, ErrPropertyInlined
	} else if j != nil {
		t.Name = j.Name
		t.Optional = j.OmitEmpty
	}

	g, err := ParseGotsTag(val)
	if err != nil {
		return &t, nil
	}

	if g.Omit {
		return nil, ErrPropertyIgnored
	}

	if len(g.Type) > 0 {
		t.Type = g.Type
	}

	if g.Optional {
		t.Optional = true
	}

	return &t, nil
}
