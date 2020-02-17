package bundler

import (
	"path/filepath"
	"strings"

	"github.com/colinjfw/jbld/pkg/compiler"
)

// BundleMapper sets up.
func BundleMapper(b *Bundler) ([]*Bundle, error) {
	s := &bundleMapper{b: b}
	s.buildMap()
	return s.run()
}

type bundleMapper struct {
	b       *Bundler
	fileMap map[string]compiler.File
}

func (s *bundleMapper) nameEntrypoint(name string) string {
	for _, e := range s.b.Entrypoints {
		if e.Path == name {
			return e.Name
		}
	}
	base := filepath.Base(name)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

func (s *bundleMapper) buildMap() {
	s.fileMap = map[string]compiler.File{}
	for _, f := range s.b.Manifest.Files {
		s.fileMap[f.Name] = f
	}
}

func (s *bundleMapper) traverse(f compiler.File, emit func(f compiler.File)) {
	emit(f)
	for _, imp := range f.Object.Imports {
		s.traverse(s.fileMap[imp.Resolved], emit)
	}
}

func (s *bundleMapper) traverseCollect(f compiler.File) (out []compiler.File) {
	seen := map[string]bool{}
	s.traverse(f, func(n compiler.File) {
		if !seen[n.Name] {
			out = append(out, n)
		}
		seen[n.Name] = true
	})
	return
}

func (s *bundleMapper) run() ([]*Bundle, error) {
	var bundles []*Bundle
	for _, entry := range s.b.Manifest.Config.Entrypoints {
		bundles = append(bundles, &Bundle{
			Primary: true,
			Name:    s.nameEntrypoint(entry),
			Main:    entry,
			Files:   s.traverseCollect(s.fileMap[entry]),
			Bundles: []string{},
			Resolve: s.b.Manifest.Resolve,
		})
	}
	return bundles, nil
}
