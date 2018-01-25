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
	app.Action = func(c *cli.Context) error {
		return nil
	}

	app.Run(os.Args)

	// Everything after here needs to be refactored out to a subpackage
	parsers.WalkDirectory("", []parsers.LicenseMatch{})
}
