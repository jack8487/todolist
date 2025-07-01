package api

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"todolist/internal/middleware"
	"todolist/internal/model"
	"todolist/internal/service"
)

// TaskHandler 任务处理器
type TaskHandler struct {
	taskService service.TaskService
}

// NewTaskHandler 创建任务处理器
func NewTaskHandler(taskService service.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

// parseDateString 解析各种格式的日期字符串
func parseDateString(dateStr string) (time.Time, error) {
	// 预处理日期字符串
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return time.Time{}, nil
	}

	// 标准化分隔符
	re := regexp.MustCompile(`[./\-]`)
	dateStr = re.ReplaceAllString(dateStr, "-")

	// 移除多余的空格
	dateStr = regexp.MustCompile(`\s+`).ReplaceAllString(dateStr, " ")

	// 尝试解析不同格式的日期
	formats := []string{
		"2006-01-02",                // YYYY-MM-DD
		"2006-1-2",                  // YYYY-M-D
		"06-01-02",                  // YY-MM-DD
		"06-1-2",                    // YY-M-D
		"2006-01-02 15:04:05",       // YYYY-MM-DD HH:mm:ss
		"2006-01-02 15:04",          // YYYY-MM-DD HH:mm
		"2006-01-02T15:04:05Z07:00", // ISO 8601
		"01-02-2006",                // MM-DD-YYYY
		"1-2-2006",                  // M-D-YYYY
		"January 2, 2006",           // Month D, YYYY
		"Jan 2, 2006",               // Mon D, YYYY
		"2 Jan 2006",                // D Mon YYYY
		"02 Jan 06",                 // DD Mon YY
		"2006年01月02日",               // 中文格式
		"2006年1月2日",                 // 中文格式（简化）
	}

	// 处理中文日期格式
	if strings.Contains(dateStr, "年") {
		dateStr = strings.ReplaceAll(dateStr, "年", "-")
		dateStr = strings.ReplaceAll(dateStr, "月", "-")
		dateStr = strings.ReplaceAll(dateStr, "日", "")
	}

	var lastErr error
	for _, format := range formats {
		t, err := time.Parse(format, dateStr)
		if err == nil {
			// 处理两位数年份
			if t.Year() < 100 {
				if t.Year() < 50 {
					t = t.AddDate(2000, 0, 0)
				} else {
					t = t.AddDate(1900, 0, 0)
				}
			}

			// 验证日期是否合理
			if t.Year() < 1900 || t.Year() > 9999 {
				continue
			}

			// 验证是否为未来100年内的日期
			now := time.Now()
			if t.After(now.AddDate(100, 0, 0)) {
				continue
			}

			return t, nil
		}
		lastErr = err
	}

	// 尝试解析时间戳
	if timestamp, err := strconv.ParseInt(dateStr, 10, 64); err == nil {
		// 处理秒级时间戳
		if len(dateStr) == 10 {
			return time.Unix(timestamp, 0), nil
		}
		// 处理毫秒级时间戳
		if len(dateStr) == 13 {
			return time.Unix(timestamp/1000, 0), nil
		}
	}

	if lastErr != nil {
		return time.Time{}, lastErr
	}
	return time.Time{}, nil
}

// Create godoc
// @Summary 创建任务
// @Description 创建新任务
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body CreateTaskRequest true "任务信息"
// @Success 200 {object} Response{data=model.Task} "创建成功"
// @Failure 400 {object} Response{} "请求参数错误"
// @Failure 401 {object} Response{} "未授权"
// @Failure 500 {object} Response{} "服务器内部错误"
// @Router /tasks [post]
func (h *TaskHandler) Create(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "请求参数错误",
			Error:   err.Error(),
		})
		return
	}

	var dueDate time.Time
	if req.DueDate != "" {
		var err error
		dueDate, err = parseDateString(req.DueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Code:    400,
				Message: "日期格式错误",
				Error:   "支持的日期格式：YYYY-MM-DD、YYYY/MM/DD、YYYY年MM月DD日、MM/DD/YYYY等",
			})
			return
		}
	}

	task := &model.Task{
		UserID:      middleware.GetUserID(c),
		Title:       req.Title,
		Description: req.Description,
		DueDate:     dueDate,
		Status:      model.TaskStatusTodo,
	}

	if err := h.taskService.Create(task); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "创建任务失败",
			Error:   err.Error(),
		})
		return
	}

	// 转换状态为文本形式
	response := struct {
		*model.Task
		Status string `json:"status"`
	}{
		Task:   task,
		Status: task.GetStatusText(),
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "创建任务成功",
		Data:    response,
	})
}

