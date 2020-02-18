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
	modLink       = "module.exports=%s;"
	resolvePrefix = "  Object.assign(__modules.resolve,"
	resolveSuffix = ");\n"
	mainStr       = "  __modules.main(%s, %s);\n"
)

// BundleID references a Bundle.
type BundleID struct{ FullName, URL, Name, Type, Hash string }

// MarshalJSON implements the marshalling interface.
func (b BundleID) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.URL)
}

// BundleCreate maps a request to create a new bundle object.
type BundleCreate struct {
	Config   Config
	Manifest *compiler.Manifest
	Name     string
	Type     string
	Main     string
	Files    []compiler.File
}

// NewBundle constructs a new bundle.
func NewBundle(r BundleCreate) *Bundle {
	b := &Bundle{
		BundleID: BundleID{
			Name: r.Name,
			Type: r.Type,
		},
		BaseURL:   r.Config.BaseURL,
		AssetPath: r.Config.AssetPath,
		Main:      r.Main,
		Files:     r.Files,
		Resolve:   r.Manifest.Resolve,
	}
	b.setHash()
	b.FullName = b.Name + "-" + b.Hash + "." + b.Type
	b.URL = filepath.Join(b.BaseURL, b.AssetPath, b.FullName)
	return b
}

// Bundle represents an individual bundle.
type Bundle struct {
	BundleID
	BaseURL   string
	AssetPath string
	Main      string
	Files     []compiler.File
	Bundles   []BundleID
	Resolve   map[string]string
}

// AddDependent adds a dependent chunk.
func (b *Bundle) AddDependent(id BundleID) *Bundle {
	b.Bundles = append(b.Bundles, id)
	return b
}

// Hash provides a sha256 for a bundle.
func (b *Bundle) setHash() {
	h := sha256.New()
	for _, f := range b.Files {
		h.Write([]byte(f.Object.Hash))
	}
	for _, s := range b.Bundles {
		h.Write([]byte(s.Type + "/" + s.Name))
	}
	b.Hash = hex.EncodeToString(h.Sum(nil))
}

// Run executes the bundler process.
func (b *Bundle) Run(srcDir, outputDir string) error {
	dstDir := filepath.Join(outputDir, b.AssetPath)
	t1 := time.Now()
	os.MkdirAll(dstDir, 0700)
	src := filepath.Join(dstDir, b.FullName)

	var w *bufio.Writer
	{
		f, err := os.OpenFile(src, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0700)
		if err != nil {
			return err
		}
		defer f.Close()
		w = bufio.NewWriter(f)
	}

	if b.Type == "js" {
		if err := b.bundleJS(w, srcDir); err != nil {
			return err
		}
	} else if b.Type == "css" {
		if err := b.bundleRaw(w, srcDir); err != nil {
			return err
		}
	} else {
		panic("TODO: Handle URLs")
	}

	log.Printf("bundler: entrypoint %s bundled %d files in %v",
		b.Name, len(b.Files), time.Since(t1),
	)
	return w.Flush()
}

func (b *Bundle) bundleRaw(w *bufio.Writer, srcDir string) error {
	for _, f := range b.Files {
		if err := b.bundleRawFile(srcDir, w, f); err != nil {
			return err
		}
	}
	return nil
}

func (b *Bundle) bundleJS(w *bufio.Writer, srcDir string) error {
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
		if f.Object.Type == "js" {
			if err := b.bundleFile(srcDir, w, f); err != nil {
				return err
			}
		} else {
			if err := b.bundleJSLink(srcDir, w, f); err != nil {
				return err
			}
		}
	}
	if err := b.bundleMain(w); err != nil {
		return err
	}
	if err := b.bundleBodyEnd(w); err != nil {
		return err
	}
	return nil
}

func (b *Bundle) bundleMain(w *bufio.Writer) error {
	if b.Main == "" {
		return nil
	}
	bundleIds, err := json.Marshal(b.Bundles)
	if err != nil {
		return err
	}
	_, err = w.WriteString(fmt.Sprintf(mainStr, string(bundleIds), strconv.Quote(b.Main)))
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

func (b *Bundle) bundleJSLink(srcDir string, w *bufio.Writer, file compiler.File) error {
	src := filepath.Join(srcDir, file.Name)
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = w.WriteString(fmt.Sprintf(modPrefix, strconv.Quote(file.Name)))
	if err != nil {
		return err
	}
	_, err = w.WriteString(fmt.Sprintf(modLink, strconv.Quote(
		filepath.Join(b.BaseURL, b.AssetPath, file.Name),
	)))
	if err != nil {
		return err
	}
	_, err = w.WriteString(modSuffix)
	if err != nil {
		return err
	}
	return nil
}

func (b *Bundle) bundleFile(srcDir string, w *bufio.Writer, file compiler.File) error {
	src := filepath.Join(srcDir, file.Name)
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()
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

func (b *Bundle) bundleRawFile(srcDir string, w *bufio.Writer, file compiler.File) error {
	src := filepath.Join(srcDir, file.Name)
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(w, f)
	if err != nil {
		return err
	}
	_, err = w.WriteString("\n")
	return err
}
