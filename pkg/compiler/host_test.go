package compiler

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHost(t *testing.T) {
	cwd, _ := os.Getwd()
	h := NewHost(Config{
		HostJS:     cwd + "/../../js/host.js",
		ConfigFile: cwd + "/testdata/config.jsbld.js",
		SourceDir:  cwd + "/testdata/src",
		OutputDir:  cwd + "/testdata/lib",
	})
	defer h.Close()

	expected := []Import{{Kind: "static", Name: "file2", Resolved: "file2.js"}}
	imports, err := h.Run(Source{
		Name:    "file.js",
		Plugins: []string{"test"},
	})
	require.NoError(t, err)
	require.Equal(t, expected, imports)
}
