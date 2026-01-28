package service

import (
	"errors"
	"gateway-service/internal/define"
	"gateway-service/internal/dto"
	"gateway-service/internal/models"
	"strings"
	"time"
)

// ServiceService 网关服务
type ServiceService struct{}

// PageListServices 分页获取服务列表
func (s *ServiceService) PageListServices(in *dto.ServiceListInput) ([]models.ServiceInfo, int, int, int64, error) {
	var offset int
	// 分页 默认第一页,默认size大小为10
	pageNo, pageSize, offset := define.PageHelper(in.PageNo, in.PageSize)

	list, total, err := models.GetServicePage(offset, pageSize, in.Info)

	//数据总数
	return list, pageNo, pageSize, total, err
}

// GetService 获取服务详情
func (s *ServiceService) GetServiceDetail(serviceID int64) (*models.ServiceDetail, error) {
	return models.GetServiceDetailById(uint(serviceID))
}

// CreateService 创建服务
func (s *ServiceService) CreateService(loadType int, serviceName, serviceDesc string) (*models.ServiceInfo, error) {
	service := &models.ServiceInfo{
		LoadType:    loadType,
		ServiceName: serviceName,
		ServiceDesc: serviceDesc,
	}

	if err := models.CreateService(service); err != nil {
		return nil, err
	}

	return service, nil
}

// UpdateService 更新服务
func (s *ServiceService) UpdateService(serviceID int64, loadType int, serviceName, serviceDesc string) (*models.ServiceInfo, error) {
	service, err := models.GetServiceById(uint(serviceID))
	if err != nil {
		return nil, err
	}

	service.LoadType = loadType
	service.ServiceName = serviceName
	service.ServiceDesc = serviceDesc

	if err := models.UpdateService(service); err != nil {
		return nil, err
	}

	return service, nil
}

// DeleteService 删除服务
func (s *ServiceService) DeleteService(serviceID int64) error {

	// 删除关联规则 (理论上serviceinfo被逻辑删除,gorm就找不到对应的联动数据),以便快速恢复
	//models.DeleteHttpRuleByServiceId(serviceID)
	//models.DeleteTcpRuleByServiceId(serviceID)
	//models.DeleteGrpcRuleByServiceId(serviceID)
	//models.DeleteAccessControlByServiceId(serviceID)
	//models.DeleteLoadBalanceByServiceId(serviceID)
	// 逻辑删除服务serviinfo即可,以便快速恢复
	return models.DeleteService(uint(serviceID))
}

// CreateHTTPService 创建 HTTP 服务
func (s *ServiceService) CreateHTTPService(input *dto.ServiceAddHTTPInput) (*models.ServiceInfo, error) {
	// 验证 IP 列表与权重列表数量一致
	if len(strings.Split(input.IpList, ",")) != len(strings.Split(input.WeightList, ",")) {
		return nil, errors.New("IP列表与权重列表数量不一致")
	}

	// 开始事务
	tx := models.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// 检查服务名是否已存在
	exists, err := models.ServiceNameExists(input.ServiceName)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if exists {
		tx.Rollback()
		return nil, errors.New("服务已存在")
	}

	// 检查规则是否已存在
	ruleExists, err := models.HttpRuleExists(input.RuleType, input.Rule)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if ruleExists {
		tx.Rollback()
		return nil, errors.New("服务接入前缀或域名已存在")
	}

	// 创建服务
	serviceInfo := &models.ServiceInfo{
		LoadType:    0, // HTTP
		ServiceName: input.ServiceName,
		ServiceDesc: input.ServiceDesc,
	}
	if err := tx.Create(serviceInfo).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 创建 HTTP 规则
	httpRule := &models.HttpRule{
		ServiceID:      int64(serviceInfo.ID),
		RuleType:       input.RuleType,
		Rule:           input.Rule,
		NeedHttps:      input.NeedHttps,
		NeedStripUri:   input.NeedStripUri,
		NeedWebsocket:  input.NeedWebsocket,
		UrlRewrite:     input.UrlRewrite,
		HeaderTransfor: input.HeaderTransfor,
	}
	if err := tx.Create(httpRule).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 创建访问控制
	accessControl := &models.AccessControl{
		ServiceID:         int64(serviceInfo.ID),
		OpenAuth:          input.OpenAuth,
		BlackList:         input.BlackList,
		WhiteList:         input.WhiteList,
		ClientIPFlowLimit: input.ClientipFlowLimit,
		ServiceFlowLimit:  input.ServiceFlowLimit,
	}
	if err := tx.Create(accessControl).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 创建负载均衡
	loadBalance := &models.LoadBalance{
		ServiceID:              int64(serviceInfo.ID),
		RoundType:              input.RoundType,
		IpList:                 input.IpList,
		WeightList:             input.WeightList,
		UpstreamConnectTimeout: input.UpstreamConnectTimeout,
		UpstreamHeaderTimeout:  input.UpstreamHeaderTimeout,
		UpstreamIdleTimeout:    input.UpstreamIdleTimeout,
		UpstreamMaxIdle:        input.UpstreamMaxIdle,
	}
	if err := tx.Create(loadBalance).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return serviceInfo, nil
}

