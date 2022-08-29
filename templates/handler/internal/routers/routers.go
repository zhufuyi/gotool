package routers

import (
	"net/http"

	"github.com/zhufuyi/goctl/templates/handler/config"
	"github.com/zhufuyi/goctl/templates/handler/docs"
	"github.com/zhufuyi/goctl/templates/handler/internal/handler"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/zhufuyi/pkg/gin/middleware"
	"github.com/zhufuyi/pkg/gin/middleware/metrics"
	"github.com/zhufuyi/pkg/gin/middleware/ratelimiter"
	"github.com/zhufuyi/pkg/gin/validator"
	"github.com/zhufuyi/pkg/logger"
)

// NewRouter 实例化路由
func NewRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.Cors())

	// request id
	r.Use(middleware.RequestID())

	// 日志
	r.Use(middleware.Logging(
		middleware.WithLog(logger.Get()),
		middleware.WithIgnoreRoutes("/metrics"), // 忽略路由
	))

	// 指标 metrics
	if config.Get().EnableMetrics {
		r.Use(metrics.Metrics(r,
			metrics.WithIgnoreStatusCodes(http.StatusNotFound), // 忽略404状态码
		))
	}

	// 限流 limit
	if config.Get().EnableLimit {
		opts := []ratelimiter.Option{
			ratelimiter.WithQPS(config.Get().Ratelimiter.QPSLimit),
			ratelimiter.WithBurst(config.Get().Ratelimiter.MaxLimit),
		}
		if config.Get().Ratelimiter.IsIP() {
			opts = append(opts, ratelimiter.WithIP())
		}
		r.Use(ratelimiter.QPS(opts...))
	}

	// 链路跟踪 trace
	if config.Get().EnableTracing {
		r.Use(middleware.Tracing(config.Get().ServiceName))
	}

	// 性能分析 profile
	if config.Get().EnableProfile {
		pprof.Register(r)
	}

	// 校验器
	binding.Validator = validator.Init()

	// 注册swagger路由，通过swag init生成代码
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiV1 := r.Group("/api/v1")

	// 注册路由
	userExampleRouter(apiV1, handler.NewUserExampleHandler())

	return r
}
