package conf

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type AproxyConfig struct {
	Listen          string
	WebDir          string
	LoginHost       string
	AproxyUrlPrefix string
	Session         struct {
		Cookie     string
		Domain     string
		Expiration int64
		Redis      struct {
			Addr     string
			Password string
			Db       int64
		}
	}
	Db struct {
		Mongo struct {
			Servers []string
			Db      string
		}
	}
}

var aproxyConfig AproxyConfig

func LoadAproxyConfig(tomlFile string) {
	if _, err := toml.DecodeFile(tomlFile, &aproxyConfig); err != nil {
		panic(fmt.Sprintf("Load config file [%s] faild: %s",
			tomlFile, err))
	}
}

func Config() *AproxyConfig {
	return &aproxyConfig
}
