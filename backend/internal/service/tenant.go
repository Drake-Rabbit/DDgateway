package service

import (
	"fmt"
	"gateway-service/internal/models"
)

// TenantService 租户服务
type TenantService struct{}

// ListTenants 获取所有租户
func (s *TenantService) ListTenants() ([]models.Tenant, error) {
	return models.GetTenants()
}

// GetTenant 根据ID获取租户
func (s *TenantService) GetTenant(id string) (*models.Tenant, error) {
	var uintId uint
	if _, err := fmt.Sscanf(id, "%d", &uintId); err != nil {
		return nil, err
	}
	return models.GetTenantById(uintId)
}

// CreateTenant 创建租户
func (s *TenantService) CreateTenant(name, code, description string) (*models.Tenant, error) {
	tenant := &models.Tenant{
		Name:        name,
		Code:        code,
		Description: description,
		Status:      "active",
	}
	if err := models.CreateTenant(tenant); err != nil {
		return nil, err
	}
	return tenant, nil
}

// UpdateTenant 更新租户
func (s *TenantService) UpdateTenant(id string, name, description, status string) (*models.Tenant, error) {
	var uintId uint
	if _, err := fmt.Sscanf(id, "%d", &uintId); err != nil {
		return nil, err
	}

	tenant, err := models.GetTenantById(uintId)
	if err != nil {
		return nil, err
	}

	if name != "" {
		tenant.Name = name
	}
	if description != "" {
		tenant.Description = description
	}
	if status != "" {
		tenant.Status = status
	}

	if err := models.UpdateTenant(tenant); err != nil {
		return nil, err
	}
	return tenant, nil
}

// DeleteTenant 删除租户
func (s *TenantService) DeleteTenant(id string) error {
	var uintId uint
	if _, err := fmt.Sscanf(id, "%d", &uintId); err != nil {
		return err
	}
	return models.DeleteTenant(uintId)
}
