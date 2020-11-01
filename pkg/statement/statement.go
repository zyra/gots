package statement

import (
	"strings"
	"sync"
)

// A statement is the base type to generate TypeScript syntax
// Statements are used to generate simple TypeScript code,
// which can be grouped to declare variables, types, or complex functions
type Statement struct {
	// Statement prefix
	Prefix string

	// Statement suffix
	Suffix string

	// Separator to use when joining children
	Separator string

	// Value is typically used when the statement has a single child or a single value
	Value string

	// Statement children
	Children []*Statement

	// Mutex to lock children slice
	cMtx *sync.RWMutex
}

// Create a new blank statement
func NewStatement() *Statement {
	return &Statement{
		Children: make([]*Statement, 0),
		cMtx:     new(sync.RWMutex),
	}
}

// Create a new statement with the provided children slice
func NewStatementWithChildren(children []*Statement) *Statement {
	return &Statement{
		Children: children,
		cMtx:     new(sync.RWMutex),
	}
}

// Create a new statement with the provided value
func NewStatementWithValue(val string) *Statement {
	s := NewStatement()
	s.Value = val
	return s
}

// Append a new child statement
func (d *Statement) Append(c *Statement) *Statement {
	d.cMtx.Lock()
	d.Children = append(d.Children, c)
	d.cMtx.Unlock()
	return d
}

// Converts the statement to a TypeScript code block
// Results in a string in the following format (without spaces in between):
// [prefix] [value] [children [joined by separator]] [suffix]
func (d *Statement) String() string {
	ss := make([]string, 0, 4)

	if len(d.Prefix) > 0 {
		ss = append(ss, d.Prefix)
	}

	if len(d.Value) > 0 {
		ss = append(ss, d.Value)
	}

	cl := len(d.Children)
	if cl > 0 {
		cs := make([]string, cl)
		for i := 0; i < cl; i++ {
			cs[i] = d.Children[i].String()
		}
		ss = append(ss, strings.Join(cs, d.Separator))
	}

	if len(d.Suffix) > 0 {
		ss = append(ss, d.Suffix)
	}

	return strings.Join(ss, "")
}
