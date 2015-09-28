package auth

import (
	// "fmt"
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
	C_NAME_Role = "Role"
)

type Role struct {
	sync.RWMutex `bson:"-"`

	Id    string   `bson:"_id,omitempty"`
	Name  string   `bson:"Name"`
	Desc  string   `bson:"Desc"`
	Allow []string `bson:"Allow"`
	Deny  []string `bson:"Deny"`

	CreatedTime time.Time `bson:"CreatedTime"`
	UpdatedTime time.Time `bson:"UpdatedTime"`

	inited   bool
	al       int
	dl       int
	allowReg []*regexp.Regexp
	denyReg  []*regexp.Regexp
}

func (self *Role) Init() {
	self.Lock()
	defer self.Unlock()
	if self.inited {
		return
	}
	self.init()
}

func (self *Role) init() {
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
	self.al = len(self.allowReg)
	self.dl = len(self.denyReg)
	return
}

// must Init() before use this.
func (self *Role) HasPermission(url string) bool {
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
	return false
}

func GetRoleByID(id string) (*Role, error) {
	c := db.MDB().C(C_NAME_Role)
	var r Role
	err := c.Find(bson.M{"_id": id}).One(&r)
	if err != nil {
		if err != mgo.ErrNotFound {
			log.Printf("Get Role [%s] Error: %s", id, err)
		}
		return &r, err
	}
	return &r, nil
}

func GetAllRole() ([]Role, error) {
	c := db.MDB().C(C_NAME_Role)
	rs := []Role{}
	err := c.Find(nil).All(&rs)
	if err != nil {
		log.Printf("GetAllRole Error: %s", err)
		return rs, err
	}
	return rs, nil
}

func InsertRole(role *Role) error {
	c := db.MDB().C(C_NAME_Role)
	role.Id = bson.NewObjectId().Hex()
	err := c.Insert(role)
	return err
}

func UpdateRole(id string, role *Role) error {
	c := db.MDB().C(C_NAME_Role)
	change := bson.M{"$set": bson.M{
		"Name":        role.Name,
		"Desc":        role.Desc,
		"Allow":       role.Allow,
		"Deny":        role.Deny,
		"UpdatedTime": time.Now()}}
	err := c.UpdateId(id, change)
	return err
}

func DeleteRole(id string) error {
	c := db.MDB().C(C_NAME_Role)
	err := c.RemoveId(id)
	return err
}
