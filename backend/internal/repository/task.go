package repository

import (
	"gorm.io/gorm"

	"todolist/internal/model"
)

// TaskRepository 任务仓库接口
type TaskRepository interface {
	// Create 创建任务
	Create(task *model.Task) error
	// Update 更新任务
	Update(task *model.Task) error
	// Delete 删除任务
	Delete(taskID int) error
	// GetByID 根据ID获取任务
	GetByID(taskID int) (*model.Task, error)
	// GetByUserID 获取用户的任务列表
	GetByUserID(userID int, status string, page, pageSize int) ([]*model.Task, int64, error)
}

// taskRepository 任务仓库实现
type taskRepository struct {
	db *gorm.DB
}

// NewTaskRepository 创建任务仓库实例
func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{
		db: db,
	}
}

// Create 创建任务
func (r *taskRepository) Create(task *model.Task) error {
	return r.db.Create(task).Error
}

// Update 更新任务
func (r *taskRepository) Update(task *model.Task) error {
	return r.db.Save(task).Error
}

// Delete 删除任务
func (r *taskRepository) Delete(taskID int) error {
	return r.db.Delete(&model.Task{}, taskID).Error
}

// GetByID 根据ID获取任务
func (r *taskRepository) GetByID(taskID int) (*model.Task, error) {
	var task model.Task
	err := r.db.First(&task, taskID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

// GetByUserID 获取用户的任务列表
func (r *taskRepository) GetByUserID(userID int, status string, page, pageSize int) ([]*model.Task, int64, error) {
	var tasks []*model.Task
	var total int64

	query := r.db.Model(&model.Task{}).Where("user_id = ?", userID)
	if status != "" {
		// 将状态文本转换为数字
		var statusInt int
		switch status {
		case "todo":
			statusInt = model.TaskStatusTodo
		case "in_progress":
			statusInt = model.TaskStatusInProgress
		case "done":
			statusInt = model.TaskStatusDone
		}
		query = query.Where("status = ?", statusInt)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&tasks).Error
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// UpdateStatus 更新任务状态
func (r *taskRepository) UpdateStatus(id int, status bool) error {
	return r.db.Model(&model.Task{}).Where("id = ?", id).Update("status", status).Error
}
