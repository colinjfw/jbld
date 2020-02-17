package host

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHost(t *testing.T) {
	cwd, _ := os.Getwd()
	h := NewHost(cwd+"/../../lib/host.js", cwd+"/testdata/config.jbld.js")
	defer h.Close()

	resp := PingResponse{}
	err := h.Run("ping", PingRequest{Version: "unreleased"}, &resp)
	require.NoError(t, err)
	require.Equal(t, PingResponse{Version: "unreleased"}, resp)
}
