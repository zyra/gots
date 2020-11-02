package parser

import (
	"github.com/zyra/gots/pkg/parser/golang"
	"github.com/zyra/gots/pkg/parser/reader"
)

func (p *Parser) parse() ([]*reader.Package, error) {
	r := golang.Reader{}
	packages, err := r.Read(&reader.ReadConfig{
		Dir:       p.RootDir,
		Recursive: p.Recursive,
	})
	if err != nil {
		return nil, err
	}

	for _, pkg := range packages {
		if err := pkg.EachFile(func(f *reader.File) (bool, error) {
			for _, iface := range f.Interfaces {
				newProps := make([]*reader.Property, 0)
				for _, prop := range iface.Properties {
					if !prop.Inline {
						newProps = append(newProps, prop)
						continue
					}

					fromPkg := prop.Type.From
					if fromPkg == "" {
						if err := pkg.EachInterface(func(f *reader.File, iface *reader.Interface) (bool, error) {
							if iface.Name != prop.Type.Name {
								return true, nil
							}
							newProps = append(newProps, iface.Properties...)
							return false, nil
						}); err != nil {
							return false, err
						}
					} else {
						for _, pkg2 := range packages {
							if pkg2.Name != fromPkg {
								continue
							}

							if err := pkg2.EachInterface(func(f *reader.File, iface2 *reader.Interface) (bool, error) {
								if iface2.Name != prop.Type.Name {
									return true, nil
								}
								newProps = append(newProps, iface2.Properties...)
								return false, nil
							}); err != nil {
								return false, err
							}
						}
					}
				}
				iface.Properties = newProps
			}
			return true, nil
		}); err != nil {
			return nil, err
		}
	}

	return packages, nil
}
