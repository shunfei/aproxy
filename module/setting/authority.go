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

type AuthorityResource struct {
	BaseResource
}

func (self *AuthorityResource) Get(ctx *rfweb.Context) {
	res := RespData{}
	id := ctx.Get("id")
	if id == "all" {
		authoritys, err := auth.GetAllAuthority()
		if err != nil {
			res.Error = err.Error()
		} else {
			res.Success = true
			res.Data = authoritys
		}
	} else if id != "" {
		authority, err := auth.GetAuthorityByID(id)
		if err != nil {
			res.Error = err.Error()
		} else {
			res.Success = true
			res.Data = authority
		}
	} else {
		email := ctx.Get("email")
		if email != "" {
			authority, err := auth.GetAuthorityByEmail(email)
			if err != nil {
				res.Error = err.Error()
			} else {
				res.Success = true
				res.Data = authority
			}
		}
	}

	util.WriteJson(ctx.W, res)
}

// add new authority
func (self *AuthorityResource) Post(ctx *rfweb.Context) {
	res := RespData{}
	authority, err := getAuthorityFromBody(ctx.R)
	if err != nil {
		res.Error = err.Error()
	} else {
		err = auth.InsertAuthority(authority)
		if err != nil {
			res.Error = err.Error()
		} else {
			res.Data = authority
			res.Success = true
			// proxy.RemoveBackendConfCache()
		}
	}
	util.WriteJson(ctx.W, res)
}

// update authority
func (self *AuthorityResource) Put(ctx *rfweb.Context) {
	res := RespData{}
	authority, err := getAuthorityFromBody(ctx.R)
	if err != nil {
		res.Error = err.Error()
	} else {
		err = auth.UpdateAuthority(authority.Id, authority)
		if err != nil {
			res.Error = err.Error()
		} else {
			res.Data = authority
			res.Success = true
			// proxy.RemoveBackendConfCache()
		}
	}
	util.WriteJson(ctx.W, res)
}

// delete authority
func (self *AuthorityResource) Delete(ctx *rfweb.Context) {
	res := RespData{}
	id := ctx.Get("id")
	if len(id) < 1 {
		res.Error = "no id"
	} else {
		err := auth.DeleteAuthority(id)
		if err == nil {
			res.Success = true
			res.Data = id
		} else {
			res.Error = err.Error()
		}
	}
	util.WriteJson(ctx.W, res)
}

func getAuthorityFromBody(r *http.Request) (*auth.Authority, error) {
	authority := &auth.Authority{}
	err := util.DecodeJsonBody(r.Body, &authority)
	if err != nil {
		return authority, err
	}

	authority.Email = strings.TrimSpace(authority.Email)
	if len(authority.Email) < 1 {
		return authority, errors.New("[Email] must not empty.")
	}

	allow := []string{}
	for _, a := range authority.Allow {
		a = strings.TrimSpace(a)
		if a != "" {
			allow = append(allow, a)
		}
	}
	if len(allow) > 0 {
		authority.Allow = allow
	}
	deny := []string{}
	for _, d := range authority.Deny {
		d = strings.TrimSpace(d)
		if d != "" {
			deny = append(deny, d)
		}
	}
	if len(deny) > 0 {
		authority.Deny = deny
	}
	if len(deny) < 1 && len(allow) < 1 && len(authority.Roles) < 1 {
		return authority, errors.New("[Deny], [Allow] and [Roles] can't all be empty.")
	}

	authority.CreatedTime = time.Now()
	authority.UpdatedTime = authority.CreatedTime
	return authority, nil
}