// UpdateHTTPService 更新 HTTP 服务
func (s *ServiceService) UpdateHTTPService(input *dto.ServiceUpdateHTTPInput) error {
	// 验证 IP 列表与权重列表数量一致
	if len(strings.Split(input.IpList, ",")) != len(strings.Split(input.WeightList, ",")) {
		return errors.New("IP列表与权重列表数量不一致")
	}

	// 获取服务基本信息
	serviceInfo, err := models.GetServiceById(uint(input.ID))
	if err != nil {
		return errors.New("服务不存在")
	}

	serviceDetail, err := models.GetServiceDetailById(serviceInfo.ID)
	if err != nil {
		return errors.New("该服务详情不存在")
	}

	// 开始事务
	tx := models.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 更新服务信息
	info := serviceDetail.Info
	info.ServiceDesc = input.ServiceDesc
	if err := tx.Save(info).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 更新 HTTP 规则
	httpRule := serviceDetail.HTTPRule
	httpRule.NeedHttps = input.NeedHttps
	httpRule.NeedStripUri = input.NeedStripUri
	httpRule.NeedWebsocket = input.NeedWebsocket
	httpRule.UrlRewrite = input.UrlRewrite
	httpRule.HeaderTransfor = input.HeaderTransfor
	if err := tx.Save(httpRule).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 更新访问控制
	accessControl := serviceDetail.AccessControl
	accessControl.OpenAuth = input.OpenAuth
	accessControl.BlackList = input.BlackList
	accessControl.WhiteList = input.WhiteList
	accessControl.ClientIPFlowLimit = input.ClientipFlowLimit
	accessControl.ServiceFlowLimit = input.ServiceFlowLimit
	if err := tx.Save(accessControl).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 更新负载均衡
	loadbalance := serviceDetail.LoadBalance
	loadbalance.RoundType = input.RoundType
	loadbalance.IpList = input.IpList
	loadbalance.WeightList = input.WeightList
	loadbalance.UpstreamConnectTimeout = input.UpstreamConnectTimeout
	loadbalance.UpstreamHeaderTimeout = input.UpstreamHeaderTimeout
	loadbalance.UpstreamIdleTimeout = input.UpstreamIdleTimeout
	loadbalance.UpstreamMaxIdle = input.UpstreamMaxIdle
	if err := tx.Save(loadbalance).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}

// ServiceStat 获取服务统计
func (s *ServiceService) ServiceStat(input *dto.ServiceStatInput) (output *dto.ServiceStatOutput, error error) {

	// 获取服务基本信息
	serviceInfo, err := models.GetServiceById(uint(input.ID))
	if err != nil {
		return nil, errors.New("服务不存在")
	}

	// 获取服务详情信息
	_, err = models.GetServiceDetailById(serviceInfo.ID)
	if err != nil {
		return nil, errors.New("该服务详情不存在")
	}

	//今日每小时的统计
	todayList := []int64{}
	for i := 0; i < time.Now().Hour(); i++ {
		todayList = append(todayList, 0)
	}
	//昨日每小时的统计
	yesterdayList := []int64{}
	for i := 0; i < 23; i++ {
		yesterdayList = append(todayList, 0)
	}

	return &dto.ServiceStatOutput{
		Today:     todayList,
		Yesterday: yesterdayList,
	}, nil
}

