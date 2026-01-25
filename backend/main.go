package main

import (
	"gateway-service/internal/config"
	"gateway-service/internal/database"
	"gateway-service/internal/models"
	"gateway-service/internal/router"
	"log"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 设置全局 DB
	models.SetDB(db)

	// 设置路由
	r := router.SetupRouter(cfg)

	// 启动服务器
	addr := ":" + cfg.Server.Port
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
