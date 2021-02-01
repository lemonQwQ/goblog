package bootstrap

import (
	"goblog/app/models/user"
	"goblog/pkg/model"
	"goblog/pkg/models/article"
	"time"

	"gorm.io/gorm"
)

// SetupDB 初始化数据库和 ORM
func SetupDB() {

	// 连接并设置数据库
	db := model.ConnectDB()

	
	sqlDB, _ := db.DB()

	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// 创建和维护数据表结构
	migration(db)
}

func migration(db *gorm.DB) {
	db.AutoMigrate(
		&user.User{},
		&article.Article{},
	)
}
