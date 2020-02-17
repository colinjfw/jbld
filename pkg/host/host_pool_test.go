package host

import (
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHostPool(t *testing.T) {
	cwd, _ := os.Getwd()
	h := NewHostPool(3, cwd+"/../../lib/host.js", cwd+"/testdata/config.jbld.js")
	defer h.Close()

	wg := sync.WaitGroup{}
	wg.Add(5)

	for i := 0; i < 5; i++ {
		go func() {
			resp := PingResponse{}
			err := h.Run("ping", PingRequest{Version: "unreleased"}, &resp)
			require.NoError(t, err)
			require.Equal(t, PingResponse{Version: "unreleased"}, resp)
			wg.Done()
		}()
	}
	wg.Wait()
}
