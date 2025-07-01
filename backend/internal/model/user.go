package model

import "time"

/*
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
*/

type User struct {
	ID           int       `json:"id" gorm:"primaryKey;autoIncrement" validate:"-"`                                   // 自增主键，无需验证
	Username     string    `json:"username" gorm:"type:varchar(50);unique;not null" validate:"required,min=3,max=50"` // 用户名必填，3-50字符
	PasswordHash string    `json:"password_hash" gorm:"type:varchar(255);not null" validate:"required,min=6"`         // 密码哈希，至少6字符
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime" validate:"-"`                                     // 自动设置创建时间
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime" validate:"-"`                                     // 自动更新时间
}
