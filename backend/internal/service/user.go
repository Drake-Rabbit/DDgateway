package service

import (
	"errors"
	"fmt"
	"gateway-service/internal/dto"
	"gateway-service/internal/models"
	"strconv"
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

// UpdatePasswordService 更新密码
func (s *UserService) UpdatePassword(input *dto.ChangePasswordInput) error {
	//int化string类型的userID
	var userID, err = strconv.Atoi(input.UserID)
	if err != nil {
		return err
	}

	//1.先检查用户旧密码是否正确
	trueUser, err := models.GetUserById(uint(userID))

	isright := trueUser.CheckPassword(input.OldPassword)
	if !isright {
		return errors.New("旧密码输入错误")
	}
	//2.去model更新新密码
	trueUser.Password = input.NewPassword
	if err := trueUser.HashPassword(); err != nil {
		return errors.New("密码加密失败")
	}
	//3.数据库更新
	if err := models.UpdateUser(trueUser); err != nil {
		return errors.New("密码更新失败")
	}
	return nil
}
