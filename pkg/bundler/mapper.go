package bundler

import (
	"path/filepath"
	"strings"

	"github.com/colinjfw/jbld/pkg/compiler"
)

// BundleMapper sets up a set of bundles to compile. Bundles are traversed
// starting with the entrypoint and selecting all
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

func (s *bundleMapper) traverseCollect(f compiler.File) map[string][]compiler.File {
	seen := map[string]bool{}
	out := map[string][]compiler.File{}
	s.traverse(f, func(n compiler.File) {
		if !seen[n.Name] {
			out[n.Object.Type] = append(out[n.Object.Type], n)
		}
		seen[n.Name] = true
	})
	return out
}

func (s *bundleMapper) run() ([]*Bundle, error) {
	var bundles []*Bundle
	for _, entry := range s.b.Manifest.Config.Entrypoints {
		collect := s.traverseCollect(s.fileMap[entry])
		common := []BundleID{}
		for typ := range collect {
			common = append(common, BundleID{
				Type: typ,
				Name: s.nameEntrypoint(entry),
			})
		}
		for typ, files := range collect {
			bundles = append(bundles, &Bundle{
				Primary: true,
				Type:    typ,
				Name:    s.nameEntrypoint(entry),
				Main:    entry,
				Files:   files,
				Bundles: common,
				Resolve: s.b.Manifest.Resolve,
			})
		}
	}
	return bundles, nil
}
