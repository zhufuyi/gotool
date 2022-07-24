package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/zhufuyi/goctl/templates/web/configs"
	"github.com/zhufuyi/goctl/templates/web/global"
	"github.com/zhufuyi/goctl/templates/web/routers"

	"github.com/gin-gonic/gin"
	"github.com/zhufuyi/pkg/logger"
	"github.com/zhufuyi/pkg/mysql"
)

var (
	configFile string
	server     *http.Server
)

func main() {
	initServer()

	wait()
}

func initServer() {
	// 初始化配置
	initConfigs()

	// 初始化日志，输出到文件
	initLog()

	// 初始化mysql
	initMysql()

	// 初始化web服务
	initWeb()
}

func initConfigs() {
	flag.StringVar(&configFile, "c", "../configs/conf.yml", "配置文件")
	flag.Parse()
	err := configs.Init(configFile)
	if err != nil {
		panic(err.Error())
	}
	global.Conf = configs.Get()
}

func initLog() {
	l, err := logger.Init(
		logger.WithLevel(global.Conf.Log.Level),
		logger.WithFormat(global.Conf.Log.Format),
		//logger.WithSave(
		//	global.Conf.Log.IsSave,
		//	logger.WithFileName(global.Conf.Log.LogFileConfig.Filename),
		//	logger.WithFileMaxAge(global.Conf.Log.LogFileConfig.MaxAge),
		//	logger.WithFileMaxBackups(global.Conf.Log.LogFileConfig.MaxBackups),
		//	logger.WithFileMaxAge(global.Conf.Log.LogFileConfig.MaxAge),
		//	logger.WithFileIsCompression(global.Conf.Log.LogFileConfig.IsCompression),
		//),
	)
	if err != nil {
		panic(err)
	}
	_ = l
	//global.Logger = l
}

func initMysql() {
	db, err := mysql.Init(global.Conf.MysqlURL, mysql.WithLog(logger.Get()))
	if err != nil {
		panic(err)
	}
	global.MysqlDB = db
}

func initWeb() {
	if global.Conf.RunMode == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	router := routers.NewRouter()
	server = &http.Server{
		Addr:           ":" + strconv.Itoa(global.Conf.ServerPort),
		Handler:        router,
		ReadTimeout:    time.Second * time.Duration(global.Conf.ReadTimeout),
		WriteTimeout:   time.Second * time.Duration(global.Conf.WriteTimeout),
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("s.ListenAndServe err", logger.Err(err))
		}
	}()
}

func wait() {
	logger.Infof("'%s' server start successful, port:%d", global.Conf.ServerName, global.Conf.ServerPort)
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Infof("'%s' server shutting down ......", global.Conf.ServerName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("'%s' server forced to shutdown, err: %v", global.Conf.ServerName, err)
	}

	logger.Infof("'%s' server exited", global.Conf.ServerName)
}
