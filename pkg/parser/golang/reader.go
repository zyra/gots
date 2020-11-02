package golang

import (
	"fmt"
	"github.com/zyra/gots/pkg/parser/reader"
	goparser "go/parser"
	"go/token"
)

type Reader struct{}

func (r *Reader) Read(rc *reader.ReadConfig) ([]*reader.Package, error) {
	fSet := token.NewFileSet()
	dirs, err := rc.Directories()
	if err != nil {
		return nil, err
	}

	out := make([]*reader.Package, len(dirs))

	for i, it := range dirs {
		pkgs, err := goparser.ParseDir(fSet, it, nil, goparser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("failed to parse file %v: %v", it, err)
		}

		for _, p := range pkgs {
			np := NewPackage(p, it)
			if err := np.Parse(); err != nil {
				return nil, err
			}
			out[i] = &np.Package
		}
	}

	return out, nil
}
