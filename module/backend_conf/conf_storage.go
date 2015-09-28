package backend_conf

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"aproxy/module/db"
)

const (
	C_NAME_BackendConf = "BackendConf"
)

var backendConfStorag BackendConfStorager

type BackendConf struct {
	Id   string `bson:"_id,omitempty"`
	Desc string `bson:"Desc"`
	// request hostname,
	// e.g: www.abc.com, 192.168.10.33
	HostName string `bson:"HostName"`
	// proxy to which host,
	// e.g: http://localhost:8081/
	UpStreams []string `bson:"UpStreams"`
	// authentication type:
	//     0. public: every one
	//     1. login: login user
	//     2. auth: login & has permitted
	AuthType int `bson:"AuthType"`

	CreatedTime time.Time `bson:"CreatedTime"`
	UpdatedTime time.Time `bson:"UpdatedTime"`
}

type BackendConfStorager interface {
	Get(hostname string) (BackendConf, error)
	GetAll() ([]BackendConf, error)
	Insert(BackendConf) error
	Update(id string, bc BackendConf) error
	Delete(id string) error
}

// using mongodb for BackendConf Storage
type MongoBackendConfStorage struct {
}

func (self *MongoBackendConfStorage) Get(hostname string) (BackendConf, error) {
	c := db.MDB().C(C_NAME_BackendConf)
	var bc BackendConf
	err := c.Find(bson.M{"HostName": hostname}).One(&bc)
	if err != nil {
		if err != mgo.ErrNotFound {
			log.Printf("MongoBackendConfStorage Get [%s] Error: %s", hostname, err)
		}
		return bc, err
	}
	return bc, nil
}

func (self *MongoBackendConfStorage) GetAll() ([]BackendConf, error) {
	c := db.MDB().C(C_NAME_BackendConf)
	bcs := []BackendConf{}
	err := c.Find(nil).All(&bcs)
	if err != nil {
		log.Printf("MongoBackendConfStorage GetAll Error: %s", err)
		return bcs, err
	}
	return bcs, nil
}

func (self *MongoBackendConfStorage) Insert(bc BackendConf) error {
	bc.Id = bson.NewObjectId().Hex()
	c := db.MDB().C(C_NAME_BackendConf)
	err := c.Insert(bc)
	return err
}

func (self *MongoBackendConfStorage) Update(id string, bc BackendConf) error {
	c := db.MDB().C(C_NAME_BackendConf)
	change := bson.M{"$set": bson.M{
		"AuthType":    bc.AuthType,
		"Desc":        bc.Desc,
		"HostName":    bc.HostName,
		"UpStreams":   bc.UpStreams,
		"UpdatedTime": time.Now()}}
	err := c.UpdateId(id, change)
	return err
}

func (self *MongoBackendConfStorage) Delete(id string) error {
	c := db.MDB().C(C_NAME_BackendConf)
	err := c.RemoveId(id)
	return err
}

//
//

func SetBackendConfStorage(bcs BackendConfStorager) {
	if bcs == nil {
		panic("SetBackendConfStorage: BackendConfStorager MUST NOT (nil)!")
	}
	backendConfStorag = bcs
}

func SetBackendConfStorageToMongo() error {
	SetBackendConfStorage(&MongoBackendConfStorage{})
	return nil
}

//
//

func Get(hostname string) (BackendConf, error) {
	hostname = strings.ToLower(hostname)
	return backendConfStorag.Get(hostname)
}

func GetAll() ([]BackendConf, error) {
	return backendConfStorag.GetAll()
}

func Insert(bc BackendConf) error {
	err := validBackendConf(bc)
	if err != nil {
		return err
	}
	ebc, err2 := Get(strings.ToLower(bc.HostName))
	if err2 == nil && ebc.HostName == bc.HostName {
		return errors.New(fmt.Sprintf("host [%s] has exist.", bc.HostName))
	}
	return backendConfStorag.Insert(bc)
}

func Update(id string, bc BackendConf) error {
	err := validBackendConf(bc)
	if err != nil {
		return err
	}
	return backendConfStorag.Update(id, bc)
}

func Delete(id string) error {
	return backendConfStorag.Delete(id)
}

//
//

func validBackendConf(bc BackendConf) error {
	err := validHostName(bc.HostName)
	if err != nil {
		return err
	}
	for _, upstream := range bc.UpStreams {
		err = validUpstream(upstream)
		if err != nil {
			return err
		}
	}
	return nil
}

func validHostName(hostname string) error {
	if strings.Index(hostname, "http://") == 0 ||
		strings.Index(hostname, "https://") == 0 {
		return errors.New(fmt.Sprintf("hostname [%s] don't need http:// .", hostname))
	}
	url_, err := url.Parse("http://" + hostname)
	if err != nil {
		return err
	} else if url_.Host == "" {
		return errors.New(fmt.Sprintf("hostname [%s] is wrong.", hostname))
	} else if url_.Host != hostname {
		return errors.New(fmt.Sprintf("hostname [%s] must not contains url-path.", hostname))
	}
	return nil
}

func validUpstream(upstream string) error {
	url_, err := url.Parse(upstream)
	if err != nil {
		return err
	}
	if strings.Index(url_.Scheme, "http") != 0 {
		return errors.New(fmt.Sprintf("upstream [%s] must contains http:// or https:// .", upstream))
	}
	if url_.Host == "" {
		return errors.New(fmt.Sprintf("upstream [%s] has the wrong host .", upstream))
	}
	return nil
}
