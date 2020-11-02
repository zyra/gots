package golang

import (
	"fmt"
	"github.com/zyra/gots/pkg/parser/reader"
	"go/ast"
	"go/token"
)

type File struct {
	ast        *ast.File
	structs    []*Struct
	interfaces []*Interface
	types      []*TypeAlias
	constants  []*Const

	reader.File
}

func NewFile(file *ast.File) *File {
	return &File{
		ast:        file,
		structs:    make([]*Struct, 0),
		interfaces: make([]*Interface, 0),
		types:      make([]*TypeAlias, 0),
		constants:  make([]*Const, 0),
	}
}

// Inspect function
func (f *File) inspect(node ast.Node) (bool, error) {
	if node == nil {
		return true, nil
	}

	switch n := node.(type) {
	case *ast.File, *ast.Ident:
		return true, nil

	case *ast.FuncDecl, *ast.CommentGroup:
		return false, nil

	case *ast.GenDecl:
		if n.Specs == nil {
			return false, nil
		}

		switch n.Tok {
		case token.CONST:
			// TODO handle constants with no explicit values (numeric, iota+n ..etc)
			// 		it will need to be handled with access to all specs + list index..etc
			// 		since we need to check other props to know where the counter starts
			for i := range n.Specs {
				if spec, ok := n.Specs[i].(*ast.ValueSpec); ok {
					c, err := ConstFromValueSpec(spec)
					if err != nil {
						//return false, fmt.Errorf("failed to parse const: %v", err)
						continue
					}
					f.constants = append(f.constants, c)
				}
			}

		case token.TYPE:
			for i := range n.Specs {
				if spec, ok := n.Specs[i].(*ast.TypeSpec); ok {
					if err := f.ParseTypeSpec(spec); err != nil {
						return false, fmt.Errorf("failed to parse type: %v", err)
					}
				}
			}
		}

		return false, nil
	}

	return true, nil
}

// Parse file and populate type arrays
func (f *File) Parse() error {
	i := NewInspector(f.inspect)
	ast.Walk(i, f.ast)
	if err := i.Error(); err != nil {
		return err
	}

	for _, c := range f.constants {
		f.Constants = append(f.Constants, &c.Constant)
	}

	for _, s := range f.structs {
		f.Interfaces = append(f.Interfaces, &s.Interface)
	}

	return nil
}

func (f *File) ParseTypeSpec(spec *ast.TypeSpec) error {
	switch spec.Type.(type) {
	case *ast.StructType:
		st, err := ParseStruct(spec)
		if err != nil {
			return err
		}
		f.structs = append(f.structs, st)
		return nil
	case *ast.InterfaceType:
		it, err := ParseInterface(spec)
		if err != nil {
			return err
		}
		f.interfaces = append(f.interfaces, it)
		return nil

	case *ast.ArrayType, *ast.MapType, *ast.Ident:
		t := ParseTypeAlias(spec)
		f.types = append(f.types, t)
		return nil

	case *ast.SelectorExpr:
		return nil
	}

	return fmt.Errorf("unhandled type %t", spec.Type)
}
