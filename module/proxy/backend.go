package proxy

import (
	"sync"

	bkconf "aproxy/module/backend_conf"
	"github.com/mailgun/oxy/forward"
	"github.com/mailgun/oxy/roundrobin"
)

type Backend struct {
	Conf bkconf.BackendConf
	Fwd  *forward.Forwarder     `json:"-"`
	Lb   *roundrobin.RoundRobin `json:"-"`
}
type Backends struct {
	sync.RWMutex
	Backends map[string]Backend
}

var backends Backends

func init() {
	backends.Backends = map[string]Backend{}
}
