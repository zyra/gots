package ts

type Constant struct {
	// Constant name
	Name string `json:"name"`

	// Constant type
	Type Type `json:"type"`

	// Constant value
	Value string `json:"value"`
}
