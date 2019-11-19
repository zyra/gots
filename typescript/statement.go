package typescript

import "strings"

type Statement struct {
	prefix    string
	suffix    string
	separator string
	value     string
	children  []*Statement
}

func newStatement() *Statement {
	return &Statement{
		children: make([]*Statement, 0),
	}
}

func newStatementWithChildren(children []*Statement) *Statement {
	return &Statement{
		children: children,
	}
}

func newStatementWithValue(val string) *Statement {
	s := newStatement()
	s.value = val
	return s
}

func (d *Statement) append(c *Statement) *Statement {
	d.children = append(d.children, c)
	return d
}

func (d *Statement) String() string {
	ss := make([]string, 0, 4)

	if d.prefix != "" {
		ss = append(ss, d.prefix)
	}

	if d.value != "" {
		ss = append(ss, d.value)
	}

	if d.children != nil {
		cl := len(d.children)
		cs := make([]string, cl, cl)
		for i := range d.children {
			cs[i] = d.children[i].String()
		}
		ss = append(ss, strings.Join(cs, d.separator))
	}

	if d.suffix != "" {
		ss = append(ss, d.suffix)
	}

	return strings.Join(ss, "")
}
