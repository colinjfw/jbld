package compiler

// NewHostPool initializes a pool of w hosts.
func NewHostPool(w int, js, config string) Host {
	hp := &hostPool{js: js, config: config, count: w}
	hp.create()
	return hp
}

type hostPool struct {
	js     string
	config string
	count  int
	hosts  chan Host
}

func (h *hostPool) create() {
	h.hosts = make(chan Host, h.count)
	for i := 0; i < h.count; i++ {
		h.hosts <- NewHost(h.js, h.config)
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

func (h *hostPool) Run(s Source) ([]Import, error) {
	in := <-h.hosts
	imp, err := in.Run(s)
	h.hosts <- in
	return imp, err
}
