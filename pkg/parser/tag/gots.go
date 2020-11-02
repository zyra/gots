package tag

import (
	"errors"
	"regexp"
)

type ParseGotsTagResult struct {
	// Type override
	Type string

	// Optional property
	Optional bool

	// Omit property even if it's exported as JSON
	Omit bool
}

var gotsTagRgx = regexp.MustCompile(`(?i)gots:"([a-z0-9_\-,:\[\]]+)?"`)
var gotsInnerTagRgx = regexp.MustCompile(`(?i)(type|optional|-)?:?([a-z0-9_\[\]]+)?`)

var ErrGotsTagNotFound = errors.New("couldn't find gots tag")
var ErrBlankGotsTag = errors.New("gots tag is empty")
var ErrInvalidGotsTag = errors.New("invalid gots tag")

func ParseGotsTag(tag string) (*ParseGotsTagResult, error) {
	m := gotsTagRgx.FindAllStringSubmatch(tag, -1)
	if len(m) == 0 {
		return nil, ErrGotsTagNotFound
	}
	if len(m[0]) < 2 || len(m[0][1]) == 0 {
		return nil, ErrBlankGotsTag
	}

	t := ParseGotsTagResult{}
	m = gotsInnerTagRgx.FindAllStringSubmatch(m[0][1], -1)

	for i := range m {
		switch m[i][1] {
		case "type":
			t.Type = m[i][2]
		case "optional":
			t.Optional = true
		case "-":
			t.Omit = true
			return &t, nil
		}
	}

	return &t, nil
}
