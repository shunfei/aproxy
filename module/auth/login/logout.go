package login

import (
	"aproxy/lib/rfweb"
)

type LogoutResource struct {
	rfweb.BaseResource
}

func (self *LogoutResource) Get(ctx *rfweb.Context) {
	session := ctx.Session()
	session.Clear(ctx.W)
	redirectToLogin(ctx.W, ctx.R, false)
}
