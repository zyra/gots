package ts

// Enum value opts
type EnumValue struct {
	// Enum value name
	Name string `json:"name"`

	// Enum value
	Value string `json:"value"`
}

// options for TypeScript enum
type Enum struct {
	Name   string       `json:"name"`
	Values []*EnumValue `json:"values"`
}
