package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"

	"github.com/colinjfw/jbld/pkg/bundler"
	"github.com/colinjfw/jbld/pkg/compiler"
)

func main() {
	var (
		jsRoot      = os.Getenv("JS_ROOT")
		output      string
		configFile  string
		entrypoints string
	)
	flag.StringVar(&output, "out", "./dist", "output directory")
	flag.StringVar(&configFile, "config", "./config.jbld.js", "config directory")
	flag.StringVar(&entrypoints, "entrypoints", "index.js", "app entrypoints")

	bundleConfig := bundler.Config{
		OutputDir: filepath.Join(abs(output), "bundle"),
	}
	compileConfig := compiler.Config{
		HostJS:      filepath.Join(jsRoot, "host.js"),
		OutputDir:   filepath.Join(abs(output), "dist"),
		ConfigFile:  abs(configFile),
		Entrypoints: strings.Split(entrypoints, " "),
	}
}

func abs(path string) string {
	output, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return output
}