// UpdateTcpService 更新 TCP 服务
func (s *ServiceService) UpdateTcpService(input *dto.ServiceUpdateTcpInput) error {
	//1.验证 IP 与权重数量一致
	ipList := strings.Split(input.IpList, ",")
	weightList := strings.Split(input.WeightList, ",")
	if len(ipList) != len(weightList) {
		return errors.New("IP列表与权重列表数量不一致")
	}

	//2. 获取服务基本信息
	serviceInfo, err := models.GetServiceById(uint(input.ID))
	if err != nil {
		return errors.New("服务不存在")
	}

	//3. 获取服务详情
	serviceDetail, err := models.GetServiceDetailById(serviceInfo.ID)
	if err != nil {
		return errors.New("该服务详情不存在")
	}

	//4. 检查端口是否变更以及是否被占用
	if serviceDetail.TCPRule != nil && serviceDetail.TCPRule.Port != input.Port {
		// 验证新端口是否被占用 (TCP)
		var tcpRuleCount int64
		err = models.DB.Model(&models.TcpRule{}).Where("port = ? AND service_id != ?", input.Port, input.ID).Count(&tcpRuleCount).Error
		if err != nil {
			return err
		}
		if tcpRuleCount > 0 {
			return errors.New("TCP端口已被占用")
		}

		// 验证新端口是否被占用 (GRPC)
		var grpcRuleCount int64
		err = models.DB.Model(&models.GrpcRule{}).Where("port = ? AND service_id != ?", input.Port, input.ID).Count(&grpcRuleCount).Error
		if err != nil {
			return err
		}
		if grpcRuleCount > 0 {
			return errors.New("端口已被GRPC服务占用")
		}
	}

	//----------开始事务----------------
	tx := models.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	//5. 更新服务信息serviceinfo(服务名和服务描述)
	info := serviceDetail.Info
	info.ServiceDesc = input.ServiceDesc
	if err := tx.Save(info).Error; err != nil {
		tx.Rollback()
		return err
	}

	//6. 更新 TCP 规则
	var tcpRule *models.TcpRule
	if serviceDetail.TCPRule != nil {
		tcpRule = serviceDetail.TCPRule
	}

	//7.更新service_tcp_rule(连接的serviceId和此tcp的port端口)
	info = serviceDetail.Info
	tcpRule.ServiceID = int64(info.ID)
	tcpRule.Port = input.Port
	if err := tx.Save(tcpRule).Error; err != nil {
		tx.Rollback()
		return err
	}

	//8. 更新访问控制
	var accessControl *models.AccessControl
	if serviceDetail.AccessControl != nil {
		accessControl = serviceDetail.AccessControl
	} else {
		accessControl = &models.AccessControl{}
	}
	accessControl.ServiceID = int64(info.ID)
	accessControl.OpenAuth = input.OpenAuth
	accessControl.BlackList = input.BlackList
	accessControl.WhiteList = input.WhiteList
	accessControl.WhiteHostName = input.WhiteHostName
	accessControl.ClientIPFlowLimit = input.ClientIPFlowLimit
	accessControl.ServiceFlowLimit = input.ServiceFlowLimit
	if err := tx.Save(accessControl).Error; err != nil {
		tx.Rollback()
		return err
	}

	//9. 更新负载均衡
	var loadBalance *models.LoadBalance
	if serviceDetail.LoadBalance != nil {
		loadBalance = serviceDetail.LoadBalance
	} else {
		loadBalance = &models.LoadBalance{}
	}
	loadBalance.ServiceID = int64(info.ID)
	loadBalance.RoundType = input.RoundType
	loadBalance.IpList = input.IpList
	loadBalance.WeightList = input.WeightList
	loadBalance.ForbidList = input.ForbidList
	if err := tx.Save(loadBalance).Error; err != nil {
		tx.Rollback()
		return err
	}

	//10. 提交事务
	return tx.Commit().Error
}

