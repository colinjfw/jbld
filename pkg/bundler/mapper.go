package bundler

import (
	"path/filepath"
	"strings"

	"github.com/colinjfw/jbld/pkg/compiler"
)

type bundleMapper struct {
	Config
	manifest *compiler.Manifest
	fileMap  map[string]compiler.File
}

func (s *bundleMapper) nameEntrypoint(name string) string {
	base := filepath.Base(name)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

func (s *bundleMapper) buildMap() {
	s.fileMap = map[string]compiler.File{}
	for _, f := range s.manifest.Files {
		s.fileMap[f.Name] = f
	}
}

func (s *bundleMapper) traverse(f compiler.File, emit func(f compiler.File)) {
	emit(f)
	for _, imp := range f.Object.Imports {
		s.traverse(s.fileMap[imp.Resolved], emit)
	}
}

func (s *bundleMapper) traverseCollect(f compiler.File) ([]compiler.File, []compiler.File, []compiler.File) {
	seen := map[string]bool{}
	js := []compiler.File{}
	css := []compiler.File{}
	url := []compiler.File{}
	s.traverse(f, func(n compiler.File) {
		if seen[n.Name] {
			return
		}

		seen[n.Name] = true
		js = append(js, n)
		switch n.Object.Type {
		case "css":
			css = append(css, n)
		case "js":
		default:
			url = append(url, n)
		}
	})
	return js, css, url
}

func (s *bundleMapper) run() ([]*Bundle, error) {
	s.buildMap()
	var bundles []*Bundle
	for _, entry := range s.manifest.Config.Entrypoints {
		jsFiles, cssFiles, urlFiles := s.traverseCollect(s.fileMap[entry])
		name := s.nameEntrypoint(entry)
		js := NewBundle(BundleCreate{
			Manifest: s.manifest,
			Config:   s.Config,
			Type:     "js",
			Name:     name,
			Main:     entry,
			Files:    jsFiles,
		})
		bundles = append(bundles, js)

		if len(cssFiles) > 0 {
			css := NewBundle(BundleCreate{
				Manifest: s.manifest,
				Config:   s.Config,
				Type:     "css",
				Name:     name,
				Main:     entry,
				Files:    cssFiles,
			})
			bundles = append(bundles, css)
			js.AddDependent(css.BundleID)
		}
		if len(urlFiles) > 0 {
			url := NewBundle(BundleCreate{
				Manifest: s.manifest,
				Config:   s.Config,
				Type:     "url",
				Name:     name,
				Main:     entry,
				Files:    urlFiles,
			})
			bundles = append(bundles, url)
			js.AddDependent(url.BundleID)
		}
	}
	return bundles, nil
}
