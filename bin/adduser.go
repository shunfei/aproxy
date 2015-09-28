package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"aproxy/conf"
	"aproxy/module/auth"
	"aproxy/module/db"
)

var (
	inited = false

	confFile = flag.String("c", "aproxy.toml", "aproxy config file path")

	action     = flag.String("action", "", "action to do: adduser, setadmin")
	email      = flag.String("email", "", "email address")
	pwd        = flag.String("pwd", "", "password")
	adminLevel = flag.Int("adminlevel", 0, "50:System Administrator, 99:Super Administrator")
)

func main() {
	flag.Parse()

	semail := strings.TrimSpace(*email)
	semail = strings.ToLower(semail)
	spwd := strings.TrimSpace(*pwd)

	var err error
	saction := *action
	switch saction {
	case "adduser":
		initMongo()
		auth.SetUserStorageToMongo()
		err = addUser(semail, spwd)
		if err == nil {
			fmt.Println("success add user:", semail)
		}
	case "setadmin":
		initMongo()
		err = setAdmin(semail, *adminLevel)
		if err == nil {
			fmt.Printf("success set %s to admin\n", semail)
		}
	default:
		if saction == "" {
			err = fmt.Errorf("please enter action")
		} else {
			err = fmt.Errorf("wrong action [%s]", saction)
		}
	}
	if err != nil {
		log.Fatalln("add user error: ", err.Error())
	}
}

func initMongo() {
	if inited {
		return
	}
	conf.LoadAproxyConfig(*confFile)
	config := conf.Config()
	mgoConf := config.Db.Mongo
	err := db.InitMongoDB(mgoConf.Servers, mgoConf.Db)
	if err != nil {
		log.Fatalln("Can not set to MongoDB backend config storage.", mgoConf.Servers)
	}
	inited = true
}

func addUser(email, pwd string) error {
	if email == "" || pwd == "" {
		return fmt.Errorf("please enter email and pwd")
	}
	user := auth.User{}
	user.Email = email
	user.Name = email
	user.Pwd = pwd
	user.CreatedTime = time.Now()
	user.UpdatedTime = user.CreatedTime
	err := auth.InsertUser(user)
	return err
}

func setAdmin(email string, level int) error {
	if level != 50 && level != 99 {
		return fmt.Errorf("adminlevel must be 50 or 99")
	}
	authority, err := auth.GetAuthorityByEmail(email)
	if err != nil {
		return fmt.Errorf("query Authority for %s got error: %s",
			email, err.Error())
	}
	if authority != nil {
		authority.AdminLevel = level
		err = auth.UpdateAuthority(authority.Id, authority)
	} else {
		authority = &auth.Authority{}
		authority.Email = email
		authority.AdminLevel = level
		err = auth.InsertAuthority(authority)
	}
	return err
}
