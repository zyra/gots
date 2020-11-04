package parser

import (
	"fmt"
	"github.com/zyra/gots/pkg/parser/golang"
	"github.com/zyra/gots/pkg/parser/reader"
	"github.com/zyra/gots/pkg/statement"
)

func (p *Parser) parse() ([]*reader.Package, error) {
	r := golang.Reader{}
	err := r.Read(&reader.ReadConfig{
		Dir:       p.RootDir,
		Recursive: p.Recursive,
	})
	if err != nil {
		return nil, err
	}

	for _, pkg := range r.Packages() {
		for _, iface := range pkg.Structs {
			newProps := make([]*reader.Property, 0)

			for _, prop := range iface.Properties {
				fromPkg := prop.Type.From
				if !prop.Inline {
					if len(fromPkg) > 0 {
						for _, pkg := range r.Packages() {
							if pkg.Name == fromPkg {
								goto appendProp
							}
						}

						prop.Type.Generic = true
					}

				appendProp:
					newProps = append(newProps, prop)
					continue
				}

				if fromPkg == "" {
					for _, iface2 := range pkg.Structs {
						if iface2.Name == prop.Type.Name {
							newProps = append(newProps, iface2.Properties...)
							break
						}
					}
				} else {
				PkgLoop:
					for _, pkg := range r.Packages() {
						if pkg.Name == fromPkg {
							for _, iface2 := range pkg.Structs {
								if iface2.Name == prop.Type.Name {
									newProps = append(newProps, iface2.Properties...)
									break PkgLoop
								}
							}
						}
					}
				}
			}

			iface.Properties = newProps
		}

	}

	return r.Packages(), nil
}

func tsToGo(in string) string {
	switch in {
	case "uint8", "uint16", "uint32", "int32", "int64", "int", "uint", "float", "float32", "float64":
		return "number"

	case "bool":
		return "boolean"
	}

	return in
}

func typeRef(t *reader.Type) *statement.Statement {
	if t.Map {
		return statement.Literal(fmt.Sprintf("{[key: string]: %s}", typeRef(t.MapValue)))
	}

	name := tsToGo(t.Name)
	gs := make([]*statement.Statement, 0)
	if t.Generic {
		gs = append(gs, statement.Literal("any"))
	} else if len(t.From) > 0 {
		gs = append(gs, statement.Literal(t.From), statement.Dot(name))
	} else {
		gs = append(gs, statement.Literal(name))
	}

	if t.Array {
		gs = append(gs, statement.Literal("[]"))
	}

	return statement.Group(gs...)
}

func (p *Parser) Run() error {
	packages, err := p.parse()
	if err != nil {
		return err
	}

	for _, pkg := range packages {
		tsw := statement.NewWriter()
		p.pkgIndex[pkg.Name] = pkg
		p.pkgTsw[pkg.Name] = tsw

		for _, t := range pkg.Types {
			tsw.Export().Type(t.Name, typeRef(t.AliasedType))
		}

		for _, iface := range pkg.Interfaces {
			props := make([]*statement.Statement, 0)
			for _, p := range iface.Properties {
				var s *statement.Statement
				if p.Optional {
					s = statement.OptionalProperty(p.Name, typeRef(p.Type).String())
				} else {
					s = statement.Property(p.Name, typeRef(p.Type).String())
				}
				props = append(props, s)
			}
			tsw.Export().Interface(iface.Name, props...)
		}

		for _, iface := range pkg.Structs {
			props := make([]*statement.Statement, 0)
			for _, p := range iface.Properties {
				var s *statement.Statement
				if p.Optional {
					s = statement.OptionalProperty(p.Name, typeRef(p.Type).String())
				} else {
					s = statement.Property(p.Name, typeRef(p.Type).String())
				}
				props = append(props, s)
			}
			tsw.Export().Interface(iface.Name, props...)
		}

		for _, c := range pkg.Consts {
			if c.Type.EnumValue {
				tsw.Export().Const(statement.Property(c.Name, typeRef(c.Type).String()), statement.Group(statement.Literal(c.Type.Name), statement.Dot(c.Value)))
			} else {
				tsw.Export().Const(statement.Property(c.Name, typeRef(c.Type).String()), statement.Literal(c.Value))
			}
		}

		for _, e := range pkg.Enums {
			vals := make([]*statement.Statement, len(e.Values))
			for i, v := range e.Values {
				vals[i] = statement.EnumValue(v.Name, statement.Literal(v.Value))
			}
			tsw.Export().Enum(e.Name, vals...)
		}
	}

	return nil
}
