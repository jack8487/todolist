package service

import (
	"errors"
	"time"

	"todolist/internal/model"
	"todolist/internal/repository"
)

var (
	ErrTaskNotFound       = errors.New("任务不存在")
	ErrTaskAccessDenied   = errors.New("无权访问该任务")
	ErrInvalidDueDate     = errors.New("无效的截止日期")
	ErrEmptyTitle         = errors.New("任务标题不能为空")
	ErrTitleTooLong       = errors.New("任务标题不能超过100个字符")
	ErrDescriptionTooLong = errors.New("任务描述不能超过500个字符")
)

// TaskService 任务服务接口
type TaskService interface {
	// Create 创建任务
	Create(task *model.Task) error
	// Update 更新任务
	Update(task *model.Task) error
	// Delete 删除任务
	Delete(taskID, userID int) error
	// Get 获取任务详情
	Get(taskID, userID int) (*model.Task, error)
	// List 获取任务列表
	List(userID int, status string, page, pageSize int) ([]*model.Task, int64, error)
}

// taskService 任务服务实现
type taskService struct {
	taskRepo repository.TaskRepository
}

// NewTaskService 创建任务服务实例
func NewTaskService(taskRepo repository.TaskRepository) TaskService {
	return &taskService{
		taskRepo: taskRepo,
	}
}

// Create 创建任务
func (s *taskService) Create(task *model.Task) error {
	// 验证任务标题
	if err := s.validateTaskTitle(task.Title); err != nil {
		return err
	}

	// 验证任务描述
	if err := s.validateTaskDescription(task.Description); err != nil {
		return err
	}

	// 设置创建和更新时间
	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now

	return s.taskRepo.Create(task)
}

// Get 获取任务详情
func (s *taskService) Get(taskID, userID int) (*model.Task, error) {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}

	// 验证任务所有权
	if task.UserID != userID {
		return nil, ErrTaskAccessDenied
	}

	return task, nil
}

// List 获取任务列表
func (s *taskService) List(userID int, status string, page, pageSize int) ([]*model.Task, int64, error) {
	return s.taskRepo.GetByUserID(userID, status, page, pageSize)
}

// Update 更新任务
func (s *taskService) Update(task *model.Task) error {
	// 获取任务
	oldTask, err := s.Get(task.ID, task.UserID)
	if err != nil {
		return err
	}

	// 验证任务标题
	if task.Title != "" {
		if err := s.validateTaskTitle(task.Title); err != nil {
			return err
		}
		oldTask.Title = task.Title
	}

	// 验证任务描述
	if task.Description != "" {
		if err := s.validateTaskDescription(task.Description); err != nil {
			return err
		}
		oldTask.Description = task.Description
	}

	// 验证截止日期
	if !task.DueDate.IsZero() {
		if err := s.validateDueDate(task.DueDate); err != nil {
			return err
		}
		oldTask.DueDate = task.DueDate
	}

	// 更新状态
	if task.Status >= model.TaskStatusTodo && task.Status <= model.TaskStatusDone {
		oldTask.Status = task.Status
	}

	oldTask.UpdatedAt = time.Now()
	return s.taskRepo.Update(oldTask)
}

// Delete 删除任务
func (s *taskService) Delete(taskID, userID int) error {
	// 验证任务所有权
	_, err := s.Get(taskID, userID)
	if err != nil {
		return err
	}

	return s.taskRepo.Delete(taskID)
}

// validateTaskTitle 验证任务标题
func (s *taskService) validateTaskTitle(title string) error {
	if title == "" {
		return ErrEmptyTitle
	}
	if len(title) > 100 {
		return ErrTitleTooLong
	}
	return nil
}

// validateTaskDescription 验证任务描述
func (s *taskService) validateTaskDescription(description string) error {
	if len(description) > 500 {
		return ErrDescriptionTooLong
	}
	return nil
}

// validateDueDate 验证截止日期
func (s *taskService) validateDueDate(dueDate time.Time) error {
	if dueDate.IsZero() {
		return nil // 允许空的截止日期
	}
	// 可以添加其他日期验证逻辑，比如不允许过去的日期
	return nil
}
