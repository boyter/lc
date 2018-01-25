package main

import (
	"github.com/boyter/license-checker/parsers"
	"github.com/briandowns/spinner"
	"github.com/urfave/cli"
	"os"
	"time"
)

//go:generate go run scripts/include.go
func main() {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = "golang-license-checker"
	app.Version = "1.0"
	app.Usage = "Check directory for licenses and list what license(s) a file is under"
	app.UsageText = "golang-licence-checker [global options] DIRECTORY"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "format, f",
			Usage:       "Set output format, supports cli, json or `csv`",
			Destination: &parsers.Format,
			Value:       "cli",
		},
		cli.StringFlag{
			Name:        "output, o",
			Usage:       "Set output file `FILE`",
			Destination: &parsers.FileOutput,
			Value:       "./output",
		},
		cli.StringFlag{
			Name:        "confidence, c",
			Usage:       "Set required confidence level for licence matching should be number between 0 and 1 `0.85`",
			Value:       "0.85",
			Destination: &parsers.Confidence,
		},
		cli.StringFlag{
			Name:  "deepguess, dg",
			Usage: "Should attempt to deep guess the licence false or true `true`",
			Value: "true",
		},
	}
	app.Action = func(c *cli.Context) error {
		parsers.DirPath = c.Args().Get(0)
		parsers.Process()
		return nil
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Writer = os.Stderr
	s.Start()
	s.Suffix = " processing"
	app.Run(os.Args)
	s.Stop()

}
