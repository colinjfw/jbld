package bundler

import (
	"os"
	"testing"

	"github.com/colinjfw/jbld/pkg/compiler"
	"github.com/stretchr/testify/require"
)

func TestBundler(t *testing.T) {
	cwd, _ := os.Getwd()
	os.RemoveAll("./testdata/out2")

	c := &compiler.Compiler{Config: compiler.Config{
		HostJS:      cwd + "/../../lib/host.js",
		ConfigFile:  cwd + "/testdata/config.jbld.js",
		SourceDir:   cwd + "/testdata/src",
		OutputDir:   cwd + "/testdata/out/target",
		Entrypoints: []string{"index.js"},
		Workers:     1,
	}}

	m, err := c.Run()
	require.NoError(t, err)

	b := &Bundler{
		Config: Config{
			OutputDir: cwd + "/testdata/out/bundle",
			AssetPath: "static",
			BaseURL:   "/",
		},
		Manifest: m,
	}
	err = b.Run()
	require.NoError(t, err)
}
