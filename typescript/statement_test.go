package typescript

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type testSuite struct {
	suite.Suite
}

func (t *testSuite) TestInterface() {
	s := Interface("ITest",
		Property("name", "string"),
	)

	e := `interface ITest {
  name: string;
}
`
	t.EqualValues(e, s.String())
}

func (t *testSuite) TestBlock() {
	s := Block()
	t.EqualValues(" {\n  \n}\n", s.String())
}

func (t *testSuite) TestParam() {
	s := Param("name", "string")
	t.Equal("name: string", s.String())
}

func (t *testSuite) TestParams() {
	s := Params(
		Param("name", "string"),
		Param("age", "number"),
	)

	t.Equal("(name: string, age: number)", s.String())
}

func (t *testSuite) TestFunction() {
	s := Function("test")
	t.Equal("function test", s.String())

	e := `function test(a: number, b: number): number {
  return a + b;
}
`

	s = Group(
		s,
		Params(
			Param("a", "number"),
			Param("b", "number"),
		),
		ReturnType("number"),
		Block(
			Return(
				Group(
					Literal("a + b"),
				),
			),
		),
	)

	t.Equal(e, s.String())
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}
