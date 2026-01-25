package models

import "strings"

// LoadBalance 负载均衡结构体
type LoadBalance struct {
	ID                     int64  `gorm:"primary_key" json:"id"`
	ServiceID              int64  `gorm:"column:service_id" json:"service_id" description:"服务id	"`
	CheckMethod            int    `gorm:"column:check_method" json:"check_method" description:"检查方法 tcpchk=检测端口是否握手成功	"`
	CheckTimeout           int    `gorm:"column:check_timeout" json:"check_timeout" description:"check超时时间"`
	CheckInterval          int    `gorm:"column:check_interval" json:"check_interval" description:"检查间隔, 单位s"`
	RoundType              int    `gorm:"column:round_type" json:"round_type" description:"轮询方式 round/weight_round/random/ip_hash"`
	IpList                 string `gorm:"column:ip_list" json:"ip_list"`
	WeightList             string `gorm:"column:weight_list" json:"weight_list" description:"权重列表"`
	ForbidList             string `gorm:"column:forbid_list" json:"forbid_list" description:"禁用ip列表"`
	UpstreamConnectTimeout int    `gorm:"column:upstream_connect_timeout" json:"upstream_connect_timeout"`
	UpstreamHeaderTimeout  int    `gorm:"column:upstream_header_timeout" json:"upstream_header_timeout" description:"下游获取header超时, 单位s"`
	UpstreamIdleTimeout    int    `gorm:"column:upstream_idle_timeout" json:"upstream_idle_timeout" description:"下游链接最大空闲时间, 单位s"`
	UpstreamMaxIdle        int    `gorm:"column:upstream_max_idle" json:"upstream_max_idle" description:"下游最大空闲链接数"`
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
