package typescript

func Group(children ...*Statement) *Statement {
	return newStatementWithChildren(children)
}

func Interface(name string, properties ...*Statement) *Statement {
	b := Block(properties...)

	b.separator = ";\n  "

	if len(properties) > 0 {
		properties[len(properties)-1].suffix += ";"
	}

	s := newStatementWithChildren([]*Statement{b})
	s.prefix = "interface "
	s.value = name
	return s
}

func Block(content ...*Statement) *Statement {
	s := newStatementWithChildren(content)
	s.prefix = " {\n  "
	s.suffix = "\n}\n"
	s.separator = "\n  "
	return s
}

func Property(name string, propType string) *Statement {
	c := Colon()
	c.suffix = " "
	return Group(
		Literal(name),
		c,
		Literal(propType),
	)
}

func OptionalProperty(name string, propType string) *Statement {
	c := Colon()
	c.suffix = " "
	return Group(
		Literal(name),
		QuestionMark(),
		c,
		Literal(propType),
	)
}

func Colon() *Statement {
	return newStatementWithValue(":")
}

func ReturnType(name string) *Statement {
	s := newStatementWithValue(name)
	s.prefix = ": "
	return s
}

func Dot(propName string) *Statement {
	s := newStatementWithValue(propName)
	s.prefix = "."
	return s
}

func Call(params ...*Statement) *Statement {
	s := newStatementWithChildren(params)
	s.prefix = "("
	s.suffix = ")"
	return s
}

func Async() *Statement {
	return newStatementWithValue("async ")
}

func Await() *Statement {
	return newStatementWithValue("await ")
}

func Function(name string) *Statement {
	s := newStatementWithValue(name)
	s.prefix = "function "
	return s
}

func Params(params ...*Statement) *Statement {
	s := newStatementWithChildren(params)
	s.separator = ", "
	s.prefix = "("
	s.suffix = ")"
	return s
}

func Param(name, paramType string) *Statement {
	return Property(name, paramType)
}

func Arrow() *Statement {
	return newStatementWithValue("=>")
}

func ArrowFunction(params []*Statement, content []*Statement) *Statement {
	return newStatementWithChildren([]*Statement{
		Params(params...),
		Arrow(),
		Block(content...),
	})
}

func Class(name string) *Statement {
	s := newStatementWithValue(name)
	s.prefix = "class "
	return s
}

func Export() *Statement {
	return newStatementWithValue("export ")
}

func Import() *Statement {
	return newStatementWithValue("import ")
}

func As(typeName string) *Statement {
	s := newStatementWithValue(typeName)
	s.prefix = "as "
	return s
}

func Namespace(name string) *Statement {
	s := newStatementWithValue(name)
	s.prefix = "namespace "
	return s
}

func Literal(value string) *Statement {
	return newStatementWithValue(value)
}

func Return(value *Statement) *Statement {
	s := newStatementWithChildren([]*Statement{value})
	s.prefix = "return "
	s.suffix = ";"
	return s
}

func QuestionMark() *Statement {
	return Literal("?")
}

func Const(prop *Statement, value *Statement) *Statement {
	g := Group(prop, value)
	g.prefix = "const "
	g.separator = " = "
	g.suffix = ";\n"
	return g
}

func Type(name string, value *Statement) *Statement {
	g := Group(
		Literal(name),
		value,
	)
	g.prefix = "type "
	g.separator = " = "
	g.suffix = ";\n"
	return g
}
