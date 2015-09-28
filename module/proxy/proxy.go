package proxy

import (
	// "log"
	"net/http"
	"strings"

	"github.com/mailgun/oxy/forward"
	"github.com/mailgun/oxy/roundrobin"
	"github.com/mailgun/oxy/testutils"

	"aproxy/lib/rfweb"
	"aproxy/module/auth"
	"aproxy/module/auth/login"
	bkconf "aproxy/module/backend_conf"
	"aproxy/module/constant"
)

func Proxy(w http.ResponseWriter, r *http.Request) {
	if b, ok := getBackend(r); ok {
		ctx := rfweb.NewContext(w, r)
		status := auth.CheckPermission(b.Conf.AuthType, ctx)
		if status == constant.PERMISSION_STATUS_OK {
			b.Lb.ServeHTTP(w, r)
		} else if status == constant.PERMISSION_STATUS_NEED_LOGIN {
			login.RedirectToLogin(w, r)
		} else if status == constant.PERMISSION_STATUS_NO_PERMISSION {
			http.Error(w, "no permission", http.StatusForbidden)
		}

	} else {
		http.NotFound(w, r)
	}
}

func getBackend(r *http.Request) (Backend, bool) {
	var b = Backend{}
	host := strings.ToLower(r.Host)
	if b, ok := backends.Backends[host]; ok {
		return b, true
	}
	if bc, ok := getBackendConf(host); ok {
		b.Conf = bc
		b.Fwd, _ = forward.New(forward.PassHostHeader(true))
		b.Lb, _ = roundrobin.New(b.Fwd)
		for _, upstream := range bc.UpStreams {
			b.Lb.UpsertServer(testutils.ParseURI(upstream))
		}
		backends.Lock()
		backends.Backends[host] = b
		backends.Unlock()
		return b, true
	}
	return b, false
}

func getBackendConf(hostname string) (bkconf.BackendConf, bool) {
	bc, err := bkconf.Get(hostname)
	if err == nil {
		return bc, true
	}
	return bc, false
}

func RemoveBackendConfCache() {
	backends.Lock()
	defer backends.Unlock()
	backends.Backends = map[string]Backend{}
}
