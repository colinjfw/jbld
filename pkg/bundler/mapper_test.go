package bundler

import (
	"testing"

	"github.com/colinjfw/jbld/pkg/compiler"
	"github.com/stretchr/testify/require"
)

func TestMapper(t *testing.T) {
	m := &compiler.Manifest{
		Config: compiler.Config{Entrypoints: []string{"index.js"}},
		Files: []compiler.File{
			{
				Source: compiler.Source{Name: "index.js"},
				Object: compiler.Object{
					Type: "js", Hash: "test",
					Imports: []compiler.Import{
						{Name: "index.css", Resolved: "index.css"},
						{Name: "index.svg", Resolved: "index.svg"},
					},
				},
			},
			{
				Source: compiler.Source{Name: "index.css"},
				Object: compiler.Object{Type: "css", Hash: "test"},
			},
			{
				Source: compiler.Source{Name: "index.svg"},
				Object: compiler.Object{Type: "svg", Hash: "test"},
			},
		},
	}
	bm := &bundleMapper{
		manifest: m,
	}
	bundles, err := bm.run()
	require.NoError(t, err)

	require.Equal(t, 3, len(bundles))
	require.Equal(t, "index", bundles[0].Name)
	require.Equal(t, "js", bundles[0].Type)
	require.Equal(t, "index", bundles[1].Name)
	require.Equal(t, "css", bundles[1].Type)
	require.Equal(t, "index", bundles[2].Name)
	require.Equal(t, "url", bundles[2].Type)
}
