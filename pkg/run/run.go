package run

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/colinjfw/jbld/pkg/bundler"
	"github.com/colinjfw/jbld/pkg/compiler"
	"github.com/radovskyb/watcher"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

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
func Run(conf *Options) error {
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

func LoadOptions(opts string) (*Options, error) {
	conf := &Options{}
	if err := json.Unmarshal([]byte(opts), conf); err != nil {
		return nil, err
	}
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
