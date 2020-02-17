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
	for _, e := range s.Entrypoints {
		if e.Path == name {
			return e.Name
		}
	}
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
	for _, entry := range s.Manifest.Config.Entrypoints {
		var group []*Bundle
		for typ, files := range s.traverseCollect(s.fileMap[entry]) {
			name := s.nameEntrypoint(entry)
			group = append(group,
				NewBundle(BundleCreate{
					Manifest: s.Manifest,
					Config:   s.Config,
					Type:     typ,
					Name:     name,
					Main:     entry,
					Files:    files,
				}),
			)
		}
		for i, b := range group {
			for i2, b2 := range group {
				if i == i2 {
					continue
				}
				b.AddDependent(b2.BundleID)
			}
		}
		bundles = append(bundles, group...)
	}
	return bundles, nil
}
