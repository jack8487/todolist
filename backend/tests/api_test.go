package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"todolist/internal/api"
	"todolist/internal/model"
	"todolist/pkg/jwt"
)

// MockTaskService 模拟任务服务
type MockTaskService struct {
	mock.Mock
}

func (m *MockTaskService) Create(task *model.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskService) Update(task *model.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskService) Delete(taskID, userID int) error {
	args := m.Called(taskID, userID)
	return args.Error(0)
}

func (m *MockTaskService) Get(taskID, userID int) (*model.Task, error) {
	args := m.Called(taskID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Task), args.Error(1)
}

func (m *MockTaskService) List(userID int, status string, page, pageSize int) ([]*model.Task, int64, error) {
	args := m.Called(userID, status, page, pageSize)
	return args.Get(0).([]*model.Task), args.Get(1).(int64), args.Error(2)
}

func setupTestRouter(taskService *MockTaskService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := api.NewTaskHandler(taskService)
	handler.RegisterRoutes(r)

	return r
}

func TestTaskHandler_Create(t *testing.T) {
	taskService := new(MockTaskService)
	router := setupTestRouter(taskService)

	tests := []struct {
		name       string
		reqBody    map[string]interface{}
		setupAuth  func(r *http.Request)
		setupMock  func()
		wantStatus int
	}{
		{
			name: "创建成功",
			reqBody: map[string]interface{}{
				"title":       "Test Task",
				"description": "Test Description",
				"due_date":    time.Now().Add(24 * time.Hour).Format("2006-01-02"),
			},
			setupAuth: func(r *http.Request) {
				token, _ := jwt.GenerateToken(1, "testuser")
				r.Header.Set("Authorization", "Bearer "+token)
			},
			setupMock: func() {
				taskService.On("Create", mock.AnythingOfType("*model.Task")).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "无效的日期格式",
			reqBody: map[string]interface{}{
				"title":       "Test Task",
				"description": "Test Description",
				"due_date":    "invalid-date",
			},
			setupAuth: func(r *http.Request) {
				token, _ := jwt.GenerateToken(1, "testuser")
				r.Header.Set("Authorization", "Bearer "+token)
			},
			setupMock:  func() {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			body, _ := json.Marshal(tt.reqBody)
			req := httptest.NewRequest(http.MethodPost, "/api/tasks", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			tt.setupAuth(req)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			taskService.AssertExpectations(t)
		})
	}
}

func TestTaskHandler_Update(t *testing.T) {
	taskService := new(MockTaskService)
	router := setupTestRouter(taskService)

	tests := []struct {
		name       string
		taskID     string
		reqBody    map[string]interface{}
		setupAuth  func(r *http.Request)
		setupMock  func()
		wantStatus int
	}{
		{
			name:   "更新成功",
			taskID: "1",
			reqBody: map[string]interface{}{
				"title":       "Updated Task",
				"description": "Updated Description",
				"status":      "done",
			},
			setupAuth: func(r *http.Request) {
				token, _ := jwt.GenerateToken(1, "testuser")
				r.Header.Set("Authorization", "Bearer "+token)
			},
			setupMock: func() {
				taskService.On("Update", mock.AnythingOfType("*model.Task")).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "无效的任务ID",
			taskID: "invalid",
			reqBody: map[string]interface{}{
				"title": "Updated Task",
			},
			setupAuth: func(r *http.Request) {
				token, _ := jwt.GenerateToken(1, "testuser")
				r.Header.Set("Authorization", "Bearer "+token)
			},
			setupMock:  func() {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			body, _ := json.Marshal(tt.reqBody)
			req := httptest.NewRequest(http.MethodPut, "/api/tasks/"+tt.taskID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			tt.setupAuth(req)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			taskService.AssertExpectations(t)
		})
	}
}

func TestTaskHandler_List(t *testing.T) {
	taskService := new(MockTaskService)
	router := setupTestRouter(taskService)

	tests := []struct {
		name       string
		query      string
		setupAuth  func(r *http.Request)
		setupMock  func()
		wantStatus int
	}{
		{
			name:  "获取列表成功",
			query: "?page=1&page_size=10&status=todo",
			setupAuth: func(r *http.Request) {
				token, _ := jwt.GenerateToken(1, "testuser")
				r.Header.Set("Authorization", "Bearer "+token)
			},
			setupMock: func() {
				tasks := []*model.Task{{ID: 1, Title: "Test Task"}}
				taskService.On("List", 1, "todo", 1, 10).Return(tasks, int64(1), nil)
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			req := httptest.NewRequest(http.MethodGet, "/api/tasks"+tt.query, nil)
			tt.setupAuth(req)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			taskService.AssertExpectations(t)
		})
	}
}
