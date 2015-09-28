package session

import (
	// "encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	// "io"
	"crypto/rand"
	"net/http"
	"net/url"
	"time"
)

var (
	sessionExp        = int64(60 * 60 * 24 * 7)
	sessionCookieName = "aproxysid"
	sessionDomain     = ""
)

type Session struct {
	sessionId string

	storage SessionStorager
}

func NewSession(sid string) *Session {
	s := &Session{}
	s.sessionId = sid
	s.storage = sessionStorage
	return s
}

func (self *Session) Get(key string) (string, error) {
	return self.storage.Get(self.sessionId, key)
}

func (self *Session) Set(key, val string, exp ...int64) error {
	iexp := sessionExp
	if len(exp) > 0 && exp[0] > 0 {
		iexp = exp[0]
	}
	return self.storage.Set(self.sessionId, key, val, iexp)
}

func (self *Session) GetStuct(key string, val interface{}) error {
	sval, err := self.Get(key)
	if err != nil {
		return err
	}
	if len(sval) < 1 {
		return nil
	}
	err = json.Unmarshal([]byte(sval), val)
	return err
}

func (self *Session) SetStuct(key string, val interface{}, exp ...int64) error {
	sval, err := json.Marshal(val)
	if err != nil {
		return err
	}
	err = self.Set(key, string(sval), exp...)
	return err
}

func (self *Session) Clear(w http.ResponseWriter) error {
	err := self.storage.Clear(self.sessionId)
	c := &http.Cookie{
		Name:    sessionCookieName,
		Expires: time.Now().Add(-10 * time.Second),
		Path:    "/",
		Domain:  sessionDomain,
	}
	http.SetCookie(w, c)
	return err
}

func NewSessionId() (string, error) {
	b := make([]byte, 20)
	n, err := rand.Read(b)
	if n != len(b) || err != nil {
		return "", fmt.Errorf("Could not successfully read from the system CSPRNG.")
	}
	return hex.EncodeToString(b), nil
}

func WriteSessionId(w http.ResponseWriter, sid string, exp int64) {
	if len(sid) == 0 {
		return
	}

	cookie := &http.Cookie{
		Name:     sessionCookieName,
		Value:    url.QueryEscape(sid),
		HttpOnly: true,
		Path:     "/",
		Domain:   sessionDomain,
	}
	if exp > 0 {
		expiration := time.Now()
		expiration = expiration.Add(time.Second * time.Duration(exp))
		cookie.Expires = expiration
	}
	http.SetCookie(w, cookie)
}

func GetSession(w http.ResponseWriter, r *http.Request) (*Session, error) {
	sid := ""
	sidCookie, err := r.Cookie(sessionCookieName)
	if err == nil {
		sid, _ = url.QueryUnescape(sidCookie.Value)
	}
	if sid == "" {
		sid, err = NewSessionId()
		if err == nil {
			WriteSessionId(w, sid, sessionExp)
		}
	}
	s := NewSession(sid)
	return s, nil
}

func InitSessionServer(domain, cookieName string, exp int64) {
	sessionDomain = domain
	sessionCookieName = cookieName
	sessionExp = exp
}
