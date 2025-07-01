package model

import "time"

/*
CREATE TABLE tasks (

	id BIGINT PRIMARY KEY AUTO_INCREMENT,
	user_id BIGINT NOT NULL,
	title VARCHAR(100) NOT NULL,
	description TEXT,
	status TINYINT DEFAULT 0, -- 0: 未完成, 1: 已完成
	due_date TIMESTAMP,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

);
*/

// Task 任务模型
type Task struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	UserID      int       `json:"user_id" gorm:"not null"`
	Title       string    `json:"title" gorm:"size:100;not null"`
	Description string    `json:"description" gorm:"size:500"`
	Status      int       `json:"status" gorm:"type:tinyint;not null;default:0"`
	DueDate     time.Time `json:"due_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// 任务状态常量
const (
	TaskStatusTodo       = 0 // 待办
	TaskStatusInProgress = 1 // 进行中
	TaskStatusDone       = 2 // 已完成
)

// GetStatusText 获取状态文本
func (t *Task) GetStatusText() string {
	switch t.Status {
	case TaskStatusTodo:
		return "todo"
	case TaskStatusInProgress:
		return "in_progress"
	case TaskStatusDone:
		return "done"
	default:
		return "unknown"
	}
}

// SetStatusFromText 从文本设置状态
func (t *Task) SetStatusFromText(status string) {
	switch status {
	case "todo":
		t.Status = TaskStatusTodo
	case "in_progress":
		t.Status = TaskStatusInProgress
	case "done":
		t.Status = TaskStatusDone
	default:
		t.Status = TaskStatusTodo
	}
}
