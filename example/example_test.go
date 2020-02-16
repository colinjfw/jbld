package example

import (
	"os"
	"testing"

	"github.com/colinjfw/jbld/pkg/compiler"
	"github.com/stretchr/testify/require"
)

func TestCompiler(t *testing.T) {
	os.RemoveAll("./lib")
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

	t.Run("Normal", func(t *testing.T) {
		require.NoError(t, c.Run())
	})
	t.Run("Cached", func(t *testing.T) {
		require.NoError(t, c.Run())
	})
}
