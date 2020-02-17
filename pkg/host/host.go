package host

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
	"sync"
)

// PingRequest is a valid command.
type PingRequest struct {
	Version string `json:"version"`
}

// PingResponse is a valid response.
type PingResponse struct {
	Version string `json:"version"`
}

// Host implements an extension host.
type Host interface {
	Close() error
	Run(method string, in, out interface{}) error
}

// NewHost initializes a new host.
func NewHost(js, arg string) Host {
	return &host{
		arg:    arg,
		hostJS: js,
	}
}

type host struct {
	hostJS string
	arg    string
	lock   sync.Mutex
	stdin  io.Writer
	stdout *bufio.Reader
	proc   *exec.Cmd
}

func (h *host) Close() error {
	h.lock.Lock()
	defer h.lock.Unlock()
	return h.close()
}

func (h *host) open() error {
	if h.proc != nil { // Already open.
		return nil
	}

	cmd := exec.Command("node", h.hostJS, h.arg)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	go func() {
		s := bufio.NewScanner(stderr)
		for s.Scan() {
			log.Printf("host: js/stderr: %s", s.Text())
		}
	}()

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	h.stdin = stdin
	h.stdout = bufio.NewReader(stdout)

	err = cmd.Start()
	if err != nil {
		return err
	}

	// log.Printf("host: open - node %s %s", h.js, conf)

	h.proc = cmd
	return nil
}

func (h *host) close() error {
	if h.proc == nil {
		return nil
	}
	h.proc.Process.Kill()
	err := h.proc.Wait()

	h.stdin = nil
	h.stdout = nil
	h.proc = nil
	return err
}

// Run implements the Host interface.
func (h *host) Run(method string, in, out interface{}) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	if err := h.open(); err != nil {
		return err
	}

	reqBody, err := json.Marshal(in)
	if err != nil {
		return err
	}
	req := struct {
		Req    json.RawMessage `json:"req"`
		Method string          `json:"method"`
	}{
		Req:    reqBody,
		Method: method,
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	// log.Printf("compiler: host - wrote %s", string(data))

	_, err = h.stdin.Write(append(reqBytes, '\n'))
	if err != nil {
		h.close()
		return err
	}

	respBytes, _, err := h.stdout.ReadLine()
	if err != nil {
		h.close()
		return err
	}

	// log.Printf("compiler: host - received %s", string(respBytes))

	resp := struct {
		Res json.RawMessage `json:"res"`
		Err string          `json:"err"`
	}{}
	err = json.Unmarshal(respBytes, &resp)
	if err != nil {
		return err
	}
	if resp.Err != "" {
		return fmt.Errorf("host: response err: %s", resp.Err)
	}
	return json.Unmarshal(resp.Res, out)
}
