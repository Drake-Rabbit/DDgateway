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

// GetAppList 获取应用列表
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
	appID := ctx.Query("app_id")
	if appID == "" {
		response.BadRequest(ctx, "App ID is required")
		return
	}

	app, err := c.appService.GetApp(appID)
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

	if input.AppID == "" || input.Name == "" {
		response.BadRequest(ctx, "App ID and Name are required")
		return
	}

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
	appID := ctx.Query("app_id")
	if appID == "" {
		response.BadRequest(ctx, "App ID is required")
		return
	}

	err := c.appService.DeleteApp(appID)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}

	response.Success(ctx, gin.H{
		"message": "App deleted successfully",
	})
}

// SearchApp 搜索应用
func (c *AppController) SearchApp(ctx *gin.Context) {
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

	list, total, err := c.appService.SearchApp(&input)
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
			RealQpd:   0,
			RealQps:   0,
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

// BatchCreateApp 批量创建应用
func (c *AppController) BatchCreateApp(ctx *gin.Context) {
	var inputs []*dto.APPAddHttpInput
	if err := ctx.ShouldBindJSON(&inputs); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	if len(inputs) == 0 {
		response.BadRequest(ctx, "No apps to create")
		return
	}

	apps, err := c.appService.BatchCreateApp(inputs)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}

	response.Success(ctx, gin.H{
		"created_count": len(apps),
		"apps":          apps,
		"message":       "Apps created successfully",
	})
}

// BatchUpdateApp 批量更新应用
func (c *AppController) BatchUpdateApp(ctx *gin.Context) {
	var request struct {
		AppIDs  []string               `json:"app_ids" binding:"required"`
		Updates map[string]interface{} `json:"updates" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	if len(request.AppIDs) == 0 {
		response.BadRequest(ctx, "App IDs are required")
		return
	}

	err := c.appService.BatchUpdateApp(request.AppIDs, request.Updates)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}

	response.Success(ctx, gin.H{
		"updated_count": len(request.AppIDs),
		"message":       "Apps updated successfully",
	})
}

// BatchDeleteApp 批量删除应用
func (c *AppController) BatchDeleteApp(ctx *gin.Context) {
	var request struct {
		AppIDs []string `json:"app_ids" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	if len(request.AppIDs) == 0 {
		response.BadRequest(ctx, "App IDs are required")
		return
	}

	err := c.appService.BatchDeleteApp(request.AppIDs)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}

	response.Success(ctx, gin.H{
		"deleted_count": len(request.AppIDs),
		"message":       "Apps deleted successfully",
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
