package bundler

import (
	"path/filepath"
	"strings"

	"github.com/colinjfw/jbld/pkg/compiler"
)

// MapRequest represents a configuration for mapping.
type MapRequest struct {
	Config
	Manifest *compiler.Manifest
}

// BundleMapper sets up a set of bundles to compile. Bundles are traversed
// starting with the entrypoint and selecting all
func BundleMapper(m MapRequest) ([]*Bundle, error) {
	s := &bundleMapper{MapRequest: m}
	s.buildMap()
	return s.run()
}

type bundleMapper struct {
	MapRequest
	fileMap map[string]compiler.File
}

func (s *bundleMapper) nameEntrypoint(name string) string {
	base := filepath.Base(name)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

func (s *bundleMapper) buildMap() {
	s.fileMap = map[string]compiler.File{}
	for _, f := range s.Manifest.Files {
		s.fileMap[f.Name] = f
	}
}

func (s *bundleMapper) traverse(f compiler.File, emit func(f compiler.File)) {
	emit(f)
	for _, imp := range f.Object.Imports {
		s.traverse(s.fileMap[imp.Resolved], emit)
	}
}

func (s *bundleMapper) traverseCollect(f compiler.File) ([]compiler.File, []compiler.File) {
	seen := map[string]bool{}
	js := []compiler.File{}
	css := []compiler.File{}
	s.traverse(f, func(n compiler.File) {
		if !seen[n.Name] {
			js = append(js, n)
			if n.Object.Type == "css" {
				css = append(css, n)
			}
		}
		seen[n.Name] = true
	})
	return js, css
}

func (s *bundleMapper) run() ([]*Bundle, error) {
	var bundles []*Bundle
	for _, entry := range s.Manifest.Config.Entrypoints {
		jsFiles, cssFiles := s.traverseCollect(s.fileMap[entry])
		name := s.nameEntrypoint(entry)
		js := NewBundle(BundleCreate{
			Manifest: s.Manifest,
			Config:   s.Config,
			Type:     "js",
			Name:     name,
			Main:     entry,
			Files:    jsFiles,
		})
		bundles = append(bundles, js)

		if len(cssFiles) > 0 {
			css := NewBundle(BundleCreate{
				Manifest: s.Manifest,
				Config:   s.Config,
				Type:     "css",
				Name:     name,
				Main:     entry,
				Files:    cssFiles,
			})
			bundles = append(bundles, css)
			js.AddDependent(css.BundleID)
		}
	}
	return bundles, nil
}
