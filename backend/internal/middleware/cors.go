package middleware

import (
	"github.com/gin-gonic/gin"
)

// 允许的源列表
var allowedOrigins = []string{
	"http://localhost:5173",
	"http://localhost:5174",
	"http://localhost:3000",
	// 可以添加更多允许的源
	// "http://example.com",
	// "https://example.com",
}

// CORS 中间件用于处理跨域请求
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 检查请求的源是否在允许列表中
		allowOrigin := ""
		for _, allowed := range allowedOrigins {
			if allowed == origin {
				allowOrigin = origin
				break
			}
		}

		// 如果源在允许列表中，设置对应的头
		if allowOrigin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		}

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
