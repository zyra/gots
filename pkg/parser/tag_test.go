package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseTag(t *testing.T) {
	tests := map[string]struct {
		n, t string
		e, o bool
	}{
		`json:"name"`: {
			n: "name",
		},
		`json:"name,omitempty"`: {
			n: "name",
			o: true,
		},
		`json:"name" gots:"type:string"`: {
			n: "name",
			t: "string",
		},
		`json:"name" gots:"type:string,optional"`: {
			n: "name",
			t: "string",
			o: true,
		},
		`json:"name" gots:"optional"`: {
			n: "name",
			o: true,
		},
		`json:"-"`: {
			e: true,
		},
		`json:"name" gots:"-"`: {
			e: true,
		},
		`gots:"optional"`: {
			o: true,
		},
	}

	a := assert.New(t)

	for k, v := range tests {
		t, err := ParseTag(k)
		if v.e {
			a.Errorf(err, "%s did not throw an error", k)
		} else if a.NoErrorf(err, "failed to parse tag: %s", k) {
			a.Equalf(v.n, t.Name, "name mismatch for %s", k)
			a.Equalf(v.t, t.Type, "type mismatch for %s", k)
			a.Equalf(v.o, t.Optional, "optional mismatch for %s", k)
		}
	}
}
