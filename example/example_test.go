package example

import (
	"os"
	"testing"

	"github.com/colinjfw/jbld/pkg/compiler"
	"github.com/stretchr/testify/require"
)

func TestCompiler(t *testing.T) {
	os.RemoveAll("./lib")

	c := &compiler.Compiler{
		HostJS:      "../js/host.js",
		ConfigFile:  "config.jbld.js",
		SourceDir:   "./src",
		OutputDir:   "./lib",
		Entrypoints: []string{"index.js"},
		Plugins:     []string{"babel"},
		Workers:     1,
	}

	t.Run("Normal", func(t *testing.T) {
		require.NoError(t, c.Run())
	})
	t.Run("Cached", func(t *testing.T) {
		require.NoError(t, c.Run())
	})
}
