package main_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	_ "github.com/colinjfw/jbld/pkg/run"
	"github.com/stretchr/testify/require"
)

func TestCompiler(t *testing.T) {
	os.RemoveAll("./dist")
	cwd, _ := os.Getwd()

	build := exec.Command("./build.sh")
	build.Stderr = os.Stderr
	build.Stdout = os.Stdout
	build.Dir = filepath.Join(cwd, "..", "..", "lib")
	require.NoError(t, build.Run())

	run := exec.Command("../../lib/cli.js")
	run.Stderr = os.Stderr
	run.Stdout = os.Stdout
	require.NoError(t, run.Run())
}
