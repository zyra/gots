package parser

type Config struct {
	// Root directory to scan for go source files
	RootDir string

	// Name of file to output typescript code to
	OutFileName string

	// Whether to scan directories recursively
	Recursive bool
}
