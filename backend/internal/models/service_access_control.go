package models

// AccessControl 访问控制结构体
type AccessControl struct {
	ID                int64  `gorm:"primary_key" json:"id"`
	ServiceID         int64  `gorm:"column:service_id" json:"service_id"`
	OpenAuth          int    `gorm:"column:open_auth" json:"open_auth"`
	BlackList         string `gorm:"column:black_list" json:"black_list"`
	WhiteList         string `gorm:"column:white_list" json:"white_list"`
	WhiteHostName     string `gorm:"column:white_host_name" json:"white_host_name"`
	ClientIPFlowLimit int    `gorm:"column:clientip_flow_limit" json:"clientip_flow_limit"`
	ServiceFlowLimit  int    `gorm:"column:service_flow_limit" json:"service_flow_limit"`
}

// TableName 设置表名
func (*AccessControl) TableName() string {
	return "gateway_service_access_control"
}

// CreateAccessControl 创建访问控制
func CreateAccessControl(ac *AccessControl) error {
	return DB.Create(ac).Error
}

// GetAccessControlById 根据ID获取访问控制
func GetAccessControlById(id int64) (*AccessControl, error) {
	var ac AccessControl
	err := DB.First(&ac, id).Error
	return &ac, err
}

// GetAccessControlByServiceId 根据服务ID获取访问控制
func GetAccessControlByServiceId(serviceID int64) (*AccessControl, error) {
	var ac AccessControl
	err := DB.Where("service_id = ?", serviceID).First(&ac).Error
	return &ac, err
}

// UpdateAccessControl 更新访问控制
func UpdateAccessControl(ac *AccessControl) error {
	return DB.Save(ac).Error
}

// DeleteAccessControl 删除访问控制
func DeleteAccessControl(id int64) error {
	return DB.Delete(&AccessControl{}, id).Error
}

// DeleteAccessControlByServiceId 根据服务ID删除访问控制
func DeleteAccessControlByServiceId(serviceID int64) error {
	return DB.Where("service_id = ?", serviceID).Delete(&AccessControl{}).Error
}
