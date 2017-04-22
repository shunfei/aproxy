package oauth

import (
	"net/http"
	"strings"
)

var providers map[string]Oauther
var providersNames []string

func init() {
	providers = map[string]Oauther{}
	providersNames = []string{}
}

type Oauther interface {
	// return providers name list
	Providers() []string
	Login(providerName string, w http.ResponseWriter, r *http.Request) error
	Callback(providerName string, w http.ResponseWriter, r *http.Request) (email string, err error)
}

func Register(o Oauther) {
	ps := o.Providers()
	if ps == nil || len(ps) < 1 {
		return
	}
	for _, p := range ps {
		providers[strings.ToLower(p)] = o
		providersNames = append(providersNames, p)
	}
}

func GetOauther(providerName string) Oauther {
	providerName = strings.ToLower(providerName)
	return providers[providerName]
}

func GetProviderNameList() []string {
	return providersNames
}