// CreateTcpService 创建 TCP 服务
func (s *ServiceService) CreateTcpService(d *dto.ServiceAddTcpInput) (*models.ServiceInfo, error) {
	//1.判断服务名是否存在
	exists, err := models.ServiceNameExists(d.ServiceName)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("服务名已存在")
	}

	//2.验证端口是否被占用 (TCP)
	var tcpRuleCount int64
	err = models.DB.Model(&models.TcpRule{}).Where("port = ?", d.Port).Count(&tcpRuleCount).Error
	if err != nil {
		return nil, err
	}
	if tcpRuleCount > 0 {
		return nil, errors.New("TCP端口已被占用")
	}

	//3.验证端口是否被占用 (GRPC)
	var grpcRuleCount int64
	err = models.DB.Model(&models.GrpcRule{}).Where("port = ?", d.Port).Count(&grpcRuleCount).Error
	if err != nil {
		return nil, err
	}
	if grpcRuleCount > 0 {
		return nil, errors.New("端口已被GRPC服务占用")
	}

	//4.验证 IP 与权重数量一致
	ipList := strings.Split(d.IpList, ",")
	weightList := strings.Split(d.WeightList, ",")
	if len(ipList) != len(weightList) {
		return nil, errors.New("IP列表与权重列表数量不一致")
	}

	//5. 开始事务
	tx := models.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	//6. 创建服务
	serviceInfo := &models.ServiceInfo{
		LoadType:    define.LoadTypeTCP,
		ServiceName: d.ServiceName,
		ServiceDesc: d.ServiceDesc,
	}
	if err := tx.Create(serviceInfo).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	//7. 创建 TCP 规则
	tcpRule := &models.TcpRule{
		ServiceID: int64(serviceInfo.ID),
		Port:      d.Port,
	}
	if err := tx.Create(tcpRule).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	//8. 创建访问控制
	accessControl := &models.AccessControl{
		ServiceID:         int64(serviceInfo.ID),
		OpenAuth:          d.OpenAuth,
		BlackList:         d.BlackList,
		WhiteList:         d.WhiteList,
		WhiteHostName:     d.WhiteHostName,
		ClientIPFlowLimit: d.ClientIPFlowLimit,
		ServiceFlowLimit:  d.ServiceFlowLimit,
	}
	if err := tx.Create(accessControl).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	//9. 创建负载均衡
	loadBalance := &models.LoadBalance{
		ServiceID:  int64(serviceInfo.ID),
		RoundType:  d.RoundType,
		IpList:     d.IpList,
		WeightList: d.WeightList,
		ForbidList: d.ForbidList,
	}
	if err := tx.Create(loadBalance).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	//10. 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return serviceInfo, nil
}

func (s *ServiceService) CreateGrpcService(d *dto.ServiceAddGrpcInput) (*models.ServiceInfo, error) {
	//1.判断服务名是否存在
	exists, err := models.ServiceNameExists(d.ServiceName)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("服务名已存在")
	}

	//2.验证端口是否被占用 (TCP)
	var tcpRuleCount int64
	err = models.DB.Model(&models.TcpRule{}).Where("port = ?", d.Port).Count(&tcpRuleCount).Error
	if err != nil {
		return nil, err
	}
	if tcpRuleCount > 0 {
		return nil, errors.New("TCP端口已被占用")
	}

	//3.验证端口是否被占用 (GRPC)
	var grpcRuleCount int64
	err = models.DB.Model(&models.GrpcRule{}).Where("port = ?", d.Port).Count(&grpcRuleCount).Error
	if err != nil {
		return nil, err
	}
	if grpcRuleCount > 0 {
		return nil, errors.New("端口已被GRPC服务占用")
	}

	//4.验证 IP 与权重数量一致
	ipList := strings.Split(d.IpList, ",")
	weightList := strings.Split(d.WeightList, ",")
	if len(ipList) != len(weightList) {
		return nil, errors.New("IP列表与权重列表数量不一致")
	}

	//5. 开始事务
	tx := models.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	//6. 先创建创建服务serviceinfo
	serviceInfo := &models.ServiceInfo{
		LoadType:    define.LoadTypeGRPC,
		ServiceName: d.ServiceName,
		ServiceDesc: d.ServiceDesc,
	}
	if err := tx.Create(serviceInfo).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	//7. 创建 load balance配置
	loadBalance := &models.LoadBalance{
		ServiceID:  int64(serviceInfo.ID),
		RoundType:  d.RoundType,
		IpList:     d.IpList,
		WeightList: d.WeightList,
		ForbidList: d.ForbidList,
	}
	if err := tx.Create(loadBalance).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	//8.创建grpc rule
	grpcRule := &models.GrpcRule{
		ServiceID:      int64(serviceInfo.ID),
		Port:           d.Port,
		HeaderTransfor: d.HeaderTransfor,
	}
	if err := tx.Create(grpcRule).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	//9.创建准入规则accessControl
	accessControl := &models.AccessControl{
		ServiceID:         int64(serviceInfo.ID),
		OpenAuth:          d.OpenAuth,
		BlackList:         d.BlackList,
		WhiteList:         d.WhiteList,
		WhiteHostName:     d.WhiteHostName,
		ClientIPFlowLimit: d.ClientIPFlowLimit,
		ServiceFlowLimit:  d.ServiceFlowLimit,
	}
	if err := tx.Create(accessControl).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	//10.提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return serviceInfo, nil
}

