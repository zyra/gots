package statement

type TypeRef struct {
	// Literal type name
	Name string

	// Type source
	// If blank, we assume that the type exists in the current directory
	Source string
}
