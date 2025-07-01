package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"todolist/pkg/jwt"
)

const (
	// ContextKeyUserID 用户ID的上下文键
	ContextKeyUserID = "user_id"
	// ContextKeyUsername 用户名的上下文键
	ContextKeyUsername = "username"
)

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization 头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未提供认证信息",
			})
			c.Abort()
			return
		}

		// 处理 Bearer 前缀
		tokenString := authHeader
		if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			tokenString = authHeader[7:]
		}

		// 验证 token
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "无效的认证信息",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}

// GetUserID 从上下文中获取用户ID
func GetUserID(c *gin.Context) int {
	userID, _ := c.Get("user_id")
	return userID.(int)
}

// GetUsername 从上下文中获取用户名
func GetUsername(c *gin.Context) string {
	username, exists := c.Get(ContextKeyUsername)
	if !exists {
		return ""
	}
	return username.(string)
}

// MustGetUserID 从上下文中获取用户ID，如果不存在则panic
func MustGetUserID(c *gin.Context) int {
	userID, exists := c.Get(ContextKeyUserID)
	if !exists {
		panic("用户ID未找到在上下文中")
	}
	return userID.(int)
}

// MustGetUsername 从上下文中获取用户名，如果不存在则panic
func MustGetUsername(c *gin.Context) string {
	username, exists := c.Get(ContextKeyUsername)
	if !exists {
		panic("用户名未找到在上下文中")
	}
	return username.(string)
}
