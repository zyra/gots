package ts

// Typescript property options
type Property struct {
	// Property name
	Name string `json:"name"`

	// Property type
	Type Type `json:"type"`

	// Property is optional
	Optional bool `json:"optional"`
}

// Options for interface that another interface is extending
type ExtendsInterface struct {
	// Interface name
	Name string `json:"name"`

	// Name of package containing this interface
	From string `json:"from"`
}

// Typescript interface options
type Interface struct {
	// Interface name
	Name string `json:"name"`

	// Interface properties
	Properties []*Property `json:"properties,omitempty"`

	// Extends another interface
	Extends *ExtendsInterface `json:"extends,omitempty"`
}
