package reader

import (
	"go/ast"
	"go/token"
)

// Package represents a group of files that exist in the same file/module
type Package struct {
	Name  string   `json:"name"`
	Files []string `json:"files"`

	PkgAst   *ast.Package `json:"-"`
	FilesAst []*ast.File  `json:"-"`

	Consts     map[string]*Constant  `json:"constants"` // constants
	Structs    map[string]*Interface `json:"structs"`
	Interfaces map[string]*Interface `json:"interfaces"` // interfaces
	Types      map[string]*TypeAlias `json:"types"`      // types
	Enums      map[string]*Enum      `json:"enums"`      // enums
	EnumValues map[string]*Constant  `json:"enumValues"`

	Decls          []ast.Decl
	ConstDecls     []*ast.GenDecl
	StructDecls    []*ast.TypeSpec
	InterfaceDecls []*ast.TypeSpec
	TypeDecls      []*ast.TypeSpec

	Path string `json:"path"` // original source path
}

func NewPackage(path string, pkg *ast.Package) *Package {
	p := &Package{
		Name:       pkg.Name,
		PkgAst:     pkg,
		Files:      make([]string, 0, len(pkg.Files)),
		Consts:     make(map[string]*Constant),
		Structs:    make(map[string]*Interface),
		Interfaces: make(map[string]*Interface),
		Types:      make(map[string]*TypeAlias),
		Enums:      make(map[string]*Enum),
		EnumValues: make(map[string]*Constant),

		Decls:          make([]ast.Decl, 0),
		ConstDecls:     make([]*ast.GenDecl, 0),
		StructDecls:    make([]*ast.TypeSpec, 0),
		InterfaceDecls: make([]*ast.TypeSpec, 0),
		TypeDecls:      make([]*ast.TypeSpec, 0),

		Path: path,
	}

	for k, v := range pkg.Files {
		p.FilesAst = append(p.FilesAst, v)
		p.Files = append(p.Files, k)
		p.Decls = append(p.Decls, v.Decls...)
	}

	// loop through file declarations
	for _, decl := range p.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			{
				// Generic declaration

				switch d.Tok {
				case token.IMPORT, token.COMMENT, token.PACKAGE, token.VAR:
					continue // TODO remove after debugging

				case token.CONST:
					{
						// const declaration
						p.ConstDecls = append(p.ConstDecls, d)
					}

				case token.TYPE:
					{
						// type declaration
						for _, spec := range d.Specs {
							switch spec := spec.(type) {
							case *ast.TypeSpec:
								if !ast.IsExported(spec.Name.Name) {
									continue
								}

								switch spec.Type.(type) {
								case *ast.Ident, *ast.ArrayType, *ast.SelectorExpr:
									p.TypeDecls = append(p.TypeDecls, spec)
								case *ast.StructType:
									p.StructDecls = append(p.StructDecls, spec)
								case *ast.InterfaceType:
									p.InterfaceDecls = append(p.InterfaceDecls, spec)
								case *ast.FuncType:
									continue
								default:
									continue
								}

							default:
								continue // TODO remove
							}
						}
					}

				default:
					continue // TODO remove after debugging
				}

			}

		case *ast.FuncDecl:
			continue

		default:
			panic("TODO: REMOVE THIS")
		}
	}

	return p
}
