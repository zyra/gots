package golang

import (
	"fmt"
	"github.com/zyra/gots/pkg/parser/reader"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"strconv"
)

type Reader struct {
	packages []*reader.Package
}

func (r *Reader) Packages() []*reader.Package {
	return r.packages
}

func (r *Reader) Read(rc *reader.ReadConfig) error {
	// get list of all the directories that we should scan for go packages/code
	dirs, err := rc.Directories()
	if err != nil {
		return err
	}

	fSet := token.NewFileSet()
	r.packages = make([]*reader.Package, 0)

	for _, dir := range dirs {
		// Load AST from the current dir
		packages, err := goparser.ParseDir(fSet, dir, nil, goparser.ParseComments)
		if err != nil {
			return fmt.Errorf("failed to parse file %v: %v", dir, err)
		}

		// loop through packages
		// should be 1 package in most cases
		for _, pkg := range packages {
			p := reader.NewPackage(dir, pkg)
			r.packages = append(r.packages, p)

			// process types
			for _, d := range p.TypeDecls {
				p.Types[d.Name.Name] = ParseTypeAlias(d)
			}

			// interfaces
			for _, it := range p.InterfaceDecls {
				p.Interfaces[it.Name.Name], err = ParseInterface(it)
				if err != nil {
					return fmt.Errorf("failed to parse interface %s: %v", it.Name.Name, err)
				}
			}

			// structs
			for _, s := range p.StructDecls {
				p.Structs[s.Name.Name], err = ParseStruct(s)
				if err != nil {
					return fmt.Errorf("failed to parse struct %s: %v", s.Name.Name, err)
				}
			}

			// process consts
			for _, d := range p.ConstDecls {
				// const declaration
				consts := make([]*reader.Constant, len(d.Specs))
				for i, spec := range d.Specs {
					spec := spec.(*ast.ValueSpec)
					//name := spec.Names[0].Name

					c, err := ConstFromValueSpec(spec)
					if err != nil {
						continue
						//return fmt.Errorf("failed to read const %s: %v", name, err)
					}
					consts[i] = c

					if len(c.Value) == 0 {
						// this value is most likely part of an enum

						if i == 0 || len(consts) < i {
							panic("shouldn't happen -- TODO remove/fix this")
						}

						prev := consts[i-1]
						if !prev.EnumValue {
							panic("shouldn't happen -- TODO remove/fix this")
						}

						intVal, err := strconv.Atoi(prev.Value)
						if err != nil {
							panic("shouldn't happen -- TODO remove/fix this")
						}

						c.Value = strconv.Itoa(intVal + 1)
						c.Type = prev.Type
						c.EnumName = prev.EnumName
						c.EnumValue = true
					}

					if c.Type == nil || c.Type.IsEmpty() {
						if res, ok := p.EnumValues[c.Value]; ok {
							c.Type = &reader.Type{
								Name:      res.Type.Name,
								From:      res.Type.From,
								EnumValue: true,
							}
						} else {
							for _, file := range p.FilesAst {
								res := file.Scope.Lookup(c.Value)
								if res != nil {
									if tt, ok := res.Type.(ast.Expr); ok {
										c.Type = TypeFromExpr(tt)
									} else if tt, ok := res.Decl.(*ast.ValueSpec); ok {
										// const refers to an enum value
										cc, err := ConstFromValueSpec(tt)
										if err != nil {
											return fmt.Errorf("failed to parse const: %v", err)
										}

										c.Type = cc.Type
										c.Type.EnumValue = true
									} else {
										panic("TODO: remove/fix")
									}
									break
								}
							}
						}
					}

					if res, ok := p.Enums[c.Type.Name]; ok {
						c.EnumName = res.Name
						c.EnumValue = true
						res.Values = append(res.Values, c)
						p.EnumValues[c.Name] = c
						continue
					}

					if res, ok := p.Types[c.Type.Name]; ok {
						c.EnumValue = true
						c.EnumName = c.Type.Name
						p.EnumValues[c.Name] = c

						p.Enums[res.Name] = &reader.Enum{
							Name:   res.Name,
							Values: []*reader.Constant{c},
						}
						delete(p.Types, res.Name)
						continue
					}

					p.Consts[c.Name] = c
				}
			}
		}
	}

	return nil
}
