package repository

import (
	"fmt"
	"log"

	"todolist/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() error {
	var err error

	// 获取MySQL配置
	mysqlConfig := config.GlobalConfig.MySQL

	// 配置GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 开发环境下打印SQL
	}

	// 连接数据库
	dsn := mysqlConfig.DSN()
	DB, err = gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	// 获取底层的sqlDB
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取sqlDB失败: %v", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(mysqlConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(mysqlConfig.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(mysqlConfig.ConnMaxLifetime)

	log.Println("数据库连接成功")
	return nil
}
