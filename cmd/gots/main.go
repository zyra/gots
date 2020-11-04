package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/zyra/gots/pkg/parser"
	"github.com/zyra/gots/pkg/parser/reader"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path"
)

var AppVersion = "0.2.0"

func main() {
	app := cli.NewApp()
	app.Name = "GoTS"
	app.Description = "Parses Go types and generates TypeScript code"
	app.Version = AppVersion

	wd, _ := os.Getwd()

	app.Commands = []*cli.Command{
		{
			Name:    "export",
			Aliases: []string{"e"},
			Usage:   "Export definitions",
			Action: func(ctx *cli.Context) error {
				var config reader.Config

				configPath := ctx.String("config")
				if len(configPath) > 0 {
					cf, err := ioutil.ReadFile(configPath)
					if err != nil {
						return fmt.Errorf("failed to read config from path %v: %v", configPath, err)
					}
					if err := yaml.Unmarshal(cf, &config); err != nil {
						return fmt.Errorf("failed to parse config file: %v", err)
					}
				} else {
					config.RootDir = ctx.String("dir")
					config.Recursive = ctx.Bool("recursive")
					config.Output = reader.Output{
						Mode:        reader.OutputMode(ctx.String("mode")),
						Path:        ctx.String("output-path"),
						AIOFileName: ctx.String("aio-file"),
					}
				}

				p := parser.New(&config)

				if err := p.Run(); err != nil {
					return err
				}

				if ctx.Bool("stdout") {
					p.Print()
				} else {
					return p.WriteToFile()
				}

				return nil
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "dir",
					Aliases: []string{"d"},
					Usage:   "Base directory to lookup exportable definitions",
					Value:   wd,
				},
				&cli.StringFlag{
					Name:    "output-path",
					Aliases: []string{"output", "o"},
					EnvVars: []string{"GOTS_OUTPUT_PATH"},
					Usage:   "Output directory path",
					Value:   path.Join(wd, "gots"),
				},
				&cli.StringFlag{
					Name:    "aio-file",
					Aliases: []string{"af", "f"},
					Usage:   "Filename to export generated code to. Used only if output mode is aio",
					EnvVars: []string{"GOTS_AIO_FILE"},
					Value:   "models.ts",
				},
				&cli.BoolFlag{
					Name:    "recursive",
					Aliases: []string{"R"},
					Usage:   "Scan directories recursively",
					Value:   false,
				},
				&cli.StringFlag{
					Name:      "config",
					Aliases:   []string{"c"},
					Usage:     "Config file",
					EnvVars:   []string{"GOTS_CONFIG"},
					TakesFile: true,
				},
				&cli.StringFlag{
					Name:    "mode",
					Aliases: []string{"m"},
					EnvVars: []string{"GOTS_OUTPUT_MODE"},
					Value:   string(reader.Packages),
					Usage:   "Output mode: aio, packages, or mirror",
				},
				&cli.BoolFlag{
					Name: "stdout",
					Aliases: []string{"O"},
					Usage: "Print to stdout instead of file",
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
