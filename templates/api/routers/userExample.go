package routers

import (
	v1 "github.com/zhufuyi/goctl/templates/api/routers/api/v1"

	"github.com/gin-gonic/gin"
)

func init() {
	routerFns = append(routerFns, userExampleRouter)
}

func userExampleRouter(group *gin.RouterGroup) {
	group.POST("/userExample", v1.CreateUserExample)
	//group.DELETE("/userExample", v1.DeleteUserExample)
	group.DELETE("/userExamples", v1.DeleteUserExamples)
	group.PUT("/userExample", v1.UpdateUserExample)
	group.GET("/userExample/:id", v1.GetUserExample)
	group.GET("/userExamples", v1.GetUserExamples)
	group.POST("/userExamples", v1.GetUserExamples2) // 通过post查询多条记录
}
