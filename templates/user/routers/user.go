package routers

import (
	"github.com/zhufuyi/goctl/templates/user/routers/api/v1"

	"github.com/gin-gonic/gin"
	"github.com/zhufuyi/goctl/utils/middleware"
)

func init() {
	routerFns = append(routerFns, userRouter)
}

func userRouter(group *gin.RouterGroup) {
	// 不需要鉴权api
	group.POST("/auth/register", v1.Register)    // 注册
	group.GET("/auth/activate", v1.ActivateUser) // 注册确认，通过收件邮箱点击链接
	group.POST("/auth/login", v1.Login)          // 登录

	// 需要鉴权api
	group.POST("/user/logout/:id", middleware.Auth(), v1.Logout) // 登出
	group.PUT("/user/:id", middleware.Auth(), v1.UpdateUser)
	group.GET("/user/:id", middleware.Auth(), v1.GetUser)
	group.DELETE("/user/:id", middleware.Auth(), v1.DeleteUser)

	// 管理员鉴权
	group.GET("/users", middleware.AuthAdmin(), v1.GetUsers)
	group.POST("/users", middleware.AuthAdmin(), v1.GetUsers2) // 通过post查询多条记录
}
