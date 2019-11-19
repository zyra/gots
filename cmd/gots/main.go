package main

import (
	"github.com/urfave/cli/v2"
	"github.com/zyra/gots/parser"
	"log"
	"os"
)

var AppVersion = "0.0.1"

func main() {
	app := cli.NewApp()
	app.Name = "GoTS"
	app.Description = "Export Go definitions to TypeScript"
	app.Version = AppVersion

	wd, _ := os.Getwd()

	app.Commands = []*cli.Command{
		{
			Name:    "export",
			Aliases: []string{"e"},
			Usage:   "Export definitions",
			Action: func(ctx *cli.Context) error {
				o := ctx.String("outfile")

				p := parser.New(&parser.Config{
					BaseDir:     ctx.String("dir"),
					OutFileName: o,
					Subdirs:     ctx.Bool("sub-dir"),
				})

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
					Name:    "outfile",
					Aliases: []string{"o"},
					Usage:   "Output file. If not specified, stdout will be used.",
				},
				&cli.BoolFlag{
					Name:    "sub-dirs",
					Aliases: []string{"s"},
					Usage:   "Scan sub-directories as well",
					Value:   false,
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
