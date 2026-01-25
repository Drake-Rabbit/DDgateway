package controller

import (
	"gateway-service/internal/service"
	"gateway-service/pkg/response"

	"github.com/gin-gonic/gin"
)

// TenantController 租户控制器
type TenantController struct {
	tenantService *service.TenantService
}

// NewTenantController 创建租户控制器
func NewTenantController() *TenantController {
	return &TenantController{
		tenantService: &service.TenantService{},
	}
}

// ListTenants 获取所有租户
func (c *TenantController) ListTenants(ctx *gin.Context) {
	tenants, err := c.tenantService.ListTenants()
	if err != nil {
		response.InternalError(ctx, "Failed to fetch tenants")
		return
	}

	response.Success(ctx, tenants)
}

// GetTenant 获取租户详情
func (c *TenantController) GetTenant(ctx *gin.Context) {
	id := ctx.Param("id")

	tenant, err := c.tenantService.GetTenant(id)
	if err != nil {
		response.NotFound(ctx, "Tenant not found")
		return
	}

	response.Success(ctx, tenant)
}

// CreateTenant 创建租户
func (c *TenantController) CreateTenant(ctx *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Code        string `json:"code" binding:"required"`
		Description string `json:"description"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	tenant, err := c.tenantService.CreateTenant(req.Name, req.Code, req.Description)
	if err != nil {
		response.InternalError(ctx, "Failed to create tenant")
		return
	}

	response.Success(ctx, tenant)
}

// UpdateTenant 更新租户
func (c *TenantController) UpdateTenant(ctx *gin.Context) {
	id := ctx.Param("id")

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	tenant, err := c.tenantService.UpdateTenant(id, req.Name, req.Description, req.Status)
	if err != nil {
		response.NotFound(ctx, "Tenant not found")
		return
	}

	response.Success(ctx, tenant)
}

// DeleteTenant 删除租户
func (c *TenantController) DeleteTenant(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.tenantService.DeleteTenant(id)
	if err != nil {
		response.NotFound(ctx, "Tenant not found")
		return
	}

	response.Success(ctx, gin.H{"message": "Tenant deleted successfully"})
}
