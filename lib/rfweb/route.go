package rfweb

import (
	"fmt"
	"regexp"
	//"path"

	"aproxy/lib/util"
)

var (
	regPathParse *regexp.Regexp = regexp.MustCompile("/?\\{[\\w\\-_]+\\}") // matched like this: /{controller}
)

// Route config
//      var rt = &Route {
//          Pattern: "/users/{id}"
//      }
// and then, must init the router
//      rt.Init()
// and then, you can use it
//      rt.Match("/users/1")
//
type Route struct {
	Pattern  string // url pattern config, eg. /users/{id}
	Resource Resourcer

	rePath *regexp.Regexp
	inited bool
}

func (self *Route) Init() {
	if self.inited {
		return
	}
	if self.Pattern == "" {
		panic("Route: Pattern must be set")
	}
	if self.Resource == nil {
		panic("Route: Resource must be set")
	}

	r := regPathParse.ReplaceAllStringFunc(self.Pattern, func(s string) string {
		slash := ""
		if s[0] == '/' {
			slash = "/"
			s = s[1:]
		}
		name, reg, need := s[1:len(s)-1], "[^\\?#/]+", "?"
		if slash != "" {
			//  / => /?
			slash = slash + "?"
		}
		//(?P<name>re)
		return fmt.Sprintf("%s(?P<%s>%s)%s", slash, name, reg, need)
	})
	if r != "" && r[len(r)-1] == '/' {
		r = r + "?"
	}
	self.rePath = regexp.MustCompile("^" + r + "$")
	self.inited = true
}

func (self *Route) Match(url string) (params map[string]string, matched bool) {
	if !self.inited {
		self.Init()
	}

	params, matched = util.NamedRegexpGroup(url, self.rePath)

	return
}

type RouteTable struct {
	Routes []*Route
}

func (rt *RouteTable) Match(url string) (route *Route, params map[string]string, matched bool) {
	if url == "" {
		return
	}
	for _, route = range rt.Routes {
		params, matched = route.Match(url)
		if matched {
			return
		}
	}
	return
}

func (rt *RouteTable) AddRoute(route *Route) {
	route.Init()
	rt.Routes = append(rt.Routes, route)
}

// add a new route
func (rt *RouteTable) Map(url string, resource Resourcer) {

	route := &Route{
		Pattern:  url,
		Resource: resource,
	}
	route.Init()
	rt.Routes = append(rt.Routes, route)
}
