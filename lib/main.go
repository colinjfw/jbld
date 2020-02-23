package main

import (
	"os"

	"github.com/colinjfw/jbld/pkg/run"
)

func main() {
	opts, err := run.LoadOptions(os.Args[1])
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	err = run.Run(opts)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
