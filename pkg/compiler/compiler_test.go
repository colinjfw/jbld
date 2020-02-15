package compiler

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompiler(t *testing.T) {
	os.RemoveAll("./testdata/lib")

	c := &Compiler{
		HostJS:      "../../js/host.js",
		ConfigFile:  "./testdata/config.jsbld.js",
		SourceDir:   "./testdata/src",
		OutputDir:   "./testdata/lib",
		Entrypoints: []string{"file.js"},
		Plugins:     []string{"test"},
		Workers:     1,
	}

	t.Run("Normal", func(t *testing.T) {
		require.NoError(t, c.Run())
	})
	t.Run("Cached", func(t *testing.T) {
		require.NoError(t, c.Run())
	})
}
