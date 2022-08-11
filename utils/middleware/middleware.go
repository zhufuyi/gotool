package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zhufuyi/goctl/templates/user/global"
	"github.com/zhufuyi/pkg/gin/errcode"
	"github.com/zhufuyi/pkg/gin/render"
	"github.com/zhufuyi/pkg/jwt"
	"github.com/zhufuyi/pkg/logger"
)

// Auth 普通用户权限
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if global.Conf.IsUseAuth { // 控制是否鉴权
			authorization := c.GetHeader("Authorization")
			if len(authorization) < 20 {
				logger.Error("authorization is illegal", logger.String("authorization", authorization))
				render.Error(c, errcode.Unauthorized)
				c.Abort()
				return
			}
			token := authorization[7:] // 去掉Bearer 前缀
			claims, err := jwt.VerifyToken(token)
			if err != nil {
				logger.Error("VerifyToken error", logger.Err(err))
				render.Error(c, errcode.Unauthorized)
				c.Abort()
				return
			}
			c.Set("uid", claims.Uid)
		}

		c.Next()
	}
}

// AuthAdmin 管理员权限
func AuthAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if global.Conf.IsUseAuth { // 控制是否鉴权
			authorization := c.GetHeader("Authorization")
			if len(authorization) < 20 {
				logger.Error("authorization is illegal", logger.String("authorization", authorization))
				render.Error(c, errcode.Unauthorized)
				c.Abort()
				return
			}
			token := authorization[7:] // 去掉Bearer 前缀
			claims, err := jwt.VerifyToken(token)
			if err != nil {
				logger.Error("VerifyToken error", logger.Err(err))
				render.Error(c, errcode.Unauthorized)
				c.Abort()
				return
			}

			// 判断是否为管理员
			if claims.Role != "admin" {
				logger.Error("prohibition of access", logger.String("uid", claims.Uid), logger.String("role", claims.Role))
				render.Error(c, errcode.Forbidden)
				c.Abort()
				return
			}
			c.Set("uid", claims.Uid)
		}
		c.Next()
	}
}
