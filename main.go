// SPDX-License-Identifier: MIT
// SPDX-License-Identifier: Unlicense
package main

import (
	"github.com/boyter/lc/processor"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	rootCmd := &cobra.Command{
		Use: "lc",
		Long: "Check directory/file for licenses and list what license(s) a file is under.\n" +
			"Version " + processor.Version + "\n" +
			"Ben Boyter <ben@boyter.org>",
		Version: processor.Version,
		Run: func(cmd *cobra.Command, args []string) {
			// process hher
		},
	}

	flags := rootCmd.PersistentFlags()

	flags.BoolVar(
		&processor.IncludeBinaryFiles,
		"binary",
		false,
		"set to disable binary file detection",
	)
	flags.BoolVar(
		&processor.IgnoreIgnoreFile,
		"no-ignore",
		false,
		"disables .ignore file logic",
	)
	flags.BoolVar(
		&processor.IgnoreGitIgnore,
		"no-gitignore",
		false,
		"disables .gitignore file logic",
	)
	flags.BoolVar(
		&processor.IncludeHidden,
		"hidden",
		false,
		"include hidden files",
	)
	flags.StringSliceVarP(
		&processor.AllowListExtensions,
		"include-ext",
		"i",
		[]string{},
		"limit to file extensions case sensitive [comma separated list: e.g. go,java,js,C,cpp]",
	)
	flags.StringVarP(
		&processor.Format,
		"format",
		"f",
		"tabular",
		"set output format [progress, tabular, json, spdx, xlsx, csv]",
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
