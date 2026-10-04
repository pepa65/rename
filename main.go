package main

import (
	"fmt"

	"github.com/pepa65/rename/src"
)

const version = "0.2.6"

func main() {
	args := rename.ParseArgs()
	err := rename.Run(args)
	if err != nil {
		fmt.Println(err)
	}
}
