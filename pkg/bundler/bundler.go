package bundler

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/colinjfw/jbld/pkg/compiler"
)

// Manifest represents the output manifest.
type Manifest struct {
	Config      Config              `json:"config"`
	Entrypoints map[string][]string `json:"entrypoints"`
	Bundles     []string            `json:"bundles"`
	HTML        string              `json:"html"`
}

// Config represents Bundler configuration.
type Config struct {
	BaseURL     string       `json:"baseUrl"`
	OutputDir   string       `json:"outputDir"`
	Entrypoints []Entrypoint `json:"entrypoints"`
}

// Entrypoint configures a compilation target entrypoint.
type Entrypoint struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// Bundler implements a bundling.
type Bundler struct {
	Config
	Manifest *compiler.Manifest `json:"manifest"`
}

// Run will execute the bundler process by calling the BundleMapper function and
// then running the set of optimizers.
func (b *Bundler) Run() error {
	bundles, err := BundleMapper(MapRequest{
		Manifest: b.Manifest,
		Config:   b.Config,
	})
	if err != nil {
		return err
	}

	for _, bn := range bundles {
		err = bn.Run(b.Manifest.Config.OutputDir, b.Config.OutputDir)
		if err != nil {
			return err
		}
	}
	return b.writeManifest(bundles)
}

func (b *Bundler) writeManifest(bundles []*Bundle) error {
	m := Manifest{Config: b.Config, Entrypoints: map[string][]string{}}
	m.HTML = b.html(bundles)
	for _, bn := range bundles {
		m.Bundles = append(m.Bundles, bn.URL)
		m.Entrypoints[bn.Name] = append(m.Entrypoints[bn.Name], bn.URL)
	}
	src := filepath.Join(b.OutputDir, ".jbld-bundle-manifest")
	f, err := os.OpenFile(src, os.O_CREATE|os.O_RDWR, 0700)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	return enc.Encode(m)
}

func (b *Bundler) html(bundles []*Bundle) string {
	chunkMap := map[string]bool{}
	for _, bn := range bundles {
		chunkMap[bn.URL] = true
	}
	chunks, _ := json.Marshal(chunkMap)
	out := []string{
		fmt.Sprintf("<script>window.__chunks=%s</script>", string(chunks)),
	}
	for _, bn := range bundles {
		switch bn.Type {
		case "css":
			out = append(out, fmt.Sprintf(
				"<link rel=\"stylesheet\" type=\"text/css\" href=\"%s\">",
				bn.URL,
			))
		case "js":
			out = append(out, fmt.Sprintf(
				"<script type=\"application/javascript\" src=\"%s\"></script>",
				bn.URL,
			))
		}
	}
	return strings.Join(out, "\n")
}
