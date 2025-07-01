package service

import (
	"errors"
	"time"

	"todolist/internal/model"
	"todolist/internal/repository"
	"todolist/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound     = errors.New("用户不存在")
	ErrUserExists       = errors.New("用户已存在")
	ErrInvalidPassword  = errors.New("密码错误")
	ErrPasswordTooShort = errors.New("密码长度不能小于6位")
)

// UserService 用户服务接口
type UserService interface {
	Register(username, password string) error
	Login(username, password string) (string, error)
	GetUserByID(id int) (*model.User, error)
	UpdatePassword(id int, oldPassword, newPassword string) error
}

// userService 用户服务实现
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// Register 用户注册
func (s *userService) Register(username, password string) error {
	// 检查用户名是否已存在
	existingUser, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return ErrUserExists
	}

	// 验证密码长度
	if len(password) < 6 {
		return ErrPasswordTooShort
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 创建用户
	user := &model.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return s.userRepo.Create(user)
}

// Login 用户登录
func (s *userService) Login(username, password string) (string, error) {
	// 获取用户
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrUserNotFound
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", ErrInvalidPassword
	}

	// 生成 JWT token
	token, err := jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetUserByID 根据ID获取用户信息
func (s *userService) GetUserByID(id int) (*model.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// UpdatePassword 更新用户密码
func (s *userService) UpdatePassword(id int, oldPassword, newPassword string) error {
	// 获取用户
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	// 验证旧密码
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword))
	if err != nil {
		return ErrInvalidPassword
	}

	// 验证新密码长度
	if len(newPassword) < 6 {
		return ErrPasswordTooShort
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 更新用户信息
	user.PasswordHash = string(hashedPassword)
	user.UpdatedAt = time.Now()

	return s.userRepo.Update(user)
}
