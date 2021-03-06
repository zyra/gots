package reader

import (
	"io/ioutil"
	"path/filepath"
)

// Opts for TypeScript type alias
type TypeAlias struct {
	Name        string `json:"name"`        // Alias name
	AliasedType *Type   `json:"aliasedType"` // Aliased type
}

// options for TypeScript enum
type Enum struct {
	Name   string       `json:"name"`
	Values []*Constant `json:"values"`
}


type Constant struct {
	Name  string `json:"name"`  // Constant name
	Type  *Type   `json:"type"`  // Constant type
	Value string `json:"value"` // Constant value

	EnumValue bool   `json:"enumValue"` // This const defines an enum value
	EnumName  string `json:"enumName"`  // Enum name that this const defines a value for
}

type ReadConfig struct {
	Dir       string `json:"dir"`
	Recursive bool   `json:"recursive"`
}

func scanDirs(path string, recursive bool) ([]string, error) {
	contents, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	infos := make([]string, 0)
	for _, it := range contents {
		if !it.IsDir() {
			continue
		}

		path := filepath.Join(path, it.Name())
		infos = append(infos, path)

		if !recursive {
			continue
		}

		res, err := scanDirs(path, true)
		if err != nil {
			return nil, err
		}

		infos = append(infos, res...)
	}

	return infos, nil
}

func scanFiles(path string, recursive bool) ([]string, error) {
	contents, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	infos := make([]string, 0)
	for _, it := range contents {
		path := filepath.Join(path, it.Name())
		if !it.IsDir() {
			infos = append(infos, path)
			continue
		}

		if !recursive {
			continue
		}

		res, err := scanFiles(path, true)
		if err != nil {
			return nil, err
		}

		infos = append(infos, res...)
	}

	return infos, nil
}

// Returns a list of files matching the provided config
func (rc *ReadConfig) Files() ([]string, error) {
	return scanFiles(rc.Dir, rc.Recursive)
}

func (rc *ReadConfig) Directories() ([]string, error) {
	if !rc.Recursive {
		return []string{rc.Dir}, nil
	}

	return scanDirs(rc.Dir, true)
}

// Root scope that contains uniform description of interfaces, types, constants, enums... etc
type Spec interface {
	Packages() []*Package // Returns the packages that were parsed
}

type ReadFn func(config *ReadConfig) (Spec, error)

type WriteConfig struct {
	Output     Output       `json:"output"`
	Transforms []*Transform `json:"transforms,omitempty"`
	Packages   []*Package   `json:"packages,omitempty"`
}

type Writer interface {
	// Write package(s) with the provided output config and transforms
	// Return a map with output data for each package
	Write(spec Spec, config *WriteConfig) (map[string][]byte, error)
}
