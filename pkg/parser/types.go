package parser

import "strings"

// Package import options
type Import struct {
	Package string `json:"package"`
	Alias   string `json:"alias"`
}

// Type options (for variables and properties)
type Type struct {
	// Base type name (e.g string, ObjectID...etc)
	Name string `json:"name"`

	// Import opts if the type is external
	Import *Import `json:"importOpts"`

	Array bool `json:"array"`
}

func (t *Type) String() string {
	if t.Array {
		return strings.Join([]string{t.Name, "[]"}, "")
	}
	return t.Name
}

// Property options
type Property struct {
	// Property name
	Name string `json:"name"`

	// Property type
	Type *Type `json:"typeOpts"`

	// Whether the property is optional
	Optional bool `json:"optional"`

	Inline bool `json:"inline"`
}

// Struct options
type Struct struct {
	// Type name
	Name string `json:"name"`

	// Struct properties
	Properties []*Property `json:"properties"`
}

// Interface options
type Interface struct {
	// Interface name
	Name string `json:"name"`
}

// Type alias options
type TypeAlias struct {
	// Type name
	Name string `json:"name"`

	// Aliased type data
	Type *Type `json:"type"`
}

// Const options
type Const struct {
	// Constant name
	Name string `json:"name"`

	// Constant type data
	Type *Type `json:"type"`

	// Constant value
	Value string `json:"value"`
}
