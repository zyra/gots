package tag

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJsonTagResult_Empty(t *testing.T) {
	a := assert.New(t)
	r := JsonTagResult{}
	a.True(r.Empty())

	r.Name = "test"
	a.False(r.Empty())

	r.OmitEmpty = true
	a.False(r.Empty())

	r.Inline = true
	a.False(r.Empty())

	r.OmitEmpty = false
	a.False(r.Empty())

	r.Name = ""
	a.False(r.Empty())
}

func TestJsonTagResult_Validate(t *testing.T) {
	a := assert.New(t)
	r := &JsonTagResult{
		Name:      "test",
		OmitEmpty: false,
		Inline:    false,
	}
	a.NoError(r.Validate())

	r.OmitEmpty = true
	a.NoError(r.Validate())

	r.Inline = true
	a.Error(r.Validate())

	r.OmitEmpty = false
	a.Error(r.Validate())

	r.Name = ""
	a.NoError(r.Validate())

	r.Inline = false
	a.Error(r.Validate())
}

func TestParseJsonTag(t *testing.T) {
	a := assert.New(t)
	tests := map[string]struct {
		n       string
		o, i, e bool
	}{
		`json:"name"`:           {n: "name"},
		`json:"name,omitempty"`: {n: "name", o: true},
		`json:",omitempty"`:     {o: true},
		`json:"-"`:              {e: true},
		`json:""`:               {e: true},
		`json:",inline"`:        {i: true},
	}

	for k, v := range tests {
		t, err := ParseJsonTag(k)
		if v.e {
			a.Errorf(err, "%s did not throw an error", k)
		} else if a.NoErrorf(err, "failed to parse tag: %s", k) {
			a.Equalf(v.n, t.Name, "name mismatch for %s", k)
			a.Equalf(v.o, t.OmitEmpty, "omitempty mismatch for %s", k)
			a.Equalf(v.i, t.Inline, "inline mismatch for %s", k)
		}
	}
}
