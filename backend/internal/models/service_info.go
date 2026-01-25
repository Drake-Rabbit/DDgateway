package models

import (
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

//func (i *ServiceInfo) ServiceDetail(serviceID uint) (*ServiceDetail, error) {
//	detail := &ServiceDetail{}
//
//	// 获取基本信息
//	service, err := GetServiceById(serviceID)
//	if err != nil {
//		return nil, err
//	}
//
//	detail.Info = service
//
//	// 根据负载类型获取对应规则
//	switch service.LoadType {
//	case 0: // HTTP
//		if rule, err := GetHttpRuleByServiceId(int64(serviceID)); err == nil {
//			detail.HTTPRule = rule
//		}
//	case 1: // TCP
//		if rule, err := GetTcpRuleByServiceId(int64(serviceID)); err == nil {
//			detail.TCPRule = rule
//		}
//	case 2: // GRPC
//		if rule, err := GetGrpcRuleByServiceId(int64(serviceID)); err == nil {
//			detail.GRPCRule = rule
//		}
//	}
//
//	// 获取负载均衡
//	if lb, err := GetLoadBalanceByServiceId(int64(serviceID)); err == nil {
//		detail.LoadBalance = lb
//	}
//
//	// 获取访问控制
//	if ac, err := GetAccessControlByServiceId(int64(serviceID)); err == nil {
//		detail.AccessControl = ac
//	}
//
//	return detail, nil
//}

// CreateService 创建服务
func CreateService(service *ServiceInfo) error {
	return DB.Create(service).Error
}

// GetServices 获取所有服务
func GetServices() ([]ServiceInfo, error) {
	var list []ServiceInfo
	err := DB.Find(&list).Error
	return list, err
}

// GetServiceById 根据ID获取服务
func GetServiceById(id uint) (*ServiceInfo, error) {
	var service ServiceInfo
	err := DB.First(&service, id).Error
	return &service, err
}

// GetServiceByName 根据服务名获取服务
func GetServiceByName(name string) (*ServiceInfo, error) {
	var service ServiceInfo
	err := DB.Where("service_name = ?", name).First(&service).Error
	return &service, err
}

// GetServicePage 分页获取服务
func GetServicePage(offset, limit int, keyword string) ([]ServiceInfo, int64, error) {
	var list []ServiceInfo
	var total int64

	query := DB.Model(&ServiceInfo{})
	query = query.Select("*")

	if keyword != "" {
		query = query.Where("service_name LIKE ? OR service_desc LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("id DESC").Offset(offset).Limit(limit).Find(&list).Error
	return list, total, err
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
