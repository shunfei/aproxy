package login

import (
	"net/http"
	"net/url"
	"path"
	"strings"

	"aproxy/lib/rfweb"
	"aproxy/lib/util"
	"aproxy/module/auth"
	"aproxy/module/constant"
)

var loginUrl = ""
var loginHost = ""

type RespData struct {
	Success bool        `json:"success"`
	Error   string      `json:"error"`
	Data    interface{} `json:"data"`
}

type LoginResource struct {
	rfweb.BaseResource
}

func (self LoginResource) Post(ctx *rfweb.Context) {
	res := RespData{}
	email := strings.ToLower(ctx.Get("email"))
	email = strings.TrimSpace(email)
	pwd := strings.TrimSpace(ctx.Get("pwd"))
	// remember := ctx.Get("remember")

	user, err := auth.LoginUser(email, pwd)
	if err != nil {
		res.Error = err.Error()
	} else {
		res.Success = true
		// res.Data = ctx.Get("returnurl")
		session := ctx.Session()
		session.SetStuct(constant.SS_KEY_USER, user)
	}
	util.WriteJson(ctx.W, res)
}

func redirectToLogin(w http.ResponseWriter, r *http.Request, hasReturnUrl bool) {
	scheme := "http://"
	if r.TLS != nil {
		scheme = "https://"
	}
	returnurl := ""
	if hasReturnUrl {
		returnurl = scheme + r.Host + r.RequestURI
	}
	tourl := loginUrl + "?returnurl=" + url.QueryEscape(returnurl)
	w.Header().Set("Location", tourl)
	w.WriteHeader(http.StatusFound)
}

func RedirectToLogin(w http.ResponseWriter, r *http.Request) {
	redirectToLogin(w, r, true)
}

func InitLoginServer(host, urlPrefix string) {
	if len(host) > 0 && string(host[len(host)-1]) == "/" {
		host = host[:len(host)-1]
	}
	loginUrl = host + path.Join("/", urlPrefix, "login.html")
	r, err := url.Parse(loginUrl)
	if err == nil {
		loginHost = r.Host
	}
}
