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
	"time"

	"github.com/colinjfw/jbld/pkg/queue"
)

// Compiler represents the compilation process.
type Compiler struct {
	HostJS      string
	ConfigFile  string
	Entrypoints []string
	SourceDir   string
	OutputDir   string
	Plugins     []string
	Workers     int
}

// Run executes the compiler.
func (c *Compiler) Run() error {
	t1 := time.Now()
	h := NewHostPool(c.Workers, c.HostJS, c.ConfigFile)
	defer h.Close()
	count, err := queue.Run(c.Workers, c.Entrypoints, func(f string) ([]string, error) {
		o, err := c.process(f, h)
		if err != nil {
			return nil, err
		}
		return o.ImportFiles(), nil
	})
	log.Printf("compiler: finished files=%d err=%v in=%v", count, err, time.Since(t1))
	return err
}

func (c *Compiler) process(file string, host Host) (File, error) {
	s := Source{
		Src:     filepath.Join(c.SourceDir, file),
		Dst:     filepath.Join(c.OutputDir, file),
		Plugins: c.Plugins,
	}

	obj, err := readObj(s.Dst)
	if err != nil {
		return File{}, err
	}
	// TODO: Also diff s.Plugins vs obj.Plugins.
	hash, err := hashFile(s.Src)
	if err != nil {
		return File{}, err
	}
	if hash == obj.Hash {
		log.Printf("compiler: process - cached: %s -> %s plugins=%v", s.Src, s.Dst, s.Plugins)
		return File{
			Source: s,
			Object: obj,
		}, nil
	}

	log.Printf("compiler: process - compiling: %s -> %s plugins=%v", s.Src, s.Dst, s.Plugins)
	obj.Hash = hash
	obj.Plugins = c.Plugins
	obj.Imports, err = host.Run(s)
	if err != nil {
		return File{}, err
	}

	err = writeObj(s.Dst, obj)
	if err != nil {
		return File{}, err
	}
	return File{
		Source: s,
		Object: obj,
	}, nil
}

func hashFile(src string) (string, error) {
	h := sha256.New()
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
