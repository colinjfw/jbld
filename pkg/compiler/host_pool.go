package compiler

// NewHostPool initializes a pool of w hosts.
func NewHostPool(c Config) Host {
	hp := &hostPool{config: c}
	hp.create()
	return hp
}

type hostPool struct {
	config Config
	hosts  chan Host
}

func (h *hostPool) create() {
	h.hosts = make(chan Host, h.config.Workers)
	for i := 0; i < h.config.Workers; i++ {
		h.hosts <- NewHost(h.config)
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
