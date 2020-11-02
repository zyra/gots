package parser

import (
	"fmt"
	"github.com/zyra/gots/pkg/parser/reader"
	"github.com/zyra/gots/pkg/statement"
	"io/ioutil"
	"log"
	"path/filepath"
)

type Parser struct {
	*reader.Config
	pkgIndex map[string]string
	tsw      *statement.Writer
}

func New(config *reader.Config) *Parser {
	if !filepath.IsAbs(config.RootDir) {
		if d, err := filepath.Abs(config.RootDir); err != nil {
			log.Panicf("cannot convert base directory to absolute path: %s\n", err.Error())
		} else {
			config.RootDir = d
		}
	}

	return &Parser{
		Config:   config,
		pkgIndex: make(map[string]string),
		tsw:      statement.NewWriter(),
	}
}

func (p *Parser) Run() error {
	packages, err := p.parse()
	if err != nil {
		return err
	}

	for _, pkg := range packages {
		pkg.EachTypeAlias(func(f *reader.File, t *reader.TypeAlias) (bool, error) {
			p.tsw.Export().Type(t.Name, statement.Literal(t.AliasedType.TSType()))

			return true, nil
		})

		pkg.EachInterface(func(f *reader.File, iface *reader.Interface) (bool, error) {
			props := make([]*statement.Statement, 0)
			for _, p := range iface.Properties {
				var s *statement.Statement
				if p.Optional {
					s = statement.OptionalProperty(p.Name, p.Type.TSType())
				} else {
					s = statement.Property(p.Name, p.Type.TSType())
				}
				props = append(props, s)
			}
			p.tsw.Export().Interface(iface.Name, props...)

			return true, nil
		})

		pkg.EachConstant(func(f *reader.File, c *reader.Constant) (bool, error) {
			p.tsw.Export().Const(statement.Property(c.Name, c.Type.TSType()), statement.Literal(c.Value))
			return true, nil
		})
	}

	return nil
}

func (p *Parser) String() string {
	return p.tsw.String()
}

func (p *Parser) WriteToFile() error {
	//if p.OutFileName == "" {
	//	return errors.New("output filename was not specified")
	//}

	return ioutil.WriteFile(p.Output.AIOFileName, []byte(p.String()), 0644)
}

func (p *Parser) Print() {
	fmt.Print(p.String())
}
