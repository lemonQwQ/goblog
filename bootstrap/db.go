package bootstrap

import (
	"goblog/pkg/model"
	"time"
)

// SetupDB 初始化数据库和 ORM
func SetupDB() {

	db := model.ConnectDB()

	sqlDB, _ := db.DB()

	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
}
