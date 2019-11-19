package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

func (p *Parser) parse() {
	fset := token.NewFileSet()

	var pkgs map[string]*ast.Package
	var err error

	if pkgs, err = parser.ParseDir(fset, p.BaseDir, nil, parser.ParseComments); err != nil {
		log.Panicf("unable to parse base directory: %s\n", err.Error())
	}

	for k := range pkgs {
		for fk := range pkgs[k].Files {
			p.wg.Add(1)
			go p.parseFile(pkgs[k].Files[fk])
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
		Name:  spec.Names[0].Name,
		Type:  parseType(spec.Type),
		Value: spec.Values[0].(*ast.BasicLit).Value,
	}

	p.cMtx.Lock()
	defer p.cMtx.Unlock()

	p.constants = append(p.constants, c)
}

func (p *Parser) parseStruct(spec *ast.TypeSpec) {
	var s *ast.StructType
	var ok bool

	if s, ok = spec.Type.(*ast.StructType); ! ok {
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
			return
		}

		tag, err := parseTags(f.Tag.Value)

		if err != nil {
			return
		}

		if tag.Type == "" {
			tag.Type = parseType(f.Type)
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
		t := parseType(spec.Type)
		p.tMtx.Lock()
		defer p.tMtx.Unlock()
		p.types = append(p.types, &TypeDef{
			Name: spec.Name.Name,
			Type: t,
		})
		return
	default:
		panic("??")
	}
}
