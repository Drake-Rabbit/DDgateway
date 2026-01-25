package models

// GrpcRule gRPC规则结构体
type GrpcRule struct {
	ID             int64  `gorm:"primary_key" json:"id"`
	ServiceID      int64  `gorm:"column:service_id" json:"service_id"`
	Port           int    `gorm:"column:port" json:"port"`
	HeaderTransfor string `gorm:"column:header_transfor" json:"header_transfor"`
}

// TableName 设置表名
func (*GrpcRule) TableName() string {
	return "gateway_service_grpc_rule"
}

// CreateGrpcRule 创建gRPC规则
func CreateGrpcRule(rule *GrpcRule) error {
	return DB.Create(rule).Error
}

// GetGrpcRuleById 根据ID获取gRPC规则
func GetGrpcRuleById(id int64) (*GrpcRule, error) {
	var rule GrpcRule
	err := DB.First(&rule, id).Error
	return &rule, err
}

// GetGrpcRuleByServiceId 根据服务ID获取gRPC规则
func GetGrpcRuleByServiceId(serviceID int64) (*GrpcRule, error) {
	var rule GrpcRule
	err := DB.Where("service_id = ?", serviceID).First(&rule).Error
	return &rule, err
}

// UpdateGrpcRule 更新gRPC规则
func UpdateGrpcRule(rule *GrpcRule) error {
	return DB.Save(rule).Error
}

// DeleteGrpcRule 删除gRPC规则
func DeleteGrpcRule(id int64) error {
	return DB.Delete(&GrpcRule{}, id).Error
}

// DeleteGrpcRuleByServiceId 根据服务ID删除gRPC规则
func DeleteGrpcRuleByServiceId(serviceID int64) error {
	return DB.Where("service_id = ?", serviceID).Delete(&GrpcRule{}).Error
}
