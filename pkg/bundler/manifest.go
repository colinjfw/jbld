package bundler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
)

// Manifest represents the output manifest.
type Manifest struct {
	BaseURL     string              `json:"baseURL"`
	Entrypoints map[string][]string `json:"entrypoints"`
	Bundles     []string            `json:"bundles"`
	BundleTypes map[string][]string `json:"bundleTypes"`
}

// ConfigJS sets up require JS configuration.
func (m *Manifest) ConfigJS() template.HTML {
	chunks, _ := json.Marshal(m.Bundles)
	return template.HTML(
		fmt.Sprintf("<script>window.__chunks=%s</script>", string(chunks)),
	)
}

// JS returns a set of JS elements.
func (m *Manifest) JS() template.HTML {
	out := []string{}
	for _, url := range m.Bundles {
		switch filepath.Ext(url) {
		case ".js":
			out = append(out, fmt.Sprintf(
				"<script type=\"application/javascript\" src=\"%s\"></script>",
				url,
			))
		}
	}
	return template.HTML(strings.Join(out, ""))
}

// CSS returns a set of link elements..
func (m *Manifest) CSS() template.HTML {
	out := []string{}
	for _, url := range m.Bundles {
		switch filepath.Ext(url) {
		case ".css":
			out = append(out, fmt.Sprintf(
				"<link rel=\"stylesheet\" type=\"text/css\" href=\"%s\">",
				url,
			))
		}
	}
	return template.HTML(strings.Join(out, ""))
}
