package models

import (
	"gateway-service/internal/dto"
	"gorm.io/gorm"
)

// ServiceInfo 网关服务结构体
type ServiceInfo struct {
	gorm.Model
	LoadType    int    `gorm:"column:load_type" json:"load_type" description:"负载类型 0=http 1=tcp 2=grpc"`
	ServiceName string `gorm:"column:service_name" json:"service_name"`
	ServiceDesc string `gorm:"column:service_desc" json:"service_desc"`
}

// TableName 设置表名
func (*ServiceInfo) TableName() string {
	return "services"
}

// CreateService 创建服务
func CreateService(service *ServiceInfo) error {
	return DB.Create(service).Error
}

// GetServices 获取所有服务
//func GetServices() ([]ServiceInfo, error) {
//	var list []ServiceInfo
//	err := DB.Find(&list).Error
//	return list, err
//}

// GetServiceById 根据ID获取服务
func GetServiceById(id uint) (*ServiceInfo, error) {
	var service ServiceInfo
	err := DB.First(&service, id).Error
	return &service, err
}

// GetServicePage 分页获取服务
func GetServicePage(params *dto.ServiceListInput) ([]ServiceInfo, int64, error) {
	var list []ServiceInfo
	var total int64

	query := DB.Model(&ServiceInfo{})
	query = query.Select("*")

	if params.Info != "" {
		query = query.Where("service_name LIKE ? OR service_desc LIKE ?", "%"+params.Info+"%", "%"+params.Info+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (params.PageNo - 1) * params.PageSize
	err := query.Offset(offset).Limit(params.PageSize).
		Order("id desc").Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}

	return list, total, nil
}

// UpdateService 更新服务
func UpdateService(service *ServiceInfo) error {
	return DB.Save(service).Error
}

// DeleteService 删除服务
func DeleteService(id uint) error {
	return DB.Delete(&ServiceInfo{}, id).Error
}

// ServiceNameExists 服务名是否存在
func ServiceNameExists(name string) (bool, error) {
	var count int64
	err := DB.Model(&ServiceInfo{}).Where("service_name = ?", name).Count(&count).Error
	return count > 0, err
}
