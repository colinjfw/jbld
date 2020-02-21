package compiler

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

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
	Name string `json:"name"`
}

// Object holds information about a file.
type Object struct {
	// Type represents a fileType.
	Type string `json:"type"`
	// Hash is used for cache busting.
	Hash string `json:"hash"`
	// Imports is the list of file imports.
	Imports []Import `json:"imports"`
}

// Config represents a compiler configuration.
type Config struct {
	HostJS      string   `json:"hostJs"`
	ConfigFile  string   `json:"configFile"`
	Entrypoints []string `json:"entrypoints"`
	SourceDir   string   `json:"sourceDir"`
	OutputDir   string   `json:"outputDir"`
	Workers     int      `json:"workers"`
}

// Version describes a unique hash for this config.
func (c Config) Version() string {
	h := sha256.New()
	err := json.NewEncoder(h).Encode(c)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(h.Sum(nil))
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

// Manifest represents result of compilation work.
type Manifest struct {
	Version string            `json:"version"`
	Resolve map[string]string `json:"resolve"`
	Config  Config            `json:"config"`
	Files   []File            `json:"files"`
}
