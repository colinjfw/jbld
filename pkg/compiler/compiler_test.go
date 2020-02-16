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
		HostJS:      cwd + "/../../lib/host.js",
		ConfigFile:  cwd + "/testdata/config.jsbld.js",
		SourceDir:   cwd + "/testdata/src",
		OutputDir:   cwd + "/testdata/lib",
		Entrypoints: []string{"file.js"},
		Workers:     1,
	}}

	t.Run("Normal", func(t *testing.T) {
		_, err := c.Run()
		require.NoError(t, err)
	})
	t.Run("Cached", func(t *testing.T) {
		_, err := c.Run()
		require.NoError(t, err)
	})
}
