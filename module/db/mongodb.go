package db

import (
	"strings"
	"time"

	"gopkg.in/mgo.v2"
)

var (
	_mongoStorage *mongoStorage
)

type mongoStorage struct {
	Servers   []string
	Db        string
	dbSession *mgo.Session
}

func newMongoStorage(servers []string, db string) (*mongoStorage, error) {
	storage := &mongoStorage{}
	storage.Servers = servers
	storage.Db = db
	session, err := mgo.Dial(strings.Join(servers, ","))
	if err == nil {
		storage.dbSession = session
		go autoReconnect(session)
	}
	return storage, err
}

func InitMongoDB(servers []string, db string) error {
	var err error
	_mongoStorage, err = newMongoStorage(servers, db)
	return err
}

func MDB() *mgo.Database {
	return _mongoStorage.dbSession.DB(_mongoStorage.Db)
}

func autoReconnect(session *mgo.Session) {
	var err error
	for {
		err = session.Ping()
		if err != nil {
			// fmt.Println("Loss connection to MongoDB !!")
			session.Refresh()
			// err = session.Ping()
			// if err == nil {
			// 	fmt.Println("Reconnect to MongoDB successful.")
			// } else {
			// 	fmt.Println("Reconnect to MongoDB faild !!")
			// }
		}
		time.Sleep(time.Second * 10)
	}
}
