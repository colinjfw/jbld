package main

import (
	"os"

	"github.com/colinjfw/jbld/pkg/run"
)

func main() {
	cwd, _ := os.Getwd()
	err := run.Run("../../lib/host.js", cwd+"/config.jbld.js")
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
