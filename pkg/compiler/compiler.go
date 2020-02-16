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

	"github.com/colinjfw/jbld/pkg/queue"
)

// Compiler represents the compilation process.
type Compiler struct {
	Config
}

// Run executes the compiler.
func (c *Compiler) Run() error {
	t1 := time.Now()

	h := NewHostPool(c.Config)
	defer h.Close()

	fw := &fileWriter{config: c.Config}
	count, err := queue.Run(c.Workers, c.Entrypoints,
		func(f string) ([]string, error) {
			o, err := c.process(f, h)
			if err != nil {
				return nil, err
			}
			fw.write(o)
			return o.ImportFiles(), nil
		},
	)

	if ferr := fw.flush(); ferr != nil {
		return ferr
	}

	log.Printf(
		"compiler: finished files=%d err=%v in=%v",
		count, err, time.Since(t1),
	)
	return err
}

func (c *Compiler) process(file string, host Host) (File, error) {
	src := filepath.Join(c.SourceDir, file)
	dst := filepath.Join(c.OutputDir, file)

	s := Source{
		Name:    file,
		Plugins: c.Plugins,
	}

	obj, err := readObj(dst)
	if err != nil {
		return File{}, err
	}
	// TODO: Also diff s.Plugins vs obj.Plugins.
	hash, err := hashFile(src)
	if err != nil {
		return File{}, err
	}
	if hash == obj.Hash {
		log.Printf(
			"compiler: process - cached: %s plugins=%v",
			s.Name, s.Plugins,
		)
		return File{
			Source: s,
			Object: obj,
		}, nil
	}

	log.Printf(
		"compiler: process - compiling: %s plugins=%v",
		s.Name, s.Plugins,
	)
	obj.Hash = hash
	obj.Plugins = c.Plugins
	obj.Imports, err = host.Run(s)
	if err != nil {
		return File{}, err
	}

	err = writeObj(dst, obj)
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

type fileWriter struct {
	config Config
	files  []File
	lock   sync.Mutex
}

func (fw *fileWriter) write(f File) {
	fw.lock.Lock()
	fw.files = append(fw.files, f)
	fw.lock.Unlock()
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
	return enc.Encode(struct {
		Files  []File `json:"files"`
		Config Config `json:"config"`
	}{Files: fw.files, Config: fw.config})
}
