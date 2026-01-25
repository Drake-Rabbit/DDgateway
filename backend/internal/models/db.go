package models

import "gorm.io/gorm"

var DB *gorm.DB

// SetDB 设置全局数据库连接
func SetDB(db *gorm.DB) {
	DB = db
}
