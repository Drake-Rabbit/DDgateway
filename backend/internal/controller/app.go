package controller

import (
	"gateway-service/internal/dto"
	"gateway-service/internal/service"
	"gateway-service/pkg/response"

	"github.com/gin-gonic/gin"
)

// AppController 应用控制器
type AppController struct {
	appService *service.AppService
}

// NewAppController 创建应用控制器
func NewAppController() *AppController {
	return &AppController{
		appService: service.NewAppService(),
	}
}

// GetAppList 获取应用列表,带搜索条件
func (c *AppController) GetAppList(ctx *gin.Context) {
	var input dto.APPListInput
	if err := ctx.ShouldBindQuery(&input); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	// 设置默认分页参数
	if input.PageSize <= 0 {
		input.PageSize = 10
	}
	if input.PageNo <= 0 {
		input.PageNo = 1
	}

	list, total, err := c.appService.GetAppList(&input)
	if err != nil {
		response.InternalError(ctx, "Database error")
		return
	}

	// 转换成输出格式
	var outputItems []dto.APPListItemOutput
	for _, item := range list {
		outputItems = append(outputItems, dto.APPListItemOutput{
			ID:        int64(item.ID),
			AppID:     item.AppID,
			Name:      item.Name,
			Secret:    item.Secret,
			WhiteIPS:  item.WhiteIPS,
			Qpd:       item.Qpd,
			Qps:       item.Qps,
			RealQpd:   0, // 实际请求量需要从缓存或其他地方获取
			RealQps:   0, // 实时QPS需要从缓存或其他地方获取
			UpdatedAt: item.UpdatedAt,
			CreatedAt: item.CreatedAt,
		})
	}

	response.Success(ctx, gin.H{
		"list":  outputItems,
		"total": total,
		"page":  input.PageNo,
		"size":  input.PageSize,
	})
}

// GetAppDetail 获取应用详情
func (c *AppController) GetAppDetail(ctx *gin.Context) {
	var input dto.APPDetailInput
	if err := ctx.ShouldBindQuery(&input); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	app, err := c.appService.GetAppDetailByID(int(input.ID))
	if err != nil {
		response.NotFound(ctx, "App not found")
		return
	}

	response.Success(ctx, gin.H{
		"id":         app.ID,
		"app_id":     app.AppID,
		"name":       app.Name,
		"secret":     app.Secret,
		"white_ips":  app.WhiteIPS,
		"qpd":        app.Qpd,
		"qps":        app.Qps,
		"created_at": app.CreatedAt,
		"updated_at": app.UpdatedAt,
	})
}

// CreateApp 创建应用
func (c *AppController) CreateApp(ctx *gin.Context) {
	var input dto.APPAddHttpInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	//调用service的创建应用方法
	_, err := c.appService.CreateApp(&input)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}

	response.Success(ctx, gin.H{
		"message": "App created successfully",
	})
}

// UpdateApp 更新应用
func (c *AppController) UpdateApp(ctx *gin.Context) {
	var input dto.APPUpdateHttpInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	if input.AppID == "" {
		response.BadRequest(ctx, "App ID is required")
		return
	}

	//调用service的更新应用方法
	err := c.appService.UpdateApp(&input)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}

	response.Success(ctx, gin.H{
		"message": "App updated successfully",
	})
}

// DeleteApp 删除应用
func (c *AppController) DeleteApp(ctx *gin.Context) {
	var input dto.APPDeleteInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	err := c.appService.DeleteApp(int(input.ID))
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}

	response.Success(ctx, gin.H{
		"message": "App deleted successfully",
	})
}

// GetAppStats 获取应用统计信息
func (c *AppController) GetAppStats(ctx *gin.Context) {
	stats, err := c.appService.GetAppStats()
	if err != nil {
		response.InternalError(ctx, "Failed to get app statistics")
		return
	}

	response.Success(ctx, stats)
}
