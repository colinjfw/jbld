package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	cwd, _ := os.Getwd()

	build := exec.Command("./build.sh")
	build.Stderr = os.Stderr
	build.Stdout = os.Stdout
	build.Dir = filepath.Join(cwd, "..", "..", "lib")
	if err := build.Run(); err != nil {
		panic(err)
	}

	run := exec.Command("../../lib/cli.js", "--watch", "--serve", ":3000")
	run.Stderr = os.Stderr
	run.Stdout = os.Stdout
	if err := run.Run(); err != nil {
		panic(err)
	}
}
