package statement

// TODO automatically generate

func (w *Writer) Interface(name string, properties ...*Statement) *Writer {
	return w.AppendStatement(Interface(name, properties...))
}

func (w *Writer) Block(content ...*Statement) *Writer {
	return w.AppendStatement(Block(content...))
}

func (w *Writer) Property(name string, propType string) *Writer {
	return w.AppendStatement(Property(name, propType))
}

func (w *Writer) Colon() *Writer {
	return w.AppendStatement(Colon())
}

func (w *Writer) Dot(propName string) *Writer {
	return w.AppendStatement(Dot(propName))
}

func (w *Writer) Call(params ...*Statement) *Writer {
	return w.AppendStatement(Call(params...))
}

func (w *Writer) Async() *Writer {
	return w.AppendStatement(Async())
}

func (w *Writer) Await() *Writer {
	return w.AppendStatement(Await())
}

func (w *Writer) Function(name string) *Writer {
	return w.AppendStatement(Function(name))
}

func (w *Writer) Params(params ...*Statement) *Writer {
	return w.AppendStatement(Params(params...))
}

func (w *Writer) Param(name, paramType string) *Writer {
	return w.AppendStatement(Param(name, paramType))
}

func (w *Writer) Arrow() *Writer {
	return w.AppendStatement(Arrow())
}

func (w *Writer) ArrowFunction(params []*Statement, content []*Statement) *Writer {
	return w.AppendStatement(ArrowFunction(params, content))
}

func (w *Writer) Class(name string) *Writer {
	return w.AppendStatement(Class(name))
}

func (w *Writer) Export() *Writer {
	return w.AppendStatement(Export())
}

func (w *Writer) Import() *Writer {
	return w.AppendStatement(Import())
}

func (w *Writer) As(typeName string) *Writer {
	return w.AppendStatement(As(typeName))
}

func (w *Writer) Namespace(name string) *Writer {
	return w.AppendStatement(Namespace(name))
}

func (w *Writer) Literal(value string) *Writer {
	return w.AppendStatement(Literal(value))
}

func (w *Writer) Return(value *Statement) *Writer {
	return w.AppendStatement(Return(value))
}

func (w *Writer) QuestionMark() *Writer {
	return w.AppendStatement(QuestionMark())
}

func (w *Writer) Const(prop *Statement, value *Statement) *Writer {
	return w.AppendStatement(Const(prop, value))
}

func (w *Writer) Type(name string, value *Statement) *Writer {
	return w.AppendStatement(Type(name, value))
}
