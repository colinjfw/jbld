package example

import (
	"os"
	"testing"

	"github.com/colinjfw/jbld/pkg/bundler"
	"github.com/colinjfw/jbld/pkg/compiler"
	"github.com/stretchr/testify/require"
)

func TestCompiler(t *testing.T) {
	os.RemoveAll("./lib")
	os.RemoveAll("./dist")
	cwd, _ := os.Getwd()

	c := &compiler.Compiler{Config: compiler.Config{
		HostJS:      cwd + "/../js/host.js",
		ConfigFile:  cwd + "/config.jbld.js",
		SourceDir:   cwd,
		OutputDir:   cwd + "/lib",
		Entrypoints: []string{"src/index.js"},
		Plugins:     []string{"babel"},
		Workers:     2,
	}}
	b := &bundler.Bundler{
		OutputDir:    "./dist",
		RuntimeJS:    "../js/runtime.js",
		ManifestPath: "./lib/.jbld-manifest",
	}

	t.Run("Normal", func(t *testing.T) {
		_, err := c.Run()
		require.NoError(t, err)
	})
	t.Run("Cached", func(t *testing.T) {
		_, err := c.Run()
		require.NoError(t, err)
	})
	t.Run("Bundler", func(t *testing.T) {
		require.NoError(t, b.Run())
	})
}
