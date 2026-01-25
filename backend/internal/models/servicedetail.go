package models

// ServiceDetail 服务详情结构体
type ServiceDetail struct {
	Info          *ServiceInfo   `json:"info"`
	HTTPRule      *HttpRule      `json:"http_rule"`
	TCPRule       *TcpRule       `json:"tcp_rule"`
	GRPCRule      *GrpcRule      `json:"grpc_rule"`
	LoadBalance   *LoadBalance   `json:"load_balance"`
	AccessControl *AccessControl `json:"access_control"`
}

// GetServiceDetailById 获取服务详情
func GetServiceDetailById(serviceID uint) (*ServiceDetail, error) {
	detail := &ServiceDetail{}

	// 获取基本信息
	service, err := GetServiceById(serviceID)
	if err != nil {
		return nil, err
	}
	//if service.ServiceName == "" {
	//	return nil, errors.New("服务不存在")
	//}

	detail.Info = service

	// 根据负载类型获取对应规则
	switch service.LoadType {
	case 0: // HTTP
		if rule, err := GetHttpRuleByServiceId(int64(serviceID)); err == nil {
			detail.HTTPRule = rule
		}
	case 1: // TCP
		if rule, err := GetTcpRuleByServiceId(int64(serviceID)); err == nil {
			detail.TCPRule = rule
		}
	case 2: // GRPC
		if rule, err := GetGrpcRuleByServiceId(int64(serviceID)); err == nil {
			detail.GRPCRule = rule
		}
	}

	// 获取负载均衡
	if lb, err := GetLoadBalanceByServiceId(int64(serviceID)); err == nil {
		detail.LoadBalance = lb
	}

	// 获取访问控制
	if ac, err := GetAccessControlByServiceId(int64(serviceID)); err == nil {
		detail.AccessControl = ac
	}

	return detail, nil
}

////服务详情是否存在count
//func IsServiceDetail(serviceID uint)  error{
//	return DB.Model(&ServiceDetail{}).Where("service_id = ?", serviceID).Count(&ServiceDetail{}).Error
//}
