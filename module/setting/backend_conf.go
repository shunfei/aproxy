package setting

import (
	"errors"
	// "log"
	"net/http"
	"strings"
	"time"

	"aproxy/lib/rfweb"
	"aproxy/lib/util"
	bkconf "aproxy/module/backend_conf"
	"aproxy/module/proxy"
)

type BackendConfResource struct {
	BaseResource
}

func (self *BackendConfResource) Get(ctx *rfweb.Context) {
	res := RespData{}
	hostname := ctx.Get("hostname")
	if hostname == "all" {
		bcs, err := bkconf.GetAll()
		if err != nil {
			res.Error = err.Error()
		} else {
			res.Success = true
			res.Data = bcs
		}
	} else {
		bc, err := bkconf.Get(hostname)
		if err != nil {
			res.Error = err.Error()
		} else {
			res.Success = true
			res.Data = bc
		}
	}

	util.WriteJson(ctx.W, res)
}

// add new backend config
func (self *BackendConfResource) Post(ctx *rfweb.Context) {
	res := RespData{}
	bc, err := getBackendConfFromBody(ctx.R)
	if err != nil {
		res.Error = err.Error()
	} else {
		err = bkconf.Insert(bc)
		if err != nil {
			res.Error = err.Error()
		} else {
			res.Data = bc
			res.Success = true
			proxy.RemoveBackendConfCache()
		}
	}
	util.WriteJson(ctx.W, res)
}

// update backend config
func (self *BackendConfResource) Put(ctx *rfweb.Context) {
	res := RespData{}
	bc, err := getBackendConfFromBody(ctx.R)
	if err != nil {
		res.Error = err.Error()
	} else {
		err = bkconf.Update(bc.Id, bc)
		if err != nil {
			res.Error = err.Error()
		} else {
			res.Data = bc
			res.Success = true
			proxy.RemoveBackendConfCache()
		}
	}
	util.WriteJson(ctx.W, res)
}

// delete role
func (self *BackendConfResource) Delete(ctx *rfweb.Context) {
	res := RespData{}
	id := ctx.Get("id")
	if len(id) < 1 {
		res.Error = "no id"
	} else {
		err := bkconf.Delete(id)
		if err == nil {
			res.Success = true
			res.Data = id
		} else {
			res.Error = err.Error()
		}
	}
	util.WriteJson(ctx.W, res)
}

func getBackendConfFromBody(r *http.Request) (bkconf.BackendConf, error) {
	bc := bkconf.BackendConf{}
	err := util.DecodeJsonBody(r.Body, &bc)
	if err != nil {
		return bc, err
	}
	bc.HostName = strings.ToLower(strings.TrimSpace(bc.HostName))
	if bc.HostName == "" {
		return bc, errors.New("[hostname] must not empty.")
	}
	upstreams := []string{}
	for _, upstream := range bc.UpStreams {
		upstream = strings.TrimSpace(upstream)
		if upstream != "" {
			upstreams = append(upstreams, upstream)
		}
	}
	if len(upstreams) > 0 {
		bc.UpStreams = upstreams
	} else {
		return bc, errors.New("[upstreams] must not empty.")
	}
	bc.CreatedTime = time.Now()
	bc.UpdatedTime = bc.CreatedTime
	return bc, nil
}
