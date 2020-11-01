package parser

import (
	"go/ast"
	"sync"
)

func NewPackage(pkg *ast.Package) *Package {
	files := make([]*File, 0, len(pkg.Files))
	for _, fv := range pkg.Files {
		files = append(files, NewFile(fv))
	}

	return &Package{
		ast:        pkg,
		wg:         new(sync.WaitGroup),
		fMtx:       new(sync.RWMutex),
		files:      files,
		iMtx:       new(sync.RWMutex),
		interfaces: make([]*Interface, 0),
		tMtx:       new(sync.RWMutex),
		types:      make([]*TypeAlias, 0),
		sMtx:       new(sync.RWMutex),
		structs:    make([]*Struct, 0),
		cMtx:       new(sync.RWMutex),
		constants:  make([]*Const, 0),
		errChan:    make(chan error, 1024),
	}
}

type Package struct {
	ast *ast.Package

	wg *sync.WaitGroup

	errChan chan error

	fMtx  *sync.RWMutex
	files []*File

	iMtx       *sync.RWMutex
	interfaces []*Interface

	tMtx  *sync.RWMutex
	types []*TypeAlias

	sMtx    *sync.RWMutex
	structs []*Struct

	cMtx      *sync.RWMutex
	constants []*Const
}

func (pkg *Package) Parse() {
	pkg.wg.Add(len(pkg.files))
	for i := range pkg.files {
		go func(f *File) {
			defer pkg.wg.Done()
			if err := f.Parse(); err != nil {
				pkg.errChan <- err
			}
		}(pkg.files[i])
	}

	pkg.wg.Wait()

	for _, f := range pkg.files {
		for _, s := range f.structs {
			pkg.AddStruct(s)
		}
		for _, it := range f.interfaces {
			pkg.AddInterface(it)
		}
		for _, c := range f.constants {
			pkg.AddConst(c)
		}
		for _, t := range f.types {
			pkg.AddTypeAlias(t)
		}
	}
}

func (pkg *Package) AddFile(it *File) {
	pkg.fMtx.Lock()
	pkg.files = append(pkg.files, it)
	pkg.fMtx.Unlock()
}

func (pkg *Package) AddInterface(it *Interface) {
	pkg.iMtx.Lock()
	pkg.interfaces = append(pkg.interfaces, it)
	pkg.iMtx.Unlock()
}

func (pkg *Package) AddTypeAlias(it *TypeAlias) {
	pkg.tMtx.Lock()
	pkg.types = append(pkg.types, it)
	pkg.tMtx.Unlock()
}

func (pkg *Package) AddConst(it *Const) {
	pkg.cMtx.Lock()
	pkg.constants = append(pkg.constants, it)
	pkg.cMtx.Unlock()
}

func (pkg *Package) AddStruct(it *Struct) {
	pkg.sMtx.Lock()
	pkg.structs = append(pkg.structs, it)
	pkg.sMtx.Unlock()
}
