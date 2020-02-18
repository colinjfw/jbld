package run

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/colinjfw/jbld/pkg/bundler"
	"github.com/colinjfw/jbld/pkg/compiler"
	"github.com/colinjfw/jbld/pkg/host"
	"github.com/radovskyb/watcher"
)

// Options represents a unioned configuration.
type Options struct {
	Watch    bool              `json:"watch"`
	Serve    string            `json:"serve"`
	Mode     string            `json:"mode"`
	Env      map[string]string `json:"env"`
	Bundler  bundler.Config    `json:"bundler"`
	Compiler compiler.Config   `json:"compiler"`
}

// Run executes a full pipeline.
func Run(hostJS, configFile string) error {
	conf, err := loadOptions(hostJS, configFile)
	if err != nil {
		return err
	}

	if conf.Serve != "" {
		go serve(conf.Serve, conf.Bundler.OutputDir)
	}
	if conf.Watch {
		return watch(conf, func() {
			run(conf)
		})
	}
	return run(conf)
}

func defaultDir(cwd, val, def string) string {
	if val == "" {
		return def
	}
	if val[0] == '/' {
		return val
	}
	return filepath.Join(cwd, val)
}

func (c *Options) withDefaults() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if c.Mode == "" {
		c.Mode = "production"
	}
	if c.Bundler.AssetPath == "" {
		c.Bundler.AssetPath = "static"
	}
	if c.Bundler.BaseURL == "" {
		c.Bundler.BaseURL = "/"
	}
	if c.Serve != "" {
		c.Watch = true
	}
	c.Bundler.OutputDir = defaultDir(cwd, c.Bundler.OutputDir, filepath.Join(cwd, "dist", "bundle"))
	c.Compiler.OutputDir = defaultDir(cwd, c.Compiler.OutputDir, filepath.Join(cwd, "dist", "target"))
	c.Compiler.SourceDir = defaultDir(cwd, c.Compiler.SourceDir, filepath.Join(cwd, "src"))

	if c.Compiler.Workers == 0 {
		c.Compiler.Workers = 5
	}
	if len(c.Compiler.Entrypoints) == 0 {
		c.Compiler.Entrypoints = []string{"index.js"}
	}
}

func loadOptions(hostJS, configFile string) (*Options, error) {
	h := host.NewHost(hostJS, configFile)
	conf := &Options{}
	if err := h.Run("options", struct{}{}, conf); err != nil {
		h.Close()
		return nil, err
	}
	h.Close()

	conf.withDefaults()
	conf.Compiler.HostJS = hostJS
	conf.Compiler.ConfigFile = configFile

	os.Setenv("NODE_ENV", conf.Mode)
	for k, v := range conf.Env {
		os.Setenv(k, v)
	}

	data, _ := json.MarshalIndent(conf, "", "  ")
	log.Printf("run: configuration\n%s", string(data))
	return conf, nil
}

func run(conf *Options) error {
	comp := &compiler.Compiler{
		Config: conf.Compiler,
	}
	manifest, err := comp.Run()
	if err != nil {
		return err
	}
	bund := &bundler.Bundler{
		Config:   conf.Bundler,
		Manifest: manifest,
	}
	return bund.Run()
}

func serve(listen, dir string) {
	log.Printf("run: serving %s on %s", dir, listen)
	http.Handle("/", http.FileServer(http.Dir(dir)))
	if err := http.ListenAndServe(listen, nil); err != nil {
		log.Fatalf("run: server failed - %v", err)
	}
}

func watch(o *Options, cb func()) error {
	w := watcher.New()
	w.SetMaxEvents(1)
	if err := w.AddRecursive(o.Compiler.SourceDir); err != nil {
		return err
	}

	cb()
	go func() {
		log.Printf("run: watcher started - %s", o.Compiler.SourceDir)
		for {
			select {
			case e, ok := <-w.Event:
				if !ok {
					return
				}
				if strings.HasPrefix(e.Path, o.Bundler.OutputDir) {
					continue
				}
				cb()
			case err, ok := <-w.Error:
				if !ok {
					return
				}
				log.Println("watcher: error -", err)
			}
		}
	}()
	return w.Start(500 * time.Millisecond)
}
