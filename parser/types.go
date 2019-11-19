package parser

type Property struct {
	Name, Type string
	Optional   bool
}

type Interface struct {
	Name       string
	Properties []*Property
}

type TypeDef struct {
	Name, Type string
}

type Constant struct {
	Name, Type, Value string
}
