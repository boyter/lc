package main

import (
	"github.com/boyter/golang-license-checker/parsers"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = "golang-license-checker"
	app.Version = "1.0"
	app.Usage = "Check directory for licenses and list what license(s) a file is under"

	app.Commands = []cli.Command{
		{
			Name:    "process",
			Aliases: []string{"cf"},
			Usage:   "tasks for building and deploying cloudformation templates",
			Subcommands: []cli.Command{
				{
					Name:      "generate",
					Usage:     "compile a new cf template from a config",
					UsageText: "gostacks cloudformation generate [command options] [stack]",
					Action:    parsers.Process,
					Flags:     parsers.Generate_Flags,
				},
			},
		},
	}

	app.Run(os.Args)

	// Everything after here needs to be refactored out to a subpackage
	// parsers.Process(nil)
}
