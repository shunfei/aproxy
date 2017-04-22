package github

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"path"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	githubAuthorizeUrl = "https://github.com/login/oauth/authorize"
	githubTokenUrl     = "https://github.com/login/oauth/access_token"
)

var (
	oauthCfg    *oauth2.Config
	redirectUrl string
	scopes      []string
)

func init() {
	scopes = []string{"user:email"}
}

type GithubOauther struct {
}

func (self GithubOauther) Providers() []string {
	return []string{
		"Github",
	}
}

func (self GithubOauther) Login(providerName string, w http.ResponseWriter, r *http.Request) error {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	authUrl := oauthCfg.AuthCodeURL(state)
	// redirect
	http.Redirect(w, r, authUrl, http.StatusFound)
	return nil
}

func (self GithubOauther) Callback(providerName string, w http.ResponseWriter, r *http.Request) (string, error) {
	tkn, err := oauthCfg.Exchange(oauth2.NoContext, r.URL.Query().Get("code"))
	if err != nil {
		return "", errors.New("there was an issue getting your token: " + err.Error())
	}

	if !tkn.Valid() {
		return "", errors.New("Github oauth retreived invalid token." + err.Error())
	}

	client := github.NewClient(oauthCfg.Client(oauth2.NoContext, tkn))
	// opt := &github.ListOptions{}
	emails, _, err := client.Users.ListEmails(oauth2.NoContext, nil)
	if err != nil {
		return "", errors.New("get github email faild: " + err.Error())
	}
	email := ""
	if len(emails) > 0 {
		email = emails[0].GetEmail()
	}

	return email, nil
}

// @host: include http://, like http://abc.com
func InitGithubOauther(urlPrefix, host, clientID, clientSecret string) {
	redirectUrl = host + path.Join(urlPrefix, "/api/oauth/callback?provider=github")
	oauthCfg = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  githubAuthorizeUrl,
			TokenURL: githubTokenUrl,
		},
		RedirectURL: redirectUrl,
		Scopes:      scopes,
	}
}
