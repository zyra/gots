package statement

// Create a statement with the provided statements as it's children
func Group(children ...*Statement) *Statement {
	return NewStatementWithChildren(children)
}

// Create a typescript interface with the provided name and properties
func Interface(name string, properties ...*Statement) *Statement {
	b := Block(properties...)

	b.Separator = ";\n  "

	if len(properties) > 0 {
		properties[len(properties)-1].Suffix += ";"
	}

	s := NewStatementWithChildren([]*Statement{b})
	s.Prefix = "interface "
	s.Value = name
	return s
}

// Create a block statement with the provided statements joined by a new line
func Block(content ...*Statement) *Statement {
	s := NewStatementWithChildren(content)
	s.Prefix = " {\n  "
	s.Suffix = "\n}\n"
	s.Separator = "\n  "
	return s
}

// Create a typescript property with the provided name and type
func Property(name string, propType string) *Statement {
	c := Colon()
	c.Suffix = " "
	return Group(
		Literal(name),
		c,
		Literal(propType),
	)
}

func OptionalProperty(name string, propType string) *Statement {
	c := Colon()
	c.Suffix = " "
	return Group(
		Literal(name),
		QuestionMark(),
		c,
		Literal(propType),
	)
}

func Colon() *Statement {
	return NewStatementWithValue(":")
}

func ReturnType(name string) *Statement {
	s := NewStatementWithValue(name)
	s.Prefix = ": "
	return s
}

func Dot(propName string) *Statement {
	s := NewStatementWithValue(propName)
	s.Prefix = "."
	return s
}

func Call(params ...*Statement) *Statement {
	s := NewStatementWithChildren(params)
	s.Prefix = "("
	s.Suffix = ")"
	return s
}

func Async() *Statement {
	return NewStatementWithValue("async ")
}

func Await() *Statement {
	return NewStatementWithValue("await ")
}

func Function(name string) *Statement {
	s := NewStatementWithValue(name)
	s.Prefix = "function "
	return s
}

func Params(params ...*Statement) *Statement {
	s := NewStatementWithChildren(params)
	s.Separator = ", "
	s.Prefix = "("
	s.Suffix = ")"
	return s
}

func Param(name, paramType string) *Statement {
	return Property(name, paramType)
}

func Arrow() *Statement {
	return NewStatementWithValue("=>")
}

func ArrowFunction(params []*Statement, content []*Statement) *Statement {
	return NewStatementWithChildren([]*Statement{
		Params(params...),
		Arrow(),
		Block(content...),
	})
}

func Class(name string) *Statement {
	s := NewStatementWithValue(name)
	s.Prefix = "class "
	return s
}

func Export() *Statement {
	return NewStatementWithValue("export ")
}

func Import() *Statement {
	return NewStatementWithValue("import ")
}

func As(typeName string) *Statement {
	s := NewStatementWithValue(typeName)
	s.Prefix = "as "
	return s
}

func Namespace(name string) *Statement {
	s := NewStatementWithValue(name)
	s.Prefix = "namespace "
	return s
}

func Literal(value string) *Statement {
	return NewStatementWithValue(value)
}

func Return(value *Statement) *Statement {
	s := NewStatementWithChildren([]*Statement{value})
	s.Prefix = "return "
	s.Suffix = ";"
	return s
}

func QuestionMark() *Statement {
	return Literal("?")
}

func Const(prop *Statement, value *Statement) *Statement {
	g := Group(prop, value)
	g.Prefix = "const "
	g.Separator = " = "
	g.Suffix = ";\n"
	return g
}

func Type(name string, value *Statement) *Statement {
	g := Group(
		Literal(name),
		value,
	)
	g.Prefix = "type "
	g.Separator = " = "
	g.Suffix = ";\n"
	return g
}
