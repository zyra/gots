package reader

import "fmt"

// Typescript type options
type Type struct {
	Name string `json:"name"` // Type name
	From string `json:"from"` // Package name that contains this type

	EnumValue bool   `json:"enumValue"`

	Generic bool `json:"generic"` // Type is generic

	Array bool `json:"array"` // Type is an array

	Map      bool  `json:"bool"`     // Type is a map
	MapKey   *Type `json:"mapKey"`   // Map key type
	MapValue *Type `json:"mapValue"` // Map value type

	// Whether this type is a pointer
	Pointer bool `json:"pointer"`
}

func (t *Type) IsEmpty() bool {
	return !t.Generic && !t.Map && !t.Array && len(t.Name) == 0
}

func (t *Type) TSType() string {
	if t.Map {
		return fmt.Sprintf("{[key: string]: %s}", t.MapValue.TSType())
	}

	if t.Array {
		clone := *t
		clone.Array = false
		return fmt.Sprintf("%s[]", clone.TSType())
	}

	if t.Generic || t.From != "" {
		return "any"
	}

	return t.Name
}
