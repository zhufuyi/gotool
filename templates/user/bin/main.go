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

	"github.com/zhufuyi/goctl/templates/user/configs"
	"github.com/zhufuyi/goctl/templates/user/global"
	"github.com/zhufuyi/goctl/templates/user/routers"

	"github.com/gin-gonic/gin"
	"github.com/zhufuyi/pkg/email"
	"github.com/zhufuyi/pkg/jwt"
	"github.com/zhufuyi/pkg/logger"
	"github.com/zhufuyi/pkg/mysql"
	"github.com/zhufuyi/pkg/snowflake"
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
	initConfigs()   // 初始化配置
	initLog()       // 初始化日志
	initSnowflake() // 初始化id
	initJWT()       // 初始化jwt
	initEmail()     // 初始化邮件
	initMysql()     // 初始化mysql
	initWeb()       // 初始化web服务

	logger.Infof("%v initialization successfully", []string{
		"log",
		"snowflake",
		"jwt",
		"email",
		"mysql",
		"web",
	})
}

func initConfigs() {
	flag.StringVar(&configFile, "c", "../configs/conf.yml", "配置文件")
	//flag.StringVar(&configFile, "c", "templates/demo/configs/conf.yml", "配置文件")
	flag.Parse()
	err := configs.Init(configFile)
	if err != nil {
		panic(err)
	}
	global.Conf = configs.Get()
}

func initLog() {
	l, err := logger.Init(
		logger.WithLevel(global.Conf.Log.Level),
		logger.WithFormat(global.Conf.Log.Format),
	)
	if err != nil {
		panic(err)
	}
	_ = l
	//global.Logger = l

	logger.Info("config info", logger.Any("config", configs.ShowConfig()))
}

func initSnowflake() {
	err := snowflake.Init(1)
	if err != nil {
		panic(err)
	}
}

func initJWT() {
	jwt.Init(
		jwt.WithSigningKey(global.Conf.Jwt.SigningKey),
		jwt.WithExpire(time.Duration(global.Conf.Jwt.Expire)*time.Second),
	)
}

func initEmail() {
	client, err := email.Init(global.Conf.Email.Sender, global.Conf.Email.Password)
	if err != nil {
		panic(err)
	}

	global.EmailCli = client
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
		logger.Infof("start up '%s' service, port:%d", global.Conf.ServerName, global.Conf.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
}

func wait() {
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
