package main

import (
	"fmt"

	"github.com/pepa65/rename/src"
)

func main() {
	args := rename.ParseArgs()
	err := rename.Run(args)
	if err != nil {
		fmt.Println(err)
	}
}
