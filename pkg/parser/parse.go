package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"path/filepath"
)

func (p *Parser) parse() {
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
			log.Panicf("unable to parse base directory: %s\n", err.Error())
		}

		for pk := range pkgs {
			for fk := range pkgs[pk].Files {
				p.wg.Add(1)
				go p.parseFile(pkgs[pk].Files[fk])
			}
		}
	}
}

func (p *Parser) parseFile(file *ast.File) {
	defer p.wg.Done()

	ast.Inspect(file, func(node ast.Node) bool {
		if node == nil {
			return true
		}

		switch node.(type) {
		case *ast.File, *ast.Ident:
			return true

		case *ast.FuncDecl, *ast.CommentGroup:
			return false

		case *ast.GenDecl:
			n := node.(*ast.GenDecl)

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
						p.wg.Add(1)
						go func() {
							defer p.wg.Done()
							p.parseConst(spec)
						}()
					}
				}

			case token.TYPE:
				for i := range n.Specs {
					if spec, ok := n.Specs[i].(*ast.TypeSpec); ok {
						p.wg.Add(1)
						go func() {
							defer p.wg.Done()
							p.parseTypeSpec(spec)
						}()
					}
				}
			}

			return false
		}

		return true
	})
}

func (p *Parser) parseConst(spec *ast.ValueSpec) {
	c := &Constant{
		Name: spec.Names[0].Name,
	}

	if spec.Type != nil {
		c.Type = ParseTagType(spec.Type)
		c.Value = spec.Values[0].(*ast.BasicLit).Value
	} else {
		switch spec.Values[0].(type) {
		case *ast.CallExpr:
			if val, ok := spec.Values[0].(*ast.CallExpr).Args[0].(*ast.BasicLit); ok {
				c.Type = ParseTypeFromToken(val.Kind)
				c.Value = val.Value
			} else {
				panic("Unhandled case")
			}
		case *ast.BasicLit:
			v := spec.Values[0].(*ast.BasicLit)
			c.Type = ParseTypeFromToken(v.Kind)
			c.Value = v.Value
		default:
			panic("Unhandled case")
		}
	}

	if c.Value == "" {
		panic("Unhandled case")
	}

	p.cMtx.Lock()
	defer p.cMtx.Unlock()

	p.constants = append(p.constants, c)
}

func (p *Parser) parseStruct(spec *ast.TypeSpec) {
	var s *ast.StructType
	var ok bool

	if s, ok = spec.Type.(*ast.StructType); !ok {
		return
	}

	if s.Fields == nil {
		return
	}

	nf := s.Fields.NumFields()

	if nf == 0 {
		return
	}

	props := make([]*Property, 0, nf)

	for _, f := range s.Fields.List {
		if f.Tag == nil || f.Tag.Value == "" {
			continue
		}

		tag, err := ParseTag(f.Tag.Value)

		if err != nil {
			continue
		}

		if tag.Type == "" {
			tag.Type = ParseTagType(f.Type)
		}

		props = append(props, &Property{
			Name:     tag.Name,
			Type:     tag.Type,
			Optional: tag.Optional,
		})
	}

	p.iMtx.Lock()
	defer p.iMtx.Unlock()
	p.interfaces = append(p.interfaces, &Interface{
		Name:       spec.Name.Name,
		Properties: props,
	})
}

func (p *Parser) parseTypeSpec(spec *ast.TypeSpec) {
	switch spec.Type.(type) {
	case *ast.StructType:
		p.parseStruct(spec)
	case *ast.Ident:
		t := ParseTagType(spec.Type)

		if spec.Name.Name == t {
			return
		}

		p.tMtx.Lock()
		defer p.tMtx.Unlock()
		p.types = append(p.types, &TypeDef{
			Name: spec.Name.Name,
			Type: t,
		})
		return
	case *ast.InterfaceType:
		return
	case *ast.SelectorExpr:
		st := spec.Type.(*ast.SelectorExpr)
		t := "any"
		if xv, ok := st.X.(*ast.Ident); !ok {
			panic("unhandled case")
		} else {
			p.pMtx.Lock()
			defer p.pMtx.Unlock()
			if _, ok := p.pkgIndex[xv.Name]; ok {
				t = st.Sel.Name
			}
		}

		if spec.Name.Name == t {
			return
		}

		p.tMtx.Lock()
		defer p.tMtx.Unlock()
		p.types = append(p.types, &TypeDef{
			Name: spec.Name.Name,
			Type: t,
		})
		return
	default:
		panic(spec.Type)
	}
}
