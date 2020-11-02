package ts

// Typescript package options
// A package defines a set of interfaces, types, and constants that exist in the same file/module
type Package struct {
	// constants
	Constants []*Constant `json:"constants"`

	// interfaces
	Interfaces []*Interface `json:"interfaces"`

	// types
	Types []*Type `json:"types"`
}
