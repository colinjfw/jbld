package bundler

import "errors"

// Optimizer function.
type Optimizer func(b *Bundler, bs []*Bundle) ([]*Bundle, error)

// GetOptimizer loads an optimizer by name.
func GetOptimizer(name string) (Optimizer, error) {
	switch name {
	case "CommonChunk":
		return CommonChunkOptimizer, nil
	default:
		return nil, errors.New("optimizer does not exist")
	}
}

// CommonChunkOptimizer splits a set of bundles into chunks.
func CommonChunkOptimizer(b *Bundler, bundles []*Bundle) ([]*Bundle, error) {
	// TODO
	return bundles, nil
}
