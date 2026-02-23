package http_proxy_middleware

import (
	"gateway-service/internal/models"
	"gateway-service/internal/reverse_proxy"
	"gateway-service/pkg/response"
	"github.com/gin-gonic/gin"
)

// 匹配接入方式 基于请求信息
func HTTPReverseProxyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			response.Error(c, "service not found")
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*models.ServiceDetail)

		//各个app服务的负载均衡器中间件
		lb, err := models.LoadBalancerHandler.GetLoadBalance(serviceDetail)
		if err != nil {
			response.Error(c, "LoadBalancerHandler  ERROR")
			c.Abort()
			return
		}

		//各个app服务的连接池中间件
		trans, err := models.TransportorHandler.GetTrans(serviceDetail)
		if err != nil {
			response.Error(c, "TransporterHandler ERROR")
			c.Abort()
			return
		}
		//middleware.ResponseSuccess(c,"ok")
		//return
		//创建 reverseproxy
		//使用 reverseproxy.ServerHTTP(c.Request,c.Response)
		proxy := reverse_proxy.NewLoadBalanceReverseProxy(c, lb, trans)
		proxy.ServeHTTP(c.Writer, c.Request)
		c.Abort()
		return
	}
}
