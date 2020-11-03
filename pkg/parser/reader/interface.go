package reader

// property options
type Property struct {
	Name     string `json:"name"`     // Property name
	Type     Type   `json:"type"`     // Property type
	Optional bool   `json:"optional"` // Property is optional

	Inline bool `json:"inline"`
}

// interface extension opts
type ExtendsInterface struct {
	Name string `json:"name"` // Interface name
	From string `json:"from"` // Name of package containing this interface
}

// interface options
type Interface struct {
	Name       string            `json:"name"`                 // Interface name
	Properties []*Property       `json:"properties,omitempty"` // Interface properties
	Extends    *ExtendsInterface `json:"extends,omitempty"`    // Extends another interface
}

