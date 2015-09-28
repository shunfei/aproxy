package rfweb

import (
	"testing"

	. "gopkg.in/bluesuncorp/assert.v1"
)

var rs = BaseResource{}

var r1 = &Route{
	Pattern:  "/users/{id}",
	Resource: rs,
}

var r2 = &Route{
	Pattern:  "/backends/{hostname}",
	Resource: rs,
}

var r3 = &Route{
	Pattern:  "/pages/{name}",
	Resource: rs,
}

func initRoute() {
	r1.Init()
	r2.Init()
	r3.Init()
}

var routerTestData = []struct {
	Route   *Route
	Url     string
	Matched bool
	Params  map[string]string
}{
	{
		r1,
		"/users/3",
		true,
		map[string]string{
			"id": "3",
		},
	},
	{
		r1,
		"/post/3",
		false,
		nil,
	},
	{
		r2,
		"/backends/abc.com",
		true,
		map[string]string{
			"hostname": "abc.com",
		},
	},
	{
		r3,
		"/pages/about",
		true,
		map[string]string{
			"name": "about",
		},
	},
}

func TestRouteInit(t *testing.T) {
	defer func() {
		if x := recover(); x != nil {
			t.Errorf("Must no panic, but got panic: \n\t%s", x)
		}
	}()
	route := new(Route)
	route.Pattern = "/users/{id}"
	route.Resource = rs
	route.Init()

	params, ok := route.Match("/users/2")
	Equal(t, ok, true)
	Equal(t, params["id"], "2")
}

func TestRouteMatch(t *testing.T) {
	initRoute()

	for _, td := range routerTestData {
		params, ok := td.Route.Match(td.Url)
		if ok != td.Matched {
			t.Errorf("url [%s] not match", td.Url)
		}
		Equal(t, ok, td.Matched)
		if td.Matched && ok {
			for k, p := range td.Params {
				Equal(t, params[k], p)
			}
		}
	}
}

func TestRouteTable(t *testing.T) {
	var rt *RouteTable
	rt = &RouteTable{Routes: make([]*Route, 0, 10)}

	rt.Map("/post/{id}", rs)
	rt.AddRoute(r2)
	rt.AddRoute(r1)

	route, params, ok := rt.Match("/p")
	Equal(t, ok, false)
	NotEqual(t, route, nil)

	route, params, ok = rt.Match("/post/save")
	Equal(t, ok, true)
	Equal(t, params["id"], "save")
	NotEqual(t, route, nil)

	route, params, ok = rt.Match("/users/2")
	Equal(t, ok, true)
	Equal(t, params["id"], "2")
	Equal(t, route, r1)
}
