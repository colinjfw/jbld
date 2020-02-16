package compiler

import (
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHostPool(t *testing.T) {
	cwd, _ := os.Getwd()
	h := NewHostPool(Config{
		Workers:    3,
		HostJS:     cwd + "/../../lib/host.js",
		ConfigFile: cwd + "/testdata/config.jsbld.js",
		SourceDir:  cwd + "/testdata/src",
		OutputDir:  cwd + "/testdata/lib",
	})
	defer h.Close()

	expected := []Import{{Kind: "static", Name: "file2", Resolved: "file2.js"}}
	wg := sync.WaitGroup{}
	wg.Add(5)

	for i := 0; i < 5; i++ {
		go func() {
			imports, err := h.Run(Source{
				Name: "file.js",
			})
			require.NoError(t, err)
			require.Equal(t, expected, imports)
			wg.Done()
		}()
	}
	wg.Wait()
}
