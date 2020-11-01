package parser

// Configures what types to export
type TypesConfig struct {
	// Go interfaces
	Interfaces bool `json:"interfaces"`

	// Variable declarations that have default values
	Variables bool `json:"variables"`

	// Const declarations
	Constants bool `json:"constants"`

	// Type aliases
	Aliases bool `json:"aliases"`

	// Go structs
	Structs bool `json:"structs"`

	// Generate enum from type alias and constant combinations
	Enums bool `json:"enums"`
}

// Parser options
type Config struct {
	// Root directory to scan for go source files
	RootDir string `json:"rootDir"`

	Types TypesConfig `json:"types,omitempty"`

	// Output options
	Output Output `json:"output,omitempty"`

	// Whether to scan directories recursively
	Recursive bool `json:"recursive,omitempty"`

	// Custom global transformation
	// Allows applying gots options globally to specific types instead of using gots tags in code
	Transforms []*Transform `json:"transforms"`

	// Paths/globs to include
	// If unset, all files found in root directory will be parsed (and subdirectories if recursive is set)
	Include []string `json:"include,omitempty"`

	// Paths/globs to exclude
	Exclude []string `json:"exclude,omitempty"`
}

// Enum for output modes
type OutputMode string

const (
	// output all parsed types in a single file
	AIO OutputMode = "aio"

	// output each package into a separate file
	Packages OutputMode = "packages"

	// output directories and files that match the go source code
	Mirror OutputMode = "mirror"
)

// Output config
type Output struct {
	// Output mode
	// Defaults to packages
	Mode OutputMode `json:"mode,omitempty"`

	// Root directory to output file(s) to
	// Default to current directory
	Path string `json:"path,omitempty"`

	// Filename to use when using AIO mode
	// Defaults to `types.ts`
	AIOFileName string `json:"aioFileName,omitempty"`
}

// Transform config
type Transform struct {
	// Regex pattern to replace
	Pattern string `json:"pattern,omitempty"`

	// Exact type to replace
	// Ignored if pattern is set
	Type string `json:"type,omitempty"`

	// Type to use instead of pattern/from
	Target string `json:"target,omitempty"`

	// Force matches to be optional
	// If not set, it will fall back to default behaviour (parse tags)
	Optional bool `json:"optional,omitempty"`
}
