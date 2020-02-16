package compiler

// Import represents a mapping from a file to an import.
type Import struct {
	// Kind can be 'async' or 'static' or other plugin provided types.
	Kind string `json:"kind"`
	// Name represents the original name inside the file.
	Name string `json:"name"`
	// Resolved represents the resolved path relative to Source or Output dirs.
	Resolved string `json:"resolved"`
}

// Source represents a source file.
type Source struct {
	Name    string   `json:"name"`
	Plugins []string `json:"plugins"`
}

// Object holds information about a file.
type Object struct {
	// Hash is used for cache busting.
	Hash string `json:"hash"`
	// Imports is the list of file imports.
	Imports []Import `json:"imports"`
	// Plugins represents the plugins used on this file.
	Plugins []string `json:"plugins"`
}

// Config represents a compiler configuration.
type Config struct {
	HostJS      string   `json:"hostJS"`
	ConfigFile  string   `json:"configFile"`
	Entrypoints []string `json:"entrypoints"`
	SourceDir   string   `json:"sourceDir"`
	OutputDir   string   `json:"outputDir"`
	Plugins     []string `json:"plugins"`
	Workers     int      `json:"workers"`
}

// File holds a compiled file.
type File struct {
	Source
	Object Object `json:"object"`
}

// ImportFiles represents a list of all named imports.
func (f File) ImportFiles() (out []string) {
	for _, i := range f.Object.Imports {
		out = append(out, i.Resolved)
	}
	return
}
