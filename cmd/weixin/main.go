package main

import (
	"encoding/json"
	"flag"
	"fmt"
	log "github.com/bububa/factorlog"
	"github.com/bububa/goconfig/config"
	"github.com/bububa/pg"
	"github.com/bububa/weixin"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"time"
)

type App struct {
	Token  string
	Id     string
	Secret string
}

var (
	logFlag    = flag.String("log", "", "set log path")
	configFile = flag.String("config", "", "set config file")
	logger     *log.FactorLog
)

func init() {
	rand.Seed(int64(time.Second))
}

func SetGlobalLogger(logPath string) *log.FactorLog {
	sfmt := `%{Color "red:white" "CRITICAL"}%{Color "red" "ERROR"}%{Color "yellow" "WARN"}%{Color "green" "INFO"}%{Color "cyan" "DEBUG"}%{Color "blue" "TRACE"}[%{Date} %{Time}] [%{SEVERITY}:%{ShortFile}:%{Line}] %{Message}%{Color "reset"}`
	logger := log.New(os.Stdout, log.NewStdFormatter(sfmt))
	if len(logPath) > 0 {
		logf, err := os.OpenFile(logPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0640)
		if err != nil {
			return logger
		}
		logger = log.New(logf, log.NewStdFormatter(sfmt))
	}
	logger.SetSeverities(log.INFO | log.WARN | log.ERROR | log.FATAL | log.CRITICAL)
	return logger
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()

	logger = SetGlobalLogger(*logFlag)

	weixin.SetLogger(logger)

	if *configFile == "" {
		logger.Fatal("need a config file")
	}
	cfg, err := config.ReadDefault(*configFile)
	if err != nil {
		logger.Fatal(err)
	}
	appCfg, err := cfg.String("weixin", "apps")
	if err != nil {
		logger.Fatal("need set apps in config file in json format")
	}

	appPort, err := cfg.Int("weixin", "port")
	if err != nil {
		logger.Fatal("need set port in config file")
	}

	pgHost, _ := cfg.String("pg", "host")
	pgPort, _ := cfg.String("pg", "port")
	pgUser, _ := cfg.String("pg", "user")
	pgPassword, _ := cfg.String("pg", "passwd")
	pgDbname, _ := cfg.String("pg", "dbname")

	pgOptions := &pg.Options{
		Host:     pgHost,
		Port:     pgPort,
		User:     pgUser,
		Password: pgPassword,
		Database: pgDbname,
		PoolSize: 10,
	}

	pgDb := pg.Connect(pgOptions)
	defer pgDb.Close()

	apps := make(map[string]*App)
	json.Unmarshal([]byte(appCfg), &apps)
	for id, app := range apps {
		mux := weixin.New(id, app.Token, app.Id, app.Secret)
		mux.SetDB(pgDb)
		mux.HandleFunc(weixin.MsgTypeText, MsgTxt)
		mux.HandleFunc(weixin.MsgTypeImage, MsgImage)
		mux.HandleFunc(weixin.MsgTypeVoice, MsgVoice)
		mux.HandleFunc(weixin.MsgTypeVideo, MsgVideo)
		mux.HandleFunc(weixin.MsgTypeLocation, MsgLocation)
		mux.HandleFunc(weixin.MsgTypeLink, MsgLink)
		mux.HandleFunc(weixin.MsgTypeEventSubscribe, Subscribe)
		mux.HandleFunc(weixin.MsgTypeEventUnsubscribe, Unsubscribe)
		mux.HandleFunc("/subscribers", GetSubscribers)
		mux.HandleFunc("/groups", GetGroups)
		mux.HandleFunc("/group/create", CreateGroup)
		mux.HandleFunc("/group/changename", ChangeGroupName)
		mux.HandleFunc("/user", GetUser)
		mux.HandleFunc("/user/group", GetUserGroup)
		mux.HandleFunc("/user/changegroup", ChangeUserGroup)
		mux.HandleFunc("/menu/create", CreateMenu)
		http.Handle("/"+id+"/", mux) // 注册接收微信服务器数据的接口URI
	}
	err = http.ListenAndServe(fmt.Sprintf(":%d", appPort), nil) // 启动接收微信数据服务器
	if err != nil {
		logger.Fatal(err)
	}
}
