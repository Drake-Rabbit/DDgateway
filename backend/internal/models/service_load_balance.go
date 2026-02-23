package models

import (
	"fmt"
	"gateway-service/internal/define"
	"gateway-service/internal/reverse_proxy/load_balance"
	"net/http"
	"strings"
	"sync"
	"time"
)

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
	UpstreamConnectTimeout int    `gorm:"column:upstream_connect_timeout" json:"upstream_connect_timeout" description:""`
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

func (b *LoadBalance) GetWeightListByModel() []string {
	return strings.Split(b.WeightList, ",")
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

// 内存中有服务管理器-那也对应的有一个负载均衡器的管理器
type LoadBalanceManager struct {
	LoadBalanceMap    map[string]*LoadBalancerItem
	LoadBanlanceSlice []*LoadBalancerItem
	Locker            sync.RWMutex
}

type LoadBalancerItem struct {
	LoadBanlance load_balance.LoadBalance
	ServiceName  string
}

func NewLoadBalancer() *LoadBalanceManager {
	return &LoadBalanceManager{
		LoadBalanceMap:    map[string]*LoadBalancerItem{},
		LoadBanlanceSlice: []*LoadBalancerItem{},
		Locker:            sync.RWMutex{},
	}
}

var LoadBalancerHandler *LoadBalanceManager

func init() {
	LoadBalancerHandler = NewLoadBalancer()
}

// 得到对应的负载均衡策略
func (l *LoadBalanceManager) GetLoadBalance(service *ServiceDetail) (load_balance.LoadBalance, error) {
	//如果内存中有,则直接返回
	for _, lbrItem := range l.LoadBanlanceSlice {
		if lbrItem.ServiceName == service.Info.ServiceName {
			return lbrItem.LoadBanlance, nil
		}
	}

	// 判断是否是http
	schema := "http://"
	if service.HTTPRule.NeedHttps == 1 {
		schema = "https://"
	}
	// tcp和grpc
	if service.Info.LoadType == define.LoadTypeTCP || service.Info.LoadType == define.LoadTypeGRPC {
		schema = ""
	}
	// 得到ip列表
	ipList := service.LoadBalance.GetIPListByModel()
	// 得到权重列表
	weightList := service.LoadBalance.GetWeightListByModel()
	// 得到ip配置
	ipConf := map[string]string{}
	for ipIndex, ipItem := range ipList {
		ipConf[ipItem] = weightList[ipIndex]
	}
	//负载均衡检查机制启动
	mConf, err := load_balance.NewLoadBalanceCheckConf(fmt.Sprintf("%s%s", schema, "%s"), ipConf)
	if err != nil {
		return nil, err
	}
	// 得到负载均衡
	lb := load_balance.LoadBanlanceFactorWithConf(load_balance.LbType(service.LoadBalance.RoundType), mConf)

	// 添加到内存中
	lbItem := &LoadBalancerItem{
		LoadBanlance: lb,
		ServiceName:  service.Info.ServiceName,
	}
	l.LoadBanlanceSlice = append(l.LoadBanlanceSlice, lbItem)

	l.Locker.Lock()
	defer l.Locker.Unlock()
	l.LoadBalanceMap[service.Info.ServiceName] = lbItem
	return lb, nil
}

var TransportorHandler *Transportor

// 连接池transport 每个服务 一个连接池
type Transportor struct {
	TransportMap   map[string]*TransportItem
	TransportSlice []*TransportItem
	Locker         sync.RWMutex
}

type TransportItem struct {
	Trans       *http.Transport
	ServiceName string
}

func init() {
	TransportorHandler = NewTransportor()
}

// 创建连接池
func NewTransportor() *Transportor {
	return &Transportor{
		TransportMap:   map[string]*TransportItem{},
		TransportSlice: []*TransportItem{},
		Locker:         sync.RWMutex{},
	}
}

func (t *Transportor) GetTrans(service *ServiceDetail) (*http.Transport, error) {
	//for循环找到服务对应的连接池
	for _, item := range t.TransportSlice {
		if item.ServiceName == service.Info.ServiceName {
			return item.Trans, nil
		}
	}
	//如果没有找到,则新建连接池到内存
	trans := &http.Transport{
		MaxIdleConns:          service.LoadBalance.UpstreamMaxIdle,
		MaxIdleConnsPerHost:   service.LoadBalance.UpstreamMaxIdle,
		IdleConnTimeout:       time.Duration(service.LoadBalance.UpstreamIdleTimeout) * time.Second,
		ResponseHeaderTimeout: time.Duration(service.LoadBalance.UpstreamHeaderTimeout) * time.Second,
		//TLSHandshakeTimeout:   10 * time.Second,
		//ExpectContinueTimeout: 1 * time.Second,

	}
	//添加到内存中
	//save to map and slice
	transItem := &TransportItem{
		Trans:       trans,
		ServiceName: service.Info.ServiceName,
	}
	t.TransportSlice = append(t.TransportSlice, transItem)
	t.Locker.Lock()
	defer t.Locker.Unlock()
	t.TransportMap[service.Info.ServiceName] = transItem

	return trans, nil
}
