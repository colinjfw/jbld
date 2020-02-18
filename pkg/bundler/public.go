package bundler

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"os/exec"
	"path/filepath"
)

func writePublicFolder(conf Config) error {
	if conf.Public.Dir == "" {
		return nil
	}
	return exec.Command("cp", "-r", conf.Public.Dir, conf.OutputDir).Run()
}

func writeHTMLSources(conf Config, m *Manifest) error {
	for _, f := range conf.Public.HTML {
		dst := filepath.Join(conf.OutputDir, f)
		data, err := ioutil.ReadFile(dst)
		if err != nil {
			return err
		}
		tpl, err := template.New("").Parse(string(data))
		if err != nil {
			return err
		}
		buf := bytes.NewBuffer(nil)
		err = tpl.Execute(buf, m)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(dst, buf.Bytes(), 0700)
		if err != nil {
			return err
		}
	}
	return nil
}
