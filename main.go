package main

import (
	"github.com/boyter/lc/parsers"
	// "github.com/pkg/profile"
	"github.com/urfave/cli"
	"os"
)

//go:generate go run scripts/include.go
func main() {
	// defer profile.Start().Stop()

	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = parsers.ToolName
	app.Version = parsers.ToolVersion
	app.Usage = "Check directory for licenses and list what license(s) a file is under"
	app.UsageText = "lc [global options] [DIRECTORY]"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "format, f",
			Usage:       "Set output format, supports progress, tabular, json, spdx or `csv`",
			Destination: &parsers.Format,
			Value:       "progress",
		},
		cli.StringFlag{
			Name:        "output, o",
			Usage:       "Set output file if not set will print to stdout `FILE`",
			Destination: &parsers.FileOutput,
		},
		cli.StringFlag{
			Name:        "confidence, c",
			Usage:       "Set required confidence level for licence matching between 0 and 1 `0.85`",
			Value:       "0.85",
			Destination: &parsers.Confidence,
		},
		cli.StringFlag{
			Name:        "deepguess, dg",
			Usage:       "Should attempt to deep guess the licence false or true `true`",
			Value:       "true",
			Destination: &parsers.DeepGuess,
		},
		cli.StringFlag{
			Name:        "filesize, fs",
			Usage:       "How large a file in bytes should be processed `50000`",
			Value:       "50000",
			Destination: &parsers.MaxSize,
		},
		cli.StringFlag{
			Name:        "licensefiles, lf",
			Usage:       "Possible license files to inspect for over-arching license as comma seperated list `copying,readme`",
			Value:       "license,copying,readme",
			Destination: &parsers.PossibleLicenceFiles,
		},
		cli.StringFlag{
			Name:        "pathblacklist, pbl",
			Usage:       "Which directories should be ignored as comma seperated list `.git,.hg,.svn`",
			Value:       ".git,.hg,.svn",
			Destination: &parsers.PathBlacklist,
		},
		cli.StringFlag{
			Name:        "extblacklist, xbl",
			Usage:       "Which file extensions should be ignored as comma seperated list `gif,jpg,png`",
			Value:       "woff,eot,cur,dm,xpm,emz,db,scc,idx,mpp,dot,pspimage,stl,dml,wmf,rvm,resources,tlb,docx,doc,xls,xlsx,ppt,pptx,msg,vsd,chm,fm,book,dgn,blines,cab,lib,obj,jar,pdb,dll,bin,out,elf,so,msi,nupkg,pyc,ttf,woff2,jpg,jpeg,png,gif,bmp,psd,tif,tiff,yuv,ico,xls,xlsx,pdb,pdf,apk,com,exe,bz2,7z,tgz,rar,gz,zip,zipx,tar,rpm,bin,dmg,iso,vcd,mp3,flac,wma,wav,mid,m4a,3gp,flv,mov,mp4,mpg,rm,wmv,avi,m4v,sqlite,class,rlib,ncb,suo,opt,o,os,pch,pbm,pnm,ppm,pyd,pyo,raw,uyv,uyvy,xlsm,swf",
			Destination: &parsers.ExtentionBlacklist,
		},
		cli.StringFlag{
			Name:        "documentname, dn",
			Usage:       "Only used if you specify SPDX as an output `gif,jpg,png`",
			Value:       "Unknown",
			Destination: &parsers.DocumentName,
		},
		cli.StringFlag{
			Name:        "packagename, pn",
			Usage:       "Only used if you specify SPDX as an output should be the",
			Value:       "Unknown",
			Destination: &parsers.PackageName,
		},
	}
	app.Action = func(c *cli.Context) error {
		parsers.DirPath = c.Args().Get(0)
		parsers.Process()
		return nil
	}

	app.Run(os.Args)
}
