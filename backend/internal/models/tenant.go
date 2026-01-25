package models

import (
	"time"

	"gorm.io/gorm"
)

// Tenant 租户结构体
type Tenant struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Code        string         `gorm:"size:50;uniqueIndex;not null" json:"code"`
	Description string         `gorm:"size:255" json:"description"`
	Status      string         `gorm:"size:20;default:'active'" json:"status"`
	Users       []User         `gorm:"foreignKey:TenantID" json:"users,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName 设置表名
func (*Tenant) TableName() string {
	return "tenants"
}

// CreateTenant 创建租户
func CreateTenant(tenant *Tenant) error {
	return DB.Create(tenant).Error
}

// GetTenants 获取所有租户
func GetTenants() ([]Tenant, error) {
	var list []Tenant
	err := DB.Find(&list).Error
	return list, err
}

// GetTenantById 根据ID获取租户
func GetTenantById(id uint) (*Tenant, error) {
	var tenant Tenant
	err := DB.Preload("Users").First(&tenant, id).Error
	return &tenant, err
}

// GetTenantByCode 根据代码获取租户
func GetTenantByCode(code string) (*Tenant, error) {
	var tenant Tenant
	err := DB.Where("code = ?", code).First(&tenant).Error
	return &tenant, err
}

// UpdateTenant 更新租户
func UpdateTenant(tenant *Tenant) error {
	return DB.Save(tenant).Error
}

// DeleteTenant 删除租户
func DeleteTenant(id uint) error {
	return DB.Delete(&Tenant{}, id).Error
}

// TenantExists 检查租户是否存在
func TenantExists(code string) (bool, error) {
	var count int64
	err := DB.Model(&Tenant{}).Where("code = ?", code).Count(&count).Error
	return count > 0, err
}
