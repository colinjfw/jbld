package bundler

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/colinjfw/jbld/pkg/compiler"
)

const (
	modPrefix     = "__modules.define(%s, function(module, exports, require) {"
	modSuffix     = "});\n"
	resolvePrefix = "__modules.resolve="
	resolveSuffix = ";\n"
)

// Config represents Bundler configuration.
type Config struct {
	OutputDir string
}

// Bundler implements a bundling.
type Bundler struct {
	Config
	Manifest *compiler.Manifest
}

// Run executes the bundler process.
func (b *Bundler) Run() error {
	t1 := time.Now()
	if b.Manifest == nil {
		return errors.New("bundler: manifest is nil")
	}
	os.MkdirAll(b.OutputDir, 0700)
	src := filepath.Join(b.OutputDir, "bundle.js")

	var w *bufio.Writer
	{
		f, err := os.OpenFile(src, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0700)
		if err != nil {
			return err
		}
		defer f.Close()
		w = bufio.NewWriter(f)
	}

	if err := b.bundleRuntime(w); err != nil {
		return err
	}
	if err := b.bundleResolve(w); err != nil {
		return err
	}
	for _, f := range b.Manifest.Files {
		if err := b.bundleFile(w, f); err != nil {
			return err
		}
	}

	log.Printf(
		"bundler: bundled %d files in %v",
		len(b.Manifest.Files), time.Since(t1),
	)
	return w.Flush()
}

func (b *Bundler) bundleRuntime(w *bufio.Writer) error {
	_, err := w.WriteString(runtime)
	return err
}

func (b *Bundler) bundleResolve(w *bufio.Writer) error {
	js, err := json.Marshal(b.Manifest.Resolve)
	if err != nil {
		return err
	}
	_, err = w.WriteString(resolvePrefix)
	if err != nil {
		return err
	}
	_, err = w.Write(js)
	if err != nil {
		return err
	}
	_, err = w.WriteString(resolveSuffix)
	if err != nil {
		return err
	}
	return err
}

func (b *Bundler) bundleFile(w *bufio.Writer, file compiler.File) error {
	src := filepath.Join(b.Manifest.Config.OutputDir, file.Name)
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	_, err = w.WriteString(fmt.Sprintf(modPrefix, strconv.Quote(file.Name)))
	if err != nil {
		return err
	}
	_, err = io.Copy(w, f)
	if err != nil {
		return err
	}
	_, err = w.WriteString(modSuffix)
	if err != nil {
		return err
	}
	return nil
}