// Update godoc
// @Summary 更新任务
// @Description 更新任务信息
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "任务ID"
// @Param request body UpdateTaskRequest true "任务信息"
// @Success 200 {object} Response{data=model.Task} "更新成功"
// @Failure 400 {object} Response{} "请求参数错误"
// @Failure 401 {object} Response{} "未授权"
// @Failure 500 {object} Response{} "服务器内部错误"
// @Router /tasks/{id} [put]
func (h *TaskHandler) Update(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的任务ID",
		})
		return
	}

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "请求参数错误",
			Error:   err.Error(),
		})
		return
	}

	var dueDate time.Time
	if req.DueDate != "" {
		var err error
		dueDate, err = parseDateString(req.DueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Code:    400,
				Message: "日期格式错误",
				Error:   "支持的日期格式：YYYY-MM-DD、YYYY/MM/DD、YYYY年MM月DD日、MM/DD/YYYY等",
			})
			return
		}
	}

	task := &model.Task{
		ID:          taskID,
		UserID:      middleware.GetUserID(c),
		Title:       req.Title,
		Description: req.Description,
		DueDate:     dueDate,
	}

	// 设置状态
	if req.Status != "" {
		task.SetStatusFromText(req.Status)
	}

	if err := h.taskService.Update(task); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "更新任务失败",
			Error:   err.Error(),
		})
		return
	}

	// 转换状态为文本形式
	response := struct {
		*model.Task
		Status string `json:"status"`
	}{
		Task:   task,
		Status: task.GetStatusText(),
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "更新任务成功",
		Data:    response,
	})
}

// Delete godoc
// @Summary 删除任务
// @Description 删除指定任务
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "任务ID"
// @Success 200 {object} Response{} "删除成功"
// @Failure 400 {object} Response{} "请求参数错误"
// @Failure 401 {object} Response{} "未授权"
// @Failure 500 {object} Response{} "服务器内部错误"
// @Router /tasks/{id} [delete]
func (h *TaskHandler) Delete(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的任务ID",
		})
		return
	}

	if err := h.taskService.Delete(taskID, middleware.GetUserID(c)); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "删除任务失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "删除任务成功",
	})
}

// Get godoc
// @Summary 获取任务详情
// @Description 获取指定任务的详细信息
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "任务ID"
// @Success 200 {object} Response{data=model.Task} "获取成功"
// @Failure 400 {object} Response{} "请求参数错误"
// @Failure 401 {object} Response{} "未授权"
// @Failure 500 {object} Response{} "服务器内部错误"
// @Router /tasks/{id} [get]
func (h *TaskHandler) Get(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的任务ID",
		})
		return
	}

	task, err := h.taskService.Get(taskID, middleware.GetUserID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "获取任务失败",
			Error:   err.Error(),
		})
		return
	}

	// 转换状态为文本形式
	response := struct {
		*model.Task
		Status string `json:"status"`
	}{
		Task:   task,
		Status: task.GetStatusText(),
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取任务成功",
		Data:    response,
	})
}

// List godoc
// @Summary 获取任务列表
// @Description 获取当前用户的任务列表
// @Tags 任务管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param status query string false "任务状态" Enums(todo,in_progress,done)
// @Success 200 {object} Response{data=ListTasksResponse} "获取成功"
// @Failure 401 {object} Response{} "未授权"
// @Failure 500 {object} Response{} "服务器内部错误"
// @Router /tasks [get]
func (h *TaskHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	tasks, total, err := h.taskService.List(middleware.GetUserID(c), status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "获取任务列表失败",
			Error:   err.Error(),
		})
		return
	}

	// 转换状态为文本形式
	var responseTasks []struct {
		*model.Task
		Status string `json:"status"`
	}
	for _, task := range tasks {
		responseTasks = append(responseTasks, struct {
			*model.Task
			Status string `json:"status"`
		}{
			Task:   task,
			Status: task.GetStatusText(),
		})
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取任务列表成功",
		Data: ListTasksResponse{
			Total: total,
			Items: responseTasks,
		},
	})
}

// RegisterRoutes 注册路由
func (h *TaskHandler) RegisterRoutes(r *gin.Engine) {
	tasks := r.Group("/api/tasks")
	tasks.Use(middleware.AuthMiddleware())
	{
		tasks.POST("", h.Create)
		tasks.PUT("/:id", h.Update)
		tasks.DELETE("/:id", h.Delete)
		tasks.GET("/:id", h.Get)
		tasks.GET("", h.List)
	}
}

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=100"`
	Description string `json:"description" binding:"max=500"`
	DueDate     string `json:"due_date"` // 移除 datetime 验证，我们将手动验证
}

// UpdateTaskRequest 更新任务请求
type UpdateTaskRequest struct {
	Title       string `json:"title" binding:"omitempty,min=1,max=100"`
	Description string `json:"description" binding:"max=500"`
	Status      string `json:"status" binding:"omitempty,oneof=todo in_progress done"`
	DueDate     string `json:"due_date"` // 移除 datetime 验证，我们将手动验证
}

// ListTasksResponse 任务列表响应
type ListTasksResponse struct {
	Total int64       `json:"total"`
	Items interface{} `json:"items"`
}
