package golang

import (
	"github.com/zyra/gots/pkg/parser/reader"
	"go/ast"
)

type Package struct {
	reader.Package
	ast   *ast.Package
	files []*File
}

func NewPackage(pkg *ast.Package, path string) *Package {
	p := Package{}

	p.Name = pkg.Name
	p.ast = pkg
	p.SourcePath = path
	p.Files = make([]*reader.File, 0)

	files := make([]*File, 0, len(pkg.Files))
	for i := range pkg.Files {
		f := NewFile(pkg.Files[i])
		f.SourcePath = i
		files = append(files, f)
		p.Files = append(p.Files, &f.File)
	}

	p.files = files
	return &p
}

func (pkg *Package) Parse() error {
	for _, it := range pkg.files {
		if err := it.Parse(); err != nil {
			return err
		}
	}
	return nil
}
