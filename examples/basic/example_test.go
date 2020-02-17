package example

import (
	"os"
	"testing"

	"github.com/colinjfw/jbld/pkg/bundler"
	"github.com/colinjfw/jbld/pkg/compiler"
	"github.com/stretchr/testify/require"
)

func TestCompiler(t *testing.T) {
	os.RemoveAll("./lib")
	os.RemoveAll("./dist")
	cwd, _ := os.Getwd()

	c := &compiler.Compiler{Config: compiler.Config{
		HostJS:      cwd + "/../../lib/host.js",
		ConfigFile:  cwd + "/config.jbld.js",
		SourceDir:   cwd,
		OutputDir:   cwd + "/lib",
		Entrypoints: []string{"src/index.js"},
		Workers:     2,
	}}
	var err error
	var m *compiler.Manifest

	t.Run("Normal", func(t *testing.T) {
		m, err = c.Run()
		require.NoError(t, err)
	})
	t.Run("Cached", func(t *testing.T) {
		m, err = c.Run()
		require.NoError(t, err)
	})

	b := &bundler.Bundler{Config: bundler.Config{
		BaseURL:   "/public/",
		OutputDir: "./dist",
	}}
	b.Manifest = m
	t.Run("Bundler", func(t *testing.T) {
		require.NoError(t, b.Run())
	})
}
