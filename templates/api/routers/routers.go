package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/zhufuyi/pkg/gin/middleware"
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
	binding.Validator = validator.Init()

	apiv1 := r.Group("/api/v1")
	for _, fn := range routerFns {
		fn(apiv1)
	}

	return r
}
