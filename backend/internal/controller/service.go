package controller

import (
	"fmt"
	"gateway-service/internal/define"
	"gateway-service/internal/dto"
	"gateway-service/internal/service"
	"gateway-service/pkg/response"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ServiceController 服务控制器
type ServiceController struct {
	serviceService *service.ServiceService
}

// NewServiceController 创建服务控制器
func NewServiceController() *ServiceController {
	return &ServiceController{
		serviceService: &service.ServiceService{},
	}
}

// ListServices 获取服务列表
func (c *ServiceController) ListServices(ctx *gin.Context) {
	var input dto.ServiceListInput
	if err := ctx.ShouldBindQuery(&input); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	results, pageNo, pageSize, total, err := c.serviceService.PageListServices(&input)
	if err != nil {
		response.InternalError(ctx, "Database error")
		return
	}

	// 转换成输出格式
	var ServiceListItemOutput []*dto.ServiceListItemOutput
	for _, listItem := range results {
		serviceDetail, err := c.serviceService.GetServiceDetail(int64(listItem.ID))
		if err != nil {
			response.BadRequest(ctx, "err in ServiceDetail found")
		}
		//1、http后缀接入 clusterIP+clusterPort+path
		//2、http域名接入 domain
		//3、tcp、grpc接入 clusterIP+servicePort
		serviceAddr := "unknow"
		clusterIP := os.Getenv("cluster_ip")
		clusterPort := os.Getenv("cluster_port")
		clusterSSLPort := os.Getenv("cluster_ssl_port")

		//
		if serviceDetail.Info.LoadType == define.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == define.HTTPRuleTypePrefixURL &&
			serviceDetail.HTTPRule.NeedHttps == 1 {
			serviceAddr = fmt.Sprintf("%s:%s%s", clusterIP, clusterSSLPort, serviceDetail.HTTPRule.Rule)
		}
		if serviceDetail.Info.LoadType == define.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == define.HTTPRuleTypePrefixURL &&
			serviceDetail.HTTPRule.NeedHttps == 0 {
			serviceAddr = fmt.Sprintf("%s:%s%s", clusterIP, clusterPort, serviceDetail.HTTPRule.Rule)
		}
		if serviceDetail.Info.LoadType == define.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == define.HTTPRuleTypeDomain {
			serviceAddr = serviceDetail.HTTPRule.Rule
		}
		if serviceDetail.Info.LoadType == define.LoadTypeTCP {
			serviceAddr = fmt.Sprintf("%s:%d", clusterIP, serviceDetail.TCPRule.Port)
		}
		if serviceDetail.Info.LoadType == define.LoadTypeGRPC {
			serviceAddr = fmt.Sprintf("%s:%d", clusterIP, serviceDetail.GRPCRule.Port)
		}
		ipList := serviceDetail.LoadBalance.GetIPListByModel()
		//流量计数器
		//counter, err := define.FlowCounterHandler.GetCounter(define.FlowServicePrefix + listItem.ServiceName)

		outItem := dto.ServiceListItemOutput{
			ID:          int64(listItem.ID),
			LoadType:    listItem.LoadType,
			ServiceName: listItem.ServiceName,
			ServiceDesc: listItem.ServiceDesc,
			ServiceAddr: serviceAddr,
			//Qps:         counter.QPS,
			//Qpd:         counter.TotalCount,
			TotalNode: len(ipList),
		}
		ServiceListItemOutput = append(ServiceListItemOutput, &outItem)
	}
	response.Success(ctx, gin.H{
		"list":  ServiceListItemOutput,
		"total": total,
		"page":  pageNo,
		"size":  pageSize,
	})
}

// ServiceDetail 获取服务详情
func (c *ServiceController) ServiceDetail(ctx *gin.Context) {

	var input dto.ServiceDeleteInput
	if err := ctx.ShouldBindQuery(&input); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	//service, err := models.GetServiceDetailById(uint(input.ID))
	detail, err := c.serviceService.GetServiceDetail(input.ID)
	if err != nil {
		response.NotFound(ctx, "Service not found")
		return
	}

	response.Success(ctx, detail)
}

// UpdateService 更新服务
func (c *ServiceController) UpdateService(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "Invalid service ID")
		return
	}

	var req struct {
		LoadType    int    `json:"load_type" binding:"required"`
		ServiceName string `json:"service_name" binding:"required"`
		ServiceDesc string `json:"service_desc"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	serviceinfo, err := c.serviceService.UpdateService(id, req.LoadType, req.ServiceName, req.ServiceDesc)
	if err != nil {
		response.NotFound(ctx, "Service not found")
		return
	}

	response.Success(ctx, serviceinfo)
}

// DeleteService 删除服务
func (c *ServiceController) DeleteService(ctx *gin.Context) {
	var input dto.ServiceDeleteInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	if err := c.serviceService.DeleteService(input.ID); err != nil {
		response.NotFound(ctx, "Service 删除失败")
		return
	}

	response.Success(ctx, gin.H{"message": "Service deleted successfully"})
}

// ServiceStat 获取服务统计
func (c *ServiceController) ServiceStat(ctx *gin.Context) {
	params := &dto.ServiceStatInput{}
	if err := ctx.ShouldBindQuery(&params); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	serviceStat, err := c.serviceService.ServiceStat(params)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Success(ctx, gin.H{
		"today":     serviceStat.Today,
		"yesterday": serviceStat.Yesterday,
	})
}

// CreateHTTP 创建 HTTP 服务
func (c *ServiceController) CreateHTTP(ctx *gin.Context) {
	var input dto.ServiceAddHTTPInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	serviceinfo, err := c.serviceService.CreateHTTPService(&input)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}

	response.Success(ctx, serviceinfo)
}

// UpdateHTTP 更新 HTTP 服务
func (c *ServiceController) UpdateHTTP(ctx *gin.Context) {
	var input dto.ServiceUpdateHTTPInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	err := c.serviceService.UpdateHTTPService(&input)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}

	response.Success(ctx, nil)
}

// CreateTcp 创建 TCP 服务
func (c *ServiceController) CreateTcp(ctx *gin.Context) {
	var input dto.ServiceAddTcpInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	_, err := c.serviceService.CreateTcpService(&input)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}

	response.Success(ctx, nil)
}

// UpdateTcp 更新 TCP 服务
func (c *ServiceController) UpdateTcp(ctx *gin.Context) {
	var input dto.ServiceUpdateTcpInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	err := c.serviceService.UpdateTcpService(&input)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}

	response.Success(ctx, nil)
}

// CreateGrpc 创建 gRPC 服务
func (c *ServiceController) CreateGrpc(ctx *gin.Context) {
	var input dto.ServiceAddGrpcInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	_, err := c.serviceService.CreateGrpcService(&input)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}

	response.Success(ctx, gin.H{"message": "gRPC service created"})
}

// UpdateGrpc 更新 gRPC 服务
func (c *ServiceController) UpdateGrpc(ctx *gin.Context) {
	var input dto.ServiceUpdateGrpcInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	err := c.serviceService.UpdateGrpcService(&input)
	if err != nil {
		response.Error(ctx, "grpc创建失败:"+err.Error())
		return
	}

	response.Success(ctx, gin.H{"message": "gRPC service updated"})
}
