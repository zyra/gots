package parser

import (
	"fmt"
	"github.com/zyra/gots/pkg/parser/golang"
	"go/ast"
	"go/token"
)

type File struct {
	ast        *ast.File
	pkg        *Package
	structs    []*golang.Struct
	interfaces []*golang.Interface
	types      []*golang.TypeAlias
	constants  []*golang.Const
}

func NewFile(file *ast.File) *File {
	return &File{
		ast:        file,
		structs:    make([]*golang.Struct, 0),
		interfaces: make([]*golang.Interface, 0),
		types:      make([]*golang.TypeAlias, 0),
		constants:  make([]*golang.Const, 0),
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
					c, err := golang.ConstFromValueSpec(spec)
					if err != nil {
						return false, fmt.Errorf("failed to parse const: %v", err)
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
	i := Inspector{inspect: f.inspect}
	ast.Walk(&i, f.ast)
	return i.Error()
}

func (f *File) ParseTypeSpec(spec *ast.TypeSpec) error {
	switch spec.Type.(type) {
	case *ast.StructType:
		st, err := golang.ParseStruct(spec)
		if err != nil {
			return err
		}
		f.structs = append(f.structs, st)
		return nil
	case *ast.InterfaceType:
		it, err := golang.ParseInterface(spec)
		if err != nil {
			return err
		}
		f.interfaces = append(f.interfaces, it)
		return nil

	case *ast.ArrayType, *ast.MapType, *ast.Ident:
		t := golang.ParseTypeAlias(spec)
		f.types = append(f.types, t)
		return nil

	case *ast.SelectorExpr:
		return nil
	}

	return fmt.Errorf("unhandled type %t", spec.Type)
}
