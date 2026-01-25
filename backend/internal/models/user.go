package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 用户结构体
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TenantID  uint           `gorm:"not null;index" json:"tenant_id"`
	Tenant    Tenant         `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Username  string         `gorm:"size:50;not null;uniqueIndex:idx_tenant_user" json:"username"`
	Email     string         `gorm:"size:100;not null;uniqueIndex:idx_tenant_email" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	Role      string         `gorm:"size:20;default:'user'" json:"role"`
	Status    string         `gorm:"size:20;default:'active'" json:"status"`
	LastLogin *time.Time     `json:"last_login,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName 设置表名
func (*User) TableName() string {
	return "users"
}

// HashPassword 密码加密
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword 校验密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// CreateUser 创建用户
func CreateUser(user *User) error {
	return DB.Create(user).Error
}

// GetUsers 获取所有用户
func GetUsers() ([]User, error) {
	var list []User
	err := DB.Preload("Tenant").Find(&list).Error
	return list, err
}

// GetUserById 根据ID获取用户
func GetUserById(id uint) (*User, error) {
	var user User
	err := DB.Preload("Tenant").First(&user, id).Error
	return &user, err
}

// GetUserByUsername 根据用户名获取用户
func GetUserByUsername(username string) (*User, error) {
	var user User
	err := DB.Preload("Tenant").Where("username = ?", username).First(&user).Error
	return &user, err
}

// GetUserByEmail 根据邮箱获取用户
func GetUserByEmail(email string) (*User, error) {
	var user User
	err := DB.Preload("Tenant").Where("email = ?", email).First(&user).Error
	return &user, err
}

// GetUsersByTenantId 根据租户ID获取用户列表
func GetUsersByTenantId(tenantID uint) ([]User, error) {
	var users []User
	err := DB.Preload("Tenant").Where("tenant_id = ?", tenantID).Find(&users).Error
	return users, err
}

// UpdateUser 更新用户
func UpdateUser(user *User) error {
	return DB.Save(user).Error
}

// DeleteUser 删除用户
func DeleteUser(id uint) error {
	return DB.Delete(&User{}, id).Error
}

// UsernameExists 用户名是否存在
func UsernameExists(tenantID uint, username string) (bool, error) {
	var count int64
	err := DB.Model(&User{}).Where("tenant_id = ? AND username = ?", tenantID, username).Count(&count).Error
	return count > 0, err
}

// EmailExists 邮箱是否存在
func EmailExists(tenantID uint, email string) (bool, error) {
	var count int64
	err := DB.Model(&User{}).Where("tenant_id = ? AND email = ?", tenantID, email).Count(&count).Error
	return count > 0, err
}
