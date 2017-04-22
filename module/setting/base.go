package setting

import (
	"net/http"

	"aproxy/lib/rfweb"
	"aproxy/lib/util"
	"aproxy/module/auth"
	"aproxy/module/auth/login"
)

var (
	fileServer http.Handler

	inited          = false
	staticFileDir   = "./"
	AproxyUrlPrefix = "/-_-aproxy-_-/"
)

type RespData struct {
	Success bool        `json:"success"`
	Error   string      `json:"error"`
	Data    interface{} `json:"data"`
}

type BaseResource struct {
	rfweb.BaseResource
}

// check permission
func (self *BaseResource) OnHandleBegin(ctx *rfweb.Context) bool {
	user := auth.GetLoginedUser(ctx)
	errMsg := ""
	if user == nil || user.Email == "" {
		errMsg = "please login first."
	} else {
		authority, err := auth.GetAuthorityByEmail(user.Email)
		if err != nil {
			errMsg = "can't get authority, error: " + err.Error()
		} else if authority == nil || authority.AdminLevel < 10 {
			errMsg = "you don't has permission."
		}
	}
	if errMsg != "" {
		isXHR := ctx.R.Header.Get("X-Requested-With") == "XMLHttpRequest"
		if isXHR {
			res := RespData{
				Error: errMsg,
			}
			util.WriteJson(ctx.W, res)
		} else {
			http.Error(ctx.W, errMsg, http.StatusForbidden)
		}
		return false
	}
	return true
}

func (self *BaseResource) OnHandleEnd(ctx *rfweb.Context) {

}

func InitSettingServer(webDir, aproxyUrlPrefix string) {
	inited = true
	staticFileDir = webDir
	if aproxyUrlPrefix != "" {
		AproxyUrlPrefix = aproxyUrlPrefix
	}
	fileServer = http.FileServer(http.Dir(staticFileDir))
}

func StaticServer(w http.ResponseWriter, r *http.Request) {
	// check permission
	if r.RequestURI == AproxyUrlPrefix ||
		r.RequestURI == AproxyUrlPrefix+"index.html" {
		ctx := rfweb.NewContext(w, r)
		user := auth.GetLoginedUser(ctx)
		errMsg := ""
		if user == nil {
			login.RedirectToLogin(w, r)
			return
		} else {
			authority, err := auth.GetAuthorityByEmail(user.Email)
			if err != nil {
				errMsg = "can't get authority, error: " + err.Error()
			} else if authority == nil || authority.AdminLevel < 10 {
				errMsg = "you don't has permission."
			}
		}
		if errMsg != "" {
			http.Error(ctx.W, errMsg, http.StatusForbidden)
			return
		}
	}

	http.StripPrefix(AproxyUrlPrefix,
		fileServer).ServeHTTP(w, r)
}

func NewApiApp() *rfweb.App {
	app := rfweb.NewApp(AproxyUrlPrefix + "api/")
	app.Resource("backends/{hostname}", &BackendConfResource{})
	app.Resource("role/{id}", &RoleResource{})
	app.Resource("authority/{id}", &AuthorityResource{})

	app.Resource("users/{email}", &UserResource{})
	app.Resource("user/login", &login.LoginResource{})
	app.Resource("user/logout", &login.LogoutResource{})

	app.Resource("oauth/list", &login.OauthListResource{})
	app.Resource("oauth/login", &login.OauthLoginResource{})
	app.Resource("oauth/callback", &login.OauthCallbackResource{})

	return app
}
