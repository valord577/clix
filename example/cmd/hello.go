package cmd

import (
	"fmt"
	"os"

	"github.com/valord577/clix"
)

// @author valor.

var helloCmd = &clix.Command{
	Name: "hello",

	Summary:     "print hello message.",
	Description: "You can write a long description.",

	Run: printHello,
}

var (
	showDetail bool
)

func init() {
	helloCmd.FlagBoolVar(&showDetail, "d", false, "show more messages")
}

func printHello(c *clix.Command, args []string) error {
	fmt.Fprintf(os.Stderr, "hello world.\n")
	if showDetail {
		fmt.Fprintf(os.Stderr, "more messages: cmd: %s | args: %v\n", c.Name, args)
	}
	return nil
}
