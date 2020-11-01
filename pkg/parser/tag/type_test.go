package tag

import (
	"github.com/stretchr/testify/assert"
	"go/token"
	"testing"
)

func TestType_String(t *testing.T) {
	a := assert.New(t)
	tests := map[Type]string{
		TypeString:  "string",
		TypeBool:    "bool",
		TypeNumber:  "number",
		TypeUnknown: "unknown",
	}

	for k, v := range tests {
		a.Equal(v, k.String())
	}
}

func TestTypeFromToken(t *testing.T) {
	a := assert.New(t)
	tests := map[token.Token]Type{
		token.INT:       TypeNumber,
		token.FLOAT:     TypeNumber,
		token.STRING:    TypeString,
		token.INTERFACE: TypeUnknown,
	}
	for k, v := range tests {
		a.Equal(v, TypeFromToken(k))
	}
}
