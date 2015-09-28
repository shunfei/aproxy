package db

// import (
// 	"errors"
// 	"fmt"
// 	"log"
// 	"net/url"
// 	"strings"
// 	"time"

// 	"gopkg.in/mgo.v2"
// 	"gopkg.in/mgo.v2/bson"
// )

// type DBer interface {
// 	GetById(table, id string, val interface{}) error
// 	GetOne(table string, query map[string]interface{}, val interface{}) error
// 	GetAll(table string, vals interface{}) error
// 	Update(table, id string, val interface{}) error
// 	Insert(table string, val interface{}) error
// }

// type MongoDB struct {
// 	Servers   []string
// 	Db        string
// 	dbSession *mgo.Session
// }

// func (self *MongoDB) GetById(table, id string, val interface{}) error {

// }
