package rfweb

import (
	"net/http"

	"aproxy/lib/rfweb/session"
)

type Context struct {
	W         http.ResponseWriter
	R         *http.Request
	UrlParams map[string]string
	Data      map[interface{}]interface{}

	session *session.Session
}

func (self Context) Get(key string) string {
	res := self.UrlParams[key]
	if res == "" {
		res = self.R.FormValue(key)
	}
	return res
}

func (self Context) Session() *session.Session {
	if self.session == nil {
		self.session, _ = session.GetSession(self.W, self.R)
	}
	return self.session
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	ctx := &Context{}
	ctx.W = w
	ctx.R = r
	ctx.UrlParams = map[string]string{}
	ctx.Data = make(map[interface{}]interface{})
	return ctx
}

type Resourcer interface {
	Get(*Context)
	Post(*Context)
	Put(*Context)
	Delete(*Context)
	Head(*Context)
	Patch(*Context)
	Options(*Context)

	// before handle request,
	// will execute this first.
	// if return false, will stop handle the request,
	// and end the http request.
	OnHandleBegin(*Context) bool
	OnHandleEnd(*Context)
}

type BaseResource struct {
}

func (self *BaseResource) Get(ctx *Context) {
	http.NotFound(ctx.W, ctx.R)
}

func (self *BaseResource) Put(ctx *Context) {
	http.NotFound(ctx.W, ctx.R)
}

func (self *BaseResource) Post(ctx *Context) {
	http.NotFound(ctx.W, ctx.R)
}

func (self *BaseResource) Delete(ctx *Context) {
	http.NotFound(ctx.W, ctx.R)
}

func (self *BaseResource) Head(ctx *Context) {
	http.NotFound(ctx.W, ctx.R)
}

func (self *BaseResource) Patch(ctx *Context) {
	http.NotFound(ctx.W, ctx.R)
}

func (self *BaseResource) Options(ctx *Context) {
	http.NotFound(ctx.W, ctx.R)
}

func (self *BaseResource) OnHandleBegin(ctx *Context) bool {
	return true
}
func (self *BaseResource) OnHandleEnd(ctx *Context) {

}
