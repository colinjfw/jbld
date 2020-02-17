package bundler

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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
	bodyPrefix    = "(function(){\n"
	bodySuffix    = "})();\n"
	modPrefix     = "  __modules.define(%s, function(module, exports, require) {"
	modSuffix     = "});\n"
	resolvePrefix = "  Object.assign(__modules.resolve,"
	resolveSuffix = ");\n"
	mainStr       = "  __modules.require(%s);\n"
)

// Bundle represents an individual bundle.
type Bundle struct {
	Primary bool              `json:"primary"`
	Main    string            `json:"main"`
	Name    string            `json:"name"`
	Files   []compiler.File   `json:"files"`
	Bundles []string          `json:"bundles"`
	Resolve map[string]string `json:"resolve"`
}

// FullName provides the full bundle name.
func (b *Bundle) FullName() string {
	return b.Name + "." + b.Hash() + ".js"
}

// Hash provides a sha256 for a bundle.
func (b *Bundle) Hash() string {
	h := sha256.New()
	for _, f := range b.Files {
		h.Write([]byte(f.Object.Hash))
	}
	for _, s := range b.Bundles {
		h.Write([]byte(s))
	}
	return hex.EncodeToString(h.Sum(nil))
}

// Run executes the bundler process.
func (b *Bundle) Run(srcDir, dstDir string) error {
	t1 := time.Now()
	os.MkdirAll(dstDir, 0700)
	src := filepath.Join(dstDir, b.FullName())

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
	if err := b.bundleBodyStart(w); err != nil {
		return err
	}
	if err := b.bundleResolve(w); err != nil {
		return err
	}
	for _, f := range b.Files {
		if err := b.bundleFile(srcDir, w, f); err != nil {
			return err
		}
	}
	if err := b.bundleMain(w); err != nil {
		return err
	}
	if err := b.bundleBodyEnd(w); err != nil {
		return err
	}

	log.Printf(
		"bundler: bundled %d files in %v",
		len(b.Files), time.Since(t1),
	)
	return w.Flush()
}

func (b *Bundle) bundleMain(w *bufio.Writer) error {
	if b.Main == "" {
		return nil
	}
	_, err := w.WriteString(fmt.Sprintf(mainStr, strconv.Quote(b.Main)))
	return err
}

func (b *Bundle) bundleBodyStart(w *bufio.Writer) error {
	_, err := w.WriteString(bodyPrefix)
	return err
}
func (b *Bundle) bundleBodyEnd(w *bufio.Writer) error {
	_, err := w.WriteString(bodySuffix)
	return err
}

func (b *Bundle) bundleRuntime(w *bufio.Writer) error {
	_, err := w.WriteString(runtime)
	return err
}

func (b *Bundle) bundleResolve(w *bufio.Writer) error {
	js, err := json.Marshal(b.Resolve)
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

func (b *Bundle) bundleFile(srcDir string, w *bufio.Writer, file compiler.File) error {
	src := filepath.Join(srcDir, file.Name)
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
