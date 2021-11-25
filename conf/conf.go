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
	AuditLogPath    string
	Session         struct {
		Cookie     string
		Domain     string
		Expiration int64
		Redis      struct {
			Addr     string
			Password string
			Db       int
		}
	}
	Db struct {
		Mongo struct {
			Servers []string
			Db      string
		}
	}
	Oauth struct {
		Open   bool
		Github struct {
			Open         bool
			ClientID     string
			ClientSecret string
		}
	}
}

var aproxyConfig AproxyConfig

func LoadAproxyConfig(tomlFile string) error {
	if _, err := toml.DecodeFile(tomlFile, &aproxyConfig); err != nil {
		return fmt.Errorf("Load config file [%s] faild: %s",
			tomlFile, err)
	}
	return nil
}

func Config() *AproxyConfig {
	return &aproxyConfig
}
