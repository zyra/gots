package parser

import (
	"errors"
	"fmt"
	"github.com/zyra/gots/pkg/statement"
	"go/ast"
	"io/ioutil"
	"log"
	"path/filepath"
	"sync"
)

type Parser struct {
	*Config

	files []*ast.File

	wg sync.WaitGroup

	iMtx       sync.RWMutex
	interfaces []*Interface

	tMtx  sync.RWMutex
	types []*TypeDef

	cMtx      sync.RWMutex
	constants []*Constant

	pMtx     sync.Mutex
	pkgIndex map[string]string

	tsw *statement.Writer
}

func New(config *Config) *Parser {
	if !filepath.IsAbs(config.RootDir) {
		if d, err := filepath.Abs(config.RootDir); err != nil {
			log.Panicf("cannot convert base directory to absolute path: %s\n", err.Error())
		} else {
			config.RootDir = d
		}
	}

	return &Parser{
		Config: config,
		tsw:    statement.NewWriter(),
	}
}

func (p *Parser) Run() {
	p.parse()
	p.wg.Wait()
}

func (p *Parser) GenerateTS() {
	p.iMtx.RLock()
	defer p.iMtx.RUnlock()

	p.cMtx.RLock()
	defer p.cMtx.RUnlock()

	p.tMtx.RLock()
	defer p.tMtx.RUnlock()

	for _, it := range p.types {
		p.tsw.Export().Type(it.Name, statement.Literal(it.Type))
	}

	var lp int

	for _, it := range p.interfaces {
		lp = len(it.Properties)

		if lp == 0 {
			continue
		}

		properties := make([]*statement.Statement, lp, lp)

		for i, itt := range it.Properties {
			if itt.Optional {
				properties[i] = statement.OptionalProperty(itt.Name, itt.Type)
			} else {
				properties[i] = statement.Property(itt.Name, itt.Type)
			}
		}

		p.tsw.Export().Interface(it.Name, properties...)
	}

	for _, it := range p.constants {
		p.tsw.Export().Const(statement.Property(it.Name, it.Type), statement.Literal(it.Value))
	}
}

func (p *Parser) String() string {
	return p.tsw.String()
}

func (p *Parser) WriteToFile() error {
	if p.OutFileName == "" {
		return errors.New("output filename was not specified")
	}

	return ioutil.WriteFile(p.OutFileName, []byte(p.String()), 0644)
}

func (p *Parser) Print() {
	fmt.Print(p.String())
}
