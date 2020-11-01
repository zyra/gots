package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"path/filepath"
)

func (p *Parser) parse() error {
	fset := token.NewFileSet()

	var pkgs map[string]*ast.Package
	var err error

	pkgIndex := make(map[string]string)

	var scanDir func(path string)

	scanDir = func(path string) {
		contents, err := ioutil.ReadDir(path)

		if err != nil {
			log.Panicf("unable to read directory %s: %s\n", path, err.Error())
		}

		for _, it := range contents {
			if it.IsDir() {
				scanDir(filepath.Join(path, it.Name()))
			}
		}

		if pkgs, err = parser.ParseDir(fset, path, nil, parser.PackageClauseOnly); err != nil {
			log.Panicf("unable to scan directory %s: %s\n", path, err.Error())
		} else {
			for k := range pkgs {
				pkgIndex[k] = path
			}
		}
	}

	scanDir(p.RootDir)

	p.pkgIndex = pkgIndex

	for _, v := range pkgIndex {
		if pkgs, err = parser.ParseDir(fset, v, nil, parser.ParseComments); err != nil {
			return fmt.Errorf("unable to parse base directory: %s\n", err.Error())
		}

		p.wg.Add(len(pkgs))

		for _, pv := range pkgs {
			pkg := NewPackage(pv)
			p.pkgs = append(p.pkgs, pkg)

			go func(pkg *Package) {
				defer p.wg.Done()
				pkg.Parse()
			}(pkg)
		}
	}

	p.wg.Wait()

	for _, pkg := range p.pkgs {
		for _, file := range pkg.files {
			p.files = append(p.files, file)
		}
		for _, c := range pkg.constants {
			p.constants = append(p.constants, c)
		}
		for _, it := range pkg.interfaces {
			p.interfaces = append(p.interfaces, it)
		}
		for _, st := range pkg.structs {
			p.structs = append(p.structs, st)
		}
		for _, t := range pkg.types {
			p.types = append(p.types, t)
		}
	}

	return nil
}
