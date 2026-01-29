package controller

import (
	"gateway-service/internal/define"
	"gateway-service/internal/dto"
	"gateway-service/internal/service"
	"gateway-service/pkg/response"
	"github.com/gin-gonic/gin"
	"time"
)

// DashboradController 应用控制器
type DashboradController struct {
	DashboradService *service.DashboradService
	serviceService   *service.ServiceService
	appService       *service.AppService
}

// NewDashboradController 创建应用控制器
func NewDashboradController() *DashboradController {
	return &DashboradController{
		DashboradService: &service.DashboradService{},
		serviceService:   &service.ServiceService{},
		appService:       &service.AppService{},
	}
}

// 各个服务总数统计
func (c DashboradController) PanelGruopData(context *gin.Context) {

	var serviceService *service.ServiceService
	_, serviceNum, err := serviceService.PageListServices(&dto.ServiceListInput{PageNo: 1, PageSize: 1})
	if err != nil {
		response.Error(context, "获取service服务总数统计错误")
	}
	_, appNum, err := c.appService.GetAppList(&dto.APPListInput{PageNo: 1, PageSize: 1})
	if err != nil {
		response.Error(context, "获取app服务总数统计错误")
	}
	//中间件取数据
	//counter, err := public.FlowCounterHandler.GetCounter(public.FlowTotal)
	//if err != nil {
	//	middleware.ResponseError(c, 2003, err)
	//	return
	//}

	out := &dto.PanelGroupDataOutput{
		ServiceNum:      serviceNum,
		AppNum:          appNum,
		TodayRequestNum: 0, //counter.TotalCount,
		CurrentQPS:      0, //counter.QPS,
	}

	response.Success(context, out)
}

func (c DashboradController) FlowStat(context *gin.Context) {
	todayList := []int64{}
	currentTime := time.Now()
	for i := 0; i <= currentTime.Hour(); i++ {
		//dateTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), i, 0, 0, 0, lib.TimeLocation)
		//hourData, _ := counter.GetHourData(dateTime)
		todayList = append(todayList, 0)
	}

	yesterdayList := []int64{}
	//yesterTime := currentTime.Add(-1 * time.Duration(time.Hour*24))
	for i := 0; i <= 23; i++ {
		//dateTime := time.Date(yesterTime.Year(), yesterTime.Month(), yesterTime.Day(), i, 0, 0, 0, lib.TimeLocation)
		//hourData, _ := counter.GetHourData(dateTime)
		yesterdayList = append(yesterdayList, 0)
	}
	response.Success(context, &dto.ServiceStatOutput{
		Today:     todayList,
		Yesterday: yesterdayList,
	})
}

func (c DashboradController) ServiceStat(context *gin.Context) {

	//获取仪表盘各个服务type分类总数
	stat_list, err := c.serviceService.GroupByLoadType()
	if err != nil {
		response.Error(context, "获取服务分类统计错误")
	}

	//1.
	legend := []string{}
	for index, item := range stat_list {
		name, ok := define.LoadTypeMap[item.LoadType]
		if !ok {
			response.Error(context, "获取服务分类统计错误:load_tpye 不存在")
		}
		stat_list[index].Name = name
		legend = append(legend, name)
	}

	out := &dto.DashServiceStatOutput{
		Legend: legend,
		Data:   stat_list,
	}

	response.Success(context, out)
}
