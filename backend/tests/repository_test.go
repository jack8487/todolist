package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"todolist/internal/model"
	"todolist/internal/repository"
)

func TestUserRepository(t *testing.T) {
	// 初始化测试数据库
	db := initTestDB(t)
	defer cleanupTestDB(t, db)

	userRepo := repository.NewUserRepository(db)

	// 测试创建用户
	t.Run("测试创建用户", func(t *testing.T) {
		user := &model.User{
			Username:     "test_user",
			PasswordHash: "hashed_password",
		}

		err := userRepo.Create(user)
		if err != nil {
			t.Errorf("创建用户失败: %v", err)
		}
		if user.ID == 0 {
			t.Error("用户ID未生成")
		}

		// 测试获取用户
		found, err := userRepo.GetByID(user.ID)
		if err != nil {
			t.Errorf("获取用户失败: %v", err)
		}
		if found == nil {
			t.Error("未找到创建的用户")
		}
		if found.Username != user.Username {
			t.Errorf("用户名不匹配: 期望 %s, 实际 %s", user.Username, found.Username)
		}
	})

	// 测试通过用户名获取用户
	t.Run("测试通过用户名获取用户", func(t *testing.T) {
		username := "test_user_2"
		user := &model.User{
			Username:     username,
			PasswordHash: "hashed_password",
		}

		err := userRepo.Create(user)
		if err != nil {
			t.Errorf("创建用户失败: %v", err)
		}
		defer userRepo.Delete(user.ID)

		found, err := userRepo.GetByUsername(username)
		if err != nil {
			t.Errorf("通过用户名获取用户失败: %v", err)
		}
		if found == nil {
			t.Error("未找到用户")
		}
		if found.Username != username {
			t.Errorf("用户名不匹配: 期望 %s, 实际 %s", username, found.Username)
		}
	})
}

func TestTaskRepository(t *testing.T) {
	// 初始化测试数据库
	db := initTestDB(t)
	defer cleanupTestDB(t, db)

	taskRepo := repository.NewTaskRepository(db)

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
		err := taskRepo.Create(task)
		assert.NoError(t, err)
		assert.NotZero(t, task.ID)
	})

	// 测试获取任务
	t.Run("测试获取任务", func(t *testing.T) {
		found, err := taskRepo.GetByID(task.ID)
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, task.Title, found.Title)
		assert.Equal(t, task.Description, found.Description)
	})

	// 测试更新任务
	t.Run("测试更新任务", func(t *testing.T) {
		task.Title = "Updated Task"
		task.Description = "Updated Description"
		task.Status = model.TaskStatusDone

		err := taskRepo.Update(task)
		assert.NoError(t, err)

		found, err := taskRepo.GetByID(task.ID)
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, "Updated Task", found.Title)
		assert.Equal(t, "Updated Description", found.Description)
		assert.Equal(t, model.TaskStatusDone, found.Status)
	})

	// 测试获取用户任务列表
	t.Run("测试获取用户任务列表", func(t *testing.T) {
		tasks, total, err := taskRepo.GetByUserID(1, "", 1, 10)
		assert.NoError(t, err)
		assert.NotZero(t, total)
		assert.NotEmpty(t, tasks)
	})

	// 测试删除任务
	t.Run("测试删除任务", func(t *testing.T) {
		err := taskRepo.Delete(task.ID)
		assert.NoError(t, err)

		found, err := taskRepo.GetByID(task.ID)
		assert.NoError(t, err)
		assert.Nil(t, found)
	})
}
