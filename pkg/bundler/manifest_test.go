package bundler

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestManifest(t *testing.T) {
	m := &Manifest{
		Bundles: []string{"foo.js", "foo.css"},
	}

	require.Equal(t, `<script>window.__chunks=["foo.js","foo.css"]</script>`, string(m.ConfigJS()))
	require.Equal(t, `<script type="application/javascript" src="foo.js"></script>`, string(m.JS()))
	require.Equal(t, `<link rel="stylesheet" type="text/css" href="foo.css">`, string(m.CSS()))
}
