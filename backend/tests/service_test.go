package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"todolist/internal/model"
	"todolist/internal/repository"
	"todolist/internal/service"
)

func setupTestService(t *testing.T) (service.UserService, service.TaskService) {
	// 初始化测试数据库
	db := initTestDB(t)
	defer cleanupTestDB(t, db)

	// 创建仓储实例
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	// 创建服务实例
	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo)

	return userService, taskService
}

func TestUserService(t *testing.T) {
	userService, _ := setupTestService(t)

	// 测试用户注册
	t.Run("测试用户注册", func(t *testing.T) {
		err := userService.Register("test_user", "password123")
		assert.NoError(t, err)
	})

	// 测试用户登录
	t.Run("测试用户登录", func(t *testing.T) {
		token, err := userService.Login("test_user", "password123")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	// 测试错误密码登录
	t.Run("测试错误密码登录", func(t *testing.T) {
		_, err := userService.Login("test_user", "wrong_password")
		assert.Error(t, err)
	})
}

func TestTaskService(t *testing.T) {
	_, taskService := setupTestService(t)

	// 创建测试任务
	task := &model.Task{
		UserID:      1,
		Title:       "Test Task",
		Description: "Test Description",
		Status:      model.TaskStatusTodo,
		DueDate:     time.Now().Add(24 * time.Hour),
	}

	// 测试创建任务
	t.Run("测试创建任务", func(t *testing.T) {
		err := taskService.Create(task)
		assert.NoError(t, err)
		assert.NotZero(t, task.ID)
	})

	// 测试获取任务
	t.Run("测试获取任务", func(t *testing.T) {
		found, err := taskService.Get(task.ID, task.UserID)
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, task.Title, found.Title)
	})

	// 测试更新任务
	t.Run("测试更新任务", func(t *testing.T) {
		task.Title = "Updated Task"
		task.Status = model.TaskStatusDone

		err := taskService.Update(task)
		assert.NoError(t, err)

		found, err := taskService.Get(task.ID, task.UserID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Task", found.Title)
		assert.Equal(t, model.TaskStatusDone, found.Status)
	})

	// 测试获取用户任务列表
	t.Run("测试获取用户任务列表", func(t *testing.T) {
		tasks, total, err := taskService.List(task.UserID, "", 1, 10)
		assert.NoError(t, err)
		assert.NotZero(t, total)
		assert.NotEmpty(t, tasks)
	})

	// 测试删除任务
	t.Run("测试删除任务", func(t *testing.T) {
		err := taskService.Delete(task.ID, task.UserID)
		assert.NoError(t, err)

		found, err := taskService.Get(task.ID, task.UserID)
		assert.NoError(t, err)
		assert.Nil(t, found)
	})
}
