package parser

import (
	"github.com/stretchr/testify/suite"
	"github.com/zyra/gots/pkg/parser/reader"
	"os"
	"path/filepath"
	"testing"
)

type parserTestSuite struct {
	suite.Suite
	config *reader.Config
	parser *Parser
}

func (t *parserTestSuite) SetupSuite() {
	if wd, err := os.Getwd(); err != nil {
		t.FailNow("unable to get working directory")
	} else {
		t.config = &reader.Config{
			RootDir: filepath.Join(wd, "../../example/pkg"),
			Types:   reader.TypesConfig{},
			Output: reader.Output{
				Mode:        reader.AIO,
				AIOFileName: "test_results.ts",
			},
			Recursive:  true,
			Transforms: nil,
			Include:    nil,
			Exclude:    nil,
		}
	}
}

func (t *parserTestSuite) SetupTest() {
	t.parser = New(t.config)
}

func (t *parserTestSuite) TestRun() {
	t.parser.Run()
	t.parser.GenerateTS()
	out := t.parser.String()
	t.EqualValues(`export interface TestModel {
  id: string;
  name?: string;
  yearsAlive?: number;
}
`, out)
}

func TestParser(t *testing.T) {
	suite.Run(t, new(parserTestSuite))
}
