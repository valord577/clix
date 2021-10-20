package cmd

import "github.com/valord577/clix"

// @author valor.

var rootCmd = &clix.Command{
	Name: "example",

	Summary: "example for `github.com/valord577/cli`",
}

func init() {
	rootCmd.AddCmd(helloCmd)
}

// Execute executes root command
func Execute() error {
	return rootCmd.Execute()
}
