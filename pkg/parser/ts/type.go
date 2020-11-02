package ts

// Typescript type options
type Type struct {
	// Type name
	Name string `json:"name"`

	// Package name that contains this type
	From string `json:"from"`

	// Type is an array
	Array bool `json:"array"`

	// Type is a map
	Map bool `json:"bool"`

	// Map key type
	MapKey *Type `json:"mapKey"`

	// Map value type
	MapValue *Type `json:"mapValue"`

	// Type is generic
	Generic bool `json:"generic"`
}
