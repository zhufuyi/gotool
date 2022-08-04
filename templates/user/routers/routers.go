package routers

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/zhufuyi/goctl/templates/user/global"
	"github.com/zhufuyi/pkg/gin/middleware"
	"github.com/zhufuyi/pkg/gin/middleware/ratelimiter"
	"github.com/zhufuyi/pkg/gin/validator"
	"github.com/zhufuyi/pkg/logger"
)

type routerFn func(*gin.RouterGroup)

var routerFns []routerFn

// NewRouter 实例化路由
func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logging(middleware.WithLog(logger.Get())))
	r.Use(middleware.Cors())
	// ip, 自定义qps=10, burst=100
	r.Use(ratelimiter.QPS(
		ratelimiter.WithIP(),
		ratelimiter.WithQPS(10),
		ratelimiter.WithBurst(100),
	))
	binding.Validator = validator.Init()

	if global.Conf.IsEnableProfile {
		pprof.Register(r, "/"+global.Conf.ServerName)
	}

	apiv1 := r.Group("/api/v1")
	for _, fn := range routerFns {
		fn(apiv1)
	}

	return r
}