func (s *ServiceService) UpdateGrpcService(input *dto.ServiceUpdateGrpcInput) error {
	//1.验证 IP 与权重数量一致
	ipList := strings.Split(input.IpList, ",")
	weightList := strings.Split(input.WeightList, ",")
	if len(ipList) != len(weightList) {
		return errors.New("IP列表与权重列表数量不一致")
	}

	//2. 获取服务基本信息
	serviceInfo, err := models.GetServiceById(uint(input.ID))
	if err != nil {
		return errors.New("服务不存在")
	}

	//3. 获取服务详情
	serviceDetail, err := models.GetServiceDetailById(serviceInfo.ID)
	if err != nil {
		return errors.New("该服务详情不存在")
	}

	//----------开始事务----------------
	tx := models.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	//4. 更新服务信息serviceinfo(服务名和服务描述)
	info := serviceDetail.Info
	info.ServiceDesc = input.ServiceDesc
	if err := tx.Save(info).Error; err != nil {
		tx.Rollback()
		return err
	}

	//5.更新grp_rule
	var grpcRule *models.GrpcRule
	if serviceDetail.GRPCRule != nil {
		grpcRule = serviceDetail.GRPCRule
	} else {
		grpcRule = &models.GrpcRule{}
	}

	grpcRule.ServiceID = int64(info.ID)
	grpcRule.Port = input.Port
	grpcRule.HeaderTransfor = input.HeaderTransfor
	if err := tx.Save(grpcRule).Error; err != nil {
		tx.Rollback()
		return err
	}

	//6.更新访问控制access_control
	var accessControl *models.AccessControl
	if serviceDetail.AccessControl != nil {
		accessControl = serviceDetail.AccessControl
	}
	accessControl.ServiceID = int64(info.ID)
	accessControl.OpenAuth = input.OpenAuth
	accessControl.BlackList = input.BlackList
	accessControl.WhiteList = input.WhiteList
	accessControl.WhiteHostName = input.WhiteHostName
	accessControl.ClientIPFlowLimit = input.ClientIPFlowLimit
	accessControl.ServiceFlowLimit = input.ServiceFlowLimit
	if err := tx.Save(accessControl).Error; err != nil {
		tx.Rollback()
		return err
	}

	//7.更新负载均衡load_blance
	var loadBalance *models.LoadBalance
	if serviceDetail.LoadBalance != nil {
		loadBalance = serviceDetail.LoadBalance
	} else {
		loadBalance = &models.LoadBalance{}
	}
	loadBalance.ServiceID = int64(info.ID)
	loadBalance.RoundType = input.RoundType
	loadBalance.IpList = input.IpList
	loadBalance.WeightList = input.WeightList
	loadBalance.ForbidList = input.ForbidList
	if err := tx.Save(loadBalance).Error; err != nil {
		tx.Rollback()
		return err
	}

	//8.提交事务
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
