package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"aproxy/conf"
	"aproxy/lib/rfweb/session"
	"aproxy/module/auth"
	"aproxy/module/auth/login"
	bkconf "aproxy/module/backend_conf"
	"aproxy/module/db"
	"aproxy/module/proxy"
	"aproxy/module/setting"
)

var (
	confFile = flag.String("c", "aproxy.toml", "aproxy config file path")
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.Parse()
	conf.LoadAproxyConfig(*confFile)

	config := conf.Config()

	if !checkWebDir(config.WebDir) {
		os.Exit(1)
		return
	}

	mgoConf := config.Db.Mongo
	err := db.InitMongoDB(mgoConf.Servers, mgoConf.Db)
	if err != nil {
		log.Fatalln("Can not set to MongoDB backend config storage.", mgoConf.Servers)
	}
	// Set backend-config storage to MongoDB
	bkconf.SetBackendConfStorageToMongo()
	// Set user storage to MongoDB
	auth.SetUserStorageToMongo()

	// session
	ssConf := config.Session
	session.InitSessionServer(ssConf.Domain, ssConf.Cookie, ssConf.Expiration)
	session.SetSessionStoragerToRedis(ssConf.Redis.Addr,
		ssConf.Redis.Password, ssConf.Redis.Db)

	// login
	login.InitLoginServer(config.LoginHost, config.AproxyUrlPrefix)

	// setting manager
	setting.InitSettingServer(config.WebDir, config.AproxyUrlPrefix)

	lhost := config.Listen
	mux := http.NewServeMux()
	// setting
	setPre := setting.AproxyUrlPrefix
	apiApp := setting.NewApiApp()
	mux.HandleFunc(apiApp.UrlPrefix, apiApp.ServeHTTP)
	mux.HandleFunc(setPre, setting.StaticServer)
	// proxy
	mux.HandleFunc("/", proxy.Proxy)
	s := &http.Server{
		Addr:    lhost,
		Handler: mux,
	}
	log.Println("Starting aproxy on " + lhost)
	err = s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func checkWebDir(webDir string) bool {
	absPath, _ := filepath.Abs(webDir)
	_, err := os.Stat(absPath)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		log.Println("webdir is not exist:", absPath)
		log.Println("please change the webdir in your aproxy config file.")
		return false
	}
	return true
}
