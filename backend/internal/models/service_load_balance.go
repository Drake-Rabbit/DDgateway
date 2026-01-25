package models

import "strings"

// LoadBalance 负载均衡结构体
type LoadBalance struct {
	ID                     int64  `gorm:"primary_key" json:"id"`
	ServiceID              int64  `gorm:"column:service_id" json:"service_id"`
	CheckMethod            int    `gorm:"column:check_method" json:"check_method"`
	CheckTimeout           int    `gorm:"column:check_timeout" json:"check_timeout"`
	CheckInterval          int    `gorm:"column:check_interval" json:"check_interval"`
	RoundType              int    `gorm:"column:round_type" json:"round_type"`
	IpList                 string `gorm:"column:ip_list" json:"ip_list"`
	WeightList             string `gorm:"column:weight_list" json:"weight_list"`
	ForbidList             string `gorm:"column:forbid_list" json:"forbid_list"`
	UpstreamConnectTimeout int    `gorm:"column:upstream_connect_timeout" json:"upstream_connect_timeout"`
	UpstreamHeaderTimeout  int    `gorm:"column:upstream_header_timeout" json:"upstream_header_timeout"`
	UpstreamIdleTimeout    int    `gorm:"column:upstream_idle_timeout" json:"upstream_idle_timeout"`
	UpstreamMaxIdle        int    `gorm:"column:upstream_max_idle" json:"upstream_max_idle"`
}

// TableName 设置表名
func (*LoadBalance) TableName() string {
	return "gateway_service_load_balance"
}

func (b *LoadBalance) GetIPListByModel() []string {
	return strings.Split(b.IpList, ",")
}

// CreateLoadBalance 创建负载均衡
func CreateLoadBalance(lb *LoadBalance) error {
	return DB.Create(lb).Error
}

// GetLoadBalanceById 根据ID获取负载均衡
func GetLoadBalanceById(id int64) (*LoadBalance, error) {
	var lb LoadBalance
	err := DB.First(&lb, id).Error
	return &lb, err
}

// GetLoadBalanceByServiceId 根据服务ID获取负载均衡
func GetLoadBalanceByServiceId(serviceID int64) (*LoadBalance, error) {
	var lb LoadBalance
	err := DB.Where("service_id = ?", serviceID).First(&lb).Error
	return &lb, err
}

// UpdateLoadBalance 更新负载均衡
func UpdateLoadBalance(lb *LoadBalance) error {
	return DB.Save(lb).Error
}

// DeleteLoadBalance 删除负载均衡
func DeleteLoadBalance(id int64) error {
	return DB.Delete(&LoadBalance{}, id).Error
}

// DeleteLoadBalanceByServiceId 根据服务ID删除负载均衡
func DeleteLoadBalanceByServiceId(serviceID int64) error {
	return DB.Where("service_id = ?", serviceID).Delete(&LoadBalance{}).Error
}
