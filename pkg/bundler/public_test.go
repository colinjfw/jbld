package bundler

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPublic(t *testing.T) {
	os.RemoveAll("./testdata/out/public")

	m := &Manifest{
		Bundles: []string{"foo.js", "foo.css"},
	}

	err := writePublicFolder(Config{
		Public:    PublicConfig{Dir: "./testdata/public"},
		OutputDir: "./testdata/out/public",
	})
	require.NoError(t, err)

	err = writeHTMLSources(Config{
		Public:    PublicConfig{HTML: []string{"index.html"}},
		OutputDir: "./testdata/out/public",
	}, m)
	require.NoError(t, err)
}
