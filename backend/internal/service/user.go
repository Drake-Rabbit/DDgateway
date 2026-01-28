package service

import (
	"errors"
	"fmt"
	"gateway-service/internal/models"
	"time"
)

// UserService 用户服务
type UserService struct{}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username   string
	Email      string
	Password   string
	TenantCode string
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string
	Password string
}

// ListUsers 获取用户列表
func (s *UserService) ListUsers(tenantID int) ([]models.User, error) {
	return models.GetUsersByTenantId(uint(tenantID))
}

// GetUser 获取用户
func (s *UserService) GetUser(userID string) (*models.User, error) {
	var uintId uint
	if _, err := fmt.Sscanf(userID, "%d", &uintId); err != nil {
		return nil, err
	}

	user, err := models.GetUserById(uintId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(userID string, email, role, status string) (*models.User, error) {
	var uintId uint
	if _, err := fmt.Sscanf(userID, "%d", &uintId); err != nil {
		return nil, err
	}

	user, err := models.GetUserById(uintId)
	if err != nil {
		return nil, err
	}

	if email != "" {
		user.Email = email
	}
	if role != "" {
		user.Role = role
	}
	if status != "" {
		user.Status = status
	}

	if err := models.UpdateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(userID string) error {
	var uintId uint
	if _, err := fmt.Sscanf(userID, "%d", &uintId); err != nil {
		return err
	}

	user, err := models.GetUserById(uintId)
	if err != nil {
		return err
	}

	return models.DeleteUser(user.ID)
}

// Register 用户注册
func (s *UserService) Register(req RegisterRequest) (*models.User, error) {

	// 检查用户名是否存在
	exists, err := models.UsernameExists(req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	// 检查邮箱是否存在
	exists, err = models.EmailExists(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     "user",
		Status:   "active",
	}

	if err := user.HashPassword(); err != nil {
		return nil, errors.New("failed to hash password")
	}

	if err := models.CreateUser(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	return user, nil
}

// Login 用户登录
func (s *UserService) Login(req LoginRequest) (*models.User, error) {
	user, err := models.GetUserByUsername(req.Username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if user.Status != "active" {
		return nil, errors.New("user account is inactive")
	}

	if !user.CheckPassword(req.Password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// UpdateLastLogin 更新最后登录时间
func (s *UserService) UpdateLastLogin(user *models.User) error {
	now := time.Now()
	user.LastLogin = &now
	return models.UpdateUser(user)
}
