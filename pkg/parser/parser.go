package parser

import (
	"fmt"
	"github.com/zyra/gots/pkg/parser/godef"
	"github.com/zyra/gots/pkg/statement"
	"io/ioutil"
	"log"
	"path/filepath"
	"sync"
)

type Parser struct {
	*Config

	wg *sync.WaitGroup

	iMtx       *sync.RWMutex
	interfaces []*godef.Interface

	tMtx  *sync.RWMutex
	types []*godef.TypeAlias

	cMtx      *sync.RWMutex
	constants []*godef.Const

	sMtx    *sync.RWMutex
	structs []*godef.Struct

	pMtx *sync.RWMutex
	pkgs []*Package

	fMtx  *sync.RWMutex
	files []*File

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
		Config:     config,
		wg:         new(sync.WaitGroup),
		iMtx:       new(sync.RWMutex),
		interfaces: make([]*godef.Interface, 0),
		tMtx:       new(sync.RWMutex),
		types:      make([]*godef.TypeAlias, 0),
		cMtx:       new(sync.RWMutex),
		constants:  make([]*godef.Const, 0),
		sMtx:       new(sync.RWMutex),
		structs:    make([]*godef.Struct, 0),
		pMtx:       new(sync.RWMutex),
		pkgs:       make([]*Package, 0),
		pkgIndex:   make(map[string]string),
		tsw:        statement.NewWriter(),
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
		p.tsw.Export().Type(it.Name, statement.Literal(it.Type.Name))
	}

	var lp int

	for _, it := range p.structs {
		lp = len(it.Properties)

		if lp == 0 {
			continue
		}

		properties := make([]*statement.Statement, lp, lp)

		for i, itt := range it.Properties {
			if itt.Optional {
				properties[i] = statement.OptionalProperty(itt.Name, itt.Type.Name)
			} else {
				properties[i] = statement.Property(itt.Name, itt.Type.Name)
			}
		}

		p.tsw.Export().Interface(it.Name, properties...)
	}

	for _, it := range p.constants {
		if it.Type == nil {
			it.Type = &godef.Type{}
		}
		p.tsw.Export().Const(statement.Property(it.Name, it.Type.Name), statement.Literal(it.Value))
	}
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
