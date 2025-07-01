package main

import (
	"testing"
	"time"

	"todolist/config"
	"todolist/pkg/jwt"
)

func TestJWT(t *testing.T) {
	// 加载配置
	err := config.LoadConfig("../config/config.yaml")
	if err != nil {
		t.Fatalf("加载配置失败: %v", err)
	}

	// 测试生成令牌
	t.Run("测试生成令牌", func(t *testing.T) {
		token, err := jwt.GenerateToken(1, "test_user")
		if err != nil {
			t.Errorf("生成令牌失败: %v", err)
		}
		if token == "" {
			t.Error("生成的令牌为空")
		}
	})

	// 测试解析和验证令牌
	t.Run("测试解析和验证令牌", func(t *testing.T) {
		// 生成测试令牌
		token, _ := jwt.GenerateToken(1, "test_user")

		// 测试解析令牌
		claims, err := jwt.ParseToken(token)
		if err != nil {
			t.Errorf("解析令牌失败: %v", err)
		}
		if claims.UserID != 1 {
			t.Errorf("用户ID不匹配: 期望 1, 实际 %d", claims.UserID)
		}
		if claims.Username != "test_user" {
			t.Errorf("用户名不匹配: 期望 test_user, 实际 %s", claims.Username)
		}

		// 测试验证令牌
		if !jwt.ValidateToken(token) {
			t.Error("验证令牌失败")
		}
	})

	// 测试过期令牌
	t.Run("测试过期令牌", func(t *testing.T) {
		// 修改配置使令牌立即过期
		originalExpireHours := config.GlobalConfig.JWT.ExpireHours
		config.GlobalConfig.JWT.ExpireHours = -1 * time.Hour

		// 生成过期令牌
		token, _ := jwt.GenerateToken(1, "test_user")

		// 还原配置
		config.GlobalConfig.JWT.ExpireHours = originalExpireHours

		// 验证过期令牌
		if jwt.ValidateToken(token) {
			t.Error("验证过期令牌应该失败")
		}

		// 解析过期令牌
		_, err := jwt.ParseToken(token)
		if err != jwt.ErrTokenExpired {
			t.Errorf("期望错误 ErrTokenExpired, 实际得到: %v", err)
		}
	})

	// 测试无效令牌
	t.Run("测试无效令牌", func(t *testing.T) {
		invalidToken := "invalid.token.string"
		if jwt.ValidateToken(invalidToken) {
			t.Error("验证无效令牌应该失败")
		}

		_, err := jwt.ParseToken(invalidToken)
		if err != jwt.ErrTokenMalformed {
			t.Errorf("期望错误 ErrTokenMalformed, 实际得到: %v", err)
		}
	})

	// 测试从令牌获取信息
	t.Run("测试从令牌获取信息", func(t *testing.T) {
		token, _ := jwt.GenerateToken(1, "test_user")

		// 测试获取用户ID
		userID, err := jwt.GetUserIDFromToken(token)
		if err != nil {
			t.Errorf("获取用户ID失败: %v", err)
		}
		if userID != 1 {
			t.Errorf("用户ID不匹配: 期望 1, 实际 %d", userID)
		}

		// 测试获取用户名
		username, err := jwt.GetUsernameFromToken(token)
		if err != nil {
			t.Errorf("获取用户名失败: %v", err)
		}
		if username != "test_user" {
			t.Errorf("用户名不匹配: 期望 test_user, 实际 %s", username)
		}
	})
}
