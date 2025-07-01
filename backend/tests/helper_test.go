package main

import (
	"fmt"
	"log"
	"testing"

	"todolist/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// initTestDB 初始化测试数据库连接
func initTestDB(t *testing.T) *gorm.DB {
	// 加载配置
	err := config.LoadConfig("../config/config.yaml")
	if err != nil {
		t.Fatalf("加载配置失败: %v", err)
	}

	// 获取MySQL配置
	mysqlConfig := config.GlobalConfig.MySQL

	// 使用测试数据库
	mysqlConfig.Database = fmt.Sprintf("%s_test", mysqlConfig.Database)

	// 配置GORM
	gormConfig := &gorm.Config{}

	// 连接数据库
	dsn := mysqlConfig.DSN()
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		t.Fatalf("连接数据库失败: %v", err)
	}

	log.Printf("成功连接到测试数据库: %s\n", mysqlConfig.Database)
	return db
}

// cleanupTestDB 清理测试数据
func cleanupTestDB(t *testing.T, db *gorm.DB) {
	// 删除所有任务
	err := db.Exec("DELETE FROM tasks").Error
	if err != nil {
		t.Fatalf("清理任务表失败: %v", err)
	}

	// 删除所有用户
	err = db.Exec("DELETE FROM users").Error
	if err != nil {
		t.Fatalf("清理用户表失败: %v", err)
	}
}
