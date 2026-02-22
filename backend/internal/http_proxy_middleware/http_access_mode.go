package http_proxy_middleware

import (
	"fmt"
	"gateway-service/internal/define"
	"gateway-service/internal/models"
	"gateway-service/pkg/response"
	"github.com/gin-gonic/gin"
)

// http请求的准入中间件
// 匹配接入方式 基于请求信息
func HTTPAccessModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//
		s, err := models.ServiceManagerHandler.HTTPAccessMode(c)
		if err != nil {
			response.Error(c, err.Error())
			c.Abort()
			return
		}
		// 设置上下文的服务信息
		fmt.Println("match service:", define.ObjToJson(s))
		c.Set("service", s)
		c.Next()
		// todo 获取服务规则
		// todo 获取客户端IP
		// todo 获取客户端信息
		// todo 获取客户端IP白名单
		// todo 获取客户端IP黑名单
		// todo 获取客户端IP访问限制
		// todo 获取客户端IP访问限制
	}
}
