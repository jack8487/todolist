package main

import (
	"testing"
	"time"

	"todolist/config"
)

// 测试配置加载
func TestLoadConfig(t *testing.T) {
	// 1. 测试配置文件加载
	err := config.LoadConfig("../config/config.yaml")
	if err != nil {
		t.Fatalf("加载配置文件失败: %v", err)
	}

	// 2. 测试服务器配置
	t.Run("测试服务器配置", func(t *testing.T) {
		if config.GlobalConfig.Server.Port != 8080 {
			t.Errorf("服务器端口配置错误: 期望 8080, 实际 %d", config.GlobalConfig.Server.Port)
		}
		if config.GlobalConfig.Server.Mode != "debug" {
			t.Errorf("服务器模式配置错误: 期望 debug, 实际 %s", config.GlobalConfig.Server.Mode)
		}
	})

	// 3. 测试MySQL配置
	t.Run("测试MySQL配置", func(t *testing.T) {
		mysql := config.GlobalConfig.MySQL
		if mysql.Host != "localhost" {
			t.Errorf("MySQL主机配置错误: 期望 localhost, 实际 %s", mysql.Host)
		}
		if mysql.Port != 3306 {
			t.Errorf("MySQL端口配置错误: 期望 3306, 实际 %d", mysql.Port)
		}
		if mysql.Username != "root" {
			t.Errorf("MySQL用户名配置错误: 期望 root, 实际 %s", mysql.Username)
		}
		if mysql.Password != "1234" {
			t.Errorf("MySQL密码配置错误: 期望 1234, 实际 %s", mysql.Password)
		}
		if mysql.Database != "todolist" {
			t.Errorf("MySQL数据库配置错误: 期望 todolist, 实际 %s", mysql.Database)
		}

		// 测试DSN格式
		expectedDSN := "root:1234@tcp(localhost:3306)/todolist?charset=utf8mb4&parseTime=true&loc=Local"
		if dsn := mysql.DSN(); dsn != expectedDSN {
			t.Errorf("MySQL DSN格式错误:\n期望: %s\n实际: %s", expectedDSN, dsn)
		}
	})

	// 4. 测试Redis配置
	t.Run("测试Redis配置", func(t *testing.T) {
		redis := config.GlobalConfig.Redis
		if redis.Host != "localhost" {
			t.Errorf("Redis主机配置错误: 期望 localhost, 实际 %s", redis.Host)
		}
		if redis.Port != 6379 {
			t.Errorf("Redis端口配置错误: 期望 6379, 实际 %d", redis.Port)
		}

		// 测试Redis地址格式
		expectedAddr := "localhost:6379"
		if addr := redis.Addr(); addr != expectedAddr {
			t.Errorf("Redis地址格式错误: 期望 %s, 实际 %s", expectedAddr, addr)
		}
	})

	// 5. 测试JWT配置
	t.Run("测试JWT配置", func(t *testing.T) {
		jwt := config.GlobalConfig.JWT
		if jwt.SecretKey != "jack" {
			t.Errorf("JWT密钥配置错误: 期望 jack, 实际 %s", jwt.SecretKey)
		}
		if jwt.ExpireHours != 24*time.Hour {
			t.Errorf("JWT过期时间配置错误: 期望 24h, 实际 %v", jwt.ExpireHours)
		}
		if jwt.Issuer != "todolist" {
			t.Errorf("JWT签发者配置错误: 期望 todolist, 实际 %s", jwt.Issuer)
		}
	})

	// 6. 测试日志配置
	t.Run("测试日志配置", func(t *testing.T) {
		log := config.GlobalConfig.Log
		if log.Level != "debug" {
			t.Errorf("日志级别配置错误: 期望 debug, 实际 %s", log.Level)
		}
		if log.Filename != "logs/todolist.log" {
			t.Errorf("日志文件名配置错误: 期望 logs/todolist.log, 实际 %s", log.Filename)
		}
		if !log.Compress {
			t.Error("日志压缩配置错误: 期望 true, 实际 false")
		}
	})
}
