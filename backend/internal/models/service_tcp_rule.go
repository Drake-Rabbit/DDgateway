package models

// TcpRule TCP规则结构体
type TcpRule struct {
	ID        int64 `gorm:"primary_key" json:"id"`
	ServiceID int64 `gorm:"column:service_id" json:"service_id"`
	Port      int   `gorm:"column:port" json:"port"`
}

// TableName 设置表名
func (*TcpRule) TableName() string {
	return "gateway_service_tcp_rule"
}

// CreateTcpRule 创建TCP规则
func CreateTcpRule(rule *TcpRule) error {
	return DB.Create(rule).Error
}

// GetTcpRuleById 根据ID获取TCP规则
func GetTcpRuleById(id int64) (*TcpRule, error) {
	var rule TcpRule
	err := DB.First(&rule, id).Error
	return &rule, err
}

// GetTcpRuleByServiceId 根据服务ID获取TCP规则
func GetTcpRuleByServiceId(serviceID int64) (*TcpRule, error) {
	var rule TcpRule
	err := DB.Where("service_id = ?", serviceID).First(&rule).Error
	return &rule, err
}

// UpdateTcpRule 更新TCP规则
func UpdateTcpRule(rule *TcpRule) error {
	return DB.Save(rule).Error
}

// DeleteTcpRule 删除TCP规则
func DeleteTcpRule(id int64) error {
	return DB.Delete(&TcpRule{}, id).Error
}

// DeleteTcpRuleByServiceId 根据服务ID删除TCP规则
func DeleteTcpRuleByServiceId(serviceID int64) error {
	return DB.Where("service_id = ?", serviceID).Delete(&TcpRule{}).Error
}
