package setting

import (
	"errors"
	// "log"
	"net/http"
	"strings"
	"time"

	"aproxy/lib/rfweb"
	"aproxy/lib/util"
	"aproxy/module/auth"
)

type RoleResource struct {
	BaseResource
}

func (self *RoleResource) Get(ctx *rfweb.Context) {
	res := RespData{}
	id := ctx.Get("id")
	if id == "all" {
		roles, err := auth.GetAllRole()
		if err != nil {
			res.Error = err.Error()
		} else {
			res.Success = true
			res.Data = roles
		}
	} else {
		role, err := auth.GetRoleByID(id)
		if err != nil {
			res.Error = err.Error()
		} else {
			res.Success = true
			res.Data = role
		}
	}

	util.WriteJson(ctx.W, res)
}

// add new role
func (self *RoleResource) Post(ctx *rfweb.Context) {
	res := RespData{}
	role, err := getRoleFromBody(ctx.R)
	if err != nil {
		res.Error = err.Error()
	} else {
		err = auth.InsertRole(role)
		if err != nil {
			res.Error = err.Error()
		} else {
			res.Data = role
			res.Success = true
			// proxy.RemoveBackendConfCache()
		}
	}
	util.WriteJson(ctx.W, res)
}

// update role config
func (self *RoleResource) Put(ctx *rfweb.Context) {
	res := RespData{}
	role, err := getRoleFromBody(ctx.R)
	if err != nil {
		res.Error = err.Error()
	} else {
		err = auth.UpdateRole(role.Id, role)
		if err != nil {
			res.Error = err.Error()
		} else {
			res.Data = role
			res.Success = true
			// proxy.RemoveBackendConfCache()
		}
	}
	util.WriteJson(ctx.W, res)
}

// delete role
func (self *RoleResource) Delete(ctx *rfweb.Context) {
	res := RespData{}
	id := ctx.Get("id")
	if len(id) < 1 {
		res.Error = "no id"
	} else {
		err := auth.DeleteRole(id)
		if err == nil {
			res.Success = true
			res.Data = id
		} else {
			res.Error = err.Error()
		}
	}
	util.WriteJson(ctx.W, res)
}

func getRoleFromBody(r *http.Request) (*auth.Role, error) {
	role := &auth.Role{}
	err := util.DecodeJsonBody(r.Body, &role)
	if err != nil {
		return role, err
	}

	role.Name = strings.TrimSpace(role.Name)
	if len(role.Name) < 1 {
		return role, errors.New("[Name] must not empty.")
	}

	allow := []string{}
	for _, a := range role.Allow {
		a = strings.TrimSpace(a)
		if a != "" {
			allow = append(allow, a)
		}
	}
	if len(allow) > 0 {
		role.Allow = allow
	}
	deny := []string{}
	for _, d := range role.Deny {
		d = strings.TrimSpace(d)
		if d != "" {
			deny = append(deny, d)
		}
	}
	if len(deny) > 0 {
		role.Deny = deny
	}
	if len(deny) < 1 && len(allow) < 1 {
		return role, errors.New("[Deny] and [Allow] can't all be empty.")
	}

	role.CreatedTime = time.Now()
	role.UpdatedTime = role.CreatedTime
	return role, nil
}
