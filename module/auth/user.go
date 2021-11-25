package auth

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"aproxy/lib/util"
	"aproxy/module/db"
)

const (
	C_NAME_User = "User"
)

var userStorager UserStorager

type User struct {
	Id    string `bson:"_id,omitempty"`
	Name  string `bson:"Name"`
	Email string `bson:"Email"`
	Desc  string `bson:"Desc"` // description
	Pwd   string `bson:"Pwd" json:"-"`

	CreatedTime time.Time `bson:"CreatedTime"`
	UpdatedTime time.Time `bson:"UpdatedTime"`
}

type UserStorager interface {
	Login(email, pwd string) (*User, error)
	GetByEmail(email string) (*User, error)
	GetAll() ([]User, error)
	// add new user.
	// user.Pwd field has encrypted.
	Insert(user User) error
	Update(id string, user User) error
	Delete(id string) error
}

type MongoUserStorage struct {
}

func (self *MongoUserStorage) Login(email, pwd string) (*User, error) {
	if email == "" || pwd == "" {
		return nil, errors.New("please enter email and password.")
	}
	user, err := self.GetByEmail(email)
	if err != nil {
		return nil, errors.New("query user by email got error: " + err.Error())
	}
	if user == nil || user.Email == "" {
		return nil, errors.New("email or password wrong.")
	}
	err = util.CompareHashAndPassword([]byte(user.Pwd), []byte(pwd))
	if err != nil {
		return nil, errors.New("email or password wrong.")
	}
	return user, nil
}

func (self *MongoUserStorage) Insert(user User) error {
	user.Id = bson.NewObjectId().Hex()
	c := db.MDB().C(C_NAME_User)
	err := c.Insert(user)
	return err
}

func (self *MongoUserStorage) Update(id string, user User) error {
	if len(user.Pwd) < 10 {
		return nil
	}
	c := db.MDB().C(C_NAME_User)
	change := bson.M{
		"Pwd":         user.Pwd,
		"UpdatedTime": time.Now()}
	err := c.UpdateId(id, bson.M{"$set": change})
	return err
}

func (self *MongoUserStorage) Delete(id string) error {
	c := db.MDB().C(C_NAME_User)
	err := c.RemoveId(id)
	return err
}

func (self *MongoUserStorage) GetByEmail(email string) (*User, error) {
	c := db.MDB().C(C_NAME_User)
	var user User
	err := c.Find(bson.M{"Email": email}).One(&user)
	if err != nil {
		if err != mgo.ErrNotFound {
			log.Printf("MongoUserStorage GetByEmail [%s] Error: %s", email, err)
			return nil, err
		} else {
			return nil, nil
		}
	}
	return &user, nil
}

func (self *MongoUserStorage) GetAll() ([]User, error) {
	c := db.MDB().C(C_NAME_User)
	users := []User{}
	err := c.Find(nil).All(&users)
	if err != nil {
		log.Printf("User GetAll Error: %s", err)
		return users, err
	}
	return users, nil
}

func validUser(user User) error {
	return nil
}

//
//

func SetUserStorage(us UserStorager) {
	if us == nil {
		panic("SetUserStorage: UserStorager MUST NOT (nil)!")
	}
	userStorager = us
}

func SetUserStorageToMongo() error {
	SetUserStorage(&MongoUserStorage{})
	return nil
}

//
//

func LoginUser(email, pwd string) (*User, error) {
	return userStorager.Login(email, pwd)
}

func GetUserByEmail(email string) (*User, error) {
	email = strings.ToLower(email)
	return userStorager.GetByEmail(email)
}

func GetAllUsers() ([]User, error) {
	return userStorager.GetAll()
}

func InsertUser(user User) error {
	err := validUser(user)
	if err != nil {
		return err
	}
	euser, err2 := userStorager.GetByEmail(user.Email)
	if err2 == nil && euser != nil && euser.Email == user.Email {
		return errors.New(fmt.Sprintf("email [%s] has exist.", user.Email))
	}
	pwd, err := util.CryptPassword([]byte(user.Pwd))
	if err != nil {
		return err
	}
	user.Pwd = string(pwd)
	return userStorager.Insert(user)
}

func UpdateUser(id string, user User) error {
	if len(user.Pwd) > 0 {
		pwd, err := util.CryptPassword([]byte(user.Pwd))
		if err != nil {
			return err
		}
		user.Pwd = string(pwd)
	}
	return userStorager.Update(id, user)
}

func DeleteUser(id string) error {
	if len(id) < 1 {
		return errors.New("wrong id")
	}
	return userStorager.Delete(id)
}
