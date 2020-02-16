package compiler

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompiler(t *testing.T) {
	cwd, _ := os.Getwd()
	os.RemoveAll("./testdata/lib")

	c := &Compiler{Config: Config{
		HostJS:      cwd + "/../../js/host.js",
		ConfigFile:  cwd + "/testdata/config.jsbld.js",
		SourceDir:   cwd + "/testdata/src",
		OutputDir:   cwd + "/testdata/lib",
		Entrypoints: []string{"file.js"},
		Plugins:     []string{"test"},
		Workers:     1,
	}}

	t.Run("Normal", func(t *testing.T) {
		require.NoError(t, c.Run())
	})
	t.Run("Cached", func(t *testing.T) {
		require.NoError(t, c.Run())
	})
}
