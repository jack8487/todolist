package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"todolist/internal/service"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Register godoc
// @Summary 用户注册
// @Description 注册新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "注册信息"
// @Success 200 {object} Response{data=string} "注册成功"
// @Failure 400 {object} Response{} "请求参数错误"
// @Failure 500 {object} Response{} "服务器内部错误"
// @Router /users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "请求参数错误",
			Error:   err.Error(),
		})
		return
	}

	err := h.userService.Register(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "注册失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "注册成功",
	})
}

// Login godoc
// @Summary 用户登录
// @Description 用户登录并获取令牌
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录信息"
// @Success 200 {object} Response{data=string} "登录成功"
// @Failure 400 {object} Response{} "请求参数错误"
// @Failure 500 {object} Response{} "服务器内部错误"
// @Router /users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "请求参数错误",
			Error:   err.Error(),
		})
		return
	}

	token, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "登录失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "登录成功",
		Data:    token,
	})
}

// UpdatePassword godoc
// @Summary 更新密码
// @Description 更新用户密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body UpdatePasswordRequest true "密码更新信息"
// @Success 200 {object} Response{} "密码更新成功"
// @Failure 400 {object} Response{} "请求参数错误"
// @Failure 401 {object} Response{} "未授权"
// @Failure 500 {object} Response{} "服务器内部错误"
// @Router /users/password [put]
func (h *UserHandler) UpdatePassword(c *gin.Context) {
	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "请求参数错误",
			Error:   err.Error(),
		})
		return
	}

	userID := c.GetInt("user_id")
	err := h.userService.UpdatePassword(userID, req.OldPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "更新密码失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "更新密码成功",
	})
}

// RegisterRoutes 注册路由
func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	users := r.Group("/api/users")
	{
		users.POST("/register", h.Register)
		users.POST("/login", h.Login)
		users.PUT("/password", h.UpdatePassword)
	}
}

// Response API 响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdatePasswordRequest 更新密码请求
type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}
