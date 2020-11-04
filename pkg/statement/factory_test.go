package statement

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type factorySuite struct {
	suite.Suite
}

func (a *factorySuite) TestGroup() {
	c := NewStatement()
	s := Group(c)
	if a.NotNil(c) && a.NotNil(s) && a.Len(s.Children, 1) {
		a.Same(c, s.Children[0])
	}
}

func (a *factorySuite) TestProperty() {
	s := Property("test", "string")
	a.Equal(`test: string`, s.String())
}

func (a *factorySuite) TestOptionalProperty() {
	s := OptionalProperty("test", "string")
	a.Equal("test?: string", s.String())
}

func (a *factorySuite) TestInterface() {
	s := Interface("Test",
		Property("name", "string"),
		Property("email", "string"),
	)
	a.Equal(`interface Test {
  name: string;
  email: string;
}
`, s.String())
}

func (a *factorySuite) TestDot() {
	a.Equal(".test", Dot("test").String())
	a.Equal(".", Dot("").String())
}

func (a *factorySuite) TestLiteralString() {
	a.Equal(`"test"`, LiteralString("test").String())
}

func (a *factorySuite) TestLiteral() {
	a.Equal(`"test"`, Literal(`"test"`).String())
	a.Equal(`true`, Literal(`true`).String())
	a.Equal(`1`, Literal(`1`).String())
}

func (a *factorySuite) TestEnumValue() {
	a.Equal(`Test = "test"`, EnumValue("Test", LiteralString("test")).String())
	a.Equal(`Test = 1`, EnumValue("Test", Literal("1")).String())
}

func (a *factorySuite) TestEnum() {
	s := Enum("Mode",
		EnumValue("Test", LiteralString("test")),
		EnumValue("Development", LiteralString("dev")),
		EnumValue("Production", LiteralString("prod")),
	)
	a.Equal(`enum Mode {
  Test = "test",
  Development = "dev",
  Production = "prod"
}
`, s.String())

	a.Equal(`enum Mode {
  Test = "test"
}
`, Enum("Mode", EnumValue("Test", LiteralString("test"))).String())
}

func TestFactory(t *testing.T) {
	suite.Run(t, new(factorySuite))
}
