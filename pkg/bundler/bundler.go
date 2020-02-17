package bundler

import (
	"github.com/colinjfw/jbld/pkg/compiler"
)

// Config represents Bundler configuration.
type Config struct {
	OutputDir   string       `json:"outputDir"`
	Entrypoints []Entrypoint `json:"entrypoints"`
	Optimizers  []string     `json:"optimizers"`
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
	bundles, err := BundleMapper(b)
	if err != nil {
		return err
	}
	for _, optim := range b.Optimizers {
		op, err := GetOptimizer(optim)
		if err != nil {
			return err
		}
		bundles, err = op(b, bundles)
	}
	for _, bn := range bundles {
		err = bn.Run(b.Manifest.Config.OutputDir, b.Config.OutputDir)
		if err != nil {
			return err
		}
	}
	return nil
}
