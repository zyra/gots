package parser

import (
	"go/ast"
	"sync"
)

type Package struct {

	files []*ast.File

	wg *sync.WaitGroup

	iMtx       sync.RWMutex
	interfaces []*Interface

	tMtx  sync.RWMutex
	types []*TypeDef

	cMtx      sync.RWMutex
	constants []*Constant
}
