package main

import (
	"os"

	"github.com/colinjfw/jbld/pkg/run"
)

func main() {
	err := run.Run(os.Args[1])
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
