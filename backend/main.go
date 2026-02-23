package main

import (
	"gateway-service/internal/config"
	"gateway-service/internal/database"
	"gateway-service/internal/http_proxy_router"
	"gateway-service/internal/models"
	"gateway-service/internal/router"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 加载配置
	cfg := config.Load()
	config.InitViperConfig()

	// 初始化mysql数据库
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 设置全局 DB
	models.SetDB(db)

	//初始化redis
	err = database.InitRedisClient()
	if err != nil {
		log.Fatal("Failed to connect to redis:", err)
	}

	// 后台dashboard设置路由
	r := router.SetupRouter(cfg)

	// 启动admin_dashboard管理面板
	go func() {
		addr := ":" + cfg.Server.Port
		log.Printf("Server starting on %s", addr)
		if err := r.Run(addr); err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}()

	//启动代理服务器
	err = models.ServiceManagerHandler.LoadOnce()
	if err != nil {
		log.Fatal("Failed to load service manager:", err)
	}

	go func() {
		http_proxy_router.HttpServerRun()
	}()

	//监听信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	//停止代理服务器
	http_proxy_router.HttpServerStop()
}
