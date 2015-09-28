package auth

import (
	"log"

	"aproxy/lib/rfweb"
	"aproxy/module/constant"
)

// check permission, return error code:
//     0: ok
//     1: need login
//     2: do not has permission
func CheckPermission(authType int, ctx *rfweb.Context) int {
	if authType < constant.AUTH_TYPE_PUBLIC {
		return constant.PERMISSION_STATUS_OK
	}
	user := GetLoginedUser(ctx)
	if user == nil {
		return constant.PERMISSION_STATUS_NEED_LOGIN
	}
	if authType == constant.AUTH_TYPE_LOGIN {
		return constant.PERMISSION_STATUS_OK
	}
	// TODO: need cache
	authority, err := GetAuthorityByEmail(user.Email)
	if err != nil || authority == nil {
		return constant.PERMISSION_STATUS_NO_PERMISSION
	}
	authority.Init()
	rurl := ctx.R.Host + ctx.R.RequestURI
	if !authority.HasPermission(rurl) {
		return constant.PERMISSION_STATUS_NO_PERMISSION
	}
	return constant.PERMISSION_STATUS_OK
}

func GetLoginedUser(ctx *rfweb.Context) *User {
	r, ok := ctx.Data[constant.CTX_KEY_USER]
	if ok && r != nil {
		return r.(*User)
	}
	session := ctx.Session()
	var user User
	err := session.GetStuct(constant.SS_KEY_USER, &user)
	if err != nil {
		log.Println("[ERROR]", err)
		return nil
	}
	if len(user.Id) > 0 {
		ctx.Data[constant.CTX_KEY_USER] = &user
	} else {
		return nil
	}
	return &user
}
