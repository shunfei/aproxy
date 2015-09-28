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

// with auth.User,
// need marshal exclude `Pwd` field,
// but unmarshal include `Pwd` field.
// but the json pkg can't do that,
// so create this struct
type FormUser struct {
	Id    string
	Email string
	Pwd   string
}

type UserResource struct {
	BaseResource
}

func (self *UserResource) Get(ctx *rfweb.Context) {
	res := RespData{}
	email := ctx.Get("email")
	if email == "all" {
		users, err := auth.GetAllUsers()
		if err != nil {
			res.Error = err.Error()
		} else {
			res.Success = true
			res.Data = users
		}
	} else {
		user, err := auth.GetUserByEmail(email)
		if err != nil {
			res.Error = err.Error()
		} else {
			res.Success = true
			res.Data = user
		}
	}

	util.WriteJson(ctx.W, res)
}

// add new authority
func (self *UserResource) Post(ctx *rfweb.Context) {
	res := RespData{}
	user, err := getUserFromBody(ctx.R, true)
	if err != nil {
		res.Error = err.Error()
	} else {
		err = auth.InsertUser(*user)
		if err != nil {
			res.Error = err.Error()
		} else {
			res.Data = user
			res.Success = true
			// proxy.RemoveBackendConfCache()
		}
	}
	util.WriteJson(ctx.W, res)
}

// update user
func (self *UserResource) Put(ctx *rfweb.Context) {
	res := RespData{}
	user, err := getUserFromBody(ctx.R, false)
	if err != nil {
		res.Error = err.Error()
	} else {
		err = auth.UpdateUser(user.Id, *user)
		if err != nil {
			res.Error = err.Error()
		} else {
			res.Data = user
			res.Success = true
			// proxy.RemoveBackendConfCache()
		}
	}
	util.WriteJson(ctx.W, res)
}

// delete user
func (self *UserResource) Delete(ctx *rfweb.Context) {
	res := RespData{}
	id := ctx.Get("id")
	if len(id) < 1 {
		res.Error = "no id"
	} else {
		err := auth.DeleteUser(id)
		if err == nil {
			res.Success = true
			res.Data = id
		} else {
			res.Error = err.Error()
		}
	}
	util.WriteJson(ctx.W, res)
}

func getUserFromBody(r *http.Request, needPwd bool) (*auth.User, error) {
	fuser := FormUser{}
	user := &auth.User{}
	err := util.DecodeJsonBody(r.Body, &fuser)
	if err != nil {
		return user, err
	}
	user.Id = fuser.Id
	user.Email = fuser.Email
	user.Pwd = fuser.Pwd

	user.Email = strings.TrimSpace(user.Email)
	if len(user.Email) < 1 {
		return user, errors.New("[Email] must not empty.")
	}
	user.Pwd = strings.TrimSpace(user.Pwd)
	if needPwd && len(user.Pwd) < 1 {
		return user, errors.New("[Password] must not empty.")
	}

	user.CreatedTime = time.Now()
	user.UpdatedTime = user.CreatedTime
	return user, nil
}
