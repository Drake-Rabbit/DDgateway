package models

// HttpRule HTTP规则结构体
type HttpRule struct {
	ID             int64  `gorm:"primary_key" json:"id"`
	ServiceID      int64  `gorm:"column:service_id" json:"service_id"`
	RuleType       int    `gorm:"column:rule_type" json:"rule_type"`
	Rule           string `gorm:"column:rule" json:"rule"`
	NeedHttps      int    `gorm:"column:need_https" json:"need_https"`
	NeedWebsocket  int    `gorm:"column:need_websocket" json:"need_websocket"`
	NeedStripUri   int    `gorm:"column:need_strip_uri" json:"need_strip_uri"`
	UrlRewrite     string `gorm:"column:url_rewrite" json:"url_rewrite"`
	HeaderTransfor string `gorm:"column:header_transfor" json:"header_transfor"`
}

// TableName 设置表名
func (*HttpRule) TableName() string {
	return "gateway_service_http_rule"
}

// CreateHttpRule 创建HTTP规则
func CreateHttpRule(rule *HttpRule) error {
	return DB.Create(rule).Error
}

// GetHttpRuleById 根据ID获取HTTP规则
func GetHttpRuleById(id int64) (*HttpRule, error) {
	var rule HttpRule
	err := DB.First(&rule, id).Error
	return &rule, err
}

// GetHttpRuleByServiceId 根据服务ID获取HTTP规则
func GetHttpRuleByServiceId(serviceID int64) (*HttpRule, error) {
	var rule HttpRule
	err := DB.Where("service_id = ?", serviceID).First(&rule).Error
	return &rule, err
}

// UpdateHttpRule 更新HTTP规则
func UpdateHttpRule(rule *HttpRule) error {
	return DB.Save(rule).Error
}

// DeleteHttpRule 删除HTTP规则
func DeleteHttpRule(id int64) error {
	return DB.Delete(&HttpRule{}, id).Error
}

// DeleteHttpRuleByServiceId 根据服务ID删除HTTP规则
func DeleteHttpRuleByServiceId(serviceID int64) error {
	return DB.Where("service_id = ?", serviceID).Delete(&HttpRule{}).Error
}

// HttpRuleExists 检查规则是否存在（根据rule_type和rule）
func HttpRuleExists(ruleType int, rule string) (bool, error) {
	var count int64
	err := DB.Model(&HttpRule{}).Where("rule_type = ? AND rule = ?", ruleType, rule).Count(&count).Error
	return count > 0, err
}
