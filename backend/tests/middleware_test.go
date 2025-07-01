package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"todolist/internal/middleware"
	"todolist/pkg/jwt"
)

func TestAuthMiddleware(t *testing.T) {
	// 设置测试模式
	gin.SetMode(gin.TestMode)

	// 创建测试用例
	tests := []struct {
		name       string
		setupAuth  func(r *http.Request)
		checkAuth  func(t *testing.T, rec *httptest.ResponseRecorder)
		expectCode int
	}{
		{
			name: "无认证头",
			setupAuth: func(r *http.Request) {
				// 不设置认证头
			},
			checkAuth: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, rec.Code)
			},
			expectCode: http.StatusUnauthorized,
		},
		{
			name: "认证头格式错误",
			setupAuth: func(r *http.Request) {
				r.Header.Set("Authorization", "InvalidFormat")
			},
			checkAuth: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, rec.Code)
			},
			expectCode: http.StatusUnauthorized,
		},
		{
			name: "无效的token",
			setupAuth: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer invalid.token.here")
			},
			checkAuth: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, rec.Code)
			},
			expectCode: http.StatusUnauthorized,
		},
		{
			name: "有效的token",
			setupAuth: func(r *http.Request) {
				token, _ := jwt.GenerateToken(1, "testuser")
				r.Header.Set("Authorization", "Bearer "+token)
			},
			checkAuth: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rec.Code)
			},
			expectCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试路由
			r := gin.New()
			r.Use(middleware.AuthMiddleware())
			r.GET("/test", func(c *gin.Context) {
				userID := middleware.GetUserID(c)
				username := middleware.GetUsername(c)
				c.JSON(http.StatusOK, gin.H{
					"user_id":  userID,
					"username": username,
				})
			})

			// 创建测试请求
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			tt.setupAuth(req)

			// 创建响应记录器
			rec := httptest.NewRecorder()

			// 执行请求
			r.ServeHTTP(rec, req)

			// 检查结果
			tt.checkAuth(t, rec)
		})
	}
}

func TestGetUserFunctions(t *testing.T) {
	// 设置测试模式
	gin.SetMode(gin.TestMode)

	t.Run("测试获取用户信息函数", func(t *testing.T) {
		// 创建gin上下文
		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		// 测试未设置用户信息的情况
		assert.Equal(t, 0, middleware.GetUserID(c))
		assert.Equal(t, "", middleware.GetUsername(c))

		// 设置用户信息
		c.Set(middleware.ContextKeyUserID, 1)
		c.Set(middleware.ContextKeyUsername, "testuser")

		// 测试正常获取
		assert.Equal(t, 1, middleware.GetUserID(c))
		assert.Equal(t, "testuser", middleware.GetUsername(c))
	})

	t.Run("测试Must函数", func(t *testing.T) {
		// 创建gin上下文
		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		// 测试未设置用户信息时的panic
		assert.Panics(t, func() {
			middleware.MustGetUserID(c)
		})
		assert.Panics(t, func() {
			middleware.MustGetUsername(c)
		})

		// 设置用户信息
		c.Set(middleware.ContextKeyUserID, 1)
		c.Set(middleware.ContextKeyUsername, "testuser")

		// 测试正常获取
		assert.NotPanics(t, func() {
			userID := middleware.MustGetUserID(c)
			assert.Equal(t, 1, userID)
		})
		assert.NotPanics(t, func() {
			username := middleware.MustGetUsername(c)
			assert.Equal(t, "testuser", username)
		})
	})
}
