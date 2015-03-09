package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bububa/goconfig/config"
	"github.com/bububa/pg"
	"github.com/bububa/weixin"
	log "github.com/kdar/factorlog"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	_CONFIG_FILE = "/var/code/go/weixin.cfg"
)

type App struct {
	Token  string
	Id     string
	Secret string
}

var (
	logFlag     = flag.String("log", "", "set log path")
	configFile  = flag.String("config", "", "set config file")
	logger      *log.FactorLog
	scenesMap   *weixin.ScenesMap
	apiHost     string
	callbackUrl map[string]string
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	logger = SetGlobalLogger(*logFlag)
	weixin.SetLogger(logger)
	flag.Parse()
	rand.Seed(int64(time.Second))

	scenesMap = &weixin.ScenesMap{
		Mutex:  new(sync.Mutex),
		Scenes: make(map[uint64]weixin.SceneParams),
	}
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
	logger.SetSeverities(log.INFO | log.WARN | log.ERROR | log.DEBUG | log.FATAL | log.CRITICAL | log.TRACE)
	return logger
}

func main() {
	conf := _CONFIG_FILE
	if *configFile != "" {
		conf = *configFile
	}
	cfg, _ := config.ReadDefault(conf)

	appCfg, err := cfg.String("weixin", "apps")
	if err != nil {
		logger.Fatal("need set apps in config file in json format")
	}
	appPort, err := cfg.Int("weixin", "port")
	if err != nil {
		logger.Fatal("need set port in config file")
	}

	apiHost, _ = cfg.String("weixin", "apihost")
	callback, _ := cfg.String("weixin", "callbackurl")
	callbackUrl = make(map[string]string)
	json.Unmarshal([]byte(callback), &callbackUrl)

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
		mux.SetDb(pgDb)
		mux.HandleFunc(weixin.MsgTypeText, MsgTxt)                  // 文本事件
		mux.HandleFunc(weixin.MsgTypeImage, MsgImage)               // 图片事件
		mux.HandleFunc(weixin.MsgTypeVoice, MsgVoice)               // 语音事件
		mux.HandleFunc(weixin.MsgTypeVideo, MsgVideo)               // 视频事件
		mux.HandleFunc(weixin.MsgTypeLocation, MsgLocation)         // 地理位置事件
		mux.HandleFunc(weixin.MsgTypeLink, MsgLink)                 // 链接事件
		mux.HandleFunc(weixin.MsgTypeEventScan, MsgScan)            // 扫描qrcode事件
		mux.HandleFunc(weixin.MsgTypeEventSubscribe, Subscribe)     // 订阅事件
		mux.HandleFunc(weixin.MsgTypeEventUnsubscribe, Unsubscribe) // 取消订阅事件
		mux.HandleFunc("/subscribers", GetSubscribers)              // 更新订阅信息
		mux.HandleFunc("/groups", GetGroups)                        // 获取组
		mux.HandleFunc("/group/create", CreateGroup)                // 创建组
		mux.HandleFunc("/group/changename", ChangeGroupName)        // 组改名
		mux.HandleFunc("/user", GetUser)                            // 获取用户信息
		mux.HandleFunc("/user/group", GetUserGroup)                 // 获取用户组
		mux.HandleFunc("/user/changegroup", ChangeUserGroup)        // 更改用户组
		mux.HandleFunc("/menu/create", CreateMenu)                  // 创建自定义菜单
		mux.HandleFunc("/menu/delete", DeleteMenu)                  // 删除自定义菜单
		mux.HandleFunc("/qrcode/create", CreateQrcode)              // 生成qrcode
		mux.HandleFunc("/qrcode/show", ShowQrcode)                  // 显示qrcode
		mux.HandleFunc("/qrcode/temp/create", CreateTempQrcode)     // 生成临时qrcode
		mux.HandleFunc("/qrcode/temp/show", ShowTempQrcode)         // 显示临时qrcode
		mux.HandleFunc("/message/custom/send", SendCustomMessage)   // 推送消息
		http.Handle("/"+id+"/", mux)                                // 注册接收微信服务器数据的接口URI
	}
	logger.Tracef("Server running at %d", appPort)
	err = http.ListenAndServe(fmt.Sprintf(":%d", appPort), nil) // 启动接收微信数据服务器
	if err != nil {
		logger.Fatal(err)
	}
}
