package reader

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// Typescript type options
type Type struct {
	Name string `json:"name"` // Type name
	From string `json:"from"` // Package name that contains this type

	Generic bool `json:"generic"` // Type is generic

	Array bool `json:"array"` // Type is an array

	Map      bool  `json:"bool"`     // Type is a map
	MapKey   *Type `json:"mapKey"`   // Map key type
	MapValue *Type `json:"mapValue"` // Map value type

	// Whether this type is a pointer
	Pointer bool `json:"pointer"`
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

// Opts for TypeScript type alias
type TypeAlias struct {
	Name        string `json:"name"`        // Alias name
	AliasedType Type   `json:"aliasedType"` // Aliased type
}

// Enum value opts
type EnumValue struct {
	Key   string `json:"name"`  // Enum value name
	Value string `json:"value"` // Enum value
}

// options for TypeScript enum
type Enum struct {
	Name   string       `json:"name"`
	Values []*EnumValue `json:"values"`
}

type Constant struct {
	Name  string `json:"name"`  // Constant name
	Type  Type   `json:"type"`  // Constant type
	Value string `json:"value"` // Constant value
}

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

// contains interfaces, types, and constants that were found in the same file
type File struct {
	Constants  []*Constant  `json:"constants"`  // constants
	Interfaces []*Interface `json:"interfaces"` // interfaces
	Types      []*TypeAlias `json:"types"`      // types
	SourcePath string       `json:"-"`          // original source path
}

func (f *File) EachInterface(fx func(iface *Interface) (bool, error)) error {
	for i := range f.Interfaces {
		if ok, err := fx(f.Interfaces[i]); err != nil {
			return err
		} else if !ok {
			return nil
		}
	}
	return nil
}

func (f *File) EachTypeAlias(fx func(t *TypeAlias) (bool, error)) error {
	for i := range f.Types {
		if ok, err := fx(f.Types[i]); err != nil {
			return err
		} else if !ok {
			return nil
		}
	}
	return nil
}

func (f *File) EachConst(fx func(t *Constant) (bool, error)) error {
	for i := range f.Constants {
		if ok, err := fx(f.Constants[i]); err != nil {
			return err
		} else if !ok {
			return nil
		}
	}
	return nil
}

// Package represents a group of files that exist in the same file/module
type Package struct {
	Name       string  `json:"name"`  // Package name
	Files      []*File `json:"files"` // Files found in this package
	SourcePath string  `json:"-"`     // original source path
}

// Loop through all package files
func (pkg *Package) EachFile(fx func(f *File) (bool, error)) error {
	for i := range pkg.Files {
		if ok, err := fx(pkg.Files[i]); err != nil {
			return err
		} else if !ok {
			return nil
		}
	}
	return nil
}

// Loop through all package files
func (pkg *Package) EachInterface(fx func(f *File, iface *Interface) (bool, error)) error {
	return pkg.EachFile(func(f *File) (bool, error) {
		return true, f.EachInterface(func(iface *Interface) (bool, error) {
			return fx(f, iface)
		})
	})
}

// Loop through all package constants
func (pkg *Package) EachConstant(fx func(f *File, c *Constant) (bool, error)) error {
	return pkg.EachFile(func(f *File) (bool, error) {
		return true, f.EachConst(func(t *Constant) (bool, error) {
			return fx(f, t)
		})
	})
}

// Loop through all package constants
func (pkg *Package) EachTypeAlias(fx func(f *File, t *TypeAlias) (bool, error)) error {
	return pkg.EachFile(func(f *File) (bool, error) {
		return true, f.EachTypeAlias(func(t *TypeAlias) (bool, error) {
			return fx(f, t)
		})
	})
}

// Get all package constants
func (pkg *Package) Constants() []*Constant {
	ca := make([]*Constant, 0)
	_ = pkg.EachConstant(func(f *File, c *Constant) (bool, error) {
		ca = append(ca, c)
		return true, nil
	})
	return ca
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

type Reader interface {
	// Read file(s) and generate Package objects
	Read(config *ReadConfig) ([]*Package, error)
}

type WriteConfig struct {
	Output     Output       `json:"output"`
	Transforms []*Transform `json:"transforms,omitempty"`
	Packages   []*Package   `json:"packages,omitempty"`
}

type Writer interface {
	// Write package(s) with the provided output config and transforms
	// Return a map with output data for each package
	Write(config *WriteConfig) (map[string][]byte, error)
}
