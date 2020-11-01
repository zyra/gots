package parser

import (
	"errors"
	"fmt"
	"github.com/zyra/gots/pkg/parser/tag"
	"go/ast"
	"go/token"
	"sync"
)

type File struct {
	ast *ast.File
	pkg *Package

	wg *sync.WaitGroup

	sMtx    *sync.RWMutex
	structs []*Struct

	iMtx       *sync.RWMutex
	interfaces []*Interface

	tMtx  *sync.RWMutex
	types []*TypeAlias

	cMtx      *sync.RWMutex
	constants []*Const
}

func NewFile(file *ast.File) *File {
	return &File{
		ast:        file,
		wg:         new(sync.WaitGroup),
		sMtx:       new(sync.RWMutex),
		structs:    make([]*Struct, 0),
		iMtx:       new(sync.RWMutex),
		interfaces: make([]*Interface, 0),
		tMtx:       new(sync.RWMutex),
		types:      make([]*TypeAlias, 0),
		cMtx:       new(sync.RWMutex),
		constants:  make([]*Const, 0),
	}
}

func (f *File) Parse() error {
	var parseErr error
	ast.Inspect(f.ast, func(node ast.Node) bool {
		if node == nil {
			return true
		}

		switch n := node.(type) {
		case *ast.File, *ast.Ident:
			return true

		case *ast.FuncDecl, *ast.CommentGroup:
			return false

		case *ast.GenDecl:
			if n.Specs == nil {
				return false
			}

			switch n.Tok {
			case token.CONST:
				// TODO handle constants with no explicit values (numeric, iota+n ..etc)
				// 		it will need to be handled with access to all specs + list index..etc
				// 		since we need to check other props to know where the counter starts
				for i := range n.Specs {
					if spec, ok := n.Specs[i].(*ast.ValueSpec); ok {
						f.wg.Add(1)
						go func(spec *ast.ValueSpec) {
							defer f.wg.Done()

							c, err := ConstFromValueSpec(spec)
							if err != nil {
								parseErr = err
								return
							}
							f.constants = append(f.constants, c)
						}(spec)
					}
				}

			case token.TYPE:
				for i := range n.Specs {
					if spec, ok := n.Specs[i].(*ast.TypeSpec); ok {
						f.wg.Add(1)
						go func(spec *ast.TypeSpec) {
							defer f.wg.Done()
							if err := f.ParseTypeSpec(spec); err != nil {
								parseErr = err
								return
							}
						}(spec)
					}
				}
			}

			return false
		}

		return true
	})
	f.wg.Wait()
	return parseErr
}

func (f *File) ParseTypeAlias(spec *ast.TypeSpec) (*TypeAlias, error) {
	ident, ok := spec.Type.(*ast.Ident)
	var isArray bool
	if !ok {
		if arrIdent, ok := spec.Type.(*ast.ArrayType); ok {
			isArray = true
			if ident, ok = arrIdent.Elt.(*ast.Ident); ok {
				goto parse
			}

			if star, ok := arrIdent.Elt.(*ast.StarExpr); ok {
				if ident, ok = star.X.(*ast.Ident); ok {
					goto parse
				}
			}
		}

		return nil, errors.New("invalid TypeSpec")
	}

parse:
	t := ParseType(ident)
	t.Array = isArray

	if spec.Name.Name == t.Name {
		return nil, nil
	}

	ta := &TypeAlias{
		Name: spec.Name.Name,
		Type: t,
	}

	return ta, nil
}

func (f *File) ParseTypeSpec(spec *ast.TypeSpec) error {
	switch spec.Type.(type) {
	case *ast.StructType:
		st, err := f.ParseStruct(spec)
		if err != nil {
			return err
		}
		f.sMtx.Lock()
		f.structs = append(f.structs, st)
		f.sMtx.Unlock()
		return nil
	case *ast.Ident:
		out, err := f.ParseTypeAlias(spec)
		if err != nil {
			return err
		}
		f.tMtx.Lock()
		f.types = append(f.types, out)
		f.tMtx.Unlock()
		return nil

	case *ast.InterfaceType:
		it, err := f.ParseInterface(spec)
		if err != nil {
			return err
		}
		f.iMtx.Lock()
		f.interfaces = append(f.interfaces, it)
		f.iMtx.Unlock()
		return nil

	case *ast.SelectorExpr:
		return nil

	case *ast.ArrayType:
		out, err := f.ParseTypeAlias(spec)
		if err != nil {
			return err
		}
		out.Type.Array = true
		f.tMtx.Lock()
		f.types = append(f.types, out)
		f.tMtx.Unlock()
		return nil
	}

	return fmt.Errorf("unhandled type %t", spec.Type)
}

func (f *File) ParseInterface(spec *ast.TypeSpec) (*Interface, error) {
	_, ok := spec.Type.(*ast.InterfaceType)
	if !ok {
		return nil, errors.New("invalid TypeSpec")
	}
	name := spec.Name.Name
	return &Interface{Name: name}, nil
}

func (f *File) ParseStruct(spec *ast.TypeSpec) (*Struct, error) {
	s, ok := spec.Type.(*ast.StructType)
	if !ok {
		return nil, errors.New("invalid TypeSpec")
	}

	itName := spec.Name.Name
	if !ast.IsExported(itName) {
		return nil, ErrNotExported
	}

	var props []*Property
	if s.Fields != nil {
		nf := s.Fields.NumFields()
		if nf != 0 {
			props = make([]*Property, 0, nf)
			for i := range s.Fields.List {
				f := s.Fields.List[i]
				t, err := tag.ParseTag(f.Tag.Value)

				propType := ParseType(f.Type)

				prop := Property{}

				if err != nil {
					if err == tag.ErrJsonIgnored || err == tag.ErrJsonTagNotPresent || err == tag.ErrPropertyIgnored {
						continue
					}

					if err == tag.ErrPropertyInlined {
						prop.Inline = true
					} else {
						return nil, fmt.Errorf("failed to parse tag: %v", err)
					}
				} else {
					if t.Type != "" {
						propType.Name = t.Type
						propType.Import = nil
					}

					prop.Name = t.Name
					prop.Optional = t.Optional
				}

				prop.Type = propType
				props = append(props, &prop)
			}
		}
	}

	st := Struct{
		Name:       itName,
		Properties: props,
	}

	return &st, nil
}

func (f *File) Render() {

}
