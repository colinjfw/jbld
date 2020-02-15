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
	Src     string   `json:"src"`
	Dst     string   `json:"dst"`
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

// File holds a compiled file.
type File struct {
	Object
	Source
}

// ImportFiles represents a list of all named imports.
func (f File) ImportFiles() (out []string) {
	for _, i := range f.Imports {
		out = append(out, i.Resolved)
	}
	return
}
