package parser

import (
	"fmt"
	"github.com/zyra/gots/pkg/parser/reader"
	"github.com/zyra/gots/pkg/statement"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Parser struct {
	*reader.Config
	pkgIndex map[string]*reader.Package
	pkgTsw   map[string]*statement.Writer
}

func New(config *reader.Config) *Parser {
	if !filepath.IsAbs(config.RootDir) {
		if d, err := filepath.Abs(config.RootDir); err != nil {
			log.Panicf("cannot convert base directory to absolute path: %s\n", err.Error())
		} else {
			config.RootDir = d
		}
	}

	return &Parser{
		Config:   config,
		pkgIndex: make(map[string]*reader.Package),
		pkgTsw:   make(map[string]*statement.Writer),
	}
}

func (p *Parser) String() string {
	strs := make([]string, len(p.pkgTsw))
	for _, tsw := range p.pkgTsw {
		strs = append(strs, tsw.String())
	}
	return strings.Join(strs, "\n")
}

func (p *Parser) WriteToFile() error {
	switch p.Config.Output.Mode {
	case reader.AIO:
		outPath := p.Output.AIOFileName
		if !filepath.IsAbs(outPath) {
			wd, err := os.Getwd()
			if err == nil {
				outPath = filepath.Join(wd, outPath)
			}
		}
		dirPath := filepath.Dir(outPath)
		_ = os.MkdirAll(dirPath, 0744)

		strs := make([]string, len(p.pkgTsw))
		for _, tsw := range p.pkgTsw {
			strs = append(strs, tsw.String())
		}
		return ioutil.WriteFile(p.Output.AIOFileName, []byte(strings.Join(strs, "\n")), 0644)
	case reader.Packages:
		for n, tsw := range p.pkgTsw {
			outPath := filepath.Join(p.Output.Path, fmt.Sprintf("%s.ts", n))
			if !filepath.IsAbs(outPath) {
				wd, err := os.Getwd()
				if err == nil {
					outPath = filepath.Join(wd, outPath)
				}
			}
			dirPath := filepath.Dir(outPath)
			_ = os.MkdirAll(dirPath, 0744)

			if err := ioutil.WriteFile(outPath, []byte(tsw.String()), 0644); err != nil {
				return err
			}
		}
		return nil
	case reader.Mirror:
		for n, pkg := range p.pkgIndex {
			relPath := strings.Replace(pkg.Path, p.RootDir, "", 1)
			outPath := filepath.Join(p.Output.Path, relPath, fmt.Sprintf("%s.ts", n))
			if !filepath.IsAbs(outPath) {
				wd, err := os.Getwd()
				if err == nil {
					outPath = filepath.Join(wd, outPath)
				}
			}
			dirPath := filepath.Dir(outPath)
			_ = os.MkdirAll(dirPath, 0744)
			if err := ioutil.WriteFile(outPath, []byte(p.pkgTsw[n].String()), 0644); err != nil {
				return err
			}
		}
	}

	return fmt.Errorf("invalid output mode %s", p.Config.Output.Mode)
}

func (p *Parser) Print() {
	fmt.Print(p.String())
}
