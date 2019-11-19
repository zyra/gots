package typescript

import (
	"bytes"
)

type Writer struct {
	buff       *bytes.Buffer
	statements []*Statement
}

func NewWriter() *Writer {
	return &Writer{
		buff:       bytes.NewBuffer(make([]byte, 0)),
		statements: make([]*Statement, 0),
	}
}

func (w *Writer) MarshalBinary() ([]byte, error) {
	return w.buff.Bytes(), nil
}

func (w *Writer) UnmarshalBinary(data []byte) error {
	w.buff.Reset()
	_, err := w.buff.Write(data)
	return err
}

func (w *Writer) AppendStatement(s *Statement) *Writer {
	w.statements = append(w.statements, s)
	return w
}

func (w *Writer) AppendStatements(s ...*Statement) *Writer {
	w.statements = append(w.statements, s...)
	return w
}

func (w *Writer) String() string {
	return Group(w.statements...).String()
}
