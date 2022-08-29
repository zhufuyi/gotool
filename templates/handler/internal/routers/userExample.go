package routers

import (
	"github.com/zhufuyi/goctl/templates/handler/internal/handler"

	"github.com/gin-gonic/gin"
)

func userExampleRouter(group *gin.RouterGroup, h handler.UserExampleHandler) {
	group.POST("/userExample", h.Create)
	group.DELETE("/userExample/:id", h.DeleteByID)
	group.PUT("/userExample/:id", h.UpdateByID)
	group.GET("/userExample/:id", h.GetByID)
	group.POST("/userExamples", h.List) // 通过post任意列组合查询
}
