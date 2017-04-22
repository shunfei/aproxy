package login

import (
	"aproxy/lib/rfweb"
	"aproxy/lib/util"
	"aproxy/module/auth"
	"aproxy/module/constant"
	"aproxy/module/oauth"
	"io"
	"net/http"
)

type OauthListResource struct {
	rfweb.BaseResource
}

func (self OauthListResource) Get(ctx *rfweb.Context) {
	res := RespData{}

	providers := oauth.GetProviderNameList()
	res.Success = true
	res.Data = providers

	util.WriteJson(ctx.W, res)
}

type OauthLoginResource struct {
	rfweb.BaseResource
}

func (self OauthLoginResource) Get(ctx *rfweb.Context) {
	returnurl := ctx.Get("returnurl")
	if returnurl != "" {
		ctx.Session().Set("returnurl", returnurl)
	}
	providerName := ctx.Get("provider")
	provider := oauth.GetOauther(providerName)
	if provider == nil {
		http.Error(ctx.W, "Can't find oauth provider.", http.StatusForbidden)
		return
	}
	err := provider.Login(providerName, ctx.W, ctx.R)
	if err != nil {
		http.Error(ctx.W, "oAuth login faild: "+err.Error(), http.StatusInternalServerError)
	}
}

type OauthCallbackResource struct {
	rfweb.BaseResource
}

func (self OauthCallbackResource) Get(ctx *rfweb.Context) {
	providerName := ctx.Get("provider")
	provider := oauth.GetOauther(providerName)
	if provider == nil {
		http.Error(ctx.W, "Can't find oauth provider.", http.StatusForbidden)
		return
	}
	email, err := provider.Callback(providerName, ctx.W, ctx.R)
	if err != nil {
		http.Error(ctx.W, "oAuth login faild: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if email == "" {
		http.Error(ctx.W, "oAuth login get email faild.", http.StatusInternalServerError)
		return
	}
	user := &auth.User{}
	user.Name = email
	user.Email = email
	user.Id = email

	session := ctx.Session()
	err = session.SetStuct(constant.SS_KEY_USER, user)
	if err != nil {
		http.Error(ctx.W, "set user info to session faild: "+err.Error(), http.StatusInternalServerError)
		return
	}

	returnurl, _ := session.Get("returnurl")
	if returnurl != "" {
		ctx.W.Header().Set("Location", returnurl)
		ctx.W.WriteHeader(http.StatusFound)
	} else {
		io.WriteString(ctx.W, "Login success with "+email)
	}
}
