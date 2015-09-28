package rfweb

import (
	"net/http"
	"path"
)

const (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	HEAD    = "HEAD"
	PATCH   = "PATCH"
	OPTIONS = "OPTIONS"
)

type App struct {
	UrlPrefix string

	rt RouteTable
}

func NewApp(urlPrefix string) *App {
	app := &App{}
	app.UrlPrefix = urlPrefix
	app.rt = RouteTable{}
	return app
}

func (self *App) Resource(urlPath string, resource Resourcer) {
	self.rt.Map(
		path.Join(self.UrlPrefix, urlPath),
		resource)
}

func (self *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route, params, matched := self.rt.Match(r.URL.Path)
	if !matched {
		http.NotFound(w, r)
		return
	}
	ctx := NewContext(w, r)
	ctx.UrlParams = params
	if !route.Resource.OnHandleBegin(ctx) {
		return
	}
	switch r.Method {
	case GET:
		route.Resource.Get(ctx)
	case POST:
		route.Resource.Post(ctx)
	case PUT:
		route.Resource.Put(ctx)
	case DELETE:
		route.Resource.Delete(ctx)
	case PATCH:
		route.Resource.Patch(ctx)
	case OPTIONS:
		route.Resource.Options(ctx)
	case HEAD:
		route.Resource.Head(ctx)
	default:
		http.NotFound(w, r)
	}
	route.Resource.OnHandleEnd(ctx)
}
