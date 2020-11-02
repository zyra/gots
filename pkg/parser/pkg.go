package parser

import (
	"go/ast"
	"sync"
)

func NewPackage(pkg *ast.Package) *PkgReader {
	//files := make([]*golang.PkgFile, 0, len(pkg.Files))
	//for _, fv := range pkg.Files {
	//	files = append(files, golang.NewFile(fv))
	//}

	return &PkgReader{
		ast: pkg,
		//wg:         new(sync.WaitGroup),
		//fMtx:       new(sync.RWMutex),
		//files:      files,
		//iMtx:       new(sync.RWMutex),
		//interfaces: make([]*golang.Interface, 0),
		//tMtx:       new(sync.RWMutex),
		//types:      make([]*golang.TypeAlias, 0),
		//sMtx:       new(sync.RWMutex),
		//structs:    make([]*golang.Struct, 0),
		//cMtx:       new(sync.RWMutex),
		//constants:  make([]*golang.Const, 0),
		errChan: make(chan error, 1024),
	}
}

type PkgReader struct {
	ast *ast.Package

	wg *sync.WaitGroup

	errChan chan error

	//fMtx  *sync.RWMutex
	//files []*golang.PkgFile
	//
	//iMtx       *sync.RWMutex
	//interfaces []*golang.Interface
	//
	//tMtx  *sync.RWMutex
	//types []*golang.TypeAlias
	//
	//sMtx    *sync.RWMutex
	//structs []*golang.Struct
	//
	//cMtx      *sync.RWMutex
	//constants []*golang.Const
}

func (pkg *PkgReader) Parse() {
	//pkg.wg.Add(len(pkg.files))
	//for i := range pkg.files {
	//	go func(f *golang.PkgFile) {
	//		defer pkg.wg.Done()
	//		if err := f.Parse(); err != nil {
	//			pkg.errChan <- err
	//		}
	//	}(pkg.files[i])
	//}
	//
	//pkg.wg.Wait()
	//
	//for _, f := range pkg.files {
	//	for _, s := range f.structs {
	//		pkg.AddStruct(s)
	//	}
	//	for _, it := range f.interfaces {
	//		pkg.AddInterface(it)
	//	}
	//	for _, c := range f.constants {
	//		pkg.AddConst(c)
	//	}
	//	for _, t := range f.types {
	//		pkg.AddTypeAlias(t)
	//	}
	//}
}

//
//func (pkg *PkgReader) AddFile(it *golang.PkgFile) {
//	//pkg.fMtx.Lock()
//	//pkg.files = append(pkg.files, it)
//	//pkg.fMtx.Unlock()
//}
//
//func (pkg *PkgReader) AddInterface(it *golang.Interface) {
//	//pkg.iMtx.Lock()
//	//pkg.interfaces = append(pkg.interfaces, it)
//	//pkg.iMtx.Unlock()
//}
//
//func (pkg *PkgReader) AddTypeAlias(it *golang.TypeAlias) {
//	//pkg.tMtx.Lock()
//	//pkg.types = append(pkg.types, it)
//	//pkg.tMtx.Unlock()
//}
//
//func (pkg *PkgReader) AddConst(it *golang.Const) {
//	//pkg.cMtx.Lock()
//	//pkg.constants = append(pkg.constants, it)
//	//pkg.cMtx.Unlock()
//}
//
//func (pkg *PkgReader) AddStruct(it *golang.Struct) {
//	//pkg.sMtx.Lock()
//	//pkg.structs = append(pkg.structs, it)
//	//pkg.sMtx.Unlock()
//}
