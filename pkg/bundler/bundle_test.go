package bundler

import (
	"os"
	"os/exec"
	"testing"

	"github.com/colinjfw/jbld/pkg/compiler"
	"github.com/stretchr/testify/require"
)

func TestBundle(t *testing.T) {
	t.Run("JS", func(t *testing.T) {
		cwd, _ := os.Getwd()
		b := &Bundle{
			BundleID:  BundleID{Type: "js", Name: "test", FullName: "test.js"},
			BaseURL:   "/",
			AssetPath: "",
			Main:      "index.js",
			Bundles:   []BundleID{},
			Resolve:   map[string]string{},
			Files: []compiler.File{
				{
					Source: compiler.Source{Name: "index.js"},
					Object: compiler.Object{Type: "js", Hash: "test"},
				},
			},
		}
		err := b.Run(cwd+"/testdata/src", cwd+"/testdata/out")
		require.NoError(t, err)

		out, err := exec.Command("node", "./testdata/out/test.js").CombinedOutput()
		require.NoError(t, err)
		require.Equal(t, "HELLO\n", string(out))
	})

	t.Run("CSS", func(t *testing.T) {
		cwd, _ := os.Getwd()
		b := &Bundle{
			BundleID:  BundleID{Type: "css", Name: "test", FullName: "test.css"},
			BaseURL:   "/",
			AssetPath: "",
			Bundles:   []BundleID{},
			Resolve:   map[string]string{},
			Files: []compiler.File{
				{
					Source: compiler.Source{Name: "index.css"},
					Object: compiler.Object{Type: "css", Hash: "test"},
				},
			},
		}
		err := b.Run(cwd+"/testdata/src", cwd+"/testdata/out")
		require.NoError(t, err)
	})

	t.Run("URL", func(t *testing.T) {
		cwd, _ := os.Getwd()
		b := &Bundle{
			BundleID:  BundleID{Type: "url", Name: "url", FullName: "url.svg"},
			BaseURL:   "/",
			AssetPath: "",
			Bundles:   []BundleID{},
			Resolve:   map[string]string{},
			Files: []compiler.File{
				{
					Source: compiler.Source{Name: "url.svg"},
					Object: compiler.Object{Type: "svg", Hash: "test"},
				},
			},
		}
		err := b.Run(cwd+"/testdata/src", cwd+"/testdata/out")
		require.NoError(t, err)
	})

	t.Run("JSImport", func(t *testing.T) {
		cwd, _ := os.Getwd()
		b := &Bundle{
			BundleID:  BundleID{Type: "js", Name: "import-test", FullName: "import-test.js"},
			BaseURL:   "/",
			AssetPath: "",
			Main:      "url-import.js",
			Bundles:   []BundleID{},
			Resolve:   map[string]string{},
			Files: []compiler.File{
				{
					Source: compiler.Source{Name: "url-import.js"},
					Object: compiler.Object{Type: "js", Hash: "test"},
				},
				{
					Source: compiler.Source{Name: "url.svg"},
					Object: compiler.Object{Type: "svg", Hash: "test"},
				},
				{
					Source: compiler.Source{Name: "index.css"},
					Object: compiler.Object{Type: "css", Hash: "test"},
				},
			},
		}
		err := b.Run(cwd+"/testdata/src", cwd+"/testdata/out")
		require.NoError(t, err)

		out, err := exec.Command("node", "./testdata/out/import-test.js").CombinedOutput()
		require.NoError(t, err)
		require.Equal(t, "/url-test.svg\n", string(out))
	})
}
