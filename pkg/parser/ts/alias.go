package ts

// Opts for TypeScript type alias
type TypeAlias struct {
	// Alias name
	Name string `json:"name"`

	// Aliased type
	Type *Type `json:"type"`
}
