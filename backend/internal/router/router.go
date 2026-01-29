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
	userController := controller.NewUserController()
	serviceController := controller.NewServiceController()
	appController := controller.NewAppController()
	dashboardController := controller.NewDashboradController()

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

			// 用户路由
			users := protected.Group("/users")
			{
				users.GET("/lists", userController.ListUsers)
				users.GET("/detail/", userController.GetUser)
				users.POST("/update/", userController.UpdateUser)
				users.POST("/delete/", userController.DeleteUser)
			}

			// app租户应用路由
			apps := protected.Group("/apps")
			{
				apps.GET("/list", appController.GetAppList)
				apps.GET("/detail", appController.GetAppDetail)
				apps.POST("/create", appController.CreateApp)
				apps.POST("/update", appController.UpdateApp)
				apps.POST("/delete", appController.DeleteApp)

				// 增强的应用管理功能
				apps.GET("/stats", appController.GetAppStats)
			}

			// 服务路由
			services := protected.Group("/services")
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

			// 仪表盘
			dashboard := protected.Group("/dashboard")
			{
				//最上方基础的服务和app统计数
				dashboard.GET("/panelGruopData", dashboardController.PanelGruopData)
				//flowstat请求数对比,昨日今日请求数
				dashboard.GET("/flowstat", dashboardController.FlowStat)
				//仪表盘各个服务type分类总数
				dashboard.GET("/serviceStat", dashboardController.ServiceStat)
			}
		}
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{"status": "ok"})
	})

	return r
}
