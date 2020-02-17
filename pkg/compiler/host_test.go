package compiler

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHost(t *testing.T) {
	cwd, _ := os.Getwd()
	h := NewHost(Config{
		HostJS:     cwd + "/../../lib/host.js",
		ConfigFile: cwd + "/testdata/config.jsbld.js",
		SourceDir:  cwd + "/testdata/src",
		OutputDir:  cwd + "/testdata/lib",
	})
	defer h.Close()

	expected := HostResponse{
		Type:    "js",
		Imports: []Import{{Kind: "static", Name: "file2", Resolved: "file2.js"}},
	}
	resp, err := h.Run(Source{
		Name: "file.js",
	})
	require.NoError(t, err)
	require.Equal(t, expected, resp)
}
