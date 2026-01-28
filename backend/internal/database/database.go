package database

import (
	"fmt"
	"gateway-service/internal/config"
	"gateway-service/internal/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)
	//关闭对mysql外键的创建
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, err
	}
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate tables
	err = db.AutoMigrate(
		&models.User{},
		&models.ServiceInfo{},
		&models.AccessControl{},
		&models.LoadBalance{},
		&models.HttpRule{},
		&models.TcpRule{},
		&models.GrpcRule{},
		&models.App{},
	)
	if err != nil {
		return nil, err
	}

	log.Println("Database connected and migrated successfully")

	DB = db
	return db, nil
}
