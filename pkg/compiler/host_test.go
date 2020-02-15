package compiler

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHost(t *testing.T) {
	h := &host{
		js:     "../../js/host.js",
		config: "./testdata/config.jsbld.js",
	}
	defer h.Close()

	expected := []Import{{Kind: "static", Name: "file2", Resolved: "file2.js"}}
	imports, err := h.Run(Source{
		Src:     "testdata/src/file.js",
		Dst:     "testdata/lib/file.js",
		Plugins: []string{"test"},
	})
	require.NoError(t, err)
	require.Equal(t, expected, imports)
}
