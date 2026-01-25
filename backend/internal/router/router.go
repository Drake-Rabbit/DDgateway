package router

import (
	"gateway-service/internal/config"
	"gateway-service/internal/controller"
	"gateway-service/internal/middleware"
	"gateway-service/pkg/response"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置所有路由
func SetupRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORS())
	response.Setup(r)

	// 创建控制器实例
	authController := controller.NewAuthController(cfg)
	tenantController := controller.NewTenantController()
	userController := controller.NewUserController()
	serviceController := controller.NewServiceController()

	// API 路由
	api := r.Group("/api")
	{
		// 公开路由 - 认证
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}

		// 受保护的路由 - 需要认证
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			// 租户路由
			tenants := protected.Group("/tenants")
			tenants.Use(middleware.TenantMiddleware())
			{
				tenants.GET("/list", tenantController.ListTenants)
				tenants.GET("/detail/:id", tenantController.GetTenant)
				tenants.POST("/create", tenantController.CreateTenant)
				tenants.POST("/update", tenantController.UpdateTenant)
				tenants.POST("/delete", tenantController.DeleteTenant)
			}

			// 用户路由
			users := protected.Group("/users")
			users.Use(middleware.TenantMiddleware())
			{
				users.GET("/lists", userController.ListUsers)
				users.GET("/detail/", userController.GetUser)
				users.POST("/update/", userController.UpdateUser)
				users.POST("/delete/", userController.DeleteUser)
			}

			// 服务路由
			services := protected.Group("/services")
			services.Use(middleware.TenantMiddleware())
			{
				// 通用接口
				services.GET("/service_list", serviceController.ListServices)
				services.POST("/service_delete", serviceController.DeleteService)
				services.GET("/service_detail", serviceController.ServiceDetail)
				//服务统计
				services.GET("/service_stat", serviceController.ServiceStat)

				// HTTP 服务
				services.POST("/service_add_http", serviceController.CreateHTTP)
				services.POST("/service_update_http", serviceController.UpdateHTTP)

				// TCP 服务
				services.POST("/service_add_tcp", serviceController.CreateTcp)
				services.POST("/service_update_tcp", serviceController.UpdateTcp)

				// gRPC 服务
				services.POST("/service_add_grpc", serviceController.CreateGrpc)
				services.POST("/service_update_grpc", serviceController.UpdateGrpc)
			}
		}
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{"status": "ok"})
	})

	return r
}
