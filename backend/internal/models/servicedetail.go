package models

import (
	"gateway-service/internal/dto"
	"sync"
)

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

// ServiceManagerHandler 服务管理实例
var serviceManagerHandler *ServiceManager

func init() {
	serviceManagerHandler = NewServiceManager()
}

// ServiceManager 服务管理结构体
type ServiceManager struct {
	ServiceMap   map[string]*ServiceDetail
	ServiceSlice []*ServiceDetail
	Locker       sync.RWMutex
	init         sync.Once
	err          error
}

func NewServiceManager() *ServiceManager {
	sm := &ServiceManager{
		ServiceMap:   make(map[string]*ServiceDetail),
		ServiceSlice: []*ServiceDetail{},
		Locker:       sync.RWMutex{},
		init:         sync.Once{},
	}
	return sm
}

// 一次性加载服务配置信息,放在内存中
func (sm *ServiceManager) LoadOnce() error {

	sm.init.Do(func() {

		//1.db中获取service配置信息
		params := &dto.ServiceListInput{PageNo: 1, PageSize: 99999}
		list, _, err := GetServicePage(params)
		if err != nil {
			sm.err = err
			return
		}
		sm.Locker.Lock()
		defer sm.Locker.Unlock()
		//2.遍历service,配置详情信息存储在map中
		for _, listItem := range list {
			serviceDetail, err := GetServiceDetailById(listItem.ID)
			//fmt.Println("serviceDetail")
			//fmt.Println(public.Obj2Json(serviceDetail))
			if err != nil {
				sm.err = err
				return
			}
			sm.ServiceMap[listItem.ServiceName] = serviceDetail
			sm.ServiceSlice = append(sm.ServiceSlice, serviceDetail)
		}
	})

	return sm.err
}
