package auth

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"aproxy/module/db"
)

const (
	C_NAME_Authority = "Authority"
)

type Authority struct {
	sync.RWMutex `bson:"-"`

	Id    string   `bson:"_id,omitempty"`
	Email string   `bson:"Email"`
	Desc  string   `bson:"Desc"` // description
	Roles []string `bson:"Roles"`
	Allow []string `bson:"Allow"`
	Deny  []string `bson:"Deny"`

	// 0: not a administrator
	// 50: system administrator
	// 99: super administrator
	// Currently used to determine whether
	// there set permissions, this field is
	// mainly reserved for future expansion.
	AdminLevel int `bson:"AdminLevel"`

	CreatedTime time.Time `bson:"CreatedTime"`
	UpdatedTime time.Time `bson:"UpdatedTime"`

	inited     bool
	rl, al, dl int
	allowReg   []*regexp.Regexp
	denyReg    []*regexp.Regexp
	roles      []*Role
}

func (self *Authority) Init() {
	self.Lock()
	defer self.Unlock()
	if self.inited {
		return
	}
	self.init()
}

func (self *Authority) init() {
	if self.allowReg == nil {
		self.allowReg = []*regexp.Regexp{}
	}
	if self.denyReg == nil {
		self.denyReg = []*regexp.Regexp{}
	}
	for _, allow := range self.Allow {
		allow = strings.Replace(allow, ".", "\\.", -1)
		allow = strings.Replace(allow, "*", ".*", -1)
		self.allowReg = append(self.allowReg,
			regexp.MustCompile(allow))
	}
	for _, deny := range self.Deny {
		deny = strings.Replace(deny, ".", "\\.", -1)
		deny = strings.Replace(deny, "*", ".*", -1)
		self.denyReg = append(self.denyReg,
			regexp.MustCompile(deny))
	}
	for _, id := range self.Roles {
		// TODO: need cache
		r, err := GetRoleByID(id)
		if err == nil && r != nil && r.Id != "" {
			// must init role
			r.Init()
			self.roles = append(self.roles, r)
		}
	}
	self.al = len(self.allowReg)
	self.dl = len(self.denyReg)
	self.rl = len(self.roles)
	return
}

// must Init() before use this.
func (self *Authority) HasPermission(url string) bool {
	for i := 0; i < self.dl; i++ {
		if self.denyReg[i].MatchString(url) {
			return false
		}
	}
	for i := 0; i < self.al; i++ {
		if self.allowReg[i].MatchString(url) {
			return true
		}
	}
	for i := 0; i < self.rl; i++ {
		if self.roles[i].HasPermission(url) {
			return true
		}
	}
	return false
}

func GetAuthorityByID(id string) (*Authority, error) {
	c := db.MDB().C(C_NAME_Authority)
	var r Authority
	err := c.Find(bson.M{"_id": id}).One(&r)
	if err != nil {
		if err != mgo.ErrNotFound {
			log.Printf("Get Authority [%s] Error: %s", id, err)
		}
		return &r, err
	}
	return &r, nil
}

// if not found, will return (nil, nil)
func GetAuthorityByEmail(email string) (*Authority, error) {
	email = strings.ToLower(email)
	c := db.MDB().C(C_NAME_Authority)
	var r *Authority
	err := c.Find(bson.M{"Email": email}).One(&r)
	if err != nil {
		if err != mgo.ErrNotFound {
			log.Printf("Get Authority [%s] Error: %s", email, err)
		} else {
			err = nil
		}
		return r, err
	}
	return r, nil
}

func GetAllAuthority() ([]Authority, error) {
	c := db.MDB().C(C_NAME_Authority)
	rs := []Authority{}
	err := c.Find(nil).All(&rs)
	if err != nil {
		log.Printf("GetAllAuthority Error: %s", err)
		return rs, err
	}
	return rs, nil
}

func InsertAuthority(a *Authority) error {
	a.Email = strings.ToLower(a.Email)
	ea, err := GetAuthorityByEmail(a.Email)
	if err != nil {
		return fmt.Errorf("check authority error: %s", err.Error())
	}
	if ea != nil && ea.Email == a.Email {
		return fmt.Errorf("authority for [%s] existed.", a.Email)
	}
	c := db.MDB().C(C_NAME_Authority)
	a.Id = bson.NewObjectId().Hex()
	err = c.Insert(a)
	return err
}

func UpdateAuthority(id string, a *Authority) error {
	c := db.MDB().C(C_NAME_Authority)
	change := bson.M{"$set": bson.M{
		"Desc":        a.Desc,
		"AdminLevel":  a.AdminLevel,
		"Roles":       a.Roles,
		"Allow":       a.Allow,
		"Deny":        a.Deny,
		"UpdatedTime": time.Now()}}
	err := c.UpdateId(id, change)
	return err
}

func DeleteAuthority(id string) error {
	c := db.MDB().C(C_NAME_Authority)
	err := c.RemoveId(id)
	return err
}
