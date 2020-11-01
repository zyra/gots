package parser

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zyra/gots/pkg/parser/tag"
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
			RootDir:     filepath.Join(wd, "test_fixture"),
			OutFileName: "test_result.ts",
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

func TestParseTagsOld(t *testing.T) {
	a := assert.New(t)
	var tagVal string
	var pt *tag.Tag
	var err error

	tagVal = "`json:\"name\""
	pt, err = tag.ParseTag(tagVal)

	if !a.NoError(err) {
		return
	}

	a.Equal("name", pt.Name)
	a.Equal("", pt.Type)
	a.Equal(false, pt.Optional)

	tagVal = "`json:\"name,omitempty\""
	pt, err = tag.ParseTag(tagVal)

	if !a.NoError(err) {
		return
	}

	a.Equal("name", pt.Name)
	a.Equal("", pt.Type)
	a.Equal(true, pt.Optional)

	tagVal = "`json:\"name\" gots:\"type:string\"`"
	pt, err = tag.ParseTag(tagVal)

	if !a.NoError(err) {
		return
	}

	a.Equal("name", pt.Name)
	a.Equal("string", pt.Type)
	a.Equal(false, pt.Optional)

	tagVal = "`json:\"name\" gots:\"name:nom,type:string\"`"
	pt, err = tag.ParseTag(tagVal)

	if !a.NoError(err) {
		return
	}

	a.Equal("nom", pt.Name)
	a.Equal("string", pt.Type)
	a.Equal(false, pt.Optional)

	tagVal = "`json:\"name\" gots:\"name:nom,type:string,optional\"`"
	pt, err = tag.ParseTag(tagVal)

	if !a.NoError(err) {
		return
	}

	a.Equal("nom", pt.Name)
	a.Equal("string", pt.Type)
	a.Equal(true, pt.Optional)

	tagVal = "`json:\"-\" gots:\"type:string\"`"
	pt, err = tag.ParseTag(tagVal)

	if !a.Error(err) {
		return
	}

	a.Equal(tag.ErrJsonIgnored, err)

	tagVal = "`bson:\"id\" gots:\"type:string\"`"
	pt, err = tag.ParseTag(tagVal)

	if !a.Error(err) {
		return
	}

	a.Equal(tag.ErrJsonTagNotPresent, err)
}

func TestParser(t *testing.T) {
	suite.Run(t, new(parserTestSuite))
}
