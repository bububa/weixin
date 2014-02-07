package main

import (
	"encoding/json"
	"flag"
	"fmt"
	log "github.com/bububa/factorlog"
	"github.com/bububa/goconfig/config"
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

	apps := make(map[string]*App)
	json.Unmarshal([]byte(appCfg), &apps)
	for id, app := range apps {
		mux := weixin.New(app.Token, app.Id, app.Secret)
		mux.HandleFunc(weixin.MsgTypeText, Echo)
		mux.HandleFunc(weixin.MsgTypeEventSubscribe, Subscribe)
		http.Handle("/"+id, mux) // 注册接收微信服务器数据的接口URI
	}
	err = http.ListenAndServe(fmt.Sprintf(":%d", appPort), nil) // 启动接收微信数据服务器
	if err != nil {
		logger.Fatal(err)
	}
}
