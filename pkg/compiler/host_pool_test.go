package compiler

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHostPool(t *testing.T) {
	h := &hostPool{
		count:  3,
		js:     "../../js/host.js",
		config: "./testdata/config.jsbld.js",
	}
	h.create()
	defer h.Close()

	expected := []Import{{Kind: "static", Name: "file2", Resolved: "file2.js"}}
	wg := sync.WaitGroup{}
	wg.Add(5)

	for i := 0; i < 5; i++ {
		go func() {
			imports, err := h.Run(Source{
				Src:     "testdata/src/file.js",
				Dst:     "testdata/lib/file.js",
				Plugins: []string{"test"},
			})
			require.NoError(t, err)
			require.Equal(t, expected, imports)
			wg.Done()
		}()
	}
	wg.Wait()
}
