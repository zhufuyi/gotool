package main

import (
	"context"
	"flag"
	"strconv"
	"time"

	"github.com/zhufuyi/goctl/templates/http_server/config"
	"github.com/zhufuyi/goctl/templates/http_server/internal/model"
	"github.com/zhufuyi/goctl/templates/http_server/internal/server"

	"github.com/zhufuyi/pkg/logger"
	"github.com/zhufuyi/pkg/tracer"
)

var (
	version    string
	configFile string
)

// @title projectExample api docs
// @description http server api docs
// @schemes http https
// @version v0.0.0
// @host localhost:8080
func main() {
	inits := registerInits()
	servers := registerServers()
	closes := registerCloses(servers)

	s := server.NewApp(inits, servers, closes)
	s.Run()
}

// -------------------------------- 注册app初始化 ---------------------------------
func registerInits() []server.Init {
	// 初始化配置文件，必须优先执行，后面的初始化需要依赖配置
	flag.StringVar(&configFile, "c", "config/conf.yml", "配置文件")
	flag.StringVar(&version, "version", "", "服务版本号")
	flag.Parse()
	err := config.Init(configFile)
	if err != nil {
		panic("init config error: " + err.Error())
	}
	if version != "" {
		config.Get().Version = version
	}

	var inits []server.Init

	// 初始化日志
	inits = append(inits, func() {
		_, err := logger.Init(
			logger.WithLevel(config.Get().Log.Level),
			logger.WithFormat(config.Get().Log.Format),
		)
		if err != nil {
			panic("init logger error: " + err.Error())
		}
		logger.Info("config info", logger.Any("config", config.Show()))
	})

	// 初始化数据库
	inits = append(inits, func() {
		model.InitMysql()
		model.InitRedis()
	})

	// 初始化链路跟踪
	inits = append(inits, func() {
		if config.Get().EnableTracing { // 根据配置是否开启链路跟踪
			exporter, err := tracer.NewJaegerAgentExporter(config.Get().Jaeger.AgentHost, config.Get().Jaeger.AgentPort)
			if err != nil {
				panic("init trace error:" + err.Error())
			}
			resource := tracer.NewResource(
				tracer.WithServiceName(config.Get().ServiceName),
				tracer.WithEnvironment(config.Get().ServiceEnv),
				tracer.WithServiceVersion(version),
			)

			tracer.Init(exporter, resource, config.Get().Jaeger.SamplingRate) // 如果SamplingRate=0.5表示只采样50%
		}
	})

	return inits
}

// -------------------------------- 注册app服务 ---------------------------------
func registerServers() []server.IServer {
	var servers []server.IServer

	// 创建http服务
	servers = append(servers, server.NewHTTPServer(
		":"+strconv.Itoa(config.Get().ServicePort),
		time.Second*time.Duration(config.Get().ReadTimeout),
		time.Second*time.Duration(config.Get().WriteTimeout),
	))

	return servers
}

// -------------------------- 注册app需要释放的资源  -------------------------------------------
func registerCloses(servers []server.IServer) []server.Close {
	var closes []server.Close

	// 关闭服务
	for _, server := range servers {
		closes = append(closes, server.Stop)
	}

	// 关闭mysql
	closes = append(closes, func() error {
		return model.CloseMysql()
	})

	// 关闭redis
	closes = append(closes, func() error {
		return model.CloseRedis()
	})

	// 关闭trace
	closes = append(closes, func() error {
		if config.Get().EnableTracing {
			ctx, _ := context.WithTimeout(context.Background(), 2*time.Second) //nolint
			return tracer.Close(ctx)
		}
		return nil
	})

	return closes
}
