package main

import (
	"github.com/boyter/golang-license-checker/parsers"
	"github.com/urfave/cli"
	"os"
)

//go:generate go run scripts/include.go
func main() {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = "golang-license-checker"
	app.Version = "1.0"
	app.Usage = "Check directory for licenses and list what license(s) a file is under"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "format, f",
			Usage: "Set output format, supports cli, json or `csv`",
		},
		cli.StringFlag{
			Name:        "confidence, c",
			Usage:       "Set required confidence level for licence matching defaults to `0.85`",
			Destination: &parsers.Confidence,
		},
		cli.StringFlag{
			Name:  "deepguess, dg",
			Usage: "Should attempt to deep guess the licence false or true defaults to `true`",
		},
	}
	app.Action = func(c *cli.Context) error {
		parsers.DirPath = c.Args().Get(0)
		parsers.Process()
		return nil
	}

	app.Run(os.Args)
}
