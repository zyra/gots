package parser

import (
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"testing"
)

type parserTestSuite struct {
	suite.Suite
	config *Config
	parser *Parser
}

func (t *parserTestSuite) SetupSuite() {
	if wd, err := os.Getwd(); err != nil {
		t.FailNow("unable to get working directory")
	} else {
		t.config = &Config{
			RootDir: filepath.Join(wd, "../../example/pkg"),
			Types:   TypesConfig{},
			Output: Output{
				Mode:        AIO,
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
