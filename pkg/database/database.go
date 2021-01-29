package database

import (
	"database/sql"
	"goblog/pkg/logger"
	"time"

	"github.com/go-sql-driver/mysql"
)

// DB 数据库对象
var DB *sql.DB

// Initialize 初始化数据库
func Initialize() {
	initDB()
	createTables()
}

func initDB() {
	var err error
	config := mysql.Config{
		User:                 "root",
		Passwd:               "12345678",
		Addr:                 "127.0.0.1:3306",
		Net:                  "tcp",
		DBName:               "goblog",
		AllowNativePasswords: true,
	}

	// 准备数据库连接池
	// DSN => Data Source Name
	DB, err = sql.Open("mysql", config.FormatDSN())
	logger.LogError(err)

	// 设置最大连接数
	DB.SetMaxOpenConns(100)
	// 设置最大空闲连接数
	DB.SetMaxIdleConns(25)

	// 设置每个链接的过期时间
	DB.SetConnMaxLifetime(5 * time.Minute)

	// 尝试连接，失败报错
	err = DB.Ping()
	logger.LogError(err)
}

func createTables() {
	createArticlesSQL := `create table if not exists articles(
		id bigint(20) primary key auto_increment not null,
		title varchar(255) collate utf8mb4_unicode_ci not null,
		body longtext collate utf8mb4_unicode_ci
	);
	`
	_, err := DB.Exec(createArticlesSQL)
	logger.LogError(err)
}
