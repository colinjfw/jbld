package host

// NewHostPool initializes a pool of w hosts.
func NewHostPool(w int, js, arg string) Host {
	hp := &hostPool{workers: w, hostJS: js, arg: arg}
	hp.create()
	return hp
}

type hostPool struct {
	workers int
	hostJS  string
	arg     string
	hosts   chan Host
}

func (h *hostPool) create() {
	h.hosts = make(chan Host, h.workers)
	for i := 0; i < h.workers; i++ {
		h.hosts <- NewHost(h.hostJS, h.arg)
	}
}

func (h *hostPool) Close() error {
	close(h.hosts)
	var err error
	for in := range h.hosts {
		inerr := in.Close()
		if inerr != nil {
			err = inerr
		}
	}
	return err
}

func (h *hostPool) Run(m string, i, o interface{}) error {
	in := <-h.hosts
	err := in.Run(m, i, o)
	h.hosts <- in
	return err
}
