package tag

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseGotsTag(t *testing.T) {
	a := assert.New(t)
	tests := map[string]struct {
		t      string
		op, om bool
		e      error
	}{
		``:                            {e: ErrGotsTagNotFound},
		`gots:""`:                     {e: ErrBlankGotsTag},
		`gots:"-"`:                    {om: true},
		`gots:",optional"`:            {op: true},
		`gots:"type:string"`:          {t: "string"},
		`gots:"type:string,optional"`: {t: "string", op: true},
		`gots:"optional,type:string"`: {t: "string", op: true},
		`gots:"optional,type:"`:       {op: true},
		`gots:"type:"`:                {},
	}

	for k, v := range tests {
		t, err := ParseGotsTag(k)
		if v.e != nil {
			if a.Errorf(err, "%s did not throw an error", k) {
				a.Equal(v.e, err)
			}
		} else if a.NoErrorf(err, "failed to parse tag: %s", k) {
			a.Equalf(v.t, t.Type, "type mismatch for: %s", k)
			a.Equalf(v.om, t.Omit, "omit mismatch for: %s", k)
			a.Equalf(v.op, t.Optional, "optional mismatch for: %s", k)
		}
	}
}
