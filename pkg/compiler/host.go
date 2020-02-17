package compiler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
	"sync"
)

// HostResponse returns a response from the host.
type HostResponse struct {
	Type    string   `json:"type"`
	Imports []Import `json:"imports"`
}

// Host implements an extension host.
type Host interface {
	Close() error
	Run(Source) (HostResponse, error)
}

// NewHost initializes a new host.
func NewHost(c Config) Host {
	return &host{
		config: c,
	}
}

type host struct {
	config Config
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

	arg, err := json.Marshal(h.config)
	if err != nil {
		return err
	}

	cmd := exec.Command("node", h.config.HostJS, string(arg))
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
func (h *host) Run(s Source) (HostResponse, error) {
	h.lock.Lock()
	defer h.lock.Unlock()

	if err := h.open(); err != nil {
		return HostResponse{}, err
	}

	data, err := json.Marshal(s)
	if err != nil {
		return HostResponse{}, err
	}

	// log.Printf("compiler: host - wrote %s", string(data))

	_, err = h.stdin.Write(append(data, '\n'))
	if err != nil {
		h.close()
		return HostResponse{}, err
	}

	respBytes, _, err := h.stdout.ReadLine()
	if err != nil {
		h.close()
		return HostResponse{}, err
	}

	// log.Printf("compiler: host - received %s", string(respBytes))

	resp := struct {
		HostResponse
		Err string `json:"err"`
	}{}
	err = json.Unmarshal(respBytes, &resp)
	if err != nil {
		return HostResponse{}, err
	}
	if resp.Err != "" {
		return HostResponse{}, fmt.Errorf("host: response err: %s", resp.Err)
	}
	return resp.HostResponse, nil
}
