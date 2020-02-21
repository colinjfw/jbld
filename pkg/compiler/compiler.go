package compiler

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/colinjfw/jbld/pkg/host"
	"github.com/colinjfw/jbld/pkg/queue"
)

// ProcessRequest is sent to the Configuration JS class.
type ProcessRequest struct {
	Name   string `json:"name"`
	Src    string `json:"src"`
	Dst    string `json:"dst"`
	SrcDir string `json:"srcDir"`
	DstDir string `json:"dstDir"`
}

// ProcessResponse is expected from the Configuration JS class.
type ProcessResponse struct {
	Type    string   `json:"type"`
	Imports []Import `json:"imports"`
}

// Compiler represents the compilation process.
type Compiler struct {
	Config
}

// Run executes the compiler.
func (c *Compiler) Run() (*Manifest, error) {
	t1 := time.Now()

	h := host.NewHostPool(c.Workers, c.HostJS, c.ConfigFile)
	defer h.Close()

	confHash, err := hashFile(c.ConfigFile)
	if err != nil {
		return nil, err
	}

	fw := &fileWriter{config: c.Config, resolve: map[string]string{}}
	count, err := queue.Run(c.Workers, c.Entrypoints,
		func(f string) ([]string, error) {
			o, err := c.process(confHash, f, h)
			if err != nil {
				return nil, err
			}
			fw.write(o)
			return o.ImportFiles(), nil
		},
	)
	if err != nil {
		return nil, err
	}

	if ferr := fw.flush(); ferr != nil {
		return nil, ferr
	}

	log.Printf("compiler: compiled %d files in %v", count, time.Since(t1))
	return fw.manifest(), nil
}

func (c *Compiler) process(confHash, file string, h host.Host) (File, error) {
	src := filepath.Join(c.SourceDir, file)
	dst := filepath.Join(c.OutputDir, file)

	s := Source{
		Name: file,
	}

	obj, err := readObj(dst)
	if err != nil {
		return File{}, err
	}
	hash, err := hashFile(src, confHash)
	if err != nil {
		return File{}, err
	}
	if hash == obj.Hash {
		log.Printf("compiler: process - cached: %s", s.Name)
		return File{
			Source: s,
			Object: obj,
		}, nil
	}

	log.Printf("compiler: process - compiling: %s", s.Name)
	resp := ProcessResponse{}
	err = h.Run("process", ProcessRequest{
		Name:   file,
		Src:    src,
		Dst:    dst,
		SrcDir: c.SourceDir,
		DstDir: c.OutputDir,
	}, &resp)
	if err != nil {
		return File{}, err
	}
	obj.Hash = hash
	obj.Imports = resp.Imports
	obj.Type = resp.Type

	err = writeObj(dst, obj)
	if err != nil {
		return File{}, err
	}
	return File{
		Source: s,
		Object: obj,
	}, nil
}

func hashFile(src string, parts ...string) (string, error) {
	h := sha256.New()
	for _, p := range parts {
		h.Write([]byte(p))
	}
	f, err := os.Open(src)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func readObj(dst string) (Object, error) {
	data, err := ioutil.ReadFile(dst + ".o")
	if os.IsNotExist(err) {
		return Object{}, nil
	}
	if err != nil {
		return Object{}, err
	}
	o := Object{}
	err = json.Unmarshal(data, &o)
	if err != nil {
		return Object{}, err
	}
	return o, nil
}

func writeObj(dst string, o Object) error {
	data, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dst+".o", data, 0700)
}

type fileWriter struct {
	config  Config
	files   []File
	resolve map[string]string
	lock    sync.Mutex
}

func (fw *fileWriter) write(f File) {
	fw.lock.Lock()
	fw.files = append(fw.files, f)
	for _, imp := range f.Object.Imports {
		if imp.Name != imp.Resolved {
			fw.resolve[imp.Name] = imp.Resolved
		}
	}
	fw.lock.Unlock()
}

func (fw *fileWriter) manifest() *Manifest {
	return &Manifest{
		Version: fw.config.Version(),
		Config:  fw.config,
		Files:   fw.files,
		Resolve: fw.resolve,
	}
}

func (fw *fileWriter) flush() error {
	fw.lock.Lock()
	defer fw.lock.Unlock()

	os.MkdirAll(fw.config.OutputDir, 0700)

	dst := filepath.Join(fw.config.OutputDir, ".jbld-manifest")
	f, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0700)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(fw.manifest())
}
