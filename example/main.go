package main

import (
	"fmt"
	"os"

	"github.com/valord577/clix/example/cmd"
)

// @author valor.

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
}
