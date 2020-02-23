package bundler

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/colinjfw/jbld/pkg/compiler"
)

// PublicConfig confiugures the public html plugin.
type PublicConfig struct {
	Dir  string   `json:"dir"`
	HTML []string `json:"html"`
}

// Config represents Bundler configuration.
type Config struct {
	BaseURL   string       `json:"baseUrl"`
	OutputDir string       `json:"outputDir"`
	AssetPath string       `json:"assetPath"`
	Public    PublicConfig `json:"public"`
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
	os.RemoveAll(b.OutputDir)
	if err := writePublicFolder(b.Config); err != nil {
		return errors.New("bundler: public folder not found")
	}

	bundles, err := (&bundleMapper{
		manifest: b.Manifest,
		Config:   b.Config,
	}).run()
	if err != nil {
		return err
	}

	for _, bn := range bundles {
		err = bn.Run(b.Manifest.Config.OutputDir, b.Config.OutputDir)
		if err != nil {
			return err
		}
	}
	m := b.manifest(bundles)
	if err := writeHTMLSources(b.Config, m); err != nil {
		return err
	}
	return b.writeManifest(m)
}

func (b *Bundler) manifest(bundles []*Bundle) *Manifest {
	m := &Manifest{
		BaseURL:     b.BaseURL,
		Bundles:     []string{},
		Entrypoints: map[string][]string{},
		BundleTypes: map[string][]string{},
	}
	for _, bn := range bundles {
		m.Bundles = append(m.Bundles, bn.URL)
		m.Entrypoints[bn.Name] = append(m.Entrypoints[bn.Name], bn.URL)
		m.BundleTypes[bn.Type] = append(m.BundleTypes[bn.Type], bn.URL)
	}
	return m
}

func (b *Bundler) writeManifest(m *Manifest) error {
	src := filepath.Join(b.OutputDir, "asset-manifest.json")
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
