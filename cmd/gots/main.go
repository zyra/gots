package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/zyra/gots/pkg/parser"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path"
)

var AppVersion = "0.0.1"

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
				o := ctx.String("outfile")

				var config parser.Config

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
					config.Types = parser.TypesConfig{
						Interfaces: ctx.Bool("interfaces"),
						Constants:  ctx.Bool("constants"),
						Aliases:    ctx.Bool("type-aliases"),
						Structs:    ctx.Bool("structs"),
						Enums:      ctx.Bool("enums"),
					}
					config.Output = parser.Output{
						Mode:        parser.OutputMode(ctx.String("mode")),
						Path:        ctx.String("output-path"),
						AIOFileName: ctx.String("aio-file"),
					}
				}

				p := parser.New(&config)

				p.Run()
				p.GenerateTS()

				if o == "" {
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
				&cli.BoolFlag{
					Name:    "interfaces",
					EnvVars: []string{"GOTS_PARSE_INTERFACES"},
					Value:   true,
					Usage:   "Parse interfaces and use them in generated code where applicable",
				},
				&cli.BoolFlag{
					Name:    "structs",
					EnvVars: []string{"GOTS_PARSE_STRUCTS"},
					Value:   true,
					Usage:   "Parse structs and export them as TypeScript interfaces",
				},
				&cli.BoolFlag{
					Name:    "constants",
					Aliases: []string{"consts"},
					EnvVars: []string{"GOTS_PARSE_CONSTANTS"},
					Value:   true,
					Usage:   "Parse and export constants",
				},
				&cli.BoolFlag{
					Name:    "enums",
					Aliases: []string{"enums"},
					EnvVars: []string{"GOTS_PARSE_ENUMS"},
					Value:   true,
					Usage:   "Automatically detect enums and export them as TypeScript enums",
				},
				&cli.BoolFlag{
					Name:    "type-aliases",
					Aliases: []string{"aliases", "types"},
					EnvVars: []string{"GOTS_PARSE_TYPE_ALIASES"},
					Value:   true,
					Usage:   "Parse and export type aliases",
				},
				&cli.StringFlag{
					Name:    "mode",
					Aliases: []string{"m"},
					EnvVars: []string{"GOTS_OUTPUT_MODE"},
					Value:   string(parser.Packages),
					Usage:   "Output mode: aio, packages, or mirror",
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
