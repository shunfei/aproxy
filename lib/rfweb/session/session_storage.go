package session

import (
	"errors"
	"fmt"
	"time"

	redis "gopkg.in/redis.v5"
)

var sessionStorage SessionStorager

type SessionStorager interface {
	Get(sid, key string) (string, error)
	// Set key to hold the string value and
	// set key to timeout after a given number of seconds.
	Set(sid, key, val string, exp int64) error

	Clear(sid string) error
}

type RedisSessionStorage struct {
	Addr     string
	Password string
	DB       int

	client *redis.Client
}

// addr: "127.0.0.1:6379"
func NewRedisSessionStorage(addr, pwd string, db int) (*RedisSessionStorage, error) {
	ss := &RedisSessionStorage{
		Addr:     addr,
		Password: pwd,
		DB:       db,
	}
	ss.client = redis.NewClient(&redis.Options{
		Addr:        addr,
		Password:    pwd,
		DB:          db,
		IdleTimeout: 3 * time.Minute,
	})
	_, err := ss.client.Ping().Result()
	if err != nil {
		return nil, err
	}
	// // keep-alive
	// go func() {
	// 	for {
	// 		ss.client.Ping()
	// 		time.Sleep(10 * time.Second)
	// 	}
	// }()
	return ss, nil
}

func (self *RedisSessionStorage) Get(sid, key string) (string, error) {
	val, err := self.client.HGet(sid, key).Result()
	if err != nil {
		if err == redis.Nil {
			err = nil
		} else {
			err = errors.New("redis hget error: " + err.Error())
		}
	}
	return val, err
}

func (self *RedisSessionStorage) Set(sid, key, val string, expiration int64) error {
	_, err := self.client.HSet(sid, key, val).Result()
	if err == nil && expiration > 0 {
		self.client.Expire(sid, time.Duration(expiration)*time.Second)
	}
	return err
}

func (self *RedisSessionStorage) Clear(sid string) error {
	r := self.client.Del(sid)
	if r.Err() != nil {
		return r.Err()
	}
	return nil
}

//
//

func SetSessionStorager(ss SessionStorager) {
	if ss == nil {
		panic("[SetSessionStorager] SessionStorager must not nil.")
	}
	sessionStorage = ss
}

func SetSessionStoragerToRedis(addr, pwd string, db int) error {
	ss, err := NewRedisSessionStorage(addr, pwd, db)
	if err != nil {
		return fmt.Errorf("[SetSessionStoragerToRedis] faild: %s", err.Error())
	}
	sessionStorage = ss
	return nil
}
