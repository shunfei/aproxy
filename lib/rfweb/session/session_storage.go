package session

import (
	"fmt"
	"time"

	"gopkg.in/redis.v3"
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
	DB       int64

	client *redis.Client
}

// addr: "127.0.0.1:6379"
func NewRedisSessionStorage(addr, pwd string, db int64) (*RedisSessionStorage, error) {
	ss := &RedisSessionStorage{
		Addr:     addr,
		Password: pwd,
		DB:       db,
	}
	ss.client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       db,
	})
	_, err := ss.client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return ss, nil
}

func (self *RedisSessionStorage) Get(sid, key string) (string, error) {
	val, err := self.client.HGet(sid, key).Result()
	if err != nil && err == redis.Nil {
		err = nil
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

func SetSessionStoragerToRedis(addr, pwd string, db int64) error {
	ss, err := NewRedisSessionStorage(addr, pwd, db)
	if err != nil {
		return fmt.Errorf("[SetSessionStoragerToRedis] faild: %s", err.Error())
	}
	sessionStorage = ss
	return nil
}
