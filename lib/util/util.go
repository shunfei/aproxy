package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"aproxy/lib/crypto/bcrypt"
)

func DecodeJsonBody(r io.ReadCloser, to interface{}) error {
	defer r.Close()
	err := json.NewDecoder(r).Decode(to)
	return err
}

// return json string result
// WriteJson(obj) or WriteJson(obj, "text/html")
func WriteJson(w http.ResponseWriter, data interface{}, contentType ...string) {
	var ct string
	if len(contentType) == 1 {
		ct = contentType[0]
	} else {
		ct = "application/json"
	}
	w.Header().Set("Content-Type", ct)
	s := ""
	b, err := json.Marshal(data)
	if err != nil {
		s = `{success:false, message:"json.Marshal error"}`
	} else {
		s = string(b)
	}
	fmt.Fprint(w, s)
}

// match regexp with string, and return a named group map
// Example:
//   regexp: "(?P<name>[A-Za-z]+)-(?P<age>\\d+)"
//   string: "CGC-30"
//   return: map[string]string{ "name":"CGC", "age":"30" }
func NamedRegexpGroup(str string, reg *regexp.Regexp) (ng map[string]string, matched bool) {
	rst := reg.FindStringSubmatch(str)
	//fmt.Printf("%s => %s => %s\n\n", reg, str, rst)
	if len(rst) < 1 {
		return
	}
	ng = make(map[string]string)
	lenRst := len(rst)
	sn := reg.SubexpNames()
	for k, v := range sn {
		// SubexpNames contain the none named group,
		// so must filter v == ""
		if k == 0 || v == "" {
			continue
		}
		if k+1 > lenRst {
			break
		}
		ng[v] = rst[k]
	}
	matched = true
	return
}

func CryptPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

func CompareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
