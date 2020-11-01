package tag

import (
	"errors"
	"regexp"
)

var jsonTagRgx = regexp.MustCompile(`(?i)json:"([a-z0-9_-]+)?,?(omitempty|inline)?"`)
var ErrJsonTagNotPresent = errors.New("json tag not present")
var ErrJsonIgnored = errors.New("field is ignored")
var ErrInvalidJsonTag = errors.New("invalid json tag")

// Result from parsing a json struct property tag
type JsonTagResult struct {
	// Property name override
	Name string

	// Omitempty flag
	OmitEmpty bool

	// Inline flag
	Inline bool
}

// Checks if the result is empty
func (result *JsonTagResult) Empty() bool {
	return !result.OmitEmpty && !result.Inline && len(result.Name) == 0
}

// Validates a json tag parse result
func (result *JsonTagResult) Validate() error {
	if result.Empty() {
		return errors.New("result is empty")
	}

	if result.Inline && result.OmitEmpty {
		return errors.New("can't set omitempty and inline")
	}

	if result.Inline && len(result.Name) > 0 {
		return errors.New("can't set inline flag and have a name override")
	}

	return nil
}

// Parse json tag from a struct property
func ParseJsonTag(tags string) (*JsonTagResult, error) {
	m := jsonTagRgx.FindAllStringSubmatch(tags, -1)

	if len(m) == 0 {
		return nil, ErrJsonTagNotPresent
	}

	sm := m[0]
	if len(sm) > 2 && sm[1] == "" && sm[2] == "" {
		// empty json tag
		return nil, ErrInvalidJsonTag
	}

	if sm[1] == "-" {
		// ignored field
		return nil, ErrJsonIgnored
	}

	return &JsonTagResult{
		Name:      sm[1],
		OmitEmpty: len(sm) > 2 && sm[2] == "omitempty",
		Inline:    len(sm) > 2 && sm[2] == "inline",
	}, nil
}