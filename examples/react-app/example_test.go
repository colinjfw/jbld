package main_test

import (
	"os"
	"testing"

	"github.com/colinjfw/jbld/pkg/run"
	"github.com/stretchr/testify/require"
)

func TestCompiler(t *testing.T) {
	os.RemoveAll("./dist")
	cwd, _ := os.Getwd()

	err := run.Run("../../lib/host.js", cwd+"/config.jbld.js")
	require.NoError(t, err)
}
